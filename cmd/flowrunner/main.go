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
	"github.com/nyaruka/goflow/flows/engine"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/flows/triggers"
	"github.com/nyaruka/goflow/utils"
)

const (
	outputIndent string = "    "
)

type Output struct {
	Session json.RawMessage   `json:"session"`
	Events  []json.RawMessage `json:"events"`
}

type FlowTest struct {
	Trigger      *utils.TypedEnvelope     `json:"trigger"`
	CallerEvents [][]*utils.TypedEnvelope `json:"caller_events"`
	Outputs      []json.RawMessage        `json:"outputs"`
}

func marshalEventLog(eventLog []flows.Event) []json.RawMessage {
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
	envJSON, err := json.MarshalIndent(msg, "", outputIndent)
	if err != nil {
		return "", err
	}

	return string(clearTimestamps(envJSON)), nil
}

func eventAsJSON(event flows.Event) (string, error) {
	env, err := utils.EnvelopeFromTyped(event)
	if err != nil {
		return "", err
	}

	envJSON, err := json.MarshalIndent(env, "", outputIndent)
	if err != nil {
		return "", err
	}

	return string(clearTimestamps(envJSON)), nil
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

func clearTimestamps(input []byte) []byte {
	placeholder := "2000-01-01T00:00:00.000000000-00:00"

	replacements := map[string]interface{}{
		"arrived_on":  placeholder,
		"left_on":     placeholder,
		"exited_on":   placeholder,
		"created_on":  placeholder,
		"modified_on": placeholder,
		"expires_on":  placeholder,
		"timeout_on":  placeholder,
	}

	// unmarshal to arbitrary json
	inputJSON := make(map[string]interface{})
	err := json.Unmarshal(input, &inputJSON)
	if err != nil {
		log.Fatalf("Error unmarshalling: %s", err)
	}

	replaceMapFields(replacements, "", inputJSON)

	// return our marshalled result
	outputJSON, err := json.MarshalIndent(inputJSON, "", outputIndent)
	if err != nil {
		log.Fatalf("Error marshalling: %s", err)
	}
	return outputJSON
}

func main() {
	testdata := filepath.Join(os.Getenv("GOPATH"), "src/github.com/nyaruka/goflow/cmd/flowrunner/testdata")

	writePtr := flag.Bool("write", false, "Whether to write a _test.json file for this flow")
	contactFile := flag.String("contact", filepath.Join(testdata, "contacts/default.json"), "The location of the JSON file defining the contact to use, defaulting to test/contacts/default.json")

	flag.Parse()

	if len(flag.Args()) != 2 {
		fmt.Printf("\nUsage: runner [-write] <assets.json> flow_uuid\n\n")
		os.Exit(1)
	}

	assetsFilename := flag.Args()[0]
	startFlowUUID := flows.FlowUUID(flag.Args()[1])

	fmt.Printf("Parsing: %s\n", assetsFilename)
	assetsJSON, err := ioutil.ReadFile(assetsFilename)
	if err != nil {
		log.Fatal("Error reading assets file: ", err)
	}
	assetCache := engine.NewAssetCache(100, 5, "testing/1.0")
	if err := assetCache.Include(json.RawMessage(assetsJSON)); err != nil {
		log.Fatal("Error reading assets: ", err)
	}

	// create our environment
	env := utils.NewDefaultEnvironment()
	la, _ := time.LoadLocation("America/Los_Angeles")
	env.SetTimezone(la)

	session := engine.NewSession(assetCache, engine.NewMockAssetServer())

	contactJSON, err := ioutil.ReadFile(*contactFile)
	if err != nil {
		log.Fatal("error reading contact file: ", err)
	}
	contact, err := flows.ReadContact(session, json.RawMessage(contactJSON))
	if err != nil {
		log.Fatal("error unmarshalling contact: ", err)
	}
	flow, err := session.Assets().GetFlow(startFlowUUID)
	if err != nil {
		log.Fatal("error accessing flow: ", err)
	}

	trigger := triggers.NewManualTrigger(env, contact, flow, utils.EmptyJSONFragment, time.Now())

	// and start our flow
	err = session.Start(trigger, nil)
	if err != nil {
		log.Fatal("Error starting flow: ", err)
	}

	scanner := bufio.NewScanner(os.Stdin)

	callerEvents := make([][]flows.Event, 0)
	callerEvents = append(callerEvents, []flows.Event{})

	outputs := make([]*Output, 0)

	for session.Wait() != nil {
		outJSON, err := json.MarshalIndent(session, "", outputIndent)
		if err != nil {
			log.Fatal("Error marshalling output: ", err)
		}
		fmt.Printf("%s\n", outJSON)
		outputs = append(outputs, &Output{outJSON, marshalEventLog(session.Events())})

		// print any msg_created events
		for _, event := range session.Events() {
			if event.Type() == events.TypeMsgCreated {
				fmt.Printf(">>> %s\n", event.(*events.MsgCreatedEvent).Msg.Text())
			}
		}

		// ask for input
		fmt.Printf("<<< ")
		scanner.Scan()

		// create our event to resume with
		msg := flows.NewMsgIn(flows.MsgUUID(utils.NewUUID()), contact.URNs()[0].URN, nil, scanner.Text(), []flows.Attachment{})
		event := events.NewMsgReceivedEvent(msg)
		event.SetFromCaller(true)
		callerEvents = append(callerEvents, []flows.Event{event})

		// rebuild our session
		session, err = engine.ReadSession(assetCache, engine.NewMockAssetServer(), outJSON)
		if err != nil {
			log.Fatalf("Error unmarshalling output: %s", err)
		}

		err = session.Resume([]flows.Event{event})
		if err != nil {
			log.Print("Error resuming flow: ", err)
			break
		}
	}

	// print out our context
	outJSON, err := json.MarshalIndent(session, "", outputIndent)
	if err != nil {
		log.Fatal("Error marshalling output: ", err)
	}
	fmt.Printf("%s\n", outJSON)
	outputs = append(outputs, &Output{outJSON, marshalEventLog(session.Events())})

	// write out our test file
	if *writePtr {
		// name of the test file is the same as our assets file, just with _test.json instead of .json
		testFilename := strings.Replace(assetsFilename, ".json", "_test.json", 1)

		callerEventEnvelopes := make([][]*utils.TypedEnvelope, len(callerEvents))
		for i := range callerEvents {
			callerEventEnvelopes[i], _ = events.EventsToEnvelopes(callerEvents[i])
		}

		rawOutputs := make([]json.RawMessage, len(outputs))
		for i := range outputs {
			rawOutputs[i], err = json.Marshal(outputs[i])
			if err != nil {
				log.Fatal(err)
			}
		}

		triggerEnvelope, err := utils.EnvelopeFromTyped(trigger)
		if err != nil {
			log.Fatal(err)
		}

		flowTest := FlowTest{Trigger: triggerEnvelope, CallerEvents: callerEventEnvelopes, Outputs: rawOutputs}
		testJSON, err := json.MarshalIndent(flowTest, "", outputIndent)
		if err != nil {
			log.Fatal("Error marshalling test definition: ", err)
		}

		// write our output
		err = ioutil.WriteFile(testFilename, clearTimestamps(testJSON), 0644)
		if err != nil {
			log.Fatalf("Error writing test file to %s: %s\n", testFilename, err)
		}
	}
}
