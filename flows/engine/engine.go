package engine

import (
	"encoding/json"
	"fmt"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
)

type VisitedMap map[flows.NodeUUID]bool

const noDestination = flows.NodeUUID("")

// StartFlow starts the flow for the passed in contact, returning the created FlowRun
func StartFlow(env flows.FlowEnvironment, flow flows.Flow, contact *flows.Contact, parent flows.FlowRun, input flows.Input, extra json.RawMessage) (flows.Session, error) {
	// build our run
	run := flow.CreateRun(env, contact, parent)

	// if we got an input, set it
	if input != nil {
		run.SetInput(input)
	}

	// if we got extra, set it
	if extra != nil {
		run.SetExtra(extra)
	}

	// no first node, nothing to do (valid but weird)
	if len(flow.Nodes()) == 0 {
		run.Exit(flows.RunCompleted)
		return run.Session(), nil
	}

	initTranslations(run)

	// off to the races
	err := continueRunUntilWait(run, flow.Nodes()[0].UUID(), nil, input)
	return run.Session(), err
}

// ResumeFlow resumes our flow from the last step
func ResumeFlow(env flows.FlowEnvironment, run flows.FlowRun, event flows.Event) (flows.Session, error) {
	// to resume a flow, hydrate our run with the environment
	run.Hydrate(env)

	// no steps to resume from, nothing to do, return
	if len(run.Path()) == 0 {
		return run.Session(), nil
	}

	initTranslations(run)

	// grab the last step
	step := run.Path()[len(run.Path())-1]

	// and the last node
	node := run.Flow().GetNode(step.NodeUUID())
	if node == nil {
		err := fmt.Errorf("cannot resume at node '%s' that no longer exists", step.NodeUUID())
		run.AddError(step, err)
		return run.Session(), nil
	}

	destination, step, err := resumeNode(run, node, step, event)
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
			return run.Session(), nil
		}
		parentRun.SetSession(run.Session())
		return ResumeFlow(env, parentRun, event)
	}

	return run.Session(), nil
}

// initializes our context based on our flow and current context
func initTranslations(run flows.FlowRun) {
	// set our language based on our contact if we have one
	contact := run.Contact()
	if contact != nil {
		run.SetLanguage(contact.Language())
	} else {
		run.SetLanguage(run.Flow().Language())
	}

	// set the translations on our context
	run.SetFlowTranslations(run.Flow().Translations())
}

// Continues the flow entering the passed in flow
func continueRunUntilWait(run flows.FlowRun, destination flows.NodeUUID, step flows.Step, event flows.Event) (err error) {
	// set of uuids we've visited
	visited := make(VisitedMap)

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

		destination, step, err = enterNode(run, node, event)

		// only pass our event to the first node, it is in charge of logging it
		event = nil

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
	}

	// we are complete if we had an error or our run is active and we have no wait
	if err != nil || (run.Wait() == nil && run.Status() == flows.RunActive) {
		run.Exit(flows.RunCompleted)
	}

	return nil
}

func resumeNode(run flows.FlowRun, node flows.Node, step flows.Step, event flows.Event) (flows.NodeUUID, flows.Step, error) {
	wait := node.Wait()

	// it's an error to resume a flow at a wait that no longer exists, error
	if wait == nil {
		err := fmt.Errorf("cannot resume flow at node '%s' which no longer contains wait", node.UUID())
		run.AddError(step, err)
		return noDestination, step, nil
	}

	err := wait.End(run, step, event)
	if err != nil {
		run.AddError(step, err)
		return noDestination, step, nil
	}

	// determine our exit
	return pickNodeExit(run, node, step)
}

func enterNode(run flows.FlowRun, node flows.Node, event flows.Event) (flows.NodeUUID, flows.Step, error) {
	// create our step
	step := run.CreateStep(node)

	// log our entry event if we have one
	if event != nil {
		run.AddEvent(step, event)
	}

	// execute our actions
	if node.Actions() != nil {
		for _, action := range node.Actions() {
			err := action.Execute(run, step)
			if err != nil {
				run.AddError(step, err)
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
			return noDestination, step, nil
		}

		// can we end immediately?
		event, err := wait.GetEndEvent(run, step)
		if err != nil {
			run.AddError(step, err)
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
		return noDestination, step, err
	}

	// look up our actual exit
	if exitUUID != "" {
		// find our exit
		for _, e := range node.Exits() {
			if e.UUID() == exitUUID {
				exitName = e.Name()
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
		event := events.NewSaveFlowResult(node.UUID(), router.ResultName(), route.Match(), exitName)
		run.AddEvent(step, event)
		run.Results().Save(node.UUID(), router.ResultName(), route.Match(), exitName, *event.CreatedOn())
	}

	// log any error we received
	if err != nil {
		run.AddError(step, err)
	}

	// no exit? return no destination
	if exit == nil {
		return noDestination, step, nil
	}

	return exit.DestinationNodeUUID(), step, nil
}

func GetFlow(uuid flows.FlowUUID) (flows.Flow, error) {
	return nil, nil
}
