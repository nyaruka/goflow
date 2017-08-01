package routers

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

// TypeRandomOnce is the constant for our random once router
const TypeRandomOnce string = "random_once"

// RandomOnceRouter exits of our exits once (randomly) before taking exit
type RandomOnceRouter struct {
	BaseRouter
	Exit flows.ExitUUID `json:"exit"     validate:"required"`
}

func NewRandomOnceRouter(exit flows.ExitUUID, resultName string) *RandomOnceRouter {
	return &RandomOnceRouter{
		BaseRouter: BaseRouter{ResultName_: resultName},
		Exit:       exit,
	}
}

// Type returns the type of our router
func (r *RandomOnceRouter) Type() string { return TypeRandomOnce }

// Validate validates the parameters on this router
func (r *RandomOnceRouter) Validate(exits []flows.Exit) error {
	err := utils.Validate(r)
	if err != nil {
		// check that our exit is valid
		found := false
		for _, e := range exits {
			if r.Exit == e.UUID() {
				found = true
				break
			}
		}

		if !found {
			err = fmt.Errorf("Invalid exit uuid: '%s'", r.Exit)
		}
	}
	return err
}

// PickRoute will attempt to take a random exit it hasn't taken before. If all exits have been taken, then it will
// take the exit specified in it's Exit parameter
func (r *RandomOnceRouter) PickRoute(run flows.FlowRun, exits []flows.Exit, step flows.Step) (flows.Route, error) {
	if len(exits) == 0 {
		return flows.NoRoute, nil
	}

	// find all the exits we have taken
	exited := make(map[flows.ExitUUID]bool)
	for _, s := range run.Path() {
		if s.NodeUUID() == step.NodeUUID() {
			exited[s.ExitUUID()] = true
		}
	}

	// build up a list of the valid exits
	var validExits []flows.ExitUUID
	for i := range exits {
		// this isn't our default exit and we haven't used it yet
		if exits[i].UUID() != r.Exit && !exited[exits[i].UUID()] {
			validExits = append(validExits, exits[i].UUID())
		}
	}

	// no valid choices? exit!
	if len(validExits) == 0 {
		return flows.NewRoute(r.Exit, "0"), nil
	}

	// ok, now pick one randomly
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	exitN := random.Intn(len(validExits))
	return flows.NewRoute(validExits[exitN], string(exitN)), nil
}
