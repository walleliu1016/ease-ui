//go:build unix

package single

import (
	"os"
	"syscall"
)

func acquireAt(path string) (func() error, error) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0o644)
	if err != nil {
		return nil, err
	}
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		_ = f.Close()
		if err == syscall.EWOULDBLOCK || err == syscall.EAGAIN {
			return nil, ErrAlreadyRunning
		}
		return nil, err
	}
	// Closing the fd releases the flock automatically.
	return f.Close, nil
}
