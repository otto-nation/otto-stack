//go:build !windows

package registry

import (
	"os"
	"syscall"
	"time"
)

const (
	lockRetries = 5
	lockDelay   = 100 * time.Millisecond
)

// lockFile acquires an exclusive lock on the file
func lockFile(f *os.File) error {
	for range lockRetries {
		err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
		if err == nil {
			return nil
		}
		if err != syscall.EWOULDBLOCK {
			return err
		}
		time.Sleep(lockDelay)
	}
	return syscall.EWOULDBLOCK
}

// unlockFile releases the lock on the file
func unlockFile(f *os.File) error {
	return syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
}
