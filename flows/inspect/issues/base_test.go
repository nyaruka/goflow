package issues_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/definition"
	"github.com/nyaruka/goflow/flows/inspect/issues"
	"github.com/nyaruka/goflow/test"
	"github.com/nyaruka/goflow/utils/jsonx"
	"github.com/stretchr/testify/require"
)

func TestIssueTypes(t *testing.T) {
	assets, err := test.LoadSessionAssets("testdata/_assets.json")
	require.NoError(t, err)

	for typeName := range issues.RegisteredTypes {
		testIssueType(t, assets, typeName)
	}
}

func testIssueType(t *testing.T, sa flows.SessionAssets, typeName string) {
	testPath := fmt.Sprintf("testdata/%s.json", typeName)
	testFile, err := ioutil.ReadFile(testPath)
	require.NoError(t, err)

	tests := []struct {
		Description string          `json:"description"`
		NoAssets    bool            `json:"no_assets,omitempty"`
		Flow        json.RawMessage `json:"flow"`

		Issues json.RawMessage `json:"issues"`
	}{}

	err = json.Unmarshal(testFile, &tests)
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
		issuesJSON, _ := json.Marshal(info.Issues)

		// clone test case and populate with actual values
		actual := tc
		actual.Issues = issuesJSON

		if !test.WriteOutput {
			// check the found issues
			test.AssertEqualJSON(t, tc.Issues, actual.Issues, "issues mismatch in %s", testName)
		} else {
			tests[i] = actual
		}
	}

	if test.WriteOutput {
		actualJSON, err := jsonx.MarshalPretty(tests)
		require.NoError(t, err)

		err = ioutil.WriteFile(testPath, actualJSON, 0666)
		require.NoError(t, err)
	}
}
