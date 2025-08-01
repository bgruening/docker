package daemon

import (
	"context"

	cerrdefs "github.com/containerd/errdefs"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/events"
	"github.com/moby/moby/v2/errdefs"
	"github.com/pkg/errors"
)

// ContainerUpdate updates configuration of the container
func (daemon *Daemon) ContainerUpdate(name string, hostConfig *container.HostConfig) (container.UpdateResponse, error) {
	var warnings []string

	daemonCfg := daemon.config()
	warnings, err := daemon.verifyContainerSettings(daemonCfg, hostConfig, nil, true)
	if err != nil {
		return container.UpdateResponse{Warnings: warnings}, errdefs.InvalidParameter(err)
	}

	if err := daemon.update(name, hostConfig); err != nil {
		return container.UpdateResponse{Warnings: warnings}, err
	}

	return container.UpdateResponse{Warnings: warnings}, nil
}

func (daemon *Daemon) update(name string, hostConfig *container.HostConfig) error {
	if hostConfig == nil {
		return nil
	}

	ctr, err := daemon.GetContainer(name)
	if err != nil {
		return err
	}

	restoreConfig := false
	backupHostConfig := *ctr.HostConfig

	defer func() {
		if restoreConfig {
			ctr.Lock()
			if !ctr.RemovalInProgress && !ctr.Dead {
				ctr.HostConfig = &backupHostConfig
				ctr.CheckpointTo(context.WithoutCancel(context.TODO()), daemon.containersReplica)
			}
			ctr.Unlock()
		}
	}()

	ctr.Lock()

	if ctr.RemovalInProgress || ctr.Dead {
		ctr.Unlock()
		return errCannotUpdate(ctr.ID, errors.New(`container is marked for removal and cannot be "update"`))
	}

	if err := ctr.UpdateContainer(hostConfig); err != nil {
		restoreConfig = true
		ctr.Unlock()
		return errCannotUpdate(ctr.ID, err)
	}
	if err := ctr.CheckpointTo(context.TODO(), daemon.containersReplica); err != nil {
		restoreConfig = true
		ctr.Unlock()
		return errCannotUpdate(ctr.ID, err)
	}

	ctr.Unlock()

	// if Restart Policy changed, we need to update container monitor
	if hostConfig.RestartPolicy.Name != "" {
		ctr.UpdateMonitor(hostConfig.RestartPolicy)
	}

	defer daemon.LogContainerEvent(ctr, events.ActionUpdate)

	// If container is not running, update hostConfig struct is enough,
	// resources will be updated when the container is started again.
	// If container is running (including paused), we need to update configs
	// to the real world.
	ctr.Lock()
	isRestarting := ctr.Restarting
	tsk, err := ctr.GetRunningTask()
	ctr.Unlock()
	if cerrdefs.IsConflict(err) || isRestarting {
		return nil
	}
	if err != nil {
		return err
	}

	if err := tsk.UpdateResources(context.TODO(), toContainerdResources(hostConfig.Resources)); err != nil {
		restoreConfig = true
		// TODO: it would be nice if containerd responded with better errors here so we can classify this better.
		return errCannotUpdate(ctr.ID, errdefs.System(err))
	}

	return nil
}

func errCannotUpdate(containerID string, err error) error {
	return errors.Wrap(err, "Cannot update container "+containerID)
}
