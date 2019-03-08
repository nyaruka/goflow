package routers_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
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

	for typeName := range routers.RegisteredTypes() {
		testRouterType(t, assetsJSON, typeName, server.URL)
	}
}

type inspectionResults struct {
	Templates    []string `json:"templates"`
	Dependencies []string `json:"dependencies"`
	ResultNames  []string `json:"result_names"`
}

func testRouterType(t *testing.T, assetsJSON json.RawMessage, typeName string, testServerURL string) {
	testFile, err := ioutil.ReadFile(fmt.Sprintf("testdata/%s.json", typeName))
	require.NoError(t, err)

	tests := []struct {
		Description     string             `json:"description"`
		Router          json.RawMessage    `json:"router"`
		ValidationError string             `json:"validation_error"`
		Results         json.RawMessage    `json:"results"`
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

		testName := fmt.Sprintf("test '%s' for action type '%s'", tc.Description, typeName)

		// create unstarted session from our assets
		session, err := test.CreateSession(assetsJSON, testServerURL)
		require.NoError(t, err)

		// read the router to be tested
		router, err := routers.ReadRouter(tc.Router)
		require.NoError(t, err, "error loading router in %s", testName)
		assert.Equal(t, typeName, router.Type())

		// get a suitable "holder" flow
		flow, err := session.Assets().Flows().Get("16f6eee7-9843-4333-bad2-1d7fd636452c")
		require.NoError(t, err)

		// if not, add it to our flow
		flow.Nodes()[0].SetRouter(router)

		// if this router is expected to cause flow validation failure, check that
		err = flow.Validate(session.Assets())
		if tc.ValidationError != "" {
			rootErr := errors.Cause(err)
			assert.EqualError(t, rootErr, tc.ValidationError, "validation error mismatch in %s", testName)
			continue
		} else {
			assert.NoError(t, err, "unexpected validation error in %s", testName)
		}

		// load our contact
		contact, err := flows.ReadContact(session.Assets(), json.RawMessage(contactJSON), assets.PanicOnMissing)
		require.NoError(t, err)

		trigger := triggers.NewManualTrigger(utils.NewEnvironmentBuilder().Build(), flow.Reference(), contact, nil)
		_, err = session.Start(trigger, nil)
		require.NoError(t, err)

		// check results are what we expected
		run := session.Runs()[0]
		actualResultsJSON, _ := json.Marshal(run.Results())
		expectedResultsJSON, _ := json.Marshal(tc.Results)
		test.AssertEqualJSON(t, expectedResultsJSON, actualResultsJSON, "results mismatch in %s", testName)

		// try marshaling the router back to JSON
		routerJSON, err := json.Marshal(router)
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

		resultNames := flow.ExtractResultNames()
		assert.Equal(t, tc.Inspection.ResultNames, resultNames, "inspected result names mismatch in %s", testName)
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
