package modifiers

import (
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
)

func init() {
	RegisterType(TypeField, func() Modifier { return &FieldModifier{} })
}

// TypeField is the type of our field modifier
const TypeField string = "field"

// FieldModifier modifies a field value on the contact
type FieldModifier struct {
	baseModifier

	Field *flows.Field
	Value *flows.Value
}

// NewFieldModifier creates a new field modifier
func NewFieldModifier(field *flows.Field, value *flows.Value) *FieldModifier {
	return &FieldModifier{
		baseModifier: newBaseModifier(TypeField),
		Field:        field,
		Value:        value,
	}
}

// Apply applies this modification to the given contact
func (m *FieldModifier) Apply(assets flows.SessionAssets, contact *flows.Contact, log func(flows.Event)) bool {
	oldValue := contact.Fields().Get(m.Field)

	if !m.Value.Equals(oldValue) {
		contact.Fields().Set(m.Field, m.Value)
		log(events.NewContactFieldChangedEvent(m.Field, m.Value))
		return true
	}
	return false
}

var _ Modifier = (*FieldModifier)(nil)
