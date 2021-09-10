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

	Status flows.ContactStatus `json:"status" validate:"contact_status"`
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

	if contact.Status() != m.Status {
		contact.SetStatus(m.Status)
		log(events.NewContactStatusChanged(m.Status))
		ReevaluateGroups(env, assets, contact, log)
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
