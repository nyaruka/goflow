package modifiers

import (
	"encoding/json"
	"time"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	RegisterType(TypeTimezone, readTimezoneModifier)
}

// TypeTimezone is the type of our timezone modifier
const TypeTimezone string = "timezone"

// TimezoneModifier modifies the timezone of a contact
type TimezoneModifier struct {
	baseModifier

	timezone *time.Location
}

// NewTimezoneModifier creates a new timezone modifier
func NewTimezoneModifier(timezone *time.Location) *TimezoneModifier {
	return &TimezoneModifier{
		baseModifier: newBaseModifier(TypeTimezone),
		timezone:     timezone,
	}
}

// Apply applies this modification to the given contact
func (m *TimezoneModifier) Apply(env utils.Environment, assets flows.SessionAssets, contact *flows.Contact, log flows.EventCallback) {
	if !timezonesEqual(contact.Timezone(), m.timezone) {
		contact.SetTimezone(m.timezone)
		log(events.NewContactTimezoneChangedEvent(m.timezone))
		m.reevaluateDynamicGroups(env, assets, contact, log)
	}
}

func timezonesEqual(tz1 *time.Location, tz2 *time.Location) bool {
	return (tz1 == nil && tz2 == nil) || (tz1 != nil && tz2 != nil && tz1.String() == tz2.String())
}

var _ flows.Modifier = (*TimezoneModifier)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type timezoneModifierEnvelope struct {
	utils.TypedEnvelope
	Timezone string `json:"timezone"`
}

func readTimezoneModifier(assets flows.SessionAssets, data json.RawMessage) (flows.Modifier, error) {
	e := &timezoneModifierEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	var tz *time.Location
	if e.Timezone != "" {
		var err error
		tz, err = time.LoadLocation(e.Timezone)
		if err != nil {
			return nil, err
		}
	}

	return NewTimezoneModifier(tz), nil
}

func (m *TimezoneModifier) MarshalJSON() ([]byte, error) {
	tzName := ""
	if m.timezone != nil {
		tzName = m.timezone.String()
	}
	return json.Marshal(&timezoneModifierEnvelope{TypedEnvelope: utils.TypedEnvelope{Type: m.Type()}, Timezone: tzName})
}
