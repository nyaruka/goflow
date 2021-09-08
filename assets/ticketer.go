package assets

import (
	"fmt"

	"github.com/nyaruka/gocommon/uuids"
)

// TicketerUUID is the UUID of a ticketer
type TicketerUUID uuids.UUID

// Ticketer is a system which can open or close tickets
//
//   {
//     "uuid": "37657cf7-5eab-4286-9cb0-bbf270587bad",
//     "name": "Support Tickets",
//     "type": "mailgun"
//   }
//
// @asset ticketer
type Ticketer interface {
	UUID() TicketerUUID
	Name() string
	Type() string
}

// TicketerReference is used to reference a ticketer
type TicketerReference struct {
	UUID TicketerUUID `json:"uuid" validate:"required,uuid"`
	Name string       `json:"name"`
}

// NewTicketerReference creates a new classifier reference with the given UUID and name
func NewTicketerReference(uuid TicketerUUID, name string) *TicketerReference {
	return &TicketerReference{UUID: uuid, Name: name}
}

// Type returns the name of the asset type
func (r *TicketerReference) Type() string {
	return "ticketer"
}

// GenericUUID returns the untyped UUID
func (r *TicketerReference) GenericUUID() uuids.UUID {
	return uuids.UUID(r.UUID)
}

// Identity returns the unique identity of the asset
func (r *TicketerReference) Identity() string {
	return string(r.UUID)
}

// Variable returns whether this a variable (vs concrete) reference
func (r *TicketerReference) Variable() bool {
	return false
}

func (r *TicketerReference) String() string {
	return fmt.Sprintf("%s[uuid=%s,name=%s]", r.Type(), r.Identity(), r.Name)
}

var _ UUIDReference = (*TicketerReference)(nil)
