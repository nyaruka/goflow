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
	MaxValueLength() int

	// Convenience method to get the current time in the env timezone
	Now() time.Time

	// extensions to the engine can expect their own env values
	Extension(string) json.RawMessage

	Equal(Environment) bool
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
	maxValueLength   int
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
func (e *environment) MaxValueLength() int              { return e.maxValueLength }

func (e *environment) Now() time.Time { return Now().In(e.Timezone()) }

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
	DateFormat       DateFormat                 `json:"date_format" validate:"date_format"`
	TimeFormat       TimeFormat                 `json:"time_format" validate:"time_format"`
	Timezone         string                     `json:"timezone"`
	DefaultLanguage  Language                   `json:"default_language,omitempty" validate:"omitempty,language"`
	AllowedLanguages []Language                 `json:"allowed_languages,omitempty" validate:"omitempty,dive,language"`
	NumberFormat     *NumberFormat              `json:"number_format,omitempty"`
	DefaultCountry   Country                    `json:"default_country,omitempty" validate:"omitempty,country"`
	RedactionPolicy  RedactionPolicy            `json:"redaction_policy" validate:"omitempty,eq=none|eq=urns"`
	MaxValuelength   int                        `json:"max_value_length"`
	Extensions       map[string]json.RawMessage `json:"extensions,omitempty"`
}

// ReadEnvironment reads an environment from the given JSON
func ReadEnvironment(data json.RawMessage) (Environment, error) {
	// create new env with defaults
	env := NewEnvironmentBuilder().Build().(*environment)
	envelope := env.toEnvelope()

	if err := UnmarshalAndValidate(data, envelope); err != nil {
		return nil, err
	}

	env.dateFormat = envelope.DateFormat
	env.timeFormat = envelope.TimeFormat
	env.defaultLanguage = envelope.DefaultLanguage
	env.allowedLanguages = envelope.AllowedLanguages
	env.defaultCountry = envelope.DefaultCountry
	env.numberFormat = envelope.NumberFormat
	env.redactionPolicy = envelope.RedactionPolicy
	env.maxValueLength = envelope.MaxValuelength
	env.extensions = envelope.Extensions

	tz, err := time.LoadLocation(envelope.Timezone)
	if err != nil {
		return nil, err
	}
	env.timezone = tz

	return env, nil
}

func (e *environment) toEnvelope() *envEnvelope {
	return &envEnvelope{
		DateFormat:       e.dateFormat,
		TimeFormat:       e.timeFormat,
		Timezone:         e.timezone.String(),
		DefaultLanguage:  e.defaultLanguage,
		AllowedLanguages: e.allowedLanguages,
		DefaultCountry:   e.defaultCountry,
		NumberFormat:     e.numberFormat,
		RedactionPolicy:  e.redactionPolicy,
		MaxValuelength:   e.maxValueLength,
		Extensions:       e.extensions,
	}
}

// MarshalJSON marshals this environment into JSON
func (e *environment) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.toEnvelope())
}

//------------------------------------------------------------------------------------------
// Builder
//------------------------------------------------------------------------------------------

// EnvironmentBuilder is a builder for environments
type EnvironmentBuilder struct {
	env *environment
}

// NewEnvironmentBuilder creates a new environment builder
func NewEnvironmentBuilder() *EnvironmentBuilder {
	return &EnvironmentBuilder{
		env: &environment{
			dateFormat:       DateFormatYearMonthDay,
			timeFormat:       TimeFormatHourMinute,
			timezone:         time.UTC,
			defaultLanguage:  NilLanguage,
			allowedLanguages: nil,
			defaultCountry:   NilCountry,
			numberFormat:     DefaultNumberFormat,
			maxValueLength:   640,
			redactionPolicy:  RedactionPolicyNone,
		},
	}
}

// WithDateFormat sets the date format
func (b *EnvironmentBuilder) WithDateFormat(dateFormat DateFormat) *EnvironmentBuilder {
	b.env.dateFormat = dateFormat
	return b
}

// WithTimeFormat sets the time format
func (b *EnvironmentBuilder) WithTimeFormat(timeFormat TimeFormat) *EnvironmentBuilder {
	b.env.timeFormat = timeFormat
	return b
}

func (b *EnvironmentBuilder) WithTimezone(timezone *time.Location) *EnvironmentBuilder {
	b.env.timezone = timezone
	return b
}

func (b *EnvironmentBuilder) WithDefaultLanguage(defaultLanguage Language) *EnvironmentBuilder {
	b.env.defaultLanguage = defaultLanguage
	return b
}

func (b *EnvironmentBuilder) WithAllowedLanguages(allowedLanguages []Language) *EnvironmentBuilder {
	b.env.allowedLanguages = allowedLanguages
	return b
}

func (b *EnvironmentBuilder) WithDefaultCountry(defaultCountry Country) *EnvironmentBuilder {
	b.env.defaultCountry = defaultCountry
	return b
}

func (b *EnvironmentBuilder) WithNumberFormat(numberFormat *NumberFormat) *EnvironmentBuilder {
	b.env.numberFormat = numberFormat
	return b
}

func (b *EnvironmentBuilder) WithRedactionPolicy(redactionPolicy RedactionPolicy) *EnvironmentBuilder {
	b.env.redactionPolicy = redactionPolicy
	return b
}

func (b *EnvironmentBuilder) WithMaxValueLength(maxValueLength int) *EnvironmentBuilder {
	b.env.maxValueLength = maxValueLength
	return b
}

// Build returns the final environment
func (b *EnvironmentBuilder) Build() Environment { return b.env }
