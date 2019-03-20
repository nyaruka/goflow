package routers

import (
	"encoding/json"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"

	"github.com/pkg/errors"
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
func NewRandomRouter(resultName string, categories []*Category) *RandomRouter {
	return &RandomRouter{newBaseRouter(TypeRandom, resultName, categories)}
}

// Validate validates that the fields on this router are valid
func (r *RandomRouter) Validate(exits []flows.Exit) error {
	return r.validate(exits)
}

// PickExit determines which exit to take from a node
func (r *RandomRouter) PickExit(run flows.FlowRun, step flows.Step, logEvent flows.EventCallback) (flows.ExitUUID, error) {
	route, err := r.pickRoute(run, step)
	if err != nil {
		return "", err
	}

	// find the category
	var category *Category
	for _, c := range r.Categories_ {
		if c.UUID() == route.categoryUUID {
			category = c
			break
		}
	}

	if category == nil {
		return "", errors.Errorf("category %s is not a valid category", route.categoryUUID)
	}

	// save result if we have a result name
	if r.ResultName_ != "" {
		// localize the category name
		localizedCategory := run.GetText(utils.UUID(category.UUID()), "name", "")

		var extraJSON json.RawMessage
		if route.extra != nil {
			extraJSON, _ = json.Marshal(route.extra)
		}
		result := flows.NewResult(r.ResultName_, route.match, category.Name(), localizedCategory, step.NodeUUID(), route.input, extraJSON, utils.Now())
		run.SaveResult(result)
		logEvent(events.NewRunResultChangedEvent(result))
	}

	return category.ExitUUID(), nil
}

func (r *RandomRouter) pickRoute(run flows.FlowRun, step flows.Step) (*route, error) {
	// pick a random category
	rand := utils.RandDecimal()
	categoryNum := rand.Mul(decimal.New(int64(len(r.Categories_)), 0)).IntPart()
	categoryUUID := r.Categories_[categoryNum].UUID()
	return newRoute(nil, rand.String(), categoryUUID, nil), nil
}

// Inspect inspects this object and any children
func (r *RandomRouter) Inspect(inspect func(flows.Inspectable)) {
	inspect(r)
}
