package actions

import (
	"encoding/json"
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
	Templating *Templating    `json:"templating,omitempty" validate:"omitempty"`
	Topic      flows.MsgTopic `json:"topic,omitempty" validate:"omitempty,msg_topic"`
}

type TemplateParams struct {
	UUID   uuids.UUID `json:"uuid"`
	Values map[string][]string
}

func (p *TemplateParams) EnumerateTemplates(localization flows.Localization, include func(i18n.Language, string)) {
	for _, comp := range utils.SortedKeys(p.Values) {
		for _, v := range p.Values[comp] {
			include(i18n.NilLanguage, v)
		}
		for _, lang := range localization.Languages() {
			lvals := localization.GetItemTranslation(lang, p.UUID, comp)
			for _, v := range lvals {
				include(lang, v)
			}
		}
	}
}

func (p *TemplateParams) MarshalJSON() ([]byte, error) {
	if p == nil {
		return json.Marshal(p)
	}

	m := make(map[string]any, 1+len(p.Values))
	m["uuid"] = p.UUID
	for k, v := range p.Values {
		m[k] = v
	}
	return json.Marshal(m)
}

func (p *TemplateParams) UnmarshalJSON(d []byte) error {
	var m map[string]any
	if err := json.Unmarshal(d, &m); err != nil {
		return err
	}
	p.Values = make(map[string][]string, len(m)-1)
	for k, v := range m {
		switch typed := v.(type) {
		case string:
			if k == "uuid" {
				p.UUID = uuids.UUID(typed)
			}
		case []any:
			var l []string
			for _, j := range typed {
				l = append(l, j.(string))
			}
			p.Values[k] = l
		}
	}
	return nil
}

// Templating represents the templating that should be used if possible
type Templating struct {
	UUID      uuids.UUID                `json:"uuid" validate:"required,uuid4"`
	Template  *assets.TemplateReference `json:"template" validate:"required"`
	Variables []string                  `json:"variables,omitempty" engine:"localized,evaluated"`
	Params    *TemplateParams           `json:"params,omitempty"`
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

		var msg *flows.MsgOut
		if template != nil {
			locales := []i18n.Locale{run.Session().MergedEnvironment().DefaultLocale(), run.Session().Environment().DefaultLocale()}
			templateTranslation := template.FindTranslation(dest.Channel, locales)
			if templateTranslation != nil {
				msg = a.getTemplateMsg(run, urn, channelRef, templateTranslation, unsendableReason, logEvent)
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

// for message actions that specidy a template, this generates the template message where the message content should be
// considered just a preview of how the template will be evaluated by the channel
func (a *SendMsgAction) getTemplateMsg(run flows.Run, urn urns.URN, channelRef *assets.ChannelReference, translation *flows.TemplateTranslation, unsendableReason flows.UnsendableReason, logEvent flows.EventCallback) *flows.MsgOut {
	evaluatedParams := make(map[string][]string)

	// start by localizing and evaluating either the legacy variables or per-component params
	if len(a.Templating.Variables) > 0 {
		localizedVariables, _ := run.GetTextArray(uuids.UUID(a.Templating.UUID), "variables", a.Templating.Variables, nil)

		evaluatedVariables := make([]string, len(localizedVariables))
		for i, variable := range localizedVariables {
			sub, _ := run.EvaluateTemplate(variable, logEvent)
			evaluatedVariables[i] = sub
		}

		evaluatedParams["body"] = evaluatedVariables

	} else if a.Templating.Params != nil {
		for comp, compParams := range a.Templating.Params.Values {
			localizedCompParams, _ := run.GetTextArray(uuids.UUID(a.Templating.Params.UUID), comp, compParams, nil)
			evaluatedCompParams := make([]string, len(localizedCompParams))

			for i, variable := range localizedCompParams {
				sub, _ := run.EvaluateTemplate(variable, logEvent)
				evaluatedCompParams[i] = sub
			}
			evaluatedParams[comp] = evaluatedCompParams
		}
	}

	// next we cross reference with params defined in the template translation to get types
	params := make(map[string][]flows.TemplateParam, len(translation.Components()))

	for key, comp := range translation.Components() {
		compParams := comp.Params()
		if len(compParams) > 0 {
			params[key] = make([]flows.TemplateParam, len(compParams))
		}

		for i, tp := range compParams {
			if i < len(evaluatedParams[key]) {
				params[key][i] = flows.TemplateParam{Type: tp.Type(), Value: evaluatedParams[key][i]}
			} else {
				params[key][i] = flows.TemplateParam{Type: tp.Type(), Value: ""}
			}
		}
	}

	locale := translation.Locale()
	templating := flows.NewMsgTemplating(a.Templating.Template, params, translation.Namespace())

	// extract content for preview message
	preview := translation.Preview(templating)
	previewText := preview["body"]
	var previewQRs []string
	for _, key := range utils.SortedKeys(preview) {
		if strings.HasPrefix(key, "button.") {
			previewQRs = append(previewQRs, preview[key])
		}
	}

	return flows.NewMsgOut(urn, channelRef, previewText, nil, previewQRs, templating, flows.NilMsgTopic, locale, unsendableReason)
}
