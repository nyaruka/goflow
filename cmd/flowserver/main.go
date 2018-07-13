//go:generate statik -src=./static
package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/nyaruka/goflow/cmd/flowserver/statik"
	_ "github.com/nyaruka/goflow/contrib/transferto"
	"github.com/nyaruka/goflow/utils"

	"github.com/evalphobia/logrus_sentry"
	log "github.com/sirupsen/logrus"
)

var splash = `                ______             
   ____  ____  / __/ /___ _      ______
  / __ '/ __ \/ /_/ / __ \ | /| / / __
 / /_/ / /_/ / __/ / /_/ / |/ |/ / _
 \__, /\____/_/ /_/\____/|__/|__/ _
/____/`

var version = "Dev"

func main() {
	config := NewConfigWithPath("flowserver.toml")

	// if we have a custom version, use it
	if version != "Dev" {
		config.Version = version
	}

	fmt.Printf("%s  --- version: %s ---\n", splash, config.Version)

	// configure logging
	level, err := log.ParseLevel(config.LogLevel)
	if err != nil {
		log.Fatalf("Invalid log level '%s'", level)
	}
	log.SetLevel(level)

	// configure error reporting to Sentry if we have a DSN
	if config.SentryDSN != "" {
		hook, err := logrus_sentry.NewSentryHook(config.SentryDSN, []log.Level{log.PanicLevel, log.FatalLevel, log.ErrorLevel})
		hook.Timeout = 0
		hook.StacktraceConfiguration.Enable = true
		hook.StacktraceConfiguration.Skip = 4
		hook.StacktraceConfiguration.Context = 5
		if err != nil {
			log.Fatalf("Invalid sentry DSN: '%s': %s", config.SentryDSN, err)
		}
		log.StandardLogger().Hooks.Add(hook)
	}

	// start the server
	flowServer := NewFlowServer(config)
	flowServer.Start()

	log.WithField("comp", "server").WithField("port", config.Port).WithField("version", version).Info("listening")

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.WithField("comp", "server").WithField("signal", <-ch).Info("stopping")

	flowServer.Stop()
}

type errorResponse struct {
	Text []string `json:"errors"`
}

// writeError writes a JSON response for the passed in error
func writeError(w http.ResponseWriter, r *http.Request, status int, err error) error {
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
	w.WriteHeader(statusCode)

	respJSON, err := utils.JSONMarshal(response)
	if err != nil {
		return err
	}
	_, err = w.Write(respJSON)
	return err
}

type jsonHandlerFunc func(http.ResponseWriter, *http.Request) (interface{}, error)

func jsonHandler(handler jsonHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		value, err := handler(w, r)
		if err != nil {
			writeError(w, r, http.StatusBadRequest, err)
		} else {
			err := writeJSONResponse(w, r, http.StatusOK, value)
			if err != nil {
				log.WithError(err).Error()
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
