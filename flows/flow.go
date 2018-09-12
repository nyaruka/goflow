package flows

import (
	"encoding/json"
	"sync"

	"github.com/nyaruka/goflow/assets"
)

var FlowReader func(json.RawMessage) (Flow, error)

func SetFlowReader(r func(json.RawMessage) (Flow, error)) {
	FlowReader = r
}

// FlowAssets provides access to flow assets, tho unlike other asset managers it lazy loads from its source
type FlowAssets struct {
	byUUID map[assets.FlowUUID]Flow

	mutex  sync.Mutex
	source assets.AssetSource
}

func NewFlowAssets(source assets.AssetSource) *FlowAssets {
	return &FlowAssets{
		byUUID: make(map[assets.FlowUUID]Flow),
		source: source,
	}
}

func (s *FlowAssets) Get(uuid assets.FlowUUID) (Flow, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	flow := s.byUUID[uuid]
	if flow != nil {
		return flow, nil
	}

	asset, err := s.source.Flow(uuid)
	if err != nil {
		return nil, err
	}

	flow, err = FlowReader(asset.Definition())
	if err != nil {
		return nil, err
	}

	s.byUUID[flow.UUID()] = flow
	return flow, nil
}
