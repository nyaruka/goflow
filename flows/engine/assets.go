package engine

import (
	"encoding/json"
	"fmt"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/definition"
	"github.com/nyaruka/goflow/utils"
)

type assets struct {
	flows    map[flows.FlowUUID]flows.Flow
	channels map[flows.ChannelUUID]flows.Channel
}

func NewAssets(flowList []flows.Flow, channelList []flows.Channel) flows.Assets {
	flowMap := make(map[flows.FlowUUID]flows.Flow, len(channelList))
	for _, flow := range flowList {
		flowMap[flow.UUID()] = flow
	}

	channelMap := make(map[flows.ChannelUUID]flows.Channel, len(channelList))
	for _, channel := range channelList {
		channelMap[channel.UUID()] = channel
	}

	return &assets{flows: flowMap, channels: channelMap}
}

func (e *assets) GetFlow(uuid flows.FlowUUID) (flows.Flow, error) {
	flow, exists := e.flows[uuid]
	if exists {
		return flow, nil
	}
	return nil, fmt.Errorf("unable to find flow with UUID: %s", uuid)
}

func (e *assets) GetChannel(uuid flows.ChannelUUID) (flows.Channel, error) {
	channel, exists := e.channels[uuid]
	if exists {
		return channel, nil
	}
	return nil, fmt.Errorf("unable to find channel with UUID: %s %d", uuid)
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type assetsEnvelope struct {
	Flows    []json.RawMessage `json:"flows"               validate:"required"`
	Channels []json.RawMessage `json:"channels,omitempty"`
}

func ReadAssets(data json.RawMessage) (flows.Assets, error) {
	var envelope assetsEnvelope
	err := json.Unmarshal(data, &envelope)
	if err != nil {
		return nil, err
	}

	err = utils.ValidateAs(&envelope, "assets")
	if err != nil {
		return nil, err
	}

	flowList := make([]flows.Flow, len(envelope.Flows))
	channelList := make([]flows.Channel, len(envelope.Channels))

	for f := range envelope.Flows {
		flow, err := definition.ReadFlow(envelope.Flows[f])
		if err != nil {
			return nil, err
		}
		flowList[f] = flow
	}
	for c := range envelope.Channels {
		channel, err := flows.ReadChannel(envelope.Channels[c])
		if err != nil {
			return nil, err
		}
		channelList[c] = channel
	}

	return NewAssets(flowList, channelList), nil
}
