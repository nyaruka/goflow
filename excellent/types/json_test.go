package types_test

import (
	"testing"

	"github.com/nyaruka/goflow/excellent"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/utils"

	"github.com/stretchr/testify/assert"
)

func TestXJSON(t *testing.T) {
	jobj := types.JSONToXValue([]byte(`{"foo": "x", "bar": null}`)).(types.XJSONObject)
	assert.Equal(t, `{"foo": "x", "bar": null}`, jobj.String())
	assert.Equal(t, `json object`, jobj.Describe())

	jarr := types.JSONToXValue([]byte(`["one", "two", "three"]`)).(types.XJSONArray)
	assert.Equal(t, `["one", "two", "three"]`, jarr.String())
	assert.Equal(t, 3, jarr.Length())
	assert.Equal(t, types.NewXText("two"), jarr.Index(1))
	assert.True(t, types.IsXError(jarr.Index(7)))
	assert.Equal(t, `json array`, jarr.Describe())

	num := types.JSONToXValue([]byte(`37.27903`)).(types.XNumber)
	assert.Equal(t, num, types.RequireXNumberFromString(`37.27903`))

	jerr := types.JSONToXValue([]byte(`fish`)).(types.XError)
	assert.Equal(t, `Unknown value type`, jerr.Error())
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

		{[]byte(`{"arr": ["one", "two"]}`), "j.arr", types.NewXJSONArray([]byte(`["one", "two"]`)), false},
		{[]byte(`{"arr": {"foo": "bar"}}`), "j.arr", types.NewXJSONObject([]byte(`{"foo": "bar"}`)), false},

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
