package routers

import (
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

// TypeFirst is the type for FirstRouters
const TypeFirst string = "first"

// FirstRouter is a simple router that always takes the first exit
type FirstRouter struct {
	BaseRouter
}

func NewFirstRouter(resultName string) *FirstRouter {
	return &FirstRouter{BaseRouter: BaseRouter{ResultName_: resultName}}
}

// Type returns the type of this router
func (r *FirstRouter) Type() string { return TypeFirst }

// Validate validates the arguments on this router
func (r *FirstRouter) Validate(exits []flows.Exit) error {
	return utils.Validate(r)
}

// PickRoute always picks the first exit if available for this router
func (r *FirstRouter) PickRoute(run flows.FlowRun, exits []flows.Exit, step flows.Step) (*string, flows.Route, error) {
	if len(exits) == 0 {
		return nil, flows.NoRoute, nil
	}

	return nil, flows.NewRoute(exits[0].UUID(), ""), nil
}
