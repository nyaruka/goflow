package legacy_test

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/flows/definition/legacy"
	"github.com/nyaruka/goflow/test"

	"github.com/buger/jsonparser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var legacyActionHolderDef = `
{
	"base_language": "eng",
	"entry": "10e483a8-5ffb-4c4f-917b-d43ce86c1d65", 
	"flow_type": "%s",
	"action_sets": [{
		"uuid": "10e483a8-5ffb-4c4f-917b-d43ce86c1d65",
		"y": 100, 
		"x": 100, 
		"destination": null, 
		"exit_uuid": "cfcf5cef-49f9-41a6-886b-f466575a3045",
		"actions": [%s]
	}],
	"metadata": {
		"uuid": "50c3706e-fedb-42c0-8eab-dda3335714b7",
		"name": "TestFlow"
	}
}
`

var legacyTestHolderDef = `
{
	"base_language": "eng",
	"entry": "10e483a8-5ffb-4c4f-917b-d43ce86c1d65",
	"flow_type": "F",
	"rule_sets": [{
		"uuid": "10e483a8-5ffb-4c4f-917b-d43ce86c1d65",
		"rules": [{
			"test": %s, 
			"category": {"eng": "All Responses"}, 
			"destination": null, 
			"uuid": "c072ecb5-0686-40ea-8ed3-898dc1349783", 
			"destination_type": null
		}],
		"ruleset_type": "wait_message", 
		"label": "Name", 
		"operand": "@step.value", 
		"finished_key": null,
		"response_type": "", 
		"y": 0, 
		"x": 100, 
		"config": {}
	}],
	"metadata": {
		"uuid": "50c3706e-fedb-42c0-8eab-dda3335714b7",
		"name": "TestFlow"
	}
}
`

var legacyRuleSetHolderDef = `
{
	"base_language": "eng",
	"entry": "10e483a8-5ffb-4c4f-917b-d43ce86c1d65",
	"flow_type": "F", 
	"rule_sets": [%s],
	"action_sets": [
		{
			"uuid": "5b977652-91e3-48be-8e86-7c8094b4aa8f",
			"x": 0, "y": 2200, 
			"destination": null, 
			"exit_uuid": "cfcf5cef-49f9-41a6-886b-f466575a3045",
			"actions": []
		},
		{
			"uuid": "833fc698-d590-42dc-93e1-39e701b7e8e4",
			"x": 0, "y": 2400, 
			"destination": null, 
			"exit_uuid": "da3e7eaf-c087-4e80-97b5-0b2e217fcc93",
			"actions": []
		},
		{
			"uuid": "42ff72d3-5f4d-4dbf-89c9-8a97864dabcd",
			"x": 0, "y": 2600, 
			"destination": null, 
			"exit_uuid": "6a8cb81b-1b59-4cfb-b00e-575ccbafd3ba",
			"actions": []
		}
	],
	"metadata": {
		"uuid": "50c3706e-fedb-42c0-8eab-dda3335714b7",
		"name": "TestFlow"
	}
}
`

type FlowMigrationTest struct {
	Legacy   json.RawMessage `json:"legacy"`
	Expected json.RawMessage `json:"expected"`
}

type ActionMigrationTest struct {
	LegacyAction         json.RawMessage `json:"legacy_action"`
	LegacyFlowType       string          `json:"legacy_flow_type"`
	ExpectedAction       json.RawMessage `json:"expected_action"`
	ExpectedLocalization json.RawMessage `json:"expected_localization"`
}

type TestMigrationTest struct {
	LegacyTest           json.RawMessage `json:"legacy_test"`
	ExpectedCase         json.RawMessage `json:"expected_case"`
	ExpectedLocalization json.RawMessage `json:"expected_localization"`
	ExpectedUI           json.RawMessage `json:"expected_ui"`
}

type RuleSetMigrationTest struct {
	LegacyRuleSet        json.RawMessage `json:"legacy_ruleset"`
	ExpectedNode         json.RawMessage `json:"expected_node"`
	ExpectedLocalization json.RawMessage `json:"expected_localization"`
	ExpectedUI           json.RawMessage `json:"expected_ui"`
}

func TestFlowMigration(t *testing.T) {
	data, err := os.ReadFile("testdata/flows.json")
	require.NoError(t, err)

	var tests []FlowMigrationTest
	err = jsonx.Unmarshal(data, &tests)
	require.NoError(t, err)

	defer uuids.SetGenerator(uuids.DefaultGenerator)

	for _, tc := range tests {
		uuids.SetGenerator(uuids.NewSeededGenerator(123456))

		migratedFlowJSON, err := legacy.MigrateDefinition(tc.Legacy, "https://myfiles.com")
		require.NoError(t, err)

		test.AssertEqualJSON(t, tc.Expected, migratedFlowJSON, "migrated flow produced unexpected JSON")
	}
}

func TestActionMigration(t *testing.T) {
	data, err := os.ReadFile("testdata/actions.json")
	require.NoError(t, err)

	var tests []ActionMigrationTest
	err = jsonx.Unmarshal(data, &tests)
	require.NoError(t, err)

	for _, tc := range tests {
		if tc.LegacyFlowType == "" {
			tc.LegacyFlowType = "F"
		}

		legacyFlowJSON := fmt.Sprintf(legacyActionHolderDef, tc.LegacyFlowType, string(tc.LegacyAction))
		migratedFlowJSON, err := legacy.MigrateDefinition(json.RawMessage(legacyFlowJSON), "https://myfiles.com")
		require.NoError(t, err)

		migratedActionJSON, _, _, err := jsonparser.Get(migratedFlowJSON, "nodes", "[0]", "actions", "[0]")
		require.NoError(t, err)

		test.AssertEqualJSON(t, tc.ExpectedAction, migratedActionJSON, "migrated action produced unexpected JSON")

		migratedLocalizationJSON, _, _, err := jsonparser.Get(migratedFlowJSON, "localization")
		require.NoError(t, err)

		test.AssertEqualJSON(t, tc.ExpectedLocalization, migratedLocalizationJSON, "migrated localization produced unexpected JSON")
	}
}

func TestTestMigration(t *testing.T) {
	data, err := os.ReadFile("testdata/tests.json")
	require.NoError(t, err)

	var tests []TestMigrationTest
	err = jsonx.Unmarshal(data, &tests)
	require.NoError(t, err)

	defer uuids.SetGenerator(uuids.DefaultGenerator)

	for _, tc := range tests {
		uuids.SetGenerator(uuids.NewSeededGenerator(123456))

		legacyFlowJSON := fmt.Sprintf(legacyTestHolderDef, string(tc.LegacyTest))
		migratedFlowJSON, err := legacy.MigrateDefinition(json.RawMessage(legacyFlowJSON), "https://myfiles.com")
		require.NoError(t, err)

		migratedRouterJSON, _, _, err := jsonparser.Get(migratedFlowJSON, "nodes", "[0]", "router")
		require.NoError(t, err)

		migratedCaseJSON, _, _, err := jsonparser.Get(migratedRouterJSON, "cases", "[0]")
		require.NoError(t, err)

		test.AssertEqualJSON(t, tc.ExpectedCase, migratedCaseJSON, "migrated test produced unexpected JSON")

		migratedLocalizationJSON, _, _, err := jsonparser.Get(migratedFlowJSON, "localization")
		require.NoError(t, err)

		test.AssertEqualJSON(t, tc.ExpectedLocalization, migratedLocalizationJSON, "migrated localization produced unexpected JSON")

		if tc.ExpectedUI != nil {
			migratedCaseUUID, _, _, err := jsonparser.Get(migratedCaseJSON, "uuid")
			require.NoError(t, err)

			fmt.Println(string(migratedFlowJSON))

			migratedCaseUIJSON, _, _, err := jsonparser.Get(migratedFlowJSON, "_ui", "nodes", "10e483a8-5ffb-4c4f-917b-d43ce86c1d65", "config", "cases", string(migratedCaseUUID))
			require.NoError(t, err)

			test.AssertEqualJSON(t, tc.ExpectedUI, migratedCaseUIJSON, "migrated ruleset produced unexpected UI JSON")
		}
	}
}

func TestRuleSetMigration(t *testing.T) {
	data, err := os.ReadFile("testdata/rulesets.json")
	require.NoError(t, err)

	var tests []RuleSetMigrationTest
	err = jsonx.Unmarshal(data, &tests)
	require.NoError(t, err)

	defer uuids.SetGenerator(uuids.DefaultGenerator)

	for _, tc := range tests {
		uuids.SetGenerator(uuids.NewSeededGenerator(123456))

		legacyFlowJSON := fmt.Sprintf(legacyRuleSetHolderDef, string(tc.LegacyRuleSet))
		migratedFlowJSON, err := legacy.MigrateDefinition(json.RawMessage(legacyFlowJSON), "https://myfiles.com")
		require.NoError(t, err)

		migratedNodeJSON, _, _, err := jsonparser.Get(migratedFlowJSON, "nodes", "[0]")
		require.NoError(t, err)

		test.AssertEqualJSON(t, tc.ExpectedNode, migratedNodeJSON, "migrated ruleset produced unexpected JSON")

		migratedNodeUUID, _, _, err := jsonparser.Get(migratedNodeJSON, "uuid")
		require.NoError(t, err)

		migratedNodeUIJSON, _, _, err := jsonparser.Get(migratedFlowJSON, "_ui", "nodes", string(migratedNodeUUID))
		require.NoError(t, err)

		test.AssertEqualJSON(t, tc.ExpectedUI, migratedNodeUIJSON, "migrated ruleset produced unexpected UI JSON")

		migratedLocalizationJSON, _, _, err := jsonparser.Get(migratedFlowJSON, "localization")
		require.NoError(t, err)

		test.AssertEqualJSON(t, tc.ExpectedLocalization, migratedLocalizationJSON, "migrated localization produced unexpected JSON")
	}
}

func TestIsPossibleDefinition(t *testing.T) {
	// try reading empty JSON
	assert.False(t, legacy.IsPossibleDefinition([]byte(`{}`)))
	assert.True(t, legacy.IsPossibleDefinition([]byte(`{"flow_type":"M"}`)))

	// try with new flow
	assert.False(t, legacy.IsPossibleDefinition([]byte(`{
		"uuid": "76f0a02f-3b75-4b86-9064-e9195e1b3a02",
		"name": "Simple",
		"spec_version": "13.0",
		"language": "eng",
		"type": "messaging",
		"nodes": []
	}`)))

	// try with legacy flow
	assert.True(t, legacy.IsPossibleDefinition([]byte(`{
		"metadata": {
			"uuid": "76f0a02f-3b75-4b86-9064-e9195e1b3a02",
			"name": "Simple",
			"revision": 1
		},
		"base_language": "eng",
		"flow_type": "F",
		"entry"
		"version": 11,
		"action_sets": [],
		"rule_sets": []
	}`)))

	// try with legacy flow without metadata or flow_type
	assert.True(t, legacy.IsPossibleDefinition([]byte(`{
		"base_language": "eng",
		"entry"
		"version": 11,
		"action_sets": [],
		"rule_sets": []
	}`)))
}

func TestMigrateDefinition(t *testing.T) {
	migrated, err := legacy.MigrateDefinition([]byte(`{
		"flow_type": "S", 
		"action_sets": [],
		"rule_sets": [],
		"base_language": "eng",
		"metadata": {
			"uuid": "061be894-4507-470c-a20b-34273bf915be",
			"name": "Survey"
		}
	}`), "")

	assert.NoError(t, err)
	test.AssertEqualJSON(t, []byte(`{
		"uuid": "061be894-4507-470c-a20b-34273bf915be",
		"name": "Survey",
		"spec_version": "13.0.0",
		"type": "messaging_offline",
		"expire_after_minutes": 0,
		"language": "eng",
		"localization": {},
		"nodes": [],
		"revision": 0,
		"_ui": {
			"nodes": {},
			"stickies": {}
		}
	}`), migrated, "migrated flow mismatch")
}
