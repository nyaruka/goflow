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

// NumberFormat describes how numbers should be parsed and formatted
type NumberFormat struct {
	DecimalSymbol       string `json:"decimal_symbol"`
	DigitGroupingSymbol string `json:"digit_grouping_symbol"`
}

// DefaultNumberFormat is the default number formatting, e.g. 1,234.567
var DefaultNumberFormat = &NumberFormat{DecimalSymbol: `.`, DigitGroupingSymbol: `,`}

// Environment defines the environment that the Excellent function is running in, this includes
// the timezone the user is in as well as the preferred date and time formats.
type Environment interface {
	DateFormat() DateFormat
	TimeFormat() TimeFormat
	Timezone() *time.Location
	DefaultLanguage() Language
	AllowedLanguages() []Language
	DefaultCountry() Country
	NumberFormat() *NumberFormat
	RedactionPolicy() RedactionPolicy

	// Convenience method to get the current time in the env timezone
	Now() time.Time

	// extensions to the engine can expect their own env values
	Extension(string) json.RawMessage

	Equal(Environment) bool
}

// NewDefaultEnvironment creates a new Environment with our usual defaults in the UTC timezone
func NewDefaultEnvironment() Environment {
	return &environment{
		dateFormat:       DateFormatYearMonthDay,
		timeFormat:       TimeFormatHourMinute,
		timezone:         time.UTC,
		defaultLanguage:  NilLanguage,
		allowedLanguages: nil,
		defaultCountry:   NilCountry,
		numberFormat:     DefaultNumberFormat,
		redactionPolicy:  RedactionPolicyNone,
	}
}

// NewEnvironment creates a new Environment with the passed in date and time formats and timezone
func NewEnvironment(dateFormat DateFormat, timeFormat TimeFormat, timezone *time.Location, defaultLanguage Language, allowedLanguages []Language, defaultCountry Country, numberFormat *NumberFormat, redactionPolicy RedactionPolicy) Environment {
	return &environment{
		dateFormat:       dateFormat,
		timeFormat:       timeFormat,
		timezone:         timezone,
		defaultLanguage:  defaultLanguage,
		allowedLanguages: allowedLanguages,
		defaultCountry:   defaultCountry,
		numberFormat:     numberFormat,
		redactionPolicy:  redactionPolicy,
	}
}

type environment struct {
	dateFormat       DateFormat
	timeFormat       TimeFormat
	timezone         *time.Location
	defaultLanguage  Language
	allowedLanguages []Language
	defaultCountry   Country
	numberFormat     *NumberFormat
	redactionPolicy  RedactionPolicy
	extensions       map[string]json.RawMessage
}

func (e *environment) DateFormat() DateFormat           { return e.dateFormat }
func (e *environment) TimeFormat() TimeFormat           { return e.timeFormat }
func (e *environment) Timezone() *time.Location         { return e.timezone }
func (e *environment) DefaultLanguage() Language        { return e.defaultLanguage }
func (e *environment) AllowedLanguages() []Language     { return e.allowedLanguages }
func (e *environment) DefaultCountry() Country          { return e.defaultCountry }
func (e *environment) NumberFormat() *NumberFormat      { return e.numberFormat }
func (e *environment) RedactionPolicy() RedactionPolicy { return e.redactionPolicy }
func (e *environment) Now() time.Time                   { return Now().In(e.Timezone()) }

func (e *environment) Extension(name string) json.RawMessage {
	return e.extensions[name]
}

// Equal returns true if this instance is equal to the given instance
func (e *environment) Equal(other Environment) bool {
	asJSON1, _ := json.Marshal(e)
	asJSON2, _ := json.Marshal(other)
	return string(asJSON1) == string(asJSON2)
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type envEnvelope struct {
	DateFormat       DateFormat                 `json:"date_format" validate:"required,date_format"`
	TimeFormat       TimeFormat                 `json:"time_format" validate:"required,time_format"`
	Timezone         string                     `json:"timezone" validate:"required"`
	DefaultLanguage  Language                   `json:"default_language,omitempty" validate:"omitempty,language"`
	AllowedLanguages []Language                 `json:"allowed_languages,omitempty" validate:"omitempty,dive,language"`
	NumberFormat     *NumberFormat              `json:"number_format,omitempty"`
	DefaultCountry   Country                    `json:"default_country,omitempty" validate:"omitempty,country"`
	RedactionPolicy  RedactionPolicy            `json:"redaction_policy" validate:"omitempty,eq=none|eq=urns"`
	Extensions       map[string]json.RawMessage `json:"extensions,omitempty"`
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
	env.defaultLanguage = envelope.DefaultLanguage
	env.allowedLanguages = envelope.AllowedLanguages
	env.defaultCountry = envelope.DefaultCountry
	env.extensions = envelope.Extensions

	if envelope.NumberFormat != nil {
		env.numberFormat = envelope.NumberFormat
	}

	tz, err := time.LoadLocation(envelope.Timezone)
	if err != nil {
		return nil, err
	}
	env.timezone = tz

	env.redactionPolicy = envelope.RedactionPolicy
	if env.redactionPolicy == "" {
		env.redactionPolicy = RedactionPolicyNone
	}

	return env, nil
}

// MarshalJSON marshals this environment into JSON
func (e *environment) MarshalJSON() ([]byte, error) {
	ee := &envEnvelope{
		DateFormat:       e.dateFormat,
		TimeFormat:       e.timeFormat,
		Timezone:         e.timezone.String(),
		DefaultLanguage:  e.defaultLanguage,
		AllowedLanguages: e.allowedLanguages,
		DefaultCountry:   e.defaultCountry,
		RedactionPolicy:  e.redactionPolicy,
		Extensions:       e.extensions,
	}
	if e.numberFormat != DefaultNumberFormat {
		ee.NumberFormat = e.numberFormat
	}

	return json.Marshal(ee)
}
