package flows

import (
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
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
func (t *Topic) Asset() assets.Topic { return t.Topic }

// Reference returns a reference to this topic
func (t *Topic) Reference() *assets.TopicReference {
	if t == nil {
		return nil
	}
	return assets.NewTopicReference(t.UUID(), t.Name())
}

// Context returns the properties available in expressions
//
//   __default__:text -> the name
//   uuid:text -> the UUID of the topic
//   name:text -> the name of the topic
//
// @context topic
func (t *Topic) Context(env envs.Environment) map[string]types.XValue {

	return map[string]types.XValue{
		"__default__": types.NewXText(t.Name()),
		"uuid":        types.NewXText(string(t.UUID())),
		"name":        types.NewXText(t.Name()),
	}
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
