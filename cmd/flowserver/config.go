package main

import "github.com/nyaruka/ezconf"

// Config is our top level config for our flowserver
type Config struct {
	Port             int    `help:"the port we will run on"`
	Static           string `help:""`
	AssetCacheSize   int64  `help:"the size of our asset cache"`
	AssetCachePrune  int    `help:"when to prune our asset cache"`
	AssetServerToken string `help:"the token to use when authentication to the asset server"`
	Version          string `help:"the version to use in request and response headers"`
}

// NewDefaultConfig returns our default configuration
func NewDefaultConfig() *Config {
	return &Config{
		Port:             8800,
		AssetCacheSize:   1000,
		AssetCachePrune:  100,
		AssetServerToken: "missing_temba_token",
		Version:          "Dev",
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
