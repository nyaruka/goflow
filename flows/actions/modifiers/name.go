package modifiers

import (
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	RegisterType(TypeName, func() Modifier { return &NameModifier{} })
}

// TypeName is the type of our name modifier
const TypeName string = "name"

// NameModifier modifies the name of a contact
type NameModifier struct {
	baseModifier

	Name string `json:"name"`
}

// NewNameModifier creates a new name modifier
func NewNameModifier(name string) *NameModifier {
	return &NameModifier{
		baseModifier: newBaseModifier(TypeName),
		Name:         name,
	}
}

// Apply applies this modification to the given contact
func (m *NameModifier) Apply(env utils.Environment, assets flows.SessionAssets, contact *flows.Contact, log func(flows.Event)) {
	if contact.Name() != m.Name {
		contact.SetName(m.Name)
		log(events.NewContactNameChangedEvent(m.Name))
		m.reevaluateDynamicGroups(env, assets, contact, log)
	}
}

var _ Modifier = (*NameModifier)(nil)
