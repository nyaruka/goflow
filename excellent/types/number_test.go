package types_test

import (
	"testing"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestXNumber(t *testing.T) {
	env := envs.NewBuilder().Build()

	// test creation
	assert.Equal(t, types.RequireXNumberFromString("123"), types.NewXNumberFromInt(123))
	assert.Equal(t, types.RequireXNumberFromString("123"), types.NewXNumberFromInt64(123))
	assert.Panics(t, func() { types.RequireXNumberFromString("xxx") })

	n, err := types.NewXNumber(decimal.RequireFromString(`123.45`))
	assert.NoError(t, err)
	assert.Equal(t, types.RequireXNumberFromString("123.45"), n)

	n, err = types.NewXNumber(decimal.RequireFromString(`1234567890123456789012345678901234567890`))
	assert.EqualError(t, err, "number not in valid range")

	v := types.NewXNumberOrError(decimal.RequireFromString(`12345.67890`))
	assert.Equal(t, types.RequireXNumberFromString(`12345.67890`), v)

	v = types.NewXNumberOrError(decimal.RequireFromString(`1234567890123456789012345678901234567890`))
	assert.Equal(t, types.NewXErrorf("number not in valid range"), v)

	// test equality
	assert.True(t, types.NewXNumberFromInt(123).Equals(types.NewXNumberFromInt(123)))
	assert.False(t, types.NewXNumberFromInt(123).Equals(types.NewXNumberFromInt(124)))

	// test comparison
	assert.Equal(t, 0, types.NewXNumberFromInt(123).Compare(types.NewXNumberFromInt(123)))
	assert.Equal(t, -1, types.NewXNumberFromInt(123).Compare(types.NewXNumberFromInt(124)))
	assert.Equal(t, 1, types.NewXNumberFromInt(124).Compare(types.NewXNumberFromInt(123)))

	assert.Equal(t, `123`, types.NewXNumberFromInt64(123).Render())
	assert.Equal(t, `123`, types.NewXNumberFromInt64(123).Format(env))
	assert.Equal(t, `XNumber(123)`, types.NewXNumberFromInt64(123).String())
	assert.Equal(t, `XNumber(123.45)`, types.RequireXNumberFromString("123.45").String())

	// unmarshal with quotes
	var num types.XNumber
	err = jsonx.Unmarshal([]byte(`"23.45"`), &num)
	assert.NoError(t, err)
	assert.Equal(t, types.RequireXNumberFromString("23.45"), num)

	// unmarshal without quotes
	err = jsonx.Unmarshal([]byte(`34.56`), &num)
	assert.NoError(t, err)
	assert.Equal(t, types.RequireXNumberFromString("34.56"), num)

	// marshal (doesn't use quotes)
	data, err := jsonx.Marshal(types.RequireXNumberFromString("23.45"))
	assert.NoError(t, err)
	assert.Equal(t, []byte(`23.45`), data)
}

func TestToXNumber(t *testing.T) {
	var tests = []struct {
		value    types.XValue
		asNumber types.XNumber
		err      string
	}{
		{nil, types.XNumberZero, "unable to convert null to a number"},
		{types.NewXErrorf("Error"), types.XNumberZero, "Error"},
		{types.NewXNumberFromInt(123), types.NewXNumberFromInt(123), ""},
		{types.NewXText("15.5"), types.RequireXNumberFromString("15.5"), ""},
		{types.NewXText("  15.4  "), types.RequireXNumberFromString("15.4"), ""},
		{types.NewXObject(map[string]types.XValue{
			"__default__": types.NewXNumberFromInt(123), // should use default
			"foo":         types.NewXNumberFromInt(234),
		}), types.NewXNumberFromInt(123), ""},
		{types.NewXText("12345678901234567890"), types.XNumberZero, `unable to convert "12345678901234567890" to a number`}, // out of range
		{types.NewXText("1E100"), types.XNumberZero, `unable to convert "1E100" to a number`},                               // scientific notation not allowed
		{types.NewXText("1e100"), types.XNumberZero, `unable to convert "1e100" to a number`},                               // scientific notation not allowed
	}

	env := envs.NewBuilder().Build()

	for _, tc := range tests {
		number, err := types.ToXNumber(env, tc.value)

		if tc.err != "" {
			assert.EqualError(t, err, tc.err, "error mismatch for input %T{%s}", tc.value, tc.value)
		} else {
			assert.NoError(t, err, "unexpected error for input %T{%s}", tc.value, tc.value)
			assert.Equal(t, tc.asNumber.Native(), number.Native(), "number mismatch for input %T{%s}", tc.value, tc.value)
		}
	}
}

func TestToInteger(t *testing.T) {
	var tests = []struct {
		value     types.XValue
		asInteger int
		err       string
	}{
		{nil, 0, "unable to convert null to a number"},
		{types.NewXErrorf("Error"), 0, "Error"},
		{types.NewXNumberFromInt(123), 123, ""},
		{types.NewXText("15.5"), 15, ""},
		{types.NewXText("  15.4  "), 15, ""},
		{types.NewXObject(map[string]types.XValue{
			"__default__": types.NewXNumberFromInt(123), // should use default
			"foo":         types.NewXNumberFromInt(234),
		}), 123, ""},
		{types.NewXText("12345678901234567890"), 0, `unable to convert "12345678901234567890" to a number`}, // out of range
		{types.NewXText("1E100"), 0, `unable to convert "1E100" to a number`},                               // scientific notation not allowed
		{types.NewXText("1e100"), 0, `unable to convert "1e100" to a number`},                               // scientific notation not allowed
	}

	env := envs.NewBuilder().Build()

	for _, tc := range tests {
		integer, err := types.ToInteger(env, tc.value)

		if tc.err != "" {
			assert.EqualError(t, err, tc.err, "error mismatch for input %T{%s}", tc.value, tc.value)
		} else {
			assert.NoError(t, err, "unexpected error for input %T{%s}", tc.value, tc.value)
			assert.Equal(t, tc.asInteger, integer, "integer mismatch for input %T{%s}", tc.value, tc.value)
		}
	}
}

func TestFormatCustom(t *testing.T) {
	fmtTests := []struct {
		input       types.XNumber
		format      *envs.NumberFormat
		places      int
		groupDigits bool
		expected    string
	}{
		// zero padding for extending decimal places
		{types.RequireXNumberFromString("1"), envs.DefaultNumberFormat, 2, true, "1.00"},
		{types.RequireXNumberFromString("12"), envs.DefaultNumberFormat, 2, true, "12.00"},
		{types.RequireXNumberFromString("123"), envs.DefaultNumberFormat, 2, true, "123.00"},
		{types.RequireXNumberFromString("1234"), envs.DefaultNumberFormat, 2, true, "1,234.00"},
		{types.RequireXNumberFromString("123456789"), envs.DefaultNumberFormat, 2, true, "123,456,789.00"},

		// rounding for truncating decimal places
		{types.RequireXNumberFromString("1.9876"), envs.DefaultNumberFormat, 2, true, "1.99"},
		{types.RequireXNumberFromString("12.9876"), envs.DefaultNumberFormat, 2, true, "12.99"},
		{types.RequireXNumberFromString("123.9876"), envs.DefaultNumberFormat, 2, true, "123.99"},
		{types.RequireXNumberFromString("1234.9876"), envs.DefaultNumberFormat, 2, true, "1,234.99"},

		// rounding for truncating decimal places
		{types.RequireXNumberFromString("1.1111"), envs.DefaultNumberFormat, 0, true, "1"},
		{types.RequireXNumberFromString("12.1111"), envs.DefaultNumberFormat, 0, true, "12"},
		{types.RequireXNumberFromString("123.1111"), envs.DefaultNumberFormat, 0, true, "123"},
		{types.RequireXNumberFromString("1234.1111"), envs.DefaultNumberFormat, 0, true, "1,234"},

		{types.RequireXNumberFromString("1.9876"), envs.DefaultNumberFormat, 0, true, "2"},
		{types.RequireXNumberFromString("12.9876"), envs.DefaultNumberFormat, 0, true, "13"},
		{types.RequireXNumberFromString("123.9876"), envs.DefaultNumberFormat, 0, true, "124"},
		{types.RequireXNumberFromString("1234.9876"), envs.DefaultNumberFormat, 0, true, "1,235"},

		// places -1 means keep significant decimals
		{types.RequireXNumberFromString("1234"), envs.DefaultNumberFormat, -1, true, "1,234"},
		{types.RequireXNumberFromString("1234.000"), envs.DefaultNumberFormat, -1, true, "1,234"},
		{types.RequireXNumberFromString("1234.500"), envs.DefaultNumberFormat, -1, true, "1,234.5"},

		// grouping is optional
		{types.RequireXNumberFromString("1234"), envs.DefaultNumberFormat, 0, false, "1234"},
		{types.RequireXNumberFromString("1234.567"), envs.DefaultNumberFormat, 2, false, "1234.57"},

		// custom number format
		{types.RequireXNumberFromString("1234.567"), &envs.NumberFormat{DecimalSymbol: ",", DigitGroupingSymbol: "."}, 2, true, "1.234,57"},
	}

	for _, tc := range fmtTests {
		val := tc.input.FormatCustom(tc.format, tc.places, tc.groupDigits)

		assert.Equal(t, tc.expected, val, "format decimal failed for input=%s, places=%d", tc.input, tc.places)
	}
}
