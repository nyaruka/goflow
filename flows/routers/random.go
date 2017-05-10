package routers

import (
	"math/rand"
	"time"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

const RANDOM string = "random"

type RandomRouter struct {
	BaseRouter
}

func (r *RandomRouter) Type() string { return RANDOM }

func (r *RandomRouter) Validate(exits []flows.Exit) error {
	return utils.ValidateAll(r)
}

func (r *RandomRouter) PickRoute(run flows.FlowRun, exits []flows.Exit, step flows.Step) (flows.Route, error) {
	if len(exits) == 0 {
		return flows.NoRoute, nil
	}

	// pick a random exit
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	exitN := random.Intn(len(exits))
	return flows.NewRoute(exits[exitN].UUID(), string(exitN)), nil
}
