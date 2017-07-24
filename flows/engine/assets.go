package engine

import (
	"encoding/json"
	"fmt"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/definition"
)

type assets struct {
	flows    map[flows.FlowUUID]flows.Flow
	channels map[flows.ChannelUUID]flows.Channel
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
	Flows    []json.RawMessage `json:"flows"                validate:"required"`
	Channels []json.RawMessage `json:"channels,omitempty"`
}

func ReadAssets(data json.RawMessage) (flows.Assets, error) {
	var envelope assetsEnvelope
	if err := json.Unmarshal(data, &envelope); err != nil {
		return nil, err
	}

	flowMap := make(map[flows.FlowUUID]flows.Flow, len(envelope.Flows))
	for f := range envelope.Flows {
		flow, err := definition.ReadFlow(envelope.Flows[f])
		if err != nil {
			return nil, err
		}
		flowMap[flow.UUID()] = flow
	}

	channelMap := make(map[flows.ChannelUUID]flows.Channel, len(envelope.Channels))
	for c := range envelope.Channels {
		channel, err := flows.ReadChannel(envelope.Channels[c])
		if err != nil {
			return nil, err
		}
		channelMap[channel.UUID()] = channel
	}

	return &assets{flowMap, channelMap}, nil
}
