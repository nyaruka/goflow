//go:generate statik -src=./static
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/nyaruka/goflow/flows/engine"

	"errors"

	"github.com/koding/multiconfig"
	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
	"github.com/pressly/lg"
	"github.com/rakyll/statik/fs"
	"github.com/sirupsen/logrus"

	"strconv"

	"io"

	"io/ioutil"

	_ "github.com/nyaruka/goflow/cmd/flowserver/statik"
	"github.com/nyaruka/goflow/utils"
)

var version = "dev"
var assetCache *engine.AssetCache

func main() {
	m := multiconfig.New()
	config := new(Server)
	m.MustLoad(config)

	logger := logrus.New()

	lg.RedirectStdlogOutput(logger)
	lg.DefaultLogger = logger

	r := chi.NewRouter()

	r.Use(middleware.DefaultCompress)
	r.Use(middleware.StripSlashes)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(traceErrors(logger))
	r.Use(lg.RequestLogger(logger))
	r.Use(middleware.Timeout(60 * time.Second))

	// no static dir passed in? serve from statik
	var staticDir http.FileSystem
	var err error

	if config.Static == "" {
		staticDir, err = fs.New()
		if err != nil {
			lg.Fatal(err)
		}
		logrus.WithField("comp", "server").Info("using compiled statik assets")
	} else {
		staticDir = http.Dir(config.Static)
		logrus.WithField("comp", "server").Info("using asset dir: ", config.Static)
	}

	// root page just serves our example and "postman"" interface
	r.Get("/", templateHandler(staticDir, indexHandler))
	r.Get("/version", jsonHandler(handleVersion))

	r.Post("/flow/start", jsonHandler(handleStart))
	r.Post("/flow/resume", jsonHandler(handleResume))
	r.Post("/flow/migrate", jsonHandler(handleMigrate))
	r.Post("/expression", jsonHandler(handleExpression))

	r.NotFound(errorHandler(http.StatusNotFound, "not found"))
	r.MethodNotAllowed(errorHandler(http.StatusMethodNotAllowed, "method not allowed"))

	assetCache = engine.NewAssetCache()

	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", config.Port),
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
	logrus.WithField("comp", "server").WithField("port", "8080").WithField("version", version).Info("listening")

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

	vErrs, isValidation := err.(utils.ValidationErrors)
	if isValidation {
		status = http.StatusBadRequest
		errors = []string{}
		for i := range vErrs {
			errors = append(errors, vErrs[i].Error())
		}
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
			writeError(w, r, http.StatusInternalServerError, err)
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

func handleVersion(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	response := map[string]string{
		"version": version,
	}
	return response, nil
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
