package modifiers

import (
	"encoding/json"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
	"github.com/nyaruka/goflow/utils/jsonx"
)

func init() {
	registerType(TypeField, readFieldModifier)
}

// TypeField is the type of our field modifier
const TypeField string = "field"

// FieldModifier modifies a field value on the contact
type FieldModifier struct {
	baseModifier

	field *flows.Field
	value *flows.Value
}

// NewField creates a new field modifier
func NewField(field *flows.Field, value *flows.Value) *FieldModifier {
	return &FieldModifier{
		baseModifier: newBaseModifier(TypeField),
		field:        field,
		value:        value,
	}
}

// Apply applies this modification to the given contact
func (m *FieldModifier) Apply(env envs.Environment, assets flows.SessionAssets, contact *flows.Contact, log flows.EventCallback) {
	oldValue := contact.Fields().Get(m.field)

	if !m.value.Equals(oldValue) {
		var value *flows.Value

		// copy and truncate text value if necessary
		if m.value != nil {
			v := *m.value
			value = &v
			value.Text = types.NewXText(utils.Truncate(value.Text.Native(), env.MaxValueLength()))
		}

		contact.Fields().Set(m.field, value)
		log(events.NewContactFieldChanged(m.field, value))
		m.reevaluateGroups(env, assets, contact, log)
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

func readFieldModifier(assets flows.SessionAssets, data json.RawMessage, missing assets.MissingCallback) (flows.Modifier, error) {
	e := &fieldModifierEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	var field *flows.Field
	if e.Field != nil {
		field = assets.Fields().Get(e.Field.Key)
		if field == nil {
			missing(e.Field, nil)
			return nil, ErrNoModifier // nothing left to modify without the field
		}
	}
	return NewField(field, e.Value), nil
}

func (m *FieldModifier) MarshalJSON() ([]byte, error) {
	return jsonx.Marshal(&fieldModifierEnvelope{
		TypedEnvelope: utils.TypedEnvelope{Type: m.Type()},
		Field:         m.field.Reference(),
		Value:         m.value,
	})
}
