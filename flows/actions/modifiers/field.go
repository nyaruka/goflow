package modifiers

import (
	"encoding/json"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	RegisterType(TypeField, readFieldModifier)
}

// TypeField is the type of our field modifier
const TypeField string = "field"

// FieldModifier modifies a field value on the contact
type FieldModifier struct {
	baseModifier

	field *flows.Field
	value *flows.Value
}

// NewFieldModifier creates a new field modifier
func NewFieldModifier(field *flows.Field, value *flows.Value) *FieldModifier {
	return &FieldModifier{
		baseModifier: newBaseModifier(TypeField),
		field:        field,
		value:        value,
	}
}

// Apply applies this modification to the given contact
func (m *FieldModifier) Apply(env utils.Environment, assets flows.SessionAssets, contact *flows.Contact, log flows.EventCallback) {
	oldValue := contact.Fields().Get(m.field)

	if !m.value.Equals(oldValue) {
		contact.Fields().Set(m.field, m.value)
		log(events.NewContactFieldChangedEvent(m.field, m.value))
		m.reevaluateDynamicGroups(env, assets, contact, log)
	}
}

var _ flows.Modifier = (*FieldModifier)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type fieldModifierEnvelope struct {
	utils.TypedEnvelope
	Field *assets.FieldReference `json:"field" validate:"required"`
	Value *flows.Value           `json:"value"`
}

func readFieldModifier(assets flows.SessionAssets, data json.RawMessage) (flows.Modifier, error) {
	e := &fieldModifierEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	var field *flows.Field
	if e.Field != nil {
		var err error
		if field, err = assets.Fields().Get(e.Field.Key); err != nil {
			return nil, err
		}
	}
	return NewFieldModifier(field, e.Value), nil
}

func (m *FieldModifier) MarshalJSON() ([]byte, error) {
	return json.Marshal(&fieldModifierEnvelope{
		TypedEnvelope: utils.TypedEnvelope{Type: m.Type()},
		Field:         m.field.Reference(),
		Value:         m.value,
	})
}
