package routers

import (
	"fmt"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

var registeredTypes = map[string](func() flows.Router){}

func registerType(name string, initFunc func() flows.Router) {
	registeredTypes[name] = initFunc
}

// BaseRouter is the base class for all our router classes
type BaseRouter struct {
	// ResultName_ is the name of the which the result of this router should be saved as (if any)
	ResultName_ string `json:"result_name,omitempty"`
}

// ResultName returns the name which the result of this router should be saved as (if any)
func (r *BaseRouter) ResultName() string { return r.ResultName_ }

// RouterFromEnvelope attempts to build a router given the passed in TypedEnvelope
func RouterFromEnvelope(envelope *utils.TypedEnvelope) (flows.Router, error) {
	initFunc := registeredTypes[envelope.Type]
	if initFunc == nil {
		return nil, fmt.Errorf("unknown router type: %s", envelope.Type)
	}

	router := initFunc()
	if err := utils.UnmarshalAndValidate(envelope.Data, router, ""); err != nil {
		return nil, fmt.Errorf("unable to read router[type=%s]: %s", envelope.Type, err)
	}
	return router, nil
}
