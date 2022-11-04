package godaemon

import (
	"errors"
	"log"
	"os"
)

var errNotSupported = errors.New("daemon: Non-POSIX OS is not supported")

// Mark of daemon process - system environment variable _GO_DAEMON=1
const (
	MarkName  = "_GO_DAEMON"
	MarkValue = "1"
)

// FilePerm is the default file permissions for log and pid files.
const FilePerm = os.FileMode(0o640)

// WasReborn returns true in child process (daemon) and false in parent process.
func WasReborn() bool {
	return os.Getenv(MarkName) == MarkValue
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

// Deamonize set the current process daemonozed
func Deamonize(daemon bool) {
	if !daemon {
		return
	}

	if p, _ := new(Context).Reborn(); p != nil {
		os.Exit(0)
	}

	log.Printf("--- daemon started --")
}
