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

	// can't create empty environment
	env, err := utils.ReadEnvironment(json.RawMessage(`{}`))
	assert.Error(t, err)

	// can't create with invalid date format
	env, err = utils.ReadEnvironment(json.RawMessage(`{"date_format": "YYYYYYYYYYY", "time_format": "tt:mm:ss", "timezone": "Africa/Kigali"}`))
	assert.Error(t, err)

	// can't create with invalid time format
	env, err = utils.ReadEnvironment(json.RawMessage(`{"date_format": "DD-MM-YYYY", "time_format": "tttttt", "timezone": "Africa/Kigali"}`))
	assert.Error(t, err)

	// can't create with invalid timzeone
	env, err = utils.ReadEnvironment(json.RawMessage(`{"date_format": "DD-MM-YYYY", "time_format": "tttttt", "timezone": "Cuenca"}`))
	assert.Error(t, err)

	// can create with valid values
	env, err = utils.ReadEnvironment(json.RawMessage(`{"date_format": "DD-MM-YYYY", "time_format": "tt:mm:ss", "timezone": "Africa/Kigali"}`))
	assert.NoError(t, err)
	assert.Equal(t, utils.DateFormatDayMonthYear, env.DateFormat())
	assert.Equal(t, utils.TimeFormatHourMinuteSecond, env.TimeFormat())
	assert.Equal(t, kgl, env.Timezone())
	assert.Equal(t, utils.LanguageList{}, env.Languages())

	data, err := json.Marshal(env)
	require.NoError(t, err)
	assert.Equal(t, string(data), `{"date_format":"DD-MM-YYYY","time_format":"tt:mm:ss","timezone":"Africa/Kigali","languages":[]}`)
}
