package main

import (
	"encoding/json"
	"flag"
	"fmt"
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
		outJSON, err := utils.JSONMarshalPretty(session)
		if err != nil {
			return runResult{}, fmt.Errorf("Error marshalling output: %s", err)
		}
		outputs = append(outputs, &Output{outJSON, marshalEventLog(session.Events())})

		session, err = engine.ReadSession(assetCache, assets.NewMockAssetServer(), engine.NewDefaultConfig(), test.TestHTTPClient, outJSON)
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

	outJSON, err := utils.JSONMarshalPretty(session)
	if err != nil {
		return runResult{}, fmt.Errorf("Error marshalling output: %s", err)
	}
	outputs = append(outputs, &Output{outJSON, marshalEventLog(session.Events())})

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
			// unmarshal our expected outputs
			expectedOutputs := make([]*Output, len(flowTest.Outputs))
			for i := range expectedOutputs {
				output := &Output{}
				err := json.Unmarshal(flowTest.Outputs[i], output)
				require.NoError(t, err, "Error unmarshalling output: %s", err)

				expectedOutputs[i] = output
			}

			// read our output and test that we are the same
			if len(runResult.outputs) != len(expectedOutputs) {
				t.Errorf("Actual outputs:\n%s\n do not match expected:\n%s\n for flow '%s'", runResult.outputs, expectedOutputs, tc.assets)
				continue
			}

			for i := range runResult.outputs {
				actualOutput := runResult.outputs[i]
				expectedOutput := expectedOutputs[i]

				actualSession, err := engine.ReadSession(runResult.assetCache, assets.NewMockAssetServer(), engine.NewDefaultConfig(), test.TestHTTPClient, actualOutput.Session)
				require.NoError(t, err, "Error unmarshalling session running flow '%s': %s\n", tc.assets, err)

				expectedSession, err := engine.ReadSession(runResult.assetCache, assets.NewMockAssetServer(), engine.NewDefaultConfig(), test.TestHTTPClient, expectedOutput.Session)
				require.NoError(t, err, "Error unmarshalling expected session running flow '%s': %s\n", tc.assets, err)

				// number of runs should be the same
				if len(actualSession.Runs()) != len(expectedSession.Runs()) {
					t.Errorf("Actual runs:\n%#v\n do not match expected:\n%#v\n for flow '%s'\n", actualSession.Runs(), expectedSession.Runs(), tc.assets)
				}

				// runs should have same status and flows
				for i := range actualSession.Runs() {
					run := actualSession.Runs()[i]
					expected := expectedSession.Runs()[i]

					if run.Flow() != expected.Flow() {
						t.Errorf("Actual run flow: %s does not match expected: %s for flow '%s'", run.Flow().UUID(), expected.Flow().UUID(), tc.assets)
					}

					if run.Status() != expected.Status() {
						t.Errorf("Actual run status: %s does not match expected: %s for flow '%s'", run.Status(), expected.Status(), tc.assets)
					}
				}

				if len(actualOutput.Events) != len(expectedOutput.Events) {
					t.Errorf("Actual events:\n%#v\n do not match expected:\n%#v\n for flow '%s'\n", actualOutput.Events, expectedOutput.Events, tc.assets)
				}

				for j := range actualOutput.Events {
					event := actualOutput.Events[j]
					expected := expectedOutput.Events[j]

					// write our events as json
					eventJSON, err := rawMessageAsJSON(event)
					require.NoError(t, err, "Error marshalling event for flow '%s' and output '%s': %s\n", tc.assets, tc.output, err)

					expectedJSON, err := rawMessageAsJSON(expected)
					require.NoError(t, err, "Error marshalling expected event for flow '%s' and output '%s': %s\n", tc.assets, tc.output, err)

					if eventJSON != expectedJSON {
						t.Errorf("Got event:\n'%s'\n\nwhen expecting:\n'%s'\n\n for flow '%s' and output '%s\n", eventJSON, expectedJSON, tc.assets, tc.output)
						break
					}
				}
			}
		}
	}
}
