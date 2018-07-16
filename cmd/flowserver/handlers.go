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

type sessionRequest struct {
	Assets      *json.RawMessage `json:"assets"`
	AssetServer json.RawMessage  `json:"asset_server" validate:"required"`
	Config      *json.RawMessage `json:"config"`
}

type startRequest struct {
	sessionRequest

	Trigger *utils.TypedEnvelope   `json:"trigger" validate:"required"`
	Events  []*utils.TypedEnvelope `json:"events"`
}

type resumeRequest struct {
	sessionRequest

	Session json.RawMessage        `json:"session" validate:"required"`
	Events  []*utils.TypedEnvelope `json:"events" validate:"required,min=1"`
}

// reads the assets and asset_server section of a request
func (s *FlowServer) readAssets(request *sessionRequest) (assets.AssetServer, error) {
	// include any embedded assets
	if request.Assets != nil {
		if err := s.assetCache.Include(*request.Assets); err != nil {
			return nil, err
		}
	}

	// read and validate our asset server
	return assets.ReadAssetServer(s.config.AssetServerToken, s.httpClient, request.AssetServer)
}

// reads the configuration section of a request
func (s *FlowServer) readConfig(request *sessionRequest) (flows.EngineConfig, error) {
	var configData json.RawMessage
	if request.Config != nil {
		configData = *request.Config
	} else {
		configData = json.RawMessage(`{}`)
	}

	return engine.ReadConfig(configData, s.config.EngineDefaults())
}

// handles a request to /start
func (s *FlowServer) handleStart(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	start := &startRequest{}
	if err := unmarshalWithLimit(r.Body, start); err != nil {
		return nil, err
	}

	assetServer, err := s.readAssets(&start.sessionRequest)
	if err != nil {
		return nil, err
	}

	// build the configuration for this request
	config, err := s.readConfig(&start.sessionRequest)
	if err != nil {
		return nil, err
	}

	// build our session
	session := engine.NewSession(s.assetCache, assetServer, config, s.httpClient)

	// read our trigger
	trigger, err := triggers.ReadTrigger(session, start.Trigger)
	if err != nil {
		return nil, fmt.Errorf("unable to read trigger[type=%s]: %s", start.Trigger.Type, err)
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

// handles a request to /resume
func (s *FlowServer) handleResume(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	resume := &resumeRequest{}
	if err := unmarshalWithLimit(r.Body, resume); err != nil {
		return nil, err
	}

	assetServer, err := s.readAssets(&resume.sessionRequest)
	if err != nil {
		return nil, err
	}

	// build the configuration for this request
	config, err := s.readConfig(&resume.sessionRequest)
	if err != nil {
		return nil, err
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

// utility method to read and unmarsmal with a limit on how many bytes can be read
func unmarshalWithLimit(reader io.ReadCloser, s interface{}) error {
	body, err := ioutil.ReadAll(io.LimitReader(reader, 1048576))
	if err != nil {
		return err
	}
	if err := reader.Close(); err != nil {
		return err
	}
	if err := json.Unmarshal(body, &s); err != nil {
		return err
	}

	// validate the request
	return utils.Validate(s)
}
