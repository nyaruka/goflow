package utils

import (
	"encoding/json"
	"time"
)

type RedactionPolicy string

const (
	RedactionPolicyNone RedactionPolicy = "none"
	RedactionPolicyURNs RedactionPolicy = "urns"
)

// Environment defines the environment that the Excellent function is running in, this includes
// the timezone the user is in as well as the preferred date and time formats.
type Environment interface {
	DateFormat() DateFormat
	TimeFormat() TimeFormat
	Timezone() *time.Location
	Languages() LanguageList
	RedactionPolicy() RedactionPolicy

	// Convenience method to get the current time in the env timezone
	Now() time.Time

	// extensions to the engine can expect their own env values
	Extension(string) json.RawMessage
}

// NewDefaultEnvironment creates a new Environment with our usual defaults in the UTC timezone
func NewDefaultEnvironment() Environment {
	return &environment{
		dateFormat:      DateFormatYearMonthDay,
		timeFormat:      TimeFormatHourMinute,
		timezone:        time.UTC,
		languages:       LanguageList{},
		redactionPolicy: RedactionPolicyNone,
	}
}

// NewEnvironment creates a new Environment with the passed in date and time formats and timezone
func NewEnvironment(dateFormat DateFormat, timeFormat TimeFormat, timezone *time.Location, languages LanguageList, redactionPolicy RedactionPolicy) Environment {
	return &environment{
		dateFormat:      dateFormat,
		timeFormat:      timeFormat,
		timezone:        timezone,
		languages:       languages,
		redactionPolicy: redactionPolicy,
	}
}

type environment struct {
	dateFormat      DateFormat
	timeFormat      TimeFormat
	timezone        *time.Location
	languages       LanguageList
	redactionPolicy RedactionPolicy
	extensions      map[string]json.RawMessage
}

func (e *environment) DateFormat() DateFormat           { return e.dateFormat }
func (e *environment) TimeFormat() TimeFormat           { return e.timeFormat }
func (e *environment) Timezone() *time.Location         { return e.timezone }
func (e *environment) Languages() LanguageList          { return e.languages }
func (e *environment) RedactionPolicy() RedactionPolicy { return e.redactionPolicy }
func (e *environment) Now() time.Time                   { return Now().In(e.Timezone()) }

func (e *environment) Extension(name string) json.RawMessage {
	return e.extensions[name]
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type envEnvelope struct {
	DateFormat      DateFormat                 `json:"date_format" validate:"required,date_format"`
	TimeFormat      TimeFormat                 `json:"time_format" validate:"required,time_format"`
	Timezone        string                     `json:"timezone" validate:"required"`
	Languages       LanguageList               `json:"languages"`
	RedactionPolicy RedactionPolicy            `json:"redaction_policy" validate:"omitempty,eq=none|eq=urns"`
	Extensions      map[string]json.RawMessage `json:"extensions,omitempty"`
}

// ReadEnvironment reads an environment from the given JSON
func ReadEnvironment(data json.RawMessage) (Environment, error) {
	env := NewDefaultEnvironment().(*environment)

	var envelope envEnvelope
	if err := UnmarshalAndValidate(data, &envelope); err != nil {
		return nil, err
	}

	env.dateFormat = envelope.DateFormat
	env.timeFormat = envelope.TimeFormat
	env.extensions = envelope.Extensions

	tz, err := time.LoadLocation(envelope.Timezone)
	if err != nil {
		return nil, err
	}
	env.timezone = tz

	if envelope.Languages != nil {
		env.languages = envelope.Languages
	}

	env.redactionPolicy = envelope.RedactionPolicy
	if env.redactionPolicy == "" {
		env.redactionPolicy = RedactionPolicyNone
	}

	return env, nil
}

// MarshalJSON marshals this environment into JSON
func (e *environment) MarshalJSON() ([]byte, error) {
	ee := &envEnvelope{
		DateFormat:      e.dateFormat,
		TimeFormat:      e.timeFormat,
		Timezone:        e.timezone.String(),
		Languages:       e.languages,
		RedactionPolicy: e.redactionPolicy,
		Extensions:      e.extensions,
	}
	return json.Marshal(ee)
}
