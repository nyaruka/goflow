package main

import (
	"github.com/koding/multiconfig"
)

// FlowServerConfig is our top level configuration object
type FlowServerConfig struct {
	Port int `default:"8080"`

	Static string `default:""`

	AssetCacheSize   int64  `default:"1000"`
	AssetCachePrune  int    `default:"100"`
	AssetServerToken string `default:"missing_temba_token"`

	Version string `default:"Dev"`
}

// NewConfigWithPath returns a new instance of Loader to read from the given configuration file using our config options
func NewConfigWithPath(path string) *multiconfig.DefaultLoader {
	loaders := []multiconfig.Loader{}

	loaders = append(loaders, &multiconfig.TagLoader{})
	loaders = append(loaders, &multiconfig.TOMLLoader{Path: path})
	loaders = append(loaders, &multiconfig.EnvironmentLoader{CamelCase: true})
	loaders = append(loaders, &multiconfig.FlagLoader{CamelCase: true})
	loader := multiconfig.MultiLoader(loaders...)

	return &multiconfig.DefaultLoader{Loader: loader, Validator: multiconfig.MultiValidator(&multiconfig.RequiredValidator{})}
}

// NewTestConfig returns a new instance of our config initialized just from our defaults as defined above
func NewTestConfig() *FlowServerConfig {
	loader := &multiconfig.DefaultLoader{
		Loader:    multiconfig.MultiLoader(&multiconfig.TagLoader{}),
		Validator: multiconfig.MultiValidator(&multiconfig.RequiredValidator{}),
	}

	config := &FlowServerConfig{}
	loader.Load(config)
	return config
}
