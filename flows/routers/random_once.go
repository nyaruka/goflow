package routers

import (
	"encoding/json"
	"strconv"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
	"github.com/nyaruka/goflow/utils/jsonx"
	"github.com/nyaruka/goflow/utils/random"

	"github.com/pkg/errors"
)

func init() {
	registerType(TypeRandomOnce, readRandomOnceRouter)
}

// TypeRandomOnce is the type for a random once router
const TypeRandomOnce string = "random_once"

// RandomOnceRouter is a router which will choose one of it's exits at random, but only using each exit once
type RandomOnceRouter struct {
	baseRouter

	defaultCategoryUUID flows.CategoryUUID
}

// NewRandomOnce creates a new random once router
func NewRandomOnce(wait flows.Wait, resultName string, categories []*Category, doneCategoryUUID flows.CategoryUUID) *RandomOnceRouter {
	return &RandomOnceRouter{
		baseRouter:          newBaseRouter(TypeRandom, wait, resultName, categories),
		defaultCategoryUUID: doneCategoryUUID,
	}
}

// Validate validates that the fields on this router are valid
func (r *RandomOnceRouter) Validate(exits []flows.Exit) error {
	// check the done category is valid
	if r.defaultCategoryUUID != "" && !r.isValidCategory(r.defaultCategoryUUID) {
		return errors.Errorf("default category %s is not a valid category", r.defaultCategoryUUID)
	}

	return r.validate(exits)
}

// Route determines which exit to take from a node
func (r *RandomOnceRouter) Route(run flows.FlowRun, step flows.Step, logEvent flows.EventCallback) (flows.ExitUUID, error) {
	// find all the exits we have taken in the past
	exitsTaken := make(map[flows.ExitUUID]bool)
	for _, s := range run.Path() {
		if s.NodeUUID() == step.NodeUUID() {
			exitsTaken[s.ExitUUID()] = true
		}
	}

	// convert that to bucket categories remaining
	bucketsRemaining := make([]*Category, 0)
	for _, c := range r.categories {
		if !exitsTaken[c.exitUUID] && c.UUID() != r.defaultCategoryUUID {
			bucketsRemaining = append(bucketsRemaining, c)
		}
	}

	// take the default route if we've been to all the other bucket categories
	if len(bucketsRemaining) == 0 {
		return r.routeToCategory(run, step, r.defaultCategoryUUID, "", "", nil, logEvent)
	}

	// pick a random remaining bucket category
	categoryNum := random.IntN(len(bucketsRemaining))
	categoryUUID := bucketsRemaining[categoryNum].UUID()

	return r.routeToCategory(run, step, categoryUUID, strconv.Itoa(categoryNum), "", nil, logEvent)
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type randomOnceRouterEnvelope struct {
	baseRouterEnvelope

	DefaultCategoryUUID flows.CategoryUUID `json:"default_category_uuid" validate:"required,uuid4"`
}

func readRandomOnceRouter(data json.RawMessage) (flows.Router, error) {
	e := &randomOnceRouterEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	r := &RandomOnceRouter{
		defaultCategoryUUID: e.DefaultCategoryUUID,
	}

	if err := r.unmarshal(&e.baseRouterEnvelope); err != nil {
		return nil, err
	}

	return r, nil
}

// MarshalJSON marshals this resume into JSON
func (r *RandomOnceRouter) MarshalJSON() ([]byte, error) {
	e := &randomOnceRouterEnvelope{
		DefaultCategoryUUID: r.defaultCategoryUUID,
	}

	if err := r.marshal(&e.baseRouterEnvelope); err != nil {
		return nil, err
	}

	return jsonx.Marshal(e)
}
