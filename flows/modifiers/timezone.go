package modifiers

import (
	"time"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	registerType(TypeTimezone, readTimezone)
}

// TypeTimezone is the type of our timezone modifier
const TypeTimezone string = "timezone"

// Timezone modifies the timezone of a contact
type Timezone struct {
	baseModifier

	timezone *time.Location
}

// NewTimezone creates a new timezone modifier
func NewTimezone(timezone *time.Location) *Timezone {
	return &Timezone{
		baseModifier: newBaseModifier(TypeTimezone),
		timezone:     timezone,
	}
}

// Apply applies this modification to the given contact
func (m *Timezone) Apply(eng flows.Engine, env envs.Environment, sa flows.SessionAssets, contact *flows.Contact, log flows.EventCallback) (bool, error) {
	if !timezonesEqual(contact.Timezone(), m.timezone) {
		contact.SetTimezone(m.timezone)
		log(events.NewContactTimezoneChanged(m.timezone))
		return true, nil
	}
	return false, nil
}

func timezonesEqual(tz1 *time.Location, tz2 *time.Location) bool {
	return (tz1 == nil && tz2 == nil) || (tz1 != nil && tz2 != nil && tz1.String() == tz2.String())
}

var _ flows.Modifier = (*Timezone)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type timezoneEnvelope struct {
	utils.TypedEnvelope

	Timezone string `json:"timezone"`
}

func readTimezone(sa flows.SessionAssets, data []byte, missing assets.MissingCallback) (flows.Modifier, error) {
	e := &timezoneEnvelope{}
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

	return NewTimezone(tz), nil
}

func (m *Timezone) MarshalJSON() ([]byte, error) {
	tzName := ""
	if m.timezone != nil {
		tzName = m.timezone.String()
	}
	return jsonx.Marshal(&timezoneEnvelope{
		TypedEnvelope: utils.TypedEnvelope{Type: m.Type()},
		Timezone:      tzName,
	})
}
