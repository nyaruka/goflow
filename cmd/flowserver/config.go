package main

import (
	"github.com/nyaruka/ezconf"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/engine"
)

// Config is our top level config for our flowserver
type Config struct {
	Port                          int    `help:"the port we will run on"`
	LogLevel                      string `help:"the logging level to use"`
	Static                        string `help:""`
	AssetCacheSize                int64  `help:"the maximum size of our asset cache"`
	AssetCachePrune               int    `help:"the number of assets to prune when we reach our max size"`
	AssetServerToken              string `help:"the token to use when authentication to the asset server"`
	EngineDisableWebhooks         bool   `help:"whether to disable webhook calls from the engine"`
	EngineMaxWebhookResponseBytes int    `help:"the maximum allowed byte size of webhook responses"`
	SentryDSN                     string `help:"the DSN for reporting errors to Sentry"`
	Version                       string `help:"the version to use in request and response headers"`
}

func (c *Config) Engine() flows.EngineConfig {
	return engine.NewConfig(c.EngineDisableWebhooks, nil, c.EngineMaxWebhookResponseBytes)
}

// NewDefaultConfig returns our default configuration
func NewDefaultConfig() *Config {
	return &Config{
		Port:                          8800,
		LogLevel:                      "info",
		AssetCacheSize:                1000,
		AssetCachePrune:               100,
		AssetServerToken:              "missing_temba_token",
		EngineDisableWebhooks:         false,
		EngineMaxWebhookResponseBytes: 10000,
		Version: "Dev",
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
