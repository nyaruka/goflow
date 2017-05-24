package actions

import (
	"github.com/nyaruka/goflow/excellent"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

// TypeMsg is the type for msg actions
const TypeMsg string = "msg"

// MsgAction can be used to send a message to one or more contacts. It accepts a list of URNs, a list of groups
// and a list of contacts.
//
// The URNs and text fields may be templates. A `msg_out` event will be created for each unique urn, contact and group
// with the evaluated text.
//
//   {
//     "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//     "type": "msg",
//     "urns": ["tel:+12065551212"],
//     "text": "Hi @contact.name, are you ready to complete today's survey?"
//   }
//
// @action msg
type MsgAction struct {
	BaseAction
	URNs     []flows.URN               `json:"urns"`
	Contacts []*flows.ContactReference `json:"contacts"     validate:"dive"`
	Groups   []*flows.Group            `json:"groups"       validate:"dive"`
	Text     string                    `json:"text"         validate:"required"`
}

// Type returns the type of this action
func (a *MsgAction) Type() string { return TypeMsg }

// Validate validates whether this struct is correct
func (a *MsgAction) Validate() error {
	return utils.ValidateAll(a)
}

// Execute runs this action
func (a *MsgAction) Execute(run flows.FlowRun, step flows.Step) error {
	// TODO: customize this for receiving contacts instead of one global replace
	text, err := excellent.EvaluateTemplateAsString(run.Environment(), run.Context(), run.GetText(flows.UUID(a.UUID), "text", a.Text))
	if err != nil {
		run.AddError(step, err)
	}

	// create events for each URN
	for _, urn := range a.URNs {
		run.AddEvent(step, events.NewMsgToURN(urn, text))
	}

	for _, contact := range a.Contacts {
		run.AddEvent(step, events.NewMsgToContact(contact.UUID, text))
	}

	for _, group := range a.Groups {
		run.AddEvent(step, events.NewMsgToGroup(group.UUID(), text))
	}
	return nil
}
