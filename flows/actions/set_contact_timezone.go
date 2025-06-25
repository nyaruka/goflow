package actions

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/flows/modifiers"
)

func init() {
	registerType(TypeSetContactTimezone, func() flows.Action { return &SetContactTimezone{} })
}

// TypeSetContactTimezone is the type for the set contact timezone action
const TypeSetContactTimezone string = "set_contact_timezone"

// SetContactTimezone can be used to update the timezone of the contact. The timezone is a localizable
// template and white space is trimmed from the final value. An empty string clears the timezone.
// A [event:contact_timezone_changed] event will be created with the corresponding value.
//
//	{
//	  "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//	  "type": "set_contact_timezone",
//	  "timezone": "Africa/Kigali"
//	}
//
// @action set_contact_timezone
type SetContactTimezone struct {
	baseAction
	universalAction

	Timezone string `json:"timezone" engine:"evaluated"`
}

// NewSetContactTimezone creates a new set timezone action
func NewSetContactTimezone(uuid flows.ActionUUID, timezone string) *SetContactTimezone {
	return &SetContactTimezone{
		baseAction: newBaseAction(TypeSetContactTimezone, uuid),
		Timezone:   timezone,
	}
}

// Execute runs this action
func (a *SetContactTimezone) Execute(ctx context.Context, run flows.Run, step flows.Step, logModifier flows.ModifierCallback, logEvent flows.EventCallback) error {
	timezone, ok := run.EvaluateTemplate(a.Timezone, logEvent)
	timezone = strings.TrimSpace(timezone)

	if !ok {
		return nil
	}

	// timezone must be empty or valid timezone name
	var tz *time.Location
	var err error
	if timezone != "" {
		tz, err = time.LoadLocation(timezone)
		if err != nil {
			logEvent(events.NewError(fmt.Sprintf("unrecognized timezone: '%s'", timezone)))
			return nil
		}
	}

	a.applyModifier(run, modifiers.NewTimezone(tz), logModifier, logEvent)
	return nil
}
