package modifiers

import (
	"time"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
)

func init() {
	RegisterType(TypeTimezone, func() Modifier { return &TimezoneModifier{} })
}

// TypeTimezone is the type of our timezone modifier
const TypeTimezone string = "timezone"

// TimezoneModifier modifies the timezone of a contact
type TimezoneModifier struct {
	baseModifier

	Timezone *time.Location `json:"timezone"`
}

// NewTimezoneModifier creates a new timezone modifier
func NewTimezoneModifier(timezone *time.Location) *TimezoneModifier {
	return &TimezoneModifier{
		baseModifier: newBaseModifier(TypeTimezone),
		Timezone:     timezone,
	}
}

// Apply applies this modification to the given contact
func (m *TimezoneModifier) Apply(assets flows.SessionAssets, contact *flows.Contact) flows.Event {
	if !timezonesEqual(contact.Timezone(), m.Timezone) {
		contact.SetTimezone(m.Timezone)
		return events.NewContactTimezoneChangedEvent(m.Timezone)
	}
	return nil
}

func timezonesEqual(tz1 *time.Location, tz2 *time.Location) bool {
	return (tz1 == nil && tz2 == nil) || (tz1 != nil && tz2 != nil && tz1.String() == tz2.String())
}

var _ Modifier = (*TimezoneModifier)(nil)
