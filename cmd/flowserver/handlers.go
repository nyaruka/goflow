package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/nyaruka/goflow/assets/rest"
	"github.com/nyaruka/goflow/excellent"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/engine"
	"github.com/nyaruka/goflow/flows/resumes"
	"github.com/nyaruka/goflow/flows/triggers"
	"github.com/nyaruka/goflow/legacy"
	"github.com/nyaruka/goflow/utils"

	"github.com/pkg/errors"
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
	Session flows.Session `json:"session"`
	Events  []flows.Event `json:"events"`
}

// Starts a new engine session
//
//   {
//     "assets": [...],
//     "asset_server": {...},
//     "trigger": {...}
//   }
//
type startRequest struct {
	sessionRequest

	Trigger json.RawMessage `json:"trigger" validate:"required"`
}

// reads the assets and asset_server section of a request
func (s *FlowServer) readAssets(request *sessionRequest, cache *rest.AssetCache) (*rest.ServerSource, error) {
	// include any embedded assets
	if request.Assets != nil {
		if err := s.assetCache.Include(*request.Assets); err != nil {
			return nil, err
		}
	}

	// read and validate our asset server
	return rest.ReadServerSource(s.config.AssetServerToken, s.httpClient, cache, request.AssetServer)
}

// handles a request to /start
func (s *FlowServer) handleStart(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	request := &startRequest{}
	if err := utils.UnmarshalAndValidateWithLimit(r.Body, request, maxRequestBytes); err != nil {
		return nil, err
	}

	assetServer, err := s.readAssets(&request.sessionRequest, s.assetCache)
	if err != nil {
		return nil, err
	}

	// build the configuration for this request
	config := s.config.Engine()
	if request.Config != nil {
		config, err = engine.ReadConfig(*request.Config, config)
	}

	// build our session
	assets, err := engine.NewSessionAssets(assetServer)
	if err != nil {
		return nil, err
	}

	session := engine.NewSession(assets, config, s.httpClient)

	// read our trigger
	trigger, err := triggers.ReadTrigger(session.Assets(), request.Trigger)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read trigger")
	}

	// start our flow
	newEvents, err := session.Start(trigger)
	if err != nil {
		return nil, err
	}

	return &sessionResponse{Session: session, Events: newEvents}, nil
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

	Session json.RawMessage `json:"session" validate:"required"`
	Resume  json.RawMessage `json:"resume" validate:"required"`
}

func (s *FlowServer) handleResume(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	request := &resumeRequest{}
	if err := utils.UnmarshalAndValidateWithLimit(r.Body, request, maxRequestBytes); err != nil {
		return nil, err
	}

	assetServer, err := s.readAssets(&request.sessionRequest, s.assetCache)
	if err != nil {
		return nil, err
	}

	// build the configuration for this request
	config := s.config.Engine()
	if request.Config != nil {
		config, err = engine.ReadConfig(*request.Config, config)
	}

	// read our session
	assets, err := engine.NewSessionAssets(assetServer)
	if err != nil {
		return nil, err
	}

	session, err := engine.ReadSession(assets, config, s.httpClient, request.Session)
	if err != nil {
		return nil, err
	}

	// read our resume
	resume, err := resumes.ReadResume(session, request.Resume)
	if err != nil {
		return nil, err
	}

	// resume our session
	newEvents, err := session.Resume(resume)
	if err != nil {
		return nil, err
	}

	return &sessionResponse{Session: session, Events: newEvents}, nil
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
		return nil, errors.Errorf("missing flow element")
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
		return nil, errors.Errorf("missing context or expression element")
	}

	context := types.JSONToXValue(expression.Context)

	// evaluate it
	result, err := excellent.EvaluateTemplateAsString(utils.NewDefaultEnvironment(), context, expression.Expression, nil)
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
