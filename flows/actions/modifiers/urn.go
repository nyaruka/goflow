package modifiers

import (
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
)

func init() {
	RegisterType(TypeURN, func() Modifier { return &URNModifier{} })
}

// TypeURN is the type of our URN modifier
const TypeURN string = "urn"

// URNModification is the type of modification to make
type URNModification string

// the supported types of modification
const (
	URNAppend URNModification = "append"
)

// URNModifier modifies a URN on a contact
type URNModifier struct {
	baseModifier

	URN          urns.URN        `json:"urn"`
	Modification URNModification `json:"modification"`
}

// NewURNModifier creates a new name modifier
func NewURNModifier(urn urns.URN, modification URNModification) *URNModifier {
	return &URNModifier{
		baseModifier: newBaseModifier(TypeURN),
		URN:          urn,
	}
}

// Apply applies this modification to the given contact
func (m *URNModifier) Apply(assets flows.SessionAssets, contact *flows.Contact) flows.Event {
	contactURN := flows.NewContactURN(m.URN.Normalize(""), nil)

	if contact.AddURN(contactURN) {
		return events.NewContactURNsChangedEvent(contact.URNs().RawURNs())
	}

	return nil
}

var _ Modifier = (*URNModifier)(nil)
