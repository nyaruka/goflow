package flows

import (
	"context"
	"time"

	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/core"
	"github.com/nyaruka/goflow/core/events"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
)

// SprintUUID is the UUID of a sprint
type SprintUUID uuids.UUID

// SessionStatus represents the current status of the engine session
type SessionStatus string

const (
	// SessionStatusActive represents a session that is still active
	SessionStatusActive SessionStatus = "active"

	// SessionStatusCompleted represents a session that has run to completion
	SessionStatusCompleted SessionStatus = "completed"

	// SessionStatusWaiting represents a session which is waiting for something from the caller
	SessionStatusWaiting SessionStatus = "waiting"

	// SessionStatusFailed represents a session that encountered an unrecoverable error
	SessionStatusFailed SessionStatus = "failed"

	// SessionStatusExpired is never used by the engine but callers may put sessions into this state
	SessionStatusExpired SessionStatus = "expired"

	// SessionStatusInterrupted is never used by the engine but callers may put sessions into this state
	SessionStatusInterrupted SessionStatus = "interrupted"
)

// Segment is a movement on the flow graph from an exit to another node
type Segment interface {
	Flow() Flow
	Node() Node
	Exit() Exit
	Operand() string
	Destination() Node
	Time() time.Time
}

// Sprint is an interaction with the engine - i.e. a start or resume of a session
type Sprint interface {
	UUID() SprintUUID
	IsInitial() bool
	Events() []events.Event
	Segments() []Segment
	Flows() []Flow
}

// Session represents the session of a flow run which may contain many runs
type Session interface {
	Assets() SessionAssets

	UUID() core.SessionUUID
	Type() FlowType
	CreatedOn() time.Time

	Environment() envs.Environment
	MergedEnvironment() envs.Environment
	Contact() *Contact
	Call() *Call
	Input() Input
	Sprints() int
	Status() SessionStatus
	Trigger() Trigger
	CurrentResume() Resume
	BatchStart() bool
	PushFlow(Flow, Run, bool)

	Resume(context.Context, Resume) (Sprint, error)
	Runs() []Run
	ParentRun() RunSummary
	CurrentContext() *types.XObject
	History() *core.SessionHistory

	Engine() Engine
}

// NewChildHistory creates a new history for a child of the given session
func NewChildHistory(parent Session) *core.SessionHistory {
	parentHadInput := false
	for _, r := range parent.Runs() {
		if r.HadInput() {
			parentHadInput = true
			break
		}
	}

	return parent.History().Advance(parent.UUID(), parentHadInput)
}
