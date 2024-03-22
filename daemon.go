package godaemon

import (
	"errors"
	"os"
	"strconv"
	"strings"

	"github.com/bingoohuang/q"
)

var errNotSupported = errors.New("daemon: Non-POSIX OS is not supported")

// MarkName marks of daemon process - system environment variable _GO_DAEMON={pid}
const (
	MarkName      = "_GO_DAEMON_MARK"
	MarkParentPID = "_GO_DAEMON_PID"
	MarkValue     = "1"
)

// FilePerm is the default file permissions for log and pid files.
const FilePerm = os.FileMode(0o640)

// ClearReborn clear the reborn env.
func ClearReborn() error {
	err := os.Unsetenv(MarkName)
	return err
}

// WasReborn returns true in child process (daemon) and false in parent process.
func WasReborn() bool {
	markValue := os.Getenv(MarkName)
	return markValue == MarkValue
}

// GetParentPID returns the parent pid if forked by godaemon.
func GetParentPID() (pid int) {
	value := os.Getenv(MarkParentPID)
	pid, _ = strconv.Atoi(value)
	return pid
}

// Reborn runs second copy of current process in the given context.
// function executes separate parts of code in child process and parent process
// and provides demonization of child process. It looks similar as the
// fork-daemonization, but goroutine-safe.
// In success returns *os.Process in parent process and nil in child process.
// Otherwise, returns error.
func (d *Context) Reborn() (child *os.Process, err error) {
	return d.reborn()
}

// Search searches daemons process by given in context pid file name.
// If success returns pointer on daemons os.Process structure,
// else returns error. Returns nil if filename is empty.
func (d *Context) Search() (daemon *os.Process, err error) {
	return d.search()
}

// Release provides correct pid-file release in daemon.
func (d *Context) Release() error {
	return d.release()
}

func (d *Context) Clean() error {
	return d.clean()
}

// Option is options for Daemonize function.
type Option struct {
	Daemon        bool
	Debug         bool
	LogFileName   string
	ParentProcess func(child *os.Process)
}

// OptionFn is an option function prototype.
type OptionFn func(option *Option)

// WithDaemon set the deamon flag.
func WithDaemon(v bool) OptionFn {
	return func(o *Option) {
		o.Daemon = v
	}
}

// WithLogFileName set the log file name.
func WithLogFileName(v string) OptionFn {
	return func(o *Option) {
		o.LogFileName = v
	}
}

// WithParentProcess set the custom parent processor, if not set, the default is os.Exit(0)
func WithParentProcess(f func(child *os.Process)) OptionFn {
	return func(o *Option) {
		o.ParentProcess = f
	}
}

// Daemonize set the current process daemonized
func Daemonize(optionFns ...OptionFn) {
	option := &Option{
		Daemon: true,
	}
	for _, fn := range optionFns {
		fn(option)
	}

	if !option.Daemon {
		return
	}

	// goland 启动时，不进入后台模式
	if strings.Contains(os.Args[0], "/Caches/JetBrains") {
		return
	}

	workDir, err := os.Getwd()
	if err != nil {
		q.D("Getwd error", err)
	}

	ctx := &Context{
		WorkDir:     workDir,
		LogFileName: option.LogFileName,
	}

	child, err := ctx.Reborn()
	if err != nil {
		q.D("reborn error", err)
	}
	if child != nil {
		// 有孩子，是父进程
		if option.ParentProcess != nil {
			option.ParentProcess(child)
		} else {
			os.Exit(0)
		}
	}
	// 子进程，继续
}

type Credential struct {
	Uid         uint32   // User ID.
	Gid         uint32   // Group ID.
	Groups      []uint32 // Supplementary group IDs.
	NoSetGroups bool     // If true, don't set supplementary groups
}
