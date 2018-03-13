package engine

import (
	"encoding/json"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

type logEntry struct {
	step   flows.Step
	action flows.Action
	event  flows.Event
}

// NewLogEntry creates a new event log entry
func NewLogEntry(step flows.Step, action flows.Action, event flows.Event) flows.LogEntry {
	return &logEntry{step: step, action: action, event: event}
}

func (s *logEntry) Step() flows.Step     { return s.step }
func (s *logEntry) Action() flows.Action { return s.action }
func (s *logEntry) Event() flows.Event   { return s.event }

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type logEntryEnvelope struct {
	StepUUID   flows.StepUUID       `json:"step_uuid,omitempty"`
	ActionUUID flows.ActionUUID     `json:"action_uuid,omitempty"`
	Event      *utils.TypedEnvelope `json:"event" validate:"required"`
}

// MarshalJSON marshals this log entry into JSON
func (s *logEntry) MarshalJSON() ([]byte, error) {
	var se logEntryEnvelope

	if s.step != nil {
		se.StepUUID = s.step.UUID()
	}
	if s.action != nil {
		se.ActionUUID = s.action.UUID()
	}

	eventData, err := json.Marshal(s.event)
	if err != nil {
		return nil, err
	}

	se.Event = &utils.TypedEnvelope{Type: s.event.Type(), Data: eventData}

	return json.Marshal(se)
}
