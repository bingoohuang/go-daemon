//go:build darwin || dragonfly || freebsd || linux || netbsd || openbsd || plan9
// +build darwin dragonfly freebsd linux netbsd openbsd plan9

package godaemon

import (
	"errors"
	"syscall"
)

func lockFile(fd uintptr) error {
	err := syscall.Flock(int(fd), syscall.LOCK_EX|syscall.LOCK_NB)
	if errors.Is(err, syscall.EWOULDBLOCK) {
		err = ErrWouldBlock
	}
	return err
}

func unlockFile(fd uintptr) error {
	err := syscall.Flock(int(fd), syscall.LOCK_UN)
	if errors.Is(err, syscall.EWOULDBLOCK) {
		err = ErrWouldBlock
	}
	return err
}
