package events

import (
	"github.com/nyaruka/goflow/flows"
	"time"
)

type BaseEvent struct {
	CreatedOn_  time.Time `json:"created_on"    validate:"required"`
	FromCaller_ bool      `json:"-"`
}

func NewBaseEvent() BaseEvent {
	return BaseEvent{CreatedOn_: time.Now().UTC()}
}

func (e *BaseEvent) CreatedOn() time.Time        { return e.CreatedOn_ }
func (e *BaseEvent) SetCreatedOn(time time.Time) { e.CreatedOn_ = time }

func (e *BaseEvent) FromCaller() bool              { return e.FromCaller_ }
func (e *BaseEvent) SetFromCaller(fromCaller bool) { e.FromCaller_ = fromCaller }

type callerOnlyEvent struct{}

// AllowedOrigin determines where this event type can originate
func (e *callerOnlyEvent) AllowedOrigin() flows.EventOrigin { return flows.EventOriginCaller }

type engineOnlyEvent struct{}

// AllowedOrigin determines where this event type can originate
func (e *engineOnlyEvent) AllowedOrigin() flows.EventOrigin { return flows.EventOriginEngine }

type callerOrEngineEvent struct{}

// AllowedOrigin determines where this event type can originate
func (e *callerOrEngineEvent) AllowedOrigin() flows.EventOrigin {
	return flows.EventOriginCaller | flows.EventOriginEngine
}
