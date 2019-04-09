package types_test

import (
	"testing"
	"time"

	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/utils"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestXDate(t *testing.T) {
	d1 := types.NewXDate(utils.NewDate(2019, 2, 20))
	assert.Equal(t, `date`, d1.Describe())
	assert.Equal(t, `2019-02-20`, d1.String())

	// test equality
	assert.True(t, d1.Equals(types.NewXDate(utils.NewDate(2019, 2, 20))))
	assert.False(t, d1.Equals(types.NewXDate(utils.NewDate(2019, 2, 21))))

	// test comparisons
	assert.Equal(t, 0, types.NewXDate(utils.NewDate(2019, 2, 20)).Compare(d1))
	assert.Equal(t, 1, types.NewXDate(utils.NewDate(2019, 2, 21)).Compare(d1))
	assert.Equal(t, -1, types.NewXDate(utils.NewDate(2019, 2, 19)).Compare(d1))
}

func TestToXDate(t *testing.T) {
	var tests = []struct {
		value    types.XValue
		expected types.XDate
		hasError bool
	}{
		{nil, types.XDateZero, true},
		{types.NewXError(errors.Errorf("Error")), types.XDateZero, true},
		{types.NewXNumberFromInt(123), types.XDateZero, true},
		{types.NewXText("2018-01-20"), types.NewXDate(utils.NewDate(2018, 1, 20)), false},
		{types.NewXDate(utils.NewDate(2018, 4, 19)), types.NewXDate(utils.NewDate(2018, 4, 19)), false},
		{types.NewXDateTime(time.Date(2018, 4, 9, 17, 1, 30, 0, time.UTC)), types.NewXDate(utils.NewDate(2018, 4, 9)), false},
	}

	env := utils.NewEnvironmentBuilder().Build()

	for _, test := range tests {
		result, err := types.ToXDate(env, test.value)

		if test.hasError {
			assert.Error(t, err, "expected error for input %T{%s}", test.value, test.value)
		} else {
			assert.NoError(t, err, "unexpected error for input %T{%s}", test.value, test.value)
			assert.Equal(t, test.expected.Native(), result.Native(), "result mismatch for input %T{%s}", test.value, test.value)
		}
	}
}
