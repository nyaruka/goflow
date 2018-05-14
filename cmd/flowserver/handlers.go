package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/nyaruka/goflow/excellent"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/assets"
	"github.com/nyaruka/goflow/flows/engine"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/flows/triggers"
	"github.com/nyaruka/goflow/legacy"
	"github.com/nyaruka/goflow/utils"
)

type startRequest struct {
	Assets      *json.RawMessage       `json:"assets"`
	AssetServer json.RawMessage        `json:"asset_server" validate:"required"`
	Trigger     *utils.TypedEnvelope   `json:"trigger" validate:"required"`
	Events      []*utils.TypedEnvelope `json:"events"`
	Config      *json.RawMessage       `json:"config"`
}

func (s *FlowServer) handleStart(w http.ResponseWriter, r *http.Request) (interface{}, error) {
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
	err = utils.Validate(start)
	if err != nil {
		return nil, err
	}

	// include any embedded assets
	if start.Assets != nil {
		if err = s.assetCache.Include(*start.Assets); err != nil {
			return nil, err
		}
	}

	// read and validate our asset server
	assetServer, err := assets.ReadAssetServer(s.config.AssetServerToken, s.httpClient, start.AssetServer)
	if err != nil {
		return nil, err
	}

	// build the configuration for this request
	config := s.config.Engine()
	if start.Config != nil {
		config, err = engine.ReadConfig(*start.Config, config)
	}

	// build our session
	session := engine.NewSession(s.assetCache, assetServer, config, s.httpClient)

	// read our trigger
	trigger, err := triggers.ReadTrigger(session, start.Trigger)
	if err != nil {
		return nil, err
	}

	// read our caller events
	callerEvents, err := events.ReadEvents(start.Events)
	if err != nil {
		return nil, err
	}

	// start our flow
	err = session.Start(trigger, callerEvents)
	if err != nil {
		return nil, err
	}

	return &sessionResponse{Session: session, Events: session.Events()}, nil
}

type resumeRequest struct {
	Assets      *json.RawMessage       `json:"assets"`
	AssetServer json.RawMessage        `json:"asset_server" validate:"required"`
	Session     json.RawMessage        `json:"session" validate:"required"`
	Events      []*utils.TypedEnvelope `json:"events" validate:"required,min=1"`
	Config      *json.RawMessage       `json:"config"`
}

func (s *FlowServer) handleResume(w http.ResponseWriter, r *http.Request) (interface{}, error) {
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
	err = utils.Validate(resume)
	if err != nil {
		return nil, err
	}

	// include any embedded assets
	if resume.Assets != nil {
		if err = s.assetCache.Include(*resume.Assets); err != nil {
			return nil, err
		}
	}

	// read and validate our asset server
	assetServer, err := assets.ReadAssetServer(s.config.AssetServerToken, s.httpClient, resume.AssetServer)
	if err != nil {
		return nil, err
	}

	// build the configuration for this request
	config := s.config.Engine()
	if resume.Config != nil {
		config, err = engine.ReadConfig(*resume.Config, config)
	}

	// read our session
	session, err := engine.ReadSession(s.assetCache, assetServer, config, s.httpClient, resume.Session)
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
		return nil, err
	}

	return &sessionResponse{Session: session, Events: session.Events()}, nil
}

type migrateRequest struct {
	Flows     []json.RawMessage `json:"flows"`
	IncludeUI *bool             `json:"include_ui"`
}

func (s *FlowServer) handleMigrate(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	migrate := migrateRequest{}
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		return nil, err
	}

	if err := r.Body.Close(); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &migrate); err != nil {
		return nil, err
	}

	if migrate.Flows == nil {
		return nil, fmt.Errorf("missing flows element")
	}

	legacyFlows, err := legacy.ReadLegacyFlows(migrate.Flows)
	if err != nil {
		return nil, err
	}

	includeUI := migrate.IncludeUI == nil || *migrate.IncludeUI

	flows := make([]flows.Flow, len(legacyFlows))
	for f := range legacyFlows {
		flows[f], err = legacyFlows[f].Migrate(includeUI)
		if err != nil {
			return nil, err
		}
	}

	return flows, err
}

type expressionResponse struct {
	Result string   `json:"result"`
	Errors []string `json:"errors"`
}

type expressionRequest struct {
	Expression string          `json:"expression"`
	Context    json.RawMessage `json:"context"`
}

func (s *FlowServer) handleExpression(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	expression := expressionRequest{}
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		return nil, err
	}
	if err := r.Body.Close(); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(body, &expression); err != nil {
		return nil, err
	}

	if expression.Context == nil || expression.Expression == "" {
		return nil, fmt.Errorf("missing context or expression element")
	}

	context := types.JSONToXValue(expression.Context)

	// evaluate it
	result, err := excellent.EvaluateTemplateAsString(utils.NewDefaultEnvironment(), context, expression.Expression, false, nil)
	if err != nil {
		return expressionResponse{Result: result, Errors: []string{err.Error()}}, nil
	}

	return expressionResponse{Result: result, Errors: []string{}}, nil
}
