package flows_test

import (
	"github.com/satori/go.uuid"
	"testing"

	"github.com/nyaruka/goflow/excellent"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"

	"github.com/stretchr/testify/assert"
)

func TestGroupListResolve(t *testing.T) {
	customers := flows.NewGroup(flows.GroupUUID(uuid.NewV4().String()), "Customers", "")
	testers := flows.NewGroup(flows.GroupUUID(uuid.NewV4().String()), "Testers", "")
	males := flows.NewGroup(flows.GroupUUID(uuid.NewV4().String()), "Males", "gender = \"M\"")
	urnList := flows.NewGroupList([]*flows.Group{customers, testers, males})

	env := utils.NewDefaultEnvironment()

	testCases := []struct {
		key      string
		hasValue bool
		value    interface{}
	}{
		{"0", true, customers},
		{"1", true, testers},
		{"2", true, males},
		{"-1", true, males},
		{"3", false, nil}, // index out of range
	}
	for _, tc := range testCases {
		val := excellent.ResolveVariable(env, urnList, tc.key)

		err, isErr := val.(error)

		if tc.hasValue && isErr {
			t.Errorf("Got unexpected error resolving %s: %s", tc.key, err)
		}

		if !tc.hasValue && !isErr {
			t.Errorf("Did not get expected error resolving %s", tc.key)
		}

		if tc.hasValue {
			assert.Equal(t, tc.value, val)
		}
	}
}
