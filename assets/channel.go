package assets

import (
	"fmt"

	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/envs"
)

// ChannelUUID is the UUID of a channel
type ChannelUUID uuids.UUID

// ChannelRole is a role that a channel can perform
type ChannelRole string

// different roles that channels can perform
const (
	ChannelRoleSend    ChannelRole = "send"
	ChannelRoleReceive ChannelRole = "receive"
	ChannelRoleCall    ChannelRole = "call"
	ChannelRoleAnswer  ChannelRole = "answer"
	ChannelRoleUSSD    ChannelRole = "ussd"
)

// Channel is something that can send/receive messages.
//
//   {
//     "uuid": "14782905-81a6-4910-bc9f-93ad287b23c3",
//     "name": "My Android",
//     "address": "+593979011111",
//     "schemes": ["tel"],
//     "roles": ["send", "receive"],
//     "country": "EC"
//   }
//
// @asset channel
type Channel interface {
	UUID() ChannelUUID
	Name() string
	Address() string
	Schemes() []string
	Roles() []ChannelRole
	Parent() *ChannelReference
	Country() envs.Country
	MatchPrefixes() []string
	AllowInternational() bool
}

// ChannelReference is used to reference a channel
type ChannelReference struct {
	UUID ChannelUUID `json:"uuid" validate:"required,uuid"`
	Name string      `json:"name"`
}

// NewChannelReference creates a new channel reference with the given UUID and name
func NewChannelReference(uuid ChannelUUID, name string) *ChannelReference {
	return &ChannelReference{UUID: uuid, Name: name}
}

// Type returns the name of the asset type
func (r *ChannelReference) Type() string {
	return "channel"
}

// GenericUUID returns the untyped UUID
func (r *ChannelReference) GenericUUID() uuids.UUID {
	return uuids.UUID(r.UUID)
}

// Identity returns the unique identity of the asset
func (r *ChannelReference) Identity() string {
	return string(r.UUID)
}

// Variable returns whether this a variable (vs concrete) reference
func (r *ChannelReference) Variable() bool {
	return false
}

func (r *ChannelReference) String() string {
	return fmt.Sprintf("%s[uuid=%s,name=%s]", r.Type(), r.Identity(), r.Name)
}

var _ UUIDReference = (*ChannelReference)(nil)
