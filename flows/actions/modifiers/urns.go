package modifiers

import (
	"encoding/json"

	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	registerType(TypeURNs, readURNsModifier)
}

// TypeURNs is the type of our URNs modifier
const TypeURNs string = "urns"

// URNsModification is the type of modification to make
type URNsModification string

// the supported types of modification
const (
	URNsAppend URNsModification = "append"
	URNsRemove URNsModification = "remove"
	URNsSet    URNsModification = "set"
)

// URNsModifier modifies the URNs on a contact
type URNsModifier struct {
	baseModifier

	URNs         []urns.URN       `json:"urns" validate:"required"`
	Modification URNsModification `json:"modification" validate:"required,eq=append|eq=remove|eq=set"`
}

// NewURNs creates a new URNs modifier
func NewURNs(urnz []urns.URN, modification URNsModification) *URNsModifier {
	return &URNsModifier{
		baseModifier: newBaseModifier(TypeURNs),
		URNs:         urnz,
		Modification: modification,
	}
}

// Apply applies this modification to the given contact
func (m *URNsModifier) Apply(env envs.Environment, assets flows.SessionAssets, contact *flows.Contact, log flows.EventCallback) {
	modified := false

	if m.Modification == URNsSet {
		modified = contact.ClearURNs()
	}

	for _, urn := range m.URNs {
		// normalize the URN
		urn := urn.Normalize(string(env.DefaultCountry()))

		if m.Modification == URNsAppend || m.Modification == URNsSet {
			modified = contact.AddURN(urn, nil)
		} else {
			modified = contact.RemoveURN(urn)
		}
	}

	if modified {
		log(events.NewContactURNsChanged(contact.URNs().RawURNs()))
		m.reevaluateGroups(env, assets, contact, log)
	}
}

var _ flows.Modifier = (*URNsModifier)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

func readURNsModifier(assets flows.SessionAssets, data json.RawMessage, missing assets.MissingCallback) (flows.Modifier, error) {
	m := &URNsModifier{}
	return m, utils.UnmarshalAndValidate(data, m)
}
