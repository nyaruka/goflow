package flow

import (
	"encoding/json"
	"fmt"

	"github.com/nyaruka/goflow/flows"
)

// NewContextForContact creates a new context for the passed in contact
func NewContextForContact(contact flows.Contact, run flows.FlowRun) flows.Context {
	context := context{contact: contact, run: run}
	return &context
}

type context struct {
	contact flows.Contact
	run     flows.FlowRun
}

func (c *context) Run() flows.FlowRun { return c.run }

func (c *context) Contact() flows.Contact { return c.contact }

func (c *context) Validate() error {
	// TODO: do some validation here
	return nil
}

func (c *context) Resolve(key string) interface{} {
	switch key {

	case "channel":
		return c.run.Channel()

	case "contact":
		return c.Contact()

	case "child":
		return c.run.Child()

	case "input":
		return c.run.Input()

	case "parent":
		return c.run.Parent()

	case "path":
		return c.run.Path()

	case "run":
		return c.run

	case "webhook":
		return c.run.Webhook()

	}

	return fmt.Errorf("No field '%s' on context", key)
}

func (c *context) Default() interface{} {
	return c
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

// ReadContext decodes a context from the passed in JSON
func ReadContext(data json.RawMessage) (flows.Context, error) {
	context := &context{}
	err := json.Unmarshal(data, context)
	if err == nil {
		err = context.Validate()
	}
	return context, err
}

type contextEnvelope struct {
	Contact *contact
	Run     *run
}

func (c *context) UnmarshalJSON(data []byte) error {
	var ce contextEnvelope
	var err error

	err = json.Unmarshal(data, &ce)
	if err != nil {
		return err
	}

	c.contact = ce.Contact
	c.run = ce.Run
	return err
}

func (c *context) MarshalJSON() ([]byte, error) {
	var ce contextEnvelope

	ce.Contact = c.contact.(*contact)
	ce.Run = c.run.(*run)
	return json.Marshal(ce)
}
