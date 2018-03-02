package main

import (
	"github.com/koding/multiconfig"
)

// Config is our top level configuration object
type Config struct {
	Port int `default:"8800" toml:"port"`

	Static string `default:""`

	AssetCacheSize   int64  `default:"1000"                toml:"asset_cache_size"`
	AssetCachePrune  int    `default:"100"                 toml:"asset_cache_prune"`
	AssetServerToken string `default:"missing_temba_token" toml:"asset_server_token"`

	Version string `default:"Dev"`
}

// NewConfigWithPath returns a new instance of Loader to read from the given configuration file using our config options
func NewConfigWithPath(path string) *multiconfig.DefaultLoader {
	loaders := []multiconfig.Loader{}

	loaders = append(loaders, &multiconfig.TagLoader{})
	loaders = append(loaders, &multiconfig.TOMLLoader{Path: path})
	loaders = append(loaders, &multiconfig.EnvironmentLoader{Prefix: "FLOWSERVER", CamelCase: true})
	loaders = append(loaders, &multiconfig.FlagLoader{EnvPrefix: "FLOWSERVER", CamelCase: true})
	loader := multiconfig.MultiLoader(loaders...)

	return &multiconfig.DefaultLoader{Loader: loader, Validator: multiconfig.MultiValidator(&multiconfig.RequiredValidator{})}
}

// NewTestConfig returns a new instance of our config initialized just from our defaults as defined above
func NewTestConfig() *Config {
	loader := &multiconfig.DefaultLoader{
		Loader:    multiconfig.MultiLoader(&multiconfig.TagLoader{}),
		Validator: multiconfig.MultiValidator(&multiconfig.RequiredValidator{}),
	}

	config := &Config{}
	loader.Load(config)
	return config
}
