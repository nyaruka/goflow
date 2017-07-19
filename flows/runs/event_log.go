package runs

import (
	"encoding/json"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

type eventLogEntry struct {
	stepUUID   flows.StepUUID
	actionUUID flows.ActionUUID
	event      flows.Event
}

// NewEventLogEntry creates a new event log entry
func NewEventLogEntry(step flows.Step, action flows.Action, event flows.Event) flows.EventLogEntry {
	var actionUUID flows.ActionUUID
	if action != nil {
		actionUUID = action.UUID()
	}

	return &eventLogEntry{stepUUID: step.UUID(), actionUUID: actionUUID, event: event}
}

func (s *eventLogEntry) StepUUID() flows.StepUUID     { return s.stepUUID }
func (s *eventLogEntry) ActionUUID() flows.ActionUUID { return s.actionUUID }
func (s *eventLogEntry) Event() flows.Event           { return s.event }

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type eventLogEntryEnvelope struct {
	StepUUID   flows.StepUUID       `json:"step_uuid"   validate:"required"`
	ActionUUID *flows.ActionUUID    `json:"action_uuid"`
	Event      *utils.TypedEnvelope `json:"event"       validate:"required"`
}

func (s *eventLogEntry) MarshalJSON() ([]byte, error) {
	var se eventLogEntryEnvelope

	se.StepUUID = s.stepUUID
	if s.actionUUID != "" {
		se.ActionUUID = &s.actionUUID
	}

	eventData, err := json.Marshal(s.event)
	if err != nil {
		return nil, err
	}

	se.Event = &utils.TypedEnvelope{s.event.Type(), eventData}

	return json.Marshal(se)
}
