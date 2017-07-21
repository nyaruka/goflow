package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/definition"
	"github.com/nyaruka/goflow/flows/engine"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/flows/runs"
	"github.com/nyaruka/goflow/utils"
)

type Output struct {
	Session json.RawMessage   `json:"session"`
	Log     []json.RawMessage `json:"log"`
}

type FlowTest struct {
	Extra        json.RawMessage        `json:"extra,omitempty"`
	CallerEvents []*utils.TypedEnvelope `json:"caller_events"`
	Outputs      []json.RawMessage      `json:"outputs"`
}

func envelopesForEvents(events []flows.Event) []*utils.TypedEnvelope {
	envelopes := make([]*utils.TypedEnvelope, len(events))
	for i := range events {
		envelope, err := utils.EnvelopeFromTyped(events[i])
		if err != nil {
			log.Fatalf("Error creating envelope for %s: %s", events[i], err)
		}

		envelopes[i] = envelope
	}
	return envelopes
}

func marshalEventLog(eventLog []flows.LogEntry) []json.RawMessage {
	envelopes := make([]json.RawMessage, len(eventLog))
	for i := range eventLog {
		envelope, err := json.Marshal(eventLog[i])
		if err != nil {
			log.Fatalf("Error creating envelope for %s: %s", eventLog[i], err)
		}

		envelopes[i] = envelope
	}
	return envelopes
}

func rawMessageAsJSON(msg json.RawMessage) (string, error) {
	envJSON, err := json.MarshalIndent(msg, "", "  ")
	if err != nil {
		return "", err
	}

	return string(replaceFields(envJSON)), nil
}

func eventAsJSON(event flows.Event) (string, error) {
	env, err := utils.EnvelopeFromTyped(event)
	if err != nil {
		return "", err
	}

	envJSON, err := json.MarshalIndent(env, "", "  ")
	if err != nil {
		return "", err
	}

	return string(replaceFields(envJSON)), nil
}

func replaceArrayFields(replacements map[string]interface{}, parent string, arrFields []interface{}) {
	for _, e := range arrFields {
		switch child := e.(type) {
		case map[string]interface{}:
			replaceMapFields(replacements, parent, child)
		case []interface{}:
			replaceArrayFields(replacements, parent, child)
		}
	}
}

func replaceMapFields(replacements map[string]interface{}, parent string, mapFields map[string]interface{}) {
	for k, v := range mapFields {
		replacement, found := replacements[k]
		if found {
			mapFields[k] = replacement
			continue
		}

		if parent != "" {
			parentKey := parent + "." + k
			replacement, found = replacements[parentKey]
			if found {
				mapFields[k] = replacement
				continue
			}
		}

		switch child := v.(type) {
		case map[string]interface{}:
			replaceMapFields(replacements, k, child)
		case []interface{}:
			replaceArrayFields(replacements, k, child)
		}
	}
}

func replaceFields(input []byte) []byte {
	replacements := map[string]interface{}{
		"arrived_on":  "2000-01-01T00:00:00.000000000-00:00",
		"left_on":     "2000-01-01T00:00:00.000000000-00:00",
		"exited_on":   "2000-01-01T00:00:00.000000000-00:00",
		"created_on":  "2000-01-01T00:00:00.000000000-00:00",
		"modified_on": "2000-01-01T00:00:00.000000000-00:00",
		"expires_on":  "2000-01-01T00:00:00.000000000-00:00",
		"timesout_on": "2000-01-01T00:00:00.000000000-00:00",
		"event.uuid":  "",
		"path.uuid":   "",
		"runs.uuid":   "",
		"step_uuid":   "",
		"parent_uuid": "",
		"child_uuid":  "",
	}

	// unmarshal to arbitrary json
	inputJSON := make(map[string]interface{})
	err := json.Unmarshal(input, &inputJSON)
	if err != nil {
		log.Fatalf("Error unmarshalling: %s", err)
	}

	replaceMapFields(replacements, "", inputJSON)

	// return our marshalled result
	outputJSON, err := json.MarshalIndent(inputJSON, "", "  ")
	if err != nil {
		log.Fatalf("Error marshalling: %s", err)
	}
	return outputJSON
}

func main() {
	testdata := filepath.Join(os.Getenv("GOPATH"), "src/github.com/nyaruka/goflow/cmd/flowrunner/testdata")

	writePtr := flag.Bool("write", false, "Whether to write a _test.json file for this flow")
	contactFile := flag.String("contact", filepath.Join(testdata, "contacts/default.json"), "The location of the JSON file defining the contact to use, defaulting to test/contacts/default.json")
	channelFile := flag.String("channel", filepath.Join(testdata, "channels/default.json"), "The location of the JSON file defining the channel to use, defaulting to test/channels/default.json")

	flag.Parse()

	if len(flag.Args()) != 1 {
		fmt.Printf("\nUsage: runner [-write] <flow.json>\n\n")
		os.Exit(1)
	}

	flowFilename := flag.Args()[0]

	fmt.Printf("Parsing: %s\n", flowFilename)
	flowJSON, err := ioutil.ReadFile(flowFilename)
	if err != nil {
		log.Fatal("Error reading flow file: ", err)
	}
	runnerFlows, err := definition.ReadFlows(json.RawMessage(flowJSON))
	if err != nil {
		log.Fatal("Error reading flows: ", err)
	}

	contactJSON, err := ioutil.ReadFile(*contactFile)
	if err != nil {
		log.Fatal("Error reading contact file: ", err)
	}
	contact, err := flows.ReadContact(json.RawMessage(contactJSON))
	if err != nil {
		log.Fatal("Error unmarshalling contact: ", err)
	}

	channelJSON, err := ioutil.ReadFile(*channelFile)
	if err != nil {
		log.Fatal("Error reading channel file: ", err)
	}
	channel, err := flows.ReadChannel(json.RawMessage(channelJSON))
	if err != nil {
		log.Fatal("Error unmarshalling channel: ", err)
	}

	// create our flow environment
	baseEnv := utils.NewDefaultEnvironment()
	la, _ := time.LoadLocation("America/Los_Angeles")
	baseEnv.SetTimezone(la)
	env := engine.NewSessionEnvironment(baseEnv, runnerFlows, []flows.Channel{channel}, []*flows.Contact{contact})

	// and start our flow
	session, err := engine.StartFlow(env, runnerFlows[0], contact, nil, nil, nil, nil)
	if err != nil {
		log.Fatal("Error starting flow: ", err)
	}

	scanner := bufio.NewScanner(os.Stdin)
	inputs := make([]flows.Event, 0)
	outputs := make([]*Output, 0)

	channelUUID := flows.ChannelUUID("57f1078f-88aa-46f4-a59a-948a5739c03d")

	run := session.ActiveRun()
	for run != nil && run.Wait() != nil {
		outJSON, err := json.MarshalIndent(session, "", "  ")
		if err != nil {
			log.Fatal("Error marshalling output: ", err)
		}
		fmt.Printf("%s\n", outJSON)
		outputs = append(outputs, &Output{outJSON, marshalEventLog(session.Log())})

		// print any send_msg events
		for _, e := range session.Log() {
			if e.Event().Type() == events.TypeSendMsg {
				fmt.Printf(">>> %s\n", e.Event().(*events.SendMsgEvent).Text)
			}
		}

		// ask for input
		fmt.Printf("<<< ")
		scanner.Scan()

		// create our event to resume with
		event := events.NewMsgReceivedEvent(channelUUID, contact.UUID(), contact.URNs()[0], scanner.Text())
		event.SetFromCaller(true)
		inputs = append(inputs, event)

		// rebuild our session
		session, err = runs.ReadSession(session.Environment(), outJSON)
		if err != nil {
			log.Fatalf("Error unmarshalling output: %s", err)
		}
		baseEnv := utils.NewDefaultEnvironment()
		la, _ := time.LoadLocation("America/Los_Angeles")
		baseEnv.SetTimezone(la)

		run = session.ActiveRun()

		session, err = engine.ResumeFlow(env, run, []flows.Event{event})
		if err != nil {
			log.Print("Error resuming flow: ", err)
			break
		}

		run = session.ActiveRun()
	}

	// print out our context
	outJSON, err := json.MarshalIndent(session, "", "  ")
	if err != nil {
		log.Fatal("Error marshalling output: ", err)
	}
	fmt.Printf("%s\n", outJSON)
	outputs = append(outputs, &Output{outJSON, marshalEventLog(session.Log())})

	// write out our test file
	if *writePtr {
		// name of the test file is the same as our flow file, just with _test.json intead of .json
		testFilename := strings.Replace(flowFilename, ".json", "_test.json", 1)

		envelopes := envelopesForEvents(inputs)
		rawOutputs := make([]json.RawMessage, len(outputs))
		for i := range outputs {
			rawOutputs[i], err = json.Marshal(outputs[i])
			if err != nil {
				log.Fatal(err)
			}
		}

		flowTest := FlowTest{nil, envelopes, rawOutputs}
		testJSON, err := json.MarshalIndent(flowTest, "", "  ")
		if err != nil {
			log.Fatal("Error marshalling test definition: ", err)
		}

		// write our output
		err = ioutil.WriteFile(testFilename, replaceFields(testJSON), 0644)
		if err != nil {
			log.Fatalf("Error writing test file to %s: %s\n", testFilename, err)
		}
	}
}
