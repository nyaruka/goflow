package actions

import (
	"fmt"
	"strings"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

// TypeSetContactTimezone is the type for the set contact timezone action
const TypeSetContactTimezone string = "set_contact_timezone"

// SetContactTimezoneAction can be used to update the timezone of the contact. A `contact_timezone_changed`
// event will be created with the corresponding value.
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
	Timezone string `json:"timezone"`
}

// Type returns the type of this action
func (a *SetContactTimezoneAction) Type() string { return TypeSetContactTimezone }

// Validate validates our action is valid and has all the assets it needs
func (a *SetContactTimezoneAction) Validate(assets flows.SessionAssets) error {
	return nil
}

// Execute runs this action
func (a *SetContactTimezoneAction) Execute(run flows.FlowRun, step flows.Step, log flows.EventLog) error {
	if run.Contact() == nil {
		log.Add(events.NewFatalErrorEvent(fmt.Errorf("can't execute action in session without a contact")))
		return nil
	}

	// get our localized value if any
	template := run.GetText(utils.UUID(a.UUID()), "timezone", a.Timezone)
	timezone, err := run.EvaluateTemplateAsString(template, false)
	timezone = strings.TrimSpace(timezone)

	// if we received an error, log it
	if err != nil {
		log.Add(events.NewErrorEvent(err))
		return nil
	}

	log.Add(events.NewContactTimezoneChangedEvent(timezone))
	return nil
}
