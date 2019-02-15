package flows

import (
	"encoding/json"
	"testing"

	"github.com/nyaruka/goflow/excellent"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/utils"

	"github.com/stretchr/testify/assert"
)

func TestResults(t *testing.T) {
	var ERROR = types.NewXErrorf("any error")

	var tests = []struct {
		JSON     []byte
		lookup   string
		expected types.XValue
	}{
		{[]byte(`{}`), "key", ERROR},
		{[]byte(`{ "name": { "result_name": "Name", "value": "Ryan Lewis", "node": "uuid", "created_on": "2000-01-01T00:00:00.000000000-00:00"}}`), `results.key`, ERROR},
		{[]byte(`{ "name": { "result_name": "Name", "value": "Ryan Lewis", "node": "uuid", "created_on": "2000-01-01T00:00:00.000000000-00:00"}}`), `results.name`, types.NewXText("Ryan Lewis")},
		{[]byte(`{ "last_name": { "result_name": "Last Name", "value": "Lewis", "node": "uuid", "created_on": "2000-01-01T00:00:00.000000000-00:00"}}`), `results.last_name`, types.NewXText("Lewis")},
		{[]byte(`{ "last_name": { "result_name": "Last Name", "value": "Lewis", "node": "uuid", "created_on": "2000-01-01T00:00:00.000000000-00:00"}}`), `results["Last Name"]`, types.NewXText("Lewis")},
	}

	env := utils.NewEnvironmentBuilder().Build()
	for _, test := range tests {
		results := NewResults()
		err := json.Unmarshal(test.JSON, &results)
		if err != nil {
			t.Errorf("Error unmarshalling: '%s'", err)
			continue
		}
		context := types.NewXMap(map[string]types.XValue{"results": results})
		value := excellent.EvaluateExpression(env, context, test.lookup)

		// don't check error equality - just check that we got an error if we expected one
		if test.expected == ERROR {
			assert.True(t, types.IsXError(value), "expecting error, got %T{%s} for lookup %s", value, value, test.lookup)
		} else {
			if !types.Equals(env, value, test.expected) {
				t.Errorf("Expected: '%s' Got: '%s' for lookup: '%s' and Results:\n%s", test.expected, value, test.lookup, test.JSON)
			}
		}
	}
}
