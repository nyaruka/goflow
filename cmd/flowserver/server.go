package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/nyaruka/goflow/utils"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	log "github.com/sirupsen/logrus"
)

// FlowServer exposes several engine functions as an HTTP service
type FlowServer struct {
	config     *Config
	httpServer *http.Server
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
	r.Post("/flow/migrate", jsonHandler(s.handleMigrate))
	r.Get("/version", jsonHandler(s.handleVersion))

	r.NotFound(errorHandler(http.StatusNotFound, "not found"))
	r.MethodNotAllowed(errorHandler(http.StatusMethodNotAllowed, "method not allowed"))

	return s
}

// Start starts the flow server
func (s *FlowServer) Start() {
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
}
