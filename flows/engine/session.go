package engine

import (
	"encoding/json"
	"fmt"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/runs"
	"github.com/nyaruka/goflow/flows/waits"
	"github.com/nyaruka/goflow/utils"
)

type session struct {
	assets  flows.Assets
	env     utils.Environment
	contact *flows.Contact

	runs       []flows.FlowRun
	runsByUUID map[flows.RunUUID]flows.FlowRun
	wait       flows.Wait
	log        []flows.LogEntry
}

// NewSession creates a new session
func NewSession(assets flows.Assets) flows.Session {
	return &session{
		env:        utils.NewDefaultEnvironment(),
		assets:     assets,
		runsByUUID: make(map[flows.RunUUID]flows.FlowRun),
	}
}

func (s *session) Assets() flows.Assets                 { return s.assets }
func (s *session) Environment() utils.Environment       { return s.env }
func (s *session) SetEnvironment(env utils.Environment) { s.env = env }
func (s *session) Contact() *flows.Contact              { return s.contact }
func (s *session) SetContact(contact *flows.Contact)    { s.contact = contact }

func (s *session) CreateRun(flow flows.Flow, parent flows.FlowRun) flows.FlowRun {
	run := runs.NewRun(s, flow, s.contact, parent)
	s.addRun(run)
	return run
}

func (s *session) Runs() []flows.FlowRun { return s.runs }
func (s *session) GetRun(uuid flows.RunUUID) (flows.FlowRun, error) {
	run, exists := s.runsByUUID[uuid]
	if exists {
		return run, nil
	}
	return nil, fmt.Errorf("unable to find run with UUID: %s", uuid)
}

func (s *session) addRun(run flows.FlowRun) {
	s.runs = append(s.runs, run)
	s.runsByUUID[run.UUID()] = run
}

func (s *session) Wait() flows.Wait        { return s.wait }
func (s *session) SetWait(wait flows.Wait) { s.wait = wait }

// looks through this session's run for the one that is waiting
func (s *session) waitingRun() flows.FlowRun {
	for _, run := range s.runs {
		if run.Status() == flows.StatusWaiting {
			return run
		}
	}
	return nil
}

func (s *session) LogEvent(step flows.Step, action flows.Action, event flows.Event) {
	s.log = append(s.log, NewLogEntry(step, action, event))
}
func (s *session) Log() []flows.LogEntry { return s.log }
func (s *session) ClearLog()             { s.log = nil }

// StartFlow starts the flow for the passed in contact, returning the created FlowRun
func (s *session) StartFlow(flowUUID flows.FlowUUID, parent flows.FlowRun, callerEvents []flows.Event) error {
	flow, err := s.assets.GetFlow(flowUUID)
	if err != nil {
		return err
	}

	// create our new run
	run := s.CreateRun(flow, parent)

	// no first node, nothing to do (valid but weird)
	if len(flow.Nodes()) == 0 {
		run.Exit(flows.StatusCompleted)
		return nil
	}

	// off to the races
	return continueRunUntilWait(run, flow.Nodes()[0].UUID(), nil, callerEvents)
}

// Resume resumes a waiting session from its last step
func (s *session) Resume(callerEvents []flows.Event) error {
	// check that this session is waiting and therefore can be resumed
	if s.Wait() == nil {
		return utils.NewValidationErrors("only waiting sessions can be resumed")
	}

	// figure out which run and step we began waiting on
	run := s.waitingRun()
	step, _, err := run.PathLocation()
	if err != nil {
		return err
	}

	// apply our caller events
	for _, event := range callerEvents {
		run.ApplyEvent(step, nil, event)
	}

	// if our wait is now satified, resume the run
	if s.Wait().CanResume(run, step) {
		s.SetWait(nil)
		run.SetStatus(flows.StatusActive)

		return s.resumeRun(run)
	}

	// otherwise return to the caller
	return nil
}

func (s *session) resumeRun(run flows.FlowRun) error {
	step, node, err := run.PathLocation()
	if err != nil {
		return err
	}

	// see if this node can now pick a destination
	destination, step, err := pickNodeExit(run, node, step)
	if err != nil {
		return err
	}

	err = continueRunUntilWait(run, destination, step, nil)
	if err != nil {
		return err
	}

	// if we ran to completion and have a parent, resume that flow
	if run.Parent() != nil && run.IsComplete() {
		parentRun, err := run.Session().GetRun(run.Parent().UUID())
		if err != nil {
			return err
		}
		return s.resumeRun(parentRun)
	}

	return nil
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type sessionEnvelope struct {
	Environment json.RawMessage      `json:"environment"`
	Contact     json.RawMessage      `json:"contact"`
	Runs        []json.RawMessage    `json:"runs"`
	Wait        *utils.TypedEnvelope `json:"wait"`
}

// ReadSession decodes a session from the passed in JSON
func ReadSession(assets flows.Assets, data json.RawMessage) (flows.Session, error) {
	s := NewSession(assets).(*session)
	var envelope sessionEnvelope
	var err error

	err = json.Unmarshal(data, &envelope)
	if err != nil {
		return nil, err
	}
	if err = utils.Validate(s); err != nil {
		return nil, err
	}

	// read our environment
	s.env, err = utils.ReadEnvironment(envelope.Environment)
	if err != nil {
		return nil, err
	}

	// read our contact
	s.contact, err = flows.ReadContact(assets, envelope.Contact)
	if err != nil {
		return nil, err
	}

	// read each of our runs
	for i := range envelope.Runs {
		run, err := runs.ReadRun(s, envelope.Runs[i])
		if err != nil {
			return nil, err
		}
		s.addRun(run)
	}

	// once all runs are read, we can resolve references between runs
	err = runs.ResolveReferences(s, s.Runs())
	if err != nil {
		return nil, utils.NewValidationErrors(err.Error())
	}

	// and our wait
	if envelope.Wait != nil {
		s.wait, err = waits.WaitFromEnvelope(envelope.Wait)
		if err != nil {
			return nil, err
		}
	}

	return s, nil
}

func (s *session) MarshalJSON() ([]byte, error) {
	var envelope sessionEnvelope
	var err error

	envelope.Environment, err = json.Marshal(s.env)
	if err != nil {
		return nil, err
	}

	envelope.Contact, err = json.Marshal(s.contact)
	if err != nil {
		return nil, err
	}

	envelope.Runs = make([]json.RawMessage, len(s.runs))
	for i := range s.runs {
		envelope.Runs[i], err = json.Marshal(s.runs[i])
		if err != nil {
			return nil, err
		}
	}

	if s.wait != nil {
		if envelope.Wait, err = utils.EnvelopeFromTyped(s.wait); err != nil {
			return nil, err
		}
	}

	return json.Marshal(envelope)
}
