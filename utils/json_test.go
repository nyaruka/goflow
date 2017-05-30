package utils

import (
	"testing"
)

func TestJSON(t *testing.T) {
	var jsonTests = []struct {
		JSON     string
		lookup   string
		expected string
	}{
		{`["one", "two", "three"]`, "0", "one"},
		{`{"1": "one"}`, "1", "one"},
		{`{"arr": ["one", "two"]}`, "arr[1]", "two"},
		{`{"arr": ["one", "two"]}`, "arr.1", "two"},
		{`{"key": {"key2": "val2"}}`, "key.key2", "val2"},
		{`{"key": {"key-with-dash": "val2"}}`, `key["key-with-dash"]`, "val2"},
		{`{"key": {"key with space": "val2"}}`, `key["key with space"]`, "val2"},
	}

	env := NewDefaultEnvironment()
	for _, test := range jsonTests {
		fragment := JSONFragment(test.JSON)
		value := ResolveVariable(env, fragment, test.lookup)

		valueStr, err := ToString(env, value)
		if err != nil {
			t.Errorf("Error getting string value for lookup: '%s' and JSON fragment:\n%s", test.lookup, test.JSON)
			continue
		}
		if valueStr != test.expected {
			t.Errorf("FExpected: '%s' Got: '%s' for lookup: '%s' and JSON fragment:\n%s", test.expected, valueStr, test.lookup, test.JSON)
		}
	}
}
