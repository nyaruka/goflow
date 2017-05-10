package events

import (
	"time"

	"github.com/nyaruka/goflow/flows"
)

type BaseEvent struct {
	Run_       flows.RunUUID `json:"run"`
	CreatedOn_ *time.Time    `json:"created_on"`
}

func (e *BaseEvent) CreatedOn() *time.Time       { return e.CreatedOn_ }
func (e *BaseEvent) SetCreatedOn(time time.Time) { e.CreatedOn_ = &time }

func (e *BaseEvent) Run() flows.RunUUID       { return e.Run_ }
func (e *BaseEvent) SetRun(run flows.RunUUID) { e.Run_ = run }
