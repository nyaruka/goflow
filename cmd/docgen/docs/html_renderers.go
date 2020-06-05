package docs

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/static"
	"github.com/nyaruka/goflow/cmd/docgen/completion"
	"github.com/nyaruka/goflow/excellent/functions"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/actions"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/flows/resumes"
	"github.com/nyaruka/goflow/flows/triggers"
	"github.com/nyaruka/goflow/utils"
	"github.com/nyaruka/goflow/utils/jsonx"

	"github.com/pkg/errors"
)

var dynamicContextTypes = []string{"fields", "globals", "results", "urns"}

// function that can render a single tagged item
type renderFunc func(*strings.Builder, *TaggedItem, flows.Session, flows.Session) error

func init() {
	registerContextFunc(createItemListContextFunc("type", renderTypeDoc))
	registerContextFunc(createItemListContextFunc("operator", renderOperatorDoc))
	registerContextFunc(createItemListContextFunc("function", renderFunctionDoc))
	registerContextFunc(createItemListContextFunc("asset", renderAssetDoc))
	registerContextFunc(createItemListContextFunc("context", renderContextDoc))
	registerContextFunc(createItemListContextFunc("test", renderFunctionDoc))
	registerContextFunc(createItemListContextFunc("action", renderActionDoc))
	registerContextFunc(createItemListContextFunc("event", renderEventDoc))
	registerContextFunc(createItemListContextFunc("trigger", renderTriggerDoc))
	registerContextFunc(createItemListContextFunc("resume", renderResumeDoc))
	registerContextFunc(renderRootContext)
}

// creates a context function that renders all tagged items of a given type as a list
func createItemListContextFunc(tag string, renderer renderFunc) ContextFunc {
	return func(items map[string][]*TaggedItem, session flows.Session, voiceSession flows.Session) (map[string]string, error) {
		contextKey := fmt.Sprintf("%sDocs", tag)
		buffer := &strings.Builder{}

		for _, item := range items[tag] {
			if err := renderer(buffer, item, session, voiceSession); err != nil {
				return nil, errors.Wrapf(err, "error rendering %s:%s", item.tagName, item.tagValue)
			}
		}

		return map[string]string{contextKey: buffer.String()}, nil
	}
}

func renderAssetDoc(output *strings.Builder, item *TaggedItem, session flows.Session, voiceSession flows.Session) error {
	if len(item.examples) == 0 {
		return errors.Errorf("no examples found for asset item %s/%s", item.tagValue, item.typeName)
	}

	marshaled, err := jsonx.MarshalPretty(json.RawMessage(strings.Join(item.examples, "\n")))
	if err != nil {
		return errors.Wrap(err, "unable to marshal example")
	}

	// try to load example as part of a static asset source
	var assetSet string
	if item.typeName == "flow" {
		assetSet = fmt.Sprintf(`{"flow": %s}`, string(marshaled))
	} else {
		assetSet = fmt.Sprintf(`{"%ss": [%s]}`, item.tagValue, string(marshaled))
	}

	_, err = static.NewSource([]byte(assetSet))
	if err != nil {
		return errors.Wrap(err, "unable to load example into asset source")
	}

	output.WriteString(renderItemTitle(item))
	output.WriteString(strings.Join(item.description, "\n"))
	output.WriteString("\n")
	output.WriteString("```objectivec\n")
	output.WriteString(string(marshaled))
	output.WriteString("\n")
	output.WriteString("```\n")
	output.WriteString("\n")
	return nil
}

func renderTypeDoc(output *strings.Builder, item *TaggedItem, session flows.Session, voiceSession flows.Session) error {
	if len(item.examples) == 0 {
		return errors.Errorf("no examples found for type %s/%s", item.tagValue, item.typeName)
	}

	// check the examples
	for _, ex := range item.examples {
		if err := checkExample(session, ex); err != nil {
			return err
		}
	}

	output.WriteString(renderItemTitle(item))
	output.WriteString(strings.Join(item.description, "\n"))
	output.WriteString("\n")
	output.WriteString("```objectivec\n")
	output.WriteString(strings.Join(item.examples, "\n"))
	output.WriteString("\n")
	output.WriteString("```\n")
	output.WriteString("\n")
	return nil
}

func renderOperatorDoc(output *strings.Builder, item *TaggedItem, session flows.Session, voiceSession flows.Session) error {
	if len(item.examples) == 0 {
		return errors.Errorf("no examples found for operator %s/%s", item.tagValue, item.typeName)
	}

	// check the examples
	for _, ex := range item.examples {
		if err := checkExample(session, ex); err != nil {
			return err
		}
	}

	output.WriteString(renderItemTitle(item))
	output.WriteString(strings.Join(item.description, "\n"))
	output.WriteString("\n")
	output.WriteString("```objectivec\n")
	output.WriteString(strings.Join(item.examples, "\n"))
	output.WriteString("\n")
	output.WriteString("```\n")
	output.WriteString("\n")
	return nil
}

func renderContextDoc(output *strings.Builder, item *TaggedItem, session flows.Session, voiceSession flows.Session) error {
	// root of context is rendered separately by renderRootContext
	if item.tagValue == "root" {
		return nil
	}

	// examples are actually auto-completion property descriptors
	var defaultProp *completion.Property
	properties := make([]*completion.Property, 0, len(item.examples))
	for _, propDesc := range item.examples {
		prop := completion.ParseProperty(propDesc)
		if prop == nil {
			return errors.Errorf("invalid format for property description \"%s\"", propDesc)
		}
		if prop.Key == "__default__" {
			defaultProp = prop
		} else {
			properties = append(properties, prop)
		}
	}

	output.WriteString(renderItemTitle(item))

	if defaultProp != nil {
		output.WriteString(fmt.Sprintf("Defaults to %s (%s)\n\n", defaultProp.Help, renderPropertyType(defaultProp)))
	}

	for _, p := range properties {
		output.WriteString(fmt.Sprintf(" * `%s` %s (%s)\n", p.Key, p.Help, renderPropertyType(p)))
	}
	output.WriteString("\n")
	return nil
}

func renderRootContext(items map[string][]*TaggedItem, session flows.Session, voiceSession flows.Session) (map[string]string, error) {
	var root *TaggedItem
	for _, item := range items["context"] {
		if item.tagValue == "root" {
			root = item
			break
		}
	}

	// examples are actually auto-completion property descriptors
	properties := make([]*completion.Property, 0, len(root.examples))
	for _, propDesc := range root.examples {
		prop := completion.ParseProperty(propDesc)
		if prop == nil {
			return nil, errors.Errorf("invalid format for property description \"%s\"", propDesc)
		}
		properties = append(properties, prop)
	}

	output := &strings.Builder{}
	for _, p := range properties {
		output.WriteString(fmt.Sprintf(" * `%s` %s (%s)\n", p.Key, p.Help, renderPropertyType(p)))
	}
	output.WriteString("\n")

	return map[string]string{"contextRoot": output.String()}, nil
}

func renderPropertyType(p *completion.Property) string {
	if p.Type == "any" || utils.StringSliceContains(dynamicContextTypes, p.Type, true) {
		return p.Type
	} else if p.Type == "text" || p.Type == "number" || p.Type == "datetime" {
		return fmt.Sprintf("[type:%s]", p.Type)
	}
	return fmt.Sprintf("[context:%s]", p.Type)
}

func renderFunctionDoc(output *strings.Builder, item *TaggedItem, session flows.Session, voiceSession flows.Session) error {
	if len(item.examples) == 0 {
		return errors.Errorf("no examples found for function %s", item.tagValue)
	}

	// check the function name is a registered function
	_, exists := functions.XFUNCTIONS[item.tagValue]
	if !exists {
		return errors.Errorf("docstring function tag %s isn't a registered function", item.tagValue)
	}

	// check the examples
	for _, l := range item.examples {
		if err := checkExample(session, l); err != nil {
			return err
		}
	}

	output.WriteString(renderItemTitle(item))
	output.WriteString(strings.Join(item.description, "\n"))
	output.WriteString("\n")
	output.WriteString("```objectivec\n")
	output.WriteString(strings.Join(item.examples, "\n"))
	output.WriteString("\n")
	output.WriteString("```\n")
	output.WriteString("\n")
	return nil
}

func renderEventDoc(output *strings.Builder, item *TaggedItem, session flows.Session, voiceSession flows.Session) error {
	// try to parse our example
	exampleJSON := []byte(strings.Join(item.examples, "\n"))
	event, err := events.ReadEvent(exampleJSON)
	if err != nil {
		return errors.Wrap(err, "unable to read event")
	}

	// validate it
	err = utils.Validate(event)
	if err != nil {
		return errors.Wrap(err, "unable to validate example")
	}

	exampleJSON, err = jsonx.MarshalPretty(event)
	if err != nil {
		return errors.Wrap(err, "unable to marshal example")
	}

	output.WriteString(renderItemTitle(item))
	output.WriteString(strings.Join(item.description, "\n"))

	output.WriteString(`<div class="output_event">`)
	output.WriteString("\n\n")
	output.WriteString("```json\n")
	output.WriteString(fmt.Sprintf("%s\n", exampleJSON))
	output.WriteString("```\n")
	output.WriteString(`</div>`)

	output.WriteString("\n")

	return nil
}

func renderActionDoc(output *strings.Builder, item *TaggedItem, session flows.Session, voiceSession flows.Session) error {
	// try to parse our example
	exampleJSON := []byte(strings.Join(item.examples, "\n"))
	action, err := actions.ReadAction(exampleJSON)
	if err != nil {
		return errors.Wrap(err, "unable to read action")
	}

	// validate it
	err = utils.Validate(action)
	if err != nil {
		return errors.Wrap(err, "unable to validate example")
	}

	exampleJSON, err = jsonx.MarshalPretty(action)
	if err != nil {
		return errors.Wrap(err, "unable to marshal example")
	}

	// get the events created by this action
	events, err := eventsForAction(action, session, voiceSession)
	if err != nil {
		return errors.Wrap(err, "error running action")
	}

	output.WriteString(renderItemTitle(item))
	output.WriteString(strings.Join(item.description, "\n"))

	output.WriteString(`<div class="input_action"><h3>Action</h3>`)
	output.WriteString("\n\n")
	output.WriteString("```json\n")
	output.WriteString(fmt.Sprintf("%s\n", exampleJSON))
	output.WriteString("```\n")
	output.WriteString(`</div>`)

	output.WriteString(`<div class="output_event"><h3>Event</h3>`)
	output.WriteString("\n\n")
	output.WriteString("```json\n")
	output.WriteString(fmt.Sprintf("%s\n", events))
	output.WriteString("```\n")
	output.WriteString(`</div>`)
	output.WriteString("\n")

	return nil
}

func renderTriggerDoc(output *strings.Builder, item *TaggedItem, session flows.Session, voiceSession flows.Session) error {
	// try to parse our example
	exampleJSON := json.RawMessage(strings.Join(item.examples, "\n"))
	trigger, err := triggers.ReadTrigger(session.Assets(), exampleJSON, assets.PanicOnMissing)
	if err != nil {
		return errors.Wrap(err, "unable to read trigger")
	}

	// validate it
	err = utils.Validate(trigger)
	if err != nil {
		return errors.Wrap(err, "unable to validate example")
	}

	exampleJSON, err = jsonx.MarshalPretty(trigger)
	if err != nil {
		return errors.Wrap(err, "unable to marshal example")
	}

	output.WriteString(renderItemTitle(item))
	output.WriteString(strings.Join(item.description, "\n"))
	output.WriteString("\n")
	output.WriteString("```json\n")
	output.WriteString(fmt.Sprintf("%s\n", exampleJSON))
	output.WriteString("```\n")
	output.WriteString("\n")

	return nil
}

func renderResumeDoc(output *strings.Builder, item *TaggedItem, session flows.Session, voiceSession flows.Session) error {
	// try to parse our example
	exampleJSON := json.RawMessage(strings.Join(item.examples, "\n"))
	resume, err := resumes.ReadResume(session.Assets(), exampleJSON, assets.PanicOnMissing)
	if err != nil {
		return errors.Wrap(err, "unable to read resume")
	}

	// validate it
	if err := utils.Validate(resume); err != nil {
		return errors.Wrap(err, "unable to validate example")
	}

	exampleJSON, err = jsonx.MarshalPretty(resume)
	if err != nil {
		return errors.Wrap(err, "unable to marshal example")
	}

	output.WriteString(renderItemTitle(item))
	output.WriteString(strings.Join(item.description, "\n"))
	output.WriteString("\n")
	output.WriteString("```json\n")
	output.WriteString(fmt.Sprintf("%s\n", exampleJSON))
	output.WriteString("```\n")
	output.WriteString("\n")

	return nil
}

func renderItemTitle(item *TaggedItem) string {
	return fmt.Sprintf("<h2 class=\"item_title\"><a name=\"%[1]s:%[2]s\" href=\"#%[1]s:%[2]s\">%[3]s</a></h2>\n\n", item.tagName, item.tagValue, item.tagTitle)
}

func checkExample(session flows.Session, line string) error {
	pieces := strings.Split(line, "→")
	if len(pieces) != 2 {
		return errors.Errorf("unparseable example: %s", line)
	}

	test := strings.TrimSpace(pieces[0])
	expected := strings.Replace(strings.TrimSpace(pieces[1]), `\n`, "\n", -1)
	expected = strings.Replace(expected, `\x20`, " ", -1)

	// evaluate our expression
	val, err := session.Runs()[0].EvaluateTemplate(test)

	if expected == "ERROR" {
		if err == nil {
			return errors.Errorf("expected example '%s' to error but it didn't", strconv.Quote(test))
		}
	} else {
		if err != nil {
			return errors.Errorf("unexpected error from example '%s': %s", strconv.Quote(test), err)
		}
		if val != expected {
			return errors.Errorf("expected %s from example: %s, but got %s", strconv.Quote(expected), strconv.Quote(test), strconv.Quote(val))
		}
	}

	return nil
}

func eventsForAction(action flows.Action, msgSession flows.Session, voiceSession flows.Session) (json.RawMessage, error) {
	voiceAction := len(action.AllowedFlowTypes()) == 1 && action.AllowedFlowTypes()[0] == flows.FlowTypeVoice
	session := msgSession
	if voiceAction {
		session = voiceSession
	}

	run := session.Runs()[0]
	step := run.Path()[len(run.Path())-1]
	modifierLog := func(flows.Modifier) {}

	eventList := make([]flows.Event, 0)
	eventLog := func(e flows.Event) {
		e.SetStepUUID(step.UUID())
		eventList = append(eventList, e)
	}

	err := action.Execute(run, step, modifierLog, eventLog)
	if err != nil {
		return nil, err
	}

	eventJSON := make([]json.RawMessage, len(eventList))
	for i, event := range eventList {
		// action examples aren't supposed to generate error events - if they have, something went wrong
		if event.Type() == events.TypeError {
			errEvent := event.(*events.ErrorEvent)
			return nil, errors.Errorf("error event generated: %s", errEvent.Text)
		}

		eventJSON[i], err = jsonx.MarshalPretty(event)
		if err != nil {
			return nil, err
		}
	}
	if len(eventList) == 1 {
		return eventJSON[0], err
	}
	js, err := jsonx.MarshalPretty(eventJSON)
	if err != nil {
		return nil, err
	}
	return js, nil
}
