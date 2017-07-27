package definition

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/nyaruka/goflow/flows/routers"
	"github.com/nyaruka/goflow/utils"
)

var legacyFlowDef = `
[
	{
		"base_language": "eng",
		"rule_sets": [{
			"uuid": "2bff5c33-9d29-4cfc-8bb7-0a1b9f97d830",
			"rules": [{
				"test": %s, 
				"category": {"eng": "All Responses"}, 
				"destination": "10e483a8-5ffb-4c4f-917b-d43ce86c1d65", 
				"uuid": "c072ecb5-0686-40ea-8ed3-898dc1349783", 
				"destination_type": "A"
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
		"action_sets": [{
			"uuid": "10e483a8-5ffb-4c4f-917b-d43ce86c1d65",
			"y": 100, 
            "x": 100, 
            "destination": null, 
			"actions": [%s]
		}],
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

func TestActionMigration(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/action_migrations.json")
	if err != nil {
		t.Fatal(err)
	}

	var tests []ActionMigrationTest
	err = json.Unmarshal(data, &tests)
	if err != nil {
		t.Fatal(err)
	}

	for _, test := range tests {
		legacyFlowsJSON := fmt.Sprintf(legacyFlowDef, `{"type": "true"}`, string(test.LegacyAction))
		legacyFlows, err := ReadLegacyFlows(json.RawMessage(legacyFlowsJSON))
		if err != nil {
			t.Fatal(err)
		}

		migratedFlow := legacyFlows[0]
		migratedAction := migratedFlow.Nodes()[0].Actions()[0]
		migratedActionEnvelope, _ := utils.EnvelopeFromTyped(migratedAction)
		migratedActionRaw, _ := json.Marshal(migratedActionEnvelope)
		migratedActionJSON := formatJSON(migratedActionRaw)
		expectedActionJSON := formatJSON(test.ExpectedAction)

		if !wildcardEquals(migratedActionJSON, expectedActionJSON) {
			t.Errorf("Got action:\n%s\n\nwhen expecting:\n%s\n\n", migratedActionJSON, expectedActionJSON)
		}

		checkFlowLocalization(t, &migratedFlow, test.ExpectedLocalization)
	}
}

func TestTestMigration(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/test_migrations.json")
	if err != nil {
		t.Fatal(err)
	}

	var tests []TestMigrationTest
	err = json.Unmarshal(data, &tests)
	if err != nil {
		t.Fatal(err)
	}

	for _, test := range tests {
		legacyFlowsJSON := fmt.Sprintf(legacyFlowDef, string(test.LegacyTest), "")
		legacyFlows, err := ReadLegacyFlows(json.RawMessage(legacyFlowsJSON))
		if err != nil {
			t.Fatal(err)
		}

		migratedFlow := legacyFlows[0]
		migratedRouter := migratedFlow.Nodes()[1].Router().(*routers.SwitchRouter)

		if len(migratedRouter.Cases) == 0 {
			t.Errorf("Got no migrated case from legacy test:\n%s\n\n", string(test.LegacyTest))
		} else {
			migratedCase := migratedRouter.Cases[0]
			migratedCaseRaw, _ := json.Marshal(migratedCase)
			migratedCaseJSON := formatJSON(migratedCaseRaw)
			expectedCaseJSON := formatJSON(test.ExpectedCase)

			if !wildcardEquals(migratedCaseJSON, expectedCaseJSON) {
				t.Errorf("Got case:\n%s\n\nwhen expecting:\n%s\n\n", migratedCaseJSON, expectedCaseJSON)
			}

			checkFlowLocalization(t, &migratedFlow, test.ExpectedLocalization)
		}
	}
}

func checkFlowLocalization(t *testing.T, flow *LegacyFlow, expectedLocalization json.RawMessage) {
	actualLocalization := *flow.translations.(*flowTranslations)
	actualLocalizationRaw, _ := json.Marshal(actualLocalization)
	actualLocalizationJSON := formatJSON(actualLocalizationRaw)
	expectedLocalizationJSON := formatJSON(expectedLocalization)

	if !wildcardEquals(actualLocalizationJSON, expectedLocalizationJSON) {
		t.Errorf("Got localization:\n%s\n\nwhen expecting:\n%s\n\n", actualLocalizationJSON, expectedLocalizationJSON)
	}
}

func formatJSON(data json.RawMessage) string {
	out, _ := json.MarshalIndent(data, "", "  ")
	return string(out)
}

// checks if two strings are, ignoring any � characters
func wildcardEquals(actual string, expected string) bool {
	actualRunes := []rune(actual)
	expectedRunes := []rune(expected)
	substituted := make([]rune, len(expectedRunes))
	for c, ch := range expectedRunes {
		if ch == '�' && c < len(actualRunes) {
			substituted[c] = actualRunes[c]
		} else {
			substituted[c] = expectedRunes[c]
		}
	}
	return string(actualRunes) == string(substituted)
}
