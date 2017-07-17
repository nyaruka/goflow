package engine

import (
	"encoding/json"
	"fmt"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
)

// StartFlow starts the flow for the passed in contact, returning the created FlowRun
func StartFlow(env flows.FlowEnvironment, flow flows.Flow, contact *flows.Contact, parent flows.FlowRun, callerEvents []flows.Event, extra json.RawMessage) (flows.Session, error) {
	// build our run
	run := flow.CreateRun(env, contact, parent)

	// if we got extra, set it
	if extra != nil {
		run.SetExtra(extra)
	}

	// no first node, nothing to do (valid but weird)
	if len(flow.Nodes()) == 0 {
		run.Exit(flows.StatusCompleted)
		return run.Session(), nil
	}

	// off to the races
	err := continueRunUntilWait(run, flow.Nodes()[0].UUID(), nil, callerEvents)
	return run.Session(), err
}

// ResumeFlow resumes our flow from the last step
func ResumeFlow(env flows.FlowEnvironment, run flows.FlowRun, callerEvents []flows.Event) (flows.Session, error) {
	// to resume a flow, hydrate our run with the environment
	err := run.Hydrate(env)
	if err != nil {
		return run.Session(), err
	}

	// no steps to resume from, nothing to do, return
	if len(run.Path()) == 0 {
		return run.Session(), nil
	}

	// grab the last step
	step := run.Path()[len(run.Path())-1]

	// and the last node
	node := run.Flow().GetNode(step.NodeUUID())
	if node == nil {
		err := fmt.Errorf("cannot resume at node '%s' that no longer exists", step.NodeUUID())
		run.AddError(step, err)
		run.Exit(flows.StatusErrored)
		return run.Session(), nil
	}

	// the first event resumes the wait
	destination, step, err := resumeNode(run, node, step, callerEvents)
	if err != nil {
		return run.Session(), err
	}

	err = continueRunUntilWait(run, destination, step, nil)
	if err != nil {
		return run.Session(), err
	}

	// if we ran to completion and have a parent, resume that flow
	if run.Parent() != nil && run.IsComplete() {
		event := events.NewFlowExitedEvent(run)
		parentRun, err := env.GetRun(run.Parent().UUID())
		if err != nil {
			run.AddError(step, err)
			run.Exit(flows.StatusErrored)
			return run.Session(), nil
		}
		parentRun.SetSession(run.Session())
		return ResumeFlow(env, parentRun, []flows.Event{event})
	}

	return run.Session(), nil
}

// Continues the flow entering the passed in flow
func continueRunUntilWait(run flows.FlowRun, destination flows.NodeUUID, step flows.Step, callerEvents []flows.Event) (err error) {
	// set of uuids we've visited
	visited := make(visitedMap)

	for destination != noDestination {
		// this is a loop, we log it and stop execution
		if visited[destination] {
			err = fmt.Errorf("flow loop detected, stopping execution before entering '%s'", destination)
			break
		}

		node := run.Flow().GetNode(destination)

		if node == nil {
			err = fmt.Errorf("unable to find destination '%s'", destination)
			break
		}

		destination, step, err = enterNode(run, node, callerEvents)

		// only pass our caller events to the first node as it is responsible for handling them
		callerEvents = nil

		// mark this node as visited to prevent loops
		visited[node.UUID()] = true

		// if we have an error, break out
		if err != nil {
			break
		}
	}

	// if we have an error, log it if we have a step
	if err != nil && step != nil {
		run.AddError(step, err)
		run.Exit(flows.StatusErrored)
	}

	// mark ourselves as complete if our run is active and we have no wait
	if run.Wait() == nil && run.Status() == flows.StatusActive {
		run.Exit(flows.StatusCompleted)
	}

	return nil
}

func resumeNode(run flows.FlowRun, node flows.Node, step flows.Step, callerEvents []flows.Event) (flows.NodeUUID, flows.Step, error) {
	wait := node.Wait()

	// it's an error to resume a flow at a wait that no longer exists, error
	if wait == nil {
		err := fmt.Errorf("cannot resume flow at node '%s' which no longer contains wait", node.UUID())
		run.AddError(step, err)
		run.Exit(flows.StatusErrored)
		return noDestination, step, nil
	}

	// try to resume our wait with the first caller event
	err := wait.End(run, step, callerEvents[0])
	if err != nil {
		run.AddError(step, err)
		run.Exit(flows.StatusErrored)
		return noDestination, step, nil
	}

	// we can now apply the caller events
	for _, event := range callerEvents {
		run.ApplyEvent(step, event)
	}

	// determine our exit
	return pickNodeExit(run, node, step)
}

func enterNode(run flows.FlowRun, node flows.Node, callerEvents []flows.Event) (flows.NodeUUID, flows.Step, error) {
	// create our step
	step := run.CreateStep(node)

	// apply any caller events
	for _, event := range callerEvents {
		run.ApplyEvent(step, event)
	}

	// execute our actions
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

	// if we have a wait, execute that
	wait := node.Wait()
	if wait != nil {
		err := wait.Begin(run, step)
		if err != nil {
			run.AddError(step, err)
			run.Exit(flows.StatusErrored)
			return noDestination, step, nil
		}

		// can we end immediately?
		event, err := wait.GetEndEvent(run, step)
		if err != nil {
			run.AddError(step, err)
			run.Exit(flows.StatusErrored)
			return noDestination, step, nil
		}

		// we have to really wait, return out
		if event == nil {
			return noDestination, step, nil
		}

		// end our wait and continue onwards
		err = wait.End(run, step, event)
		if err != nil {
			run.AddError(step, err)
			run.Exit(flows.StatusErrored)
			return noDestination, step, nil
		}
	}

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
		run.ApplyEvent(step, event)
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
