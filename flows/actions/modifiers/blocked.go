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
	registerType(TypeBlocked, readBlockedModifier)
}

// TypeBlocked is the type of our blocked modifier
const TypeBlocked string = "blocked"

// BlockedModifier modifies the blocked state of a contact
type BlockedModifier struct {
	baseModifier
	State bool `json:"state"`
}

// NewBlocked creates a new blocked modifier
func NewBlocked(state bool) *BlockedModifier {
	return &BlockedModifier{
		baseModifier: newBaseModifier(TypeBlocked),
		State:        state,
	}
}

// Apply applies this modification to the given contact
func (m *BlockedModifier) Apply(env envs.Environment, assets flows.SessionAssets, contact *flows.Contact, log flows.EventCallback) {
	if m.State != contact.Blocked() {
		contact.SetBlocked(m.State)

		if m.State {
			log(events.NewContactBlocked())
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
			log(events.NewContactUnblocked())
			m.reevaluateDynamicGroups(env, assets, contact, log)
		}
	}
}

var _ flows.Modifier = (*BlockedModifier)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

func readBlockedModifier(assets flows.SessionAssets, data json.RawMessage, missing assets.MissingCallback) (flows.Modifier, error) {
	m := &BlockedModifier{}
	return m, utils.UnmarshalAndValidate(data, m)
}
