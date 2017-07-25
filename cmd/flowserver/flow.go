package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/engine"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

type flowResponse struct {
	Session flows.Session    `json:"session"`
	Log     []flows.LogEntry `json:"log"`
}

func (r *flowResponse) MarshalJSON() ([]byte, error) {
	envelope := struct {
		Session flows.Session    `json:"session"`
		Log     []flows.LogEntry `json:"log"`
	}{
		Session: r.Session,
		Log:     r.Session.Log(),
	}

	return json.Marshal(envelope)
}

type startRequest struct {
	Assets json.RawMessage        `json:"assets"           validate:"required"`
	Flow   flows.FlowUUID         `json:"flow_uuid"        validate:"required"`
	Extra  json.RawMessage        `json:"extra,omitempty"`
	Events []*utils.TypedEnvelope `json:"events"`
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

	// read our assets
	assets, err := engine.ReadAssets(start.Assets)
	if err != nil {
		return nil, err
	}

	// build our session
	session := engine.NewSession(assets)

	// read our caller events
	callerEvents, err := events.ReadEvents(start.Events)
	if err != nil {
		return nil, err
	}

	// start our flow
	err = session.StartFlow(start.Flow, nil, callerEvents, start.Extra)
	if err != nil {
		return nil, fmt.Errorf("error starting flow: %s", err)
	}

	return &flowResponse{Session: session, Log: session.Log()}, nil
}

type resumeRequest struct {
	Assets  json.RawMessage        `json:"assets"`
	Session json.RawMessage        `json:"session" validate:"required"`
	Events  []*utils.TypedEnvelope `json:"events"  validate:"required,min=1"`
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

	// read our assets
	assets, err := engine.ReadAssets(resume.Assets)
	if err != nil {
		return nil, err
	}

	// read our session
	session, err := engine.ReadSession(assets, resume.Session)
	if err != nil {
		return nil, err
	}

	// read our new caller events
	callerEvents, err := events.ReadEvents(resume.Events)
	if err != nil {
		return nil, err
	}

	// resume our flow
	err = session.Resume(callerEvents)
	if err != nil {
		return nil, fmt.Errorf("error resuming flow: %s", err)
	}

	return &flowResponse{Session: session, Log: session.Log()}, nil
}
