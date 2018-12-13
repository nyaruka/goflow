package modifiers

import (
	"encoding/json"

	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	RegisterType(TypeURN, readURNModifier)
}

// TypeURN is the type of our URN modifier
const TypeURN string = "urn"

// URNModification is the type of modification to make
type URNModification string

// the supported types of modification
const (
	URNAppend URNModification = "append"
)

// URNModifier modifies a URN on a contact
type URNModifier struct {
	baseModifier

	URN          urns.URN        `json:"urn"`
	Modification URNModification `json:"modification"`
}

// NewURNModifier creates a new name modifier
func NewURNModifier(urn urns.URN, modification URNModification) *URNModifier {
	return &URNModifier{
		baseModifier: newBaseModifier(TypeURN),
		URN:          urn,
	}
}

// Apply applies this modification to the given contact
func (m *URNModifier) Apply(env utils.Environment, assets flows.SessionAssets, contact *flows.Contact, log func(flows.Event)) {
	contactURN := flows.NewContactURN(m.URN.Normalize(""), nil)
	if contact.AddURN(contactURN) {
		log(events.NewContactURNsChangedEvent(contact.URNs().RawURNs()))
		m.reevaluateDynamicGroups(env, assets, contact, log)
	}
}

var _ Modifier = (*URNModifier)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

func readURNModifier(assets flows.SessionAssets, data json.RawMessage) (Modifier, error) {
	m := &URNModifier{}
	return m, utils.UnmarshalAndValidate(data, m)
}
