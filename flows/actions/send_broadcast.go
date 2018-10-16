package actions

import (
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	RegisterType(TypeSendBroadcast, func() flows.Action { return &SendBroadcastAction{} })
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
	BaseAction
	onlineAction

	Text         string                    `json:"text"`
	Attachments  []string                  `json:"attachments"`
	QuickReplies []string                  `json:"quick_replies,omitempty"`
	URNs         []urns.URN                `json:"urns,omitempty"`
	Contacts     []*flows.ContactReference `json:"contacts,omitempty" validate:"dive"`
	Groups       []*assets.GroupReference  `json:"groups,omitempty" validate:"dive"`
	LegacyVars   []string                  `json:"legacy_vars,omitempty"`
}

// NewSendBroadcastAction creates a new send broadcast action
func NewSendBroadcastAction(uuid flows.ActionUUID, text string, attachments []string, quickReplies []string, urns []urns.URN, contacts []*flows.ContactReference, groups []*assets.GroupReference, legacyVars []string) *SendBroadcastAction {
	return &SendBroadcastAction{
		BaseAction:   NewBaseAction(TypeSendBroadcast, uuid),
		Text:         text,
		Attachments:  attachments,
		QuickReplies: quickReplies,
		URNs:         urns,
		Contacts:     contacts,
		Groups:       groups,
		LegacyVars:   legacyVars,
	}
}

// Validate validates our action is valid and has all the assets it needs
func (a *SendBroadcastAction) Validate(assets flows.SessionAssets, context *flows.ValidationContext) error {
	return a.validateGroups(assets, a.Groups)
}

// Execute runs this action
func (a *SendBroadcastAction) Execute(run flows.FlowRun, step flows.Step) error {
	urnList, contactRefs, groupRefs, err := a.resolveContactsAndGroups(run, step, a.URNs, a.Contacts, a.Groups, a.LegacyVars)
	if err != nil {
		return err
	}

	translations := make(map[utils.Language]*events.BroadcastTranslation)
	languages := append([]utils.Language{run.Flow().Language()}, run.Flow().Localization().Languages()...)

	// evaluate the broadcast in each language we have translations for
	for _, language := range languages {
		languages := []utils.Language{language, run.Flow().Language()}

		evaluatedText, evaluatedAttachments, evaluatedQuickReplies := a.evaluateMessage(run, step, languages, a.Text, a.Attachments, a.QuickReplies)
		translations[language] = &events.BroadcastTranslation{
			Text:         evaluatedText,
			Attachments:  evaluatedAttachments,
			QuickReplies: evaluatedQuickReplies,
		}
	}

	a.log(run, step, events.NewBroadcastCreatedEvent(translations, run.Flow().Language(), urnList, contactRefs, groupRefs))

	return nil
}
