package utils

import (
	"encoding/json"
	"time"
)

// Environment defines the Environment that the Excellent function is running in, this includes
// the timezone the user is in as well as the preferred date and time formats.
type Environment interface {
	DateFormat() DateFormat
	SetDateFormat(DateFormat)

	TimeFormat() TimeFormat
	SetTimeFormat(TimeFormat)

	Timezone() *time.Location
	SetTimezone(*time.Location)

	Languages() LanguageList
	Locations() (*LocationHierarchy, error)
}

// NewDefaultEnvironment creates a new Environment with our usual defaults in the UTC timezone
func NewDefaultEnvironment() Environment {
	return &environment{
		dateFormat: DateFormat_yyyy_MM_dd,
		timeFormat: TimeFormat_HH_mm,
		timezone:   time.UTC,
		languages:  LanguageList{},
	}
}

// NewEnvironment creates a new Environment with the passed in date and time formats and timezone
func NewEnvironment(dateFormat DateFormat, timeFormat TimeFormat, timezone *time.Location, languages LanguageList) Environment {
	if timezone == nil {
		timezone = time.UTC
	}
	return &environment{
		dateFormat: dateFormat,
		timeFormat: timeFormat,
		timezone:   timezone,
		languages:  languages,
	}
}

type environment struct {
	dateFormat DateFormat
	timeFormat TimeFormat
	timezone   *time.Location
	languages  LanguageList
}

func (e *environment) DateFormat() DateFormat              { return e.dateFormat }
func (e *environment) SetDateFormat(dateFormat DateFormat) { e.dateFormat = dateFormat }

func (e *environment) TimeFormat() TimeFormat              { return e.timeFormat }
func (e *environment) SetTimeFormat(timeFormat TimeFormat) { e.timeFormat = timeFormat }

func (e *environment) Timezone() *time.Location { return e.timezone }
func (e *environment) SetTimezone(timezone *time.Location) {
	if timezone == nil {
		e.timezone = time.UTC
	} else {
		e.timezone = timezone
	}
}

func (e *environment) Languages() LanguageList                { return e.languages }
func (e *environment) Locations() (*LocationHierarchy, error) { return nil, nil }

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type envEnvelope struct {
	DateFormat DateFormat   `json:"date_format"`
	TimeFormat TimeFormat   `json:"time_format"`
	Timezone   string       `json:"timezone"`
	Languages  LanguageList `json:"languages"`
}

// ReadEnvironment reads an environment from the given JSON
func ReadEnvironment(data json.RawMessage) (*environment, error) {
	env := NewDefaultEnvironment().(*environment)

	var envelope envEnvelope
	var err error

	err = json.Unmarshal(data, &envelope)
	if err != nil {
		return nil, err
	}

	env.dateFormat = envelope.DateFormat
	env.timeFormat = envelope.TimeFormat
	tz, err := time.LoadLocation(envelope.Timezone)
	if err != nil {
		return nil, err
	}
	env.timezone = tz
	env.languages = envelope.Languages
	return env, nil
}

func (e *environment) MarshalJSON() ([]byte, error) {
	ee := envEnvelope{
		DateFormat: e.dateFormat,
		TimeFormat: e.timeFormat,
		Timezone:   e.timezone.String(),
		Languages:  e.languages,
	}
	return json.Marshal(ee)
}
