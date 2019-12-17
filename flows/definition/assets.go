package definition

import (
	"sync"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/definition/migrations"
)

// implemention of FlowAssets which provides lazy loading and validation of flows
type flowAssets struct {
	byUUID map[assets.FlowUUID]flows.Flow

	mutex  sync.Mutex
	source assets.Source

	migrationConfig *migrations.Config
}

// NewFlowAssets creates a new flow assets
func NewFlowAssets(source assets.Source, migrationConfig *migrations.Config) flows.FlowAssets {
	return &flowAssets{
		byUUID:          make(map[assets.FlowUUID]flows.Flow),
		source:          source,
		migrationConfig: migrationConfig,
	}
}

// Get returns the flow with the given UUID
func (a *flowAssets) Get(uuid assets.FlowUUID) (flows.Flow, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	flow := a.byUUID[uuid]
	if flow != nil {
		return flow, nil
	}

	asset, err := a.source.Flow(uuid)
	if err != nil {
		return nil, err
	}

	flow, err = ReadFlow(asset.Definition(), a.migrationConfig)
	if err != nil {
		return nil, err
	}

	a.byUUID[flow.UUID()] = flow
	return flow, nil
}
