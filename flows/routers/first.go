package routers

import (
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

const TypeFirst string = "first"

type FirstRouter struct {
	BaseRouter
}

func (r *FirstRouter) Type() string { return TypeFirst }

func (r *FirstRouter) Validate(exits []flows.Exit) error {
	return utils.ValidateAll(r)
}

func (r *FirstRouter) PickRoute(run flows.FlowRun, exits []flows.Exit, step flows.Step) (flows.Route, error) {
	if len(exits) == 0 {
		return flows.NoRoute, nil
	}

	return flows.NewRoute(exits[0].UUID(), ""), nil
}
