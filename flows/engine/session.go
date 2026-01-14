package engine

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/flows/inputs"
	"github.com/nyaruka/goflow/flows/resumes"
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

	runs    []flows.Run
	status  flows.SessionStatus
	input   flows.Input
	sprints int

	// state which is temporary to each call
	batchStart bool
	runsByUUID map[flows.RunUUID]*run
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
func (s *session) Sprints() int                        { return s.sprints }

func (s *session) Input() flows.Input { return s.input }
func (s *session) setInput(input flows.Input) {
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
func (s *session) getRun(uuid flows.RunUUID) (*run, error) {
	r, exists := s.runsByUUID[uuid]
	if exists {
		return r, nil
	}
	return nil, fmt.Errorf("unable to find run with UUID '%s'", uuid)
}

func (s *session) addRun(r *run) {
	s.runs = append(s.runs, r)
	s.runsByUUID[r.UUID()] = r
}

func (s *session) findCurrentChild(r *run) *run {
	// the current child of a run, is the last added run which has that run as its parent
	for i := len(s.runs) - 1; i >= 0; i-- {
		sr := s.runs[i].(*run)

		if sr.parent == r {
			return sr
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
func (s *session) currentRun() *run {
	var last flows.Run
	for _, r := range s.runs {
		if last == nil || r.ModifiedOn().After(last.ModifiedOn()) {
			last = r
		}
	}
	if last != nil {
		return last.(*run)
	}
	return nil
}

// looks through this session's run for the one that is waiting
func (s *session) waitingRun() *run {
	for _, r := range s.runs {
		if r.Status() == flows.RunStatusWaiting {
			return r.(*run)
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

	// if trigger provides input, set it
	s.setInput(s.trigger.Input(s.assets))

	// ensure groups are correct
	s.ensureQueryBasedGroups(sprint.logEvent)

	// off to the races...
	if err := s.continueUntilWait(ctx, sprint, nil, nil, nil, "", nil, trigger); err != nil {
		return sprint, err
	}

	s.sprints++

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

	s.sprints++

	return sprint, nil
}

// prepares the session for starting/resuming
func (s *session) prepareForSprint() error {
	if s.parentRun == nil {
		// if we have a trigger with a parent run, load that
		triggerWithRun, hasRun := s.trigger.(flows.TriggerWithRun)
		if hasRun {
			r, err := ReadRunSummary(s.Assets(), triggerWithRun.RunSummary(), assets.IgnoreMissing)
			if err != nil {
				return fmt.Errorf("error reading parent run from trigger: %w", err)
			}
			s.parentRun = r
		}
	}
	return nil
}

// tries to resume a waiting session with the given resume
func (s *session) tryToResume(ctx context.Context, sprint *sprint, waitingRun *run, resume flows.Resume) error {
	failSession := func(msg string, args ...any) {
		// put failure event in waiting run
		failRun(sprint, waitingRun, nil, fmt.Errorf(msg, args...))

		// but also fail any other non-exited runs
		for _, r := range s.runs {
			if r.Status() == flows.RunStatusActive || r.Status() == flows.RunStatusWaiting {
				r.Exit(flows.RunStatusFailed)
				sprint.logEvent(events.NewRunEnded(r.UUID(), r.FlowReference(), flows.RunStatusFailed))
			}
		}

		s.status = flows.SessionStatusFailed
	}

	// if flow for this run is a missing asset, we have a problem
	if waitingRun.Flow() == nil {
		failSession("can't resume run with missing flow asset")
		return nil
	}

	if s.sprints >= s.engine.Options().MaxSprintsPerSession {
		failSession("reached maximum number of sprints per session (%d)", s.engine.Options().MaxSprintsPerSession)
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
		e.SetStep(step)
		sprint.logEvent(e)
	}

	// resumes can set or clear input
	input := resume.Input(s.assets)
	s.setInput(input)

	if input != nil {
		waitingRun.recordInput()
	}

	isTimeout := false

	switch resume.Type() {
	case resumes.TypeWaitExpiration:
		waitingRun.Exit(flows.RunStatusExpired)
		sprint.logEvent(events.NewRunEnded(waitingRun.UUID(), waitingRun.FlowReference(), flows.RunStatusExpired))
	case resumes.TypeWaitTimeout:
		isTimeout = true
		fallthrough
	default:
		waitingRun.setStatus(flows.RunStatusActive)
	}

	exit, operand, err := s.findResumeExit(sprint, waitingRun, isTimeout)
	if err != nil {
		failSession(fmt.Sprintf("unable to resolve router exit: %s", err.Error()))
		return nil
	}

	// ensure groups are correct
	s.ensureQueryBasedGroups(logEvent)

	// off to the races again...
	return s.continueUntilWait(ctx, sprint, waitingRun, node, exit, operand, step, nil)
}

// finds the exit from a the current node in a run that may have been waiting or a parent paused for a child subflow
func (s *session) findResumeExit(sprint *sprint, run *run, isTimeout bool) (flows.Exit, string, error) {
	// we might have no immediate destination in this run, but continueUntilWait can resume a parent run
	if run.Status() != flows.RunStatusActive {
		return nil, "", nil
	}

	step, node, err := run.PathLocation()
	if err != nil {
		return nil, "", err
	}
	logEvent := func(e flows.Event) {
		e.SetStep(step)
		sprint.logEvent(e)
	}

	// see if this node can now pick a destination
	return s.pickNodeExit(sprint, run, node, step, isTimeout, logEvent)
}

// the main flow execution loop
func (s *session) continueUntilWait(ctx context.Context, sprint *sprint, currentRun *run, node flows.Node, exit flows.Exit, operand string, step flows.Step, trigger flows.Trigger) (err error) {
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
					sprint.logEvent(events.NewRunEnded(run.UUID(), run.FlowReference(), flows.RunStatusCompleted))
				}
			}

			// create a new run for it
			flow := s.pushedFlow.flow
			currentRun = newRun(s, s.pushedFlow.flow, currentRun)
			s.addRun(currentRun)
			sprint.logEvent(events.NewRunStarted(currentRun, s.pushedFlow.terminal))
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
				sprint.logEvent(events.NewRunEnded(currentRun.UUID(), currentRun.FlowReference(), flows.RunStatusCompleted))
			}

			parentRun := currentRun.parent

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
					failRun(sprint, currentRun, step, nil)
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
func (s *session) visitNode(ctx context.Context, sprint *sprint, r *run, node flows.Node, trigger flows.Trigger) (flows.Step, flows.Exit, string, error) {
	step := r.CreateStep(node)
	logEvent := func(e flows.Event) {
		e.SetStep(step)
		sprint.logEvent(e)
	}

	// if this is a new run based on a trigger that provided input, record that on the run
	if trigger != nil && s.input != nil {
		r.recordInput()
	}

	// execute our node's actions
	if node.Actions() != nil {
		for _, action := range node.Actions() {
			if err := action.Execute(ctx, r, step, logEvent); err != nil {
				return step, nil, "", fmt.Errorf("error executing action[type=%s,uuid=%s]: %w", action.Type(), action.UUID(), err)
			}

			// check if this action has errored the run
			if r.Status() == flows.RunStatusFailed {
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
		if wait.Begin(r, logEvent) {
			// mark ouselves as waiting and hand back to
			r.setStatus(flows.RunStatusWaiting)
			s.status = flows.SessionStatusWaiting

			return step, nil, "", nil
		}
	}

	// use our node's router to determine where to go next
	exit, operand, err := s.pickNodeExit(sprint, r, node, step, false, logEvent)
	return step, exit, operand, err
}

// picks the exit to use on the given node
func (s *session) pickNodeExit(sprint *sprint, r *run, node flows.Node, step flows.Step, isTimeout bool, logEvent flows.EventLogger) (flows.Exit, string, error) {
	var exitUUID flows.ExitUUID
	var operand string
	var err error

	if node.Router() != nil {
		if isTimeout {
			exitUUID, err = node.Router().RouteTimeout(r, step, logEvent)
		} else {
			exitUUID, operand, err = node.Router().Route(r, step, logEvent)
		}

		if err != nil {
			return nil, "", fmt.Errorf("error routing from node[uuid=%s]: %w", node.UUID(), err)
		}
		// router didn't error.. but it failed to pick a category
		if exitUUID == "" {
			failRun(sprint, r, step, fmt.Errorf("router on node[uuid=%s] failed to pick a category", node.UUID()))
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
func (s *session) ensureQueryBasedGroups(logEvent flows.EventLogger) {
	if s.contact == nil {
		return
	}

	added, removed := s.contact.ReevaluateQueryBasedGroups(s.Environment())

	// add groups changed event for the groups we were added/removed to/from
	if len(added) > 0 || len(removed) > 0 {
		logEvent(events.NewContactGroupsChanged(added, removed))
	}
}

// utility to fail the current run and log a failRun event
func failRun(sp *sprint, r *run, step flows.Step, err error) {
	if err != nil {
		evt := events.NewFailure(err)
		if step != nil {
			evt.SetStep(step)
		}

		sp.logEvent(evt)
	}

	r.Exit(flows.RunStatusFailed)

	sp.logEvent(events.NewRunEnded(r.UUID(), r.FlowReference(), flows.RunStatusFailed))
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type sessionEnvelope struct {
	UUID        flows.SessionUUID   `json:"uuid"                validate:"required"`
	Type        flows.FlowType      `json:"type"                validate:"required"`
	CreatedOn   time.Time           `json:"created_on"          validate:"required"`
	Trigger     json.RawMessage     `json:"trigger"             validate:"required"`
	ContactUUID flows.ContactUUID   `json:"contact_uuid"        validate:"required,uuid"`
	CallUUID    flows.CallUUID      `json:"call_uuid,omitempty" validate:"omitempty,uuid"`
	Runs        []json.RawMessage   `json:"runs"`
	Status      flows.SessionStatus `json:"status"              validate:"required"`
	Wait        json.RawMessage     `json:"wait,omitempty"`
	Input       json.RawMessage     `json:"input,omitempty"`
	Sprints     int                 `json:"sprints"`
}

// ReadSession decodes a session from the passed in JSON
func readSession(eng flows.Engine, sa flows.SessionAssets, data []byte, env envs.Environment, contact *flows.Contact, call *flows.Call, missing assets.MissingCallback) (flows.Session, error) {
	e := &sessionEnvelope{}
	var err error

	if err = utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, fmt.Errorf("unable to read session: %w", err)
	}

	s := &session{
		engine:    eng,
		assets:    sa,
		uuid:      e.UUID,
		type_:     e.Type,
		createdOn: e.CreatedOn,
		status:    e.Status,
		sprints:   e.Sprints,

		env:        env,
		contact:    contact,
		call:       call,
		runsByUUID: make(map[flows.RunUUID]*run),
	}

	if e.Trigger != nil {
		if s.trigger, err = triggers.Read(s.Assets(), e.Trigger, missing); err != nil {
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
		r, err := readRun(s, e.Runs[i], missing)
		if err != nil {
			return nil, fmt.Errorf("unable to read run %d: %w", i, err)
		}
		s.addRun(r)
	}

	// and our input
	if e.Input != nil {
		if s.input, err = inputs.Read(s.Assets(), e.Input, missing); err != nil {
			return nil, fmt.Errorf("unable to read input: %w", err)
		}
	}

	// older sessions won't have a sprints count but will have events and will have set legacyWaitCount
	if s.sprints == 0 {
		for _, r := range s.runsByUUID {
			s.sprints += r.legacyWaitCount
		}
	}

	return s, nil
}

// MarshalJSON marshals this session into JSON
func (s *session) MarshalJSON() ([]byte, error) {
	e := &sessionEnvelope{
		UUID:      s.uuid,
		Type:      s.type_,
		CreatedOn: s.createdOn,
		Status:    s.status,
		Sprints:   s.sprints,
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
