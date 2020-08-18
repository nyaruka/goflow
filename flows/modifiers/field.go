package modifiers

import (
	"encoding/json"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
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
	value string
}

// NewField creates a new field modifier
func NewField(field *flows.Field, value string) *FieldModifier {
	return &FieldModifier{
		baseModifier: newBaseModifier(TypeField),
		field:        field,
		value:        value,
	}
}

// Apply applies this modification to the given contact
func (m *FieldModifier) Apply(env envs.Environment, sa flows.SessionAssets, contact *flows.Contact, log flows.EventCallback) {
	oldValue := contact.Fields().Get(m.field)

	newValue := contact.Fields().Parse(env, sa.Fields(), m.field, m.value)

	// truncate text value if necessary
	if newValue != nil {
		newValue.Text = types.NewXText(utils.Truncate(newValue.Text.Native(), env.MaxValueLength()))
	}

	if !newValue.Equals(oldValue) {
		contact.Fields().Set(m.field, newValue)
		log(events.NewContactFieldChanged(m.field, newValue))
		m.reevaluateGroups(env, sa, contact, log)
	}
}

func (m *FieldModifier) Value() string {
	return m.value
}

var _ flows.Modifier = (*FieldModifier)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type fieldModifierEnvelope struct {
	utils.TypedEnvelope
	Field *assets.FieldReference `json:"field" validate:"required"`
	Value json.RawMessage        `json:"value"`
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

	value := ""

	// try unmarshaling value as string
	json.Unmarshal(e.Value, &value)

	// then try as a value object (how this modifier used to be work)
	valueAsObj := &flows.Value{}
	if json.Unmarshal(e.Value, valueAsObj) == nil {
		value = valueAsObj.Text.Native()
	}

	return NewField(field, value), nil
}

func (m *FieldModifier) MarshalJSON() ([]byte, error) {
	value, err := jsonx.Marshal(m.value)
	if err != nil {
		return nil, err
	}

	return jsonx.Marshal(&fieldModifierEnvelope{
		TypedEnvelope: utils.TypedEnvelope{Type: m.Type()},
		Field:         m.field.Reference(),
		Value:         value,
	})
}
