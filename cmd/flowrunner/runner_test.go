package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"testing"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/engine"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/flows/flow"
	"github.com/nyaruka/goflow/utils"
)

var flowTests = []struct {
	flow    string
	contact string
	channel string
	output  string
}{
	{"two_questions.json", "", "", "two_questions_test.json"},
	{"subflow.json", "", "", "subflow_test.json"},
}

var writeOutput bool

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

func runFlow(env utils.Environment, flowFilename string, contactFilename string, channelFilename string, resumeEvents []flows.Event) (flows.RunOutput, error) {
	flowJSON, err := readFile("flows/", flowFilename)
	if err != nil {
		return nil, err
	}
	runnerFlows, err := flow.ReadFlows(json.RawMessage(flowJSON))
	if err != nil {
		return nil, fmt.Errorf("Error unmarshalling flows '%s': %s", flowFilename, err)
	}

	contactJSON, err := readFile("contacts/", contactFilename)
	if err != nil {
		return nil, err
	}
	contact, err := flow.ReadContact(json.RawMessage(contactJSON))
	if err != nil {
		return nil, fmt.Errorf("Error unmarshalling contact '%s': %s", contactFilename, err)
	}

	channelJSON, err := readFile("channels/", channelFilename)
	if err != nil {
		return nil, err
	}
	_, err = flow.ReadChannel(json.RawMessage(channelJSON))
	if err != nil {
		return nil, fmt.Errorf("Error unmarshalling channel '%s': %s", channelFilename, err)
	}

	// start our contact down this flow
	flowEnv := engine.NewFlowEnvironment(env, runnerFlows)
	output, err := engine.StartFlow(flowEnv, runnerFlows[0], contact, nil)
	if err != nil {
		return nil, err
	}

	// for each of our resume events
	for i := range resumeEvents {
		activeRun := output.ActiveRun()

		// if we aren't at a wait, that's an error
		if activeRun == nil {
			return nil, fmt.Errorf("Did not stop at expected wait, have unused resume events: %#v", resumeEvents[i:])
		}

		// resume the flow
		output, err = engine.ResumeFlow(flowEnv, activeRun, resumeEvents[i])
		if err != nil {
			return nil, err
		}
	}

	return output, nil
}

func TestFlows(t *testing.T) {
	env := utils.NewDefaultEnvironment()

	for _, test := range flowTests {
		testJSON, err := readFile("flows/", test.output)
		if err != nil {
			t.Errorf("Error reading output file for flow '%s' and output '%s': %s", test.flow, test.output, err)
			continue
		}

		flowTest := FlowTest{}
		err = json.Unmarshal(json.RawMessage(testJSON), &flowTest)
		if err != nil {
			t.Errorf("Error unmarshalling output for flow '%s' and output '%s': %s", test.flow, test.output, err)
			continue
		}

		// unmarshal our resume events
		resumeEvents := make([]flows.Event, len(flowTest.ResumeEvents))
		for i := range flowTest.ResumeEvents {
			resumeEvents[i], err = events.EventFromEnvelope(flowTest.ResumeEvents[i])
			if err != nil {
				t.Errorf("Error unmarshalling resume events for flow '%s' and output '%s': %s", test.flow, test.output, err)
				continue
			}
		}

		// run our flow
		output, err := runFlow(env, test.flow, test.contact, test.channel, resumeEvents)
		if err != nil {
			t.Errorf("Error running flow for flow '%s' and output '%s': %s", test.flow, test.output, err)
			continue
		}

		if writeOutput {
			// we are writing new outputs, we write new files but don't test anything
			envelopes, err := envelopesForEvents(resumeEvents)
			if err != nil {
				log.Fatal("Error marshalling inputs: ", err)
			}

			outJSON, err := json.MarshalIndent(output, "", "  ")
			if err != nil {
				log.Fatal("Error marshalling output: ", err)
			}

			flowTest := FlowTest{envelopes, outJSON}
			testJSON, err := json.MarshalIndent(flowTest, "", "  ")
			if err != nil {
				log.Fatal("Error marshalling test definition: ", err)
			}

			// write our output
			outputFilename := deriveFilename("test/flows/", test.output)
			err = ioutil.WriteFile(outputFilename, replaceFields(testJSON), 0644)
			if err != nil {
				log.Fatalf("Error writing test file to %s: %s\n", outputFilename, err)
			}
		} else {
			// read our output and test that we are the same
			expectedOutput, err := flow.ReadRunOutput(flowTest.Output)
			if err != nil {
				t.Errorf("Error unmarshalling resume events for flow '%s' and output '%s': %s", test.flow, test.output, err)
				continue
			}

			// check the expected and actual events
			if len(output.Events()) != len(expectedOutput.Events()) {
				t.Errorf("Actual events: '%#v' do not match expected '%#v' for flow '%s' and output '%s'", output.Events(), expectedOutput.Events(), test.flow, test.output)
				continue
			}

			for i := range output.Events() {
				event := output.Events()[i]
				expected := expectedOutput.Events()[i]

				if event.Type() != expected.Type() {
					t.Errorf("Got event '%#v' when expecting: '%#v' for flow '%s' and output '%s", event, expected, test.flow, test.output)
					break
				}

				// write our events as json
				eventJSON, err := eventAsJSON(event)
				if err != nil {
					t.Errorf("Error marshalling event for flow '%s' and output '%s': %s", test.flow, test.output, err)
				}
				expectedJSON, err := eventAsJSON(expected)
				if err != nil {
					t.Errorf("Error marshalling expected event for flow '%s' and output '%s': %s", test.flow, test.output, err)
				}

				if eventJSON != expectedJSON {
					t.Errorf("Got event:\n'%s'\n\nwhen expecting:\n'%s'\n\n for flow '%s' and output '%s", eventJSON, expectedJSON, test.flow, test.output)
					break
				}
			}
		}
	}
}
