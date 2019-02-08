package utils_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/nyaruka/goflow/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnvironmentMarshaling(t *testing.T) {
	kgl, err := time.LoadLocation("Africa/Kigali")
	require.NoError(t, err)

	// can't create with invalid date format
	env, err := utils.ReadEnvironment(json.RawMessage(`{"date_format": "YYYYYYYYYYY", "time_format": "tt:mm:ss", "timezone": "Africa/Kigali"}`))
	assert.Error(t, err)

	// can't create with invalid time format
	env, err = utils.ReadEnvironment(json.RawMessage(`{"date_format": "DD-MM-YYYY", "time_format": "tttttt", "timezone": "Africa/Kigali"}`))
	assert.Error(t, err)

	// can't create with invalid language
	env, err = utils.ReadEnvironment(json.RawMessage(`{"date_format": "DD-MM-YYYY", "time_format": "tttttt", "default_language": "elvish"}`))
	assert.Error(t, err)

	// can't create with invalid country
	env, err = utils.ReadEnvironment(json.RawMessage(`{"date_format": "DD-MM-YYYY", "time_format": "tttttt", "default_country": "Narnia"}`))
	assert.Error(t, err)

	// can't create with invalid timzeone
	env, err = utils.ReadEnvironment(json.RawMessage(`{"date_format": "DD-MM-YYYY", "time_format": "tttttt", "timezone": "Cuenca"}`))
	assert.Error(t, err)

	// can't create with non-map extensions field
	env, err = utils.ReadEnvironment(json.RawMessage(`{"date_format": "DD-MM-YYYY", "time_format": "tttttt", "timezone": "Cuenca", "extensions": []}`))
	assert.Error(t, err)

	// empty environment uses all defaults
	env, err = utils.ReadEnvironment(json.RawMessage(`{}`))
	assert.NoError(t, err)
	assert.Equal(t, utils.DateFormatYearMonthDay, env.DateFormat())
	assert.Equal(t, utils.TimeFormatHourMinute, env.TimeFormat())
	assert.Equal(t, utils.DefaultNumberFormat, env.NumberFormat())
	assert.Equal(t, 640, env.MaxValueLength())

	// can create with valid values
	env, err = utils.ReadEnvironment(json.RawMessage(`{"date_format": "DD-MM-YYYY", "time_format": "tt:mm:ss", "default_language": "eng", "allowed_languages": ["eng", "fra"], "default_country": "RW", "timezone": "Africa/Kigali", "extensions": {"foo":{"bar":1234}}}`))
	assert.NoError(t, err)
	assert.Equal(t, utils.DateFormatDayMonthYear, env.DateFormat())
	assert.Equal(t, utils.TimeFormatHourMinuteSecond, env.TimeFormat())
	assert.Equal(t, kgl, env.Timezone())
	assert.Equal(t, utils.Language("eng"), env.DefaultLanguage())
	assert.Equal(t, []utils.Language{utils.Language("eng"), utils.Language("fra")}, env.AllowedLanguages())
	assert.Equal(t, utils.Country("RW"), env.DefaultCountry())
	assert.Equal(t, json.RawMessage(`{"bar":1234}`), env.Extension("foo"))

	data, err := json.Marshal(env)
	require.NoError(t, err)
	assert.Equal(t, string(data), `{"date_format":"DD-MM-YYYY","time_format":"tt:mm:ss","timezone":"Africa/Kigali","default_language":"eng","allowed_languages":["eng","fra"],"number_format":{"decimal_symbol":".","digit_grouping_symbol":","},"default_country":"RW","redaction_policy":"none","max_value_length":640,"extensions":{"foo":{"bar":1234}}}`)
}

func TestEnvironmentEqual(t *testing.T) {
	env1, err := utils.ReadEnvironment(json.RawMessage(`{"date_format": "DD-MM-YYYY", "time_format": "tt:mm:ss", "timezone": "Africa/Kigali", "extensions": {"foo":{"bar":1234}}}`))
	require.NoError(t, err)

	env2, err := utils.ReadEnvironment(json.RawMessage(`{"date_format": "DD-MM-YYYY", "time_format": "tt:mm:ss", "timezone": "Africa/Kigali", "extensions": {"foo":{"bar":1234}}}`))
	require.NoError(t, err)

	env3, err := utils.ReadEnvironment(json.RawMessage(`{"date_format": "DD-MM-YYYY", "time_format": "tt:mm:ss", "timezone": "Africa/Kigali", "extensions": {"foo":{"bar":2345}}}`))
	require.NoError(t, err)

	assert.True(t, env1.Equal(env2))
	assert.True(t, env2.Equal(env1))
	assert.False(t, env1.Equal(env3))

	// marshal and unmarshal env 1 again
	env1JSON, err := json.Marshal(env1)
	require.NoError(t, err)
	env1, err = utils.ReadEnvironment(env1JSON)
	require.NoError(t, err)

	assert.True(t, env1.Equal(env2))
}

func TestEnvironmentBuilder(t *testing.T) {
	kgl, err := time.LoadLocation("Africa/Kigali")
	require.NoError(t, err)

	env := utils.NewEnvironmentBuilder().
		WithDateFormat(utils.DateFormatDayMonthYear).
		WithTimeFormat(utils.TimeFormatHourMinuteSecond).
		WithTimezone(kgl).
		WithDefaultLanguage(utils.Language("fra")).
		WithAllowedLanguages([]utils.Language{utils.Language("fra"), utils.Language("eng")}).
		WithDefaultCountry(utils.Country("RW")).
		WithNumberFormat(&utils.NumberFormat{DecimalSymbol: "'"}).
		WithRedactionPolicy(utils.RedactionPolicyURNs).
		WithMaxValueLength(1024).
		Build()

	assert.Equal(t, utils.DateFormatDayMonthYear, env.DateFormat())
	assert.Equal(t, utils.TimeFormatHourMinuteSecond, env.TimeFormat())
	assert.Equal(t, kgl, env.Timezone())
	assert.Equal(t, utils.Language("fra"), env.DefaultLanguage())
	assert.Equal(t, []utils.Language{utils.Language("fra"), utils.Language("eng")}, env.AllowedLanguages())
	assert.Equal(t, utils.Country("RW"), env.DefaultCountry())
	assert.Equal(t, &utils.NumberFormat{DecimalSymbol: "'"}, env.NumberFormat())
	assert.Equal(t, utils.RedactionPolicyURNs, env.RedactionPolicy())
	assert.Equal(t, 1024, env.MaxValueLength())
}
