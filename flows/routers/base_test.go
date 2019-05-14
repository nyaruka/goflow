package routers_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/engine"
	"github.com/nyaruka/goflow/flows/routers"
	"github.com/nyaruka/goflow/flows/triggers"
	"github.com/nyaruka/goflow/test"
	"github.com/nyaruka/goflow/utils"

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

	server := test.NewTestHTTPServer(49993)

	for _, typeName := range routers.RegisteredTypes() {
		testRouterType(t, assetsJSON, typeName, server.URL)
	}
}

type inspectionResults struct {
	Templates    []string            `json:"templates"`
	Dependencies []string            `json:"dependencies"`
	Results      []*flows.ResultSpec `json:"results"`
}

func testRouterType(t *testing.T, assetsJSON json.RawMessage, typeName string, testServerURL string) {
	testFile, err := ioutil.ReadFile(fmt.Sprintf("testdata/%s.json", typeName))
	require.NoError(t, err)

	tests := []struct {
		Description     string             `json:"description"`
		Router          json.RawMessage    `json:"router"`
		ValidationError string             `json:"validation_error"`
		InspectionError string             `json:"inspection_error"`
		Results         json.RawMessage    `json:"results"`
		Events          []json.RawMessage  `json:"events"`
		Inspection      *inspectionResults `json:"inspection"`
	}{}

	err = json.Unmarshal(testFile, &tests)
	require.NoError(t, err)

	defer utils.SetTimeSource(utils.DefaultTimeSource)
	defer utils.SetUUIDGenerator(utils.DefaultUUIDGenerator)
	defer utils.SetRand(utils.DefaultRand)

	for _, tc := range tests {
		utils.SetTimeSource(utils.NewFixedTimeSource(time.Date(2018, 10, 18, 14, 20, 30, 123456, time.UTC)))
		utils.SetUUIDGenerator(utils.NewSeededUUID4Generator(12345))
		utils.SetRand(utils.NewSeededRand(123456))

		testName := fmt.Sprintf("test '%s' for router type '%s'", tc.Description, typeName)

		// inject the router into a suitable node in our test flow
		routerPath := []string{"flows", "[0]", "nodes", "[0]", "router"}
		assetsJSON = test.JSONReplace(assetsJSON, routerPath, tc.Router)

		// create session assets
		sa, err := test.CreateSessionAssets(assetsJSON, testServerURL)
		require.NoError(t, err)

		// now try to read the flow, and if we expect a validation error, check that
		flow, err := sa.Flows().Get("16f6eee7-9843-4333-bad2-1d7fd636452c")
		if tc.ValidationError != "" {
			rootErr := errors.Cause(err)
			assert.EqualError(t, rootErr, tc.ValidationError, "read error mismatch in %s", testName)
			continue
		} else {
			assert.NoError(t, err, "unexpected read error in %s", testName)
		}

		// if this router is expected to return a inspection failure, check that
		err = flow.Inspect(sa)
		if tc.InspectionError != "" {
			rootErr := errors.Cause(err)
			assert.EqualError(t, rootErr, tc.InspectionError, "check error mismatch in %s", testName)
			continue
		} else {
			assert.NoError(t, err, "unexpected check error in %s", testName)
		}

		// load our contact
		contact, err := flows.ReadContact(sa, json.RawMessage(contactJSON), assets.PanicOnMissing)
		require.NoError(t, err)

		trigger := triggers.NewManualTrigger(utils.NewEnvironmentBuilder().Build(), flow.Reference(), contact, nil)

		eng := engine.NewBuilder().WithDefaultUserAgent("goflow-testing").Build()
		session, _, err := eng.NewSession(sa, trigger)
		require.NoError(t, err)

		// check results are what we expected
		run := session.Runs()[0]
		actualResultsJSON, _ := json.Marshal(run.Results())
		expectedResultsJSON, _ := json.Marshal(tc.Results)
		test.AssertEqualJSON(t, expectedResultsJSON, actualResultsJSON, "results mismatch in %s", testName)

		// check events are what we expected
		actualEventsJSON, _ := json.Marshal(run.Events())
		expectedEventsJSON, _ := json.Marshal(tc.Events)
		test.AssertEqualJSON(t, expectedEventsJSON, actualEventsJSON, "events mismatch in %s", testName)

		// try marshaling the router back to JSON
		routerJSON, err := json.Marshal(flow.Nodes()[0].Router())
		test.AssertEqualJSON(t, tc.Router, routerJSON, "marshal mismatch in %s", testName)

		// finally try inspecting this router
		templates := flow.ExtractTemplates()
		assert.Equal(t, tc.Inspection.Templates, templates, "inspected templates mismatch in %s", testName)

		dependencies := flow.ExtractDependencies()
		depStrings := make([]string, len(dependencies))
		for i := range dependencies {
			depStrings[i] = dependencies[i].String()
		}
		assert.Equal(t, tc.Inspection.Dependencies, depStrings, "inspected dependencies mismatch in %s", testName)

		results := flow.ExtractResults()
		assert.Equal(t, tc.Inspection.Results, results, "inspected results mismatch in %s", testName)
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
