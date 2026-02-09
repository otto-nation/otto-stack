//go:build windows

package registry

import (
	"os"
	"syscall"
	"time"
	"unsafe"
)

const (
	lockRetries = 5
	lockDelay   = 100 * time.Millisecond
)

var (
	kernel32     = syscall.NewLazyDLL("kernel32.dll")
	lockFileEx   = kernel32.NewProc("LockFileEx")
	unlockFileEx = kernel32.NewProc("UnlockFileEx")
)

const (
	lockfileExclusiveLock   = 0x00000002
	lockfileFailImmediately = 0x00000001
	errorLockViolation      = 33
)

// lockFile acquires an exclusive lock on the file
func lockFile(f *os.File) error {
	for range lockRetries {
		overlapped := &syscall.Overlapped{}
		r1, _, err := lockFileEx.Call(
			uintptr(f.Fd()),
			uintptr(lockfileExclusiveLock|lockfileFailImmediately),
			0,
			0xFFFFFFFF,
			0xFFFFFFFF,
			uintptr(unsafe.Pointer(overlapped)),
		)
		if r1 != 0 {
			return nil
		}
		if errno, ok := err.(syscall.Errno); !ok || errno != errorLockViolation {
			return err
		}
		time.Sleep(lockDelay)
	}
	return syscall.Errno(errorLockViolation)
}

// unlockFile releases the lock on the file
func unlockFile(f *os.File) error {
	overlapped := &syscall.Overlapped{}
	r1, _, err := unlockFileEx.Call(
		uintptr(f.Fd()),
		0,
		0xFFFFFFFF,
		0xFFFFFFFF,
		uintptr(unsafe.Pointer(overlapped)),
	)
	if r1 == 0 {
		return err
	}
	return nil
}
