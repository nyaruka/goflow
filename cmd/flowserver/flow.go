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
	"github.com/nyaruka/goflow/flows/inputs"
	"github.com/nyaruka/goflow/flows/runs"
	"github.com/nyaruka/goflow/utils"
)

type flowResponse struct {
	Contact *flows.Contact `json:"contact"`
	Session flows.Session  `json:"session"`
	Events  []flows.Event  `json:"events"`
}

func (r *flowResponse) MarshalJSON() ([]byte, error) {
	envelope := struct {
		Contact *flows.Contact         `json:"contact"`
		Session flows.Session          `json:"session"`
		Events  []*utils.TypedEnvelope `json:"events"`
	}{
		Contact: r.Contact,
		Session: r.Session,
	}

	envelope.Events = make([]*utils.TypedEnvelope, len(r.Events))
	var err error
	for i := range r.Events {
		envelope.Events[i], err = utils.EnvelopeFromTyped(r.Events[i])
		if err != nil {
			return nil, err
		}
	}
	return json.Marshal(envelope)
}

type startRequest struct {
	Environment *json.RawMessage     `json:"environment"`
	Flows       json.RawMessage      `json:"flows"    validate:"required"`
	Contact     json.RawMessage      `json:"contact"  validate:"required"`
	Extra       json.RawMessage      `json:"extra,omitempty"`
	Input       *utils.TypedEnvelope `json:"input"`
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

	// read our base environment
	var env utils.Environment
	if start.Environment != nil {
		env, err = utils.ReadEnvironment(start.Environment)
		if err != nil {
			return nil, err
		}
	} else {
		env = utils.NewDefaultEnvironment()
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

	// build our flow environment
	flowEnv := engine.NewFlowEnvironment(env, startFlows, []flows.FlowRun{}, []*flows.Contact{contact})

	// read our input
	var input flows.Input
	if start.Input != nil {
		input, err = inputs.InputFromEnvelope(flowEnv, start.Input)
		if err != nil {
			return nil, err
		}
	}

	// start our flow
	session, err := engine.StartFlow(flowEnv, startFlows[0], contact, nil, input, start.Extra)
	if err != nil {
		return nil, fmt.Errorf("error starting flow: %s", err)
	}

	return &flowResponse{Contact: contact, Session: session, Events: session.Events()}, nil
}

type resumeRequest struct {
	Environment *json.RawMessage     `json:"environment"`
	Flows       json.RawMessage      `json:"flows"       validate:"required,min=1"`
	Contact     json.RawMessage      `json:"contact"     validate:"required"`
	Session     json.RawMessage      `json:"session"     validate:"required"`
	Input       *utils.TypedEnvelope `json:"input"       validate:"required"`
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

	// read our base environment
	var env utils.Environment
	if resume.Environment != nil {
		env, err = utils.ReadEnvironment(resume.Environment)
		if err != nil {
			return nil, err
		}
	} else {
		env = utils.NewDefaultEnvironment()
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

	// clear events if they passed them in
	session.ClearEvents()

	// our contact
	contact, err := flows.ReadContact(resume.Contact)
	if err != nil {
		return nil, err
	}

	// build our flow environment
	flowEnv := engine.NewFlowEnvironment(env, flowList, session.Runs(), []*flows.Contact{contact})

	// finally read our input
	input, err := inputs.InputFromEnvelope(flowEnv, resume.Input)
	if err != nil {
		return nil, err
	}

	// hydrate all our runs
	for _, run := range session.Runs() {
		err = run.Hydrate(flowEnv)
		if err != nil {
			return nil, utils.NewValidationError(err.Error())
		}
	}

	// set our contact on our run
	activeRun := session.ActiveRun()
	if activeRun == nil {
		return nil, utils.NewValidationError("session: no active run to resume")
	}

	// set the input on our active run and convert to an event
	event := input.Event(activeRun)
	activeRun.SetInput(input)

	// resume our flow
	session, err = engine.ResumeFlow(flowEnv, activeRun, event)
	if err != nil {
		return nil, fmt.Errorf("error resuming flow: %s", err)
	}

	return &flowResponse{Contact: contact, Session: session, Events: session.Events()}, nil
}
