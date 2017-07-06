package actions

import (
	"fmt"

	"github.com/nyaruka/goflow/excellent"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

// TypeSendMsg is the type for msg actions
const TypeSendMsg string = "send_msg"

// SendMsgAction can be used to send a message to one or more contacts. It accepts a list of URNs, a list of groups
// and a list of contacts.
//
// The URNs and text fields may be templates. A `send_msg` event will be created for each unique urn, contact and group
// with the evaluated text.
//
// ```
//   {
//     "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//     "type": "send_msg",
//     "urns": ["tel:+12065551212"],
//     "text": "Hi @contact.name, are you ready to complete today's survey?"
//   }
// ```
//
// @action send_msg
type SendMsgAction struct {
	BaseAction
	URNs        []flows.URN               `json:"urns"`
	Contacts    []*flows.ContactReference `json:"contacts"     validate:"dive"`
	Groups      []*flows.Group            `json:"groups"       validate:"dive"`
	Text        string                    `json:"text"`
	Attachments []string                  `json:"attachments"`
}

// Type returns the type of this action
func (a *SendMsgAction) Type() string { return TypeSendMsg }

// Validate validates whether this struct is correct
func (a *SendMsgAction) Validate() error {
	return utils.ValidateAll(a)
}

// Execute runs this action
func (a *SendMsgAction) Execute(run flows.FlowRun, step flows.Step) error {
	// TODO: customize this for receiving contacts instead of one global replace
	text, err := excellent.EvaluateTemplateAsString(run.Environment(), run.Context(), run.GetText(flows.UUID(a.UUID), "text", a.Text))
	if err != nil {
		run.AddError(step, err)
	}
	if text == "" {
		run.AddError(step, fmt.Errorf("send_msg text evaluated to empty string, skipping"))
		return nil
	}

	attachments := a.Attachments
	if attachments == nil {
		attachments = []string{}
	}

	// create events for each URN
	for _, urn := range a.URNs {
		run.AddEvent(step, events.NewSendMsgToURN(urn, text, attachments))
	}

	for _, contact := range a.Contacts {
		run.AddEvent(step, events.NewSendMsgToContact(contact.UUID, text, attachments))
	}

	for _, group := range a.Groups {
		run.AddEvent(step, events.NewSendMsgToGroup(group.UUID(), text, attachments))
	}
	return nil
}
