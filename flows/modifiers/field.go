package modifiers

import (
	"encoding/json"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/gocommon/stringsx"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	registerType(TypeField, readField)
}

// TypeField is the type of our field modifier
const TypeField string = "field"

// Field modifies a field value on the contact
type Field struct {
	baseModifier

	field *flows.Field
	value string
}

// NewField creates a new field modifier
func NewField(field *flows.Field, value string) *Field {
	return &Field{
		baseModifier: newBaseModifier(TypeField),
		field:        field,
		value:        value,
	}
}

// Apply applies this modification to the given contact
func (m *Field) Apply(eng flows.Engine, env envs.Environment, sa flows.SessionAssets, contact *flows.Contact, log flows.EventCallback) bool {
	oldValue := contact.Fields().Get(m.field)

	newValue := contact.Fields().Parse(env, sa.Fields(), m.field, m.value)

	// truncate text value if necessary
	if newValue != nil {
		newValue.Text = types.NewXText(stringsx.Truncate(newValue.Text.Native(), eng.Options().MaxFieldChars))
	}

	if !newValue.Equals(oldValue) {
		contact.Fields().Set(m.field, newValue)
		log(events.NewContactFieldChanged(m.field, newValue))
		return true
	}
	return false
}

func (m *Field) Value() string {
	return m.value
}

var _ flows.Modifier = (*Field)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type fieldEnvelope struct {
	utils.TypedEnvelope

	Field *assets.FieldReference `json:"field" validate:"required"`
	Value json.RawMessage        `json:"value"`
}

func readField(sa flows.SessionAssets, data []byte, missing assets.MissingCallback) (flows.Modifier, error) {
	e := &fieldEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	var field *flows.Field
	if e.Field != nil {
		field = sa.Fields().Get(e.Field.Key)
		if field == nil {
			missing(e.Field, nil)
			return nil, ErrNoModifier // nothing left to modify without the field
		}
	}

	value := ""

	// try unmarshaling value as string
	json.Unmarshal(e.Value, &value)

	return NewField(field, value), nil
}

func (m *Field) MarshalJSON() ([]byte, error) {
	value, err := jsonx.Marshal(m.value)
	if err != nil {
		return nil, err
	}

	return jsonx.Marshal(&fieldEnvelope{
		TypedEnvelope: utils.TypedEnvelope{Type: m.Type()},
		Field:         m.field.Reference(),
		Value:         value,
	})
}
