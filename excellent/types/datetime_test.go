package types_test

import (
	"testing"
	"time"

	"github.com/nyaruka/gocommon/dates"
	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestXDateTime(t *testing.T) {
	env := envs.NewBuilder().WithDateFormat(envs.DateFormatDayMonthYear).Build()
	env2 := envs.NewBuilder().WithDateFormat(envs.DateFormatYearMonthDay).WithAllowedLanguages([]envs.Language{"spa"}).Build()

	assert.True(t, types.NewXDateTime(time.Date(2018, 4, 9, 17, 1, 30, 123456789, time.UTC)).Truthy())

	// test stringing
	assert.Equal(t, `2018-04-09T17:01:30.123456Z`, types.NewXDateTime(time.Date(2018, 4, 9, 17, 1, 30, 123456789, time.UTC)).Render())
	assert.Equal(t, `09-04-2018 17:01`, types.NewXDateTime(time.Date(2018, 4, 9, 17, 1, 30, 123456789, time.UTC)).Format(env))
	assert.Equal(t, `XDateTime(2018, 4, 9, 17, 1, 30, 123456789, UTC)`, types.NewXDateTime(time.Date(2018, 4, 9, 17, 1, 30, 123456789, time.UTC)).String())

	asJSON, _ := types.ToXJSON(types.NewXDateTime(time.Date(2018, 4, 9, 17, 1, 30, 123456789, time.UTC)))
	assert.Equal(t, types.NewXText(`"2018-04-09T17:01:30.123456Z"`), asJSON)

	// test equality
	assert.True(t, types.NewXDateTime(time.Date(2018, 4, 9, 17, 1, 30, 0, time.UTC)).Equals(types.NewXDateTime(time.Date(2018, 4, 9, 17, 1, 30, 0, time.UTC))))
	assert.False(t, types.NewXDateTime(time.Date(2019, 4, 9, 17, 1, 30, 0, time.UTC)).Equals(types.NewXDateTime(time.Date(2018, 4, 9, 17, 1, 30, 0, time.UTC))))

	// test comparisons
	assert.Equal(t, 0, types.NewXDateTime(time.Date(2018, 4, 9, 17, 1, 30, 0, time.UTC)).Compare(types.NewXDateTime(time.Date(2018, 4, 9, 17, 1, 30, 0, time.UTC))))
	assert.Equal(t, 1, types.NewXDateTime(time.Date(2018, 4, 9, 17, 1, 31, 0, time.UTC)).Compare(types.NewXDateTime(time.Date(2018, 4, 9, 17, 1, 30, 0, time.UTC))))
	assert.Equal(t, -1, types.NewXDateTime(time.Date(2018, 4, 9, 17, 1, 29, 0, time.UTC)).Compare(types.NewXDateTime(time.Date(2018, 4, 9, 17, 1, 30, 0, time.UTC))))

	la, _ := time.LoadLocation("America/Los_Angeles")

	d1 := types.NewXDateTime(time.Date(2018, 4, 9, 17, 1, 30, 0, la))
	assert.Equal(t, `datetime`, d1.Describe())

	formatted, err := d1.FormatCustom(env, "EEE, DD-MM-YYYY", nil)
	assert.NoError(t, err)
	assert.Equal(t, "Mon, 09-04-2018", formatted)

	formatted, err = d1.FormatCustom(env2, "EEE, DD-MM-YYYY", nil)
	assert.NoError(t, err)
	assert.Equal(t, "lun, 09-04-2018", formatted)

	_, err = d1.FormatCustom(env, "YYYYYY", nil)
	assert.EqualError(t, err, "'YYYYYY' is not valid in a datetime formatting layout")

	d2 := d1.ReplaceTime(types.NewXTime(dates.NewTimeOfDay(16, 20, 30, 123456789)))
	assert.Equal(t, 2018, d2.Native().Year())
	assert.Equal(t, time.Month(4), d2.Native().Month())
	assert.Equal(t, 9, d2.Native().Day())
	assert.Equal(t, 16, d2.Native().Hour())
	assert.Equal(t, 20, d2.Native().Minute())
	assert.Equal(t, 30, d2.Native().Second())
	assert.Equal(t, 123456789, d2.Native().Nanosecond())
	assert.Equal(t, la, d2.Native().Location())

	// test unmarshaling
	var date types.XDateTime
	err = jsonx.Unmarshal([]byte(`"2018-04-09T17:01:30Z"`), &date)
	assert.NoError(t, err)
	assert.Equal(t, types.NewXDateTime(time.Date(2018, 4, 9, 17, 1, 30, 0, time.UTC)), date)

	// test marshaling
	data, err := jsonx.Marshal(types.NewXDateTime(time.Date(2018, 4, 9, 17, 1, 30, 0, time.UTC)))
	assert.NoError(t, err)
	assert.Equal(t, `"2018-04-09T17:01:30.000000Z"`, string(data))
}

func TestToXDateTime(t *testing.T) {
	var tests = []struct {
		value    types.XValue
		expected types.XDateTime
		hasError bool
	}{
		{nil, types.XDateTimeZero, true},
		{types.NewXError(errors.Errorf("Error")), types.XDateTimeZero, true},
		{types.NewXNumberFromInt(123), types.XDateTimeZero, true},
		{types.NewXText("2018-06-05"), types.NewXDateTime(time.Date(2018, 6, 5, 0, 0, 0, 0, time.UTC)), false},
		{types.NewXText("wha?"), types.XDateTimeZero, true},
		{types.NewXDate(dates.NewDate(2018, 4, 9)), types.NewXDateTime(time.Date(2018, 4, 9, 0, 0, 0, 0, time.UTC)), false},
		{types.NewXDateTime(time.Date(2018, 4, 9, 17, 1, 30, 0, time.UTC)), types.NewXDateTime(time.Date(2018, 4, 9, 17, 1, 30, 0, time.UTC)), false},
		{types.NewXObject(map[string]types.XValue{
			"__default__": types.NewXText("2018-06-05"), // should use default
			"foo":         types.NewXNumberFromInt(234),
		}), types.NewXDateTime(time.Date(2018, 6, 5, 0, 0, 0, 0, time.UTC)), false},
	}

	env := envs.NewBuilder().Build()

	for _, test := range tests {
		result, err := types.ToXDateTime(env, test.value)

		if test.hasError {
			assert.Error(t, err, "expected error for input %T{%s}", test.value, test.value)
		} else {
			assert.NoError(t, err, "unexpected error for input %T{%s}", test.value, test.value)
			assert.Equal(t, test.expected.Native(), result.Native(), "result mismatch for input %T{%s}", test.value, test.value)
		}
	}
}

func TestToXDateTimeWithTimeFill(t *testing.T) {
	dates.SetNowSource(dates.NewFixedNowSource(time.Date(2018, 9, 13, 13, 36, 30, 123456789, time.UTC)))
	defer dates.SetNowSource(dates.DefaultNowSource)

	env := envs.NewBuilder().Build()
	result, err := types.ToXDateTimeWithTimeFill(env, types.NewXText("2018/12/20"))
	assert.NoError(t, err)
	assert.Equal(t, types.NewXDateTime(time.Date(2018, 12, 20, 13, 36, 30, 123456789, time.UTC)), result)
}
