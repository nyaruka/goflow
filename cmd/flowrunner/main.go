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

	"github.com/nyaruka/goflow/assets"
	_ "github.com/nyaruka/goflow/extensions/transferto"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/engine"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/flows/triggers"
	"github.com/nyaruka/goflow/utils"
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

func normalizeJSON(data json.RawMessage) ([]byte, error) {
	var asMap map[string]interface{}
	if err := json.Unmarshal(data, &asMap); err != nil {
		return nil, err
	}

	return utils.JSONMarshalPretty(asMap)
}

func marshalEventLog(eventLog []flows.Event) []json.RawMessage {
	envelopes, err := events.EventsToEnvelopes(eventLog)
	marshaled := make([]json.RawMessage, len(envelopes))

	for i := range envelopes {
		marshaled[i], err = utils.JSONMarshal(envelopes[i])
		if err != nil {
			log.Fatalf("error creating marshaling envelope %s: %s", envelopes[i], err)
		}
	}
	return marshaled
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

	httpClient := utils.NewHTTPClient("goflow-flowrunner")

	assetsFilename := flag.Args()[0]
	startFlowUUID := flows.FlowUUID(flag.Args()[1])

	fmt.Printf("Parsing: %s\n", assetsFilename)
	assetsJSON, err := ioutil.ReadFile(assetsFilename)
	if err != nil {
		log.Fatal("Error reading assets file: ", err)
	}
	assetCache := assets.NewAssetCache(100, 5)
	if err := assetCache.Include(json.RawMessage(assetsJSON)); err != nil {
		log.Fatal("Error reading assets: ", err)
	}

	// create our environment
	la, _ := time.LoadLocation("America/Los_Angeles")
	env := utils.NewEnvironment(utils.DateFormatYearMonthDay, utils.TimeFormatHourMinute, la, utils.LanguageList{}, utils.RedactionPolicyNone)

	assets := engine.NewSessionAssets(engine.NewMockAssetServer(assetCache))
	session := engine.NewSession(assets, engine.NewDefaultConfig(), httpClient)

	contactJSON, err := ioutil.ReadFile(*contactFile)
	if err != nil {
		log.Fatal("error reading contact file: ", err)
	}
	contact, err := flows.ReadContact(session.Assets(), json.RawMessage(contactJSON))
	if err != nil {
		log.Fatal("error unmarshalling contact: ", err)
	}
	flow, err := session.Assets().GetFlow(startFlowUUID)
	if err != nil {
		log.Fatal("error accessing flow: ", err)
	}

	trigger := triggers.NewManualTrigger(env, contact, flow, nil, time.Now())

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
		outJSON, err := utils.JSONMarshalPretty(session)
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
		assets := engine.NewSessionAssets(engine.NewMockAssetServer(assetCache))
		session, err = engine.ReadSession(assets, engine.NewDefaultConfig(), httpClient, outJSON)
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
	outJSON, err := utils.JSONMarshalPretty(session)
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
			rawOutputs[i], err = utils.JSONMarshal(outputs[i])
			if err != nil {
				log.Fatal(err)
			}
		}

		triggerEnvelope, err := utils.EnvelopeFromTyped(trigger)
		if err != nil {
			log.Fatal(err)
		}

		flowTest := FlowTest{Trigger: triggerEnvelope, CallerEvents: callerEventEnvelopes, Outputs: rawOutputs}
		testJSON, err := utils.JSONMarshal(flowTest)
		if err != nil {
			log.Fatal("Error marshalling test definition: ", err)
		}

		testJSON, _ = normalizeJSON(testJSON)

		// write our output
		err = ioutil.WriteFile(testFilename, testJSON, 0644)
		if err != nil {
			log.Fatalf("Error writing test file to %s: %s\n", testFilename, err)
		}
	}
}
