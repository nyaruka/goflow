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
	foo := &struct {
		Val *types.XNumber `json:"val"`
	}{}
	err := jsonx.Unmarshal([]byte(`{"val": "23.45"}`), foo)
	assert.NoError(t, err)
	assert.Equal(t, types.RequireXNumberFromString("23.45").Native(), foo.Val.Native())

	// unmarshal without quotes
	err = jsonx.Unmarshal([]byte(`{"val": 34.56}`), foo)
	assert.NoError(t, err)
	assert.Equal(t, types.RequireXNumberFromString("34.56").Native(), foo.Val.Native())

	// marshal (doesn't use quotes)
	data, err := jsonx.Marshal(types.RequireXNumberFromString("23.45"))
	assert.NoError(t, err)
	assert.Equal(t, []byte(`23.45`), data)

	// test NewXNumber range check
	assert.False(t, types.IsXError(types.NewXNumber(decimal.New(123, 0))))
	assert.True(t, types.IsXError(types.NewXNumber(decimal.New(1, 200))))             // magnitude too large
	assert.True(t, types.IsXError(types.NewXNumber(decimal.New(1, -200))))            // magnitude too small
	assert.False(t, types.IsXError(types.NewXNumber(decimal.New(1, 50))))             // within range
	assert.False(t, types.IsXError(types.NewXNumber(decimal.New(1, -50))))            // within range
	assert.False(t, types.IsXError(types.NewXNumber(decimal.RequireFromString("0")))) // zero always ok
}

func TestToXNumberAndInteger(t *testing.T) {
	var tests = []struct {
		value     types.XValue
		asNumber  *types.XNumber
		asInteger int
		hasError  bool
	}{
		{nil, types.XNumberZero, 0, true},
		{types.NewXErrorf("Error"), types.XNumberZero, 0, true},
		{types.NewXNumberFromInt(123), types.NewXNumberFromInt(123), 123, false},
		{types.NewXText("15.5"), types.RequireXNumberFromString("15.5"), 15, false},
		{types.NewXText("  15.4  "), types.RequireXNumberFromString("15.4"), 15, false},
		{types.NewXObject(map[string]types.XValue{
			"__default__": types.NewXNumberFromInt(123), // should use default
			"foo":         types.NewXNumberFromInt(234),
		}), types.NewXNumberFromInt(123), 123, false},
		{types.NewXText("12345678901234567890"), types.RequireXNumberFromString("12345678901234567890"), 0, true},                                 // out of int range
		{types.NewXText("123456789012345678901234567890123456"), types.RequireXNumberFromString("123456789012345678901234567890123456"), 0, true}, // 36 digits, ok as number but out of int range
		{types.NewXText("1234567890123456789012345678901234567"), types.XNumberZero, 0, true},                                                     // 37 digits, too many
		{types.NewXText("1E100"), types.XNumberZero, 0, true},                                                                                     // scientific notation not allowed
		{types.NewXText("1e100"), types.XNumberZero, 0, true},                                                                                     // scientific notation not allowed
		{types.NewXText("234."), types.XNumberZero, 0, true},
		{types.NewXText("+1800567890"), types.XNumberZero, 0, true},
	}

	env := envs.NewBuilder().Build()

	for _, test := range tests {
		number, _ := types.ToXNumber(env, test.value)
		integer, err := types.ToInteger(env, test.value)

		if test.hasError {
			assert.Error(t, err.Native(), "expected error for input %T{%s}", test.value, test.value)
		} else {
			assert.NoError(t, err.Native(), "unexpected error for input %T{%s}", test.value, test.value)
			assert.Equal(t, test.asNumber.Native(), number.Native(), "number mismatch for input %T{%s}", test.value, test.value)
			assert.Equal(t, test.asInteger, integer, "integer mismatch for input %T{%s}", test.value, test.value)
		}
	}
}

func TestCheckDecimalRange(t *testing.T) {
	tests := []struct {
		value   decimal.Decimal
		isError bool
	}{
		{decimal.Zero, false},
		{decimal.New(1, 0), false},
		{decimal.New(-1, 0), false},
		{decimal.New(123, 0), false},
		{decimal.RequireFromString("123456789012345678901234567890123456"), false},                  // 36 significant digits - ok
		{decimal.RequireFromString("1234567890123456789012345678901234567"), true},                 // 37 significant digits - too many
		{decimal.RequireFromString("-1234567890123456789012345678901234567"), true},                // negative 37 significant digits
		{decimal.RequireFromString("1234567895171680000000000000000000000000"), false},             // 40 digits but only 15 significant - ok
		{decimal.RequireFromString("12345678901234567890123456789012345670000000"), true},          // 37 significant digits with trailing zeros
		{decimal.RequireFromString("0.000000000000000000000000000000000001"), false},
		{decimal.New(1, 100), false}, // 1E100 - ok magnitude
		{decimal.New(1, 200), true},  // 1E200 - too large magnitude
		{decimal.New(1, -100), false},
		{decimal.New(1, -200), true}, // 1E-200 - too small magnitude
	}

	for _, tc := range tests {
		err := types.CheckDecimalRange(tc.value)
		if tc.isError {
			assert.Error(t, err, "expected error for %s", tc.value)
		} else {
			assert.NoError(t, err, "unexpected error for %s", tc.value)
		}
	}
}

func TestFormatCustom(t *testing.T) {
	fmtTests := []struct {
		input       *types.XNumber
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
