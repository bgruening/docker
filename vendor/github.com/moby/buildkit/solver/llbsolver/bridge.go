package llbsolver

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/containerd/platforms"
	"github.com/mitchellh/hashstructure/v2"
	"github.com/moby/buildkit/cache/remotecache"
	"github.com/moby/buildkit/client"
	"github.com/moby/buildkit/client/llb/sourceresolver"
	"github.com/moby/buildkit/executor"
	resourcestypes "github.com/moby/buildkit/executor/resources/types"
	"github.com/moby/buildkit/frontend"
	gw "github.com/moby/buildkit/frontend/gateway/client"
	"github.com/moby/buildkit/identity"
	"github.com/moby/buildkit/session"
	"github.com/moby/buildkit/solver"
	"github.com/moby/buildkit/solver/errdefs"
	llberrdefs "github.com/moby/buildkit/solver/llbsolver/errdefs"
	"github.com/moby/buildkit/solver/llbsolver/provenance"
	"github.com/moby/buildkit/solver/pb"
	"github.com/moby/buildkit/sourcepolicy"
	spb "github.com/moby/buildkit/sourcepolicy/pb"
	"github.com/moby/buildkit/util/bklog"
	"github.com/moby/buildkit/util/entitlements"
	"github.com/moby/buildkit/util/flightcontrol"
	"github.com/moby/buildkit/util/progress"
	"github.com/moby/buildkit/worker"
	digest "github.com/opencontainers/go-digest"
	"github.com/pkg/errors"
)

type llbBridge struct {
	builder                   solver.Builder
	frontends                 map[string]frontend.Frontend
	resolveWorker             func() (worker.Worker, error)
	eachWorker                func(func(worker.Worker) error) error
	resolveCacheImporterFuncs map[string]remotecache.ResolveCacheImporterFunc
	cms                       map[string]solver.CacheManager
	cmsMu                     sync.Mutex
	sm                        *session.Manager

	executorOnce sync.Once
	executorErr  error
	executor     executor.Executor
}

func (b *llbBridge) Warn(ctx context.Context, dgst digest.Digest, msg string, opts frontend.WarnOpts) error {
	return b.builder.InContext(ctx, func(ctx context.Context, g session.Group) error {
		pw, ok, _ := progress.NewFromContext(ctx, progress.WithMetadata("vertex", dgst))
		if !ok {
			return nil
		}
		level := opts.Level
		if level == 0 {
			level = 1
		}
		pw.Write(identity.NewID(), client.VertexWarning{
			Vertex:     dgst,
			Level:      level,
			Short:      []byte(msg),
			SourceInfo: opts.SourceInfo,
			Range:      opts.Range,
			Detail:     opts.Detail,
			URL:        opts.URL,
		})
		return pw.Close()
	})
}

func (b *llbBridge) loadResult(ctx context.Context, def *pb.Definition, cacheImports []gw.CacheOptionsEntry, pol []*spb.Policy) (solver.CachedResultWithProvenance, error) {
	w, err := b.resolveWorker()
	if err != nil {
		return nil, err
	}
	ent, err := loadEntitlements(b.builder)
	if err != nil {
		return nil, err
	}
	srcPol, err := loadSourcePolicy(b.builder)
	if err != nil {
		return nil, err
	}
	var polEngine SourcePolicyEvaluator
	if srcPol != nil || len(pol) > 0 {
		for _, p := range pol {
			if p == nil {
				return nil, errors.Errorf("invalid nil policy")
			}
			if err := validateSourcePolicy(p); err != nil {
				return nil, err
			}
		}
		if srcPol != nil {
			pol = append([]*spb.Policy{srcPol}, pol...)
		}
		polEngine = sourcepolicy.NewEngine(pol)
	}
	var cms []solver.CacheManager
	for _, im := range cacheImports {
		cmID, err := cmKey(im)
		if err != nil {
			return nil, err
		}
		b.cmsMu.Lock()
		var cm solver.CacheManager
		if prevCm, ok := b.cms[cmID]; !ok {
			func(cmID string, im gw.CacheOptionsEntry) {
				cm = newLazyCacheManager(cmID, func() (solver.CacheManager, error) {
					var cmNew solver.CacheManager
					if err := inBuilderContext(context.TODO(), b.builder, "importing cache manifest from "+cmID, "", func(ctx context.Context, g session.Group) error {
						resolveCI, ok := b.resolveCacheImporterFuncs[im.Type]
						if !ok {
							return errors.Errorf("unknown cache importer: %s", im.Type)
						}
						ci, desc, err := resolveCI(ctx, g, im.Attrs)
						if err != nil {
							return errors.Wrapf(err, "failed to configure %v cache importer", im.Type)
						}
						cmNew, err = ci.Resolve(ctx, desc, cmID, w)
						return err
					}); err != nil {
						bklog.G(ctx).Debugf("error while importing cache manifest from cmId=%s: %v", cmID, err)
						return nil, err
					}
					return cmNew, nil
				})
			}(cmID, im)
			b.cms[cmID] = cm
		} else {
			cm = prevCm
		}
		cms = append(cms, cm)
		b.cmsMu.Unlock()
	}
	dpc := &detectPrunedCacheID{}

	edge, err := Load(ctx, def, polEngine, dpc.Load, ValidateEntitlements(ent, w.CDIManager()), WithCacheSources(cms), NormalizeRuntimePlatforms(), WithValidateCaps())
	if err != nil {
		return nil, errors.Wrap(err, "failed to load LLB")
	}

	if len(dpc.ids) > 0 {
		if err := b.eachWorker(func(w worker.Worker) error {
			return w.PruneCacheMounts(ctx, dpc.ids)
		}); err != nil {
			return nil, err
		}
	}

	res, err := b.builder.Build(ctx, edge)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (b *llbBridge) validateEntitlements(p executor.ProcessInfo) error {
	ent, err := loadEntitlements(b.builder)
	if err != nil {
		return err
	}
	v := entitlements.Values{
		NetworkHost:      p.Meta.NetMode == pb.NetMode_HOST,
		SecurityInsecure: p.Meta.SecurityMode == pb.SecurityMode_INSECURE,
	}
	return ent.Check(v)
}

func (b *llbBridge) Run(ctx context.Context, id string, rootfs executor.Mount, mounts []executor.Mount, process executor.ProcessInfo, started chan<- struct{}) (resourcestypes.Recorder, error) {
	if err := b.validateEntitlements(process); err != nil {
		return nil, err
	}

	if err := b.loadExecutor(); err != nil {
		return nil, err
	}
	return b.executor.Run(ctx, id, rootfs, mounts, process, started)
}

func (b *llbBridge) Exec(ctx context.Context, id string, process executor.ProcessInfo) error {
	if err := b.validateEntitlements(process); err != nil {
		return err
	}

	if err := b.loadExecutor(); err != nil {
		return err
	}
	return b.executor.Exec(ctx, id, process)
}

func (b *llbBridge) loadExecutor() error {
	b.executorOnce.Do(func() {
		w, err := b.resolveWorker()
		if err != nil {
			b.executorErr = err
			return
		}
		b.executor = w.Executor()
	})
	return b.executorErr
}

type resultProxy struct {
	id         string
	b          *provenanceBridge
	req        frontend.SolveRequest
	g          flightcontrol.Group[solver.CachedResult]
	mu         sync.Mutex
	released   bool
	v          solver.CachedResult
	err        error
	errResults []solver.Result
	provenance *provenance.Capture
}

func newResultProxy(b *provenanceBridge, req frontend.SolveRequest) *resultProxy {
	return &resultProxy{req: req, b: b, id: identity.NewID()}
}

func (rp *resultProxy) ID() string {
	return rp.id
}

func (rp *resultProxy) Definition() *pb.Definition {
	return rp.req.Definition
}

func (rp *resultProxy) Provenance() any {
	if rp.provenance == nil {
		return nil
	}
	return rp.provenance
}

func (rp *resultProxy) Release(ctx context.Context) (err error) {
	rp.mu.Lock()
	defer rp.mu.Unlock()
	for _, res := range rp.errResults {
		rerr := res.Release(ctx)
		if rerr != nil {
			err = rerr
		}
	}
	if rp.v != nil {
		if rp.released {
			bklog.G(ctx).Warnf("release of already released result")
		}
		rerr := rp.v.Release(ctx)
		if err != nil {
			return rerr
		}
	}
	rp.released = true
	return
}

func (rp *resultProxy) wrapError(err error) error {
	if err == nil {
		return nil
	}
	var ve *errdefs.VertexError
	if errors.As(err, &ve) {
		if rp.req.Definition.Source != nil {
			locs, ok := rp.req.Definition.Source.Locations[ve.Digest]
			if ok {
				for _, loc := range locs.Locations {
					err = errdefs.WithSource(err, &errdefs.Source{
						Info:   rp.req.Definition.Source.Infos[loc.SourceIndex],
						Ranges: loc.Ranges,
					})
				}
			}
		}
	}
	return err
}

func (rp *resultProxy) loadResult(ctx context.Context) (solver.CachedResultWithProvenance, error) {
	res, err := rp.b.loadResult(ctx, rp.req.Definition, rp.req.CacheImports, rp.req.SourcePolicies)
	var ee *llberrdefs.ExecError
	if errors.As(err, &ee) {
		ee.EachRef(func(res solver.Result) error {
			rp.errResults = append(rp.errResults, res)
			return nil
		})
		// acquire ownership so ExecError finalizer doesn't attempt to release as well
		ee.OwnerBorrowed = true
	}
	return res, err
}

func (rp *resultProxy) Result(ctx context.Context) (res solver.CachedResult, err error) {
	defer func() {
		err = rp.wrapError(err)
	}()
	return rp.g.Do(ctx, "result", func(ctx context.Context) (solver.CachedResult, error) {
		rp.mu.Lock()
		if rp.released {
			rp.mu.Unlock()
			return nil, errors.Errorf("accessing released result")
		}
		if rp.v != nil || rp.err != nil {
			rp.mu.Unlock()
			return rp.v, rp.err
		}
		rp.mu.Unlock()
		v, err := rp.loadResult(ctx)
		if err != nil {
			select {
			case <-ctx.Done():
				if errdefs.IsCanceled(ctx, err) {
					return v, err
				}
			default:
			}
		}
		rp.mu.Lock()
		if rp.released {
			if v != nil {
				v.Release(context.TODO())
			}
			rp.mu.Unlock()
			return nil, errors.Errorf("evaluating released result")
		}
		if err == nil {
			var capture *provenance.Capture
			capture, err = captureProvenance(ctx, v)
			if err != nil {
				err = errors.Errorf("failed to capture provenance: %v", err)
				v.Release(context.TODO())
				v = nil
			}
			rp.provenance = capture
		}
		rp.v = v
		rp.err = err
		rp.mu.Unlock()
		return v, err
	})
}

func (b *llbBridge) ResolveSourceMetadata(ctx context.Context, op *pb.SourceOp, opt sourceresolver.Opt) (resp *sourceresolver.MetaResponse, err error) {
	w, err := b.resolveWorker()
	if err != nil {
		return nil, err
	}
	if opt.LogName == "" {
		// TODO: better name
		opt.LogName = fmt.Sprintf("resolve image config for %s", op.Identifier)
	}
	id := op.Identifier
	if opt.Platform != nil {
		id += platforms.FormatAll(*opt.Platform)
	} else {
		id += platforms.FormatAll(platforms.DefaultSpec())
	}
	pol, err := loadSourcePolicy(b.builder)
	if err != nil {
		return nil, err
	}
	if pol != nil {
		opt.SourcePolicies = append(opt.SourcePolicies, pol)
	}

	if _, err := sourcepolicy.NewEngine(opt.SourcePolicies).Evaluate(ctx, op); err != nil {
		return nil, errors.Wrap(err, "could not resolve image due to policy")
	}

	// policy is evaluated, so we can remove it from the options
	opt.SourcePolicies = nil

	err = inBuilderContext(ctx, b.builder, opt.LogName, id, func(ctx context.Context, g session.Group) error {
		resp, err = w.ResolveSourceMetadata(ctx, op, opt, b.sm, g)
		return err
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

type lazyCacheManager struct {
	id   string
	main solver.CacheManager

	waitCh chan struct{}
	err    error
}

func (lcm *lazyCacheManager) ID() string {
	return lcm.id
}

func (lcm *lazyCacheManager) Query(inp []solver.CacheKeyWithSelector, inputIndex solver.Index, dgst digest.Digest, outputIndex solver.Index) ([]*solver.CacheKey, error) {
	lcm.wait()
	if lcm.main == nil {
		return nil, nil
	}
	return lcm.main.Query(inp, inputIndex, dgst, outputIndex)
}

func (lcm *lazyCacheManager) Records(ctx context.Context, ck *solver.CacheKey) ([]*solver.CacheRecord, error) {
	lcm.wait()
	if lcm.main == nil {
		return nil, nil
	}
	return lcm.main.Records(ctx, ck)
}

func (lcm *lazyCacheManager) Load(ctx context.Context, rec *solver.CacheRecord) (solver.Result, error) {
	if err := lcm.wait(); err != nil {
		return nil, err
	}
	return lcm.main.Load(ctx, rec)
}

func (lcm *lazyCacheManager) Save(key *solver.CacheKey, s solver.Result, createdAt time.Time) (*solver.ExportableCacheKey, error) {
	if err := lcm.wait(); err != nil {
		return nil, err
	}
	return lcm.main.Save(key, s, createdAt)
}

func (lcm *lazyCacheManager) ReleaseUnreferenced(ctx context.Context) error {
	if err := lcm.wait(); err != nil {
		return err
	}
	return lcm.main.ReleaseUnreferenced(ctx)
}

func (lcm *lazyCacheManager) wait() error {
	<-lcm.waitCh
	return lcm.err
}

func newLazyCacheManager(id string, fn func() (solver.CacheManager, error)) solver.CacheManager {
	lcm := &lazyCacheManager{id: id, waitCh: make(chan struct{})}
	go func() {
		defer close(lcm.waitCh)
		cm, err := fn()
		if err != nil {
			lcm.err = err
			return
		}
		lcm.main = cm
	}()
	return lcm
}

func cmKey(im gw.CacheOptionsEntry) (string, error) {
	if im.Type == "registry" && im.Attrs["ref"] != "" {
		return im.Attrs["ref"], nil
	}
	i, err := hashstructure.Hash(im, hashstructure.FormatV2, nil)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s:%d", im.Type, i), nil
}
