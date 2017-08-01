package routers

import (
	"math/rand"
	"time"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

// TypeRandom is the type for a random router
const TypeRandom string = "random"

// RandomRouter is a router which will exit out a random exit
type RandomRouter struct {
	BaseRouter
}

// Type returns the type of this router
func (r *RandomRouter) Type() string { return TypeRandom }

// Validate validates that the fields on this router are valid
func (r *RandomRouter) Validate(exits []flows.Exit) error {
	return utils.Validate(r)
}

// PickRoute picks a route randomly from our available exits
func (r *RandomRouter) PickRoute(run flows.FlowRun, exits []flows.Exit, step flows.Step) (flows.Route, error) {
	if len(exits) == 0 {
		return flows.NoRoute, nil
	}

	// pick a random exit
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	exitN := random.Intn(len(exits))
	return flows.NewRoute(exits[exitN].UUID(), string(exitN)), nil
}
