package definition

import (
	"strings"
	"sync"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/definition/migrations"
)

// implemention of FlowAssets which provides lazy loading and validation of flows
type flowAssets struct {
	cache map[assets.FlowUUID]flows.Flow

	mutex  sync.Mutex
	source assets.Source

	migrationConfig *migrations.Config
}

// NewFlowAssets creates a new flow assets
func NewFlowAssets(source assets.Source, migrationConfig *migrations.Config) flows.FlowAssets {
	return &flowAssets{
		cache:           make(map[assets.FlowUUID]flows.Flow),
		source:          source,
		migrationConfig: migrationConfig,
	}
}

// Get returns the flow with the given UUID
func (a *flowAssets) Get(uuid assets.FlowUUID) (flows.Flow, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	flow := a.cache[uuid]
	if flow != nil {
		return flow, nil
	}

	asset, err := a.source.FlowByUUID(uuid)
	if err != nil {
		return nil, err
	}

	flow, err = ReadAsset(asset, a.migrationConfig)
	if err != nil {
		return nil, err
	}

	a.cache[flow.UUID()] = flow
	return flow, nil
}

// FindByName tries to find a flow with the given name
func (a *flowAssets) FindByName(name string) (flows.Flow, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	for _, flow := range a.cache {
		if strings.EqualFold(flow.Name(), name) {
			return flow, nil
		}
	}

	asset, err := a.source.FlowByName(name)
	if err != nil {
		return nil, err
	}

	flow, err := ReadAsset(asset, a.migrationConfig)
	if err != nil {
		return nil, err
	}

	a.cache[flow.UUID()] = flow
	return flow, nil
}
