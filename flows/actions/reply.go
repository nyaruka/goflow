package actions

import (
	"fmt"

	"github.com/nyaruka/goflow/excellent"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
)

// TypeReply is the type for reply actions
const TypeReply string = "reply"

// ReplyAction can be used to reply to the current contact in a flow. The text field may contain templates.
//
// A `send_msg` event will be created with the evaluated text.
//
// ```
//   {
//     "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//     "type": "reply",
//     "text": "Hi @contact.name, are you ready to complete today's survey?",
//     "attachments": [],
//     "all_urns": false
//   }
// ```
//
// @action reply
type ReplyAction struct {
	BaseAction
	Text        string   `json:"text"         validate:"required"`
	Attachments []string `json:"attachments"`
	AllURNs     bool     `json:"all_urns,omitempty"`
}

// Type returns the type of this action
func (a *ReplyAction) Type() string { return TypeReply }

// Validate validates whether this struct is correct
func (a *ReplyAction) Validate(assets flows.AssetStore) error {
	return nil
}

// Execute runs this action
func (a *ReplyAction) Execute(run flows.FlowRun, step flows.Step) error {

	text, err := excellent.EvaluateTemplateAsString(run.Environment(), run.Context(), run.GetText(flows.UUID(a.UUID()), "text", a.Text))
	if err != nil {
		run.AddError(step, a, err)
	}
	if text == "" {
		run.AddError(step, a, fmt.Errorf("reply text evaluated to empty string, skipping"))
		return nil
	}

	urns := run.Contact().URNs()
	if a.AllURNs && len(urns) > 0 {
		for _, urn := range urns {
			run.ApplyEvent(step, a, events.NewSendMsgToURN(urn, text, a.Attachments))
		}
	} else {
		run.ApplyEvent(step, a, events.NewSendMsgToContact(run.Contact().UUID(), text, a.Attachments))
	}

	return nil

}
