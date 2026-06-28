//go:build windows

package single

import (
	"os"
	"syscall"
	"unsafe"
)

var (
	modkernel32      = syscall.NewLazyDLL("kernel32.dll")
	procLockFileEx   = modkernel32.NewProc("LockFileEx")
	procUnlockFileEx = modkernel32.NewProc("UnlockFileEx")
)

const (
	lockfileFailImmediately = 0x00000001
	lockfileExclusiveLock   = 0x00000002
)

func acquireAt(path string) (func() error, error) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0o644)
	if err != nil {
		return nil, err
	}
	var overlapped syscall.Overlapped
	// Lock 1 byte at offset 0; the size and offset are split into
	// low/high 32-bit halves as the Win32 API requires.
	r1, _, e1 := procLockFileEx.Call(
		uintptr(f.Fd()),
		uintptr(lockfileFailImmediately|lockfileExclusiveLock),
		0,
		1, 0,
		uintptr(unsafe.Pointer(&overlapped)),
	)
	if r1 == 0 {
		_ = f.Close()
		if e1 == syscall.Errno(0x21) /* ERROR_LOCK_VIOLATION */ ||
			e1 == syscall.Errno(0x3E5) /* ERROR_IO_PENDING */ {
			return nil, ErrAlreadyRunning
		}
		return nil, e1
	}
	return f.Close, nil
}
