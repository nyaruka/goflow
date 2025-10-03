package routers

import (
	"fmt"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/gocommon/random"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"

	"github.com/shopspring/decimal"
)

func init() {
	registerType(TypeRandom, func() flows.Router { return &Random{} })
}

// TypeRandom is the type for a random router
const TypeRandom string = "random"

// Random is a router which will exit out a random exit
type Random struct {
	baseRouter
}

// NewRandom creates a new random router
func NewRandom(wait flows.Wait, resultName string, categories []flows.Category) *Random {
	return &Random{newBaseRouter(TypeRandom, wait, resultName, categories)}
}

// Validate validates that the fields on this router are valid
func (r *Random) Validate(flow flows.Flow, exits []flows.Exit) error {
	return r.validate(flow, exits)
}

// Route determines which exit to take from a node
func (r *Random) Route(run flows.Run, step flows.Step, logEvent flows.EventLogger) (flows.ExitUUID, string, error) {
	// pick a random category
	rand := random.Decimal()
	categoryNum := rand.Mul(decimal.New(int64(len(r.categories)), 0)).IntPart()
	categoryUUID := r.categories[categoryNum].UUID()

	exit, err := r.routeToCategory(run, step, categoryUUID, fmt.Sprintf("%d", categoryNum), rand.String(), nil, logEvent)
	return exit, rand.String(), err
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

func (r *Random) UnmarshalJSON(data []byte) error {
	e := &baseEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return err
	}

	if err := r.unmarshal(e); err != nil {
		return err
	}

	return nil
}

// MarshalJSON marshals this resume into JSON
func (r *Random) MarshalJSON() ([]byte, error) {
	e := &baseEnvelope{}

	if err := r.marshal(e); err != nil {
		return nil, err
	}

	return jsonx.Marshal(e)
}
