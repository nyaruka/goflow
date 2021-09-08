package static

import (
	"github.com/nyaruka/goflow/assets"
)

// Topic is a JSON serializable implementation of a topic asset
type Topic struct {
	UUID_ assets.TopicUUID `json:"uuid" validate:"required,uuid"`
	Name_ string           `json:"name"`
}

// NewTopic creates a new topic
func NewTopic(uuid assets.TopicUUID, name string) assets.Topic {
	return &Topic{
		UUID_: uuid,
		Name_: name,
	}
}

// UUID returns the UUID of this ticketer
func (t *Topic) UUID() assets.TopicUUID { return t.UUID_ }

// Name returns the name of this ticketer
func (t *Topic) Name() string { return t.Name_ }
