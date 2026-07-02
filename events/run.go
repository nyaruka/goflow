package events

import "github.com/nyaruka/gocommon/uuids"

// RunUUID is the UUID of a flow run
type RunUUID uuids.UUID

// NewRunUUID generates a new UUID for a run
func NewRunUUID() RunUUID { return RunUUID(uuids.NewV7()) }

// RunStatus represents the current status of the flow run
type RunStatus string

const (
	// RunStatusActive represents a run that is still active
	RunStatusActive RunStatus = "active"

	// RunStatusCompleted represents a run that has run to completion
	RunStatusCompleted RunStatus = "completed"

	// RunStatusWaiting represents a run which is waiting for something from the caller
	RunStatusWaiting RunStatus = "waiting"

	// RunStatusFailed represents a run that encountered an unrecoverable error
	RunStatusFailed RunStatus = "failed"

	// RunStatusExpired represents a run that expired due to inactivity
	RunStatusExpired RunStatus = "expired"

	// RunStatusInterrupted is never used by the engine but callers may put runs into this state
	RunStatusInterrupted RunStatus = "interrupted"
)
