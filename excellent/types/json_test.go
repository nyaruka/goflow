package types_test

import (
	"testing"

	"github.com/nyaruka/goflow/excellent"
	"github.com/nyaruka/goflow/excellent/test"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/utils"

	"github.com/stretchr/testify/assert"
)

func TestJSONToXValue(t *testing.T) {
	test.AssertEqual(t, types.NewXDict(map[string]types.XValue{
		"foo": types.NewXText("x"),
		"bar": nil,
		"sub": types.NewXDict(map[string]types.XValue{
			"x": types.NewXNumberFromInt(3),
		}),
	}), types.JSONToXValue([]byte(`{"foo": "x", "bar": null, "sub": {"x": 3}}`)))

	test.AssertEqual(t, types.NewXArray(
		types.NewXText("foo"),
		types.NewXNumberFromInt(123),
		nil,
		types.NewXArray(types.NewXNumberFromInt(2), types.NewXNumberFromInt(3)),
	), types.JSONToXValue([]byte(`["foo", 123, null, [2, 3]]`)))

	test.AssertEqual(t, types.RequireXNumberFromString(`37.27903`), types.JSONToXValue([]byte(`37.27903`)))

	xerr := types.JSONToXValue([]byte(`fish`)).(types.XError)
	assert.Equal(t, `Unknown value type`, xerr.Error())
}

func TestXJSONResolve(t *testing.T) {
	var jsonTests = []struct {
		JSON       []byte
		expression string
		expected   types.XValue
		hasError   bool
	}{
		// error cases
		{nil, "j.key", nil, true},
		{[]byte(`malformed`), "j.key", nil, true},

		// different data types in an object
		{[]byte(`{"foo": "x", "bar": "one"}`), "j.bar", types.NewXText("one"), false},
		{[]byte(`{"foo": "x", "bar": 1.23}`), "j.bar", types.RequireXNumberFromString("1.23"), false},
		{[]byte(`{"foo": "x", "bar": true}`), "j.bar", types.NewXBoolean(true), false},
		{[]byte(`{"foo": "x", "bar": null}`), "j.bar", nil, false},

		// different data types in an array
		{[]byte(`["foo", "one"]`), "j[1]", types.NewXText("one"), false},
		{[]byte(`["foo", 1.23]`), "j[1]", types.RequireXNumberFromString("1.23"), false},
		{[]byte(`["foo", true]`), "j[1]", types.NewXBoolean(true), false},
		{[]byte(`["foo", null]`), "j[1]", nil, false},

		{[]byte(`["one", "two", "three"]`), "j[0]", types.NewXText("one"), false},
		{[]byte(`["escaped \"string\""]`), "j[0]", types.NewXText(`escaped "string"`), false},
		{[]byte(`{"arr": ["one", "two"]}`), "j.arr[1]", types.NewXText("two"), false},
		{[]byte(`{"arr": ["one", "two"]}`), "j.arr[1]", types.NewXText("two"), false},
		{[]byte(`{"key": {"key2": "val2"}}`), "j.key.key2", types.NewXText("val2"), false},
		{[]byte(`{"key": {"key-with-dash": "val2"}}`), `j.key["key-with-dash"]`, types.NewXText("val2"), false},
		{[]byte(`{"key": {"key with space": "val2"}}`), `j.key["key with space"]`, types.NewXText("val2"), false},

		//{[]byte(`{"arr": ["one", "two"]}`), "j.arr", types.NewXJSONArray([]byte(`["one", "two"]`)), false},
		//{[]byte(`{"arr": {"foo": "bar"}}`), "j.arr", types.NewXJSONDict([]byte(`{"foo": "bar"}`)), false},

		// resolve errors
		{[]byte(`{"foo": "x", "bar": "one"}`), "j.zed", nil, true},
		{[]byte(`["foo", null]`), "j.0", nil, true},
		{[]byte(`["foo", null]`), "j[3]", nil, true},
	}

	env := utils.NewEnvironmentBuilder().Build()
	for _, test := range jsonTests {
		fragment := types.JSONToXValue(test.JSON)
		context := types.NewXDict(map[string]types.XValue{"j": fragment})

		value := excellent.EvaluateExpression(env, context, test.expression)
		err, _ := value.(error)

		if test.hasError {
			assert.Error(t, err, "expected error resolving '%s' for '%s'", test.expression, test.JSON)
		} else {
			assert.NoError(t, err, "unexpected error resolving '%s' for '%s'", test.expression, test.JSON)
			assert.Equal(t, test.expected, value, "Actual '%s' does not match expected '%s' resolving '%s' for '%s'", value, test.expected, test.expression, test.JSON)
		}
	}
}
