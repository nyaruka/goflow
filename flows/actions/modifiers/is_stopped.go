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
	registerType(TypeIsStopped, readIsStoppedModifier)
}

// TypeIsStopped is the type of our is_stoped modifier
const TypeIsStopped string = "is_stopped"

// IsStoppedModification is the type of modification to make
type IsStoppedModification string

// the supported types of modification
const (
	ShouldStop   IsStoppedModification = "stop"
	ShouldUnstop IsStoppedModification = "unstop"
)

// IsStoppedModifier modifies the is_stopped value of a contact
type IsStoppedModifier struct {
	baseModifier
	Modification IsStoppedModification `json:"modification" validate:"eq=stop|eq=unstop"`
}

// NewIsStopped creates a new is_stopped modifier
func NewIsStopped(modification IsStoppedModification) *IsStoppedModifier {
	return &IsStoppedModifier{
		baseModifier: newBaseModifier(TypeIsStopped),
		Modification: modification,
	}
}

// Apply applies this modification to the given contact
func (m *IsStoppedModifier) Apply(env envs.Environment, assets flows.SessionAssets, contact *flows.Contact, log flows.EventCallback) {
	isStopped := m.Modification == ShouldStop

	if isStopped != contact.IsStopped() {
		contact.SetIsStopped(isStopped)
		log(events.NewContactIsStoppedChanged(isStopped))

		if isStopped {
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

var _ flows.Modifier = (*IsStoppedModifier)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

func readIsStoppedModifier(assets flows.SessionAssets, data json.RawMessage, missing assets.MissingCallback) (flows.Modifier, error) {
	m := &IsStoppedModifier{}
	return m, utils.UnmarshalAndValidate(data, m)
}
