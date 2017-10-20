//go:generate statik -src=./static
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/pressly/chi/middleware"
	"github.com/pressly/lg"
	"github.com/sirupsen/logrus"

	_ "github.com/nyaruka/goflow/cmd/flowserver/statik"
	"github.com/nyaruka/goflow/utils"
)

var version = "Dev"

func main() {
	m := NewConfigWithPath("flowserver.toml")
	config := new(FlowServerConfig)
	m.MustLoad(config)

	// if we have a custom version, use it
	if version != "Dev" {
		config.Version = version
	}

	logger := logrus.New()
	lg.RedirectStdlogOutput(logger)
	lg.DefaultLogger = logger

	flowServer := NewFlowServer(config, logger)
	flowServer.Start()

	logrus.WithField("comp", "server").WithField("port", "8080").WithField("version", version).Info("listening")

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	logrus.WithField("comp", "server").WithField("signal", <-ch).Info("stopping")

	flowServer.Stop()
}

type errorResponse struct {
	Text []string `json:"errors"`
}

// writeError writes a JSON response for the passed in error
func writeError(w http.ResponseWriter, r *http.Request, status int, err error) error {
	lg.Log(r.Context()).WithError(err).Error()
	var errors []string

	vErrs, isValidation := err.(utils.ValidationErrors)
	if isValidation {
		errors = []string{}
		for i := range vErrs {
			errors = append(errors, vErrs[i].Error())
		}
	} else {
		errors = []string{err.Error()}
	}

	return writeJSONResponse(w, r, status, &errorResponse{errors})
}

func writeJSONResponse(w http.ResponseWriter, r *http.Request, statusCode int, response interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Version", version)
	start := r.Context().Value(contextStart)
	if start != nil {
		elapsed := time.Since(start.(time.Time)).Nanoseconds()
		w.Header().Set("X-Elapsed-NS", strconv.FormatInt(elapsed, 10))
	}
	w.WriteHeader(statusCode)

	respJSON, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return err
	}
	_, err = w.Write(respJSON)
	return err
}

type jsonHandlerFunc func(http.ResponseWriter, *http.Request) (interface{}, error)

func jsonHandler(handler jsonHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// stuff our start time in our context
		r = r.WithContext(context.WithValue(r.Context(), contextStart, time.Now()))
		value, err := handler(w, r)
		if err != nil {
			writeError(w, r, http.StatusBadRequest, err)
		} else {
			err := writeJSONResponse(w, r, http.StatusOK, value)
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

type templateHandlerFunc func(http.FileSystem, http.ResponseWriter, *http.Request) error

func templateHandler(fs http.FileSystem, handler templateHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := handler(fs, w, r)
		if err != nil {
			writeError(w, r, http.StatusInternalServerError, err)
		}
	}
}

func traceErrors(logger *logrus.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			body := bytes.Buffer{}
			r.Body = ioutil.NopCloser(io.TeeReader(r.Body, &body))
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(ww, r)

			// we are returning an error of some kind, log the incoming request body
			if ww.Status() != 200 && strings.ToLower(r.Method) == "post" {
				logger.WithFields(logrus.Fields{
					"request_body": body.String(),
					"status":       ww.Status(),
					"req_id":       r.Context().Value(middleware.RequestIDKey)}).Error()
			}
		}
		return http.HandlerFunc(fn)
	}
}

type contextKey string

func (c contextKey) String() string {
	return "flowserver: " + string(c)
}

var (
	contextStart = contextKey("start")
)
