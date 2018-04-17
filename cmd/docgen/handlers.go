package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/nyaruka/goflow/excellent/functions"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/actions"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/test"
	"github.com/nyaruka/goflow/utils"
)

func handleContextDoc(output *strings.Builder, item *documentedItem, session flows.Session) error {
	if len(item.examples) == 0 {
		return fmt.Errorf("no examples found for context item %s/%s", item.tagValue, item.typeName)
	}

	// check the examples
	for _, ex := range item.examples {
		if err := checkExample(session, ex); err != nil {
			return err
		}
	}

	exampleBlock := strings.Replace(strings.Join(item.examples, "\n"), "->", "→", -1)

	output.WriteString(fmt.Sprintf("<a name=\"context:%s\"></a>\n\n", item.tagValue))
	output.WriteString(fmt.Sprintf("## %s\n\n", strings.Title(item.tagValue)))
	output.WriteString(strings.Join(item.description, "\n"))
	output.WriteString("\n")
	output.WriteString("```objectivec\n")
	output.WriteString(exampleBlock)
	output.WriteString("\n")
	output.WriteString("```\n")
	output.WriteString("\n")
	return nil
}

func handleFunctionDoc(output *strings.Builder, item *documentedItem, session flows.Session) error {
	if len(item.examples) == 0 {
		return fmt.Errorf("no examples found for function %s", item.tagValue)
	}

	// get name of function from signature to use as our anchor
	name := item.tagValue[0:strings.Index(item.tagValue, "(")]

	// check the function name is a registered function
	_, exists := functions.XFUNCTIONS[name]
	if !exists {
		return fmt.Errorf("docstring function tag %s isn't a registered function", item.tagValue)
	}

	// check the examples
	for _, l := range item.examples {
		if err := checkExample(session, l); err != nil {
			return err
		}
	}

	exampleBlock := strings.Replace(strings.Join(item.examples, "\n"), "->", "→", -1)

	output.WriteString(fmt.Sprintf("<a name=\"%s:%s\"></a>\n\n", item.tagName, name))
	output.WriteString(fmt.Sprintf("## %s\n\n", item.tagValue))
	output.WriteString(strings.Join(item.description, "\n"))
	output.WriteString("\n")
	output.WriteString("```objectivec\n")
	output.WriteString(exampleBlock)
	output.WriteString("\n")
	output.WriteString("```\n")
	output.WriteString("\n")
	return nil
}

func handleEventDoc(output *strings.Builder, item *documentedItem, session flows.Session) error {
	// try to parse our example
	exampleJSON := []byte(strings.Join(item.examples, "\n"))
	typed := &utils.TypedEnvelope{}
	err := json.Unmarshal(exampleJSON, typed)
	if err != nil {
		return fmt.Errorf("unable to parse example: %s", err)
	}

	event, err := events.EventFromEnvelope(typed)
	if err != nil {
		return fmt.Errorf("unable to parse example: %s", err)
	}

	// validate it
	err = utils.Validate(event)
	if err != nil {
		return fmt.Errorf("unable to validate example: %s", err)
	}

	typed, err = utils.EnvelopeFromTyped(event)
	if err != nil {
		return fmt.Errorf("unable to marshal example: %s", err)
	}
	exampleJSON, err = json.MarshalIndent(typed, "", "    ")
	if err != nil {
		return fmt.Errorf("unable to marshal example: %s", err)
	}

	output.WriteString(fmt.Sprintf("<a name=\"event:%s\"></a>\n\n", item.tagValue))
	output.WriteString(fmt.Sprintf("## %s\n\n", item.tagValue))
	output.WriteString(strings.Join(item.description, "\n"))

	output.WriteString(`<div class="output_event"><h3>Event</h3>`)
	output.WriteString("```json\n")
	output.WriteString(fmt.Sprintf("%s\n", exampleJSON))
	output.WriteString("```\n")
	output.WriteString(`</div>`)

	output.WriteString("\n")

	return nil
}

func handleActionDoc(output *strings.Builder, item *documentedItem, session flows.Session) error {
	// try to parse our example
	exampleJSON := []byte(strings.Join(item.examples, "\n"))
	typed := &utils.TypedEnvelope{}
	err := json.Unmarshal(exampleJSON, typed)
	action, err := actions.ActionFromEnvelope(typed)
	if err != nil {
		return fmt.Errorf("unable to parse example: %s", err)
	}

	// validate it
	err = utils.Validate(action)
	if err != nil {
		return fmt.Errorf("unable to validate example: %s", err)
	}

	typed, err = utils.EnvelopeFromTyped(action)
	if err != nil {
		return fmt.Errorf("unable to marshal example: %s", err)
	}

	exampleJSON, err = json.MarshalIndent(typed, "", "  ")
	if err != nil {
		return fmt.Errorf("unable to marshal example: %s", err)
	}

	// get the events created by this action
	events, err := eventsForAction(action)
	if err != nil {
		return fmt.Errorf("error running action %s", err)
	}

	output.WriteString(fmt.Sprintf("<a name=\"action:%s\"></a>\n\n", item.tagValue))
	output.WriteString(fmt.Sprintf("## %s\n\n", item.tagValue))
	output.WriteString(strings.Join(item.description, "\n"))

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
	output.WriteString("\n")

	return nil
}

func checkExample(session flows.Session, line string) error {
	pieces := strings.Split(line, "->")
	if len(pieces) != 2 {
		return fmt.Errorf("unparseable example: %s", line)
	}

	test := strings.TrimSpace(pieces[0])
	expected := strings.Replace(strings.TrimSpace(pieces[1]), "\\n", "\n", -1)

	// evaluate our expression
	val, err := session.Runs()[0].EvaluateTemplateAsString(test, false)

	if expected == "ERROR" {
		if err == nil {
			return fmt.Errorf("expected example '%s' to error but it didn't", test)
		}
	} else if val != expected {
		return fmt.Errorf("expected '%s' from example: '%s', but got '%s'", expected, test, val)
	}

	return nil
}

func eventsForAction(action flows.Action) (json.RawMessage, error) {
	session, err := test.CreateTestSession(49998, action)
	if err != nil {
		return nil, err
	}

	path := session.Runs()[0].Path()
	lastStep := path[len(path)-1]

	// only interested in events created on the last step
	eventLog := make([]flows.Event, 0)
	for _, event := range session.Events() {
		if event.StepUUID() == lastStep.UUID() {
			eventLog = append(eventLog, event)
		}
	}

	eventJSON := make([]json.RawMessage, len(eventLog))
	for i, event := range eventLog {
		// action examples aren't supposed to generate error events - if they have, something went wrong
		if event.Type() == events.TypeError {
			errEvent := event.(*events.ErrorEvent)
			return nil, fmt.Errorf("error event generated: %s", errEvent.Text)
		}

		// give all our example events a fixed created on time
		event.SetCreatedOn(session.Environment().Now())

		typed, err := utils.EnvelopeFromTyped(event)
		if err != nil {
			return nil, err
		}
		eventJSON[i], err = json.MarshalIndent(typed, "", "    ")
		if err != nil {
			return nil, err
		}
	}
	if len(eventLog) == 1 {
		return eventJSON[0], err
	}
	js, err := json.MarshalIndent(eventJSON, "", "    ")
	if err != nil {
		return nil, err
	}
	return js, nil
}
