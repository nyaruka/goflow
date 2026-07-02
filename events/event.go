package events

import (
	"time"

	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/utils"
)

// EventUUID is the UUID of an event
type EventUUID uuids.UUID

// NewEventUUID generates a new UUID for an event
func NewEventUUID() EventUUID { return EventUUID(uuids.NewV7()) }

// NodeUUID is a UUID of a flow node
type NodeUUID uuids.UUID

// ExitUUID is the UUID of a node exit
type ExitUUID uuids.UUID

// StepUUID is the UUID of a run step
type StepUUID uuids.UUID

// InputUUID is the UUID of an input
type InputUUID uuids.UUID

// Step is the location in a flow run where an event occurred. It is implemented by the engine's run step type -
// this narrow interface exists so that events don't have to depend on the engine's run types.
type Step interface {
	UUID() StepUUID
	NodeUUID() NodeUUID
	ExitUUID() ExitUUID
	ArrivedOn() time.Time
}

// Event describes a state change
type Event interface {
	utils.Typed

	UUID() EventUUID
	CreatedOn() time.Time
	Step() Step
	SetStep(Step)
	SetUser(*assets.UserReference, string)
}

// EventLogger is a callback invoked when an event has been generated
type EventLogger func(Event)
