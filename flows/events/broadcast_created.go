package events

import (
	"github.com/nyaruka/gocommon/i18n"
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
)

func init() {
	registerType(TypeBroadcastCreated, func() flows.Event { return &BroadcastCreated{} })
}

// TypeBroadcastCreated is a constant for outgoing message events
const TypeBroadcastCreated string = "broadcast_created"

// BroadcastCreated events are created when an action wants to send a message to other contacts.
//
//	{
//	  "uuid": "0197b335-6ded-79a4-95a6-3af85b57f108",
//	  "type": "broadcast_created",
//	  "created_on": "2006-01-02T15:04:05Z",
//	  "translations": {
//	    "eng": {
//	      "text": "hi, what's up",
//	      "attachments": [],
//	      "quick_replies": ["All good", "Got 99 problems"]
//	    },
//	    "spa": {
//	      "text": "Que pasa",
//	      "attachments": [],
//	      "quick_replies": ["Todo bien", "Tengo 99 problemas"]
//	    }
//	  },
//	  "base_language": "eng",
//	  "urns": ["tel:+12065551212"],
//	  "contacts": [{"uuid": "0e06f977-cbb7-475f-9d0b-a0c4aaec7f6a", "name": "Bob"}]
//	}
//
// @event broadcast_created
type BroadcastCreated struct {
	BaseEvent

	Translations flows.BroadcastTranslations `json:"translations" validate:"min=1,dive"`
	BaseLanguage i18n.Language               `json:"base_language" validate:"required"`
	Groups       []*assets.GroupReference    `json:"groups,omitempty" validate:"dive"`
	Contacts     []*flows.ContactReference   `json:"contacts,omitempty" validate:"dive"`
	ContactQuery string                      `json:"contact_query,omitempty"`
	URNs         []urns.URN                  `json:"urns,omitempty" validate:"dive,urn"`
}

// NewBroadcastCreated creates a new outgoing msg event for the given recipients
func NewBroadcastCreated(translations flows.BroadcastTranslations, baseLanguage i18n.Language, groups []*assets.GroupReference, contacts []*flows.ContactReference, contactQuery string, urns []urns.URN) *BroadcastCreated {
	return &BroadcastCreated{
		BaseEvent:    NewBaseEvent(TypeBroadcastCreated),
		Translations: translations,
		BaseLanguage: baseLanguage,
		Groups:       groups,
		Contacts:     contacts,
		ContactQuery: contactQuery,
		URNs:         urns,
	}
}

var _ flows.Event = (*BroadcastCreated)(nil)
