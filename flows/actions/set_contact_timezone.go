package actions

import (
	"fmt"
	"strings"
	"time"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
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
		a.logError(run, step, fmt.Errorf("can't execute action in session without a contact"))
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
			a.logError(run, step, fmt.Errorf("unrecognized timezone: '%s'", timezone))
			return nil
		}
	}

	if !timezonesEqual(run.Contact().Timezone(), tz) {
		run.Contact().SetTimezone(tz)
		a.log(run, step, events.NewContactTimezoneChangedEvent(timezone))
	}

	return nil
}

func timezonesEqual(tz1 *time.Location, tz2 *time.Location) bool {
	return (tz1 == nil && tz2 == nil) || (tz1 != nil && tz2 != nil && tz1.String() == tz2.String())
}
