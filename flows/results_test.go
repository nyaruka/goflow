package flows

import (
	"encoding/json"
	"testing"

	"github.com/nyaruka/goflow/utils"
)

func TestResults(t *testing.T) {
	var resultTests = []struct {
		JSON     []byte
		lookup   string
		expected string
	}{
		{[]byte(`{}`), "key", ""},
		{[]byte(`{ "name": { "result_name": "Name", "value": "Ryan Lewis", "node": "uuid", "created_on": "2000-01-01T00:00:00.000000000-00:00"}}`), "key", ""},
		{[]byte(`{ "name": { "result_name": "Name", "value": "Ryan Lewis", "node": "uuid", "created_on": "2000-01-01T00:00:00.000000000-00:00"}}`), "name", "Ryan Lewis"},
		{[]byte(`{ "last_name": { "result_name": "Last Name", "value": "Lewis", "node": "uuid", "created_on": "2000-01-01T00:00:00.000000000-00:00"}}`), "last_name", "Lewis"},
		{[]byte(`{ "last_name": { "result_name": "Last Name", "value": "Lewis", "node": "uuid", "created_on": "2000-01-01T00:00:00.000000000-00:00"}}`), "Last Name", "Lewis"},
	}

	env := utils.NewDefaultEnvironment()
	for _, test := range resultTests {
		results := NewResults()
		err := json.Unmarshal(test.JSON, &results)
		if err != nil {
			t.Errorf("Error unmarshalling: '%s'", err)
			continue
		}
		value := utils.ResolveVariable(env, results, test.lookup)

		valueStr, _ := utils.ToString(env, value)
		if valueStr != test.expected {
			t.Errorf("Expected: '%s' Got: '%s' for lookup: '%s' and Results:\n%s", test.expected, valueStr, test.lookup, test.JSON)
		}
	}
}
