package functions

import (
	"strings"

	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/utils"
)

// Call calls the given function with the given parameters
func Call(env utils.Environment, name string, params []types.XValue) types.XValue {
	var function XFunction
	var found bool

	function, found = XFUNCTIONS[name]
	if !found {
		return types.NewXErrorf("no function with name '%s'", name)
	}

	val := function(env, params...)

	// if function returned an error, wrap the error with the function name
	if types.IsXError(val) {
		return types.NewXErrorf("error calling %s: %s", strings.ToUpper(name), val.(types.XError).Error())
	}

	return val
}
