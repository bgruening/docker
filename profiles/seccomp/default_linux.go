package seccomp

import (
	"github.com/opencontainers/runtime-spec/specs-go"
	"golang.org/x/sys/unix"
)

func arches() []Architecture {
	return []Architecture{
		{
			Arch:      specs.ArchX86_64,
			SubArches: []specs.Arch{specs.ArchX86, specs.ArchX32},
		},
		{
			Arch:      specs.ArchAARCH64,
			SubArches: []specs.Arch{specs.ArchARM},
		},
		{
			Arch:      specs.ArchMIPS64,
			SubArches: []specs.Arch{specs.ArchMIPS, specs.ArchMIPS64N32},
		},
		{
			Arch:      specs.ArchMIPS64N32,
			SubArches: []specs.Arch{specs.ArchMIPS, specs.ArchMIPS64},
		},
		{
			Arch:      specs.ArchMIPSEL64,
			SubArches: []specs.Arch{specs.ArchMIPSEL, specs.ArchMIPSEL64N32},
		},
		{
			Arch:      specs.ArchMIPSEL64N32,
			SubArches: []specs.Arch{specs.ArchMIPSEL, specs.ArchMIPSEL64},
		},
		{
			Arch:      specs.ArchS390X,
			SubArches: []specs.Arch{specs.ArchS390},
		},
		{
			Arch:      specs.ArchRISCV64,
			SubArches: nil,
		},
	}
}

// DefaultProfile defines the allowed syscalls for the default seccomp profile.
func DefaultProfile() *Seccomp {
	nosys := uint(unix.ENOSYS)
	syscalls := []*Syscall{
		{
			LinuxSyscall: specs.LinuxSyscall{
				Names: []string{
					"accept",
					"accept4",
					"access",
					"adjtimex",
					"alarm",
					"bind",
					"brk",
					"cachestat", // kernel v6.5, libseccomp v2.5.5
					"capget",
					"capset",
					"chdir",
					"chmod",
					"chown",
					"chown32",
					"clock_adjtime",
					"clock_adjtime64",
					"clock_getres",
					"clock_getres_time64",
					"clock_gettime",
					"clock_gettime64",
					"clock_nanosleep",
					"clock_nanosleep_time64",
					"close",
					"close_range",
					"connect",
					"copy_file_range",
					"creat",
					"dup",
					"dup2",
					"dup3",
					"epoll_create",
					"epoll_create1",
					"epoll_ctl",
					"epoll_ctl_old",
					"epoll_pwait",
					"epoll_pwait2",
					"epoll_wait",
					"epoll_wait_old",
					"eventfd",
					"eventfd2",
					"execve",
					"execveat",
					"exit",
					"exit_group",
					"faccessat",
					"faccessat2",
					"fadvise64",
					"fadvise64_64",
					"fallocate",
					"fanotify_mark",
					"fchdir",
					"fchmod",
					"fchmodat",
					"fchmodat2", // kernel v6.6, libseccomp v2.5.5
					"fchown",
					"fchown32",
					"fchownat",
					"fcntl",
					"fcntl64",
					"fdatasync",
					"fgetxattr",
					"flistxattr",
					"flock",
					"fork",
					"fremovexattr",
					"fsetxattr",
					"fstat",
					"fstat64",
					"fstatat64",
					"fstatfs",
					"fstatfs64",
					"fsync",
					"ftruncate",
					"ftruncate64",
					"futex",
					"futex_requeue", // kernel v6.7, libseccomp v2.5.5
					"futex_time64",
					"futex_wait", // kernel v6.7, libseccomp v2.5.5
					"futex_waitv",
					"futex_wake", // kernel v6.7, libseccomp v2.5.5
					"futimesat",
					"getcpu",
					"getcwd",
					"getdents",
					"getdents64",
					"getegid",
					"getegid32",
					"geteuid",
					"geteuid32",
					"getgid",
					"getgid32",
					"getgroups",
					"getgroups32",
					"getitimer",
					"getpeername",
					"getpgid",
					"getpgrp",
					"getpid",
					"getppid",
					"getpriority",
					"getrandom",
					"getresgid",
					"getresgid32",
					"getresuid",
					"getresuid32",
					"getrlimit",
					"get_robust_list",
					"getrusage",
					"getsid",
					"getsockname",
					"getsockopt",
					"get_thread_area",
					"gettid",
					"gettimeofday",
					"getuid",
					"getuid32",
					"getxattr",
					"getxattrat", // kernel v6.13, libseccomp v2.6.0
					"inotify_add_watch",
					"inotify_init",
					"inotify_init1",
					"inotify_rm_watch",
					"io_cancel",
					"ioctl",
					"io_destroy",
					"io_getevents",
					"io_pgetevents",
					"io_pgetevents_time64",
					"ioprio_get",
					"ioprio_set",
					"io_setup",
					"io_submit",
					"ipc",
					"kill",
					"landlock_add_rule",
					"landlock_create_ruleset",
					"landlock_restrict_self",
					"lchown",
					"lchown32",
					"lgetxattr",
					"link",
					"linkat",
					"listen",
					"listmount", // kernel v6.8, libseccomp v2.6.0
					"listxattr",
					"listxattrat", // kernel v6.13, libseccomp v2.6.0
					"llistxattr",
					"_llseek",
					"lremovexattr",
					"lseek",
					"lsetxattr",
					"lstat",
					"lstat64",
					"madvise",
					"map_shadow_stack", // kernel v6.6, libseccomp v2.5.5
					"membarrier",
					"memfd_create",
					"memfd_secret",
					"mincore",
					"mkdir",
					"mkdirat",
					"mknod",
					"mknodat",
					"mlock",
					"mlock2",
					"mlockall",
					"mmap",
					"mmap2",
					"mprotect",
					"mq_getsetattr",
					"mq_notify",
					"mq_open",
					"mq_timedreceive",
					"mq_timedreceive_time64",
					"mq_timedsend",
					"mq_timedsend_time64",
					"mq_unlink",
					"mremap",
					"mseal", // kernel v6.9, libseccomp v2.6.0
					"msgctl",
					"msgget",
					"msgrcv",
					"msgsnd",
					"msync",
					"munlock",
					"munlockall",
					"munmap",
					"name_to_handle_at",
					"nanosleep",
					"newfstatat",
					"_newselect",
					"open",
					"openat",
					"openat2",
					"pause",
					"pidfd_open",
					"pidfd_send_signal",
					"pipe",
					"pipe2",
					"pkey_alloc",
					"pkey_free",
					"pkey_mprotect",
					"poll",
					"ppoll",
					"ppoll_time64",
					"prctl",
					"pread64",
					"preadv",
					"preadv2",
					"prlimit64",
					"process_mrelease",
					"pselect6",
					"pselect6_time64",
					"pwrite64",
					"pwritev",
					"pwritev2",
					"read",
					"readahead",
					"readlink",
					"readlinkat",
					"readv",
					"recv",
					"recvfrom",
					"recvmmsg",
					"recvmmsg_time64",
					"recvmsg",
					"remap_file_pages",
					"removexattr",
					"removexattrat", // kernel v6.13, libseccomp v2.6.0
					"rename",
					"renameat",
					"renameat2",
					"restart_syscall",
					"riscv_hwprobe", // kernel v6.12, libseccomp v2.6.0
					"rmdir",
					"rseq",
					"rt_sigaction",
					"rt_sigpending",
					"rt_sigprocmask",
					"rt_sigqueueinfo",
					"rt_sigreturn",
					"rt_sigsuspend",
					"rt_sigtimedwait",
					"rt_sigtimedwait_time64",
					"rt_tgsigqueueinfo",
					"sched_getaffinity",
					"sched_getattr",
					"sched_getparam",
					"sched_get_priority_max",
					"sched_get_priority_min",
					"sched_getscheduler",
					"sched_rr_get_interval",
					"sched_rr_get_interval_time64",
					"sched_setaffinity",
					"sched_setattr",
					"sched_setparam",
					"sched_setscheduler",
					"sched_yield",
					"seccomp",
					"select",
					"semctl",
					"semget",
					"semop",
					"semtimedop",
					"semtimedop_time64",
					"send",
					"sendfile",
					"sendfile64",
					"sendmmsg",
					"sendmsg",
					"sendto",
					"setfsgid",
					"setfsgid32",
					"setfsuid",
					"setfsuid32",
					"setgid",
					"setgid32",
					"setgroups",
					"setgroups32",
					"setitimer",
					"setpgid",
					"setpriority",
					"setregid",
					"setregid32",
					"setresgid",
					"setresgid32",
					"setresuid",
					"setresuid32",
					"setreuid",
					"setreuid32",
					"setrlimit",
					"set_robust_list",
					"setsid",
					"setsockopt",
					"set_thread_area",
					"set_tid_address",
					"setuid",
					"setuid32",
					"setxattr",
					"setxattrat", // kernel v6.13, libseccomp v2.6.0
					"shmat",
					"shmctl",
					"shmdt",
					"shmget",
					"shutdown",
					"sigaltstack",
					"signalfd",
					"signalfd4",
					"sigprocmask",
					"sigreturn",
					"socketcall",
					"socketpair",
					"splice",
					"stat",
					"stat64",
					"statfs",
					"statfs64",
					"statmount", // kernel v6.8, libseccomp v2.6.0
					"statx",
					"symlink",
					"symlinkat",
					"sync",
					"sync_file_range",
					"syncfs",
					"sysinfo",
					"tee",
					"tgkill",
					"time",
					"timer_create",
					"timer_delete",
					"timer_getoverrun",
					"timer_gettime",
					"timer_gettime64",
					"timer_settime",
					"timer_settime64",
					"timerfd_create",
					"timerfd_gettime",
					"timerfd_gettime64",
					"timerfd_settime",
					"timerfd_settime64",
					"times",
					"tkill",
					"truncate",
					"truncate64",
					"ugetrlimit",
					"umask",
					"uname",
					"unlink",
					"unlinkat",
					"uretprobe", // kernel v6.11, libseccomp v2.6.0
					"utime",
					"utimensat",
					"utimensat_time64",
					"utimes",
					"vfork",
					"vmsplice",
					"wait4",
					"waitid",
					"waitpid",
					"write",
					"writev",
				},
				Action: specs.ActAllow,
			},
		},
		{
			LinuxSyscall: specs.LinuxSyscall{
				Names: []string{
					"process_vm_readv",
					"process_vm_writev",
					"ptrace",
				},
				Action: specs.ActAllow,
			},
			Includes: &Filter{
				MinKernel: &KernelVersion{4, 8},
			},
		},
		{
			LinuxSyscall: specs.LinuxSyscall{
				Names:  []string{"socket"},
				Action: specs.ActAllow,
				Args: []specs.LinuxSeccompArg{
					{
						Index: 0,
						Value: unix.AF_VSOCK,
						Op:    specs.OpNotEqual,
					},
				},
			},
		},
		{
			LinuxSyscall: specs.LinuxSyscall{
				Names:  []string{"personality"},
				Action: specs.ActAllow,
				Args: []specs.LinuxSeccompArg{
					{
						Index: 0,
						Value: 0x0,
						Op:    specs.OpEqualTo,
					},
				},
			},
		},
		{
			LinuxSyscall: specs.LinuxSyscall{
				Names:  []string{"personality"},
				Action: specs.ActAllow,
				Args: []specs.LinuxSeccompArg{
					{
						Index: 0,
						Value: 0x0008,
						Op:    specs.OpEqualTo,
					},
				},
			},
		},
		{
			LinuxSyscall: specs.LinuxSyscall{
				Names:  []string{"personality"},
				Action: specs.ActAllow,
				Args: []specs.LinuxSeccompArg{
					{
						Index: 0,
						Value: 0x20000,
						Op:    specs.OpEqualTo,
					},
				},
			},
		},
		{
			LinuxSyscall: specs.LinuxSyscall{
				Names:  []string{"personality"},
				Action: specs.ActAllow,
				Args: []specs.LinuxSeccompArg{
					{
						Index: 0,
						Value: 0x20008,
						Op:    specs.OpEqualTo,
					},
				},
			},
		},
		{
			LinuxSyscall: specs.LinuxSyscall{
				Names:  []string{"personality"},
				Action: specs.ActAllow,
				Args: []specs.LinuxSeccompArg{
					{
						Index: 0,
						Value: 0xffffffff,
						Op:    specs.OpEqualTo,
					},
				},
			},
		},
		{
			LinuxSyscall: specs.LinuxSyscall{
				Names: []string{
					"sync_file_range2",
					"swapcontext",
				},
				Action: specs.ActAllow,
			},
			Includes: &Filter{
				Arches: []string{"ppc64le"},
			},
		},
		{
			LinuxSyscall: specs.LinuxSyscall{
				Names: []string{
					"arm_fadvise64_64",
					"arm_sync_file_range",
					"sync_file_range2",
					"breakpoint",
					"cacheflush",
					"set_tls",
				},
				Action: specs.ActAllow,
			},
			Includes: &Filter{
				Arches: []string{"arm", "arm64"},
			},
		},
		{
			LinuxSyscall: specs.LinuxSyscall{
				Names: []string{
					"arch_prctl",
				},
				Action: specs.ActAllow,
			},
			Includes: &Filter{
				Arches: []string{"amd64", "x32"},
			},
		},
		{
			LinuxSyscall: specs.LinuxSyscall{
				Names: []string{
					"modify_ldt",
				},
				Action: specs.ActAllow,
			},
			Includes: &Filter{
				Arches: []string{"amd64", "x32", "x86"},
			},
		},
		{
			LinuxSyscall: specs.LinuxSyscall{
				Names: []string{
					"s390_pci_mmio_read",
					"s390_pci_mmio_write",
					"s390_runtime_instr",
				},
				Action: specs.ActAllow,
			},
			Includes: &Filter{
				Arches: []string{"s390", "s390x"},
			},
		},
		{
			LinuxSyscall: specs.LinuxSyscall{
				Names: []string{
					"riscv_flush_icache",
				},
				Action: specs.ActAllow,
			},
			Includes: &Filter{
				Arches: []string{"riscv64"},
			},
		},
		{
			LinuxSyscall: specs.LinuxSyscall{
				Names: []string{
					"open_by_handle_at",
				},
				Action: specs.ActAllow,
			},
			Includes: &Filter{
				Caps: []string{"CAP_DAC_READ_SEARCH"},
			},
		},
		{
			LinuxSyscall: specs.LinuxSyscall{
				Names: []string{
					"bpf",
					"clone",
					"clone3",
					"fanotify_init",
					"fsconfig",
					"fsmount",
					"fsopen",
					"fspick",
					"lookup_dcookie",
					"lsm_get_self_attr", // kernel v6.8, libseccomp v2.6.0
					"lsm_list_modules",  // kernel v6.8, libseccomp v2.6.0
					"lsm_set_self_attr", // kernel v6.8, libseccomp v2.6.0
					"mount",
					"mount_setattr",
					"move_mount",
					"open_tree",
					"perf_event_open",
					"quotactl",
					"quotactl_fd",
					"setdomainname",
					"sethostname",
					"setns",
					"syslog",
					"umount",
					"umount2",
					"unshare",
				},
				Action: specs.ActAllow,
			},
			Includes: &Filter{
				Caps: []string{"CAP_SYS_ADMIN"},
			},
		},
		{
			LinuxSyscall: specs.LinuxSyscall{
				Names: []string{
					"clone",
				},
				Action: specs.ActAllow,
				Args: []specs.LinuxSeccompArg{
					{
						Index:    0,
						Value:    unix.CLONE_NEWNS | unix.CLONE_NEWUTS | unix.CLONE_NEWIPC | unix.CLONE_NEWUSER | unix.CLONE_NEWPID | unix.CLONE_NEWNET | unix.CLONE_NEWCGROUP,
						ValueTwo: 0,
						Op:       specs.OpMaskedEqual,
					},
				},
			},
			Excludes: &Filter{
				Caps:   []string{"CAP_SYS_ADMIN"},
				Arches: []string{"s390", "s390x"},
			},
		},
		{
			LinuxSyscall: specs.LinuxSyscall{
				Names: []string{
					"clone",
				},
				Action: specs.ActAllow,
				Args: []specs.LinuxSeccompArg{
					{
						Index:    1,
						Value:    unix.CLONE_NEWNS | unix.CLONE_NEWUTS | unix.CLONE_NEWIPC | unix.CLONE_NEWUSER | unix.CLONE_NEWPID | unix.CLONE_NEWNET | unix.CLONE_NEWCGROUP,
						ValueTwo: 0,
						Op:       specs.OpMaskedEqual,
					},
				},
			},
			Comment: "s390 parameter ordering for clone is different",
			Includes: &Filter{
				Arches: []string{"s390", "s390x"},
			},
			Excludes: &Filter{
				Caps: []string{"CAP_SYS_ADMIN"},
			},
		},
		{
			LinuxSyscall: specs.LinuxSyscall{
				Names: []string{
					"clone3",
				},
				Action:   specs.ActErrno,
				ErrnoRet: &nosys,
			},
			Excludes: &Filter{
				Caps: []string{"CAP_SYS_ADMIN"},
			},
		},
		{
			LinuxSyscall: specs.LinuxSyscall{
				Names: []string{
					"reboot",
				},
				Action: specs.ActAllow,
			},
			Includes: &Filter{
				Caps: []string{"CAP_SYS_BOOT"},
			},
		},
		{
			LinuxSyscall: specs.LinuxSyscall{
				Names: []string{
					"chroot",
				},
				Action: specs.ActAllow,
			},
			Includes: &Filter{
				Caps: []string{"CAP_SYS_CHROOT"},
			},
		},
		{
			LinuxSyscall: specs.LinuxSyscall{
				Names: []string{
					"delete_module",
					"init_module",
					"finit_module",
				},
				Action: specs.ActAllow,
			},
			Includes: &Filter{
				Caps: []string{"CAP_SYS_MODULE"},
			},
		},
		{
			LinuxSyscall: specs.LinuxSyscall{
				Names: []string{
					"acct",
				},
				Action: specs.ActAllow,
			},
			Includes: &Filter{
				Caps: []string{"CAP_SYS_PACCT"},
			},
		},
		{
			LinuxSyscall: specs.LinuxSyscall{
				Names: []string{
					"kcmp",
					"pidfd_getfd",
					"process_madvise",
					"process_vm_readv",
					"process_vm_writev",
					"ptrace",
				},
				Action: specs.ActAllow,
			},
			Includes: &Filter{
				Caps: []string{"CAP_SYS_PTRACE"},
			},
		},
		{
			LinuxSyscall: specs.LinuxSyscall{
				Names: []string{
					"iopl",
					"ioperm",
				},
				Action: specs.ActAllow,
			},
			Includes: &Filter{
				Caps: []string{"CAP_SYS_RAWIO"},
			},
		},
		{
			LinuxSyscall: specs.LinuxSyscall{
				Names: []string{
					"settimeofday",
					"stime",
					"clock_settime",
					"clock_settime64",
				},
				Action: specs.ActAllow,
			},
			Includes: &Filter{
				Caps: []string{"CAP_SYS_TIME"},
			},
		},
		{
			LinuxSyscall: specs.LinuxSyscall{
				Names: []string{
					"vhangup",
				},
				Action: specs.ActAllow,
			},
			Includes: &Filter{
				Caps: []string{"CAP_SYS_TTY_CONFIG"},
			},
		},
		{
			LinuxSyscall: specs.LinuxSyscall{
				Names: []string{
					"get_mempolicy",
					"mbind",
					"set_mempolicy",
					"set_mempolicy_home_node", // kernel v5.17, libseccomp v2.5.4
				},
				Action: specs.ActAllow,
			},
			Includes: &Filter{
				Caps: []string{"CAP_SYS_NICE"},
			},
		},
		{
			LinuxSyscall: specs.LinuxSyscall{
				Names: []string{
					"syslog",
				},
				Action: specs.ActAllow,
			},
			Includes: &Filter{
				Caps: []string{"CAP_SYSLOG"},
			},
		},
		{
			LinuxSyscall: specs.LinuxSyscall{
				Names: []string{
					"bpf",
				},
				Action: specs.ActAllow,
			},
			Includes: &Filter{
				Caps: []string{"CAP_BPF"},
			},
		},
		{
			LinuxSyscall: specs.LinuxSyscall{
				Names: []string{
					"perf_event_open",
				},
				Action: specs.ActAllow,
			},
			Includes: &Filter{
				Caps: []string{"CAP_PERFMON"},
			},
		},
	}

	errnoRet := uint(unix.EPERM)
	return &Seccomp{
		LinuxSeccomp: specs.LinuxSeccomp{
			DefaultAction:   specs.ActErrno,
			DefaultErrnoRet: &errnoRet,
		},
		ArchMap:  arches(),
		Syscalls: syscalls,
	}
}
