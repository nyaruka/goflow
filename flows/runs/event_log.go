package runs

import (
	"encoding/json"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

type logEntry struct {
	stepUUID   flows.StepUUID
	actionUUID flows.ActionUUID
	event      flows.Event
}

// NewLogEntry creates a new event log entry
func NewLogEntry(step flows.Step, action flows.Action, event flows.Event) flows.LogEntry {
	var actionUUID flows.ActionUUID
	if action != nil {
		actionUUID = action.UUID()
	}

	return &logEntry{stepUUID: step.UUID(), actionUUID: actionUUID, event: event}
}

func (s *logEntry) StepUUID() flows.StepUUID     { return s.stepUUID }
func (s *logEntry) ActionUUID() flows.ActionUUID { return s.actionUUID }
func (s *logEntry) Event() flows.Event           { return s.event }

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type logEntryEnvelope struct {
	StepUUID   flows.StepUUID       `json:"step_uuid"   validate:"required"`
	ActionUUID *flows.ActionUUID    `json:"action_uuid"`
	Event      *utils.TypedEnvelope `json:"event"       validate:"required"`
}

func (s *logEntry) MarshalJSON() ([]byte, error) {
	var se logEntryEnvelope

	se.StepUUID = s.stepUUID
	if s.actionUUID != "" {
		se.ActionUUID = &s.actionUUID
	}

	eventData, err := json.Marshal(s.event)
	if err != nil {
		return nil, err
	}

	se.Event = &utils.TypedEnvelope{Type: s.event.Type(), Data: eventData}

	return json.Marshal(se)
}
