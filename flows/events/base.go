package events

import (
	"time"

	"github.com/nyaruka/goflow/flows"
)

type BaseEvent struct {
	CreatedOn_  time.Time      `json:"created_on"    validate:"required"`
	StepUUID_   flows.StepUUID `json:"step_uuid"`
	FromCaller_ bool           `json:"-"`
}

func NewBaseEvent() BaseEvent {
	return BaseEvent{CreatedOn_: time.Now().UTC()}
}

func (e *BaseEvent) CreatedOn() time.Time        { return e.CreatedOn_ }
func (e *BaseEvent) SetCreatedOn(time time.Time) { e.CreatedOn_ = time }

func (e *BaseEvent) Step() flows.StepUUID        { return e.StepUUID_ }
func (e *BaseEvent) SetStep(step flows.StepUUID) { e.StepUUID_ = step }

func (e *BaseEvent) FromCaller() bool              { return e.FromCaller_ }
func (e *BaseEvent) SetFromCaller(fromCaller bool) { e.FromCaller_ = fromCaller }
