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
	registerType(TypeStopped, readStoppedModifier)
}

// TypeStopped is the type of our stopped modifier
const TypeStopped string = "stopped"

// StoppedModifier modifies the stopped state of a contact
type StoppedModifier struct {
	baseModifier

	State bool `json:"state"`
}

// NewStopped creates a new stopped modifier
func NewStopped(state bool) *StoppedModifier {
	return &StoppedModifier{
		baseModifier: newBaseModifier(TypeStopped),
		State:        state,
	}
}

// Apply applies this modification to the given contact
func (m *StoppedModifier) Apply(env envs.Environment, assets flows.SessionAssets, contact *flows.Contact, log flows.EventCallback) {
	if m.State != contact.Stopped() {
		contact.SetStopped(m.State)
		if m.State {
			log(events.NewContactStopped())
		} else {
			log(events.NewContactUnstopped())
		}
		m.reevaluateGroups(env, assets, contact, m.State, log)
	}
}

var _ flows.Modifier = (*StoppedModifier)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

func readStoppedModifier(assets flows.SessionAssets, data json.RawMessage, missing assets.MissingCallback) (flows.Modifier, error) {
	m := &StoppedModifier{}
	return m, utils.UnmarshalAndValidate(data, m)
}
