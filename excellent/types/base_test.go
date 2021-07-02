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
	"github.com/stretchr/testify/require"
)

func TestXValue(t *testing.T) {
	chi, err := time.LoadLocation("America/Chicago")
	require.NoError(t, err)

	date1 := time.Date(2017, 6, 23, 15, 30, 0, 0, time.UTC)
	date2 := time.Date(2017, 7, 18, 15, 30, 0, 0, chi)
	object1 := types.NewXObject(map[string]types.XValue{
		"foo": types.NewXText("Hello"),
		"bar": types.NewXNumberFromInt(123),
	})
	object2 := types.NewXObject(map[string]types.XValue{
		"foo": types.NewXText("World"),
		"bar": types.NewXNumberFromInt(456),
	})

	env := envs.NewBuilder().WithDateFormat(envs.DateFormatDayMonthYear).Build()

	tests := []struct {
		value     types.XValue
		marshaled string
		rendered  string
		formatted string
		asBool    bool
	}{
		{
			value:     nil,
			marshaled: `null`,
			rendered:  "",
			formatted: "",
			asBool:    false,
		}, {
			value:     types.NewXText(""),
			marshaled: `""`,
			rendered:  "",
			formatted: "",
			asBool:    false, // empty strings are false
		}, {
			value:     types.NewXText("FALSE"),
			marshaled: `"FALSE"`,
			rendered:  "FALSE",
			formatted: "FALSE",
			asBool:    false, // because it's string value is "false"
		}, {
			value:     types.NewXText("hello \"bob\""),
			marshaled: `"hello \"bob\""`,
			rendered:  "hello \"bob\"",
			formatted: "hello \"bob\"",
			asBool:    true,
		}, {
			value:     types.NewXNumberFromInt(0),
			marshaled: `0`,
			rendered:  "0",
			formatted: "0",
			asBool:    false, // because any decimal != 0 is true
		}, {
			value:     types.NewXNumberFromInt(1234),
			marshaled: `1234`,
			rendered:  "1234",
			formatted: "1,234",
			asBool:    true, // because any decimal != 0 is true
		}, {
			value:     types.RequireXNumberFromString("123.00"),
			marshaled: `123`,
			rendered:  "123",
			formatted: "123",
			asBool:    true,
		}, {
			value:     types.RequireXNumberFromString("1234.5678"),
			marshaled: `1234.5678`,
			rendered:  "1234.5678",
			formatted: "1,234.5678",
			asBool:    true,
		}, {
			value:     types.NewXBoolean(false),
			marshaled: `false`,
			rendered:  "false",
			formatted: "false",
			asBool:    false,
		}, {
			value:     types.NewXBoolean(true),
			marshaled: `true`,
			rendered:  "true",
			formatted: "true",
			asBool:    true,
		}, {
			value:     types.NewXDateTime(date1),
			marshaled: `"2017-06-23T15:30:00.000000Z"`,
			rendered:  "2017-06-23T15:30:00.000000Z",
			formatted: "23-06-2017 15:30",
			asBool:    true,
		}, {
			value:     types.NewXDateTime(date2),
			marshaled: `"2017-07-18T15:30:00.000000-05:00"`,
			rendered:  "2017-07-18T15:30:00.000000-05:00",
			formatted: "18-07-2017 20:30",
			asBool:    true,
		}, {
			value:     types.NewXArray(),
			marshaled: `[]`,
			rendered:  `[]`,
			formatted: "",
			asBool:    false,
		}, {
			value:     types.NewXArray(types.NewXNumberFromInt(1), types.NewXNumberFromInt(2)),
			marshaled: `[1,2]`,
			rendered:  `[1, 2]`,
			formatted: "1, 2",
			asBool:    true,
		}, {
			value:     types.NewXArray(types.NewXDateTime(date1), types.NewXDateTime(date2)),
			marshaled: `["2017-06-23T15:30:00.000000Z","2017-07-18T15:30:00.000000-05:00"]`,
			rendered:  `[2017-06-23T15:30:00.000000Z, 2017-07-18T15:30:00.000000-05:00]`,
			formatted: "23-06-2017 15:30, 18-07-2017 20:30",
			asBool:    true,
		}, {
			value:     types.NewXArray(object1, object2),
			marshaled: `[{"bar":123,"foo":"Hello"},{"bar":456,"foo":"World"}]`,
			rendered:  `[{bar: 123, foo: Hello}, {bar: 456, foo: World}]`,
			formatted: "- bar: 123\n  foo: Hello\n- bar: 456\n  foo: World",
			asBool:    true,
		}, {
			value:     types.XObjectEmpty,
			marshaled: `{}`,
			rendered:  `{}`,
			formatted: "",
			asBool:    false,
		}, {
			value:     types.NewXObject(map[string]types.XValue{"first": object1, "second": object2}),
			marshaled: `{"first":{"bar":123,"foo":"Hello"},"second":{"bar":456,"foo":"World"}}`,
			rendered:  `{first: {bar: 123, foo: Hello}, second: {bar: 456, foo: World}}`,
			formatted: "first:\n  bar: 123\n  foo: Hello\nsecond:\n  bar: 456\n  foo: World",
			asBool:    true,
		}, {
			value:     types.NewXObject(map[string]types.XValue{"__default__": types.NewXNumberFromInt(1), "foo": object1}),
			marshaled: `{"foo":{"bar":123,"foo":"Hello"}}`,
			rendered:  `1`,
			formatted: "1",
			asBool:    true,
		}, {
			value:     types.NewXError(errors.Errorf("it failed")), // once an error, always an error
			marshaled: `null`,
			rendered:  "",
			formatted: "",
			asBool:    false,
		},
	}
	for _, test := range tests {
		marshaled := jsonx.MustMarshal(test.value)
		rendered, _ := types.ToXText(env, test.value)
		formatted := types.Format(env, test.value)
		asBool, _ := types.ToXBoolean(test.value)

		assert.Equal(t, test.marshaled, string(marshaled), "jsonx.Marshal mismatch for %T{%s}", test.value, test.value)
		assert.Equal(t, types.NewXText(test.rendered), rendered, "ToXText mismatch for %T{%s}", test.value, test.value)
		assert.Equal(t, test.formatted, formatted, "Format mismatch for %T{%s}", test.value, test.value)
		assert.Equal(t, types.NewXBoolean(test.asBool), asBool, "ToXBool mismatch for %T{%s}", test.value, test.value)

		if types.IsXError(test.value) {
			_, xerr := types.ToXJSON(test.value)
			assert.Error(t, xerr)
		} else {
			marshaled, _ := types.ToXJSON(test.value)
			assert.Equal(t, types.NewXText(test.marshaled), marshaled, "ToXJSON mismatch for %T{%s}", test.value, test.value)
		}
	}
}

func TestEquals(t *testing.T) {
	var tests = []struct {
		x1     types.XValue
		x2     types.XValue
		result bool
	}{

		{nil, nil, true},                                         // nil == nil
		{nil, types.NewXText(""), false},                         // nil != non-nil
		{types.NewXText("1"), types.NewXNumberFromInt(1), false}, // different types are never equal

		{types.NewXArray(types.XBooleanFalse, types.NewXText("bob")), types.NewXArray(types.XBooleanFalse, types.NewXText("bob")), true},
		{types.NewXArray(types.XBooleanFalse, types.NewXText("abc")), types.NewXArray(types.XBooleanFalse, types.NewXText("bob")), false},
		{types.NewXArray(types.XBooleanFalse, types.NewXText("bob")), types.NewXArray(types.NewXText("bob"), types.XBooleanFalse), false}, // order matters

		{types.XBooleanFalse, types.XBooleanFalse, true},
		{types.XBooleanTrue, types.XBooleanFalse, false},

		{types.NewXDate(dates.NewDate(2018, 4, 9)), types.NewXDate(dates.NewDate(2018, 4, 9)), true},
		{types.NewXDate(dates.NewDate(2019, 4, 9)), types.NewXDate(dates.NewDate(2018, 4, 10)), false},

		{types.NewXDateTime(time.Date(2018, 4, 9, 17, 1, 30, 0, time.UTC)), types.NewXDateTime(time.Date(2018, 4, 9, 17, 1, 30, 0, time.UTC)), true},
		{types.NewXDateTime(time.Date(2019, 4, 9, 17, 1, 30, 0, time.UTC)), types.NewXDateTime(time.Date(2018, 4, 9, 17, 1, 30, 0, time.UTC)), false},

		{
			types.NewXObject(map[string]types.XValue{"foo": types.XBooleanFalse, "bar": types.NewXText("bob")}),
			types.NewXObject(map[string]types.XValue{"foo": types.XBooleanFalse, "bar": types.NewXText("bob")}),
			true,
		},
		{
			types.NewXObject(map[string]types.XValue{"__default__": types.XBooleanTrue, "bar": types.NewXText("bob")}),
			types.NewXObject(map[string]types.XValue{"__default__": types.XBooleanFalse, "bar": types.NewXText("bob")}),
			false, // different default
		},
		{
			types.NewXObject(map[string]types.XValue{"foo": types.XBooleanFalse, "bar": types.NewXText("bob")}),
			types.NewXObject(map[string]types.XValue{"foo": types.XBooleanFalse}),
			false, // different number of keys
		},
		{
			types.NewXObject(map[string]types.XValue{"foo": types.XBooleanFalse, "bar": types.NewXText("bob")}),
			types.NewXObject(map[string]types.XValue{"foo": types.XBooleanFalse, "baz": types.NewXText("bob")}),
			false, // different key
		},
		{
			types.NewXObject(map[string]types.XValue{"foo": types.XBooleanFalse, "bar": types.NewXText("bob")}),
			types.NewXObject(map[string]types.XValue{"foo": types.XBooleanFalse, "bar": types.NewXText("boo")}),
			false, // different value
		},

		{types.NewXError(errors.Errorf("Error")), types.NewXError(errors.Errorf("Error")), true},
		{types.NewXError(errors.Errorf("Error")), types.XDateTimeZero, false},

		{types.NewXText("bob"), types.NewXText("bob"), true},
		{types.NewXText("bob"), types.NewXText("abc"), false},

		{types.NewXTime(dates.NewTimeOfDay(10, 30, 0, 123456789)), types.NewXTime(dates.NewTimeOfDay(10, 30, 0, 123456789)), true},
		{types.NewXTime(dates.NewTimeOfDay(10, 30, 0, 123456789)), types.NewXTime(dates.NewTimeOfDay(10, 30, 0, 987654321)), false},

		{types.NewXNumberFromInt(123), types.NewXNumberFromInt(123), true},
		{types.NewXNumberFromInt(123), types.NewXNumberFromInt(124), false},
	}

	for _, tc := range tests {
		assert.Equal(t, tc.result, types.Equals(tc.x1, tc.x2), "equality mismatch for inputs '%s' and '%s'", tc.x1, tc.x2)
	}

	// test we get panic if we forgot to code Equals for a new xvalue type
	assert.Panics(t, func() { types.Equals(&XBogusType{}, &XBogusType{}) })
}

type XBogusType struct {
	types.XText
}
