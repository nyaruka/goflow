package events

import (
	"encoding/json"
	"time"

	"github.com/nyaruka/gocommon/dates"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"

	"github.com/pkg/errors"
)

var registeredTypes = map[string](func() flows.Event){}

// registers a new type of event
func registerType(name string, initFunc func() flows.Event) {
	registeredTypes[name] = initFunc
}

// base of all event types
type baseEvent struct {
	Type_      string         `json:"type" validate:"required"`
	CreatedOn_ time.Time      `json:"created_on" validate:"required"`
	StepUUID_  flows.StepUUID `json:"step_uuid,omitempty" validate:"omitempty,uuid4"`
}

// creates a new base event
func newBaseEvent(typeName string) baseEvent {
	return baseEvent{Type_: typeName, CreatedOn_: dates.Now()}
}

// Type returns the type of this event
func (e *baseEvent) Type() string { return e.Type_ }

// CreatedOn returns the created on time of this event
func (e *baseEvent) CreatedOn() time.Time { return e.CreatedOn_ }

// StepUUID returns the UUID of the step in the path where this event occurred
func (e *baseEvent) StepUUID() flows.StepUUID { return e.StepUUID_ }

// SetStepUUID sets the UUID of the step in the path where this event occurred
func (e *baseEvent) SetStepUUID(stepUUID flows.StepUUID) { e.StepUUID_ = stepUUID }

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

// ReadEvent reads a single event from the given JSON
func ReadEvent(data json.RawMessage) (flows.Event, error) {
	typeName, err := utils.ReadTypeFromJSON(data)
	if err != nil {
		return nil, err
	}

	f := registeredTypes[typeName]
	if f == nil {
		return nil, errors.Errorf("unknown type: '%s'", typeName)
	}

	event := f()
	return event, utils.UnmarshalAndValidate(data, event)
}
