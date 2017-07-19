package events

import (
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
