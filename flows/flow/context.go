package flow

import (
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
