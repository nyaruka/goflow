package flow

import (
	"encoding/json"
	"fmt"

	"github.com/nyaruka/goflow/flows"
)

type channel struct {
	uuid        flows.ChannelUUID
	name        string
	channelType flows.ChannelType
	config      string
}

func (c *channel) UUID() flows.ChannelUUID { return c.uuid }
func (c *channel) Name() string            { return c.name }
func (c *channel) Type() flows.ChannelType { return c.channelType }

func (c *channel) Resolve(key string) interface{} {
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

func (c *channel) Default() interface{} {
	return c
}

func (c *channel) String() interface{} {
	return c.name
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

// ReadChannel decodes a channel from the passed in JSON
func ReadChannel(data json.RawMessage) (flows.Channel, error) {
	channel := &channel{}
	err := json.Unmarshal(data, channel)
	if err == nil {
		// err = run.Validate()
	}
	return channel, err
}

type channelEnvelope struct {
	UUID        flows.ChannelUUID `json:"uuid"`
	Name        string            `json:"name"`
	ChannelType flows.ChannelType `json:"type"`
}

func (c *channel) UnmarshalJSON(data []byte) error {
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

func (c *channel) MarshalJSON() ([]byte, error) {
	var ce channelEnvelope

	ce.Name = c.name
	ce.UUID = c.uuid
	ce.ChannelType = c.channelType

	return json.Marshal(ce)
}
