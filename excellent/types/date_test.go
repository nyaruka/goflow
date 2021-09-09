package types_test

import (
	"testing"
	"time"

	"github.com/nyaruka/gocommon/dates"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestXDate(t *testing.T) {
	env := envs.NewBuilder().WithDateFormat(envs.DateFormatDayMonthYear).Build()
	env2 := envs.NewBuilder().WithDateFormat(envs.DateFormatYearMonthDay).WithAllowedLanguages([]envs.Language{"spa"}).Build()

	d1 := types.NewXDate(dates.NewDate(2019, 2, 20))
	assert.Equal(t, `date`, d1.Describe())
	assert.Equal(t, `2019-02-20`, d1.Render())
	assert.Equal(t, `20-02-2019`, d1.Format(env))

	assert.True(t, d1.Truthy())
	assert.Equal(t, `XDate(2019, 2, 20)`, d1.String())

	formatted, err := d1.FormatCustom(env, "EEE, DD-MM-YYYY")
	assert.NoError(t, err)
	assert.Equal(t, "Wed, 20-02-2019", formatted)

	formatted, err = d1.FormatCustom(env2, "EEE, DD-MM-YYYY")
	assert.NoError(t, err)
	assert.Equal(t, "mié, 20-02-2019", formatted)

	_, err = d1.FormatCustom(env, "YYYYYY")
	assert.EqualError(t, err, "'YYYYYY' is not valid in a date formatting layout")

	asJSON, _ := types.ToXJSON(d1)
	assert.Equal(t, types.NewXText(`"2019-02-20"`), asJSON)

	// test equality
	assert.True(t, d1.Equals(types.NewXDate(dates.NewDate(2019, 2, 20))))
	assert.False(t, d1.Equals(types.NewXDate(dates.NewDate(2019, 2, 21))))

	// test comparisons
	assert.Equal(t, 0, types.NewXDate(dates.NewDate(2019, 2, 20)).Compare(d1))
	assert.Equal(t, 1, types.NewXDate(dates.NewDate(2019, 2, 21)).Compare(d1))
	assert.Equal(t, -1, types.NewXDate(dates.NewDate(2019, 2, 19)).Compare(d1))
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
		{types.NewXText("2018-01-20"), types.NewXDate(dates.NewDate(2018, 1, 20)), false},
		{types.NewXDate(dates.NewDate(2018, 4, 19)), types.NewXDate(dates.NewDate(2018, 4, 19)), false},
		{types.NewXDateTime(time.Date(2018, 4, 9, 17, 1, 30, 0, time.UTC)), types.NewXDate(dates.NewDate(2018, 4, 9)), false},
		{types.NewXObject(map[string]types.XValue{
			"__default__": types.NewXText("2018-01-20"), // should use default
			"foo":         types.NewXNumberFromInt(234),
		}), types.NewXDate(dates.NewDate(2018, 1, 20)), false},
	}

	env := envs.NewBuilder().Build()

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
