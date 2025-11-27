package events

import (
	"fmt"
	"time"

	"github.com/nyaruka/gocommon/dates"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

var registeredTypes = map[string](func() flows.Event){}

// registers a new type of event
func registerType(name string, initFunc func() flows.Event) {
	registeredTypes[name] = initFunc
}

// BaseEvent is the base of all event types
type BaseEvent struct {
	UUID_      flows.EventUUID `json:"uuid"                validate:"required,uuid"`
	Type_      string          `json:"type"                validate:"required"`
	CreatedOn_ time.Time       `json:"created_on"          validate:"required"`

	// can be used by callers to locate events but not persisted
	StepUUID_ flows.StepUUID `json:"-"`

	// not set by engine but can be set by callers for storage of events
	User_ *assets.UserReference `json:"_user,omitempty"`
}

// NewBaseEvent creates a new base event
func NewBaseEvent(typ string) BaseEvent {
	return BaseEvent{UUID_: flows.NewEventUUID(), Type_: typ, CreatedOn_: dates.Now()}
}

// UUID returns the UUID of this event
func (e *BaseEvent) UUID() flows.EventUUID { return e.UUID_ }

// Type returns the type of this event
func (e *BaseEvent) Type() string { return e.Type_ }

// CreatedOn returns the created on time of this event
func (e *BaseEvent) CreatedOn() time.Time { return e.CreatedOn_ }

// StepUUID returns the UUID of the step in the path where this event occurred
func (e *BaseEvent) StepUUID() flows.StepUUID { return e.StepUUID_ }

// SetStepUUID sets the UUID of the step in the path where this event occurred
func (e *BaseEvent) SetStepUUID(u flows.StepUUID) { e.StepUUID_ = u }

// SetUser can be used by callers to set the user associated with this event
func (e *BaseEvent) SetUser(u *assets.UserReference) { e.User_ = u }

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

// Read reads a single event from the given JSON
func Read(data []byte) (flows.Event, error) {
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
