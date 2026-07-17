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
	cache sync.Map // assets.FlowUUID -> flows.Flow

	source          assets.Source
	migrationConfig *migrations.Config
}

// NewFlowAssets creates a new flow assets
func NewFlowAssets(source assets.Source, migrationConfig *migrations.Config) flows.FlowAssets {
	return &flowAssets{
		source:          source,
		migrationConfig: migrationConfig,
	}
}

// Get returns the flow with the given UUID
func (a *flowAssets) Get(uuid assets.FlowUUID) (flows.Flow, error) {
	if flow, ok := a.cache.Load(uuid); ok {
		return flow.(flows.Flow), nil
	}

	asset, err := a.source.FlowByUUID(uuid)
	if err != nil {
		return nil, err
	}

	flow, err := ReadAsset(asset, a.migrationConfig)
	if err != nil {
		return nil, err
	}

	return a.cached(flow), nil
}

// FindByName tries to find a flow with the given name
func (a *flowAssets) FindByName(name string) (flows.Flow, error) {
	var found flows.Flow
	a.cache.Range(func(_, v any) bool {
		if flow := v.(flows.Flow); strings.EqualFold(flow.Name(), name) {
			found = flow
			return false
		}
		return true
	})
	if found != nil {
		return found, nil
	}

	asset, err := a.source.FlowByName(name)
	if err != nil {
		return nil, err
	}

	flow, err := ReadAsset(asset, a.migrationConfig)
	if err != nil {
		return nil, err
	}

	return a.cached(flow), nil
}

// concurrent misses may both read a flow but this ensures all callers get the same instance
func (a *flowAssets) cached(flow flows.Flow) flows.Flow {
	actual, _ := a.cache.LoadOrStore(flow.UUID(), flow)
	return actual.(flows.Flow)
}
