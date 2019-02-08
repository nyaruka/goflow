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

type testXObject struct {
	foo string
	bar int
}

func NewTestXObject(foo string, bar int) *testXObject {
	return &testXObject{foo: foo, bar: bar}
}

func (v *testXObject) Resolve(env utils.Environment, key string) types.XValue {
	switch key {
	case "foo":
		return types.NewXText(v.foo)
	case "bar":
		return types.NewXNumberFromInt(v.bar)
	}
	return types.NewXResolveError(v, key)
}

// ToXJSON is called when this type is passed to @(json(...))
func (v *testXObject) ToXJSON(env utils.Environment) types.XText {
	return types.ResolveKeys(env, v, "foo", "bar").ToXJSON(env)
}

// MarshalJSON converts this type to its internal JSON representation which can differ from ToJSON
func (v *testXObject) MarshalJSON() ([]byte, error) {
	e := struct {
		Foo string `json:"foo"`
	}{
		Foo: v.foo,
	}
	return utils.JSONMarshal(e)
}

// Describe returns a representation of this type for error messages
func (v *testXObject) Describe() string { return "test" }

func (v *testXObject) Reduce(env utils.Environment) types.XPrimitive { return types.NewXText(v.foo) }

var _ types.XValue = &testXObject{}
var _ types.XResolvable = &testXObject{}

func TestXValueRequiredConversions(t *testing.T) {
	chi, err := time.LoadLocation("America/Chicago")
	require.NoError(t, err)

	date1 := time.Date(2017, 6, 23, 15, 30, 0, 0, time.UTC)
	date2 := time.Date(2017, 7, 18, 15, 30, 0, 0, chi)

	env := utils.NewEnvironmentBuilder().Environment()

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
			asText:         `[1,2]`,
			asBool:         true,
			isEmpty:        false,
		}, {
			value:          types.NewXArray(types.NewXDateTime(date1), types.NewXDateTime(date2)),
			asInternalJSON: `["2017-06-23T15:30:00Z","2017-07-18T15:30:00-05:00"]`,
			asJSON:         `["2017-06-23T15:30:00.000000Z","2017-07-18T15:30:00.000000-05:00"]`,
			asText:         `["2017-06-23T15:30:00Z","2017-07-18T15:30:00-05:00"]`,
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
			value:          types.NewEmptyXMap(),
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
	env := utils.NewEnvironmentBuilder().Environment()

	var tests = []struct {
		x1     types.XValue
		x2     types.XValue
		result bool
	}{
		{nil, nil, true},
		{nil, types.NewXText(""), false},
		{types.NewXError(errors.Errorf("Error")), types.NewXError(errors.Errorf("Error")), true},
		{types.NewXError(errors.Errorf("Error")), types.XDateTimeZero, false},
		{types.NewXText("bob"), types.NewXText("bob"), true},
		{types.NewXText("bob"), types.NewXText("abc"), false},
		{types.XBooleanFalse, types.XBooleanFalse, true},
		{types.XBooleanTrue, types.XBooleanFalse, false},
		{types.NewXNumberFromInt(123), types.NewXNumberFromInt(123), true},
		{types.NewXNumberFromInt(123), types.NewXNumberFromInt(124), false},
		{types.NewXDateTime(time.Date(2018, 4, 9, 17, 1, 30, 0, time.UTC)), types.NewXDateTime(time.Date(2018, 4, 9, 17, 1, 30, 0, time.UTC)), true},
		{types.NewXDateTime(time.Date(2019, 4, 9, 17, 1, 30, 0, time.UTC)), types.NewXDateTime(time.Date(2018, 4, 9, 17, 1, 30, 0, time.UTC)), false},
		{NewTestXObject("Hello", 123), NewTestXObject("Hello", 123), true},
		{NewTestXObject("Hello", 456), NewTestXObject("Hello", 123), true},
	}

	for _, test := range tests {
		assert.Equal(t, test.result, types.Equals(env, test.x1, test.x2), "equality mismatch for inputs '%s' and '%s'", test.x1, test.x2)
	}
}

func TestIsEmpty(t *testing.T) {
	assert.True(t, types.IsEmpty(nil))
	assert.True(t, types.IsEmpty(types.NewXArray()))
	assert.True(t, types.IsEmpty(types.NewEmptyXMap()))
	assert.True(t, types.IsEmpty(types.NewXText("")))
	assert.False(t, types.IsEmpty(types.NewXText("a")))
	assert.False(t, types.IsEmpty(types.XBooleanFalse))
	assert.False(t, types.IsEmpty(types.XBooleanTrue))
	assert.False(t, types.IsEmpty(types.NewXNumberFromInt(0)))
	assert.False(t, types.IsEmpty(types.NewXNumberFromInt(123)))
}

func TestReduce(t *testing.T) {
	env := utils.NewEnvironmentBuilder().Environment()

	assert.Nil(t, types.Reduce(env, nil))
	assert.Equal(t, types.NewXText("Hello"), types.Reduce(env, NewTestXObject("Hello", 123)))
}
