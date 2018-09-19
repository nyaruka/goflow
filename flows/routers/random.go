package routers

import (
	"fmt"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	RegisterType(TypeRandom, func() flows.Router { return &RandomRouter{} })
}

// TypeRandom is the type for a random router
const TypeRandom string = "random"

// RandomRouter is a router which will exit out a random exit
type RandomRouter struct {
	BaseRouter
}

func NewRandomRouter(resultName string) *RandomRouter {
	return &RandomRouter{BaseRouter: BaseRouter{ResultName_: resultName}}
}

// Type returns the type of this router
func (r *RandomRouter) Type() string { return TypeRandom }

// Validate validates that the fields on this router are valid
func (r *RandomRouter) Validate(exits []flows.Exit) error {
	return utils.Validate(r)
}

// PickRoute picks a route randomly from our available exits
func (r *RandomRouter) PickRoute(run flows.FlowRun, exits []flows.Exit, step flows.Step) (*string, flows.Route, error) {
	if len(exits) == 0 {
		return nil, flows.NoRoute, nil
	}

	// pick a random exit
	exitN := utils.RandIntN(len(exits))
	return nil, flows.NewRoute(exits[exitN].UUID(), fmt.Sprintf("%d", exitN), nil), nil
}
