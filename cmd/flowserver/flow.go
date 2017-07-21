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
	Contact *flows.Contact   `json:"contact"`
	Session flows.Session    `json:"session"`
	Log     []flows.LogEntry `json:"log"`
}

func (r *flowResponse) MarshalJSON() ([]byte, error) {
	envelope := struct {
		Contact *flows.Contact   `json:"contact"`
		Session flows.Session    `json:"session"`
		Log     []flows.LogEntry `json:"log"`
	}{
		Contact: r.Contact,
		Session: r.Session,
		Log:     r.Session.Log(),
	}

	return json.Marshal(envelope)
}

type startRequest struct {
	Environment *json.RawMessage       `json:"environment"`
	Flows       json.RawMessage        `json:"flows"    validate:"required"`
	Channels    []json.RawMessage      `json:"channels,omit_empty"`
	Contact     json.RawMessage        `json:"contact"  validate:"required"`
	Extra       json.RawMessage        `json:"extra,omitempty"`
	Events      []*utils.TypedEnvelope `json:"events"`
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
	flowlist, err := definition.ReadFlows(start.Flows)
	if err != nil {
		return nil, err
	}
	if len(flowlist) == 0 {
		return nil, utils.NewValidationError("flows: must include at least one flow")
	}

	// read our channels
	channelList, err := flows.ReadChannels(start.Channels)
	if err != nil {
		return nil, err
	}

	// read our contact
	contact, err := flows.ReadContact(start.Contact)
	if err != nil {
		return nil, err
	}

	// build our session environment
	sessionEnv := engine.NewSessionEnvironment(env, flowlist, channelList, []*flows.Contact{contact})

	// read our caller events
	callerEvents, err := events.ReadEvents(start.Events)
	if err != nil {
		return nil, err
	}

	// start our flow
	session, err := engine.StartFlow(sessionEnv, flowlist[0], contact, nil, callerEvents, start.Extra)
	if err != nil {
		return nil, fmt.Errorf("error starting flow: %s", err)
	}

	return &flowResponse{Contact: contact, Session: session, Log: session.Log()}, nil
}

type resumeRequest struct {
	Environment *json.RawMessage       `json:"environment"`
	Flows       json.RawMessage        `json:"flows"       validate:"required,min=1"`
	Channels    []json.RawMessage      `json:"channels,omit_empty"`
	Contact     json.RawMessage        `json:"contact"     validate:"required"`
	Session     json.RawMessage        `json:"session"     validate:"required"`
	Events      []*utils.TypedEnvelope `json:"events"      validate:"required,min=1"`
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

	// read our channels
	channelList, err := flows.ReadChannels(resume.Channels)
	if err != nil {
		return nil, err
	}

	// read our contact
	contact, err := flows.ReadContact(resume.Contact)
	if err != nil {
		return nil, err
	}

	// build our environment
	sessionEnv := engine.NewSessionEnvironment(env, flowList, channelList, []*flows.Contact{contact})

	// read our session
	session, err := runs.ReadSession(sessionEnv, resume.Session)
	if err != nil {
		return nil, err
	}
	if len(session.Runs()) == 0 {
		return nil, utils.NewValidationError("session: must include at least one run")
	}

	// clear the event log if it was passed in
	session.ClearLog()

	// read our new caller events
	callerEvents, err := events.ReadEvents(resume.Events)
	if err != nil {
		return nil, err
	}

	// set our contact on our run
	activeRun := session.ActiveRun()
	if activeRun == nil {
		return nil, utils.NewValidationError("session: no active run to resume")
	}

	// resume our flow
	session, err = engine.ResumeFlow(sessionEnv, activeRun, callerEvents)
	if err != nil {
		return nil, fmt.Errorf("error resuming flow: %s", err)
	}

	return &flowResponse{Contact: contact, Session: session, Log: session.Log()}, nil
}
