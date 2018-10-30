package engine

import (
	"encoding/json"
	"fmt"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/flows/inputs"
	"github.com/nyaruka/goflow/flows/runs"
	"github.com/nyaruka/goflow/flows/triggers"
	"github.com/nyaruka/goflow/flows/waits"
	"github.com/nyaruka/goflow/utils"

	log "github.com/sirupsen/logrus"
)

// used to spawn a new run or sub-flow in the event loop
type pushedFlow struct {
	flow      flows.Flow
	parentRun flows.FlowRun
	terminal  bool
}

type session struct {
	assets flows.SessionAssets

	// state which is maintained between engine calls
	env     utils.Environment
	trigger flows.Trigger
	contact *flows.Contact
	runs    []flows.FlowRun
	status  flows.SessionStatus
	wait    flows.Wait
	input   flows.Input

	// state which is temporary to each call
	runsByUUID map[flows.RunUUID]flows.FlowRun
	pushedFlow *pushedFlow
	flowStack  *flowStack
	newEvents  []flows.Event
	parentRun  flows.RunSummary

	engineConfig flows.EngineConfig
	httpClient   *utils.HTTPClient
}

// NewSession creates a new session
func NewSession(assets flows.SessionAssets, engineConfig flows.EngineConfig, httpClient *utils.HTTPClient) flows.Session {
	return &session{
		env:          utils.NewDefaultEnvironment(),
		assets:       assets,
		status:       flows.SessionStatusActive,
		newEvents:    []flows.Event{},
		runsByUUID:   make(map[flows.RunUUID]flows.FlowRun),
		flowStack:    newFlowStack(),
		engineConfig: engineConfig,
		httpClient:   httpClient,
	}
}

func (s *session) Assets() flows.SessionAssets { return s.assets }
func (s *session) Trigger() flows.Trigger      { return s.trigger }

func (s *session) Environment() utils.Environment       { return s.env }
func (s *session) SetEnvironment(env utils.Environment) { s.env = env }

func (s *session) Contact() *flows.Contact           { return s.contact }
func (s *session) SetContact(contact *flows.Contact) { s.contact = contact }

func (s *session) Input() flows.Input         { return s.input }
func (s *session) SetInput(input flows.Input) { s.input = input }

func (s *session) CanEnterFlow(flow flows.Flow) bool {
	return !s.flowStack.hasVisitedFlowSinceResume(flow.UUID())
}

func (s *session) PushFlow(flow flows.Flow, parentRun flows.FlowRun, terminal bool) {
	s.pushedFlow = &pushedFlow{flow: flow, parentRun: parentRun, terminal: terminal}
}

func (s *session) Runs() []flows.FlowRun { return s.runs }
func (s *session) GetRun(uuid flows.RunUUID) (flows.FlowRun, error) {
	run, exists := s.runsByUUID[uuid]
	if exists {
		return run, nil
	}
	return nil, fmt.Errorf("unable to find run with UUID '%s'", uuid)
}

func (s *session) addRun(run flows.FlowRun) {
	s.runs = append(s.runs, run)
	s.runsByUUID[run.UUID()] = run
}

func (s *session) GetCurrentChild(run flows.FlowRun) flows.FlowRun {
	// the current child of a run, is the last added run which has that run as its parent
	for r := len(s.runs) - 1; r >= 0; r-- {
		if s.runs[r].ParentInSession() == run {
			return s.runs[r]
		}
	}
	return nil
}

// ParentRun gets the parent run of this session if it was started by a flow action
func (s *session) ParentRun() flows.RunSummary {
	return s.parentRun
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

func (s *session) LogEvent(event flows.Event) {
	s.newEvents = append(s.newEvents, event)
}

func (s *session) EngineConfig() flows.EngineConfig { return s.engineConfig }
func (s *session) HTTPClient() *utils.HTTPClient    { return s.httpClient }

//------------------------------------------------------------------------------------------
// Flow execution
//------------------------------------------------------------------------------------------

// Start initializes this session with the given trigger and runs the flow to the first wait
func (s *session) Start(trigger flows.Trigger) ([]flows.Event, error) {
	s.trigger = trigger

	if err := s.prepareForSprint(); err != nil {
		return s.newEvents, err
	}

	if err := s.trigger.Initialize(s); err != nil {
		return s.newEvents, err
	}

	// off to the races...
	return s.newEvents, s.continueUntilWait(nil, noDestination, nil, trigger)
}

// Resume tries to resume a waiting session
func (s *session) Resume(resume flows.Resume) ([]flows.Event, error) {
	if err := s.prepareForSprint(); err != nil {
		return s.newEvents, err
	}

	if s.status != flows.SessionStatusWaiting {
		return s.newEvents, fmt.Errorf("only waiting sessions can be resumed")
	}

	waitingRun := s.waitingRun()
	if waitingRun == nil {
		return s.newEvents, fmt.Errorf("session doesn't contain any runs which are waiting")
	}

	// check flow is valid and has everything it needs to run
	if err := waitingRun.Flow().Validate(s.Assets(), flows.NewValidationContext()); err != nil {
		return s.newEvents, fmt.Errorf("validation failed for flow[uuid=%s]: %s", waitingRun.Flow().UUID(), err)
	}

	if err := s.tryToResume(waitingRun, resume); err != nil {
		// if we got an error, add it to the log and shut everything down
		for _, run := range s.runs {
			run.Exit(flows.RunStatusErrored)
		}
		s.status = flows.SessionStatusErrored
		s.LogEvent(events.NewErrorEvent(err))
	}

	return s.newEvents, nil
}

// prepares the session for starting/resuming
func (s *session) prepareForSprint() error {
	// clear the new events log
	s.newEvents = make([]flows.Event, 0)

	if s.parentRun == nil {
		// if we have a trigger with a parent run, load that
		triggerWithRun, hasRun := s.trigger.(flows.TriggerWithRun)
		if hasRun {
			run, err := runs.ReadRunSummary(s.Assets(), triggerWithRun.RunSummary())
			if err != nil {
				return fmt.Errorf("error reading parent run from trigger: %s", err)
			}
			s.parentRun = run
		}
	}
	return nil
}

// Resume resumes a waiting session
func (s *session) tryToResume(waitingRun flows.FlowRun, resume flows.Resume) error {
	// figure out where in the flow we began waiting on
	step, _, err := waitingRun.PathLocation()
	if err != nil {
		return err
	}

	// set up our flow stack based on the current run hierarchy
	s.flowStack = flowStackFromRun(waitingRun)

	// try to end our wait which will return and error if it can't be ended with this resume
	if err := s.wait.End(resume); err != nil {
		return err
	}
	s.wait = nil
	s.status = flows.SessionStatusActive

	// resumes are allowed to make state changes
	if err := resume.Apply(waitingRun, step); err != nil {
		return err
	}

	var destination flows.NodeUUID

	if waitingRun.Status() == flows.RunStatusWaiting {
		waitingRun.SetStatus(flows.RunStatusActive)

		destination, err = s.findResumeDestination(waitingRun)
		if err != nil {
			return err
		}
	}

	// off to the races again...
	return s.continueUntilWait(waitingRun, destination, step, nil)
}

// finds the next destination in a run that may have been waiting or a parent paused for a child subflow
func (s *session) findResumeDestination(run flows.FlowRun) (flows.NodeUUID, error) {
	step, node, err := run.PathLocation()
	if err != nil {
		return noDestination, err
	}

	// see if this node can now pick a destination
	step, destination, err := s.pickNodeExit(run, node, step)
	if err != nil {
		return noDestination, err
	}

	return destination, nil
}

// the main flow execution loop
func (s *session) continueUntilWait(currentRun flows.FlowRun, destination flows.NodeUUID, step flows.Step, trigger flows.Trigger) (err error) {
	for {
		// if we have a flow trigger handle that first to find our destination in the new flow
		if s.pushedFlow != nil {
			// if this is terminal, then we need to interrupt all other runs so we don't try to resume them
			if s.pushedFlow.terminal {
				for _, run := range s.runs {
					run.Exit(flows.RunStatusInterrupted)
				}
			}

			// create a new run for it
			flow := s.pushedFlow.flow
			currentRun = runs.NewRun(s, s.pushedFlow.flow, currentRun)
			s.addRun(currentRun)
			s.flowStack.push(flow)

			// our destination is the first node in that flow... if such a node exists
			if len(flow.Nodes()) > 0 {
				destination = flow.Nodes()[0].UUID()
			} else {
				destination = noDestination
			}

			// clear the trigger
			s.pushedFlow = nil
		}

		// if we have no destination then we're done with the current run which may have completed, expired or errored
		if destination == noDestination {
			if currentRun.ExitedOn() == nil {
				currentRun.Exit(flows.RunStatusCompleted)
			}

			parentRun := currentRun.ParentInSession()

			// switch back our parent run if it's still active
			if parentRun != nil && parentRun.Status() == flows.RunStatusActive {
				childRun := currentRun
				currentRun = parentRun
				s.flowStack.pop()

				// as long as we didn't error, we can try to resume it
				if childRun.Status() != flows.RunStatusErrored {
					if destination, err = s.findResumeDestination(currentRun); err != nil {
						currentRun.LogFatalError(step, fmt.Errorf("can't resume run as node no longer exists"))
					}
				} else {
					// if we did error then that needs to bubble back up through the run hierarchy
					step, _, _ := currentRun.PathLocation()
					currentRun.LogFatalError(step, fmt.Errorf("child run for flow '%s' ended in error, ending execution", childRun.Flow().UUID()))
				}

			} else {
				// If we have no destination and no parent, then the whole session is done. A run error bubbles up the session status.
				if currentRun.Status() == flows.RunStatusErrored {
					s.status = flows.SessionStatusErrored
				} else {
					s.status = flows.SessionStatusCompleted
				}

				// return to caller
				return nil
			}
		}

		// if we now have a destination, go there
		if destination != noDestination {
			if s.flowStack.hasVisited(destination) {
				// this is a loop, we log it and stop execution
				currentRun.LogFatalError(step, fmt.Errorf("flow loop detected, stopping execution before entering '%s'", destination))
				destination = noDestination
			} else {
				node := currentRun.Flow().GetNode(destination)
				if node == nil {
					return fmt.Errorf("unable to find destination node %s in flow %s", destination, currentRun.Flow().UUID())
				}

				step, destination, err = s.visitNode(currentRun, node, trigger)
				if err != nil {
					return err
				}

				// only want to pass this to the first node
				trigger = nil

				// if we hit a wait, also return to the caller
				if s.status == flows.SessionStatusWaiting {
					return nil
				}
			}
		}
	}
}

// visits the given node, creating a step in our current run path
func (s *session) visitNode(run flows.FlowRun, node flows.Node, trigger flows.Trigger) (flows.Step, flows.NodeUUID, error) {
	// mark this node as visited to prevent loops
	s.flowStack.visit(node.UUID())

	step := run.CreateStep(node)

	// this might be the first run of the session in which case a trigger might need to initialize the run
	if trigger != nil {
		if err := trigger.InitializeRun(run, step); err != nil {
			return step, noDestination, nil
		}
	}

	// execute our node's actions
	if node.Actions() != nil {
		for _, action := range node.Actions() {
			if log.GetLevel() >= log.DebugLevel {
				actionJSON, _ := json.Marshal(action)
				log.WithField("action_type", action.Type()).WithField("payload", string(actionJSON)).WithField("run", run.UUID()).Debug("action executing")
			}

			if err := action.Execute(run, step); err != nil {
				return step, noDestination, fmt.Errorf("error executing action[type=%s,uuid=%s]: %s", action.Type(), action.UUID(), err)
			}

			// check if this action has errored the run
			if run.Status() == flows.RunStatusErrored {
				return step, noDestination, nil
			}
		}
	}

	// a start flow action may have triggered a subflow in which case we're done on this node for now
	// and it will be resumed when the subflow finishes
	if s.pushedFlow != nil {
		return step, noDestination, nil
	}

	// our node might have wait
	wait := node.Wait()
	if wait != nil {

		// waits have the option to skip themselves
		if wait.Begin(run, step) {
			// mark ouselves as waiting and hand back to
			run.SetStatus(flows.RunStatusWaiting)
			s.wait = wait
			s.status = flows.SessionStatusWaiting

			return step, noDestination, nil
		}

		// if our wait skipped we need to pretend we were never at this node to avoid a loop error if we end up back here
		s.flowStack.unvisit(node.UUID())
	}

	// use our node's router to determine where to go next
	return s.pickNodeExit(run, node, step)
}

// picks the exit to use on the given node
func (s *session) pickNodeExit(run flows.FlowRun, node flows.Node, step flows.Step) (flows.Step, flows.NodeUUID, error) {
	var err error

	var operand *string
	route := flows.NoRoute
	router := node.Router()

	// we have a router, have it determine our exit
	var exitUUID flows.ExitUUID
	if router != nil {
		if operand, route, err = router.PickRoute(run, node.Exits(), step); err != nil {
			return nil, noDestination, fmt.Errorf("error routing from node[uuid=%s]: %s", node.UUID(), err)
		}
		exitUUID = route.Exit()
	} else if len(node.Exits()) > 0 {
		// no router, pick our first exit if we have one
		exitUUID = node.Exits()[0].UUID()
	}

	step.Leave(exitUUID)

	// look up our actual exit and localized name
	var exit flows.Exit
	var localizedExitName string

	if exitUUID != "" {
		// find our exit
		for _, e := range node.Exits() {
			if e.UUID() == exitUUID {
				localizedName := run.GetText(utils.UUID(exitUUID), "name", e.Name())
				if localizedName != e.Name() {
					localizedExitName = localizedName
				}
				exit = e
				break
			}
		}
		if exit == nil {
			return nil, noDestination, fmt.Errorf("unable to find exit with UUID '%s'", exitUUID)
		}
	}

	// no exit? return no destination
	if exit == nil {
		return step, noDestination, nil
	}

	// save our results if appropriate
	if router != nil && router.ResultName() != "" {
		var extraJSON json.RawMessage
		if route.Extra() != nil {
			extraJSON, _ = json.Marshal(route.Extra())
		}
		result := flows.NewResult(router.ResultName(), route.Match(), exit.Name(), localizedExitName, step.NodeUUID(), operand, extraJSON, utils.Now())
		run.SaveResult(result)
		event := events.NewRunResultChangedEvent(result)
		run.LogEvent(step, event)
	}

	return step, exit.DestinationNodeUUID(), nil
}

const noDestination = flows.NodeUUID("")

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type sessionEnvelope struct {
	Environment json.RawMessage     `json:"environment"`
	Trigger     json.RawMessage     `json:"trigger" validate:"required"`
	Contact     *json.RawMessage    `json:"contact,omitempty"`
	Runs        []json.RawMessage   `json:"runs"`
	Status      flows.SessionStatus `json:"status" validate:"required"`
	Wait        json.RawMessage     `json:"wait,omitempty"`
	Input       json.RawMessage     `json:"input,omitempty" validate:"omitempty"`
}

// ReadSession decodes a session from the passed in JSON
func ReadSession(assets flows.SessionAssets, engineConfig flows.EngineConfig, httpClient *utils.HTTPClient, data json.RawMessage) (flows.Session, error) {
	e := &sessionEnvelope{}
	var err error

	if err = utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, fmt.Errorf("unable to read session: %s", err)
	}

	s := NewSession(assets, engineConfig, httpClient).(*session)
	s.status = e.Status

	// read our environment
	s.env, err = utils.ReadEnvironment(e.Environment)
	if err != nil {
		return nil, fmt.Errorf("unable to read environment: %s", err)
	}

	// read our trigger
	if e.Trigger != nil {
		if s.trigger, err = triggers.ReadTrigger(s.Assets(), e.Trigger); err != nil {
			return nil, fmt.Errorf("unable to read trigger: %s", err)
		}
	}

	// read our contact
	if e.Contact != nil {
		if s.contact, err = flows.ReadContact(s.Assets(), *e.Contact, false); err != nil {
			return nil, fmt.Errorf("unable to read contact: %s", err)
		}
	}

	// read each of our runs
	for i := range e.Runs {
		run, err := runs.ReadRun(s, e.Runs[i])
		if err != nil {
			return nil, fmt.Errorf("unable to read run %d: %s", i, err)
		}
		s.addRun(run)
	}

	// and our wait and input
	if e.Wait != nil {
		s.wait, err = waits.ReadWait(e.Wait)
		if err != nil {
			return nil, fmt.Errorf("unable to read wait: %s", err)
		}
	}
	if e.Input != nil {
		if s.input, err = inputs.ReadInput(s, e.Input); err != nil {
			return nil, fmt.Errorf("unable to read input: %s", err)
		}
	}

	// TODO more and don't limit to sessions being read
	// perform some structural validation
	if s.status == flows.SessionStatusWaiting && s.wait == nil {
		return nil, fmt.Errorf("session has status of \"waiting\" but no wait object")
	}

	return s, nil
}

// MarshalJSON marshals this session into JSON
func (s *session) MarshalJSON() ([]byte, error) {
	e := &sessionEnvelope{
		Status: s.status,
	}
	var err error

	if e.Environment, err = json.Marshal(s.env); err != nil {
		return nil, err
	}
	if s.contact != nil {
		var contactJSON json.RawMessage
		contactJSON, err = json.Marshal(s.contact)
		if err != nil {
			return nil, err
		}
		e.Contact = &contactJSON
	}
	if s.trigger != nil {
		if e.Trigger, err = json.Marshal(s.trigger); err != nil {
			return nil, err
		}
	}
	if s.wait != nil {
		if e.Wait, err = json.Marshal(s.wait); err != nil {
			return nil, err
		}
	}
	if s.input != nil {
		e.Input, err = json.Marshal(s.input)
		if err != nil {
			return nil, err
		}
	}

	e.Runs = make([]json.RawMessage, len(s.runs))
	for i := range s.runs {
		e.Runs[i], err = json.Marshal(s.runs[i])
		if err != nil {
			return nil, err
		}
	}

	return json.Marshal(e)
}
