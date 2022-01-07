package actions

import (
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
)

func init() {
	registerType(TypeSendBroadcast, func() flows.Action { return &SendBroadcastAction{} })
}

// TypeSendBroadcast is the type for the send broadcast action
const TypeSendBroadcast string = "send_broadcast"

// SendBroadcastAction can be used to send a message to one or more contacts. It accepts a list of URNs, a list of groups
// and a list of contacts.
//
// The URNs and text fields may be templates. A [event:broadcast_created] event will be created for each unique urn, contact and group
// with the evaluated text.
//
//   {
//     "uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
//     "type": "send_broadcast",
//     "urns": ["tel:+12065551212"],
//     "text": "Hi @contact.name, are you ready to complete today's survey?"
//   }
//
// @action send_broadcast
type SendBroadcastAction struct {
	baseAction
	onlineAction
	otherContactsAction
	createMsgAction
}

// NewSendBroadcast creates a new send broadcast action
func NewSendBroadcast(uuid flows.ActionUUID, text string, attachments []string, quickReplies []string, urns []urns.URN, contacts []*flows.ContactReference, groups []*assets.GroupReference, legacyVars []string) *SendBroadcastAction {
	return &SendBroadcastAction{
		baseAction: newBaseAction(TypeSendBroadcast, uuid),
		otherContactsAction: otherContactsAction{
			URNs:       urns,
			Contacts:   contacts,
			Groups:     groups,
			LegacyVars: legacyVars,
		},
		createMsgAction: createMsgAction{
			Text:         text,
			Attachments:  attachments,
			QuickReplies: quickReplies,
		},
	}
}

// Execute runs this action
func (a *SendBroadcastAction) Execute(run flows.Run, step flows.Step, logModifier flows.ModifierCallback, logEvent flows.EventCallback) error {
	groupRefs, contactRefs, _, urnList, err := a.resolveRecipients(run, logEvent)
	if err != nil {
		return err
	}

	// footgun prevention
	if run.Session().BatchStart() && len(groupRefs) > 0 {
		logEvent(events.NewErrorf("can't send broadcasts to groups during batch starts"))
		return nil
	}

	translations := make(map[envs.Language]*events.BroadcastTranslation)
	languages := append([]envs.Language{run.Flow().Language()}, run.Flow().Localization().Languages()...)

	// evaluate the broadcast in each language we have translations for
	for _, language := range languages {
		languages := []envs.Language{language, run.Flow().Language()}

		evaluatedText, evaluatedAttachments, evaluatedQuickReplies := a.evaluateMessage(run, languages, a.Text, a.Attachments, a.QuickReplies, logEvent)
		translations[language] = &events.BroadcastTranslation{
			Text:         evaluatedText,
			Attachments:  evaluatedAttachments,
			QuickReplies: evaluatedQuickReplies,
		}
	}

	// if we have any recipients, log an event
	if len(urnList) > 0 || len(contactRefs) > 0 || len(groupRefs) > 0 {
		logEvent(events.NewBroadcastCreated(translations, run.Flow().Language(), groupRefs, contactRefs, urnList))
	}

	return nil
}
