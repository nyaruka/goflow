package engine

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/flows/inputs"
	"github.com/nyaruka/goflow/flows/resumes"
	"github.com/nyaruka/goflow/flows/runs"
	"github.com/nyaruka/goflow/flows/triggers"
	"github.com/nyaruka/goflow/utils"
)

// used to spawn a new run or sub-flow in the event loop
type pushedFlow struct {
	flow      flows.Flow
	parentRun flows.Run
	terminal  bool
}

type session struct {
	assets flows.SessionAssets

	// state which is maintained between engine calls
	uuid          flows.SessionUUID
	type_         flows.FlowType
	createdOn     time.Time
	env           envs.Environment
	trigger       flows.Trigger
	currentResume flows.Resume
	contact       *flows.Contact
	call          *flows.Call

	runs   []flows.Run
	status flows.SessionStatus
	input  flows.Input

	// state which is temporary to each call
	batchStart bool
	runsByUUID map[flows.RunUUID]flows.Run
	pushedFlow *pushedFlow
	parentRun  flows.RunSummary

	engine flows.Engine
}

func (s *session) Assets() flows.SessionAssets { return s.assets }
func (s *session) Trigger() flows.Trigger      { return s.trigger }
func (s *session) CurrentResume() flows.Resume { return s.currentResume }

func (s *session) UUID() flows.SessionUUID             { return s.uuid }
func (s *session) Type() flows.FlowType                { return s.type_ }
func (s *session) CreatedOn() time.Time                { return s.createdOn }
func (s *session) Environment() envs.Environment       { return s.env }
func (s *session) MergedEnvironment() envs.Environment { return flows.NewSessionEnvironment(s) }
func (s *session) Contact() *flows.Contact             { return s.contact }
func (s *session) Call() *flows.Call                   { return s.call }

func (s *session) Input() flows.Input { return s.input }
func (s *session) SetInput(input flows.Input) {
	s.input = input

	// if we have a contact, update their last seen date
	if input != nil && s.contact != nil {
		s.contact.SetLastSeenOn(input.CreatedOn())
	}
}

func (s *session) BatchStart() bool { return s.batchStart }

func (s *session) PushFlow(flow flows.Flow, parentRun flows.Run, terminal bool) {
	s.pushedFlow = &pushedFlow{flow: flow, parentRun: parentRun, terminal: terminal}
}

func (s *session) Runs() []flows.Run { return s.runs }
func (s *session) GetRun(uuid flows.RunUUID) (flows.Run, error) {
	run, exists := s.runsByUUID[uuid]
	if exists {
		return run, nil
	}
	return nil, fmt.Errorf("unable to find run with UUID '%s'", uuid)
}

func (s *session) FindStep(uuid flows.StepUUID) (flows.Run, flows.Step) {
	for _, r := range s.runs {
		for _, t := range r.Path() {
			if t.UUID() == uuid {
				return r, t
			}
		}
	}
	return nil, nil
}

func (s *session) addRun(run flows.Run) {
	s.runs = append(s.runs, run)
	s.runsByUUID[run.UUID()] = run
}

func (s *session) GetCurrentChild(run flows.Run) flows.Run {
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

func (s *session) CurrentContext() *types.XObject {
	run := s.currentRun()
	if run == nil {
		return nil
	}
	return types.NewXObject(run.RootContext(s.env))
}

// looks through this session's run for the one that was last modified
func (s *session) currentRun() flows.Run {
	var last flows.Run
	for _, r := range s.runs {
		if last == nil || r.ModifiedOn().After(last.ModifiedOn()) {
			last = r
		}
	}
	return last
}

// looks through this session's run for the one that is waiting
func (s *session) waitingRun() flows.Run {
	for _, run := range s.runs {
		if run.Status() == flows.RunStatusWaiting {
			return run
		}
	}
	return nil
}

func (s *session) History() *flows.SessionHistory {
	history := s.trigger.History()
	if history != nil {
		return history
	}
	return flows.EmptyHistory
}

func (s *session) Engine() flows.Engine { return s.engine }

//------------------------------------------------------------------------------------------
// Flow execution
//------------------------------------------------------------------------------------------

// Start initializes this session with the given trigger and runs the flow to the first wait
func (s *session) start(ctx context.Context, trigger flows.Trigger, flow flows.Flow) (flows.Sprint, error) {
	sprint := newEmptySprint(true)

	if err := s.prepareForSprint(); err != nil {
		return sprint, err
	}

	s.PushFlow(flow, nil, false)

	if err := s.trigger.Initialize(s); err != nil {
		return sprint, err
	}

	// ensure groups are correct
	s.ensureQueryBasedGroups(sprint.logEvent)

	// off to the races...
	if err := s.continueUntilWait(ctx, sprint, nil, nil, nil, "", nil, trigger); err != nil {
		return sprint, err
	}

	return sprint, nil
}

// Resume tries to resume a waiting session
func (s *session) Resume(ctx context.Context, resume flows.Resume) (flows.Sprint, error) {
	sprint := newEmptySprint(false)

	if err := s.prepareForSprint(); err != nil {
		return sprint, err
	}

	if s.status != flows.SessionStatusWaiting {
		return sprint, newError(ErrorResumeNonWaitingSession, "only waiting sessions can be resumed")
	}

	waitingRun := s.waitingRun()
	if waitingRun == nil {
		return sprint, newError(ErrorResumeNoWaitingRun, "session doesn't contain any runs which are waiting")
	}

	if err := s.tryToResume(ctx, sprint, waitingRun, resume); err != nil {
		return nil, err
	}

	return sprint, nil
}

// prepares the session for starting/resuming
func (s *session) prepareForSprint() error {
	if s.parentRun == nil {
		// if we have a trigger with a parent run, load that
		triggerWithRun, hasRun := s.trigger.(flows.TriggerWithRun)
		if hasRun {
			r, err := runs.ReadRunSummary(s.Assets(), triggerWithRun.RunSummary(), assets.IgnoreMissing)
			if err != nil {
				return fmt.Errorf("error reading parent run from trigger: %w", err)
			}
			s.parentRun = r
		}
	}
	return nil
}

// tries to resume a waiting session with the given resume
func (s *session) tryToResume(ctx context.Context, sprint *sprint, waitingRun flows.Run, resume flows.Resume) error {
	failSession := func(msg string, args ...any) {
		// put failure event in waiting run
		failRun(sprint, waitingRun, nil, fmt.Errorf(msg, args...))

		// but also fail any other non-exited runs
		for _, r := range s.runs {
			if r.Status() == flows.RunStatusActive || r.Status() == flows.RunStatusWaiting {
				r.Exit(flows.RunStatusFailed)
			}
		}

		s.status = flows.SessionStatusFailed
	}

	// if flow for this run is a missing asset, we have a problem
	if waitingRun.Flow() == nil {
		failSession("can't resume run with missing flow asset")
		return nil
	}

	if s.countWaits() >= s.engine.Options().MaxResumesPerSession {
		failSession("reached maximum number of resumes per session (%d)", s.engine.Options().MaxResumesPerSession)
		return nil
	}

	// figure out where in the flow we began waiting on
	step, node, err := waitingRun.PathLocation()
	if err != nil {
		failSession(fmt.Sprintf("unable to find resume location: %s", err.Error()))
		return nil
	}

	if node.Router() == nil || node.Router().Wait() == nil {
		failSession("can't resume from node without a router or wait")
		return nil
	}

	// check that the wait accepts this resume - not a permanent error - caller can retry with different resume
	if !node.Router().Wait().Accepts(resume) {
		return newError(ErrorResumeRejectedByWait, "resume of type %s not accepted by wait of type %s", resume.Type(), node.Router().Wait().Type())
	}

	s.status = flows.SessionStatusActive
	s.currentResume = resume
	sprint.logFlow(waitingRun.Flow())

	logEvent := func(e flows.Event) {
		waitingRun.LogEvent(step, e)
		sprint.logEvent(e)
	}

	// if resume was based on an event (e.g. msg received), log that on the run but don't repeat it in the sprint
	// events because we didn't generate it
	if resume.Event() != nil {
		waitingRun.LogEvent(nil, resume.Event())
	}

	// resumes are allowed to make state changes
	resume.Apply(waitingRun, logEvent)

	// ensure groups are correct
	s.ensureQueryBasedGroups(logEvent)

	_, isTimeout := resume.(*resumes.WaitTimeoutResume)

	exit, operand, err := s.findResumeExit(sprint, waitingRun, isTimeout)
	if err != nil {
		failSession(fmt.Sprintf("unable to resolve router exit: %s", err.Error()))
		return nil
	}

	// off to the races again...
	return s.continueUntilWait(ctx, sprint, waitingRun, node, exit, operand, step, nil)
}

// finds the exit from a the current node in a run that may have been waiting or a parent paused for a child subflow
func (s *session) findResumeExit(sprint *sprint, run flows.Run, isTimeout bool) (flows.Exit, string, error) {
	// we might have no immediate destination in this run, but continueUntilWait can resume a parent run
	if run.Status() != flows.RunStatusActive {
		return nil, "", nil
	}

	step, node, err := run.PathLocation()
	if err != nil {
		return nil, "", err
	}
	logEvent := func(e flows.Event) {
		run.LogEvent(step, e)
		sprint.logEvent(e)
	}

	// see if this node can now pick a destination
	return s.pickNodeExit(sprint, run, node, step, isTimeout, logEvent)
}

// the main flow execution loop
func (s *session) continueUntilWait(ctx context.Context, sprint *sprint, currentRun flows.Run, node flows.Node, exit flows.Exit, operand string, step flows.Step, trigger flows.Trigger) (err error) {
	var destination flows.NodeUUID
	var numNewSteps int

	for {
		// start by picking a destination node...

		// if a new flow has been pushed, find a destination there
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
			sprint.logFlow(flow)

			// our destination is the first node in that flow... if such a node exists
			if len(flow.Nodes()) > 0 {
				destination = flow.Nodes()[0].UUID()
			} else {
				destination = ""
			}

			s.pushedFlow = nil // clear the trigger

		} else if exit != nil {
			// if we're at an exit, use its destination
			destination = exit.DestinationUUID()

			// if we have a destination, record as segment
			if destination != "" {
				destNode := currentRun.Flow().GetNode(destination)
				if destNode != nil {
					sprint.logSegment(currentRun.Flow(), node, exit, operand, destNode)
				}
			}

			// clear the exit and operand
			exit, operand = nil, ""
		} else {
			destination = ""
		}

		// if we have no destination then we're done with the current run which may have completed, expired or errored
		if destination == "" {
			if currentRun.ExitedOn() == nil {
				currentRun.Exit(flows.RunStatusCompleted)
			}

			parentRun := currentRun.ParentInSession()

			// switch back our parent run if it's still active
			if parentRun != nil && parentRun.Status() == flows.RunStatusActive {
				childRun := currentRun
				currentRun = parentRun

				// as long as we didn't fail, we can try to resume it
				if childRun.Status() != flows.RunStatusFailed {
					// if flow for this run is a missing asset, we have a problem
					if currentRun.Flow() == nil {
						failRun(sprint, currentRun, nil, errors.New("can't resume run with missing flow asset"))
					} else {
						if exit, operand, err = s.findResumeExit(sprint, currentRun, false); err != nil {
							failRun(sprint, currentRun, nil, fmt.Errorf("can't resume run as node no longer exists: %w", err))
						}
					}
				} else {
					// if we did fail then that needs to bubble back up through the run hierarchy
					step, _, _ := currentRun.PathLocation()
					failRun(sprint, currentRun, step, fmt.Errorf("child run for flow '%s' ended in error, ending execution", childRun.FlowReference().UUID))
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
		if destination != "" {
			numNewSteps++

			if numNewSteps > s.engine.Options().MaxStepsPerSprint {
				// we've hit the step limit - usually a sign of a loop
				failRun(sprint, currentRun, step, fmt.Errorf("reached maximum number of steps per sprint (%d)", s.engine.Options().MaxStepsPerSprint))
			} else {
				node = currentRun.Flow().GetNode(destination)
				if node == nil {
					return fmt.Errorf("unable to find destination node %s in flow %s", destination, currentRun.Flow().UUID())
				}

				step, exit, operand, err = s.visitNode(ctx, sprint, currentRun, node, trigger)
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
func (s *session) visitNode(ctx context.Context, sprint *sprint, run flows.Run, node flows.Node, trigger flows.Trigger) (flows.Step, flows.Exit, string, error) {
	step := run.CreateStep(node)
	logEvent := func(e flows.Event) {
		run.LogEvent(step, e)
		sprint.logEvent(e)
	}

	// this might be the first run of the session in which case a trigger might need to initialize the run
	if trigger != nil {
		// if trigger was based on an event (e.g. msg received), log that on the run but don't repeat it in the sprint
		// events because we didn't generate it
		if trigger.Event() != nil {
			run.LogEvent(nil, trigger.Event())
		}
	}

	// execute our node's actions
	if node.Actions() != nil {
		for _, action := range node.Actions() {
			if err := action.Execute(ctx, run, step, sprint.logModifier, logEvent); err != nil {
				return step, nil, "", fmt.Errorf("error executing action[type=%s,uuid=%s]: %w", action.Type(), action.UUID(), err)
			}

			// check if this action has errored the run
			if run.Status() == flows.RunStatusFailed {
				return step, nil, "", nil
			}
		}
	}

	// a start flow action may have triggered a subflow in which case we're done on this node for now
	// and it will be resumed when the subflow finishes
	if s.pushedFlow != nil {
		return step, nil, "", nil
	}

	// our node might have a router with a wait
	var wait flows.Wait
	if node.Router() != nil {
		wait = node.Router().Wait()
	}

	if wait != nil {
		// waits have the option to skip themselves
		if wait.Begin(run, logEvent) {
			// mark ouselves as waiting and hand back to
			run.SetStatus(flows.RunStatusWaiting)
			s.status = flows.SessionStatusWaiting

			return step, nil, "", nil
		}
	}

	// use our node's router to determine where to go next
	exit, operand, err := s.pickNodeExit(sprint, run, node, step, false, logEvent)
	return step, exit, operand, err
}

// picks the exit to use on the given node
func (s *session) pickNodeExit(sprint *sprint, run flows.Run, node flows.Node, step flows.Step, isTimeout bool, logEvent flows.EventCallback) (flows.Exit, string, error) {
	var exitUUID flows.ExitUUID
	var operand string
	var err error

	if node.Router() != nil {
		if isTimeout {
			exitUUID, err = node.Router().RouteTimeout(run, step, logEvent)
		} else {
			exitUUID, operand, err = node.Router().Route(run, step, logEvent)
		}

		if err != nil {
			return nil, "", fmt.Errorf("error routing from node[uuid=%s]: %w", node.UUID(), err)
		}
		// router didn't error.. but it failed to pick a category
		if exitUUID == "" {
			failRun(sprint, run, step, fmt.Errorf("router on node[uuid=%s] failed to pick a category", node.UUID()))
			return nil, "", nil
		}
	} else if len(node.Exits()) > 0 {
		// no router, pick our first exit if we have one
		exitUUID = node.Exits()[0].UUID()
	}

	step.Leave(exitUUID)

	// find our exit
	for _, exit := range node.Exits() {
		if exit.UUID() == exitUUID {
			return exit, operand, nil
		}
	}

	return nil, "", nil // no where to go in the flow...
}

// ensures that our session contact is in the correct query based groups as as far as the engine is concerned
func (s *session) ensureQueryBasedGroups(logEvent flows.EventCallback) {
	if s.contact == nil {
		return
	}

	added, removed := s.contact.ReevaluateQueryBasedGroups(s.Environment())

	// add groups changed event for the groups we were added/removed to/from
	if len(added) > 0 || len(removed) > 0 {
		logEvent(events.NewContactGroupsChanged(added, removed))
	}
}

func (s *session) countWaits() int {
	waits := 0
	for _, r := range s.runs {
		for _, e := range r.Events() {
			if strings.HasSuffix(e.Type(), "_wait") {
				waits++
			}
		}
	}
	return waits
}

// utility to fail the current run and log a failRun event
func failRun(sp *sprint, run flows.Run, step flows.Step, err error) {
	event := events.NewFailure(err)
	run.Exit(flows.RunStatusFailed)
	run.LogEvent(step, event)
	sp.logEvent(event)
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type sessionEnvelope struct {
	UUID        flows.SessionUUID   `json:"uuid"                validate:"required"`
	Type        flows.FlowType      `json:"type"                validate:"required"`
	CreatedOn   time.Time           `json:"created_on"` // TODO validate:"required"`
	Trigger     json.RawMessage     `json:"trigger"             validate:"required"`
	ContactUUID flows.ContactUUID   `json:"contact_uuid"        validate:"omitempty,uuid"` // TODO validate:"required"`
	CallUUID    flows.CallUUID      `json:"call_uuid,omitempty" validate:"omitempty,uuid"` // TODO validate:"required"`
	Runs        []json.RawMessage   `json:"runs"`
	Status      flows.SessionStatus `json:"status"              validate:"required"`
	Wait        json.RawMessage     `json:"wait,omitempty"`
	Input       json.RawMessage     `json:"input,omitempty"`
}

// ReadSession decodes a session from the passed in JSON
func readSession(eng flows.Engine, sa flows.SessionAssets, data []byte, env envs.Environment, contact *flows.Contact, call *flows.Call, missing assets.MissingCallback) (flows.Session, error) {
	e := &sessionEnvelope{}
	var err error

	if err = utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, fmt.Errorf("unable to read session: %w", err)
	}

	s := &session{
		engine:     eng,
		assets:     sa,
		uuid:       e.UUID,
		type_:      e.Type,
		status:     e.Status,
		env:        env,
		contact:    contact,
		call:       call,
		runsByUUID: make(map[flows.RunUUID]flows.Run),
	}

	if e.Trigger != nil {
		if s.trigger, err = triggers.ReadTrigger(s.Assets(), e.Trigger, missing); err != nil {
			return nil, fmt.Errorf("unable to read trigger: %w", err)
		}
	}
	if e.ContactUUID != "" && contact.UUID() != e.ContactUUID {
		return nil, fmt.Errorf("session contact doesn't match provided")
	}
	if e.CallUUID != "" {
		if call == nil || call.UUID() != e.CallUUID {
			return nil, fmt.Errorf("session call doesn't match provided")
		}
	}

	// read each of our runs
	for i := range e.Runs {
		run, err := runs.ReadRun(s, e.Runs[i], missing)
		if err != nil {
			return nil, fmt.Errorf("unable to read run %d: %w", i, err)
		}
		s.addRun(run)
	}

	// and our input
	if e.Input != nil {
		if s.input, err = inputs.ReadInput(s.Assets(), e.Input, missing); err != nil {
			return nil, fmt.Errorf("unable to read input: %w", err)
		}
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

	if s.contact != nil {
		e.ContactUUID = s.contact.UUID()
	}
	if s.trigger != nil {
		if e.Trigger, err = jsonx.Marshal(s.trigger); err != nil {
			return nil, err
		}
	}
	if s.call != nil {
		e.CallUUID = s.call.UUID()
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
