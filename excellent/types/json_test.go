package types_test

import (
	"testing"

	"github.com/nyaruka/goflow/excellent"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/utils"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestJSONResolve(t *testing.T) {
	var jsonTests = []struct {
		JSON     []byte
		lookup   string
		expected interface{}
		hasError bool
	}{
		// error cases
		{nil, "key", "", true},
		{[]byte(`malformed`), "key", "", true},

		// different data types in an object
		{[]byte(`{"foo": "x", "bar": "one"}`), "bar", "one", false},
		{[]byte(`{"foo": "x", "bar": 1.23}`), "bar", decimal.RequireFromString("1.23"), false},
		{[]byte(`{"foo": "x", "bar": true}`), "bar", true, false},
		{[]byte(`{"foo": "x", "bar": null}`), "bar", nil, false},

		// different data types in an object
		{[]byte(`["foo", "one"]`), "1", "one", false},
		{[]byte(`["foo", 1.23]`), "1", decimal.RequireFromString("1.23"), false},
		{[]byte(`["foo", true]`), "1", true, false},
		{[]byte(`["foo", null]`), "1", nil, false},

		{[]byte(`["one", "two", "three"]`), "0", "one", false},
		{[]byte(`["escaped \"string\""]`), "0", `escaped "string"`, false},
		{[]byte(`{"1": "one"}`), "1", "one", false}, // map key is numerical string
		{[]byte(`{"arr": ["one", "two"]}`), "arr[1]", "two", false},
		{[]byte(`{"arr": ["one", "two"]}`), "arr.1", "two", false},
		{[]byte(`{"key": {"key2": "val2"}}`), "key.key2", "val2", false},
		{[]byte(`{"key": {"key-with-dash": "val2"}}`), `key["key-with-dash"]`, "val2", false},
		{[]byte(`{"key": {"key with space": "val2"}}`), `key["key with space"]`, "val2", false},

		{[]byte(`{"arr": ["one", "two"]}`), "arr", types.JSONArray([]byte(`["one", "two"]`)), false},
		{[]byte(`{"arr": {"foo": "bar"}}`), "arr", types.JSONFragment([]byte(`{"foo": "bar"}`)), false},
	}

	env := utils.NewDefaultEnvironment()
	for _, test := range jsonTests {
		fragment := types.JSONFragment(test.JSON)
		value := excellent.ResolveVariable(env, fragment, test.lookup)
		err, _ := value.(error)

		if test.hasError {
			assert.Error(t, err, "expected error resolving '%s' in '%s'", test.lookup, test.JSON)
		} else {
			assert.NoError(t, err, "unexpected error resolving '%s' in '%s'", test.lookup, test.JSON)
			assert.Equal(t, test.expected, value, "Actual '%s' does not match expected '%s' resolving '%s' in '%s'", value, test.expected, test.lookup, test.JSON)
		}
	}
}
