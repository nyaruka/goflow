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
	"github.com/nyaruka/goflow/flows/inputs"
	"github.com/nyaruka/goflow/flows/runs"
	"github.com/nyaruka/goflow/utils"
)

type flowResponse struct {
	Contact *flows.Contact `json:"contact"`
	Session flows.Session  `json:"session"`
}

type startRequest struct {
	Flows   json.RawMessage      `json:"flows"    validate:"required"`
	Contact json.RawMessage      `json:"contact"  validate:"required"`
	Extra   json.RawMessage      `json:"extra,omitempty"`
	Input   *utils.TypedEnvelope `json:"input"`
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

	// validate our input
	err = utils.ValidateAll(start)
	if err != nil {
		return nil, err
	}

	// read our flows
	startFlows, err := definition.ReadFlows(start.Flows)
	if err != nil {
		return nil, err
	}
	if len(startFlows) == 0 {
		return nil, utils.NewValidationError("flows: must include at least one flow")
	}

	// read our contact
	contact, err := flows.ReadContact(start.Contact)
	if err != nil {
		return nil, err
	}

	// read our input
	var input flows.Input
	if start.Input != nil {
		input, err = inputs.InputFromEnvelope(start.Input)
		if err != nil {
			return nil, err
		}
	}

	// build our environment
	env := engine.NewFlowEnvironment(utils.NewDefaultEnvironment(), startFlows, []flows.FlowRun{}, []*flows.Contact{contact})

	// start our flow
	session, err := engine.StartFlow(env, startFlows[0], contact, nil, input, start.Extra)
	if err != nil {
		return nil, fmt.Errorf("error starting flow: %s", err)
	}

	return &flowResponse{Contact: contact, Session: session}, nil
}

type resumeRequest struct {
	Flows   json.RawMessage      `json:"flows"       validate:"required,min=1"`
	Contact json.RawMessage      `json:"contact"     validate:"required"`
	Session json.RawMessage      `json:"session"     validate:"required"`
	Event   *utils.TypedEnvelope `json:"event"       validate:"required"`
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

	// validate our input
	err = utils.ValidateAll(resume)
	if err != nil {
		return nil, err
	}

	// read our flows
	flowList, err := definition.ReadFlows(resume.Flows)
	if err != nil {
		return nil, err
	}
	if len(flowList) == 0 {
		return nil, utils.NewValidationError("flows: must include at least one flow")
	}

	// read our session
	session, err := runs.ReadSession(resume.Session)
	if err != nil {
		return nil, err
	}
	if len(session.Runs()) == 0 {
		return nil, utils.NewValidationError("session: must include at least one run")
	}

	// our contact
	contact, err := flows.ReadContact(resume.Contact)
	if err != nil {
		return nil, err
	}

	// and our event
	event, err := events.EventFromEnvelope(resume.Event)
	if err != nil {
		return nil, err
	}

	// build our environment
	env := engine.NewFlowEnvironment(utils.NewDefaultEnvironment(), flowList, session.Runs(), []*flows.Contact{contact})

	// hydrate all our runs
	for _, run := range session.Runs() {
		run.Hydrate(env)
	}

	// set our contact on our run
	activeRun := session.ActiveRun()
	if activeRun == nil {
		return nil, utils.NewValidationError("session: no active run to resume")
	}

	// resume our flow
	session, err = engine.ResumeFlow(env, activeRun, event)
	if err != nil {
		return nil, fmt.Errorf("error resuming flow: %s", err)
	}

	return &flowResponse{Contact: contact, Session: session}, nil
}
