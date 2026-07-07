package events

import (
	"time"

	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/core"
	"github.com/nyaruka/goflow/utils"
)

// EventUUID is the UUID of an event
type EventUUID uuids.UUID

// NewEventUUID generates a new UUID for an event
func NewEventUUID() EventUUID { return EventUUID(uuids.NewV7()) }

// Step describes the step in a flow at which an event occurred. It is set by the engine on the events
// it generates during a sprint for the benefit of callers, but is not persisted with events.
type Step struct {
	Flow *assets.FlowReference
	Node core.NodeUUID
}

// Event describes a state change
type Event interface {
	utils.Typed

	UUID() EventUUID
	CreatedOn() time.Time
	Step() *Step
	SetStep(*Step)
	SetUser(*assets.UserReference, string)
}

// EventLogger is a callback invoked when an event has been generated
type EventLogger func(Event)
