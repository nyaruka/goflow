package runs

import (
	"encoding/json"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

type session struct {
	env    flows.SessionEnvironment
	runs   []flows.FlowRun
	events []flows.LogEntry
}

// NewSession creates a new session
func NewSession(env flows.SessionEnvironment) *session {
	return &session{env: env}
}

func (s *session) Environment() flows.SessionEnvironment { return s.env }

func (s *session) CreateRun(flow flows.Flow, contact *flows.Contact, parent flows.FlowRun) flows.FlowRun {
	run := NewRun(s, flow, contact, parent)
	s.runs = append(s.runs, run)
	return run
}

func (s *session) Runs() []flows.FlowRun { return s.runs }

func (s *session) ActiveRun() flows.FlowRun {
	var active flows.FlowRun
	mostRecent := utils.ZeroTime

	for _, run := range s.runs {
		// We are complete, therefore can't be active
		if run.IsComplete() {
			continue
		}

		// We have a child, and it isn't complete, we can't be active
		if run.Child() != nil && run.Child().Status() == flows.StatusActive {
			continue
		}

		// this is more recent than our most recent flow
		if run.ModifiedOn().After(mostRecent) {
			active = run
			mostRecent = run.ModifiedOn()
		}
	}
	return active
}

func (s *session) LogEvent(step flows.Step, action flows.Action, event flows.Event) {
	s.events = append(s.events, NewLogEntry(step, action, event))
}
func (s *session) Log() []flows.LogEntry { return s.events }
func (s *session) ClearLog()             { s.events = nil }

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type sessionEnvelope struct {
	Runs []json.RawMessage `json:"runs"`
}

// ReadSession decodes a session from the passed in JSON
func ReadSession(env flows.SessionEnvironment, data json.RawMessage) (flows.Session, error) {
	s := &session{}
	var envelope sessionEnvelope
	var err error

	err = json.Unmarshal(data, &envelope)
	if err != nil {
		return nil, err
	}

	s.runs = make([]flows.FlowRun, len(envelope.Runs))
	for i := range envelope.Runs {
		s.runs[i], err = ReadRun(s, envelope.Runs[i])
		if err != nil {
			return nil, err
		}
	}

	// once all runs are read, we can resolve references between runs
	for _, run := range s.runs {
		err = run.(*flowRun).ResolveReferences(env)
		if err != nil {
			return nil, utils.NewValidationError(err.Error())
		}
	}

	err = utils.ValidateAll(s)
	return s, err
}

func (s *session) MarshalJSON() ([]byte, error) {
	var envelope sessionEnvelope
	var err error
	envelope.Runs = make([]json.RawMessage, len(s.runs))
	for i := range s.runs {
		envelope.Runs[i], err = json.Marshal(s.runs[i])
		if err != nil {
			return nil, err
		}
	}

	return json.Marshal(envelope)
}
