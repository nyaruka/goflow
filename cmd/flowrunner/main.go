package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/engine"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/flows/flow"
	"github.com/nyaruka/goflow/utils"
)

type FlowTest struct {
	ResumeEvents []*utils.TypedEnvelope `json:"resume_events"`
	Outputs      []json.RawMessage      `json:"outputs"`
}

func envelopesForEvents(events []flows.Event) ([]*utils.TypedEnvelope, error) {
	envelopes := make([]*utils.TypedEnvelope, len(events))
	for i := range events {
		envelope, err := utils.EnvelopeFromTyped(events[i])
		if err != nil {
			return nil, err
		}

		envelopes[i] = envelope
	}
	return envelopes, nil
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

func replaceFields(input []byte) []byte {
	fields := map[string]string{
		"arrived_on":  "2000-01-01T00:00:00.000000000-00:00",
		"left_on":     "2000-01-01T00:00:00.000000000-00:00",
		"exited_on":   "2000-01-01T00:00:00.000000000-00:00",
		"created_on":  "2000-01-01T00:00:00.000000000-00:00",
		"modified_on": "2000-01-01T00:00:00.000000000-00:00",
		"expires_on":  "2000-01-01T00:00:00.000000000-00:00",
		"timesout_on": "2000-01-01T00:00:00.000000000-00:00",
		"uuid":        "",
		"step":        "",
	}

	output := bytes.Buffer{}
	for i := 0; i < len(input); i++ {
		b := input[i]
		output.WriteByte(b)

		// if this is a quote, figure out if we were part of one of our fields
		if b == '"' {
			replaceField := ""
			outputLen := output.Len()
			for f := range fields {
				field := fmt.Sprintf("\"%s\"", f)
				if outputLen < len(field) {
					continue
				}
				lastPiece := output.String()[outputLen-len(field):]
				if lastPiece == field {
					replaceField = f
					break
				}
			}

			// we are skipping this field, read until we see a newline
			if replaceField != "" {
				i++
				var addRune = ' '
				for ; i < len(input) && addRune == ' '; i++ {
					if input[i] == ',' || input[i] == '\n' {
						addRune = rune(input[i])
					}
				}
				i--
				// write our empty value
				output.WriteString(fmt.Sprintf(": \"%s\"%c", fields[replaceField], addRune))
			}
		}
	}
	return output.Bytes()
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
	runnerFlows, err := flow.ReadFlows(json.RawMessage(flowJSON))
	if err != nil {
		log.Fatal("Error reading flows: ", err)
	}

	contactJSON, err := ioutil.ReadFile(*contactFile)
	if err != nil {
		log.Fatal("Error reading contact file: ", err)
	}
	contact, err := flow.ReadContact(json.RawMessage(contactJSON))
	if err != nil {
		log.Fatal("Error unmarshalling contact: ", err)
	}

	channelJSON, err := ioutil.ReadFile(*channelFile)
	if err != nil {
		log.Fatal("Error reading channel file: ", err)
	}
	_, err = flow.ReadChannel(json.RawMessage(channelJSON))
	if err != nil {
		log.Fatal("Error unmarshalling channel: ", err)
	}

	// create our flow environment
	env := engine.NewFlowEnvironment(utils.NewDefaultEnvironment(), runnerFlows, []flows.FlowRun{}, []flows.Contact{contact})

	// and start our flow
	output, err := engine.StartFlow(env, runnerFlows[0], contact, nil)
	if err != nil {
		log.Fatal("Error starting flow: ", err)
	}

	scanner := bufio.NewScanner(os.Stdin)
	inputs := make([]flows.Event, 0)
	outputs := make([]json.RawMessage, 0)

	run := output.ActiveRun()
	for run != nil && run.Wait() != nil {
		outJSON, err := json.MarshalIndent(output, "", "  ")
		if err != nil {
			log.Fatal("Error marshalling output: ", err)
		}
		fmt.Printf("%s\n", outJSON)
		outputs = append(outputs, outJSON)

		// print any events
		for _, e := range output.Events() {
			if e.Type() == events.MSG_OUT {
				fmt.Printf(">>> %s\n", e.(*events.MsgOutEvent).Text)
			}
		}

		// ask for input
		fmt.Printf("<<< ")
		scanner.Scan()

		// create our event to resume with
		event := events.NewIncomingMsgEvent("", contact.UUID(), scanner.Text())
		inputs = append(inputs, event)

		// rebuild our output
		output, err = flow.ReadRunOutput(outJSON)
		if err != nil {
			log.Fatalf("Error unmarshalling output: %s", err)
		}
		env = engine.NewFlowEnvironment(utils.NewDefaultEnvironment(), runnerFlows, output.Runs(), []flows.Contact{contact})

		for _, run := range output.Runs() {
			err = run.Hydrate(env)
			if err != nil {
				log.Fatalf("Error hydrating run: %s", err)
			}
		}
		run = output.ActiveRun()

		output, err = engine.ResumeFlow(env, run, event)
		if err != nil {
			log.Print("Error resuming flow: ", err)
			break
		}

		run = output.ActiveRun()
	}

	// print out our context
	outJSON, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		log.Fatal("Error marshalling output: ", err)
	}
	fmt.Printf("%s\n", outJSON)
	outputs = append(outputs, outJSON)

	// write out our test file
	if *writePtr {
		// name of the test file is the same as our flow file, just with _test.json intead of .json
		testFilename := strings.Replace(flowFilename, ".json", "_test.json", 1)

		envelopes, err := envelopesForEvents(inputs)
		if err != nil {
			log.Fatal("Error marshalling inputs: ", err)
		}

		flowTest := FlowTest{envelopes, outputs}
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
