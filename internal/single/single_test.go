package single

import (
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAcquire_Exclusive(t *testing.T) {
	// Simulate a second process by opening a second fd against the
	// same file. flock is per-open-file-description, so two
	// independent OpenFile calls on the same path must contend.
	path := filepath.Join(t.TempDir(), "lock")

	release1, err := acquireAt(path)
	require.NoError(t, err)
	defer release1()

	release2, err := acquireAt(path)
	require.ErrorIs(t, err, ErrAlreadyRunning)
	assert.Nil(t, release2)
}

func TestAcquire_AfterRelease(t *testing.T) {
	path := filepath.Join(t.TempDir(), "lock")

	release1, err := acquireAt(path)
	require.NoError(t, err)
	require.NoError(t, release1())

	release2, err := acquireAt(path)
	require.NoError(t, err)
	assert.NotNil(t, release2)
	require.NoError(t, release2())
}

func TestAcquire_ReentrantSameProcess(t *testing.T) {
	// On Unix, flock is reentrant within the same fd but not
	// across different fds. Our API always opens a new fd, so
	// a second Acquire from the same goroutine must still fail.
	path := filepath.Join(t.TempDir(), "lock")

	release1, err := acquireAt(path)
	require.NoError(t, err)
	defer release1()

	_, err = acquireAt(path)
	require.Error(t, err)
	require.True(t, errors.Is(err, ErrAlreadyRunning))
}

func TestLockPath(t *testing.T) {
	path, err := LockPath()
	require.NoError(t, err)
	assert.True(t, strings.HasSuffix(path, filepath.Join("ease-ui", "singleton.lock")),
		"unexpected lock path: %q", path)
}
