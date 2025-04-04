package envs_test

import (
	"testing"
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

	// can create with valid values
	env, err = envs.ReadEnvironment([]byte(`{
		"date_format": "DD-MM-YYYY", 
		"time_format": "tt:mm:ss", 
		"allowed_languages": ["eng", "fra"], 
		"default_country": "RW", 
		"timezone": "Africa/Kigali"
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
	assert.Equal(t, envs.RedactionPolicyNone, env.RedactionPolicy())
	assert.Nil(t, env.LocationResolver())

	data, err := jsonx.Marshal(env)
	require.NoError(t, err)
	assert.Equal(t, string(data), `{"date_format":"DD-MM-YYYY","time_format":"tt:mm:ss","timezone":"Africa/Kigali","allowed_languages":["eng","fra"],"number_format":{"decimal_symbol":".","digit_grouping_symbol":","},"default_country":"RW","input_collation":"default","redaction_policy":"none"}`)
}

func TestEnvironmentEqual(t *testing.T) {
	env1, err := envs.ReadEnvironment([]byte(`{"date_format": "DD-MM-YYYY", "time_format": "tt:mm:ss", "timezone": "Africa/Kigali"}`))
	require.NoError(t, err)

	env2, err := envs.ReadEnvironment([]byte(`{"date_format": "DD-MM-YYYY", "time_format": "tt:mm:ss", "timezone": "Africa/Kigali"}`))
	require.NoError(t, err)

	env3, err := envs.ReadEnvironment([]byte(`{"date_format": "DD-MM-YYYY", "time_format": "tt:mm:ss", "timezone": "Africa/Kampala"}`))
	require.NoError(t, err)

	assert.True(t, env1.Equal(env2))
	assert.True(t, env2.Equal(env1))
	assert.False(t, env1.Equal(env3))

	// marshal and unmarshal env 1 again
	env1JSON, err := jsonx.Marshal(env1)
	require.NoError(t, err)
	env1, err = envs.ReadEnvironment(env1JSON)
	require.NoError(t, err)

	assert.True(t, env1.Equal(env2))
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
		Build()

	assert.Equal(t, envs.DateFormatDayMonthYear, env.DateFormat())
	assert.Equal(t, envs.TimeFormatHourMinuteSecond, env.TimeFormat())
	assert.Equal(t, kgl, env.Timezone())
	assert.Equal(t, []i18n.Language{i18n.Language("fra"), i18n.Language("eng")}, env.AllowedLanguages())
	assert.Equal(t, i18n.Country("RW"), env.DefaultCountry())
	assert.Equal(t, &envs.NumberFormat{DecimalSymbol: "'"}, env.NumberFormat())
	assert.Equal(t, envs.RedactionPolicyURNs, env.RedactionPolicy())
	assert.Nil(t, env.LocationResolver())
}
