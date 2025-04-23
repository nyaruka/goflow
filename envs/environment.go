package envs

import (
	"text/template"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/nyaruka/gocommon/dates"
	"github.com/nyaruka/gocommon/i18n"
	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	utils.RegisterValidatorTag("language", validateLanguage, func(validator.FieldError) string {
		return "is not a valid language code"
	})
	utils.RegisterValidatorAlias("country", "len=2", func(validator.FieldError) string {
		return "is not a valid country code"
	})
}

func validateLanguage(fl validator.FieldLevel) bool {
	_, err := i18n.ParseLanguage(fl.Field().String())
	return err == nil
}

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
	AllowedLanguages() []i18n.Language
	DefaultCountry() i18n.Country
	NumberFormat() *NumberFormat
	InputCollation() Collation
	RedactionPolicy() RedactionPolicy

	// non-marshalled properties
	LocationResolver() LocationResolver
	LLMPrompt(string) *template.Template

	// utility methods
	DefaultLanguage() i18n.Language
	DefaultLocale() i18n.Locale
	Now() time.Time // current time in the env timezone

	Equal(Environment) bool
}

type environment struct {
	dateFormat       DateFormat
	timeFormat       TimeFormat
	timezone         *time.Location
	allowedLanguages []i18n.Language
	defaultCountry   i18n.Country
	numberFormat     *NumberFormat
	redactionPolicy  RedactionPolicy
	inputCollation   Collation
	locationResolver LocationResolver
	promptResolver   PromptResolver
}

func (e *environment) DateFormat() DateFormat                   { return e.dateFormat }
func (e *environment) TimeFormat() TimeFormat                   { return e.timeFormat }
func (e *environment) Timezone() *time.Location                 { return e.timezone }
func (e *environment) AllowedLanguages() []i18n.Language        { return e.allowedLanguages }
func (e *environment) DefaultCountry() i18n.Country             { return e.defaultCountry }
func (e *environment) NumberFormat() *NumberFormat              { return e.numberFormat }
func (e *environment) InputCollation() Collation                { return e.inputCollation }
func (e *environment) RedactionPolicy() RedactionPolicy         { return e.redactionPolicy }
func (e *environment) LocationResolver() LocationResolver       { return e.locationResolver }
func (e *environment) LLMPrompt(name string) *template.Template { return e.promptResolver(name) }

// DefaultLanguage is the first allowed language
func (e *environment) DefaultLanguage() i18n.Language {
	if len(e.allowedLanguages) > 0 {
		return e.allowedLanguages[0]
	}
	return i18n.NilLanguage
}

// DefaultLocale combines the default languages and countries into a locale
func (e *environment) DefaultLocale() i18n.Locale {
	return i18n.NewLocale(e.DefaultLanguage(), e.DefaultCountry())
}

// Now gets the current time in the eonvironment's timezone
func (e *environment) Now() time.Time { return dates.Now().In(e.Timezone()) }

// Equal returns true if this instance is equal to the given instance
func (e *environment) Equal(other Environment) bool {
	asJSON1, _ := jsonx.Marshal(e)
	asJSON2, _ := jsonx.Marshal(other)
	return string(asJSON1) == string(asJSON2)
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type envEnvelope struct {
	DateFormat       DateFormat      `json:"date_format" validate:"date_format"`
	TimeFormat       TimeFormat      `json:"time_format" validate:"time_format"`
	Timezone         string          `json:"timezone"`
	AllowedLanguages []i18n.Language `json:"allowed_languages,omitempty" validate:"omitempty,dive,language"`
	NumberFormat     *NumberFormat   `json:"number_format,omitempty"`
	DefaultCountry   i18n.Country    `json:"default_country,omitempty" validate:"omitempty,country"`
	InputCollation   Collation       `json:"input_collation"`
	RedactionPolicy  RedactionPolicy `json:"redaction_policy" validate:"omitempty,eq=none|eq=urns"`
}

// ReadEnvironment reads an environment from the given JSON
func ReadEnvironment(data []byte) (Environment, error) {
	// create new env with defaults
	env := NewBuilder().Build().(*environment)
	envelope := env.toEnvelope()

	if err := utils.UnmarshalAndValidate(data, envelope); err != nil {
		return nil, err
	}

	env.dateFormat = envelope.DateFormat
	env.timeFormat = envelope.TimeFormat
	env.allowedLanguages = envelope.AllowedLanguages
	env.defaultCountry = envelope.DefaultCountry
	env.numberFormat = envelope.NumberFormat
	env.inputCollation = envelope.InputCollation
	env.redactionPolicy = envelope.RedactionPolicy

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
		AllowedLanguages: e.allowedLanguages,
		DefaultCountry:   e.defaultCountry,
		NumberFormat:     e.numberFormat,
		InputCollation:   e.inputCollation,
		RedactionPolicy:  e.redactionPolicy,
	}
}

// MarshalJSON marshals this environment into JSON
func (e *environment) MarshalJSON() ([]byte, error) {
	return jsonx.Marshal(e.toEnvelope())
}

//------------------------------------------------------------------------------------------
// Builder
//------------------------------------------------------------------------------------------

// EnvironmentBuilder is a builder for environments
type EnvironmentBuilder struct {
	env *environment
}

// NewEnvironmentBuilder creates a new environment builder
func NewBuilder() *EnvironmentBuilder {
	return &EnvironmentBuilder{
		env: &environment{
			dateFormat:       DateFormatYearMonthDay,
			timeFormat:       TimeFormatHourMinute,
			timezone:         time.UTC,
			allowedLanguages: nil,
			defaultCountry:   i18n.NilCountry,
			numberFormat:     DefaultNumberFormat,
			inputCollation:   CollationDefault,
			redactionPolicy:  RedactionPolicyNone,
			promptResolver:   EmptyPromptResolver,
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

func (b *EnvironmentBuilder) WithAllowedLanguages(allowedLanguages ...i18n.Language) *EnvironmentBuilder {
	b.env.allowedLanguages = allowedLanguages
	return b
}

func (b *EnvironmentBuilder) WithDefaultCountry(defaultCountry i18n.Country) *EnvironmentBuilder {
	b.env.defaultCountry = defaultCountry
	return b
}

func (b *EnvironmentBuilder) WithNumberFormat(numberFormat *NumberFormat) *EnvironmentBuilder {
	b.env.numberFormat = numberFormat
	return b
}

func (b *EnvironmentBuilder) WithInputCollation(col Collation) *EnvironmentBuilder {
	b.env.inputCollation = col
	return b
}

func (b *EnvironmentBuilder) WithRedactionPolicy(redactionPolicy RedactionPolicy) *EnvironmentBuilder {
	b.env.redactionPolicy = redactionPolicy
	return b
}

func (b *EnvironmentBuilder) WithLocationResolver(resolver LocationResolver) *EnvironmentBuilder {
	b.env.locationResolver = resolver
	return b
}

func (b *EnvironmentBuilder) WithPromptResolver(resolver PromptResolver) *EnvironmentBuilder {
	b.env.promptResolver = resolver
	return b
}

// Build returns the final environment
func (b *EnvironmentBuilder) Build() Environment { return b.env }
