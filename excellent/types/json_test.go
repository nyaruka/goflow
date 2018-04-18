package types_test

import (
	"testing"

	"github.com/nyaruka/goflow/excellent"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/utils"

	"github.com/stretchr/testify/assert"
)

func TestXJSONResolve(t *testing.T) {
	var jsonTests = []struct {
		JSON     []byte
		lookup   string
		expected types.XValue
		hasError bool
	}{
		// error cases
		{nil, "key", nil, true},
		{[]byte(`malformed`), "key", nil, true},

		// different data types in an object
		{[]byte(`{"foo": "x", "bar": "one"}`), "bar", types.NewXText("one"), false},
		{[]byte(`{"foo": "x", "bar": 1.23}`), "bar", types.RequireXNumberFromString("1.23"), false},
		{[]byte(`{"foo": "x", "bar": true}`), "bar", types.NewXBoolean(true), false},
		{[]byte(`{"foo": "x", "bar": null}`), "bar", nil, false},

		// different data types in an array
		{[]byte(`["foo", "one"]`), "1", types.NewXText("one"), false},
		{[]byte(`["foo", 1.23]`), "1", types.RequireXNumberFromString("1.23"), false},
		{[]byte(`["foo", true]`), "1", types.NewXBoolean(true), false},
		{[]byte(`["foo", null]`), "1", nil, false},

		{[]byte(`["one", "two", "three"]`), "0", types.NewXText("one"), false},
		{[]byte(`["escaped \"string\""]`), "0", types.NewXText(`escaped "string"`), false},
		{[]byte(`{"1": "one"}`), "1", types.NewXText("one"), false}, // map key is numerical string
		{[]byte(`{"arr": ["one", "two"]}`), "arr[1]", types.NewXText("two"), false},
		{[]byte(`{"arr": ["one", "two"]}`), "arr.1", types.NewXText("two"), false},
		{[]byte(`{"key": {"key2": "val2"}}`), "key.key2", types.NewXText("val2"), false},
		{[]byte(`{"key": {"key-with-dash": "val2"}}`), `key["key-with-dash"]`, types.NewXText("val2"), false},
		{[]byte(`{"key": {"key with space": "val2"}}`), `key["key with space"]`, types.NewXText("val2"), false},

		{[]byte(`{"arr": ["one", "two"]}`), "arr", types.NewXJSONArray([]byte(`["one", "two"]`)), false},
		{[]byte(`{"arr": {"foo": "bar"}}`), "arr", types.NewXJSONObject([]byte(`{"foo": "bar"}`)), false},
	}

	env := utils.NewDefaultEnvironment()
	for _, test := range jsonTests {
		fragment := types.JSONToXValue(test.JSON)
		value := excellent.ResolveValue(env, fragment, test.lookup)
		err, _ := value.(error)

		if test.hasError {
			assert.Error(t, err, "expected error resolving '%s' in '%s'", test.lookup, test.JSON)
		} else {
			assert.NoError(t, err, "unexpected error resolving '%s' in '%s'", test.lookup, test.JSON)
			assert.Equal(t, test.expected, value, "Actual '%s' does not match expected '%s' resolving '%s' in '%s'", value, test.expected, test.lookup, test.JSON)
		}
	}
}
