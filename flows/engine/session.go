package engine

import (
	"encoding/json"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/flows/inputs"
	"github.com/nyaruka/goflow/flows/resumes"
	"github.com/nyaruka/goflow/flows/routers/waits"
	"github.com/nyaruka/goflow/flows/runs"
	"github.com/nyaruka/goflow/flows/triggers"
	"github.com/nyaruka/goflow/utils"
	"github.com/nyaruka/goflow/utils/jsonx"

	"github.com/pkg/errors"
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
	uuid    flows.SessionUUID
	type_   flows.FlowType
	env     envs.Environment
	trigger flows.Trigger
	contact *flows.Contact
	runs    []flows.FlowRun
	status  flows.SessionStatus
	wait    flows.ActivatedWait
	input   flows.Input

	// state which is temporary to each call
	batchStart bool
	runsByUUID map[flows.RunUUID]flows.FlowRun
	pushedFlow *pushedFlow
	parentRun  flows.RunSummary

	engine flows.Engine
}

func (s *session) Assets() flows.SessionAssets { return s.assets }
func (s *session) Trigger() flows.Trigger      { return s.trigger }

func (s *session) UUID() flows.SessionUUID { return s.uuid }

func (s *session) Type() flows.FlowType         { return s.type_ }
func (s *session) SetType(type_ flows.FlowType) { s.type_ = type_ }

func (s *session) Environment() envs.Environment       { return s.env }
func (s *session) SetEnvironment(env envs.Environment) { s.env = env }

func (s *session) Contact() *flows.Contact           { return s.contact }
func (s *session) SetContact(contact *flows.Contact) { s.contact = contact }

func (s *session) Input() flows.Input         { return s.input }
func (s *session) SetInput(input flows.Input) { s.input = input }

func (s *session) BatchStart() bool { return s.batchStart }

func (s *session) PushFlow(flow flows.Flow, parentRun flows.FlowRun, terminal bool) {
	s.pushedFlow = &pushedFlow{flow: flow, parentRun: parentRun, terminal: terminal}
}

func (s *session) Runs() []flows.FlowRun { return s.runs }
func (s *session) GetRun(uuid flows.RunUUID) (flows.FlowRun, error) {
	run, exists := s.runsByUUID[uuid]
	if exists {
		return run, nil
	}
	return nil, errors.Errorf("unable to find run with UUID '%s'", uuid)
}

func (s *session) addRun(run flows.FlowRun) {
	s.runs = append(s.runs, run)
	s.runsByUUID[run.UUID()] = run
}

func (s *session) GetCurrentChild(run flows.FlowRun) flows.FlowRun {
	// the current child of a run, is the last added run which has that run as its parent
	for i := len(s.runs) - 1; i >= 0; i-- {
		if s.runs[i].ParentInSession() == run {
			return s.runs[i]
		}
	}
	return nil
}

// ParentRun gets the parent run of this session if it was started by a flow action
func (s *session) ParentRun() flows.RunSummary {
	return s.parentRun
}

func (s *session) Status() flows.SessionStatus { return s.status }
func (s *session) Wait() flows.ActivatedWait   { return s.wait }

func (s *session) CurrentContext() *types.XObject {
	run := s.currentRun()
	if run == nil {
		return nil
	}
	return types.NewXObject(run.RootContext(s.env))
}

// looks through this session's run for the one that was last modified
func (s *session) currentRun() flows.FlowRun {
	var lastRun flows.FlowRun
	for _, run := range s.runs {
		if lastRun == nil || run.ModifiedOn().After(lastRun.ModifiedOn()) {
			lastRun = run
		}
	}
	return lastRun
}

// looks through this session's run for the one that is waiting
func (s *session) waitingRun() flows.FlowRun {
	for _, run := range s.runs {
		if run.Status() == flows.RunStatusWaiting {
			return run
		}
	}
	return nil
}

func (s *session) Engine() flows.Engine { return s.engine }

//------------------------------------------------------------------------------------------
// Flow execution
//------------------------------------------------------------------------------------------

// Start initializes this session with the given trigger and runs the flow to the first wait
func (s *session) start(trigger flows.Trigger) (flows.Sprint, error) {
	sprint := NewEmptySprint()

	if err := s.prepareForSprint(); err != nil {
		return sprint, err
	}

	if err := s.trigger.Initialize(s, sprint.LogEvent); err != nil {
		return sprint, err
	}

	// off to the races...
	if err := s.continueUntilWait(sprint, nil, noDestination, nil, trigger); err != nil {
		return sprint, err
	}

	return sprint, nil
}

// Resume tries to resume a waiting session
func (s *session) Resume(resume flows.Resume) (flows.Sprint, error) {
	sprint := NewEmptySprint()

	if err := s.prepareForSprint(); err != nil {
		return sprint, err
	}

	if s.status != flows.SessionStatusWaiting {
		return sprint, errors.Errorf("only waiting sessions can be resumed")
	}

	waitingRun := s.waitingRun()
	if waitingRun == nil {
		return sprint, errors.Errorf("session doesn't contain any runs which are waiting")
	}

	if err := s.tryToResume(sprint, waitingRun, resume); err != nil {
		// if we got an error, add it to the log and shut everything down
		failure(sprint, waitingRun, nil, err)

		s.status = flows.SessionStatusFailed
	}

	return sprint, nil
}

// prepares the session for starting/resuming
func (s *session) prepareForSprint() error {
	if s.parentRun == nil {
		// if we have a trigger with a parent run, load that
		triggerWithRun, hasRun := s.trigger.(flows.TriggerWithRun)
		if hasRun {
			run, err := runs.ReadRunSummary(s.Assets(), triggerWithRun.RunSummary(), assets.IgnoreMissing)
			if err != nil {
				return errors.Wrap(err, "error reading parent run from trigger")
			}
			s.parentRun = run
		}
	}
	return nil
}

// Resume resumes a waiting session
func (s *session) tryToResume(sprint flows.Sprint, waitingRun flows.FlowRun, resume flows.Resume) error {
	// if flow for this run is a missing asset, we have a problem
	if waitingRun.Flow() == nil {
		return errors.New("can't resume run with missing flow asset")
	}

	// figure out where in the flow we began waiting on
	step, node, err := waitingRun.PathLocation()
	if err != nil {
		return err
	}

	if node.Router() == nil || node.Router().Wait() == nil {
		return errors.New("can't resume from node without a router or wait")
	}

	// try to end our wait which will return and log an error if it can't be ended with this resume
	if err := node.Router().Wait().End(resume); err != nil {
		sprint.LogEvent(events.NewError(err))
		return nil
	}
	s.wait = nil
	s.status = flows.SessionStatusActive

	logEvent := func(e flows.Event) {
		waitingRun.LogEvent(step, e)
		sprint.LogEvent(e)
	}

	// resumes are allowed to make state changes
	if err := resume.Apply(waitingRun, logEvent); err != nil {
		return err
	}

	_, isTimeout := resume.(*resumes.WaitTimeoutResume)

	destination, err := s.findResumeDestination(sprint, waitingRun, isTimeout)
	if err != nil {
		return err
	}

	// off to the races again...
	return s.continueUntilWait(sprint, waitingRun, destination, step, nil)
}

// finds the next destination in a run that may have been waiting or a parent paused for a child subflow
func (s *session) findResumeDestination(sprint flows.Sprint, run flows.FlowRun, isTimeout bool) (flows.NodeUUID, error) {
	// we might have no immediate destination in this run, but continueUntilWait can resume a parent run
	if run.Status() != flows.RunStatusActive {
		return noDestination, nil
	}

	step, node, err := run.PathLocation()
	if err != nil {
		return noDestination, err
	}
	logEvent := func(e flows.Event) {
		run.LogEvent(step, e)
		sprint.LogEvent(e)
	}

	// see if this node can now pick a destination
	destination, err := s.pickNodeExit(sprint, run, node, step, isTimeout, logEvent)
	if err != nil {
		return noDestination, err
	}

	return destination, nil
}

// the main flow execution loop
func (s *session) continueUntilWait(sprint flows.Sprint, currentRun flows.FlowRun, destination flows.NodeUUID, step flows.Step, trigger flows.Trigger) (err error) {
	numNewSteps := 0

	for {
		// if we have a flow trigger handle that first to find our destination in the new flow
		if s.pushedFlow != nil {
			// if this is terminal, then we need to mark all other runs as completed so we don't try to resume them
			if s.pushedFlow.terminal {
				for _, run := range s.runs {
					run.Exit(flows.RunStatusCompleted)
				}
			}

			// create a new run for it
			flow := s.pushedFlow.flow
			currentRun = runs.NewRun(s, s.pushedFlow.flow, currentRun)
			s.addRun(currentRun)

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

				// as long as we didn't error, we can try to resume it
				if childRun.Status() != flows.RunStatusFailed {
					// if flow for this run is a missing asset, we have a problem
					if currentRun.Flow() == nil {
						return errors.New("can't resume parent run with missing flow asset")
					}

					if destination, err = s.findResumeDestination(sprint, currentRun, false); err != nil {
						failure(sprint, currentRun, step, errors.Wrapf(err, "can't resume run as node no longer exists"))
					}
				} else {
					// if we did error then that needs to bubble back up through the run hierarchy
					step, _, _ := currentRun.PathLocation()
					failure(sprint, currentRun, step, errors.Errorf("child run for flow '%s' ended in error, ending execution", childRun.Flow().UUID()))
				}

			} else {
				// If we have no destination and no parent, then the whole session is done. A run error bubbles up the session status.
				if currentRun.Status() == flows.RunStatusFailed {
					s.status = flows.SessionStatusFailed
				} else {
					s.status = flows.SessionStatusCompleted
				}

				// return to caller
				return nil
			}
		}

		// if we now have a destination, go there
		if destination != noDestination {
			numNewSteps++

			if numNewSteps > s.Engine().MaxStepsPerSprint() {
				// we've hit the step limit - usually a sign of a loop
				failure(sprint, currentRun, step, errors.Errorf("step limit exceeded, stopping execution before entering '%s'", destination))
				destination = noDestination
			} else {
				node := currentRun.Flow().GetNode(destination)
				if node == nil {
					return errors.Errorf("unable to find destination node %s in flow %s", destination, currentRun.Flow().UUID())
				}

				step, destination, err = s.visitNode(sprint, currentRun, node, trigger)
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
func (s *session) visitNode(sprint flows.Sprint, run flows.FlowRun, node flows.Node, trigger flows.Trigger) (flows.Step, flows.NodeUUID, error) {
	step := run.CreateStep(node)
	logEvent := func(e flows.Event) {
		run.LogEvent(step, e)
		sprint.LogEvent(e)
	}

	// this might be the first run of the session in which case a trigger might need to initialize the run
	if trigger != nil {
		if err := trigger.InitializeRun(run, logEvent); err != nil {
			return step, noDestination, nil
		}
	}

	// execute our node's actions
	if node.Actions() != nil {
		for _, action := range node.Actions() {
			if err := action.Execute(run, step, sprint.LogModifier, logEvent); err != nil {
				return step, noDestination, errors.Wrapf(err, "error executing action[type=%s,uuid=%s]", action.Type(), action.UUID())
			}

			// check if this action has errored the run
			if run.Status() == flows.RunStatusFailed {
				return step, noDestination, nil
			}
		}
	}

	// a start flow action may have triggered a subflow in which case we're done on this node for now
	// and it will be resumed when the subflow finishes
	if s.pushedFlow != nil {
		return step, noDestination, nil
	}

	// our node might have a router with a wait
	var wait flows.Wait
	if node.Router() != nil {
		wait = node.Router().Wait()
	}

	if wait != nil {

		// waits have the option to skip themselves
		activatedWait := wait.Begin(run, logEvent)
		if activatedWait != nil {
			// mark ouselves as waiting and hand back to
			run.SetStatus(flows.RunStatusWaiting)
			s.wait = activatedWait
			s.status = flows.SessionStatusWaiting

			return step, noDestination, nil
		}
	}

	// use our node's router to determine where to go next
	destinationUUID, err := s.pickNodeExit(sprint, run, node, step, false, logEvent)
	return step, destinationUUID, err
}

// picks the exit to use on the given node
func (s *session) pickNodeExit(sprint flows.Sprint, run flows.FlowRun, node flows.Node, step flows.Step, isTimeout bool, logEvent flows.EventCallback) (flows.NodeUUID, error) {
	var exitUUID flows.ExitUUID
	var err error

	if node.Router() != nil {
		if isTimeout {
			exitUUID, err = node.Router().RouteTimeout(run, step, logEvent)
		} else {
			exitUUID, err = node.Router().Route(run, step, logEvent)
		}

		if err != nil {
			return noDestination, errors.Wrapf(err, "error routing from node[uuid=%s]", node.UUID())
		}
		// router didn't error.. but it failed to pick a category
		if exitUUID == "" {
			failure(sprint, run, step, errors.Errorf("router on node[uuid=%s] failed to pick a category", node.UUID()))
			return noDestination, nil
		}
	} else if len(node.Exits()) > 0 {
		// no router, pick our first exit if we have one
		exitUUID = node.Exits()[0].UUID()
	}

	step.Leave(exitUUID)

	// find our exit
	for _, exit := range node.Exits() {
		if exit.UUID() == exitUUID {
			return exit.DestinationUUID(), nil
		}
	}

	// no exit? return no destination
	return noDestination, nil
}

const noDestination = flows.NodeUUID("")

// utility to fail the session and log a failure event
func failure(sprint flows.Sprint, run flows.FlowRun, step flows.Step, err error) {
	event := events.NewFailure(err)
	if run != nil {
		run.Exit(flows.RunStatusFailed)
		run.LogEvent(step, event)
	}
	sprint.LogEvent(event)
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type sessionEnvelope struct {
	UUID        flows.SessionUUID   `json:"uuid"` // TODO validate:"required"`
	Type        flows.FlowType      `json:"type"` // TODO validate:"required"`
	Environment json.RawMessage     `json:"environment"`
	Trigger     json.RawMessage     `json:"trigger" validate:"required"`
	Contact     *json.RawMessage    `json:"contact,omitempty"`
	Runs        []json.RawMessage   `json:"runs"`
	Status      flows.SessionStatus `json:"status" validate:"required"`
	Wait        json.RawMessage     `json:"wait,omitempty"`
	Input       json.RawMessage     `json:"input,omitempty" validate:"omitempty"`
}

// ReadSession decodes a session from the passed in JSON
func readSession(eng flows.Engine, sessionAssets flows.SessionAssets, data json.RawMessage, missing assets.MissingCallback) (flows.Session, error) {
	e := &sessionEnvelope{}
	var err error

	if err = utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, errors.Wrap(err, "unable to read session")
	}

	s := &session{
		engine:     eng,
		assets:     sessionAssets,
		uuid:       e.UUID,
		type_:      e.Type,
		status:     e.Status,
		runsByUUID: make(map[flows.RunUUID]flows.FlowRun),
	}

	// read our environment
	s.env, err = envs.ReadEnvironment(e.Environment)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read environment")
	}

	// read our trigger
	if e.Trigger != nil {
		if s.trigger, err = triggers.ReadTrigger(s.Assets(), e.Trigger, missing); err != nil {
			return nil, errors.Wrap(err, "unable to read trigger")
		}
	}

	// read our contact
	if e.Contact != nil {
		if s.contact, err = flows.ReadContact(s.Assets(), *e.Contact, missing); err != nil {
			return nil, errors.Wrap(err, "unable to read contact")
		}
	}

	// read each of our runs
	for i := range e.Runs {
		run, err := runs.ReadRun(s, e.Runs[i], missing)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to read run %d", i)
		}
		s.addRun(run)
	}

	// and our wait and input
	if e.Wait != nil {
		s.wait, err = waits.ReadActivatedWait(e.Wait)
		if err != nil {
			return nil, errors.Wrap(err, "unable to read wait")
		}
	}
	if e.Input != nil {
		if s.input, err = inputs.ReadInput(s.Assets(), e.Input, missing); err != nil {
			return nil, errors.Wrap(err, "unable to read input")
		}
	}

	// TODO more and don't limit to sessions being read
	// perform some structural validation
	if s.status == flows.SessionStatusWaiting && s.wait == nil {
		return nil, errors.Errorf("session has status of \"waiting\" but no wait object")
	}

	return s, nil
}

// MarshalJSON marshals this session into JSON
func (s *session) MarshalJSON() ([]byte, error) {
	e := &sessionEnvelope{
		UUID:   s.uuid,
		Type:   s.type_,
		Status: s.status,
	}
	var err error

	if e.Environment, err = jsonx.Marshal(s.env); err != nil {
		return nil, err
	}
	if s.contact != nil {
		var contactJSON json.RawMessage
		contactJSON, err = jsonx.Marshal(s.contact)
		if err != nil {
			return nil, err
		}
		e.Contact = &contactJSON
	}
	if s.trigger != nil {
		if e.Trigger, err = jsonx.Marshal(s.trigger); err != nil {
			return nil, err
		}
	}
	if s.wait != nil {
		if e.Wait, err = jsonx.Marshal(s.wait); err != nil {
			return nil, err
		}
	}
	if s.input != nil {
		e.Input, err = jsonx.Marshal(s.input)
		if err != nil {
			return nil, err
		}
	}

	e.Runs = make([]json.RawMessage, len(s.runs))
	for i := range s.runs {
		e.Runs[i], err = jsonx.Marshal(s.runs[i])
		if err != nil {
			return nil, err
		}
	}

	return jsonx.Marshal(e)
}
