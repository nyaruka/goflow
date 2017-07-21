package runs

import (
	"encoding/json"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

type session struct {
	runs   []flows.FlowRun
	events []flows.LogEntry
}

func newSession() *session {
	session := session{}
	return &session
}

func (s *session) AddRun(run flows.FlowRun) {
	// check if we already have this run
	for _, r := range s.runs {
		if r.UUID() == run.UUID() {
			return
		}
	}
	s.runs = append(s.runs, run)
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
func ReadSession(env flows.FlowEnvironment, data json.RawMessage) (flows.Session, error) {
	s := &session{}
	var se sessionEnvelope
	var err error

	err = json.Unmarshal(data, &se)
	if err != nil {
		return nil, err
	}

	s.runs = make([]flows.FlowRun, len(se.Runs))
	for i := range se.Runs {
		s.runs[i], err = ReadRun(env, se.Runs[i])
		if err != nil {
			return nil, err
		}
		s.runs[i].SetSession(s)
	}

	err = utils.ValidateAll(s)
	return s, err
}

func (s *session) MarshalJSON() ([]byte, error) {
	var se sessionEnvelope
	var err error
	se.Runs = make([]json.RawMessage, len(s.runs))
	for i := range s.runs {
		se.Runs[i], err = json.Marshal(s.runs[i])
		if err != nil {
			return nil, err
		}
	}

	return json.Marshal(se)
}
