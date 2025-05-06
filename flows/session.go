package flows

import (
	"context"
	"time"

	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
)

// SprintUUID is the UUID of a sprint
type SprintUUID uuids.UUID

// SessionUUID is the UUID of a session
type SessionUUID uuids.UUID

// NewSessionUUID generates a new UUID for a session
func NewSessionUUID() SessionUUID {
	return SessionUUID(uuids.NewV4())
}

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
	Modifiers() []Modifier
	Events() []Event
	Segments() []Segment
}

// Session represents the session of a flow run which may contain many runs
type Session interface {
	Assets() SessionAssets

	UUID() SessionUUID
	Type() FlowType
	SetType(FlowType)

	Environment() envs.Environment
	SetEnvironment(envs.Environment)
	MergedEnvironment() envs.Environment

	Contact() *Contact
	SetContact(*Contact)

	Input() Input
	SetInput(Input)

	Status() SessionStatus
	Trigger() Trigger
	CurrentResume() Resume
	BatchStart() bool
	PushFlow(Flow, Run, bool)

	Resume(context.Context, Resume) (Sprint, error)
	Runs() []Run
	GetRun(RunUUID) (Run, error)
	FindStep(uuid StepUUID) (Run, Step)
	GetCurrentChild(Run) Run
	ParentRun() RunSummary
	CurrentContext() *types.XObject
	History() *SessionHistory

	Engine() Engine
}

// SessionHistory provides information about the sessions that caused this session
type SessionHistory struct {
	ParentUUID          SessionUUID `json:"parent_uuid"`
	Ancestors           int         `json:"ancestors"`
	AncestorsSinceInput int         `json:"ancestors_since_input"`
}

// Advance moves history forward to a new parent
func (h *SessionHistory) Advance(newParent SessionUUID, receivedInput bool) *SessionHistory {
	ancestorsSinceinput := 0
	if !receivedInput {
		ancestorsSinceinput = h.AncestorsSinceInput + 1
	}

	return &SessionHistory{
		ParentUUID:          newParent,
		Ancestors:           h.Ancestors + 1,
		AncestorsSinceInput: ancestorsSinceinput,
	}
}

// EmptyHistory is used for a session which has no history
var EmptyHistory = &SessionHistory{}

// NewChildHistory creates a new history for a child of the given session
func NewChildHistory(parent Session) *SessionHistory {
	return parent.History().Advance(parent.UUID(), sessionReceivedInput(parent))
}

// looks through a session's events to see if it received input
func sessionReceivedInput(s Session) bool {
	for _, r := range s.Runs() {
		if r.ReceivedInput() {
			return true
		}
	}
	return false
}
