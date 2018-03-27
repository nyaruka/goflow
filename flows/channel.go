package flows

import (
	"encoding/json"
	"fmt"

	"github.com/nyaruka/goflow/utils"
)

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

type channel struct {
	uuid    ChannelUUID
	name    string
	address string
	schemes []string
	roles   []ChannelRole
}

// NewChannel creates a new channel
func NewChannel(uuid ChannelUUID, name string, address string, schemes []string, roles []ChannelRole) Channel {
	return &channel{
		uuid:    uuid,
		name:    name,
		address: address,
		schemes: schemes,
		roles:   roles,
	}
}

// UUID returns the UUID of this channel
func (c *channel) UUID() ChannelUUID { return c.uuid }

// Name returns the name of this channel
func (c *channel) Name() string { return c.name }

// Address returns the address of this channel
func (c *channel) Address() string { return c.address }

// Schemes returns the supported schemes of this channel
func (c *channel) Schemes() []string { return c.schemes }

// Roles returns the roles of this channel
func (c *channel) Roles() []ChannelRole { return c.roles }

// Reference returns a reference to this channel
func (c *channel) Reference() *ChannelReference { return NewChannelReference(c.uuid, c.name) }

// SupportsScheme returns whether this channel supports the given URN scheme
func (c *channel) SupportsScheme(scheme string) bool {
	for _, s := range c.schemes {
		if s == scheme {
			return true
		}
	}
	return false
}

// HasRole returns whether this channel has the given role
func (c *channel) HasRole(role ChannelRole) bool {
	for _, r := range c.roles {
		if r == role {
			return true
		}
	}
	return false
}

// Resolve resolves the given key when this channel is referenced in an expression
func (c *channel) Resolve(key string) interface{} {
	switch key {
	case "uuid":
		return c.uuid
	case "name":
		return c.name
	case "address":
		return c.address
	}

	return fmt.Errorf("No field '%s' on channel", key)
}

// String returns the default string value for a channel, which is its name
func (c *channel) String() string {
	return c.name
}

var _ utils.VariableResolver = (*channel)(nil)

// ChannelSet defines the unordered set of all channels for a session
type ChannelSet struct {
	channels       []Channel
	channelsByUUID map[ChannelUUID]Channel
}

// NewChannelSet creates a new channel set
func NewChannelSet(channels []Channel) *ChannelSet {
	s := &ChannelSet{channels: channels, channelsByUUID: make(map[ChannelUUID]Channel, len(channels))}
	for _, channel := range s.channels {
		s.channelsByUUID[channel.UUID()] = channel
	}
	return s
}

// GetForURN returns the best channel for the given URN
func (s *ChannelSet) GetForURN(urn *ContactURN) Channel {
	// if caller has told us which channel to use for this URN, use that
	if urn.Channel() != nil {
		return urn.Channel()
	}

	// if not, return the first channel which supports this URN scheme
	scheme := urn.Scheme()
	for _, ch := range s.channels {
		if ch.HasRole(ChannelRoleSend) && ch.SupportsScheme(scheme) {
			return ch
		}
	}

	return nil
}

// FindByUUID finds the channel with the given UUID
func (s *ChannelSet) FindByUUID(uuid ChannelUUID) Channel {
	return s.channelsByUUID[uuid]
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type channelEnvelope struct {
	UUID    ChannelUUID   `json:"uuid" validate:"required,uuid4"`
	Name    string        `json:"name"`
	Address string        `json:"address"`
	Schemes []string      `json:"schemes" validate:"min=1"`
	Roles   []ChannelRole `json:"roles" validate:"min=1,dive,eq=send|eq=receive|eq=call|eq=answer|eq=ussd"`
}

// ReadChannel decodes a channel from the passed in JSON
func ReadChannel(data json.RawMessage) (Channel, error) {
	ce := channelEnvelope{}
	if err := utils.UnmarshalAndValidate(data, &ce, "channel"); err != nil {
		return nil, err
	}

	return &channel{
		uuid:    ce.UUID,
		name:    ce.Name,
		address: ce.Address,
		schemes: ce.Schemes,
		roles:   ce.Roles,
	}, nil
}

// ReadChannelSet decodes channels from the passed in JSON
func ReadChannelSet(data json.RawMessage) (*ChannelSet, error) {
	items, err := utils.UnmarshalArray(data)
	if err != nil {
		return nil, err
	}

	channels := make([]Channel, len(items))
	for c := range items {
		channels[c], err = ReadChannel(items[c])
		if err != nil {
			return nil, err
		}
	}
	return NewChannelSet(channels), nil
}
