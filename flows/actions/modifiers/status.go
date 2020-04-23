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
	registerType(TypeStatus, readStatusModifier)
}

// TypeStatus is the type of our status modifier
const TypeStatus string = "status"

// StatusModifier modifies the status of a contact
type StatusModifier struct {
	baseModifier

	Status flows.ContactStatus `json:"status" validate:"eq=active|eq=blocked|eq=stopped"`
}

// NewStatus creates a new status modifier
func NewStatus(status flows.ContactStatus) *StatusModifier {
	return &StatusModifier{
		baseModifier: newBaseModifier(TypeStatus),
		Status:       status,
	}
}

// Apply applies this modification to the given contact
func (m *StatusModifier) Apply(env envs.Environment, assets flows.SessionAssets, contact *flows.Contact, log flows.EventCallback) {

	clearStatic := false
	changed := false

	if m.Status == flows.ContactStatusBlocked && !contact.Blocked() {
		contact.SetBlocked(true)
		clearStatic = true
		changed = true
	}

	if m.Status == flows.ContactStatusStopped && !contact.Stopped() {
		contact.SetStopped(true)
		clearStatic = true
		changed = true
	}

	if m.Status == flows.ContactStatusActive {
		if contact.Blocked() {
			contact.SetBlocked(false)
			changed = true
		}

		if contact.Stopped() {
			contact.SetStopped(false)
			changed = true
		}

		clearStatic = false
	}

	if changed {
		log(events.NewContactStatusChanged(m.Status))
		m.reevaluateGroups(env, assets, contact, clearStatic, log)
	}

}

var _ flows.Modifier = (*StatusModifier)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

func readStatusModifier(assets flows.SessionAssets, data json.RawMessage, missing assets.MissingCallback) (flows.Modifier, error) {
	m := &StatusModifier{}
	return m, utils.UnmarshalAndValidate(data, m)
}
