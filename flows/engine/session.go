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

type session struct {
	assets  flows.Assets
	env     utils.Environment
	contact *flows.Contact

	runs       []flows.FlowRun
	runsByUUID map[flows.RunUUID]flows.FlowRun
	wait       flows.Wait
	log        []flows.LogEntry

	trigger *events.FlowTriggeredEvent
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

func (s *session) Trigger(trigger flows.Event) {
	s.trigger = trigger.(*events.FlowTriggeredEvent)
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

//------------------------------------------------------------------------------------------
// Flow execution
//------------------------------------------------------------------------------------------

// StartFlow starts the flow for the passed in contact, returning the created FlowRun
func (s *session) StartFlow(flowUUID flows.FlowUUID, parent flows.FlowRun, callerEvents []flows.Event) error {
	// TODO parent passed in by caller
	s.Trigger(events.NewFlowTriggeredEvent(flowUUID, ""))

	// off to the races
	return s.continueUntilWait(nil, noDestination, nil, callerEvents)
}

// Resume resumes a waiting session
func (s *session) Resume(callerEvents []flows.Event) error {
	// check that this session is waiting and therefore can be resumed
	if s.wait == nil {
		return utils.NewValidationErrors("only sessions with a wait can be resumed")
	}

	// figure out where (i.e. run and step) we began waiting on
	run := s.waitingRun()
	step, _, err := run.PathLocation()
	if err != nil {
		return err
	}

	// apply our caller events to this step
	for _, event := range callerEvents {
		run.ApplyEvent(step, nil, event)
	}

	// see if the wait is now satisfied and will let us resume
	if s.wait.CanResume(run, step) {
		run.SetStatus(flows.StatusActive)
		s.SetWait(nil)

		destination, err := s.resumeRun(run)
		if err != nil {
			return err
		}

		// off to the races again
		if err = s.continueUntilWait(run, destination, step, []flows.Event{}); err != nil {
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
	destination, step, err := pickNodeExit(run, node, step)
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
		if s.trigger != nil {
			// a flow has been triggered so switch execution to that flow
			flow, _ := s.Assets().GetFlow(s.trigger.FlowUUID)
			currentRun = s.CreateRun(flow, currentRun)

			// our destination is the first node in that flow
			if len(flow.Nodes()) > 0 {
				destination = flow.Nodes()[0].UUID()
			}

			// clear the trigger
			s.trigger = nil
		}

		// if we have not destination but we do have a parent, switch back to that run
		if destination == noDestination {
			currentRun.Exit(flows.StatusCompleted)

			// if we have a parent run, try to resume it
			if currentRun.Parent() != nil {
				currentRun, err = s.GetRun(currentRun.Parent().UUID())
				if destination, err = s.resumeRun(currentRun); err != nil {
					return err
				}
			} else {
				// if we have no destination and no parent, then we are truly finished here
				return nil
			}
		}

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

			// if we hit a wait, return to the caller
			if s.wait != nil {
				return nil
			}

			// mark this node as visited to prevent loops
			visited[node.UUID()] = true

			// only pass our caller events to the first node as it is responsible for handling them
			callerEvents = nil
		}
	}

	return nil
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
			if err != nil {
				run.AddError(step, err)
				run.Exit(flows.StatusErrored)
				return noDestination, step, nil
			}
		}
	}

	// a start flow action may have triggered a subflow in which case we're down on this node for now
	// and it will be resumed when the subflow finishes
	if s.trigger != nil {
		return noDestination, step, nil
	}

	// if our node has a wait before its router, we hand back to the caller
	wait := node.Wait()
	if wait != nil {
		wait.Apply(run, step)

		run.SetStatus(flows.StatusWaiting)
		run.Session().SetWait(wait)

		return noDestination, step, nil
	}

	// use our node's router to determine where to go next
	return pickNodeExit(run, node, step)
}

func pickNodeExit(run flows.FlowRun, node flows.Node, step flows.Step) (flows.NodeUUID, flows.Step, error) {
	var err error
	var exitUUID flows.ExitUUID
	var exit flows.Exit
	var exitName string
	route := flows.NoRoute

	router := node.Router()
	if router != nil {
		// we have a router, have it determine our exit
		route, err = router.PickRoute(run, node.Exits(), step)
		exitUUID = route.Exit()
	} else if len(node.Exits()) > 0 {
		// no router, pick our first exit if we have one
		exitUUID = node.Exits()[0].UUID()
	}

	step.Leave(exitUUID)

	// if we had an error routing, that's it, we are done
	if err != nil {
		run.AddError(step, err)
		run.Exit(flows.StatusErrored)
		return noDestination, step, err
	}

	// look up our actual exit
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
			err = fmt.Errorf("unable to find exit with uuid '%s'", exitUUID)
		}
	}

	// save our results if appropriate
	if router != nil && router.ResultName() != "" {
		event := events.NewSaveFlowResult(node.UUID(), router.ResultName(), route.Match(), exit.Name(), exitName)
		run.ApplyEvent(step, nil, event)
	}

	// log any error we received
	if err != nil {
		run.AddError(step, err)
		run.Exit(flows.StatusErrored)
	}

	// no exit? return no destination
	if exit == nil {
		return noDestination, step, nil
	}

	return exit.DestinationNodeUUID(), step, nil
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
