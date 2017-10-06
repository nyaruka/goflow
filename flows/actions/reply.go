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
func (a *ReplyAction) Validate(assets flows.SessionAssets) error {
	return nil
}

// Execute runs this action
func (a *ReplyAction) Execute(run flows.FlowRun, step flows.Step) ([]flows.Event, error) {
	log := make([]flows.Event, 0)

	text, err := excellent.EvaluateTemplateAsString(run.Environment(), run.Context(), run.GetText(flows.UUID(a.UUID()), "text", a.Text))
	if err != nil {
		log = append(log, events.NewErrorEvent(err))
	}
	if text == "" {
		log = append(log, events.NewErrorEvent(fmt.Errorf("reply text evaluated to empty string, skipping")))
		return log, nil
	}

	urns := run.Contact().URNs()
	if a.AllURNs && len(urns) > 0 {
		for _, urn := range urns {
			log = append(log, events.NewSendMsgToURN(urn, text, a.Attachments))
		}
	} else {
		log = append(log, events.NewSendMsgToContact(run.Contact().Reference(), text, a.Attachments))
	}

	return log, nil

}
