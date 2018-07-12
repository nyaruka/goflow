package routers

import (
	"fmt"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

var registeredTypes = map[string](func() flows.Router){}

// RegisterType registers a new type of router
func RegisterType(name string, initFunc func() flows.Router) {
	registeredTypes[name] = initFunc
}

// BaseRouter is the base class for all our router classes
type BaseRouter struct {
	// ResultName_ is the name of the which the result of this router should be saved as (if any)
	ResultName_ string `json:"result_name,omitempty"`
}

// ResultName returns the name which the result of this router should be saved as (if any)
func (r *BaseRouter) ResultName() string { return r.ResultName_ }

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

// ReadRouter reads a router from the given typed envelope
func ReadRouter(envelope *utils.TypedEnvelope) (flows.Router, error) {
	f := registeredTypes[envelope.Type]
	if f == nil {
		return nil, fmt.Errorf("unknown type: %s", envelope.Type)
	}

	router := f()
	return router, utils.UnmarshalAndValidate(envelope.Data, router)
}
