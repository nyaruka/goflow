package modifiers

import (
	"slices"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	registerType(TypeTicketTopic, readTicketTopic)
}

// TypeTicketTopic is the type of our topic modifier
const TypeTicketTopic string = "ticket_topic"

// TicketTopic modifies the topic of tickets
type TicketTopic struct {
	baseModifier

	ticketUUIDs []flows.TicketUUID
	topic       *flows.Topic
}

// NewTicketTopic creates a new topic modifier
func NewTicketTopic(ticketUUIDs []flows.TicketUUID, topic *flows.Topic) *TicketTopic {
	return &TicketTopic{
		baseModifier: newBaseModifier(TypeTicketTopic),
		ticketUUIDs:  ticketUUIDs,
		topic:        topic,
	}
}

// Apply applies this modification to the given contact
func (m *TicketTopic) Apply(eng flows.Engine, env envs.Environment, sa flows.SessionAssets, contact *flows.Contact, log flows.EventCallback) bool {
	modified := false

	for _, ticket := range contact.Tickets().All() {
		if slices.Contains(m.ticketUUIDs, ticket.UUID()) && ticket.Topic() != m.topic {
			ticket.SetTopic(m.topic)
			log(events.NewTicketTopicChanged(ticket.UUID(), m.topic.Reference()))
			modified = true
		}
	}
	return modified
}

var _ flows.Modifier = (*TicketTopic)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type ticketTopicEnvelope struct {
	utils.TypedEnvelope

	TicketUUIDs []flows.TicketUUID     `json:"ticket_uuids" validate:"required,dive,uuid"`
	Topic       *assets.TopicReference `json:"topic"`
}

func readTicketTopic(sa flows.SessionAssets, data []byte, missing assets.MissingCallback) (flows.Modifier, error) {
	e := &ticketTopicEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	topic := sa.Topics().Get(e.Topic.UUID)
	if topic == nil {
		missing(e.Topic, nil)
		return nil, ErrNoModifier // can't proceed without a topic
	}

	return NewTicketTopic(e.TicketUUIDs, topic), nil
}

func (m *TicketTopic) MarshalJSON() ([]byte, error) {
	return jsonx.Marshal(&ticketTopicEnvelope{
		TypedEnvelope: utils.TypedEnvelope{Type: m.Type()},
		TicketUUIDs:   m.ticketUUIDs,
		Topic:         m.topic.Reference(),
	})
}
