package engine

import (
	"encoding/json"
	"fmt"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/flows/runs"
	"github.com/nyaruka/goflow/flows/waits"
	"github.com/nyaruka/goflow/utils"
)

// used to trigger a new flow run in the event loop
type flowTrigger struct {
	flow      flows.Flow
	parentRun flows.FlowRun
}

type session struct {
	assets  flows.Assets
	env     utils.Environment
	contact *flows.Contact

	runs       []flows.FlowRun
	runsByUUID map[flows.RunUUID]flows.FlowRun
	status     flows.SessionStatus
	wait       flows.Wait
	log        []flows.LogEntry

	flowTrigger *flowTrigger
}

// NewSession creates a new session
func NewSession(assets flows.Assets) flows.Session {
	return &session{
		env:        utils.NewDefaultEnvironment(),
		assets:     assets,
		status:     flows.SessionStatusActive,
		runsByUUID: make(map[flows.RunUUID]flows.FlowRun),
	}
}

func (s *session) Assets() flows.Assets                 { return s.assets }
func (s *session) Environment() utils.Environment       { return s.env }
func (s *session) SetEnvironment(env utils.Environment) { s.env = env }
func (s *session) Contact() *flows.Contact              { return s.contact }
func (s *session) SetContact(contact *flows.Contact)    { s.contact = contact }

func (s *session) Trigger(flow flows.Flow, parentRun flows.FlowRun) {
	s.flowTrigger = &flowTrigger{flow: flow, parentRun: parentRun}
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

func (s *session) Status() flows.SessionStatus { return s.status }
func (s *session) Wait() flows.Wait            { return s.wait }

// looks through this session's run for the one that is waiting
func (s *session) waitingRun() flows.FlowRun {
	for _, run := range s.runs {
		if run.Status() == flows.RunStatusWaiting {
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

//------------------------------------------------------------------------------------------
// Flow execution
//------------------------------------------------------------------------------------------

// StartFlow beings execution of the given flow in this session
func (s *session) StartFlow(flowUUID flows.FlowUUID, callerEvents []flows.Event) error {
	flow, _ := s.Assets().GetFlow(flowUUID)

	s.Trigger(flow, nil)

	// off to the races...
	return s.continueUntilWait(nil, noDestination, nil, callerEvents)
}

// Resume resumes a waiting session
func (s *session) Resume(callerEvents []flows.Event) error {
	if s.status != flows.SessionStatusWaiting {
		return utils.NewValidationErrors("only waiting sessions can be resumed")
	}

	// figure out where (i.e. run and step) we began waiting on
	waitingRun := s.waitingRun()
	step, _, err := waitingRun.PathLocation()
	if err != nil {
		return err
	}

	// apply our caller events to this step
	for _, event := range callerEvents {
		waitingRun.ApplyEvent(step, nil, event)
	}

	// see if the wait is now satisfied and will let us resume
	if s.wait.CanResume(waitingRun, step) {
		s.status = flows.SessionStatusActive
		waitingRun.SetStatus(flows.RunStatusActive)

		destination, err := s.resumeRun(waitingRun)
		if err != nil {
			return err
		}

		// off to the races again...
		if err = s.continueUntilWait(waitingRun, destination, step, []flows.Event{}); err != nil {
			return err
		}
	}

	return nil
}

// resumes a run that may have been waiting or a parent paused for a child subflow
func (s *session) resumeRun(run flows.FlowRun) (flows.NodeUUID, error) {
	step, node, err := run.PathLocation()
	if err != nil {
		return noDestination, err
	}

	// see if this node can now pick a destination
	destination, step, err := s.pickNodeExit(run, node, step)
	if err != nil {
		return noDestination, err
	}

	return destination, nil
}

// the main flow execution loop
func (s *session) continueUntilWait(currentRun flows.FlowRun, destination flows.NodeUUID, step flows.Step, callerEvents []flows.Event) (err error) {
	// track the node UUIDs we've visited so we can watch out for loops
	visited := make(visitedMap)

	for {
		// if we have a flow trigger handle that first
		if s.flowTrigger != nil {
			// create a new run for it
			flow := s.flowTrigger.flow
			currentRun = runs.NewRun(s, s.flowTrigger.flow, s.contact, currentRun)
			s.addRun(currentRun)

			// our destination is the first node in that flow... if such a node exists
			if len(flow.Nodes()) > 0 {
				destination = flow.Nodes()[0].UUID()
			} else {
				destination = noDestination
			}

			// clear the trigger
			s.flowTrigger = nil
		}

		// if we have no destination but we do have a parent, switch back to that
		if destination == noDestination {
			currentRun.Exit(flows.RunStatusCompleted)

			// if we have a parent run, try to resume it
			if currentRun.Parent() != nil {
				currentRun, err = s.GetRun(currentRun.Parent().UUID())
				if destination, err = s.resumeRun(currentRun); err != nil {
					return err
				}
			} else {
				// if we have no destination and no parent, then we are truly finished here
				s.status = flows.SessionStatusCompleted
				return nil
			}
		}

		// if we now have a destination, go there
		if destination != noDestination {
			// this is a loop, we log it and stop execution
			if visited[destination] {
				return fmt.Errorf("flow loop detected, stopping execution before entering %s", destination)
			}

			node := currentRun.Flow().GetNode(destination)
			if node == nil {
				return fmt.Errorf("unable to find destination node %s in flow %s", destination, currentRun.Flow().UUID())
			}

			destination, step, err = s.visitNode(currentRun, node, callerEvents)
			if err != nil {
				return err
			}

			// if our current run has errored, mark the session as errored and return to caller
			if currentRun.Status() == flows.RunStatusErrored {
				s.setErrored()
				return nil
			}

			// if we hit a wait, also return to the caller
			if s.status == flows.SessionStatusWaiting {
				return nil
			}

			// mark this node as visited to prevent loops
			visited[node.UUID()] = true

			// only pass our caller events to the first node as it is responsible for handling them
			callerEvents = nil
		}
	}
}

// visits the given node, creating a step in our current run path
func (s *session) visitNode(run flows.FlowRun, node flows.Node, callerEvents []flows.Event) (flows.NodeUUID, flows.Step, error) {
	step := run.CreateStep(node)

	// apply any caller events
	for _, event := range callerEvents {
		run.ApplyEvent(step, nil, event)
	}

	// execute our node's actions
	if node.Actions() != nil {
		for _, action := range node.Actions() {
			err := action.Execute(run, step)

			if err != nil || run.Status() == flows.RunStatusErrored {
				return noDestination, step, err
			}
		}
	}

	// a start flow action may have triggered a subflow in which case we're done on this node for now
	// and it will be resumed when the subflow finishes
	if s.flowTrigger != nil {
		return noDestination, step, nil
	}

	// if our node has a wait before its router, we hand back to the caller
	wait := node.Wait()
	if wait != nil {
		wait.Apply(run, step)

		s.wait = wait
		s.status = flows.SessionStatusWaiting
		run.SetStatus(flows.RunStatusWaiting)

		return noDestination, step, nil
	}

	// use our node's router to determine where to go next
	return s.pickNodeExit(run, node, step)
}

// picks the exit to use on the given node
func (s *session) pickNodeExit(run flows.FlowRun, node flows.Node, step flows.Step) (flows.NodeUUID, flows.Step, error) {
	var err error

	route := flows.NoRoute
	router := node.Router()

	// we have a router, have it determine our exit
	var exitUUID flows.ExitUUID
	if router != nil {
		if route, err = router.PickRoute(run, node.Exits(), step); err != nil {
			return noDestination, nil, err
		}
		exitUUID = route.Exit()
	} else if len(node.Exits()) > 0 {
		// no router, pick our first exit if we have one
		exitUUID = node.Exits()[0].UUID()
	}

	step.Leave(exitUUID)

	// look up our actual exit and localized name
	var exit flows.Exit
	var exitName string

	if exitUUID != "" {
		// find our exit
		for _, e := range node.Exits() {
			if e.UUID() == exitUUID {
				localizedName := run.GetText(flows.UUID(exitUUID), "name", e.Name())
				if localizedName != e.Name() {
					exitName = localizedName
				}
				exit = e
				break
			}
		}
		if exit == nil {
			return noDestination, nil, fmt.Errorf("unable to find exit with uuid '%s'", exitUUID)
		}
	}

	// save our results if appropriate
	if router != nil && router.ResultName() != "" {
		event := events.NewSaveFlowResult(node.UUID(), router.ResultName(), route.Match(), exit.Name(), exitName)
		run.ApplyEvent(step, nil, event)
	}

	// no exit? return no destination
	if exit == nil {
		return noDestination, step, nil
	}

	return exit.DestinationNodeUUID(), step, nil
}

func (s *session) setErrored() {
	// mark any active or waiting runs as errored
	for _, run := range s.runs {
		if run.Status() == flows.RunStatusActive || run.Status() == flows.RunStatusWaiting {
			run.SetStatus(flows.RunStatusErrored)
		}
	}

	s.status = flows.SessionStatusErrored
	s.wait = nil
	s.flowTrigger = nil
}

type visitedMap map[flows.NodeUUID]bool

const noDestination = flows.NodeUUID("")

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type sessionEnvelope struct {
	Environment json.RawMessage      `json:"environment"`
	Contact     json.RawMessage      `json:"contact"`
	Runs        []json.RawMessage    `json:"runs"`
	Status      flows.SessionStatus  `json:"status"`
	Wait        *utils.TypedEnvelope `json:"wait,omitempty"`
}

// ReadSession decodes a session from the passed in JSON
func ReadSession(assets flows.Assets, data json.RawMessage) (flows.Session, error) {
	s := NewSession(assets).(*session)
	var envelope sessionEnvelope
	var err error

	if err = json.Unmarshal(data, &envelope); err != nil {
		return nil, err
	}
	if err = utils.Validate(s); err != nil {
		return nil, err
	}

	s.status = envelope.Status

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

	envelope.Status = s.status
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
