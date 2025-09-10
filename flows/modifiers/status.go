package modifiers

import (
	"github.com/nyaruka/gocommon/jsonx"
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

	status flows.ContactStatus
}

// NewStatus creates a new status modifier
func NewStatus(status flows.ContactStatus) *Status {
	return &Status{
		baseModifier: newBaseModifier(TypeStatus),
		status:       status,
	}
}

// Apply applies this modification to the given contact
func (m *Status) Apply(eng flows.Engine, env envs.Environment, sa flows.SessionAssets, contact *flows.Contact, log flows.EventCallback) bool {
	if contact.Status() != m.status {
		contact.SetStatus(m.status)
		log(events.NewContactStatusChanged(m.status))
		return true
	}
	return false
}

var _ flows.Modifier = (*Status)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type statusEnvelope struct {
	utils.TypedEnvelope

	Status flows.ContactStatus `json:"status" validate:"contact_status"`
}

func readStatus(sa flows.SessionAssets, data []byte, missing assets.MissingCallback) (flows.Modifier, error) {
	e := &statusEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	return NewStatus(e.Status), nil
}

func (m *Status) MarshalJSON() ([]byte, error) {
	return jsonx.Marshal(&statusEnvelope{
		TypedEnvelope: utils.TypedEnvelope{Type: m.Type()},
		Status:        m.status,
	})
}
