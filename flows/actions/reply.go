package actions

import (
	"github.com/nyaruka/goflow/excellent"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

// TypeReply is the type for reply actions
const TypeReply string = "reply"

// ReplyAction can be used to reply to the current contact in a flow. The text field may contain templates.
//
// A `msg_out` event will be created with the evaluated text.
//
// ```
//   {
//     "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//     "type": "reply",
//     "text": "Hi @contact.name, are you ready to complete today's survey?"
//   }
// ```
//
// @action reply
type ReplyAction struct {
	BaseAction
	Text string `json:"text"         validate:"required"`
}

// Type returns the type of this action
func (a *ReplyAction) Type() string { return TypeReply }

// Validate validates whether this struct is correct
func (a *ReplyAction) Validate() error {
	return utils.ValidateAll(a)
}

// Execute runs this action
func (a *ReplyAction) Execute(run flows.FlowRun, step flows.Step) error {
	text, err := excellent.EvaluateTemplateAsString(run.Environment(), run.Context(), run.GetText(flows.UUID(a.UUID), "text", a.Text))
	if err != nil {
		run.AddError(step, err)
	}
	run.AddEvent(step, events.NewMsgToContact(run.Contact().UUID(), text))
	return nil
}
