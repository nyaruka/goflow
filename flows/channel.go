package flows

import (
	"encoding/json"
	"fmt"

	"github.com/nyaruka/goflow/utils"
)

// Channel represents a channel in the system. Channels have a UUID, name, type and configuration
type Channel struct {
	uuid        ChannelUUID
	name        string
	channelType ChannelType
	config      string
}

// UUID returns the UUID of this channel
func (c *Channel) UUID() ChannelUUID { return c.uuid }

// Name returns the name of this channel
func (c *Channel) Name() string { return c.name }

// Type returns the type of this channel
func (c *Channel) Type() ChannelType { return c.channelType }

// Resolve satisfies our resolver interface
func (c *Channel) Resolve(key string) interface{} {
	switch key {

	case "name":
		return c.name

	case "uuid":
		return c.uuid

	case "type":
		return c.channelType
	}

	return fmt.Errorf("No field '%s' on channel", key)
}

// Default returns the default value for a channel, which is itself
func (c *Channel) Default() interface{} {
	return c
}

// String returns the default string value for a channel, which is its name
func (c *Channel) String() string {
	return c.name
}

var _ utils.VariableResolver = (*Channel)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

// ReadChannel decodes a channel from the passed in JSON
func ReadChannel(data json.RawMessage) (*Channel, error) {
	channel := &Channel{}
	err := json.Unmarshal(data, channel)
	if err == nil {
		// err = run.Validate()
	}
	return channel, err
}

type channelEnvelope struct {
	UUID        ChannelUUID `json:"uuid"`
	Name        string      `json:"name"`
	ChannelType ChannelType `json:"type"`
}

// UnmarshalJSON is our custom unmarshalling of a channel
func (c *Channel) UnmarshalJSON(data []byte) error {
	var ce channelEnvelope
	var err error

	err = json.Unmarshal(data, &ce)
	if err != nil {
		return err
	}

	c.name = ce.Name
	c.uuid = ce.UUID
	c.channelType = ce.ChannelType

	return nil
}

// MarshalJSON is our custom marshalling of a channel
func (c *Channel) MarshalJSON() ([]byte, error) {
	var ce channelEnvelope

	ce.Name = c.name
	ce.UUID = c.uuid
	ce.ChannelType = c.channelType

	return json.Marshal(ce)
}
