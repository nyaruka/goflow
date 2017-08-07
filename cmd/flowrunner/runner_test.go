package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/engine"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

var flowTests = []struct {
	assets  string
	contact string
	output  string
}{
	{"two_questions.json", "", "two_questions_test.json"},
	{"subflow.json", "", "subflow_test.json"},
	{"brochure.json", "", "brochure_test.json"},
	{"all_actions.json", "", "all_actions_test.json"},
	{"default_result.json", "", "default_result_test.json"},
	{"empty.json", "", "empty_test.json"},
	{"node_loop.json", "", "node_loop_test.json"},
	{"subflow_loop.json", "", "subflow_loop_test.json"},
	{"date_parse.json", "", "date_parse_test.json"},
	{"webhook_persists.json", "", "webhook_persists_test.json"},
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

func runFlow(env utils.Environment, assetsFilename string, contactFilename string, flowUUID flows.FlowUUID, callerEvents [][]flows.Event) (flows.Session, []*Output, error) {
	assetsJSON, err := readFile("flows/", assetsFilename)
	if err != nil {
		return nil, nil, err
	}

	// rewrite the URL on any webhook actions
	assetsJSONStr := strings.Replace(string(assetsJSON), "http://localhost", serverURL, -1)

	assets, err := engine.ReadAssets(json.RawMessage(assetsJSONStr))
	if err != nil {
		return nil, nil, fmt.Errorf("Error unmarshalling assets '%s': %s", assetsFilename, err)
	}
	if err = assets.Validate(); err != nil {
		return nil, nil, fmt.Errorf("Error validiating assets '%s': %s", assetsFilename, err)
	}

	contactJSON, err := readFile("contacts/", contactFilename)
	if err != nil {
		return nil, nil, err
	}

	contact, err := flows.ReadContact(assets, json.RawMessage(contactJSON))
	if err != nil {
		return nil, nil, fmt.Errorf("Error unmarshalling contact '%s': %s", contactFilename, err)
	}

	// start our contact down this flow
	session := engine.NewSession(assets)
	session.SetEnvironment(env)
	session.SetContact(contact)

	err = session.StartFlow(flowUUID, callerEvents[0])
	if err != nil {
		return nil, nil, err
	}

	outputs := make([]*Output, 0)

	// for each of our remaining caller events
	resumeEvents := callerEvents[1:]
	for i := range resumeEvents {
		outJSON, err := json.MarshalIndent(session, "", "  ")
		if err != nil {
			return nil, nil, fmt.Errorf("Error marshalling output: %s", err)
		}
		outputs = append(outputs, &Output{outJSON, marshalEventLog(session.Log())})

		session, err = engine.ReadSession(assets, outJSON)
		if err != nil {
			return nil, nil, fmt.Errorf("Error marshalling output: %s", err)
		}

		// if we aren't at a wait, that's an error
		if session.Wait() == nil {
			return nil, nil, fmt.Errorf("Did not stop at expected wait, have unused resume events: %#v", resumeEvents[i:])
		}
		err = session.Resume(resumeEvents[i])
		if err != nil {
			return nil, nil, err
		}
	}

	outJSON, err := json.MarshalIndent(session, "", "  ")
	if err != nil {
		return nil, nil, fmt.Errorf("Error marshalling output: %s", err)
	}
	outputs = append(outputs, &Output{outJSON, marshalEventLog(session.Log())})

	return session, outputs, nil
}

// set up a mock server for webhook actions
func newTestHTTPServer() *httptest.Server {
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cmd := r.URL.Query().Get("cmd")
		defer r.Body.Close()
		w.Header().Set("Date", "")

		switch cmd {
		case "success":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{ "ok": "true" }`))
		case "unavailable":
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`{ "errors": ["service unavailable"] }`))
		default:
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{ "errors": ["bad_request"] }`))
		}
	}))
	// manually create a listener for our test server so that our output is predictable
	l, err := net.Listen("tcp", "127.0.0.1:49999")
	if err != nil {
		log.Fatal(err)
	}
	server.Listener = l
	return server
}

func TestFlows(t *testing.T) {
	env := utils.NewDefaultEnvironment()
	la, _ := time.LoadLocation("America/Los_Angeles")
	env.SetTimezone(la)

	server := newTestHTTPServer()
	server.Start()
	defer server.Close()

	// save away our server URL so we can rewrite our URLs
	serverURL = server.URL

	for _, test := range flowTests {
		testJSON, err := readFile("flows/", test.output)
		if err != nil {
			t.Errorf("Error reading output file for flow '%s' and output '%s': %s", test.assets, test.output, err)
			continue
		}

		flowTest := FlowTest{}
		err = json.Unmarshal(json.RawMessage(testJSON), &flowTest)
		if err != nil {
			t.Errorf("Error unmarshalling output for flow '%s' and output '%s': %s", test.assets, test.output, err)
			continue
		}

		// unmarshal our caller events
		callerEvents := make([][]flows.Event, len(flowTest.CallerEvents))

		for i := range flowTest.CallerEvents {
			callerEvents[i] = make([]flows.Event, len(flowTest.CallerEvents[i]))

			for e := range flowTest.CallerEvents[i] {
				event, err := events.EventFromEnvelope(flowTest.CallerEvents[i][e])
				if err != nil {
					t.Errorf("Error unmarshalling caller events for flow '%s' and output '%s': %s", test.assets, test.output, err)
					continue
				}
				event.SetFromCaller(true)
				callerEvents[i][e] = event
			}
		}

		// run our flow
		session, outputs, err := runFlow(env, test.assets, test.contact, flowTest.FlowUUID, callerEvents)
		if err != nil {
			t.Errorf("Error running flow for flow '%s' and output '%s': %s", test.assets, test.output, err)
			continue
		}

		if writeOutput {
			// we are writing new outputs, we write new files but don't test anything
			rawOutputs := make([]json.RawMessage, len(outputs))
			for i := range outputs {
				rawOutputs[i], err = json.Marshal(outputs[i])
				if err != nil {
					log.Fatal(err)
				}
			}
			flowTest := FlowTest{FlowUUID: flowTest.FlowUUID, CallerEvents: flowTest.CallerEvents, Outputs: rawOutputs}
			testJSON, err := json.MarshalIndent(flowTest, "", "  ")
			if err != nil {
				log.Fatal("Error marshalling test definition: ", err)
			}

			// write our output
			outputFilename := deriveFilename("flows/", test.output)
			err = ioutil.WriteFile(outputFilename, replaceFields(testJSON), 0644)
			if err != nil {
				log.Fatalf("Error writing test file to %s: %s\n", outputFilename, err)
			}
		} else {
			// unmarshal our expected outputs
			expectedOutputs := make([]*Output, len(flowTest.Outputs))
			for i := range expectedOutputs {
				output := &Output{}
				err := json.Unmarshal(flowTest.Outputs[i], output)
				if err != nil {
					log.Fatalf("Error unmarshalling output: %s", err)
				}
				expectedOutputs[i] = output
			}

			// read our output and test that we are the same
			if len(outputs) != len(expectedOutputs) {
				t.Errorf("Actual outputs:\n%s\n do not match expected:\n%s\n for flow '%s'", outputs, expectedOutputs, test.assets)
				continue
			}

			for i := range outputs {
				actualOutput := outputs[i]
				expectedOutput := expectedOutputs[i]

				actualSession, err := engine.ReadSession(session.Assets(), actualOutput.Session)
				if err != nil {
					t.Errorf("Error unmarshalling session running flow '%s': %s\n", test.assets, err)
				}
				expectedSession, err := engine.ReadSession(session.Assets(), expectedOutput.Session)
				if err != nil {
					t.Errorf("Error unmarshalling expected session running flow '%s': %s\n", test.assets, err)
				}

				// number of runs should be the same
				if len(actualSession.Runs()) != len(expectedSession.Runs()) {
					t.Errorf("Actual runs:\n%#v\n do not match expected:\n%#v\n for flow '%s'\n", actualSession.Runs(), expectedSession.Runs(), test.assets)
				}

				// runs should have same status and flows
				for i := range actualSession.Runs() {
					run := actualSession.Runs()[i]
					expected := expectedSession.Runs()[i]

					if run.Flow() != expected.Flow() {
						t.Errorf("Actual run flow: %s does not match expected: %s for flow '%s'", run.Flow().UUID(), expected.Flow().UUID(), test.assets)
					}

					if run.Status() != expected.Status() {
						t.Errorf("Actual run status: %s does not match expected: %s for flow '%s'", run.Status(), expected.Status(), test.assets)
					}
				}

				if len(actualOutput.Log) != len(expectedOutput.Log) {
					t.Errorf("Actual events:\n%#v\n do not match expected:\n%#v\n for flow '%s'\n", actualOutput.Log, expectedOutput.Log, test.assets)
				}

				for j := range actualOutput.Log {
					event := actualOutput.Log[j]
					expected := expectedOutput.Log[j]

					// write our events as json
					eventJSON, err := rawMessageAsJSON(event)
					if err != nil {
						t.Errorf("Error marshalling event for flow '%s' and output '%s': %s\n", test.assets, test.output, err)
					}
					expectedJSON, err := rawMessageAsJSON(expected)
					if err != nil {
						t.Errorf("Error marshalling expected event for flow '%s' and output '%s': %s\n", test.assets, test.output, err)
					}

					if eventJSON != expectedJSON {
						t.Errorf("Got event:\n'%s'\n\nwhen expecting:\n'%s'\n\n for flow '%s' and output '%s\n", eventJSON, expectedJSON, test.assets, test.output)
						break
					}
				}
			}
		}
	}
}
