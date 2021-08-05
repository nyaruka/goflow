package flows

import (
	"github.com/nyaruka/goflow/assets"
)

// Topic represents a ticket topic
type Topic struct {
	assets.Topic
}

// NewTopic creates a new topic from the given asset
func NewTopic(asset assets.Topic) *Topic {
	return &Topic{Topic: asset}
}

// Asset returns the underlying asset
func (l *Topic) Asset() assets.Topic { return l.Topic }

// Reference returns a reference to this topic
func (l *Topic) Reference() *assets.TopicReference {
	if l == nil {
		return nil
	}
	return assets.NewTopicReference(l.UUID(), l.Name())
}

var _ assets.Topic = (*Topic)(nil)

// TopicAssets provides access to all topic assets
type TopicAssets struct {
	byUUID map[assets.TopicUUID]*Topic
}

// NewTopicAssets creates a new set of topic assets
func NewTopicAssets(topics []assets.Topic) *TopicAssets {
	s := &TopicAssets{
		byUUID: make(map[assets.TopicUUID]*Topic, len(topics)),
	}
	for _, asset := range topics {
		topic := NewTopic(asset)
		s.byUUID[topic.UUID()] = topic
	}
	return s
}

// Get returns the topic with the given UUID
func (s *TopicAssets) Get(uuid assets.TopicUUID) *Topic {
	return s.byUUID[uuid]
}
