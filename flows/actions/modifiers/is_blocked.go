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

// IsBlockedModification is the type of modification to make
type IsBlockedModification string

// the supported types of modification
const (
	ShouldBlock   IsBlockedModification = "block"
	ShouldUnblock IsBlockedModification = "unblock"
)

// IsBlockedModifier modifies the is_blocked value of a contact
type IsBlockedModifier struct {
	baseModifier
	Modification IsBlockedModification `json:"modification" validate:"eq=block|eq=unblock"`
}

// NewIsBlocked creates a new is_blocked modifier
func NewIsBlocked(modification IsBlockedModification) *IsBlockedModifier {
	return &IsBlockedModifier{
		baseModifier: newBaseModifier(TypeIsBlocked),
		Modification: modification,
	}
}

// Apply applies this modification to the given contact
func (m *IsBlockedModifier) Apply(env envs.Environment, assets flows.SessionAssets, contact *flows.Contact, log flows.EventCallback) {
	isBlocked := m.Modification == ShouldBlock
	if isBlocked != contact.IsBlocked() {
		contact.SetIsBlocked(isBlocked)
		log(events.NewContactIsBlockedChanged(isBlocked))

		if isBlocked {
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

var _ flows.Modifier = (*IsBlockedModifier)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

func readIsBlockedModifier(assets flows.SessionAssets, data json.RawMessage, missing assets.MissingCallback) (flows.Modifier, error) {
	m := &IsBlockedModifier{}
	return m, utils.UnmarshalAndValidate(data, m)
}
