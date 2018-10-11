package test

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/nyaruka/goflow/assets/static"
	"io/ioutil"
	"strings"
	"testing"
	"time"

	_ "github.com/nyaruka/goflow/extensions/transferto"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/engine"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/flows/resumes"
	"github.com/nyaruka/goflow/flows/triggers"
	"github.com/nyaruka/goflow/utils"

	diff "github.com/sergi/go-diff/diffmatchpatch"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var flowTests = []struct {
	assets string
	output string
}{
	{"airtime.json", "airtime_disabled_test.json"},
	{"all_actions.json", "all_actions_test.json"},
	{"brochure.json", "brochure_test.json"},
	{"date_parse.json", "date_parse_test.json"},
	{"default_result.json", "default_result_test.json"},
	{"dynamic_groups_correction.json", "dynamic_groups_correction_test.json"},
	{"dynamic_groups.json", "dynamic_groups_test.json"},
	{"empty.json", "empty_test.json"},
	{"initial_wait.json", "initial_wait_test.json"},
	{"legacy_extra.json", "legacy_extra_test.json"},
	{"no_contact.json", "no_contact_test.json"},
	{"node_loop.json", "node_loop_test.json"},
	{"redact_urns.json", "redact_urns_test.json"},
	{"resthook.json", "resthook_test.json"},
	{"router_tests.json", "router_tests_test.json"},
	{"subflow_loop.json", "subflow_loop_test.json"},
	{"subflow_other.json", "subflow_other_test.json"},
	{"subflow.json", "subflow_test.json"},
	{"triggered.json", "triggered_test.json"},
	{"two_questions.json", "two_questions_test.json"},
	{"webhook_migrated.json", "webhook_migrated_test.json"},
	{"webhook_persists.json", "webhook_persists_test.json"},
}

var writeOutput bool
var serverURL = ""

func init() {
	flag.BoolVar(&writeOutput, "write", false, "whether to rewrite test output")
}

func normalizeJSON(data json.RawMessage) ([]byte, error) {
	var asMap map[string]interface{}
	if err := json.Unmarshal(data, &asMap); err != nil {
		return nil, err
	}

	return utils.JSONMarshalPretty(asMap)
}

func marshalEventLog(eventLog []flows.Event) ([]json.RawMessage, error) {
	envelopes, err := events.EventsToEnvelopes(eventLog)
	marshaled := make([]json.RawMessage, len(envelopes))

	for i := range envelopes {
		marshaled[i], err = utils.JSONMarshal(envelopes[i])
		if err != nil {
			return nil, fmt.Errorf("error creating marshaling envelope %s: %s", envelopes[i], err)
		}
	}
	return marshaled, nil
}

type Output struct {
	Session json.RawMessage   `json:"session"`
	Events  []json.RawMessage `json:"events"`
}

type FlowTest struct {
	Trigger *utils.TypedEnvelope   `json:"trigger"`
	Resumes []*utils.TypedEnvelope `json:"resumes"`
	Outputs []json.RawMessage      `json:"outputs"`
}

type runResult struct {
	session flows.Session
	outputs []*Output
}

func runFlow(assetsPath string, triggerEnvelope *utils.TypedEnvelope, resumeEnvelopes []*utils.TypedEnvelope) (runResult, error) {
	// load the test specific assets
	testAssetsJSON, err := ioutil.ReadFile(fmt.Sprintf("testdata/flows/%s", assetsPath))
	if err != nil {
		return runResult{}, err
	}

	// rewrite the URL on any webhook actions
	testAssetsJSONStr := strings.Replace(string(testAssetsJSON), "http://localhost", serverURL, -1)

	source, err := static.NewStaticSource(json.RawMessage(testAssetsJSONStr))
	if err != nil {
		return runResult{}, fmt.Errorf("error reading test assets '%s': %s", assetsPath, err)
	}

	assets, _ := engine.NewSessionAssets(source)
	session := engine.NewSession(assets, engine.NewDefaultConfig(), TestHTTPClient)

	trigger, err := triggers.ReadTrigger(session, triggerEnvelope)
	if err != nil {
		return runResult{}, fmt.Errorf("error unmarshalling trigger: %s", err)
	}

	err = session.Start(trigger)
	if err != nil {
		return runResult{}, err
	}

	outputs := make([]*Output, 0)

	// try to resume the session for each of the provided resumes
	for r, resumeEnvelope := range resumeEnvelopes {
		sessionJSON, err := utils.JSONMarshalPretty(session)
		if err != nil {
			return runResult{}, fmt.Errorf("Error marshalling output: %s", err)
		}
		marshalledEvents, err := marshalEventLog(session.Events())
		if err != nil {
			return runResult{}, err
		}

		outputs = append(outputs, &Output{sessionJSON, marshalledEvents})

		session, err = engine.ReadSession(assets, engine.NewDefaultConfig(), TestHTTPClient, sessionJSON)
		if err != nil {
			return runResult{}, fmt.Errorf("Error marshalling output: %s", err)
		}

		// if we aren't at a wait, that's an error
		if session.Wait() == nil {
			return runResult{}, fmt.Errorf("Did not stop at expected wait, have unused resumes: %#v", resumeEnvelopes[r:])
		}

		resume, err := resumes.ReadResume(session, resumeEnvelope)
		if err != nil {
			return runResult{}, err
		}

		err = session.Resume(resume)
		if err != nil {
			return runResult{}, err
		}
	}

	sessionJSON, err := utils.JSONMarshalPretty(session)
	if err != nil {
		return runResult{}, fmt.Errorf("Error marshalling output: %s", err)
	}

	marshalledEvents, err := marshalEventLog(session.Events())
	if err != nil {
		return runResult{}, err
	}

	outputs = append(outputs, &Output{sessionJSON, marshalledEvents})

	return runResult{session, outputs}, nil
}

func TestFlows(t *testing.T) {
	server, err := NewTestHTTPServer(49999)
	require.NoError(t, err)

	defer server.Close()
	defer utils.SetUUIDGenerator(utils.DefaultUUIDGenerator)
	defer utils.SetTimeSource(utils.DefaultTimeSource)

	// save away our server URL so we can rewrite our URLs
	serverURL = server.URL

	for _, tc := range flowTests {
		utils.SetUUIDGenerator(utils.NewSeededUUID4Generator(123456))
		utils.SetTimeSource(utils.NewSequentialTimeSource(time.Date(2018, 7, 6, 12, 30, 0, 123456789, time.UTC)))

		testJSON, err := ioutil.ReadFile(fmt.Sprintf("testdata/flows/%s", tc.output))
		require.NoError(t, err, "Error reading output file for flow '%s' and output '%s': %s", tc.assets, tc.output, err)

		flowTest := &FlowTest{}
		err = json.Unmarshal(json.RawMessage(testJSON), &flowTest)
		require.NoError(t, err, "Error unmarshalling output for flow '%s' and output '%s': %s", tc.assets, tc.output, err)

		// run our flow
		runResult, err := runFlow(tc.assets, flowTest.Trigger, flowTest.Resumes)
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
			flowTest := &FlowTest{Trigger: flowTest.Trigger, Resumes: flowTest.Resumes, Outputs: rawOutputs}
			testJSON, err := utils.JSONMarshalPretty(flowTest)
			require.NoError(t, err, "Error marshalling test definition: %s", err)

			testJSON, _ = normalizeJSON(testJSON)

			// write our output
			outputFilename := fmt.Sprintf("testdata/flows/%s", tc.output)
			err = ioutil.WriteFile(outputFilename, testJSON, 0644)
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
