package main

import (
	"github.com/nyaruka/ezconf"
	"github.com/nyaruka/goflow/flows"
)

// Config is our top level config for our flowserver
type Config struct {
	Port                          int    `help:"the port we will run on"`
	LogLevel                      string `help:"the logging level to use"`
	Static                        string `help:""`
	AssetCacheSize                int64  `help:"the maximum size of our asset cache"`
	AssetCachePrune               int    `help:"the number of assets to prune when we reach our max size"`
	AssetServerToken              string `help:"the token to use when authentication to the asset server"`
	EngineMaxWebhookResponseBytes int    `help:"the maximum allowed byte size of webhook responses"`
	Version                       string `help:"the version to use in request and response headers"`
}

func (c *Config) MaxWebhookResponseBytes() int { return c.EngineMaxWebhookResponseBytes }

// NewDefaultConfig returns our default configuration
func NewDefaultConfig() *Config {
	return &Config{
		Port:                          8800,
		LogLevel:                      "info",
		AssetCacheSize:                1000,
		AssetCachePrune:               100,
		AssetServerToken:              "missing_temba_token",
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

var _ flows.EngineConfig = (*Config)(nil)
