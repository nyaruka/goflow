package actions

import (
	"strings"
	"time"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/actions/modifiers"

	"github.com/pkg/errors"
)

func init() {
	RegisterType(TypeSetContactTimezone, func() flows.Action { return &SetContactTimezoneAction{} })
}

// TypeSetContactTimezone is the type for the set contact timezone action
const TypeSetContactTimezone string = "set_contact_timezone"

// SetContactTimezoneAction can be used to update the timezone of the contact. The timezone is a localizable
// template and white space is trimmed from the final value. An empty string clears the timezone.
// A [event:contact_timezone_changed] event will be created with the corresponding value.
//
//   {
//     "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//     "type": "set_contact_timezone",
//     "timezone": "Africa/Kigali"
//   }
//
// @action set_contact_timezone
type SetContactTimezoneAction struct {
	BaseAction
	universalAction

	Timezone string `json:"timezone"`
}

// NewSetContactTimezoneAction creates a new set timezone action
func NewSetContactTimezoneAction(uuid flows.ActionUUID, timezone string) *SetContactTimezoneAction {
	return &SetContactTimezoneAction{
		BaseAction: NewBaseAction(TypeSetContactTimezone, uuid),
		Timezone:   timezone,
	}
}

// Validate validates our action is valid and has all the assets it needs
func (a *SetContactTimezoneAction) Validate(assets flows.SessionAssets, context *flows.ValidationContext) error {
	return nil
}

// Execute runs this action
func (a *SetContactTimezoneAction) Execute(run flows.FlowRun, step flows.Step) error {
	if run.Contact() == nil {
		a.logError(run, step, errors.Errorf("can't execute action in session without a contact"))
		return nil
	}

	timezone, err := a.evaluateLocalizableTemplate(run, "timezone", a.Timezone)
	timezone = strings.TrimSpace(timezone)

	// if we received an error, log it
	if err != nil {
		a.logError(run, step, err)
		return nil
	}

	// timezone must be empty or valid timezone name
	var tz *time.Location
	if timezone != "" {
		tz, err = time.LoadLocation(timezone)
		if err != nil {
			a.logError(run, step, errors.Errorf("unrecognized timezone: '%s'", timezone))
			return nil
		}
	}

	mod := modifiers.NewTimezoneModifier(tz)
	if mod.Apply(run.Session().Assets(), run.Contact(), func(e flows.Event) { a.log(run, step, e) }) {
		a.reevaluateDynamicGroups(run, step)
	}

	return nil
}
