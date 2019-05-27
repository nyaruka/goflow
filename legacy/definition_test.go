package legacy_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/routers"
	"github.com/nyaruka/goflow/legacy"
	"github.com/nyaruka/goflow/test"
	"github.com/nyaruka/goflow/utils"

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
			"x": 0, "y": 200, 
			"destination": null, 
			"exit_uuid": "cfcf5cef-49f9-41a6-886b-f466575a3045",
			"actions": []
		},
		{
			"uuid": "833fc698-d590-42dc-93e1-39e701b7e8e4",
			"x": 0, "y": 400, 
			"destination": null, 
			"exit_uuid": "da3e7eaf-c087-4e80-97b5-0b2e217fcc93",
			"actions": []
		},
		{
			"uuid": "42ff72d3-5f4d-4dbf-89c9-8a97864dabcd",
			"x": 0, "y": 600, 
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
}

type RuleSetMigrationTest struct {
	LegacyRuleSet        json.RawMessage `json:"legacy_ruleset"`
	ExpectedNode         json.RawMessage `json:"expected_node"`
	ExpectedLocalization json.RawMessage `json:"expected_localization"`
	ExpectedUI           json.RawMessage `json:"expected_ui"`
}

func TestFlowMigration(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/flows.json")
	require.NoError(t, err)

	var tests []FlowMigrationTest
	err = json.Unmarshal(data, &tests)
	require.NoError(t, err)

	defer utils.SetUUIDGenerator(utils.DefaultUUIDGenerator)

	for _, tc := range tests {
		utils.SetUUIDGenerator(test.NewSeededUUIDGenerator(123456))

		legacyFlow, err := legacy.ReadLegacyFlow(tc.Legacy)
		require.NoError(t, err)

		migratedFlow, err := legacyFlow.Migrate(true, "https://myfiles.com")
		require.NoError(t, err)

		migratedFlowJSON, err := json.Marshal(migratedFlow)
		require.NoError(t, err)

		test.AssertEqualJSON(t, tc.Expected, migratedFlowJSON, "migrated flow produced unexpected JSON")
	}
}

func TestActionMigration(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/actions.json")
	require.NoError(t, err)

	var tests []ActionMigrationTest
	err = json.Unmarshal(data, &tests)
	require.NoError(t, err)

	for _, tc := range tests {
		if tc.LegacyFlowType == "" {
			tc.LegacyFlowType = "F"
		}

		legacyFlowJSON := fmt.Sprintf(legacyActionHolderDef, tc.LegacyFlowType, string(tc.LegacyAction))
		legacyFlow, err := legacy.ReadLegacyFlow(json.RawMessage(legacyFlowJSON))
		require.NoError(t, err)

		migratedFlow, err := legacyFlow.Migrate(false, "https://myfiles.com")
		require.NoError(t, err)

		migratedAction := migratedFlow.Nodes()[0].Actions()[0]
		migratedActionJSON, err := utils.JSONMarshal(migratedAction)
		require.NoError(t, err)

		test.AssertEqualJSON(t, tc.ExpectedAction, migratedActionJSON, "migrated action produced unexpected JSON")

		checkFlowLocalization(t, migratedFlow, tc.ExpectedLocalization)
	}
}

func TestTestMigration(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/tests.json")
	require.NoError(t, err)

	var tests []TestMigrationTest
	err = json.Unmarshal(data, &tests)
	require.NoError(t, err)

	defer utils.SetUUIDGenerator(utils.DefaultUUIDGenerator)

	for _, tc := range tests {
		utils.SetUUIDGenerator(test.NewSeededUUIDGenerator(123456))

		legacyFlowJSON := fmt.Sprintf(legacyTestHolderDef, string(tc.LegacyTest))
		legacyFlow, err := legacy.ReadLegacyFlow(json.RawMessage(legacyFlowJSON))
		require.NoError(t, err)

		migratedFlow, err := legacyFlow.Migrate(false, "https://myfiles.com")
		require.NoError(t, err)

		migratedRouter := migratedFlow.Nodes()[0].Router().(*routers.SwitchRouter)

		if len(migratedRouter.Cases()) == 0 {
			t.Errorf("Got no migrated case from legacy test:\n%s\n\n", string(tc.LegacyTest))
		} else {
			migratedCase := migratedRouter.Cases()[0]
			migratedCaseJSON, err := utils.JSONMarshal(migratedCase)
			require.NoError(t, err)

			test.AssertEqualJSON(t, tc.ExpectedCase, migratedCaseJSON, "migrated test produced unexpected JSON")

			checkFlowLocalization(t, migratedFlow, tc.ExpectedLocalization)
		}
	}
}

func TestRuleSetMigration(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/rulesets.json")
	require.NoError(t, err)

	var tests []RuleSetMigrationTest
	err = json.Unmarshal(data, &tests)
	require.NoError(t, err)

	defer utils.SetUUIDGenerator(utils.DefaultUUIDGenerator)

	for _, tc := range tests {
		utils.SetUUIDGenerator(test.NewSeededUUIDGenerator(123456))

		legacyFlowJSON := fmt.Sprintf(legacyRuleSetHolderDef, string(tc.LegacyRuleSet))
		legacyFlow, err := legacy.ReadLegacyFlow(json.RawMessage(legacyFlowJSON))
		require.NoError(t, err)

		migratedFlow, err := legacyFlow.Migrate(true, "https://myfiles.com")
		require.NoError(t, err)

		// check we now have a new node in addition to the 3 actionsets used as destinations
		if len(migratedFlow.Nodes()) <= 3 {
			t.Errorf("Got no migrated nodes from legacy ruleset:\n%s\n\n", string(tc.LegacyRuleSet))
		} else {
			// find the new node which might be before or after the actionset nodes
			var migratedNode flows.Node
			for _, node := range migratedFlow.Nodes() {
				if node.UUID() != "5b977652-91e3-48be-8e86-7c8094b4aa8f" && node.UUID() != "833fc698-d590-42dc-93e1-39e701b7e8e4" && node.UUID() != "42ff72d3-5f4d-4dbf-89c9-8a97864dabcd" {
					migratedNode = node
				}
			}

			migratedNodeJSON, err := utils.JSONMarshal(migratedNode)
			require.NoError(t, err)

			test.AssertEqualJSON(t, tc.ExpectedNode, migratedNodeJSON, "migrated ruleset produced unexpected JSON")

			uiMap, err := utils.JSONDecodeGeneric(migratedFlow.UI())
			require.NoError(t, err)

			uiNodesMap := uiMap.(map[string]interface{})["nodes"].(map[string]interface{})

			migratedNodeUIJSON, err := utils.JSONMarshal(uiNodesMap[string(migratedNode.UUID())])
			require.NoError(t, err)

			test.AssertEqualJSON(t, tc.ExpectedUI, migratedNodeUIJSON, "migrated ruleset produced unexpected UI JSON")

			checkFlowLocalization(t, migratedFlow, tc.ExpectedLocalization)
		}
	}
}

func checkFlowLocalization(t *testing.T, flow flows.Flow, expectedLocalizationRaw json.RawMessage) {
	actualLocalizationJSON, err := utils.JSONMarshal(flow.Localization())
	require.NoError(t, err)

	expectedLocalizationJSON, _ := utils.JSONMarshal(expectedLocalizationRaw)
	require.NoError(t, err)

	test.AssertEqualJSON(t, expectedLocalizationJSON, actualLocalizationJSON, "migrated localization produced unexpected JSON")
}
