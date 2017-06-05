package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/actions"
	"github.com/nyaruka/goflow/flows/definition"
	"github.com/nyaruka/goflow/flows/engine"
	"github.com/nyaruka/goflow/utils"
)

var flowDef = `
{
	"name": "ActionFlow",
	"uuid": "50c3706e-fedb-42c0-8eab-dda3335714b7",
	"nodes": [{
		"uuid": "72a1f5df-49f9-45df-94c9-d86f7ea064e5",
		"actions": [%s]
	}]
}
`

var subflowDef = `
{
	"name": "Subflow",
	"uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d",
	"nodes": [{
		"uuid": "d9dba561-b5ee-4f62-ba44-60c4dc242b84",
		"actions": []
	}]
}
`

var contactDef = `
{
	"name": "Ryan Lewis",
	"uuid": "5d76d86b-3bb9-4d5a-b822-c9d86f5d8e4f",
	"urns": ["tel:%2B12065551212", "email:foo@bar.com"],
	"groups": [{
		"uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d",
		"name": "Registered Users"
	}],
	"fields": {
		"activation_token": {
			"field_uuid": "ee46f9c4-b094-4e1b-ab0d-d4e65b4a99f1",
			"field_name": "Activation Token",
			"value": "XFW-JEV-9QE",
			"created_on": "2017-05-24T11:31:15.035757258-05:00"
		}
	}
}
`

func eventsForAction(actionJSON []byte) (json.RawMessage, error) {
	flowDef := fmt.Sprintf(flowDef, actionJSON)

	// read our flow
	flow, err := definition.ReadFlow(json.RawMessage(flowDef))
	if err != nil {
		return nil, err
	}

	// and our subflow
	subflow, err := definition.ReadFlow(json.RawMessage(subflowDef))
	if err != nil {
		return nil, err
	}

	// create our contact
	contact, err := flows.ReadContact(json.RawMessage(contactDef))
	if err != nil {
		return nil, err
	}

	// start our flow
	env := engine.NewFlowEnvironment(utils.NewDefaultEnvironment(), []flows.Flow{subflow}, []flows.FlowRun{}, []*flows.Contact{})
	output, err := engine.StartFlow(env, flow, contact, nil, nil, nil)
	if err != nil {
		return nil, err
	}

	events := output.Events()
	eventJSON := make([]json.RawMessage, len(events))
	for i, event := range events {
		typed, err := utils.EnvelopeFromTyped(event)
		if err != nil {
			return nil, err
		}
		eventJSON[i], err = json.MarshalIndent(typed, "", "  ")
		if err != nil {
			return nil, err
		}
	}
	if len(events) == 1 {
		return eventJSON[0], err
	}
	js, err := json.MarshalIndent(eventJSON, "", "  ")
	if err != nil {
		return nil, err
	}
	return js, nil
}

func handleActionDoc(prefix string, typeName string, docString string) {
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
	err = utils.ValidateAll(action)
	if err != nil {
		log.Fatalf("unable to validate example: %s\nHad err: %s", exampleJSON, err)
	}

	// make sure types match
	if name != action.Type() {
		log.Fatalf("Mismatched types for example of %s", name)
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

		fmt.Printf("# %s\n\n", name)
		fmt.Printf("%s", strings.Join(docs, "\n"))
		if len(example) > 0 {
			fmt.Printf("```json\n")
			fmt.Printf("%s\n", exampleJSON)
			fmt.Printf("```\n")

			fmt.Printf("```json\n")
			fmt.Printf("%s\n", events)
			fmt.Printf("```\n")
		}
		fmt.Printf("\n")
	}
}
