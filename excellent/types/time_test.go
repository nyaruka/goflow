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

func TestXTime(t *testing.T) {
	env := envs.NewBuilder().Build()

	t1 := types.NewXTime(dates.NewTimeOfDay(17, 1, 30, 0))
	assert.Equal(t, `time`, t1.Describe())
	assert.True(t, t1.Truthy())
	assert.Equal(t, `17:01:30.000000`, types.NewXTime(dates.NewTimeOfDay(17, 1, 30, 0)).Render())
	assert.Equal(t, `17:01`, types.NewXTime(dates.NewTimeOfDay(17, 1, 30, 0)).Format(env))
	assert.Equal(t, `XTime(17, 1, 30, 0)`, types.NewXTime(dates.NewTimeOfDay(17, 1, 30, 0)).String())

	formatted, err := t1.FormatCustom(env, "ss")
	assert.NoError(t, err)
	assert.Equal(t, `30`, formatted)

	_, err = t1.FormatCustom(env, "ssssss")
	assert.EqualError(t, err, "'ssssss' is not valid in a time formatting layout")

	marshaled, err := jsonx.Marshal(t1)
	assert.NoError(t, err)
	assert.Equal(t, `"17:01:30.000000"`, string(marshaled))

	// test equality
	assert.True(t, t1.Equals(types.NewXTime(dates.NewTimeOfDay(17, 1, 30, 0))))
	assert.False(t, t1.Equals(types.NewXTime(dates.NewTimeOfDay(17, 1, 30, 1))))

	// test comparisons
	assert.Equal(t, 0, types.NewXTime(dates.NewTimeOfDay(17, 1, 30, 0)).Compare(t1))
	assert.Equal(t, 1, types.NewXTime(dates.NewTimeOfDay(17, 1, 31, 0)).Compare(t1))
	assert.Equal(t, -1, types.NewXTime(dates.NewTimeOfDay(17, 1, 29, 0)).Compare(t1))
}

func TestToXTime(t *testing.T) {
	var tests = []struct {
		value    types.XValue
		expected types.XTime
		hasError bool
	}{
		{nil, types.XTimeZero, true},
		{types.NewXError(errors.Errorf("Error")), types.XTimeZero, true},
		{types.NewXNumberFromInt(123), types.XTimeZero, true},
		{types.NewXNumberFromInt(23), types.NewXTime(dates.NewTimeOfDay(23, 0, 0, 0)), false},
		{types.NewXNumberFromInt(24), types.XTimeZero, false},
		{types.NewXText("10:30"), types.NewXTime(dates.NewTimeOfDay(10, 30, 0, 0)), false},
		{types.NewXText("10:30 pm"), types.NewXTime(dates.NewTimeOfDay(22, 30, 0, 0)), false},
		{types.NewXText("10"), types.NewXTime(dates.NewTimeOfDay(10, 0, 0, 0)), false},
		{types.NewXText("10 PM"), types.NewXTime(dates.NewTimeOfDay(22, 0, 0, 0)), false},
		{types.NewXText("wha?"), types.XTimeZero, true},
		{types.NewXTime(dates.NewTimeOfDay(17, 1, 30, 0)), types.NewXTime(dates.NewTimeOfDay(17, 1, 30, 0)), false},
		{types.NewXDateTime(time.Date(2018, 4, 9, 17, 1, 30, 0, time.UTC)), types.NewXTime(dates.NewTimeOfDay(17, 1, 30, 0)), false},
		{types.NewXObject(map[string]types.XValue{
			"__default__": types.NewXText("10:30"), // should use default
			"foo":         types.NewXNumberFromInt(234),
		}), types.NewXTime(dates.NewTimeOfDay(10, 30, 0, 0)), false},
	}

	env := envs.NewBuilder().Build()

	for _, test := range tests {
		result, err := types.ToXTime(env, test.value)

		if test.hasError {
			assert.Error(t, err, "expected error for input %T{%s}", test.value, test.value)
		} else {
			assert.NoError(t, err, "unexpected error for input %T{%s}", test.value, test.value)
			assert.Equal(t, test.expected.Native(), result.Native(), "result mismatch for input %T{%s}", test.value, test.value)
		}
	}
}
