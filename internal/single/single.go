// Package single enforces that only one instance of the Ease UI app
// runs per user account. It is implemented as an exclusive,
// non-blocking file lock (flock on Unix, LockFileEx on Windows);
// the kernel releases the lock automatically when the process
// exits, so a hard kill never strands a stale lockfile.
package single

import (
	"errors"
	"os"
	"path/filepath"
)

// ErrAlreadyRunning is returned by Acquire when another instance
// already holds the singleton lock for this user.
var ErrAlreadyRunning = errors.New("ease-ui: another instance is already running")

// LockPath returns the default location of the singleton lock
// file. It is created under the user's config directory and the
// parent directory is ensured to exist.
func LockPath() (string, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(base, "ease-ui")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	return filepath.Join(dir, "singleton.lock"), nil
}

// Acquire takes the singleton lock for this user. The returned
// release function must be called (typically via defer) to close
// the underlying file descriptor; the kernel will also release
// the lock automatically if the process exits without calling it.
//
// If another instance already holds the lock, Acquire returns
// ErrAlreadyRunning and a nil release function.
func Acquire() (func() error, error) {
	path, err := LockPath()
	if err != nil {
		return nil, err
	}
	return acquireAt(path)
}
