package envs_test

import (
	"testing"
	"text/template"
	"time"

	"github.com/nyaruka/gocommon/i18n"
	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/envs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnvironmentMarshaling(t *testing.T) {
	kgl, err := time.LoadLocation("Africa/Kigali")
	require.NoError(t, err)

	// can't create with invalid date format
	_, err = envs.ReadEnvironment([]byte(`{"date_format": "YYYYYYYYYYY", "time_format": "tt:mm:ss", "timezone": "Africa/Kigali"}`))
	assert.Error(t, err)

	// can't create with invalid time format
	_, err = envs.ReadEnvironment([]byte(`{"date_format": "DD-MM-YYYY", "time_format": "tttttt", "timezone": "Africa/Kigali"}`))
	assert.Error(t, err)

	// can't create with invalid language
	_, err = envs.ReadEnvironment([]byte(`{"date_format": "DD-MM-YYYY", "time_format": "tttttt", "allowed_languages": ["elvish"]}`))
	assert.Error(t, err)

	// can't create with invalid country
	_, err = envs.ReadEnvironment([]byte(`{"date_format": "DD-MM-YYYY", "time_format": "tttttt", "default_country": "Narnia"}`))
	assert.Error(t, err)

	// can't create with invalid timzeone
	_, err = envs.ReadEnvironment([]byte(`{"date_format": "DD-MM-YYYY", "time_format": "tttttt", "timezone": "Cuenca"}`))
	assert.Error(t, err)

	// empty environment uses all defaults
	env, err := envs.ReadEnvironment([]byte(`{}`))
	assert.NoError(t, err)
	assert.Equal(t, envs.DateFormatYearMonthDay, env.DateFormat())
	assert.Equal(t, envs.TimeFormatHourMinute, env.TimeFormat())
	assert.Equal(t, envs.DefaultNumberFormat, env.NumberFormat())
	assert.Equal(t, i18n.NilLanguage, env.DefaultLanguage())
	assert.Nil(t, env.AllowedLanguages())
	assert.Equal(t, i18n.NilCountry, env.DefaultCountry())
	assert.Nil(t, env.LocationResolver())
	assert.Equal(t, envs.RedactionPolicyNone, env.RedactionPolicy())
	assert.Equal(t, [4]uint32{0xA3B1C, 0xD2E3F, 0x1A2B3, 0xC0FFEE}, env.ObfuscationKey())

	// can create with valid values
	env, err = envs.ReadEnvironment([]byte(`{
		"date_format": "DD-MM-YYYY", 
		"time_format": "tt:mm:ss", 
		"allowed_languages": ["eng", "fra"], 
		"default_country": "RW", 
		"timezone": "Africa/Kigali",
		"redaction_policy": "urns",
		"obfuscation_key": [123456, 234567, 345678, 456789]
	}`))
	assert.NoError(t, err)
	assert.Equal(t, envs.DateFormatDayMonthYear, env.DateFormat())
	assert.Equal(t, envs.TimeFormatHourMinuteSecond, env.TimeFormat())
	assert.Equal(t, kgl, env.Timezone())
	assert.Equal(t, i18n.Language("eng"), env.DefaultLanguage())
	assert.Equal(t, []i18n.Language{i18n.Language("eng"), i18n.Language("fra")}, env.AllowedLanguages())
	assert.Equal(t, i18n.Country("RW"), env.DefaultCountry())
	assert.Equal(t, i18n.Locale("eng-RW"), env.DefaultLocale())
	assert.Equal(t, envs.CollationDefault, env.InputCollation())
	assert.Equal(t, envs.RedactionPolicyURNs, env.RedactionPolicy())
	assert.Equal(t, [4]uint32{123456, 234567, 345678, 456789}, env.ObfuscationKey())
	assert.Nil(t, env.LocationResolver())

	data, err := jsonx.Marshal(env)
	require.NoError(t, err)
	assert.Equal(t, string(data), `{"date_format":"DD-MM-YYYY","time_format":"tt:mm:ss","timezone":"Africa/Kigali","allowed_languages":["eng","fra"],"number_format":{"decimal_symbol":".","digit_grouping_symbol":","},"default_country":"RW","input_collation":"default","redaction_policy":"urns","obfuscation_key":[123456,234567,345678,456789]}`)
}

func TestEnvironmentBuilder(t *testing.T) {
	kgl, err := time.LoadLocation("Africa/Kigali")
	require.NoError(t, err)

	env := envs.NewBuilder().
		WithDateFormat(envs.DateFormatDayMonthYear).
		WithTimeFormat(envs.TimeFormatHourMinuteSecond).
		WithTimezone(kgl).
		WithAllowedLanguages("fra", "eng").
		WithDefaultCountry(i18n.Country("RW")).
		WithNumberFormat(&envs.NumberFormat{DecimalSymbol: "'"}).
		WithRedactionPolicy(envs.RedactionPolicyURNs).
		WithObfuscationKey([4]uint32{123456, 234567, 345678, 456789}).
		WithPromptResolver(envs.NewPromptResolver(map[string]*template.Template{"hello": template.Must(template.New("").Parse("Say hello"))})).
		Build()

	assert.Equal(t, envs.DateFormatDayMonthYear, env.DateFormat())
	assert.Equal(t, envs.TimeFormatHourMinuteSecond, env.TimeFormat())
	assert.Equal(t, kgl, env.Timezone())
	assert.Equal(t, []i18n.Language{i18n.Language("fra"), i18n.Language("eng")}, env.AllowedLanguages())
	assert.Equal(t, i18n.Country("RW"), env.DefaultCountry())
	assert.Equal(t, &envs.NumberFormat{DecimalSymbol: "'"}, env.NumberFormat())
	assert.Equal(t, envs.RedactionPolicyURNs, env.RedactionPolicy())
	assert.Equal(t, [4]uint32{123456, 234567, 345678, 456789}, env.ObfuscationKey())
	assert.Nil(t, env.LocationResolver())
	assert.Nil(t, env.LLMPrompt("xxxx"))
	assert.NotNil(t, env.LLMPrompt("hello"))

	// using defaults
	env = envs.NewBuilder().Build()

	assert.Equal(t, envs.DateFormatYearMonthDay, env.DateFormat())
	assert.Equal(t, envs.TimeFormatHourMinute, env.TimeFormat())
	assert.Equal(t, time.UTC, env.Timezone())
	assert.Equal(t, []i18n.Language(nil), env.AllowedLanguages())
	assert.Equal(t, i18n.NilCountry, env.DefaultCountry())
	assert.Equal(t, &envs.NumberFormat{DecimalSymbol: ".", DigitGroupingSymbol: ","}, env.NumberFormat())
	assert.Equal(t, envs.RedactionPolicyNone, env.RedactionPolicy())
	assert.Equal(t, [4]uint32{670492, 863807, 107187, 12648430}, env.ObfuscationKey())
	assert.Nil(t, env.LocationResolver())
	assert.Nil(t, env.LLMPrompt("xxxx"))
	assert.Nil(t, env.LLMPrompt("hello"))
}
