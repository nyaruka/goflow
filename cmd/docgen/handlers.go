package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/actions"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

func handleContextDoc(output *bytes.Buffer, tag string, typeName string, docString string, session flows.Session) error {
	parsed := parseDocString(docString, tag)
	if len(parsed.examples) == 0 {
		return fmt.Errorf("no examples found for context item %s", parsed.tagValue)
	}

	// check our examples
	for _, ex := range parsed.examples {
		if err := checkExample(session, ex); err != nil {
			return err
		}
	}

	exampleBlock := strings.Replace(strings.Join(parsed.examples, "\n"), "->", "→", -1)

	output.WriteString(fmt.Sprintf("<a name=\"context:%s\"></a>\n\n", parsed.tagValue))
	output.WriteString(fmt.Sprintf("## %s\n\n", strings.Title(parsed.tagValue)))
	output.WriteString(strings.Join(parsed.description, "\n"))
	output.WriteString("\n")
	output.WriteString("```objectivec\n")
	output.WriteString(exampleBlock)
	output.WriteString("\n")
	output.WriteString("```\n")
	output.WriteString("\n")
	return nil
}

func handleFunctionDoc(output *bytes.Buffer, tag string, typeName string, docString string, session flows.Session) error {
	parsed := parseDocString(docString, tag)
	if len(parsed.examples) == 0 {
		return fmt.Errorf("no examples found for function %s", parsed.tagValue)
	}

	name := parsed.tagValue[0:strings.Index(parsed.tagValue, "(")]

	if len(parsed.description) > 0 && strings.HasPrefix(parsed.description[0], typeName) {
		parsed.description[0] = strings.Replace(parsed.description[0], typeName, name, 1)
	}

	// check our examples
	for _, l := range parsed.examples {
		if err := checkExample(session, l); err != nil {
			return err
		}
	}

	exampleBlock := strings.Replace(strings.Join(parsed.examples, "\n"), "->", "→", -1)

	output.WriteString(fmt.Sprintf("<a name=\"%ss:%s\"></a>\n\n", tag[1:], name))
	output.WriteString(fmt.Sprintf("## %s\n\n", parsed.tagValue))
	output.WriteString(strings.Join(parsed.description, "\n"))
	output.WriteString("\n")
	output.WriteString("```objectivec\n")
	output.WriteString(exampleBlock)
	output.WriteString("\n")
	output.WriteString("```\n")
	output.WriteString("\n")
	return nil
}

func handleEventDoc(output *bytes.Buffer, prefix string, typeName string, docString string, session flows.Session) error {
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
	if err != nil {
		return fmt.Errorf("unable to parse example: %s\nHas err: %s", exampleJSON, err)
	}

	event, err := events.EventFromEnvelope(typed)
	if err != nil {
		return fmt.Errorf("unable to parse example: %s\nHas err: %s", exampleJSON, err)
	}

	// make sure types match
	if name != event.Type() {
		return fmt.Errorf("mismatched event types for example of %s", name)
	}

	// validate it
	err = utils.Validate(event)
	if err != nil {
		return fmt.Errorf("unable to validate example: %s\nHad err: %s", exampleJSON, err)
	}

	typed, err = utils.EnvelopeFromTyped(event)
	if err != nil {
		return fmt.Errorf("unable to marshal example: %s\nHad err: %s", exampleJSON, err)
	}
	exampleJSON, err = json.MarshalIndent(typed, "", "    ")
	if err != nil {
		return fmt.Errorf("unable to marshal example: %s\nHad err: %s", exampleJSON, err)
	}

	if name != "" {
		if len(docs) > 0 && strings.HasPrefix(docs[0], typeName) {
			docs[0] = strings.Replace(docs[0], typeName, name, 1)
		}

		output.WriteString(fmt.Sprintf("<a name=\"events:%s\"></a>\n\n", name))
		output.WriteString(fmt.Sprintf("## %s\n\n", name))
		output.WriteString(fmt.Sprintf("%s", strings.Join(docs, "\n")))
		if len(example) > 0 {
			output.WriteString(`<div class="output_event"><h3>Event</h3>`)
			output.WriteString("```json\n")
			output.WriteString(fmt.Sprintf("%s\n", exampleJSON))
			output.WriteString("```\n")
			output.WriteString(`</div>`)
		}
		output.WriteString(fmt.Sprintf("\n"))
	}
	return nil
}

func handleActionDoc(output *bytes.Buffer, prefix string, typeName string, docString string, session flows.Session) error {
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
		return fmt.Errorf("unable to parse example: %s: %s", exampleJSON, err)
	}

	// validate it
	err = utils.Validate(action)
	if err != nil {
		return fmt.Errorf("unable to validate example: %s: %s", exampleJSON, err)
	}

	// make sure types match
	if name != action.Type() {
		return fmt.Errorf("mismatched action types for example of %s", name)
	}

	typed, err = utils.EnvelopeFromTyped(action)
	if err != nil {
		return fmt.Errorf("unable to marshal example %s: %s", exampleJSON, err)
	}

	exampleJSON, err = json.MarshalIndent(typed, "", "  ")
	if err != nil {
		return fmt.Errorf("unable to marshal example %s: %s", exampleJSON, err)
	}

	// get the events created by this action
	events, err := eventsForAction(action)
	if err != nil {
		return fmt.Errorf("error running action %s: %s", exampleJSON, err)
	}

	if name != "" {
		if len(docs) > 0 && strings.HasPrefix(docs[0], typeName) {
			docs[0] = strings.Replace(docs[0], typeName, name, 1)
		}

		output.WriteString(fmt.Sprintf("<a name=\"actions:%s\"></a>\n\n", name))
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
	return nil
}

type parsedDocs struct {
	tagValue    string   // value after @tag
	examples    []string // any indented line
	description []string // any other line
}

func parseDocString(docString string, tag string) *parsedDocs {
	var tagValue string
	examples := make([]string, 0)
	description := make([]string, 0)

	for _, l := range strings.Split(docString, "\n") {
		trimmed := strings.TrimSpace(l)

		if strings.HasPrefix(l, tag) {
			tagValue = l[len(tag)+1:]
		} else if strings.HasPrefix(l, "  ") {
			examples = append(examples, trimmed)
		} else {
			description = append(description, l)
		}
	}

	return &parsedDocs{tagValue: tagValue, examples: examples, description: description}
}

func checkExample(session flows.Session, line string) error {
	pieces := strings.Split(line, "->")
	if len(pieces) != 2 {
		return fmt.Errorf("unparseable example: %s", line)
	}

	test, expected := strings.TrimSpace(pieces[0]), strings.TrimSpace(pieces[1])

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
	session, err := createExampleSession(action)
	if err != nil {
		return nil, err
	}

	// only interested in events after the new action
	eventLog := session.Events()[2:]

	eventJSON := make([]json.RawMessage, len(eventLog))
	for i, event := range eventLog {
		// action examples aren't supposed to generate error events - if they have, something went wrong
		if event.Type() == events.TypeError {
			errEvent := event.(*events.ErrorEvent)
			return nil, fmt.Errorf("error event generated: %s", errEvent.Text)
		}

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
