package definition

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"regexp"

	"strings"

	"github.com/buger/jsonparser"
	"github.com/nyaruka/goflow/flows/routers"
	"github.com/nyaruka/goflow/utils"
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
		legacyFlowsJSON := fmt.Sprintf(legacyActionHolderDef, string(test.LegacyAction))
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

		checkFlowLocalization(t, &migratedFlow, test.ExpectedLocalization, migratedActionRaw)
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
		legacyFlowsJSON := fmt.Sprintf(legacyTestHolderDef, string(test.LegacyTest))
		legacyFlows, err := ReadLegacyFlows(json.RawMessage(legacyFlowsJSON))
		if err != nil {
			t.Fatal(err)
		}

		migratedFlow := legacyFlows[0]
		migratedRouter := migratedFlow.Nodes()[0].Router().(*routers.SwitchRouter)

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

			checkFlowLocalization(t, &migratedFlow, test.ExpectedLocalization, migratedCaseRaw)
		}
	}
}

func TestRuleSetMigration(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/ruleset_migrations.json")
	if err != nil {
		t.Fatal(err)
	}

	var tests []RuleSetMigrationTest
	err = json.Unmarshal(data, &tests)
	if err != nil {
		t.Fatal(err)
	}

	for _, test := range tests {
		legacyFlowsJSON := fmt.Sprintf(legacyRuleSetHolderDef, string(test.LegacyRuleSet))
		legacyFlows, err := ReadLegacyFlows(json.RawMessage(legacyFlowsJSON))
		if err != nil {
			t.Fatal(err)
		}

		migratedFlow := legacyFlows[0]

		if len(migratedFlow.Nodes()) == 0 {
			t.Errorf("Got no migrated nodes from legacy ruleset:\n%s\n\n", string(test.LegacyRuleSet))
		} else {
			migratedNode := migratedFlow.Nodes()[0]
			migratedNodeRaw, _ := json.Marshal(migratedNode)
			migratedNodeJSON := formatJSON(migratedNodeRaw)
			expectedNodeJSON := formatJSON(test.ExpectedNode)

			if !wildcardEquals(migratedNodeJSON, expectedNodeJSON) {
				t.Errorf("Got node:\n%s\n\nwhen expecting:\n%s\n\n", migratedNodeJSON, expectedNodeJSON)
			}

			checkFlowLocalization(t, &migratedFlow, test.ExpectedLocalization, migratedNodeRaw)
		}
	}
}

func checkFlowLocalization(t *testing.T, flow *LegacyFlow, expectedLocalizationRaw json.RawMessage, substitutionSource json.RawMessage) {
	actualLocalization := *flow.translations.(*flowTranslations)
	actualLocalizationRaw, _ := json.Marshal(actualLocalization)
	actualLocalizationJSON := formatJSON(actualLocalizationRaw)

	// perfrom substitutions on expected localization
	regex := regexp.MustCompile(`{{[^}]+}}`)
	expectedLocalizationStr := regex.ReplaceAllStringFunc(string(expectedLocalizationRaw), func(match string) string {
		queryKeys := strings.Split(strings.TrimSpace(match[2:len(match)-2]), ".")

		subValue, err := jsonparser.GetString(substitutionSource, queryKeys...)
		if err != nil {
			t.Fatalf("Couldn't find: '%s' in\n\n%s\n", strings.Join(queryKeys, "."), string(substitutionSource))
		}

		return subValue
	})

	// unmarshal and re-marchal expected JSON to get ordering correct after substitutions
	expectedLocalization := &flowTranslations{}
	json.Unmarshal(json.RawMessage(expectedLocalizationStr), expectedLocalization)
	expectedLocalizationRaw, _ = json.Marshal(expectedLocalization)
	expectedLocalizationJSON := formatJSON(expectedLocalizationRaw)

	if !wildcardEquals(actualLocalizationJSON, expectedLocalizationJSON) {
		t.Errorf("Got localization:\n%s\n\nwhen expecting:\n%s\n\n", actualLocalizationJSON, expectedLocalizationJSON)
	}
}

func formatJSON(data json.RawMessage) string {
	out, _ := json.MarshalIndent(data, "", "    ")
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
