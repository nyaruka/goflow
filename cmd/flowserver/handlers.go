package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/excellent"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/engine"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/flows/triggers"
	"github.com/nyaruka/goflow/legacy"
	"github.com/nyaruka/goflow/utils"
)

const (
	maxRequestBytes int64 = 1048576
)

type sessionRequest struct {
	Assets      *json.RawMessage `json:"assets"`
	AssetServer json.RawMessage  `json:"asset_server" validate:"required"`
	Config      *json.RawMessage `json:"config"`
}

type sessionResponse struct {
	Session flows.Session
	Events  []flows.Event
}

// MarshalJSON marshals this session response into JSON
func (r *sessionResponse) MarshalJSON() ([]byte, error) {
	eventEnvelopes, err := events.EventsToEnvelopes(r.Session.Events())
	if err != nil {
		return nil, err
	}
	envelope := struct {
		Session flows.Session          `json:"session"`
		Events  []*utils.TypedEnvelope `json:"events"`
	}{
		Session: r.Session,
		Events:  eventEnvelopes,
	}

	return utils.JSONMarshal(envelope)
}

// Starts a new engine session
//
//   {
//     "assets": [...],
//     "asset_server": {...},
//     "trigger": {...},
//     "events": [...]
//   }
//
type startRequest struct {
	sessionRequest

	Trigger *utils.TypedEnvelope   `json:"trigger" validate:"required"`
	Events  []*utils.TypedEnvelope `json:"events"`
}

// reads the assets and asset_server section of a request
func (s *FlowServer) readAssets(request *sessionRequest, cache *assets.AssetCache) (*assets.AssetServer, error) {
	// include any embedded assets
	if request.Assets != nil {
		if err := s.assetCache.Include(*request.Assets); err != nil {
			return nil, err
		}
	}

	// read and validate our asset server
	return assets.ReadAssetServer(s.config.AssetServerToken, s.httpClient, cache, request.AssetServer)
}

// handles a request to /start
func (s *FlowServer) handleStart(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	start := &startRequest{}
	if err := utils.UnmarshalAndValidateWithLimit(r.Body, start, maxRequestBytes); err != nil {
		return nil, err
	}

	assetServer, err := s.readAssets(&start.sessionRequest, s.assetCache)
	if err != nil {
		return nil, err
	}

	// build the configuration for this request
	config := s.config.Engine()
	if start.Config != nil {
		config, err = engine.ReadConfig(*start.Config, config)
	}

	// build our session
	assets, err := engine.NewSessionAssets(assetServer)
	if err != nil {
		return nil, err
	}

	session := engine.NewSession(assets, config, s.httpClient)

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

// Resumes an existing engine session
//
//   {
//     "assets": [...],
//     "asset_server": {...},
//     "session": {"uuid": "468621a8-32e6-4cd2-afc1-04416f7151f0", "runs": [...], ...},
//     "events": [...]
//   }
//
type resumeRequest struct {
	sessionRequest

	Session json.RawMessage        `json:"session" validate:"required"`
	Events  []*utils.TypedEnvelope `json:"events" validate:"required,min=1"`
}

func (s *FlowServer) handleResume(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	resume := &resumeRequest{}
	if err := utils.UnmarshalAndValidateWithLimit(r.Body, resume, maxRequestBytes); err != nil {
		return nil, err
	}

	assetServer, err := s.readAssets(&resume.sessionRequest, s.assetCache)
	if err != nil {
		return nil, err
	}

	// build the configuration for this request
	config := s.config.Engine()
	if resume.Config != nil {
		config, err = engine.ReadConfig(*resume.Config, config)
	}

	// read our session
	assets, err := engine.NewSessionAssets(assetServer)
	if err != nil {
		return nil, err
	}

	session, err := engine.ReadSession(assets, config, s.httpClient, resume.Session)
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

// Migrates a legacy flow to the new flow definition specification
//
//   {
//     "flow": {"uuid": "468621a8-32e6-4cd2-afc1-04416f7151f0", "action_sets": [], ...},
//     "include_ui": false
//   }
//
type migrateRequest struct {
	Flow          json.RawMessage `json:"flow"`
	CollapseExits *bool           `json:"collapse_exits"`
	IncludeUI     *bool           `json:"include_ui"`
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

	if migrate.Flow == nil {
		return nil, fmt.Errorf("missing flow element")
	}

	legacyFlow, err := legacy.ReadLegacyFlow(migrate.Flow)
	if err != nil {
		return nil, err
	}

	collapseExits := migrate.CollapseExits == nil || *migrate.CollapseExits
	includeUI := migrate.IncludeUI == nil || *migrate.IncludeUI

	return legacyFlow.Migrate(collapseExits, includeUI)
}

// Evaluates an expression
//
//   {
//     "expression": "@(upper(foo.bar))",
//     "context": {"foo": {"bar": "Hello!"}}
//   }
//
type expressionRequest struct {
	Expression string          `json:"expression"`
	Context    json.RawMessage `json:"context"`
}

type expressionResponse struct {
	Result string   `json:"result"`
	Errors []string `json:"errors"`
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

// Returns the current version number
func (s *FlowServer) handleVersion(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	response := map[string]string{
		"version": version,
	}
	return response, nil
}
