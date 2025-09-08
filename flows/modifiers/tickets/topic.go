package tickets

import (
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
)

// TypeTopic is the type of our topic modifier
const TypeTopic string = "topic"

// Topic modifies the topic of a ticket
type Topic struct {
	baseModifier

	topic *flows.Topic
}

// NewTopic creates a new topic modifier
func NewTopic(topic *flows.Topic) *Topic {
	return &Topic{
		baseModifier: newBaseModifier(TypeTopic),
		topic:        topic,
	}
}

// Apply applies this modification to the given ticket
func (m *Topic) Apply(eng flows.Engine, env envs.Environment, sa flows.SessionAssets, contact *flows.Contact, ticket *flows.Ticket, log flows.EventCallback) bool {
	if ticket.Topic() != m.topic {
		ticket.SetTopic(m.topic)
		log(events.NewTicketTopicChanged(ticket.UUID(), m.topic.Reference()))
		return true
	}
	return false
}

var _ flows.TicketModifier = (*Topic)(nil)
