package legacy_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/routers"
	"github.com/nyaruka/goflow/legacy"
	"github.com/nyaruka/goflow/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var legacyActionHolderDef = `
[
	{
		"base_language": "eng",
		"entry": "10e483a8-5ffb-4c4f-917b-d43ce86c1d65", 
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
]
`

var legacyTestHolderDef = `
[
	{
		"base_language": "eng",
		"entry": "10e483a8-5ffb-4c4f-917b-d43ce86c1d65",
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
]
`

var legacyRuleSetHolderDef = `
[
	{
		"base_language": "eng",
		"entry": "10e483a8-5ffb-4c4f-917b-d43ce86c1d65",
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
]
`

type ActionMigrationTest struct {
	LegacyAction         json.RawMessage `json:"legacy_action"`
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
}

func TestActionMigration(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/migrations/actions.json")
	require.NoError(t, err)

	var tests []ActionMigrationTest
	err = json.Unmarshal(data, &tests)
	require.NoError(t, err)

	for _, test := range tests {
		legacyFlowsJSON := fmt.Sprintf(legacyActionHolderDef, string(test.LegacyAction))
		legacyFlows, err := readLegacyTestFlows(legacyFlowsJSON)
		require.NoError(t, err)

		migratedFlow, err := legacyFlows[0].Migrate()
		require.NoError(t, err)

		migratedAction := migratedFlow.Nodes()[0].Actions()[0]
		migratedActionEnvelope, _ := utils.EnvelopeFromTyped(migratedAction)
		migratedActionRaw, _ := json.Marshal(migratedActionEnvelope)
		migratedActionJSON, _ := utils.JSONMarshalPretty(migratedActionRaw)
		expectedActionJSON, _ := utils.JSONMarshalPretty(test.ExpectedAction)

		assert.Equal(t, string(expectedActionJSON), string(migratedActionJSON))

		checkFlowLocalization(t, migratedFlow, test.ExpectedLocalization)
	}
}

func TestTestMigration(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/migrations/tests.json")
	require.NoError(t, err)

	var tests []TestMigrationTest
	err = json.Unmarshal(data, &tests)
	require.NoError(t, err)

	defer utils.SetUUIDGenerator(utils.DefaultUUIDGenerator)

	for _, test := range tests {
		utils.SetUUIDGenerator(utils.NewSeededUUID4Generator(123456))

		legacyFlowsJSON := fmt.Sprintf(legacyTestHolderDef, string(test.LegacyTest))
		legacyFlows, err := readLegacyTestFlows(legacyFlowsJSON)
		require.NoError(t, err)

		migratedFlow, err := legacyFlows[0].Migrate()
		require.NoError(t, err)

		migratedRouter := migratedFlow.Nodes()[0].Router().(*routers.SwitchRouter)

		if len(migratedRouter.Cases) == 0 {
			t.Errorf("Got no migrated case from legacy test:\n%s\n\n", string(test.LegacyTest))
		} else {
			migratedCase := migratedRouter.Cases[0]
			migratedCaseRaw, _ := json.Marshal(migratedCase)
			migratedCaseJSON, _ := utils.JSONMarshalPretty(migratedCaseRaw)
			expectedCaseJSON, _ := utils.JSONMarshalPretty(test.ExpectedCase)

			assert.Equal(t, string(expectedCaseJSON), string(migratedCaseJSON))

			checkFlowLocalization(t, migratedFlow, test.ExpectedLocalization)
		}
	}
}

func TestRuleSetMigration(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/migrations/rulesets.json")
	require.NoError(t, err)

	var tests []RuleSetMigrationTest
	err = json.Unmarshal(data, &tests)
	require.NoError(t, err)

	defer utils.SetUUIDGenerator(utils.DefaultUUIDGenerator)

	for _, test := range tests {
		utils.SetUUIDGenerator(utils.NewSeededUUID4Generator(123456))

		legacyFlowsJSON := fmt.Sprintf(legacyRuleSetHolderDef, string(test.LegacyRuleSet))
		legacyFlows, err := readLegacyTestFlows(legacyFlowsJSON)
		require.NoError(t, err)

		migratedFlow, err := legacyFlows[0].Migrate()
		require.NoError(t, err)

		// check we now have a new node in addition to the 3 actionsets used as destinations
		if len(migratedFlow.Nodes()) <= 3 {
			t.Errorf("Got no migrated nodes from legacy ruleset:\n%s\n\n", string(test.LegacyRuleSet))
		} else {
			// find the new node which might be before or after the actionset nodes
			var migratedNode flows.Node
			for _, node := range migratedFlow.Nodes() {
				if node.UUID() != "5b977652-91e3-48be-8e86-7c8094b4aa8f" && node.UUID() != "833fc698-d590-42dc-93e1-39e701b7e8e4" && node.UUID() != "42ff72d3-5f4d-4dbf-89c9-8a97864dabcd" {
					migratedNode = node
				}
			}

			migratedNodeRaw, _ := json.Marshal(migratedNode)
			migratedNodeJSON, _ := utils.JSONMarshalPretty(migratedNodeRaw)
			expectedNodeJSON, _ := utils.JSONMarshalPretty(test.ExpectedNode)

			assert.Equal(t, string(expectedNodeJSON), string(migratedNodeJSON))

			checkFlowLocalization(t, migratedFlow, test.ExpectedLocalization)
		}
	}
}

func readLegacyTestFlows(flowsJSON string) ([]*legacy.Flow, error) {
	var legacyFlows []json.RawMessage
	json.Unmarshal(json.RawMessage(flowsJSON), &legacyFlows)
	return legacy.ReadLegacyFlows(legacyFlows)
}

func checkFlowLocalization(t *testing.T, flow flows.Flow, expectedLocalizationRaw json.RawMessage) {
	actualLocalizationJSON, _ := utils.JSONMarshalPretty(flow.Localization())
	expectedLocalizationJSON, _ := utils.JSONMarshalPretty(expectedLocalizationRaw)

	assert.Equal(t, string(expectedLocalizationJSON), string(actualLocalizationJSON))
}

func TestTranslations(t *testing.T) {
	translations := []map[utils.Language]string{
		{"eng": "Yes", "fra": "Oui"},
		{"eng": "No", "fra": "Non"},
		{"eng": "Maybe"},
		{"eng": "Never", "fra": "Jamas"},
	}
	assert.Equal(t, map[utils.Language][]string{
		"eng": {"Yes", "No", "Maybe", "Never"},
		"fra": {"Oui", "Non", "", "Jamas"},
	}, legacy.TransformTranslations(translations))
}
