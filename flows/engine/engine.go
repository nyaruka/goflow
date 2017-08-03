package engine

import (
	"fmt"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
)

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
	if run.Session().Wait() == nil && run.Status() == flows.StatusActive {
		run.Exit(flows.StatusCompleted)
	}

	return nil
}

func enterNode(run flows.FlowRun, node flows.Node, callerEvents []flows.Event) (flows.NodeUUID, flows.Step, error) {
	// create our step
	step := run.CreateStep(node)

	// apply any caller events
	for _, event := range callerEvents {
		run.ApplyEvent(step, nil, event)
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

	// if we have a wait, apply that
	wait := node.Wait()
	if wait != nil {
		wait.Apply(run, step)

		run.SetStatus(flows.StatusWaiting)
		run.Session().SetWait(wait)

		return noDestination, step, nil
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
