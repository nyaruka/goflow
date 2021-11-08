package functions

import (
	"strings"

	"github.com/nyaruka/goflow/excellent/types"
)

// XFUNCTIONS is our map of functions available in Excellent which aren't tests
var XFUNCTIONS = map[string]*types.XFunction{}

// RegisterXFunction registers a new function in Excellent
func RegisterXFunction(name string, f types.XFunc) {
	XFUNCTIONS[name] = types.NewXFunction(name, f)
}

// Lookup returns the function with the given name (case-insensitive) or nil
func Lookup(name string) *types.XFunction {
	return XFUNCTIONS[strings.ToLower(name)]
}
