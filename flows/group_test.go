package flows_test

import (
	"testing"

	"github.com/nyaruka/goflow/excellent"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/test"
	"github.com/nyaruka/goflow/utils"

	"github.com/stretchr/testify/assert"
)

func TestGroupListResolve(t *testing.T) {
	customers := test.NewGroup("Customers", "")
	testers := test.NewGroup("Testers", "")
	males := test.NewGroup("Males", "gender = \"M\"")
	groups := flows.NewGroupList([]*flows.Group{customers, testers, males})

	env := utils.NewEnvironmentBuilder().Build()
	context := types.NewXMap(map[string]types.XValue{"groups": groups})

	testCases := []struct {
		expression string
		hasValue   bool
		value      interface{}
	}{
		{"groups[0]", true, customers},
		{"groups[1]", true, testers},
		{"groups[2]", true, males},
		{"groups[-1]", true, males},
		{"groups[3]", false, nil}, // index out of range
	}
	for _, tc := range testCases {
		value := excellent.EvaluateExpression(env, context, tc.expression)
		err, isErr := value.(error)

		if tc.hasValue && isErr {
			t.Errorf("Got unexpected error resolving %s: %s", tc.expression, err)
		}

		if !tc.hasValue && !isErr {
			t.Errorf("Did not get expected error resolving %s", tc.expression)
		}

		if tc.hasValue {
			assert.Equal(t, tc.value, value)
		}
	}
}
