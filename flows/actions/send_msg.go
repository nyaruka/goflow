package actions

import (
	"fmt"
	"strings"

	"github.com/nyaruka/gocommon/i18n"
	"github.com/nyaruka/gocommon/stringsx"
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
//	    "template": {
//	      "uuid": "3ce100b7-a734-4b4e-891b-350b1279ade2",
//	      "name": "revive_issue"
//	    },
//	    "components": [
//	      {
//	      	"uuid": "2131d674-9afb-41e8-bcb8-5910c2faec2f",
//	      	"name": "body",
//	      	"params": ["@contact.name"]
//	      }
//	    ]
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

type TemplatingComponent struct {
	UUID   uuids.UUID `json:"uuid" validate:"required,uuid4"`
	Name   string     `json:"name" validate:"required"`
	Params []string   `json:"params" engine:"localized,evaluated"`
}

// LocalizationUUID gets the UUID which identifies this object for localization
func (c *TemplatingComponent) LocalizationUUID() uuids.UUID { return c.UUID }

// Templating represents the templating that should be used if possible
type Templating struct {
	Template   *assets.TemplateReference `json:"template" validate:"required"`
	Components []*TemplatingComponent    `json:"components,omitempty"`
}

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
	// start by localizing and evaluating the param values - for now these are per-component
	evaluatedParams := make(map[string][]string)
	totalParams := 0
	for _, comp := range a.Templating.Components {
		localizedCompParams, _ := run.GetTextArray(comp.UUID, "params", comp.Params, nil)
		evaluatedCompParams := make([]string, len(localizedCompParams))
		for i, variable := range localizedCompParams {
			sub, _ := run.EvaluateTemplate(variable, logEvent)
			evaluatedCompParams[i] = sub
			totalParams++
		}
		evaluatedParams[comp.Name] = evaluatedCompParams
	}

	// Turn that into a single list of variable values using the order of components on translation.
	// Eventually this is what will be stored in the flow definition.
	evaluatedVariables := make([]string, 0, totalParams)
	for _, comp := range translation.Components() {
		evaluatedVariables = append(evaluatedVariables, evaluatedParams[comp.Name()]...)
	}

	// cross-reference with asset to get variable types
	variables := make([]*flows.TemplatingVariable, len(translation.Variables()))
	for i, v := range translation.Variables() {
		if i < len(evaluatedVariables) {
			variables[i] = &flows.TemplatingVariable{Type: v.Type(), Value: evaluatedVariables[i]}
		} else {
			variables[i] = &flows.TemplatingVariable{Type: v.Type(), Value: ""}
		}
	}

	// create a list of components that have variables
	components := make([]*flows.TemplatingComponent, 0, len(translation.Components()))
	for _, comp := range translation.Components() {
		if len(comp.Variables()) > 0 {
			tc := &flows.TemplatingComponent{Type: comp.Type(), Name: comp.Name(), Variables: comp.Variables()}

			// TODO remove when courier updated
			tc.Params = make([]*flows.TemplatingVariable, 0, len(comp.Variables()))
			for _, vn := range utils.SortedKeys(comp.Variables()) {
				vi := comp.Variables()[vn]
				tc.Params = append(tc.Params, variables[vi])
			}

			components = append(components, tc)
		}
	}

	// the message we return is an approximate preview of what the channel will send using the template
	var previewParts []string
	var previewQRs []string

	for _, comp := range translation.Components() {
		previewContent := comp.Content()
		for key, index := range comp.Variables() {
			previewContent = strings.ReplaceAll(previewContent, fmt.Sprintf("{{%s}}", key), variables[index].Value)
		}

		if previewContent != "" {
			if comp.Type() == "header" || comp.Type() == "body" || comp.Type() == "footer" {
				previewParts = append(previewParts, previewContent)
			} else if strings.HasPrefix(comp.Type(), "button/") {
				previewQRs = append(previewQRs, stringsx.TruncateEllipsis(previewContent, maxQuickReplyLength))
			}
		}
	}

	previewText := strings.Join(previewParts, "\n\n")
	locale := translation.Locale()
	templating := flows.NewMsgTemplating(a.Templating.Template, translation.Namespace(), components, variables)

	return flows.NewMsgOut(urn, channelRef, previewText, nil, previewQRs, templating, flows.NilMsgTopic, locale, unsendableReason)
}
