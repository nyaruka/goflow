package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/nyaruka/goflow/flows/actions"
	"github.com/nyaruka/goflow/utils"
)

func handleActionDoc(output *bytes.Buffer, prefix string, typeName string, docString string) {
	lines := strings.Split(docString, "\n")
	name := ""

	docs := make([]string, 0, len(lines))
	example := make([]string, 0, len(lines))
	inExample := false
	for _, l := range lines {
		if strings.HasPrefix(l, prefix) {
			name = l[len(prefix)+1:]
		} else if strings.HasPrefix(l, "```") {
			inExample = !inExample
		} else if inExample {
			example = append(example, l[2:])
		} else {
			docs = append(docs, l)
		}
	}

	// try to parse our example
	exampleJSON := []byte(strings.Join(example, "\n"))
	typed := &utils.TypedEnvelope{}
	err := json.Unmarshal(exampleJSON, typed)
	action, err := actions.ActionFromEnvelope(typed)
	if err != nil {
		log.Fatalf("unable to parse example: %s\nHas err: %s", exampleJSON, err)
	}

	// validate it
	err = utils.Validate(action)
	if err != nil {
		log.Fatalf("unable to validate example: %s\nHad err: %s", exampleJSON, err)
	}

	// make sure types match
	if name != action.Type() {
		log.Fatalf("mismatched action types for example of %s", name)
	}

	typed, err = utils.EnvelopeFromTyped(action)
	if err != nil {
		log.Fatalf("unable to marshal example: %s\nHad err: %s", exampleJSON, err)
	}

	exampleJSON, err = json.MarshalIndent(typed, "", "  ")
	if err != nil {
		log.Fatalf("unable to marshal example: %s\nHad err: %s", exampleJSON, err)
	}

	// get the events created by this action
	events, err := eventsForAction(exampleJSON)
	if err != nil {
		//log.Fatalf("Error running action: %s\nHas err: %s", exampleJSON, err)
		events = json.RawMessage(fmt.Sprintf("error: %s", err.Error()))
	}

	if name != "" {
		if len(docs) > 0 && strings.HasPrefix(docs[0], typeName) {
			docs[0] = strings.Replace(docs[0], typeName, name, 1)
		}

		output.WriteString(fmt.Sprintf("## %s\n\n", name))
		output.WriteString(fmt.Sprintf("%s", strings.Join(docs, "\n")))
		if len(example) > 0 {
			output.WriteString(`<div class="input_action"><h3>Action</h3>`)
			output.WriteString("```json\n")
			output.WriteString(fmt.Sprintf("%s\n", exampleJSON))
			output.WriteString("```\n")
			output.WriteString(`</div>`)

			output.WriteString(`<div class="output_event"><h3>Event</h3>`)
			output.WriteString("```json\n")
			output.WriteString(fmt.Sprintf("%s\n", events))
			output.WriteString("```\n")
			output.WriteString(`</div>`)
		}
		output.WriteString(fmt.Sprintf("\n"))
	}
}
