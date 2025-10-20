package actions

import (
	"context"
	"fmt"

	"github.com/nyaruka/gocommon/i18n"
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
)

func init() {
	registerType(TypeSendMsg, func() flows.Action { return &SendMsg{} })
}

// TypeSendMsg is the type for the send message action
const TypeSendMsg string = "send_msg"

// SendMsg can be used to reply to the current contact in a flow. The text field may contain templates. The action
// will attempt to find pairs of URNs and channels which can be used for sending. If it can't find such a pair, it will
// create a message without a channel or URN.
//
// A [event:msg_created] event will be created with the evaluated text.
//
//	{
//	  "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//	  "type": "send_msg",
//	  "text": "Hi @contact.name, are you ready to complete today's survey?",
//	  "attachments": [],
//	  "all_urns": false,
//	  "template": {
//	    "uuid": "3ce100b7-a734-4b4e-891b-350b1279ade2",
//	    "name": "revive_issue"
//	  },
//	  "template_variables": ["@contact.name"]
//	}
//
// @action send_msg
type SendMsg struct {
	baseAction
	universalAction
	createMsgAction

	AllURNs           bool                      `json:"all_urns,omitempty"`
	Template          *assets.TemplateReference `json:"template,omitempty"`
	TemplateVariables []string                  `json:"template_variables,omitempty" engine:"localized,evaluated"`
}

// NewSendMsg creates a new send msg action
func NewSendMsg(uuid flows.ActionUUID, text string, attachments []string, quickReplies []string, allURNs bool) *SendMsg {
	return &SendMsg{
		baseAction: newBaseAction(TypeSendMsg, uuid),
		createMsgAction: createMsgAction{
			Text:         text,
			Attachments:  attachments,
			QuickReplies: quickReplies,
		},
		AllURNs: allURNs,
	}
}

// Execute runs this action
func (a *SendMsg) Execute(ctx context.Context, run flows.Run, step flows.Step, log flows.EventLogger) error {
	content, lang := a.evaluateMessage(run, nil, a.Text, a.Attachments, a.QuickReplies, log)

	// determine if this message can be sent - unsendable messages are still created for history's sake
	unsendableReason := flows.NilUnsendableReason
	contactStatus := run.Contact().Status()
	if contactStatus == flows.ContactStatusBlocked {
		unsendableReason = flows.UnsendableReasonContactBlocked
	} else if contactStatus == flows.ContactStatusStopped {
		unsendableReason = flows.UnsendableReasonContactStopped
	} else if contactStatus == flows.ContactStatusArchived {
		unsendableReason = flows.UnsendableReasonContactArchived
	} else {
		var err error
		unsendableReason, err = run.Session().Engine().Options().CheckSendable(run.Session().Assets(), run.Contact(), content)
		if err != nil {
			return fmt.Errorf("error checking if message is sendable: %w", err)
		}
	}

	sa := run.Session().Assets()

	var template *flows.Template
	var templateVariables []string
	if a.Template != nil {
		template = sa.Templates().Get(a.Template.UUID)

		if template != nil {
			templateVariables = make([]string, len(a.TemplateVariables))
			for i, varExp := range a.TemplateVariables {
				v, _ := run.EvaluateTemplate(varExp, log)
				templateVariables[i] = v
			}
		}
	}

	destinations := run.Contact().ResolveDestinations(a.AllURNs)
	locale := currentLocale(run, lang)

	// create a new message for each URN+channel destination
	for _, dest := range destinations {
		urn := dest.URN.URN()
		channelRef := assets.NewChannelReference(dest.Channel.UUID(), dest.Channel.Name())
		var msg *flows.MsgOut

		if template != nil {
			locales := []i18n.Locale{run.Session().MergedEnvironment().DefaultLocale(), run.Session().Environment().DefaultLocale()}
			translation := template.FindTranslation(dest.Channel, locales)
			if translation != nil {
				templating := template.Templating(translation, templateVariables)

				// the message we return is an approximate preview of what the channel will send using the template
				preview := translation.Preview(templating.Variables)
				locale := translation.Locale()

				msg = flows.NewMsgOut(urn, channelRef, preview, templating, locale, unsendableReason)
			}
		}

		if msg == nil {
			msg = flows.NewMsgOut(urn, channelRef, content, nil, locale, unsendableReason)
		}

		log(events.NewMsgCreated(msg, "", ""))
	}

	// if we couldn't find a destination, create a msg without a URN or channel and it's up to the caller
	// to handle that as they want
	if len(destinations) == 0 {
		msg := flows.NewMsgOut(urns.NilURN, nil, content, nil, locale, flows.UnsendableReasonNoDestination)
		log(events.NewMsgCreated(msg, "", ""))
	}

	return nil
}

func (a *SendMsg) Inspect(dependency func(assets.Reference), local func(string), result func(*flows.ResultInfo)) {
	if a.Template != nil {
		dependency(a.Template)
	}
}
