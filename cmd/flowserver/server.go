package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/nyaruka/goflow/excellent"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/assets"
	"github.com/nyaruka/goflow/flows/engine"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/flows/triggers"
	"github.com/nyaruka/goflow/legacy"
	"github.com/nyaruka/goflow/utils"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rakyll/statik/fs"
	log "github.com/sirupsen/logrus"
)

type FlowServer struct {
	config     *Config
	httpServer *http.Server
	assetCache *assets.AssetCache
	httpClient *utils.HTTPClient
}

// NewFlowServer creates a new flow server instance
func NewFlowServer(config *Config) *FlowServer {
	r := chi.NewRouter()
	r.Use(middleware.DefaultCompress)
	r.Use(middleware.StripSlashes)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(traceErrors())
	r.Use(requestLogger())
	r.Use(middleware.Timeout(60 * time.Second))

	// no static dir passed in? serve from statik
	var staticDir http.FileSystem
	var err error

	if config.Static == "" {
		staticDir, err = fs.New()
		if err != nil {
			log.Fatal(err)
		}
		log.WithField("comp", "server").Info("using compiled statik assets")
	} else {
		staticDir = http.Dir(config.Static)
		log.WithField("comp", "server").Info("using asset dir: ", config.Static)
	}

	s := &FlowServer{
		config: config,
		httpServer: &http.Server{
			Addr:         fmt.Sprintf(":%d", config.Port),
			Handler:      r,
			ReadTimeout:  60 * time.Second,
			WriteTimeout: 60 * time.Second,
		},
		httpClient: utils.NewHTTPClient("goflow/" + config.Version),
	}

	// root page just serves our example and "postman"" interface
	r.Get("/", templateHandler(staticDir, indexHandler))
	r.Get("/version", jsonHandler(s.handleVersion))
	r.Post("/flow/start", jsonHandler(s.handleStart))
	r.Post("/flow/resume", jsonHandler(s.handleResume))
	r.Post("/flow/migrate", jsonHandler(s.handleMigrate))
	r.Post("/expression", jsonHandler(s.handleExpression))

	r.NotFound(errorHandler(http.StatusNotFound, "not found"))
	r.MethodNotAllowed(errorHandler(http.StatusMethodNotAllowed, "method not allowed"))

	return s
}

// Start starts the flow server
func (s *FlowServer) Start() {
	s.assetCache = assets.NewAssetCache(s.config.AssetCacheSize, s.config.AssetCachePrune)

	go func() {
		err := s.httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.WithFields(log.Fields{
				"comp": "server",
				"err":  err,
			}).Error()
		}
	}()
}

// Stop stops the flow server
func (s *FlowServer) Stop() {
	s.httpServer.Shutdown(context.Background())
	s.assetCache.Shutdown()
}

func (s *FlowServer) handleVersion(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	response := map[string]string{
		"version": version,
	}
	return response, nil
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

type startRequest struct {
	Assets      *json.RawMessage       `json:"assets"`
	AssetServer json.RawMessage        `json:"asset_server" validate:"required"`
	Trigger     *utils.TypedEnvelope   `json:"trigger" validate:"required"`
	Events      []*utils.TypedEnvelope `json:"events"`
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

	// build our session
	session := engine.NewSession(s.assetCache, assetServer, s.config, s.httpClient)

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

	// read our session
	session, err := engine.ReadSession(s.assetCache, assetServer, s.config, s.httpClient, resume.Session)
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
	Flows []json.RawMessage `json:"flows"`
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

	flows := make([]flows.Flow, len(legacyFlows))
	for f := range legacyFlows {
		flows[f], err = legacyFlows[f].Migrate(true)
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
