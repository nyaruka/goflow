package definition

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/nyaruka/goflow/utils"
)

var legacyActionHolderDef = `
[
	{
		"base_language": "eng",
		"action_sets": [{
			"uuid": "72a1f5df-49f9-45df-94c9-d86f7ea064e5",
			"y": 0, 
            "x": 100, 
            "destination": "2bff5c33-9d29-4cfc-8bb7-0a1b9f97d830", 
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

func formatJSON(data json.RawMessage) string {
	out, _ := json.MarshalIndent(data, "", "  ")
	return string(out)
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

		if migratedActionJSON != expectedActionJSON {
			t.Errorf("Got action:\n%s\n\nwhen expecting:\n%s\n\n", migratedActionJSON, expectedActionJSON)
		}

		migratedLocalization := *migratedFlow.translations.(*flowTranslations)
		migratedLocalizationRaw, _ := json.Marshal(migratedLocalization)
		migratedLocalizationJSON := formatJSON(migratedLocalizationRaw)
		expectedLocalizationJSON := formatJSON(test.ExpectedLocalization)

		if migratedLocalizationJSON != expectedLocalizationJSON {
			t.Errorf("Got localization:\n%s\n\nwhen expecting:\n%s\n\n", migratedLocalizationJSON, expectedLocalizationJSON)
		}

	}
}
