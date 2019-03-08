package routers

import (
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"

	"github.com/shopspring/decimal"
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

// NewRandomRouter creates a new random router
func NewRandomRouter(resultName string) *RandomRouter {
	return &RandomRouter{newBaseRouter(TypeRandom, resultName)}
}

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
	rand := utils.RandDecimal()
	exitNum := rand.Mul(decimal.New(int64(len(exits)), 0)).IntPart()
	return nil, flows.NewRoute(exits[exitNum].UUID(), rand.String(), nil), nil
}

// Inspect inspects this object and any children
func (r *RandomRouter) Inspect(inspect func(flows.Inspectable)) {
	inspect(r)
}
