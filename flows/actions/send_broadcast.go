package actions

import (
	"context"

	"github.com/nyaruka/gocommon/i18n"
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
)

func init() {
	registerType(TypeSendBroadcast, func() flows.Action { return &SendBroadcast{} })
}

// TypeSendBroadcast is the type for the send broadcast action
const TypeSendBroadcast string = "send_broadcast"

// SendBroadcast can be used to send a message to one or more contacts. It accepts a list of URNs, a list of groups
// and a list of contacts.
//
// The URNs and text fields may be templates. A [event:broadcast_created] event will be created for each unique urn, contact and group
// with the evaluated text.
//
//	{
//	  "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//	  "type": "send_broadcast",
//	  "urns": ["tel:+12065551212"],
//	  "text": "Hi @contact.name, are you ready to complete today's survey?",
//	  "template": {
//	    "uuid": "3ce100b7-a734-4b4e-891b-350b1279ade2",
//	    "name": "revive_issue"
//	  },
//	  "template_variables": ["@contact.name"]
//	}
//
// @action send_broadcast
type SendBroadcast struct {
	baseAction
	onlineAction
	otherContactsAction
	createMsgAction
}

// NewSendBroadcast creates a new send broadcast action
func NewSendBroadcast(uuid flows.ActionUUID, text string, attachments []string, quickReplies []string, groups []*assets.GroupReference, contacts []*flows.ContactReference, contactQuery string, urns []urns.URN, legacyVars []string, template *assets.TemplateReference, templateVariables []string) *SendBroadcast {
	return &SendBroadcast{
		baseAction: newBaseAction(TypeSendBroadcast, uuid),
		otherContactsAction: otherContactsAction{
			Groups:       groups,
			Contacts:     contacts,
			ContactQuery: contactQuery,
			URNs:         urns,
			LegacyVars:   legacyVars,
		},
		createMsgAction: createMsgAction{
			Text:              text,
			Attachments:       attachments,
			QuickReplies:      quickReplies,
			Template:          template,
			TemplateVariables: templateVariables,
		},
	}
}

// Execute runs this action
func (a *SendBroadcast) Execute(ctx context.Context, run flows.Run, step flows.Step, log flows.EventLogger) error {
	groupRefs, contactRefs, contactQuery, urnList, err := a.resolveRecipients(run, log)
	if err != nil {
		return err
	}

	// footgun prevention
	if run.Session().BatchStart() && (len(groupRefs) > 0 || contactQuery != "") {
		log(events.NewError("Can't send broadcasts to groups during batch starts", ""))
		return nil
	}

	translations := make(flows.BroadcastTranslations)
	languages := append([]i18n.Language{run.Flow().Language()}, run.Flow().Localization().Languages()...)

	// evaluate the broadcast in each language we have translations for
	for _, language := range languages {
		languages := []i18n.Language{language, run.Flow().Language()}

		content, _ := a.evaluateMessage(run, languages, a.Text, a.Attachments, a.QuickReplies, log)
		translations[language] = content
	}

	// if we don't have any recipients, noop
	if !(len(urnList) > 0 || len(groupRefs) > 0 || len(contactRefs) > 0 || a.ContactQuery != "") {
		return nil
	}

	// evaluate template variables if we have a template
	var templateRef *assets.TemplateReference
	var templateVariables []string
	if a.Template != nil {
		sa := run.Session().Assets()
		template := sa.Templates().Get(a.Template.UUID)

		if template != nil {
			templateRef = a.Template
			templateVariables = make([]string, len(a.TemplateVariables))
			for i, varExp := range a.TemplateVariables {
				v, _ := run.EvaluateTemplate(varExp, log)
				templateVariables[i] = v
			}
		}
	}

	log(events.NewBroadcastCreated(translations, run.Flow().Language(), groupRefs, contactRefs, contactQuery, urnList, templateRef, templateVariables))
	return nil
}

func (a *SendBroadcast) Inspect(dependency func(assets.Reference), local func(string), result func(*flows.ResultInfo)) {
	a.otherContactsAction.Inspect(dependency, local, result)
	a.createMsgAction.Inspect(dependency, local, result)
}
