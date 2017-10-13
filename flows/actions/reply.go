package actions

import (
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
	BaseMsgAction
	AllURNs bool `json:"all_urns,omitempty"`
}

// Type returns the type of this action
func (a *ReplyAction) Type() string { return TypeReply }

// Validate validates whether this struct is correct
func (a *ReplyAction) Validate(assets flows.SessionAssets) error {
	return nil
}

// Execute runs this action
func (a *ReplyAction) Execute(run flows.FlowRun, step flows.Step, log flows.ActionLog) error {
	evaluatedText, evaluatedAttachments := a.evaluateMessage(run, step, log)

	urns := run.Contact().URNs()
	if a.AllURNs && len(urns) > 0 {
		for _, urn := range urns {
			log.Add(events.NewSendMsgToURN(urn, evaluatedText, evaluatedAttachments))
		}
	} else {
		log.Add(events.NewSendMsgToContact(run.Contact().Reference(), evaluatedText, evaluatedAttachments))
	}

	return nil
}
