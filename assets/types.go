package assets

import (
	"github.com/nyaruka/goflow/utils"
)

// ChannelUUID is the UUID of a channel
type ChannelUUID utils.UUID

const NilChannelUUID ChannelUUID = ChannelUUID("")

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

// Channel is something that can send/receive messages
type Channel interface {
	UUID() ChannelUUID
	Name() string
	Address() string
	Schemes() []string
	Roles() []ChannelRole
	ParentUUID() ChannelUUID
	Country() string
	MatchPrefixes() []string
}

// GroupUUID is the UUID of a group
type GroupUUID utils.UUID

// Group is a set of contacts
type Group interface {
	UUID() GroupUUID
	Name() string
	Query() string
}

// LabelUUID is the UUID of a label
type LabelUUID utils.UUID

// Label is something that can be applied a message
type Label interface {
	UUID() LabelUUID
	Name() string
}

// AssetSource is a source of assets
type AssetSource interface {
	Channels() ([]Channel, error)
	Groups() ([]Group, error)
	Labels() ([]Label, error)
	HasLocations() bool
}
