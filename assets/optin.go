package assets

import (
	"fmt"

	"github.com/nyaruka/gocommon/uuids"
)

// OptInUUID is the UUID of an opt in
type OptInUUID uuids.UUID

// OptIn are channel specific opt-ins
//
//	{
//	  "uuid": "8925c76f-926b-4a63-a6eb-ab69e7a6b79b",
//	  "name": "Joke Of The Day",
//	  "channel": {
//	     "uuid": "204e5af9-42c3-4d46-8aab-ce204dff25b4",
//	     "name": "Facebook"
//	  }
//	}
//
// @asset optin
type OptIn interface {
	UUID() OptInUUID
	Name() string
	Channel() *ChannelReference
}

// OptInReference is used to reference an opt in
type OptInReference struct {
	UUID OptInUUID `json:"uuid" validate:"required,uuid"`
	Name string    `json:"name"`
}

// NewOptInReference creates a new optin reference with the given UUID and name
func NewOptInReference(uuid OptInUUID, name string) *OptInReference {
	return &OptInReference{UUID: uuid, Name: name}
}

// Type returns the name of the asset type
func (r *OptInReference) Type() string {
	return "optin"
}

// GenericUUID returns the untyped UUID
func (r *OptInReference) GenericUUID() uuids.UUID {
	return uuids.UUID(r.UUID)
}

// Identity returns the unique identity of the asset
func (r *OptInReference) Identity() string {
	return string(r.UUID)
}

// Variable returns whether this a variable (vs concrete) reference
func (r *OptInReference) Variable() bool {
	return false
}

func (r *OptInReference) String() string {
	return fmt.Sprintf("%s[uuid=%s,name=%s]", r.Type(), r.Identity(), r.Name)
}

var _ UUIDReference = (*OptInReference)(nil)
