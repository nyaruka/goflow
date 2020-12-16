package routers

import (
	"encoding/json"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/gocommon/random"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"

	"github.com/shopspring/decimal"
)

func init() {
	registerType(TypeRandom, readRandomRouter)
}

// TypeRandom is the type for a random router
const TypeRandom string = "random"

// RandomRouter is a router which will exit out a random exit
type RandomRouter struct {
	baseRouter
}

// NewRandom creates a new random router
func NewRandom(wait flows.Wait, resultName string, categories []flows.Category) *RandomRouter {
	return &RandomRouter{newBaseRouter(TypeRandom, wait, resultName, categories)}
}

// Validate validates that the fields on this router are valid
func (r *RandomRouter) Validate(flow flows.Flow, exits []flows.Exit) error {
	return r.validate(flow, exits)
}

// Route determines which exit to take from a node
func (r *RandomRouter) Route(run flows.FlowRun, step flows.Step, logEvent flows.EventCallback) (flows.ExitUUID, error) {
	// pick a random category
	rand := random.Decimal()
	categoryNum := rand.Mul(decimal.New(int64(len(r.categories)), 0)).IntPart()
	categoryUUID := r.categories[categoryNum].UUID()

	// TODO should raw rand value be iput and category number the match ?
	return r.routeToCategory(run, step, categoryUUID, rand.String(), "", nil, logEvent)
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

func readRandomRouter(data json.RawMessage) (flows.Router, error) {
	e := &baseRouterEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	r := &RandomRouter{}

	if err := r.unmarshal(e); err != nil {
		return nil, err
	}

	return r, nil
}

// MarshalJSON marshals this resume into JSON
func (r *RandomRouter) MarshalJSON() ([]byte, error) {
	e := &baseRouterEnvelope{}

	if err := r.marshal(e); err != nil {
		return nil, err
	}

	return jsonx.Marshal(e)
}
