package flows

import (
	"time"

	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/utils"
)

type EventUUID uuids.UUID

// NewEventUUID generates a new UUID for an event
func NewEventUUID() EventUUID { return EventUUID(uuids.NewV7()) }

// Event describes a state change
type Event interface {
	utils.Typed

	UUID() EventUUID
	CreatedOn() time.Time
	StepUUID() StepUUID
	SetStepUUID(StepUUID)
}

// EventCallback is a callback invoked when an event has been generated
type EventCallback func(Event)
