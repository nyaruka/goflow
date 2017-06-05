package flows

import (
	"encoding/json"
	"testing"

	"github.com/nyaruka/goflow/utils"
)

func TestFields(t *testing.T) {
	var fieldTests = []struct {
		JSON     []byte
		lookup   string
		expected string
	}{
		{[]byte(`{}`), "key", ""},
		{[]byte(`{ "name": { "field_name": "Name", "value": "Ryan Lewis", "created_on": "2000-01-01T00:00:00.000000000-00:00"}}`), "key", ""},
		{[]byte(`{ "name": { "field_name": "Name", "value": "Ryan Lewis", "created_on": "2000-01-01T00:00:00.000000000-00:00"}}`), "name", "Ryan Lewis"},
		{[]byte(`{ "last_name": { "field_name": "Last Name", "value": "Lewis", "created_on": "2000-01-01T00:00:00.000000000-00:00"}}`), "last_name", "Lewis"},
		{[]byte(`{ "last_name": { "field_name": "Last Name", "value": "Lewis", "created_on": "2000-01-01T00:00:00.000000000-00:00"}}`), "Last Name", "Lewis"},
	}

	env := utils.NewDefaultEnvironment()
	for _, test := range fieldTests {
		fields := NewFields()
		err := json.Unmarshal(test.JSON, fields)
		if err != nil {
			t.Errorf("Error unmarshalling: '%s'", err)
			continue
		}
		value := utils.ResolveVariable(env, fields, test.lookup)

		valueStr, _ := utils.ToString(env, value)
		if valueStr != test.expected {
			t.Errorf("Expected: '%s' Got: '%s' for lookup: '%s' and Fields:\n%s", test.expected, valueStr, test.lookup, test.JSON)
		}
	}
}
