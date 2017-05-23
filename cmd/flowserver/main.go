package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	validator "gopkg.in/go-playground/validator.v9"

	"errors"

	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
	"github.com/pressly/lg"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()

	lg.RedirectStdlogOutput(logger)
	lg.DefaultLogger = logger

	r := chi.NewRouter()

	r.Use(middleware.DefaultCompress)
	r.Use(middleware.StripSlashes)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(lg.RequestLogger(logger))
	r.Use(middleware.Heartbeat("/ping"))
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Post("/flow/start", jsonHandler(handleStart))
	r.Post("/flow/resume", jsonHandler(handleResume))
	r.Post("/flow/migrate", jsonHandler(handleMigrate))
	r.Post("/expression", jsonHandler(handleExpression))

	r.NotFound(errorHandler(http.StatusNotFound, "not found"))
	r.MethodNotAllowed(errorHandler(http.StatusMethodNotAllowed, "method not allowed"))

	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", 8080),
		Handler:      r,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	go func() {
		err := httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logrus.WithFields(logrus.Fields{
				"comp": "server",
				"err":  err,
			}).Error()
		}
	}()
	logrus.WithField("comp", "server").WithField("port", "8080").Info("listening")

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	logrus.WithField("comp", "server").WithField("signal", <-ch).Info("stopping")
	httpServer.Shutdown(context.Background())
}

type errorResponse struct {
	Text []string `json:"errors"`
}

// writeError writes a JSON response for the passed in error
func writeError(w http.ResponseWriter, r *http.Request, status int, err error) error {
	lg.Log(r.Context()).WithError(err).Error()
	errors := []string{err.Error()}

	vErrs, isValidation := err.(validator.ValidationErrors)
	if isValidation {
		status = http.StatusBadRequest
		errors = []string{}
		for i := range vErrs {
			errors = append(errors, fmt.Sprintf("field '%s' %s", strings.ToLower(vErrs[i].Field()), vErrs[i].Tag()))
		}
	}
	return writeJSONResponse(w, status, &errorResponse{errors})
}

func writeJSONResponse(w http.ResponseWriter, statusCode int, response interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(response)
}

type jsonHandlerFunc func(http.ResponseWriter, *http.Request) (interface{}, error)

func jsonHandler(handler jsonHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		value, err := handler(w, r)
		if err != nil {
			writeError(w, r, http.StatusInternalServerError, err)
		} else {
			err := writeJSONResponse(w, http.StatusOK, value)
			if err != nil {
				lg.Log(r.Context()).WithError(err).Error()
			}
		}
	}
}

func errorHandler(status int, err string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeError(w, r, status, errors.New(err))
	}
}
