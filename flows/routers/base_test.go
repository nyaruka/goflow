package routers_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	"github.com/nyaruka/gocommon/dates"
	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/gocommon/random"
	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/routers"
	"github.com/nyaruka/goflow/flows/triggers"
	"github.com/nyaruka/goflow/test"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var contactJSON = `{
	"uuid": "5d76d86b-3bb9-4d5a-b822-c9d86f5d8e4f",
	"name": "Ryan Lewis",
	"language": "eng",
	"timezone": "America/Guayaquil",
	"urns": [],
	"fields": {
		"gender": {
			"text": "Male"
		}
	},
	"created_on": "2018-06-20T11:40:30.123456789-00:00"
}`

func TestRouterTypes(t *testing.T) {
	assetsJSON, err := ioutil.ReadFile("testdata/_assets.json")
	require.NoError(t, err)

	for _, typeName := range routers.RegisteredTypes() {
		testRouterType(t, assetsJSON, typeName)
	}
}

func testRouterType(t *testing.T, assetsJSON json.RawMessage, typeName string) {
	testPath := fmt.Sprintf("testdata/%s.json", typeName)
	testFile, err := ioutil.ReadFile(testPath)
	require.NoError(t, err)

	tests := []struct {
		Description string          `json:"description"`
		Router      json.RawMessage `json:"router"`

		ReadError         string          `json:"read_error,omitempty"`
		DependenciesError string          `json:"dependencies_error,omitempty"`
		Results           json.RawMessage `json:"results,omitempty"`
		Events            json.RawMessage `json:"events,omitempty"`
		Templates         []string        `json:"templates,omitempty"`
		LocalizedText     []string        `json:"localizables,omitempty"`
		Inspection        json.RawMessage `json:"inspection,omitempty"`
	}{}

	err = jsonx.Unmarshal(testFile, &tests)
	require.NoError(t, err)

	defer dates.SetNowSource(dates.DefaultNowSource)
	defer uuids.SetGenerator(uuids.DefaultGenerator)
	defer random.SetGenerator(random.DefaultGenerator)

	for i, tc := range tests {
		dates.SetNowSource(dates.NewFixedNowSource(time.Date(2018, 10, 18, 14, 20, 30, 123456, time.UTC)))
		uuids.SetGenerator(uuids.NewSeededGenerator(12345))
		random.SetGenerator(random.NewSeededGenerator(123456))

		testName := fmt.Sprintf("test '%s' for router type '%s'", tc.Description, typeName)

		// inject the router into a suitable node in our test flow
		routerPath := []string{"flows", "[0]", "nodes", "[0]", "router"}
		assetsJSON = test.JSONReplace(assetsJSON, routerPath, tc.Router)

		// create session assets
		sa, err := test.CreateSessionAssets(assetsJSON, "")
		require.NoError(t, err)

		// now try to read the flow, and if we expect a read error, check that
		flow, err := sa.Flows().Get("16f6eee7-9843-4333-bad2-1d7fd636452c")
		if tc.ReadError != "" {
			rootErr := errors.Cause(err)
			assert.EqualError(t, rootErr, tc.ReadError, "read error mismatch in %s", testName)
			continue
		} else {
			assert.NoError(t, err, "unexpected read error in %s", testName)
		}

		// load our contact
		contact, err := flows.ReadContact(sa, json.RawMessage(contactJSON), assets.PanicOnMissing)
		require.NoError(t, err)

		trigger := triggers.NewBuilder(envs.NewBuilder().Build(), flow.Reference(), contact).Manual().Build()

		eng := test.NewEngine()
		session, _, err := eng.NewSession(sa, trigger)
		require.NoError(t, err)

		// clone test case and populate with actual values
		actual := tc

		// re-marshal the action
		actual.Router, err = jsonx.Marshal(flow.Nodes()[0].Router())
		require.NoError(t, err)

		run := session.Runs()[0]
		actual.Results, _ = jsonx.Marshal(run.Results())
		actual.Events, _ = jsonx.Marshal(run.Events())

		if tc.Templates != nil {
			actual.Templates = flow.ExtractTemplates()
		}
		if tc.LocalizedText != nil {
			actual.LocalizedText = flow.ExtractLocalizables()
		}
		if tc.Inspection != nil {
			actual.Inspection, _ = jsonx.Marshal(flow.Inspect(sa))
		}

		if !test.UpdateSnapshots {
			test.AssertEqualJSON(t, tc.Router, actual.Router, "marshal mismatch in %s", testName)

			// check results are what we expected
			test.AssertEqualJSON(t, tc.Results, actual.Results, "results mismatch in %s", testName)

			// check events are what we expected
			test.AssertEqualJSON(t, tc.Events, actual.Events, "events mismatch in %s", testName)

			// check extracted templates
			assert.Equal(t, tc.Templates, actual.Templates, "extracted templates mismatch in %s", testName)

			// check extracted localized text
			assert.Equal(t, tc.LocalizedText, actual.LocalizedText, "extracted localized text mismatch in %s", testName)

			// check inspection results
			test.AssertEqualJSON(t, tc.Inspection, actual.Inspection, "inspection mismatch in %s", testName)
		} else {
			tests[i] = actual
		}
	}

	if test.UpdateSnapshots {
		actualJSON, err := jsonx.MarshalPretty(tests)
		require.NoError(t, err)

		err = ioutil.WriteFile(testPath, actualJSON, 0666)
		require.NoError(t, err)
	}
}

func TestReadRouter(t *testing.T) {
	// error if no type field
	_, err := routers.ReadRouter([]byte(`{"foo": "bar"}`))
	assert.EqualError(t, err, "field 'type' is required")

	// error if we don't recognize action type
	_, err = routers.ReadRouter([]byte(`{"type": "do_the_foo", "foo": "bar"}`))
	assert.EqualError(t, err, "unknown type: 'do_the_foo'")
}
