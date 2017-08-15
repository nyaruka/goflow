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

func (a *assets) Validate() error {
	for _, flow := range a.flows {
		if err := flow.Validate(a); err != nil {
			return err
		}
	}
	return nil
}

func (a *assets) GetFlow(uuid flows.FlowUUID) (flows.Flow, error) {
	flow, exists := a.flows[uuid]
	if exists {
		return flow, nil
	}
	return nil, fmt.Errorf("unable to find flow with UUID: %s", uuid)
}

func (a *assets) GetChannel(uuid flows.ChannelUUID) (flows.Channel, error) {
	channel, exists := a.channels[uuid]
	if exists {
		return channel, nil
	}
	return nil, fmt.Errorf("unable to find channel with UUID: %s", uuid)
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
	err := utils.UnmarshalAndValidate(data, &envelope, "assets")
	if err != nil {
		return nil, err
	}

	channelList := make([]flows.Channel, len(envelope.Channels))
	flowList := make([]flows.Flow, len(envelope.Flows))

	for c := range envelope.Channels {
		channel, err := flows.ReadChannel(envelope.Channels[c])
		if err != nil {
			return nil, err
		}
		channelList[c] = channel
	}

	for f := range envelope.Flows {
		flow, err := definition.ReadFlow(envelope.Flows[f])
		if err != nil {
			return nil, err
		}
		flowList[f] = flow
	}

	return NewAssets(flowList, channelList), nil
}
