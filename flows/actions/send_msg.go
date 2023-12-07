package actions

import (
	"strings"

	"github.com/nyaruka/gocommon/i18n"
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	registerType(TypeSendMsg, func() flows.Action { return &SendMsgAction{} })
}

// TypeSendMsg is the type for the send message action
const TypeSendMsg string = "send_msg"

// SendMsgAction can be used to reply to the current contact in a flow. The text field may contain templates. The action
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
//	  "templating": {
//	    "uuid": "32c2ead6-3fa3-4402-8e27-9cc718175c5a",
//	    "template": {
//	      "uuid": "3ce100b7-a734-4b4e-891b-350b1279ade2",
//	      "name": "revive_issue"
//	    },
//	    "variables": ["@contact.name"]
//	  },
//	  "topic": "event"
//	}
//
// @action send_msg
type SendMsgAction struct {
	baseAction
	universalAction
	createMsgAction

	AllURNs    bool           `json:"all_urns,omitempty"`
	Templating *Templating    `json:"templating,omitempty" validate:"omitempty,dive"`
	Topic      flows.MsgTopic `json:"topic,omitempty" validate:"omitempty,msg_topic"`
}

// Templating represents the templating that should be used if possible
type Templating struct {
	UUID      uuids.UUID                       `json:"uuid" validate:"required,uuid4"`
	Template  *assets.TemplateReference        `json:"template" validate:"required"`
	Variables []string                         `json:"variables" engine:"localized,evaluated"`
	Params    map[string][]flows.TemplateParam `json:"params"`
}

// LocalizationUUID gets the UUID which identifies this object for localization
func (t *Templating) LocalizationUUID() uuids.UUID { return t.UUID }

// NewSendMsg creates a new send msg action
func NewSendMsg(uuid flows.ActionUUID, text string, attachments []string, quickReplies []string, allURNs bool) *SendMsgAction {
	return &SendMsgAction{
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
func (a *SendMsgAction) Execute(run flows.Run, step flows.Step, logModifier flows.ModifierCallback, logEvent flows.EventCallback) error {
	// a message to a non-active contact is unsendable but can still be created
	unsendableReason := flows.NilUnsendableReason
	if run.Contact().Status() != flows.ContactStatusActive {
		unsendableReason = flows.UnsendableReasonContactStatus
	}

	evaluatedText, evaluatedAttachments, evaluatedQuickReplies, lang := a.evaluateMessage(run, nil, a.Text, a.Attachments, a.QuickReplies, logEvent)
	locale := currentLocale(run, lang)

	destinations := run.Contact().ResolveDestinations(a.AllURNs)

	sa := run.Session().Assets()

	var template *flows.Template
	if a.Templating != nil {
		template = sa.Templates().Get(a.Templating.Template.UUID)
	}

	// create a new message for each URN+channel destination
	for _, dest := range destinations {
		urn := dest.URN.URN()
		channelRef := assets.NewChannelReference(dest.Channel.UUID(), dest.Channel.Name())

		// do we have a template defined?
		var msg *flows.MsgOut
		if template != nil {
			locales := []i18n.Locale{run.Session().MergedEnvironment().DefaultLocale(), run.Session().Environment().DefaultLocale()}
			templateTranslation := template.FindTranslation(dest.Channel, locales)
			if templateTranslation != nil {
				msg = getTemplatingMsg(a, run, urn, channelRef, templateTranslation, evaluatedAttachments, evaluatedQuickReplies, unsendableReason, logEvent)
			}
		}

		if msg == nil {
			msg = flows.NewMsgOut(urn, channelRef, evaluatedText, evaluatedAttachments, evaluatedQuickReplies, nil, a.Topic, locale, unsendableReason)
		}

		logEvent(events.NewMsgCreated(msg))
	}

	// if we couldn't find a destination, create a msg without a URN or channel and it's up to the caller
	// to handle that as they want
	if len(destinations) == 0 {
		msg := flows.NewMsgOut(urns.NilURN, nil, evaluatedText, evaluatedAttachments, evaluatedQuickReplies, nil, a.Topic, locale, flows.UnsendableReasonNoDestination)
		logEvent(events.NewMsgCreated(msg))
	}

	return nil
}

func getTemplatingMsg(action *SendMsgAction, run flows.Run, urn urns.URN, channelRef *assets.ChannelReference, templateTranslation *flows.TemplateTranslation, evaluatedAttachments []utils.Attachment, evaluatedQuickReplies []string, unsendableReason flows.UnsendableReason, logEvent flows.EventCallback) *flows.MsgOut {
	qrIndex := 0

	localizedVariables, _ := run.GetTextArray(uuids.UUID(action.Templating.UUID), "variables", action.Templating.Variables, nil)
	evaluatedVariables := make([]string, len(localizedVariables))
	for i, variable := range localizedVariables {
		sub, err := run.EvaluateTemplate(variable)
		if err != nil {
			logEvent(events.NewError(err))
		}
		evaluatedVariables[i] = sub
	}
	evaluatedText := templateTranslation.Substitute(evaluatedVariables)

	evaluatedParams := make(map[string][]flows.TemplateParam)
	for compKey, compParams := range action.Templating.Params {
		compVariables := make([]flows.TemplateParam, len(compParams))
		for i, templateParam := range compParams {
			var paramValue string
			var err error
			if strings.HasPrefix(compKey, "button.") {
				paramValue = evaluatedQuickReplies[qrIndex]
				qrIndex++
			} else if templateParam.Type() != "text" {
				paramValue = ""
				for _, att := range evaluatedAttachments {
					attType := strings.Split(att.ContentType(), "/")[0]
					if templateParam.Type() == attType {
						paramValue = att.URL()
						break
					}
				}
			} else {
				localizedParamVariables, _ := run.GetTextArray(uuids.UUID(templateParam.UUID()), "value", []string{templateParam.Value()}, nil)
				paramValue, err = run.EvaluateTemplate(localizedParamVariables[0])
				if err != nil {
					logEvent(events.NewError(err))
				}
			}
			evaluatedParam := flows.NewTemplateParam(templateParam.Type(), templateParam.UUID(), paramValue)
			compVariables[i] = evaluatedParam
		}
		evaluatedParams[compKey] = compVariables
	}
	templating := flows.NewMsgTemplating(action.Templating.Template, evaluatedVariables, templateTranslation.Namespace(), evaluatedParams)
	locale := templateTranslation.Locale()

	return flows.NewMsgOut(urn, channelRef, evaluatedText, evaluatedAttachments, evaluatedQuickReplies, templating, action.Topic, locale, unsendableReason)
}
