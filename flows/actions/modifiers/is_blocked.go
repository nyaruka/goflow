package modifiers

import (
	"encoding/json"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	registerType(TypeIsBlocked, readIsBlockedModifier)
}

// TypeIsBlocked is the type of our block modifier
const TypeIsBlocked string = "is_blocked"

// BlockModifier modifies the is_blocked value of a contact
type BlockModifier struct {
	baseModifier
	Modification bool `json:"modification" validate:"required"`
}

// NewIsBlocked creates a new is_blocked modifier
func NewIsBlocked(modification bool) *BlockModifier {
	return &BlockModifier{
		baseModifier: newBaseModifier(TypeIsBlocked),
		Modification: modification,
	}
}

// Apply applies this modification to the given contact
func (m *BlockModifier) Apply(env envs.Environment, assets flows.SessionAssets, contact *flows.Contact, log flows.EventCallback) {

	if m.Modification != contact.IsBlocked() {
		contact.SetIsBlocked(m.Modification)

		log(events.NewContactIsBlockedChanged(m.Modification))
		if m.Modification {
			diff := make([]*flows.Group, 0, len(contact.Groups().All()))

			for _, group := range contact.Groups().All() {
				contact.Groups().Remove(group)
				diff = append(diff, group)
			}
			// only generate event if contact's groups change
			if len(diff) > 0 {
				log(events.NewContactGroupsChanged(nil, diff))
			}

		} else {
			m.reevaluateDynamicGroups(env, assets, contact, log)
		}
	}
}

var _ flows.Modifier = (*BlockModifier)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

func readIsBlockedModifier(assets flows.SessionAssets, data json.RawMessage, missing assets.MissingCallback) (flows.Modifier, error) {
	m := &BlockModifier{}
	return m, utils.UnmarshalAndValidate(data, m)
}
