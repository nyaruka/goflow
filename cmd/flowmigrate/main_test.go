package main_test

import (
	"fmt"
	"strings"
	"testing"

	main "github.com/nyaruka/goflow/cmd/flowmigrate"
	"github.com/nyaruka/goflow/flows/definition"
	"github.com/nyaruka/goflow/test"
	"github.com/stretchr/testify/require"
)

func TestMigrate(t *testing.T) {
	testCases := []struct {
		input  string
		output string
	}{
		{ // a legacy flow
			input: `{
				"metadata": {
					"uuid": "76f0a02f-3b75-4b86-9064-e9195e1b3a02",
					"name": "Empty",
					"revision": 1
				},
				"base_language": "eng",
				"flow_type": "F",
				"action_sets": [],
				"rule_sets": []
			}`,
			output: fmt.Sprintf(`{
				"uuid": "76f0a02f-3b75-4b86-9064-e9195e1b3a02",
				"name": "Empty",
				"spec_version": "%s",
				"language": "eng",
				"type": "messaging",
				"revision": 1,
				"expire_after_minutes": 0,
				"localization": {},
				"nodes": [],
				"_ui": {
					"nodes": {},
					"stickies": {}
				}
			}`, definition.CurrentSpecVersion),
		},
		{ // a new flow
			input: `{
				"uuid": "76f0a02f-3b75-4b86-9064-e9195e1b3a02",
				"name": "Empty",
				"spec_version": "13.0.0",
				"language": "eng",
				"type": "messaging",
				"revision": 1,
				"expire_after_minutes": 0,
				"nodes": [],
				"_ui": {
					"nodes": {},
					"stickies": {}
				}
			}`,
			output: fmt.Sprintf(`{
				"uuid": "76f0a02f-3b75-4b86-9064-e9195e1b3a02",
				"name": "Empty",
				"spec_version": "%s",
				"language": "eng",
				"type": "messaging",
				"revision": 1,
				"expire_after_minutes": 0,
				"nodes": [],
				"_ui": {
					"nodes": {},
					"stickies": {}
				}
			}`, definition.CurrentSpecVersion),
		},
	}

	for _, tc := range testCases {
		input := strings.NewReader(tc.input)

		migrated, err := main.Migrate(input, nil, "http://temba.io/", true)
		require.NoError(t, err)

		test.AssertEqualJSON(t, []byte(tc.output), migrated, "Migrated flow mismatch")
	}
}
