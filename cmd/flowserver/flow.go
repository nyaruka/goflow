package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/definition"
	"github.com/nyaruka/goflow/flows/engine"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/flows/runs"
	"github.com/nyaruka/goflow/utils"
)

type flowResponse struct {
	Contact   *flows.Contact  `json:"contact"`
	RunOutput flows.RunOutput `json:"run_output"`
}

type startRequest struct {
	Flows   json.RawMessage `json:"flows"`
	Contact json.RawMessage `json:"contact"`
}

func handleStart(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	start := startRequest{}
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		return nil, err
	}
	if err := r.Body.Close(); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(body, &start); err != nil {
		return nil, err
	}

	if start.Flows == nil || start.Contact == nil {
		return nil, fmt.Errorf("missing contact or flows element")
	}

	// read our flows
	startFlows, err := definition.ReadFlows(start.Flows)
	if err != nil {
		return nil, fmt.Errorf("error parsing flows: %s", err)
	}

	// read our contact
	contact, err := flows.ReadContact(start.Contact)
	if err != nil {
		return nil, fmt.Errorf("error parsing contact: %s", err)
	}

	// build our environment
	env := engine.NewFlowEnvironment(utils.NewDefaultEnvironment(), startFlows, []flows.FlowRun{}, []*flows.Contact{contact})

	// start our flow
	output, err := engine.StartFlow(env, startFlows[0], contact, nil)
	if err != nil {
		return nil, fmt.Errorf("error starting flow: %s", err)
	}

	return &flowResponse{Contact: contact, RunOutput: output}, nil
}

type resumeRequest struct {
	Contact   json.RawMessage      `json:"contact"`
	Flows     json.RawMessage      `json:"flows"`
	RunOutput json.RawMessage      `json:"run_output"`
	Event     *utils.TypedEnvelope `json:"event"`
}

func handleResume(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	resume := resumeRequest{}
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		return nil, err
	}
	if err := r.Body.Close(); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(body, &resume); err != nil {
		return nil, err
	}

	if resume.Flows == nil || resume.RunOutput == nil || resume.Event == nil || resume.Contact == nil {
		return nil, err
	}

	// read our flows
	flowList, err := definition.ReadFlows(resume.Flows)
	if err != nil {
		return nil, fmt.Errorf("error parsing flows: %s", err)
	}

	// read our run
	runOutput, err := runs.ReadRunOutput(resume.RunOutput)
	if err != nil {
		return nil, fmt.Errorf("error parsing run output: %s", err)
	}

	// our contact
	contact, err := flows.ReadContact(resume.Contact)
	if err != nil {
		return nil, fmt.Errorf("error parsing run contact: %s", err)
	}

	// and our event
	event, err := events.EventFromEnvelope(resume.Event)
	if err != nil {
		return nil, fmt.Errorf("error reading event: %s", err)
	}

	// build our environment
	env := engine.NewFlowEnvironment(utils.NewDefaultEnvironment(), flowList, runOutput.Runs(), []*flows.Contact{contact})

	// hydrate all our runs
	for _, run := range runOutput.Runs() {
		run.Hydrate(env)
	}

	// set our contact on our run
	activeRun := runOutput.ActiveRun()

	// resume our flow
	output, err := engine.ResumeFlow(env, activeRun, event)
	if err != nil {
		return nil, fmt.Errorf("Error resuming flow: %s", err)
	}

	return &flowResponse{Contact: contact, RunOutput: output}, nil
}
