package flows

import (
	"encoding/json"
	"fmt"

	"github.com/nyaruka/goflow/utils"
)

type channel struct {
	uuid        ChannelUUID
	name        string
	address     string
	channelType ChannelType
}

// UUID returns the UUID of this channel
func (c *channel) UUID() ChannelUUID { return c.uuid }

// Name returns the name of this channel
func (c *channel) Name() string { return c.name }

// Name returns the address of this channel
func (c *channel) Address() string { return c.address }

// Type returns the type of this channel
func (c *channel) Type() ChannelType { return c.channelType }

// Resolve satisfies our resolver interface
func (c *channel) Resolve(key string) interface{} {
	switch key {

	case "uuid":
		return c.uuid

	case "name":
		return c.name

	case "address":
		return c.address

	case "type":
		return c.channelType
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
	UUID        ChannelUUID `json:"uuid"`
	Name        string      `json:"name"`
	Address     string      `json:"address"`
	ChannelType ChannelType `json:"type"`
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
	c := &channel{}
	err := json.Unmarshal(data, c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// UnmarshalJSON is our custom unmarshalling of a channel
func (c *channel) UnmarshalJSON(data []byte) error {
	var ce channelEnvelope
	var err error

	err = json.Unmarshal(data, &ce)
	if err != nil {
		return err
	}

	c.uuid = ce.UUID
	c.name = ce.Name
	c.address = ce.Address
	c.channelType = ce.ChannelType

	return nil
}

// MarshalJSON is our custom marshalling of a channel
func (c *channel) MarshalJSON() ([]byte, error) {
	var ce channelEnvelope

	ce.UUID = c.uuid
	ce.Name = c.name
	ce.Address = c.address
	ce.ChannelType = c.channelType

	return json.Marshal(ce)
}
