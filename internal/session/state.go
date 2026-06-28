package session

import "fmt"

type State int

const (
	StateIdle State = iota
	StateRunning
	StateAwaitingPermission
	StateDone
)

func (s State) String() string {
	switch s {
	case StateIdle:
		return "idle"
	case StateRunning:
		return "running"
	case StateAwaitingPermission:
		return "awaiting_permission"
	case StateDone:
		return "done"
	}
	return fmt.Sprintf("unknown(%d)", int(s))
}

func (s State) IsValid() bool {
	return s >= StateIdle && s <= StateDone
}
