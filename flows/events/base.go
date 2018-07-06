package events

import (
	"time"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

type baseEvent struct {
	CreatedOn_  time.Time      `json:"created_on" validate:"required"`
	StepUUID_   flows.StepUUID `json:"step_uuid,omitempty" validate:"omitempty,uuid4"`
	FromCaller_ bool           `json:"-"`
}

func newBaseEvent() baseEvent {
	return baseEvent{CreatedOn_: utils.Now()}
}

func (e *baseEvent) CreatedOn() time.Time        { return e.CreatedOn_ }
func (e *baseEvent) SetCreatedOn(time time.Time) { e.CreatedOn_ = time }

func (e *baseEvent) StepUUID() flows.StepUUID            { return e.StepUUID_ }
func (e *baseEvent) SetStepUUID(stepUUID flows.StepUUID) { e.StepUUID_ = stepUUID }

func (e *baseEvent) FromCaller() bool              { return e.FromCaller_ }
func (e *baseEvent) SetFromCaller(fromCaller bool) { e.FromCaller_ = fromCaller }

type callerOnlyEvent struct{}

// AllowedOrigin determines where this event type can originate
func (e *callerOnlyEvent) AllowedOrigin() flows.EventOrigin { return flows.EventOriginCaller }

type engineOnlyEvent struct{}

// AllowedOrigin determines where this event type can originate
func (e *engineOnlyEvent) AllowedOrigin() flows.EventOrigin { return flows.EventOriginEngine }

// Validate validates our event is valid and has all the assets it needs. We assume engine generated events are valid.
func (e *engineOnlyEvent) Validate(assets flows.SessionAssets) error {
	return nil
}

type callerOrEngineEvent struct{}

// AllowedOrigin determines where this event type can originate
func (e *callerOrEngineEvent) AllowedOrigin() flows.EventOrigin {
	return flows.EventOriginCaller | flows.EventOriginEngine
}
