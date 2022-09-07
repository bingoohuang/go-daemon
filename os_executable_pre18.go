//go:build !go1.8
// +build !go1.8

package godaemon

import (
	"github.com/kardianos/osext"
)

func osExecutable() (string, error) {
	return osext.Executable()
}
