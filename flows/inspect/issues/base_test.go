package issues_test

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/definition"
	"github.com/nyaruka/goflow/flows/inspect/issues"
	"github.com/nyaruka/goflow/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIssueTypes(t *testing.T) {
	env := envs.NewBuilder().Build()

	assets, err := test.LoadSessionAssets(env, "testdata/_assets.json")
	require.NoError(t, err)

	for typeName := range issues.RegisteredTypes {
		testIssueType(t, assets, typeName)
	}
}

func testIssueType(t *testing.T, sa flows.SessionAssets, typeName string) {
	testPath := fmt.Sprintf("testdata/%s.json", typeName)
	testFile, err := os.ReadFile(testPath)
	require.NoError(t, err)

	tests := []struct {
		Description string          `json:"description"`
		NoAssets    bool            `json:"no_assets,omitempty"`
		Flow        json.RawMessage `json:"flow"`

		Issues json.RawMessage `json:"issues"`
	}{}

	err = jsonx.Unmarshal(testFile, &tests)
	require.NoError(t, err)

	for i, tc := range tests {
		testName := fmt.Sprintf("test '%s' for modifier type '%s'", tc.Description, typeName)

		// read the flow to be checked
		flow, err := definition.ReadFlow(tc.Flow, nil)
		require.NoError(t, err, "error reading flow in %s", testName)

		var sessionAssets flows.SessionAssets
		if !tc.NoAssets {
			sessionAssets = sa
		}

		info := flow.Inspect(sessionAssets)
		issuesJSON := jsonx.MustMarshal(info.Issues)

		// clone test case and populate with actual values
		actual := tc
		actual.Issues = issuesJSON

		if !test.UpdateSnapshots {
			// check the found issues
			test.AssertEqualJSON(t, tc.Issues, actual.Issues, "issues mismatch in %s", testName)
		} else {
			tests[i] = actual
		}
	}

	if test.UpdateSnapshots {
		actualJSON, err := jsonx.MarshalPretty(tests)
		require.NoError(t, err)

		err = os.WriteFile(testPath, actualJSON, 0666)
		require.NoError(t, err)
	}
}

func TestIssues(t *testing.T) {
	env := envs.NewBuilder().Build()

	sa, err := test.LoadSessionAssets(env, "testdata/_assets.json")
	require.NoError(t, err)

	flow, err := definition.ReadFlow([]byte(`{
		"uuid": "76f0a02f-3b75-4b86-9064-e9195e1b3a02",
		"name": "Test Flow",
		"spec_version": "13.0",
		"language": "eng",
		"type": "messaging",
		"nodes": [
			{
				"uuid": "a58be63b-907d-4a1a-856b-0bb5579d7507",
				"actions": [
					{
						"uuid": "f01d693b-2af2-49fb-9e38-146eb00937e9",
						"type": "send_msg",
						"text": "You live in @fields.county and are @fields.age"
					}
				],
				"exits": [
					{
						"uuid": "118221f7-e637-4cdb-83ca-7f0a5aae98c6"
					}
				]
			}
		]
	}`), nil)
	require.NoError(t, err)

	info := flow.Inspect(sa)

	assert.Equal(t, 1, len(info.Issues))
	assert.Equal(t, issues.TypeMissingDependency, info.Issues[0].Type())
	assert.Equal(t, flows.NodeUUID("a58be63b-907d-4a1a-856b-0bb5579d7507"), info.Issues[0].NodeUUID())
	assert.Equal(t, flows.ActionUUID("f01d693b-2af2-49fb-9e38-146eb00937e9"), info.Issues[0].ActionUUID())
	assert.Equal(t, envs.NilLanguage, info.Issues[0].Language())
	assert.Equal(t, "missing field dependency 'county'", info.Issues[0].Description())
}
