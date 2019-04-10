package types_test

import (
	"testing"
	"time"

	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/utils"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestXValueRequiredConversions(t *testing.T) {
	chi, err := time.LoadLocation("America/Chicago")
	require.NoError(t, err)

	date1 := time.Date(2017, 6, 23, 15, 30, 0, 0, time.UTC)
	date2 := time.Date(2017, 7, 18, 15, 30, 0, 0, chi)
	dict1 := types.NewXDict(map[string]types.XValue{
		"foo": types.NewXText("Hello"),
		"bar": types.NewXNumberFromInt(123),
	})
	dict2 := types.NewXDict(map[string]types.XValue{
		"foo": types.NewXText("World"),
		"bar": types.NewXNumberFromInt(456),
	})

	env := utils.NewEnvironmentBuilder().Build()

	tests := []struct {
		value          types.XValue
		asJSON         string
		asInternalJSON string
		asText         string
		asBool         bool
		isEmpty        bool
	}{
		{
			value:          nil,
			asInternalJSON: `null`,
			asJSON:         `null`,
			asText:         "",
			asBool:         false,
			isEmpty:        true,
		}, {
			value:          types.NewXText(""),
			asInternalJSON: `""`,
			asJSON:         `""`,
			asText:         "",
			asBool:         false, // empty strings are false
			isEmpty:        true,
		}, {
			value:          types.NewXText("FALSE"),
			asInternalJSON: `"FALSE"`,
			asJSON:         `"FALSE"`,
			asText:         "FALSE",
			asBool:         false, // because it's string value is "false"
			isEmpty:        false,
		}, {
			value:          types.NewXText("hello \"bob\""),
			asInternalJSON: `"hello \"bob\""`,
			asJSON:         `"hello \"bob\""`,
			asText:         "hello \"bob\"",
			asBool:         true,
			isEmpty:        false,
		}, {
			value:          types.NewXNumberFromInt(0),
			asInternalJSON: `0`,
			asJSON:         `0`,
			asText:         "0",
			asBool:         false, // because any decimal != 0 is true
			isEmpty:        false,
		}, {
			value:          types.NewXNumberFromInt(123),
			asInternalJSON: `123`,
			asJSON:         `123`,
			asText:         "123",
			asBool:         true, // because any decimal != 0 is true
			isEmpty:        false,
		}, {
			value:          types.RequireXNumberFromString("123.00"),
			asInternalJSON: `123`,
			asJSON:         `123`,
			asText:         "123",
			asBool:         true,
			isEmpty:        false,
		}, {
			value:          types.RequireXNumberFromString("123.45"),
			asInternalJSON: `123.45`,
			asJSON:         `123.45`,
			asText:         "123.45",
			asBool:         true,
			isEmpty:        false,
		}, {
			value:          types.NewXBoolean(false),
			asInternalJSON: `false`,
			asJSON:         `false`,
			asText:         "false",
			asBool:         false,
			isEmpty:        false,
		}, {
			value:          types.NewXBoolean(true),
			asInternalJSON: `true`,
			asJSON:         `true`,
			asText:         "true",
			asBool:         true,
			isEmpty:        false,
		}, {
			value:          types.NewXDateTime(date1),
			asInternalJSON: `"2017-06-23T15:30:00Z"`,
			asJSON:         `"2017-06-23T15:30:00.000000Z"`,
			asText:         "2017-06-23T15:30:00.000000Z",
			asBool:         true,
			isEmpty:        false,
		}, {
			value:          types.NewXDateTime(date2),
			asInternalJSON: `"2017-07-18T15:30:00-05:00"`,
			asJSON:         `"2017-07-18T15:30:00.000000-05:00"`,
			asText:         "2017-07-18T15:30:00.000000-05:00",
			asBool:         true,
			isEmpty:        false,
		}, {
			value:          types.NewXArray(),
			asInternalJSON: `[]`,
			asJSON:         `[]`,
			asText:         `[]`,
			asBool:         false,
			isEmpty:        true,
		}, {
			value:          types.NewXArray(types.NewXNumberFromInt(1), types.NewXNumberFromInt(2)),
			asInternalJSON: `[1,2]`,
			asJSON:         `[1,2]`,
			asText:         `[1, 2]`,
			asBool:         true,
			isEmpty:        false,
		}, {
			value:          types.NewXArray(types.NewXDateTime(date1), types.NewXDateTime(date2)),
			asInternalJSON: `["2017-06-23T15:30:00Z","2017-07-18T15:30:00-05:00"]`,
			asJSON:         `["2017-06-23T15:30:00.000000Z","2017-07-18T15:30:00.000000-05:00"]`,
			asText:         `[2017-06-23T15:30:00.000000Z, 2017-07-18T15:30:00.000000-05:00]`,
			asBool:         true,
			isEmpty:        false,
		}, {
			value:          types.NewXArray(dict1, dict2),
			asInternalJSON: `[{"bar":123,"foo":"Hello"},{"bar":456,"foo":"World"}]`,
			asJSON:         `[{"bar":123,"foo":"Hello"},{"bar":456,"foo":"World"}]`,
			asText:         `[{bar: 123, foo: Hello}, {bar: 456, foo: World}]`,
			asBool:         true,
			isEmpty:        false,
		}, {
			value:          types.XDictEmpty,
			asInternalJSON: `{}`,
			asJSON:         `{}`,
			asText:         `{}`,
			asBool:         false,
			isEmpty:        true,
		}, {
			value:          types.NewXDict(map[string]types.XValue{"first": dict1, "second": dict2}),
			asInternalJSON: `{"first":{"bar":123,"foo":"Hello"},"second":{"bar":456,"foo":"World"}}`,
			asJSON:         `{"first":{"bar":123,"foo":"Hello"},"second":{"bar":456,"foo":"World"}}`,
			asText:         `{first: {bar: 123, foo: Hello}, second: {bar: 456, foo: World}}`,
			asBool:         true,
			isEmpty:        false,
		}, {
			value:          types.NewXError(errors.Errorf("it failed")), // once an error, always an error
			asInternalJSON: "",
			asJSON:         "",
			asText:         "",
			asBool:         false,
			isEmpty:        false,
		},
	}
	for _, test := range tests {
		asInternalJSON, _ := utils.JSONMarshal(test.value)
		asJSON, _ := types.ToXJSON(env, test.value)
		asText, _ := types.ToXText(env, test.value)
		asBool, _ := types.ToXBoolean(env, test.value)

		assert.Equal(t, test.asInternalJSON, string(asInternalJSON), "json.Marshal failed for %T{%s}", test.value, test.value)
		assert.Equal(t, types.NewXText(test.asJSON), asJSON, "ToXJSON failed for %T{%s}", test.value, test.value)
		assert.Equal(t, types.NewXText(test.asText), asText, "ToXText failed for %T{%s}", test.value, test.value)
		assert.Equal(t, types.NewXBoolean(test.asBool), asBool, "ToXBool failed for %T{%s}", test.value, test.value)
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

		{types.NewXDate(utils.NewDate(2018, 4, 9)), types.NewXDate(utils.NewDate(2018, 4, 9)), true},
		{types.NewXDate(utils.NewDate(2019, 4, 9)), types.NewXDate(utils.NewDate(2018, 4, 10)), false},

		{types.NewXDateTime(time.Date(2018, 4, 9, 17, 1, 30, 0, time.UTC)), types.NewXDateTime(time.Date(2018, 4, 9, 17, 1, 30, 0, time.UTC)), true},
		{types.NewXDateTime(time.Date(2019, 4, 9, 17, 1, 30, 0, time.UTC)), types.NewXDateTime(time.Date(2018, 4, 9, 17, 1, 30, 0, time.UTC)), false},

		{
			types.NewXDict(map[string]types.XValue{"foo": types.XBooleanFalse, "bar": types.NewXText("bob")}),
			types.NewXDict(map[string]types.XValue{"foo": types.XBooleanFalse, "bar": types.NewXText("bob")}),
			true,
		},
		{
			types.NewXDict(map[string]types.XValue{"foo": types.XBooleanFalse, "bar": types.NewXText("bob")}),
			types.NewXDict(map[string]types.XValue{"foo": types.XBooleanFalse}),
			false, // different number of keys
		},
		{
			types.NewXDict(map[string]types.XValue{"foo": types.XBooleanFalse, "bar": types.NewXText("bob")}),
			types.NewXDict(map[string]types.XValue{"foo": types.XBooleanFalse, "baz": types.NewXText("bob")}),
			false, // different key
		},
		{
			types.NewXDict(map[string]types.XValue{"foo": types.XBooleanFalse, "bar": types.NewXText("bob")}),
			types.NewXDict(map[string]types.XValue{"foo": types.XBooleanFalse, "bar": types.NewXText("boo")}),
			false, // different value
		},

		{types.NewXError(errors.Errorf("Error")), types.NewXError(errors.Errorf("Error")), true},
		{types.NewXError(errors.Errorf("Error")), types.XDateTimeZero, false},

		{types.NewXText("bob"), types.NewXText("bob"), true},
		{types.NewXText("bob"), types.NewXText("abc"), false},

		{types.NewXTime(utils.NewTimeOfDay(10, 30, 0, 123456789)), types.NewXTime(utils.NewTimeOfDay(10, 30, 0, 123456789)), true},
		{types.NewXTime(utils.NewTimeOfDay(10, 30, 0, 123456789)), types.NewXTime(utils.NewTimeOfDay(10, 30, 0, 987654321)), false},

		{types.NewXNumberFromInt(123), types.NewXNumberFromInt(123), true},
		{types.NewXNumberFromInt(123), types.NewXNumberFromInt(124), false},
	}

	for _, tc := range tests {
		assert.Equal(t, tc.result, types.Equals(tc.x1, tc.x2), "equality mismatch for inputs '%s' and '%s'", tc.x1, tc.x2)
	}
}

func TestIsEmpty(t *testing.T) {
	assert.True(t, types.IsEmpty(nil))
	assert.True(t, types.IsEmpty(types.NewXArray()))
	assert.True(t, types.IsEmpty(types.XDictEmpty))
	assert.True(t, types.IsEmpty(types.NewXText("")))
	assert.False(t, types.IsEmpty(types.NewXText("a")))
	assert.False(t, types.IsEmpty(types.XBooleanFalse))
	assert.False(t, types.IsEmpty(types.XBooleanTrue))
	assert.False(t, types.IsEmpty(types.NewXNumberFromInt(0)))
	assert.False(t, types.IsEmpty(types.NewXNumberFromInt(123)))
}
