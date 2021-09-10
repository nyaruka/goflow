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
	registerType(TypeURN, readURNModifier)
}

// TypeURN is the type of our URN modifier
const TypeURN string = "urn"

// URNModification is the type of modification to make
type URNModification string

// the supported types of modification
const (
	URNAppend URNModification = "append"
	URNRemove URNModification = "remove"
)

// URNModifier modifies a URN on a contact. This has been replaced by URNsModifier but is kept here for now
// to support processing of old Surveyor submissions.
type URNModifier struct {
	baseModifier

	URN          urns.URN        `json:"urn" validate:"required"`
	Modification URNModification `json:"modification" validate:"required,eq=append|eq=remove"`
}

// NewURN creates a new name modifier
func NewURN(urn urns.URN, modification URNModification) *URNModifier {
	return &URNModifier{
		baseModifier: newBaseModifier(TypeURN),
		URN:          urn,
		Modification: modification,
	}
}

// Apply applies this modification to the given contact
func (m *URNModifier) Apply(env envs.Environment, assets flows.SessionAssets, contact *flows.Contact, log flows.EventCallback) {
	urn := m.URN.Normalize(string(env.DefaultCountry()))
	modified := false

	if m.Modification == URNAppend {
		modified = contact.AddURN(urn, nil)
	} else {
		modified = contact.RemoveURN(urn)
	}

	if modified {
		log(events.NewContactURNsChanged(contact.URNs().RawURNs()))
		ReevaluateGroups(env, assets, contact, log)
	}
}

var _ flows.Modifier = (*URNModifier)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

func readURNModifier(assets flows.SessionAssets, data json.RawMessage, missing assets.MissingCallback) (flows.Modifier, error) {
	m := &URNModifier{}
	return m, utils.UnmarshalAndValidate(data, m)
}
