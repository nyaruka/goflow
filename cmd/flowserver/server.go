package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/nyaruka/goflow/assets/rest"
	"github.com/nyaruka/goflow/utils"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rakyll/statik/fs"
	log "github.com/sirupsen/logrus"
)

// FlowServer exposes several engine functions as an HTTP service
type FlowServer struct {
	config     *Config
	httpServer *http.Server
	assetCache *rest.AssetCache
	httpClient *utils.HTTPClient
}

// NewFlowServer creates a new flow server instance
func NewFlowServer(config *Config) *FlowServer {
	r := chi.NewRouter()
	r.Use(middleware.DefaultCompress)
	r.Use(middleware.StripSlashes)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(panicRecovery)
	r.Use(requestLogger)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(corsAllowedOrigins(config))

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

	// set up the routes and handlers
	r.Get("/", templateHandler(staticDir, indexHandler))
	r.Post("/flow/start", jsonHandler(s.handleStart))
	r.Post("/flow/resume", jsonHandler(s.handleResume))
	r.Post("/flow/migrate", jsonHandler(s.handleMigrate))
	r.Post("/expression", jsonHandler(s.handleExpression))
	r.Get("/version", jsonHandler(s.handleVersion))

	r.NotFound(errorHandler(http.StatusNotFound, "not found"))
	r.MethodNotAllowed(errorHandler(http.StatusMethodNotAllowed, "method not allowed"))

	return s
}

// Start starts the flow server
func (s *FlowServer) Start() {
	s.assetCache = rest.NewAssetCache(s.config.AssetCacheSize, s.config.AssetCachePrune)

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
