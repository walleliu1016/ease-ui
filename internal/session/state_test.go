package session

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestState_String(t *testing.T) {
	assert.Equal(t, "idle", StateIdle.String())
	assert.Equal(t, "running", StateRunning.String())
	assert.Equal(t, "awaiting_permission", StateAwaitingPermission.String())
}

func TestState_IsValid(t *testing.T) {
	assert.True(t, StateIdle.IsValid())
	assert.True(t, StateRunning.IsValid())
	assert.True(t, StateAwaitingPermission.IsValid())
	assert.False(t, State(99).IsValid())
}
