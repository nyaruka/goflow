package routers

import (
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

func init() {
	RegisterType(TypeRandomOnce, func() flows.Router { return &RandomOnceRouter{} })
}

// TypeRandomOnce is the constant for our random once router
const TypeRandomOnce string = "random_once"

// RandomOnceRouter exits of our exits once (randomly) before taking exit
type RandomOnceRouter struct {
	BaseRouter
	Default flows.ExitUUID `json:"default_exit_uuid" validate:"required,uuid4"`
}

// NewRandomOnceRouter creates a new random-once router
func NewRandomOnceRouter(defaultExit flows.ExitUUID, resultName string) *RandomOnceRouter {
	return &RandomOnceRouter{
		BaseRouter: newBaseRouter(TypeRandomOnce, resultName),
		Default:    defaultExit,
	}
}

// Validate validates the parameters on this router
func (r *RandomOnceRouter) Validate(exits []flows.Exit) error {
	// check that our default exit is valid
	found := false
	for _, e := range exits {
		if r.Default == e.UUID() {
			found = true
			break
		}
	}

	if !found {
		return errors.Errorf("default exit %s is not a valid exit", r.Default)
	}
	return nil
}

// PickRoute will attempt to take a random exit it hasn't taken before. If all exits have been taken, then it will
// take the exit specified in it's Exit parameter
func (r *RandomOnceRouter) PickRoute(run flows.FlowRun, exits []flows.Exit, step flows.Step) (*string, flows.Route, error) {
	if len(exits) == 0 {
		return nil, flows.NoRoute, nil
	}

	// find all the exits we have taken
	takenBefore := make(map[flows.ExitUUID]bool)
	for _, s := range run.Path() {
		if s.NodeUUID() == step.NodeUUID() {
			takenBefore[s.ExitUUID()] = true
		}
	}

	// build up a list of the valid exits
	var validExits []flows.ExitUUID
	for i := range exits {
		// this isn't our default exit and we haven't taken it yet
		if exits[i].UUID() != r.Default && !takenBefore[exits[i].UUID()] {
			validExits = append(validExits, exits[i].UUID())
		}
	}

	// no valid choices? exit!
	if len(validExits) == 0 {
		return nil, flows.NewRoute(r.Default, "", nil), nil
	}

	// ok, now pick one randomly
	rand := utils.RandDecimal()
	exitNum := rand.Mul(decimal.New(int64(len(validExits)), 0)).IntPart()
	return nil, flows.NewRoute(validExits[exitNum], rand.String(), nil), nil
}

// Inspect inspects this object and any children
func (r *RandomOnceRouter) Inspect(inspect func(flows.Inspectable)) {
	inspect(r)
}
