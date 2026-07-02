package events

import (
	"fmt"
	"time"

	"github.com/nyaruka/gocommon/dates"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/utils"
)

var registeredTypes = map[string](func() Event){}

// registers a new type of event
func registerType(name string, initFunc func() Event) {
	registeredTypes[name] = initFunc
}

// BaseEvent is the base of all event types
type BaseEvent struct {
	UUID_      EventUUID `json:"uuid"                validate:"required,uuid"`
	Type_      string    `json:"type"                validate:"required"`
	CreatedOn_ time.Time `json:"created_on"          validate:"required"`

	// can be used by callers to locate events but not persisted
	Step_ Step `json:"-"`

	// not set by engine but can be set by callers for storage of events
	User_ *assets.UserReference `json:"_user,omitempty"`
	Via_  string                `json:"_via,omitempty"`
}

// NewBaseEvent creates a new base event
func NewBaseEvent(typ string) BaseEvent {
	return NewBaseEventWithUUID(NewEventUUID(), typ)
}

// NewBaseEventWithUUID creates a new base event with the given pre-allocated UUID. Used when the caller
// needs to know the event's UUID before constructing the event (e.g. to pass it to an outbound provider
// call so that callbacks can be correlated back to the event).
func NewBaseEventWithUUID(uuid EventUUID, typ string) BaseEvent {
	return BaseEvent{UUID_: uuid, Type_: typ, CreatedOn_: dates.Now()}
}

// UUID returns the UUID of this event
func (e *BaseEvent) UUID() EventUUID { return e.UUID_ }

// Type returns the type of this event
func (e *BaseEvent) Type() string { return e.Type_ }

// CreatedOn returns the created on time of this event
func (e *BaseEvent) CreatedOn() time.Time { return e.CreatedOn_ }

// Step returns the step in the path where this event occurred
func (e *BaseEvent) Step() Step { return e.Step_ }

// SetStep sets the UUID of the step in the path where this event occurred
func (e *BaseEvent) SetStep(s Step) { e.Step_ = s }

// SetUser can be used by callers to set the user associated with this event
func (e *BaseEvent) SetUser(u *assets.UserReference, via string) {
	e.User_ = u
	e.Via_ = via
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

// Read reads a single event from the given JSON
func Read(data []byte) (Event, error) {
	typeName, err := utils.ReadTypeFromJSON(data)
	if err != nil {
		return nil, err
	}

	f := registeredTypes[typeName]
	if f == nil {
		return nil, fmt.Errorf("unknown type: '%s'", typeName)
	}

	event := f()
	return event, utils.UnmarshalAndValidate(data, event)
}
