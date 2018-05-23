package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/assets"
	"github.com/nyaruka/goflow/flows/engine"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/flows/triggers"
	"github.com/nyaruka/goflow/test"
	"github.com/nyaruka/goflow/utils"

	diff "github.com/sergi/go-diff/diffmatchpatch"
	"github.com/stretchr/testify/require"
)

var testServerPort = 49999

var flowTests = []struct {
	assets string
	output string
}{
	{"two_questions.json", "two_questions_test.json"},
	{"subflow.json", "subflow_test.json"},
	{"subflow_other.json", "subflow_other_test.json"},
	{"brochure.json", "brochure_test.json"},
	{"all_actions.json", "all_actions_test.json"},
	{"default_result.json", "default_result_test.json"},
	{"empty.json", "empty_test.json"},
	{"node_loop.json", "node_loop_test.json"},
	{"subflow_loop.json", "subflow_loop_test.json"},
	{"date_parse.json", "date_parse_test.json"},
	{"webhook_persists.json", "webhook_persists_test.json"},
	{"dynamic_groups.json", "dynamic_groups_test.json"},
	{"triggered.json", "triggered_test.json"},
	{"no_contact.json", "no_contact_test.json"},
	{"redact_urns.json", "redact_urns_test.json"},
}

var writeOutput bool
var serverURL = ""

func init() {
	flag.BoolVar(&writeOutput, "write", false, "whether to rewrite TestFlow output")
}

func deriveFilename(prefix string, filename string) string {
	if filename == "" {
		filename = "default.json"
	}

	if !strings.Contains(filename, "/") {
		filename = fmt.Sprintf("%s/%s%s", "testdata", prefix, filename)
	}
	return filename
}

func readFile(prefix string, filename string) ([]byte, error) {
	filename = deriveFilename(prefix, filename)
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("Error reading file '%s': %s", filename, err)
	}
	return bytes, err
}

type runResult struct {
	assetCache *assets.AssetCache
	session    flows.Session
	outputs    []*Output
}

func runFlow(assetsFilename string, triggerEnvelope *utils.TypedEnvelope, callerEvents [][]flows.Event) (runResult, error) {
	// load both the test specific assets and default assets
	defaultAssetsJSON, err := readFile("", "default.json")
	if err != nil {
		return runResult{}, err
	}
	testAssetsJSON, err := readFile("flows/", assetsFilename)
	if err != nil {
		return runResult{}, err
	}

	// rewrite the URL on any webhook actions
	testAssetsJSONStr := strings.Replace(string(testAssetsJSON), "http://localhost", serverURL, -1)

	assetCache := assets.NewAssetCache(100, 5)
	if err := assetCache.Include(defaultAssetsJSON); err != nil {
		return runResult{}, fmt.Errorf("Error reading default assets '%s': %s", assetsFilename, err)
	}
	if err := assetCache.Include(json.RawMessage(testAssetsJSONStr)); err != nil {
		return runResult{}, fmt.Errorf("Error reading test assets '%s': %s", assetsFilename, err)
	}

	session := engine.NewSession(assetCache, assets.NewMockAssetServer(), engine.NewDefaultConfig(), test.TestHTTPClient)

	trigger, err := triggers.ReadTrigger(session, triggerEnvelope)
	if err != nil {
		return runResult{}, fmt.Errorf("error unmarshalling trigger: %s", err)
	}

	err = session.Start(trigger, callerEvents[0])
	if err != nil {
		return runResult{}, err
	}

	outputs := make([]*Output, 0)

	// for each of our remaining caller events
	resumeEvents := callerEvents[1:]
	for i := range resumeEvents {
		sessionJSON, err := utils.JSONMarshalPretty(session)
		if err != nil {
			return runResult{}, fmt.Errorf("Error marshalling output: %s", err)
		}
		outputs = append(outputs, &Output{sessionJSON, marshalEventLog(session.Events())})

		session, err = engine.ReadSession(assetCache, assets.NewMockAssetServer(), engine.NewDefaultConfig(), test.TestHTTPClient, sessionJSON)
		if err != nil {
			return runResult{}, fmt.Errorf("Error marshalling output: %s", err)
		}

		// if we aren't at a wait, that's an error
		if session.Wait() == nil {
			return runResult{}, fmt.Errorf("Did not stop at expected wait, have unused resume events: %#v", resumeEvents[i:])
		}
		err = session.Resume(resumeEvents[i])
		if err != nil {
			return runResult{}, err
		}
	}

	sessionJSON, err := utils.JSONMarshalPretty(session)
	if err != nil {
		return runResult{}, fmt.Errorf("Error marshalling output: %s", err)
	}
	outputs = append(outputs, &Output{sessionJSON, marshalEventLog(session.Events())})

	return runResult{assetCache, session, outputs}, nil
}

func TestFlows(t *testing.T) {
	server, err := test.NewTestHTTPServer(testServerPort)
	require.NoError(t, err)

	defer server.Close()
	defer utils.SetUUIDGenerator(utils.DefaultUUIDGenerator)

	// save away our server URL so we can rewrite our URLs
	serverURL = server.URL

	for _, tc := range flowTests {
		utils.SetUUIDGenerator(utils.NewSeededUUID4Generator(123456))

		testJSON, err := readFile("flows/", tc.output)
		require.NoError(t, err, "Error reading output file for flow '%s' and output '%s': %s", tc.assets, tc.output, err)

		flowTest := FlowTest{}
		err = json.Unmarshal(json.RawMessage(testJSON), &flowTest)
		require.NoError(t, err, "Error unmarshalling output for flow '%s' and output '%s': %s", tc.assets, tc.output, err)

		// unmarshal our caller events
		callerEvents := make([][]flows.Event, len(flowTest.CallerEvents))

		for i := range flowTest.CallerEvents {
			callerEvents[i] = make([]flows.Event, len(flowTest.CallerEvents[i]))

			for e := range flowTest.CallerEvents[i] {
				event, err := events.EventFromEnvelope(flowTest.CallerEvents[i][e])
				require.NoError(t, err, "Error unmarshalling caller events for flow '%s' and output '%s': %s", tc.assets, tc.output, err)

				event.SetFromCaller(true)
				callerEvents[i][e] = event
			}
		}

		// run our flow
		runResult, err := runFlow(tc.assets, flowTest.Trigger, callerEvents)
		if err != nil {
			t.Errorf("Error running flow for flow '%s' and output '%s': %s", tc.assets, tc.output, err)
			continue
		}

		if writeOutput {
			// we are writing new outputs, we write new files but don't test anything
			rawOutputs := make([]json.RawMessage, len(runResult.outputs))
			for i := range runResult.outputs {
				rawOutputs[i], err = utils.JSONMarshal(runResult.outputs[i])
				require.NoError(t, err)
			}
			flowTest := FlowTest{Trigger: flowTest.Trigger, CallerEvents: flowTest.CallerEvents, Outputs: rawOutputs}
			testJSON, err := utils.JSONMarshalPretty(flowTest)
			require.NoError(t, err, "Error marshalling test definition: %s", err)

			// write our output
			outputFilename := deriveFilename("flows/", tc.output)
			err = ioutil.WriteFile(outputFilename, clearTimestamps(testJSON), 0644)
			require.NoError(t, err, "Error writing test file to %s: %s", outputFilename, err)
		} else {
			// start by checking we have the expected number of outputs
			if !assert.Equal(t, len(flowTest.Outputs), len(runResult.outputs), "wrong number of outputs for flow test %s", tc.assets) {
				continue
			}

			// then check each output
			for i, actual := range runResult.outputs {
				// unmarshal our expected outputsinto session+events
				expected := &Output{}
				err := json.Unmarshal(flowTest.Outputs[i], expected)
				require.NoError(t, err, "error unmarshalling output")

				// first the session
				if !assertEqualJSON(t, expected.Session, actual.Session, fmt.Sprintf("session is different in output[%d] for flow test %s", i, tc.assets)) {
					break
				}

				// and then each event
				for e := range actual.Events {
					if !assertEqualJSON(t, expected.Events[e], actual.Events[e], fmt.Sprintf("event[%d] is different in output[%d] for flow test %s", e, i, tc.assets)) {
						break
					}
				}
			}
		}
	}
}

// asserts that the given JSON fragments are equal
func assertEqualJSON(t *testing.T, expected json.RawMessage, actual json.RawMessage, message string) bool {
	expectedNormalized, _ := normalizeJSON(expected)
	actualNormalized, _ := normalizeJSON(actual)

	differ := diff.New()
	diffs := differ.DiffMain(string(expectedNormalized), string(actualNormalized), false)

	if len(diffs) != 1 || diffs[0].Type != diff.DiffEqual {
		assert.Fail(t, message, differ.DiffPrettyText(diffs))
		return false
	}
	return true
}

func normalizeJSON(data json.RawMessage) ([]byte, error) {
	data = clearTimestamps(data)

	var asMap map[string]interface{}
	if err := json.Unmarshal(data, &asMap); err != nil {
		return nil, err
	}

	return utils.JSONMarshalPretty(asMap)
}
