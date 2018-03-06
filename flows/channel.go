package flows

import (
	"encoding/json"
	"fmt"

	"github.com/nyaruka/goflow/utils"
)

type ChannelRole string

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
	parent  *ChannelReference
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

// Parent returns the parent of this channel
func (c *channel) Parent() *ChannelReference { return c.parent }

// Reference returns a reference to this channel
func (c *channel) Reference() *ChannelReference { return NewChannelReference(c.uuid, c.name) }

// Resolve satisfies our resolver interface
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

// Default returns the default value for a channel, which is itself
func (c *channel) Default() interface{} {
	return c
}

// String returns the default string value for a channel, which is its name
func (c *channel) String() string {
	return c.name
}

var _ utils.VariableResolver = (*channel)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type channelEnvelope struct {
	UUID    ChannelUUID       `json:"uuid" validate:"required,uuid4"`
	Name    string            `json:"name"`
	Address string            `json:"address"`
	Schemes []string          `json:"schemes" validate:"min=1"`
	Roles   []ChannelRole     `json:"roles" validate:"min=1,dive,eq=send|eq=receive|eq=call|eq=answer|eq=ussd"`
	Parent  *ChannelReference `json:"parent" validate:"omitempty,dive"`
}

// ReadChannels decodes channels from the passed in JSON
func ReadChannels(data []json.RawMessage) ([]Channel, error) {
	channels := make([]Channel, len(data))
	var err error
	for c := range data {
		channels[c], err = ReadChannel(data[c])
		if err != nil {
			return nil, err
		}
	}
	return channels, nil
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
		parent:  ce.Parent,
	}, nil
}

// MarshalJSON is our custom marshalling of a channel
func (c *channel) MarshalJSON() ([]byte, error) {
	ce := channelEnvelope{
		UUID:    c.uuid,
		Name:    c.name,
		Address: c.address,
		Schemes: c.schemes,
		Roles:   c.roles,
		Parent:  c.parent,
	}

	return json.Marshal(ce)
}
