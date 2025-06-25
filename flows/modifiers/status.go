package modifiers

import (
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	registerType(TypeStatus, readStatus)
}

// TypeStatus is the type of our status modifier
const TypeStatus string = "status"

// Status modifies the status of a contact
type Status struct {
	baseModifier

	Status flows.ContactStatus `json:"status" validate:"contact_status"`
}

// NewStatus creates a new status modifier
func NewStatus(status flows.ContactStatus) *Status {
	return &Status{
		baseModifier: newBaseModifier(TypeStatus),
		Status:       status,
	}
}

// Apply applies this modification to the given contact
func (m *Status) Apply(eng flows.Engine, env envs.Environment, sa flows.SessionAssets, contact *flows.Contact, log flows.EventCallback) bool {
	if contact.Status() != m.Status {
		contact.SetStatus(m.Status)
		log(events.NewContactStatusChanged(m.Status))
		return true
	}
	return false
}

var _ flows.Modifier = (*Status)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

func readStatus(sa flows.SessionAssets, data []byte, missing assets.MissingCallback) (flows.Modifier, error) {
	m := &Status{}
	return m, utils.UnmarshalAndValidate(data, m)
}
