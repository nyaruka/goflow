package types_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testXObject struct {
	foo string
	bar int
}

func NewTestXObject(foo string, bar int) *testXObject {
	return &testXObject{foo: foo, bar: bar}
}

// ToXJSON is called when this type is passed to @(json(...))
func (v *testXObject) ToXJSON() types.XText {
	return types.NewXMap(map[string]types.XValue{
		"foo": types.NewXText(v.foo),
		"bar": types.NewXNumberFromInt(v.bar),
	}).ToXJSON()
}

// MarshalJSON converts this type to its internal JSON representation which can differ from ToJSON
func (v *testXObject) MarshalJSON() ([]byte, error) {
	e := struct {
		Foo string `json:"foo"`
	}{
		Foo: v.foo,
	}
	return json.Marshal(e)
}

func (v *testXObject) Reduce() types.XPrimitive { return types.NewXText(v.foo) }

var _ types.XValue = &testXObject{}

func TestXValueRequiredConversions(t *testing.T) {
	chi, err := time.LoadLocation("America/Chicago")
	require.NoError(t, err)

	date1 := time.Date(2017, 6, 23, 15, 30, 0, 0, time.UTC)
	date2 := time.Date(2017, 7, 18, 15, 30, 0, 0, chi)

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
			asText:         `["1","2"]`,
			asBool:         true,
			isEmpty:        false,
		}, {
			value:          types.NewXArray(types.NewXDateTime(date1), types.NewXDateTime(date2)),
			asInternalJSON: `["2017-06-23T15:30:00Z","2017-07-18T15:30:00-05:00"]`,
			asJSON:         `["2017-06-23T15:30:00.000000Z","2017-07-18T15:30:00.000000-05:00"]`,
			asText:         `["2017-06-23T15:30:00.000000Z","2017-07-18T15:30:00.000000-05:00"]`,
			asBool:         true,
			isEmpty:        false,
		}, {
			value:          NewTestXObject("Hello", 123),
			asInternalJSON: `{"foo":"Hello"}`,
			asJSON:         `{"bar":123,"foo":"Hello"}`,
			asText:         "Hello",
			asBool:         true,
			isEmpty:        false,
		}, {
			value:          NewTestXObject("", 123),
			asInternalJSON: `{"foo":""}`,
			asJSON:         `{"bar":123,"foo":""}`,
			asText:         "",
			asBool:         false, // because it reduces to a string which itself is false
			isEmpty:        false,
		}, {
			value:          types.NewXArray(NewTestXObject("Hello", 123), NewTestXObject("World", 456)),
			asInternalJSON: `[{"foo":"Hello"},{"foo":"World"}]`,
			asJSON:         `[{"bar":123,"foo":"Hello"},{"bar":456,"foo":"World"}]`,
			asText:         `["Hello","World"]`,
			asBool:         true,
			isEmpty:        false,
		}, {
			value:          types.NewXEmptyMap(),
			asInternalJSON: `{}`,
			asJSON:         `{}`,
			asText:         `{}`,
			asBool:         false,
			isEmpty:        true,
		}, {
			value: types.NewXMap(map[string]types.XValue{
				"first":  NewTestXObject("Hello", 123),
				"second": NewTestXObject("World", 456),
			}),
			asInternalJSON: `{"first":{"foo":"Hello"},"second":{"foo":"World"}}`,
			asJSON:         `{"first":{"bar":123,"foo":"Hello"},"second":{"bar":456,"foo":"World"}}`,
			asText:         `{"first":"Hello","second":"World"}`,
			asBool:         true,
			isEmpty:        false,
		}, {
			value:          types.NewXJSONArray([]byte(`[]`)),
			asInternalJSON: `[]`,
			asJSON:         `[]`,
			asText:         `[]`,
			asBool:         false,
			isEmpty:        false,
		}, {
			value:          types.NewXJSONArray([]byte(`[5,     "x"]`)),
			asInternalJSON: `[5,"x"]`,
			asJSON:         `[5,     "x"]`,
			asText:         `[5,     "x"]`,
			asBool:         true,
			isEmpty:        false,
		}, {
			value:          types.NewXJSONObject([]byte(`{}`)),
			asInternalJSON: `{}`,
			asJSON:         `{}`,
			asText:         `{}`,
			asBool:         false,
			isEmpty:        false,
		}, {
			value:          types.NewXJSONObject([]byte(`{"foo":"World","bar":456}`)),
			asInternalJSON: `{"foo":"World","bar":456}`,
			asJSON:         `{"foo":"World","bar":456}`,
			asText:         `{"foo":"World","bar":456}`,
			asBool:         true,
			isEmpty:        false,
		}, {
			value:          types.NewXError(fmt.Errorf("it failed")), // once an error, always an error
			asInternalJSON: "",
			asJSON:         "",
			asText:         "",
			asBool:         false,
			isEmpty:        false,
		},
	}
	for _, test := range tests {
		asInternalJSON, _ := json.Marshal(test.value)
		asJSON, _ := types.ToXJSON(test.value)
		asText, _ := types.ToXText(test.value)
		asBool, _ := types.ToXBool(test.value)

		assert.Equal(t, test.asInternalJSON, string(asInternalJSON), "json.Marshal failed for %T{%s}", test.value, test.value)
		assert.Equal(t, types.NewXText(test.asJSON), asJSON, "ToXJSON failed for %T{%s}", test.value, test.value)
		assert.Equal(t, types.NewXText(test.asText), asText, "ToXText failed for %T{%s}", test.value, test.value)
		assert.Equal(t, types.NewXBoolean(test.asBool), asBool, "ToXBool failed for %T{%s}", test.value, test.value)
	}
}

func TestToXNumber(t *testing.T) {
	var tests = []struct {
		value    types.XValue
		asNumber types.XNumber
		hasError bool
	}{
		{nil, types.XNumberZero, false},
		{types.NewXErrorf("Error"), types.XNumberZero, true},
		{types.NewXNumberFromInt(123), types.NewXNumberFromInt(123), false},
		{types.NewXText("15.5"), types.RequireXNumberFromString("15.5"), false},
		{types.NewXText("lO.5"), types.RequireXNumberFromString("10.5"), false},
		{NewTestXObject("Hello", 123), types.XNumberZero, true},
		{NewTestXObject("123.45000", 123), types.RequireXNumberFromString("123.45"), false},
	}

	for _, test := range tests {
		result, err := types.ToXNumber(test.value)

		if test.hasError {
			assert.Error(t, err, "expected error for input %T{%s}", test.value, test.value)
		} else {
			assert.NoError(t, err, "unexpected error for input %T{%s}", test.value, test.value)
			assert.Equal(t, test.asNumber.Native(), result.Native(), "result mismatch for input %T{%s}", test.value, test.value)
		}
	}
}

func TestToXDate(t *testing.T) {
	var tests = []struct {
		value    types.XValue
		asNumber types.XDateTime
		hasError bool
	}{
		{nil, types.XDateTimeZero, false},
		{types.NewXError(fmt.Errorf("Error")), types.XDateTimeZero, true},
		{types.NewXNumberFromInt(123), types.XDateTimeZero, true},
		{types.NewXText("2018-06-05"), types.NewXDateTime(time.Date(2018, 6, 5, 0, 0, 0, 0, time.UTC)), false},
		{types.NewXText("wha?"), types.XDateTimeZero, true},
		{NewTestXObject("Hello", 123), types.XDateTimeZero, true},
		{NewTestXObject("2018/6/5", 123), types.NewXDateTime(time.Date(2018, 6, 5, 0, 0, 0, 0, time.UTC)), false},
	}

	env := utils.NewDefaultEnvironment()

	for _, test := range tests {
		result, err := types.ToXDate(env, test.value)

		if test.hasError {
			assert.Error(t, err, "expected error for input %T{%s}", test.value, test.value)
		} else {
			assert.NoError(t, err, "unexpected error for input %T{%s}", test.value, test.value)
			assert.Equal(t, test.asNumber.Native(), result.Native(), "result mismatch for input %T{%s}", test.value, test.value)
		}
	}
}
