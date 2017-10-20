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
	uuid "github.com/satori/go.uuid"
)

type Output struct {
	Session json.RawMessage   `json:"session"`
	Log     []json.RawMessage `json:"log"`
}

type FlowTest struct {
	Trigger      *utils.TypedEnvelope     `json:"trigger"`
	CallerEvents [][]*utils.TypedEnvelope `json:"caller_events"`
	Outputs      []json.RawMessage        `json:"outputs"`
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
		"arrived_on":      "2000-01-01T00:00:00.000000000-00:00",
		"left_on":         "2000-01-01T00:00:00.000000000-00:00",
		"exited_on":       "2000-01-01T00:00:00.000000000-00:00",
		"created_on":      "2000-01-01T00:00:00.000000000-00:00",
		"modified_on":     "2000-01-01T00:00:00.000000000-00:00",
		"expires_on":      "2000-01-01T00:00:00.000000000-00:00",
		"timeout_on":      "2000-01-01T00:00:00.000000000-00:00",
		"event.uuid":      "",
		"path.uuid":       "",
		"run.uuid":        "4213ac47-93fd-48c4-af12-7da8218ef09d",
		"runs.uuid":       "",
		"step_uuid":       "",
		"parent_uuid":     "",
		"parent_run_uuid": "",
		"child_uuid":      "",
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
	assetCache := engine.NewAssetCache(100, 5)
	if err := assetCache.Include(json.RawMessage(assetsJSON)); err != nil {
		log.Fatal("Error reading assets: ", err)
	}

	// create our environment
	env := utils.NewDefaultEnvironment()
	la, _ := time.LoadLocation("America/Los_Angeles")
	env.SetTimezone(la)

	session := engine.NewSession(assetCache, engine.NewTestAssetServer())

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

	trigger := triggers.NewManualTrigger(flow, time.Now())
	session.SetEnvironment(env)
	session.SetContact(contact)

	// and start our flow
	err = session.Start(trigger, nil)
	if err != nil {
		log.Fatal("Error starting flow: ", err)
	}

	scanner := bufio.NewScanner(os.Stdin)

	callerEvents := make([][]flows.Event, 0)
	callerEvents = append(callerEvents, []flows.Event{})

	outputs := make([]*Output, 0)

	channelUUID := flows.ChannelUUID("57f1078f-88aa-46f4-a59a-948a5739c03d")

	for session.Wait() != nil {
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
		channelRef := flows.NewChannelReference(channelUUID, "Test Channel")
		event := events.NewMsgReceivedEvent(flows.InputUUID(uuid.NewV4().String()), channelRef, contact.Reference(), contact.URNs()[0], scanner.Text(), []flows.Attachment{})
		event.SetFromCaller(true)
		callerEvents = append(callerEvents, []flows.Event{event})

		// rebuild our session
		session, err = engine.ReadSession(assetCache, engine.NewTestAssetServer(), outJSON)
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
	outJSON, err := json.MarshalIndent(session, "", "  ")
	if err != nil {
		log.Fatal("Error marshalling output: ", err)
	}
	fmt.Printf("%s\n", outJSON)
	outputs = append(outputs, &Output{outJSON, marshalEventLog(session.Log())})

	// write out our test file
	if *writePtr {
		// name of the test file is the same as our assets file, just with _test.json intead of .json
		testFilename := strings.Replace(assetsFilename, ".json", "_test.json", 1)

		callerEventEnvelopes := make([][]*utils.TypedEnvelope, len(callerEvents))
		for i := range callerEvents {
			callerEventEnvelopes[i] = envelopesForEvents(callerEvents[i])
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
