package flows

import (
	"testing"

	"github.com/nyaruka/gocommon/urns"
	"github.com/stretchr/testify/assert"
)

func TestResolve(t *testing.T) {
	urnList := URNList{"tel:+250781234567", "twitter:134252511151#billy_bob", "tel:+250781111222"}

	testCases := []struct {
		key      string
		hasValue bool
		value    interface{}
	}{
		{"0", true, urns.URN("tel:+250781234567")},
		{"1", true, urns.URN("twitter:134252511151#billy_bob")},
		{"2", true, urns.URN("tel:+250781111222")},
		{"3", false, ""}, // index out of range
		{"tel", true, URNList{"tel:+250781234567", "tel:+250781111222"}},
		{"twitter", true, URNList{"twitter:134252511151#billy_bob"}},
		{"xxxxxx", false, ""}, // not a valid scheme
	}
	for _, tc := range testCases {
		val := urnList.Resolve(tc.key)
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
