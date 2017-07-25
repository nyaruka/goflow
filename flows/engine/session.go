package engine

import (
	"encoding/json"
	"fmt"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/flows/runs"
	"github.com/nyaruka/goflow/utils"
)

type session struct {
	assets  flows.Assets
	env     utils.Environment
	contact *flows.Contact

	runs       []flows.FlowRun
	runsByUUID map[flows.RunUUID]flows.FlowRun
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

func (s *session) addRun(run flows.FlowRun) {
	s.runs = append(s.runs, run)
	s.runsByUUID[run.UUID()] = run
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

// ResumeFlow resumes our session from the last step
func (s *session) Resume(callerEvents []flows.Event) error {
	// find the active run
	run := s.ActiveRun()
	if run == nil {
		return utils.NewValidationError("session: no active run to resume")
	}

	return s.resumeRun(run, callerEvents)
}

func (s *session) resumeRun(run flows.FlowRun, callerEvents []flows.Event) error {
	// no steps to resume from, nothing to do, return
	if len(run.Path()) == 0 {
		return nil
	}

	// grab the last step
	step := run.Path()[len(run.Path())-1]

	// and the last node
	node := run.Flow().GetNode(step.NodeUUID())
	if node == nil {
		err := fmt.Errorf("cannot resume at node '%s' that no longer exists", step.NodeUUID())
		run.AddError(step, err)
		run.Exit(flows.StatusErrored)
		return nil
	}

	// the first event resumes the wait
	destination, step, err := resumeNode(run, node, step, callerEvents)
	if err != nil {
		return err
	}

	err = continueRunUntilWait(run, destination, step, nil)
	if err != nil {
		return err
	}

	// if we ran to completion and have a parent, resume that flow
	if run.Parent() != nil && run.IsComplete() {
		event := events.NewFlowExitedEvent(run)
		parentRun, err := run.Session().GetRun(run.Parent().UUID())
		if err != nil {
			run.AddError(step, err)
			run.Exit(flows.StatusErrored)
			return nil
		}
		return s.resumeRun(parentRun, []flows.Event{event})
	}

	return nil
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type sessionEnvelope struct {
	Environment json.RawMessage   `json:"environment"`
	Contact     json.RawMessage   `json:"contact"`
	Runs        []json.RawMessage `json:"runs"`
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
		return nil, utils.NewValidationError(err.Error())
	}

	err = utils.ValidateAll(s)
	return s, err
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

	return json.Marshal(envelope)
}
