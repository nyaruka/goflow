package main

import (
	"github.com/koding/multiconfig"
)

// FlowServer is our top level configuration object
type FlowServer struct {
	Port int `default:"8080"`

	Static string `default:""`

	AssetCacheSize  int64 `default:"1000"`
	AssetCachePrune int   `default:"100"`
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
