package types_test

import (
	"testing"
	"time"

	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/utils"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestXTime(t *testing.T) {
	t1 := types.NewXTime(utils.NewTimeOfDay(17, 1, 30, 0))
	assert.Equal(t, t1, t1.Reduce(utils.NewEnvironmentBuilder().Build()))
	assert.Equal(t, `time`, t1.Describe())
	assert.Equal(t, `17:01:30.000000`, types.NewXTime(utils.NewTimeOfDay(17, 1, 30, 0)).String())

	// test equality
	assert.True(t, t1.Equals(types.NewXTime(utils.NewTimeOfDay(17, 1, 30, 0))))
	assert.False(t, t1.Equals(types.NewXTime(utils.NewTimeOfDay(17, 1, 30, 1))))

	// test comparisons
	assert.Equal(t, 0, types.NewXTime(utils.NewTimeOfDay(17, 1, 30, 0)).Compare(t1))
	assert.Equal(t, 1, types.NewXTime(utils.NewTimeOfDay(17, 1, 31, 0)).Compare(t1))
	assert.Equal(t, -1, types.NewXTime(utils.NewTimeOfDay(17, 1, 29, 0)).Compare(t1))
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
		{types.NewXText("10:30"), types.NewXTime(utils.NewTimeOfDay(10, 30, 0, 0)), false},
		{types.NewXText("10:30 pm"), types.NewXTime(utils.NewTimeOfDay(22, 30, 0, 0)), false},
		{types.NewXText("10"), types.NewXTime(utils.NewTimeOfDay(10, 0, 0, 0)), false},
		{types.NewXText("10 PM"), types.NewXTime(utils.NewTimeOfDay(22, 0, 0, 0)), false},
		{types.NewXText("wha?"), types.XTimeZero, true},
		{NewTestXObject("Hello", 123), types.XTimeZero, true},
		{NewTestXObject("10:30:24", 123), types.NewXTime(utils.NewTimeOfDay(10, 30, 24, 0)), false},
		{types.NewXTime(utils.NewTimeOfDay(17, 1, 30, 0)), types.NewXTime(utils.NewTimeOfDay(17, 1, 30, 0)), false},
		{types.NewXDateTime(time.Date(2018, 4, 9, 17, 1, 30, 0, time.UTC)), types.NewXTime(utils.NewTimeOfDay(17, 1, 30, 0)), false},
	}

	env := utils.NewEnvironmentBuilder().Build()

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
