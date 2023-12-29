package godaemon

import (
	"errors"
	"log"
	"os"
	"strconv"
	"strings"
)

var errNotSupported = errors.New("daemon: Non-POSIX OS is not supported")

// MarkName marks of daemon process - system environment variable _GO_DAEMON={pid}
const (
	MarkName = "_GO_DAEMON"
)

// FilePerm is the default file permissions for log and pid files.
const FilePerm = os.FileMode(0o640)

// ClearReborn clear the reborn env.
func ClearReborn() error {
	return os.Setenv(MarkName, "")
}

// WasReborn returns true in child process (daemon) and false in parent process.
func WasReborn() bool {
	return os.Getenv(MarkName) == strconv.Itoa(os.Getppid())
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

// Option is options for Daemonize function.
type Option struct {
	Daemon      bool
	Debug       bool
	LogFileName string
}

// OptionFn is an option function prototype.
type OptionFn func(option *Option)

// WithDaemon set the deamon flag.
func WithDaemon(v bool) OptionFn {
	return func(option *Option) {
		option.Daemon = v
	}
}

// WithDebug set the debug flag.
func WithDebug(v bool) OptionFn {
	return func(option *Option) {
		option.Debug = v
	}
}

// WithLogFileName set the log file name.
func WithLogFileName(v string) OptionFn {
	return func(option *Option) {
		option.LogFileName = v
	}
}

// Daemonize set the current process daemonized
func Daemonize(optionFns ...OptionFn) {
	option := &Option{
		Daemon: true,
		Debug:  false,
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
		log.Panicf("get cwd error: %v", err)
	}

	ctx := &Context{
		WorkDir:     workDir,
		LogFileName: option.LogFileName,
	}

	if p, _ := ctx.Reborn(); p != nil {
		os.Exit(0)
	}

	if option.Debug {
		log.Printf("--- daemon started --")
	}
}
