package actions

import (
	"context"
	"strings"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/flows/modifiers"
)

func init() {
	registerType(TypeSetContactField, func() flows.Action { return &SetContactFieldAction{} })
}

// TypeSetContactField is the type for the set contact field action
const TypeSetContactField string = "set_contact_field"

// SetContactFieldAction can be used to update a field value on the contact. The value is a localizable
// template and white space is trimmed from the final value. An empty string clears the value.
// A [event:contact_field_changed] event will be created with the corresponding value.
//
//	{
//	  "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//	  "type": "set_contact_field",
//	  "field": {"key": "gender", "name": "Gender"},
//	  "value": "Female"
//	}
//
// @action set_contact_field
type SetContactFieldAction struct {
	baseAction
	universalAction

	Field *assets.FieldReference `json:"field" validate:"required"`
	Value string                 `json:"value" engine:"evaluated"`
}

// NewSetContactField creates a new set channel action
func NewSetContactField(uuid flows.ActionUUID, field *assets.FieldReference, value string) *SetContactFieldAction {
	return &SetContactFieldAction{
		baseAction: newBaseAction(TypeSetContactField, uuid),
		Field:      field,
		Value:      value,
	}
}

// Execute runs this action
func (a *SetContactFieldAction) Execute(ctx context.Context, run flows.Run, step flows.Step, logModifier flows.ModifierCallback, logEvent flows.EventCallback) error {
	value, ok := run.EvaluateTemplate(a.Value, logEvent)
	value = strings.TrimSpace(value)

	if !ok {
		return nil
	}

	fields := run.Session().Assets().Fields()
	field := fields.Get(a.Field.Key)

	if field != nil {
		a.applyModifier(run, modifiers.NewField(field, value), logModifier, logEvent)
	} else {
		logEvent(events.NewDependencyError(a.Field))
	}
	return nil
}

func (a *SetContactFieldAction) Inspect(dependency func(assets.Reference), local func(string), result func(*flows.ResultInfo)) {
	dependency(a.Field)
}
