package assets

import (
	"fmt"

	"github.com/nyaruka/gocommon/uuids"
)

// TopicUUID is the UUID of a topic
type TopicUUID uuids.UUID

// Topic categorizes tickets
//
//   {
//     "uuid": "cd48bd11-08b9-44e3-9778-8e26adf08a7a",
//     "name": "Weather"
//   }
//
// @asset topic
type Topic interface {
	UUID() TopicUUID
	Name() string
}

// TopicReference is used to reference a topic
type TopicReference struct {
	UUID TopicUUID `json:"uuid" validate:"required,uuid"`
	Name string    `json:"name"`
}

// NewTopicReference creates a new topic reference with the given UUID and name
func NewTopicReference(uuid TopicUUID, name string) *TopicReference {
	return &TopicReference{UUID: uuid, Name: name}
}

// Type returns the name of the asset type
func (r *TopicReference) Type() string {
	return "topic"
}

// GenericUUID returns the untyped UUID
func (r *TopicReference) GenericUUID() uuids.UUID {
	return uuids.UUID(r.UUID)
}

// Identity returns the unique identity of the asset
func (r *TopicReference) Identity() string {
	return string(r.UUID)
}

// Variable returns whether this a variable (vs concrete) reference
func (r *TopicReference) Variable() bool {
	return false
}

func (r *TopicReference) String() string {
	return fmt.Sprintf("%s[uuid=%s,name=%s]", r.Type(), r.Identity(), r.Name)
}

var _ UUIDReference = (*TopicReference)(nil)
