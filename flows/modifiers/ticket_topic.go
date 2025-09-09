package modifiers

import (
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

// TicketTopic modifies the topic of a ticket
type TicketTopic struct {
	baseModifier

	topic *flows.Topic
}

// NewTicketTopic creates a new topic modifier
func NewTicketTopic(topic *flows.Topic) *TicketTopic {
	return &TicketTopic{
		baseModifier: newBaseModifier(TypeTicketTopic),
		topic:        topic,
	}
}

// Apply applies this modification to the given ticket
func (m *TicketTopic) Apply(eng flows.Engine, env envs.Environment, sa flows.SessionAssets, contact *flows.Contact, ticket *flows.Ticket, log flows.EventCallback) bool {
	if ticket != nil && ticket.Topic() != m.topic {
		ticket.SetTopic(m.topic)
		log(events.NewTicketTopicChanged(ticket.UUID(), m.topic.Reference()))
		return true
	}
	return false
}

var _ flows.Modifier = (*TicketTopic)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type ticketTopicEnvelope struct {
	utils.TypedEnvelope

	Topic *assets.TopicReference `json:"topic"`
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

	return NewTicketTopic(topic), nil
}

func (m *TicketTopic) MarshalJSON() ([]byte, error) {
	return jsonx.Marshal(&ticketTopicEnvelope{
		TypedEnvelope: utils.TypedEnvelope{Type: m.Type()},
		Topic:         m.topic.Reference(),
	})
}
