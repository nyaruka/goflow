package types_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/nyaruka/goflow/excellent/types"

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
		asBool, _ := types.ToXBoolean(test.value)

		assert.Equal(t, test.asInternalJSON, string(asInternalJSON), "json.Marshal failed for %T{%s}", test.value, test.value)
		assert.Equal(t, types.NewXText(test.asJSON), asJSON, "ToXJSON failed for %T{%s}", test.value, test.value)
		assert.Equal(t, types.NewXText(test.asText), asText, "ToXText failed for %T{%s}", test.value, test.value)
		assert.Equal(t, types.NewXBoolean(test.asBool), asBool, "ToXBool failed for %T{%s}", test.value, test.value)
	}
}

func TestCompare(t *testing.T) {
	var tests = []struct {
		x1       types.XValue
		x2       types.XValue
		result   int
		hasError bool
	}{
		{nil, nil, 0, false},
		{nil, types.NewXText(""), 0, true},
		{types.NewXError(fmt.Errorf("Error")), types.NewXError(fmt.Errorf("Error")), 0, false},
		{types.NewXError(fmt.Errorf("Error")), types.XDateTimeZero, 0, true}, // type mismatch
		{types.NewXText("bob"), types.NewXText("bob"), 0, false},
		{types.NewXText("bob"), types.NewXText("cat"), -1, false},
		{types.NewXText("bob"), types.NewXText("ann"), 1, false},
		{types.NewXNumberFromInt(123), types.NewXNumberFromInt(123), 0, false},
		{types.NewXNumberFromInt(123), types.NewXNumberFromInt(124), -1, false},
		{types.NewXNumberFromInt(123), types.NewXNumberFromInt(122), 1, false},
	}

	for _, test := range tests {
		result, err := types.Compare(test.x1, test.x2)

		if test.hasError {
			assert.Error(t, err, "expected error for inputs '%s' and '%s'", test.x1, test.x2)
		} else {
			assert.NoError(t, err, "unexpected error for inputs '%s' and '%s'", test.x1, test.x2)
			assert.Equal(t, test.result, result, "result mismatch for inputs '%s' and '%s'", test.x1, test.x2)
		}
	}
}
