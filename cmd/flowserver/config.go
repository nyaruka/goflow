package main

import (
	"github.com/nyaruka/ezconf"
)

// Config is our top level config for our flowserver
type Config struct {
	Port      int    `help:"the port we will run on"`
	SentryDSN string `help:"the DSN for reporting errors to Sentry"`
	LogLevel  string `help:"the logging level to use"`
	Version   string `help:"the version to use in request and response headers"`
}

// NewDefaultConfig returns our default configuration
func NewDefaultConfig() *Config {
	return &Config{
		Port:     8800,
		LogLevel: "info",
		Version:  "Dev",
	}
}

// NewConfigWithPath returns a new instance of our config loaded from the path, environment and args
func NewConfigWithPath(path string) *Config {
	config := NewDefaultConfig()
	loader := ezconf.NewLoader(
		config,
		"flowserver", "flowserver - a self contained flow engine server",
		[]string{path},
	)
	loader.MustLoad()
	return config
}
