package routers

import (
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	RegisterType(TypeFirst, func() flows.Router { return &FirstRouter{} })
}

// TypeFirst is the type for FirstRouters
const TypeFirst string = "first"

// FirstRouter is a simple router that always takes the first exit
type FirstRouter struct {
	BaseRouter
}

func NewFirstRouter(resultName string) *FirstRouter {
	return &FirstRouter{BaseRouter: newBaseRouter(TypeFirst, resultName)}
}

// Validate validates the arguments on this router
func (r *FirstRouter) Validate(exits []flows.Exit) error {
	return utils.Validate(r)
}

// PickRoute always picks the first exit if available for this router
func (r *FirstRouter) PickRoute(run flows.FlowRun, exits []flows.Exit, step flows.Step) (*string, flows.Route, error) {
	if len(exits) == 0 {
		return nil, flows.NoRoute, nil
	}

	return nil, flows.NewRoute(exits[0].UUID(), "", nil), nil
}

// Inspect inspects this object and any children
func (r *FirstRouter) Inspect(inspect func(flows.Inspectable)) {
	inspect(r)
}
