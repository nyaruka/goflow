package flows_test

import (
	"testing"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"

	"github.com/stretchr/testify/assert"
)

func TestURNListResolve(t *testing.T) {
	urnList := flows.URNList{
		flows.NewContactURN("tel:+250781234567", nil),
		flows.NewContactURN("twitter:134252511151#billy_bob", nil),
		flows.NewContactURN("tel:+250781111222", nil),
	}

	env := utils.NewDefaultEnvironment()

	testCases := []struct {
		key      string
		hasValue bool
		value    interface{}
	}{
		{"0", true, flows.NewContactURN("tel:+250781234567", nil)},
		{"1", true, flows.NewContactURN("twitter:134252511151#billy_bob", nil)},
		{"2", true, flows.NewContactURN("tel:+250781111222", nil)},
		{"-1", true, flows.NewContactURN("tel:+250781111222", nil)},
		{"3", false, nil}, // index out of range
		{"tel", true, flows.URNList{flows.NewContactURN("tel:+250781234567", nil), flows.NewContactURN("tel:+250781111222", nil)}},
		{"twitter", true, flows.URNList{flows.NewContactURN("twitter:134252511151#billy_bob", nil)}},
		{"xxxxxx", false, ""}, // not a valid scheme
	}
	for _, tc := range testCases {
		val := utils.ResolveVariable(env, urnList, tc.key)

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
