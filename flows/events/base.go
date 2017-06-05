package events

import (
	"time"

	"github.com/nyaruka/goflow/flows"
)

type BaseEvent struct {
	CreatedOn_ *time.Time     `json:"created_on"    validate:"required"`
	StepUUID_  flows.StepUUID `json:"step_uuid"`
}

func (e *BaseEvent) CreatedOn() *time.Time       { return e.CreatedOn_ }
func (e *BaseEvent) SetCreatedOn(time time.Time) { e.CreatedOn_ = &time }

func (e *BaseEvent) Step() flows.StepUUID        { return e.StepUUID_ }
func (e *BaseEvent) SetStep(step flows.StepUUID) { e.StepUUID_ = step }
