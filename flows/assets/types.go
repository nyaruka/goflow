package assets

import (
	"encoding/json"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/definition"
)

type assetType string

const (
	assetTypeChannel           assetType = "channel"
	assetTypeField             assetType = "field"
	assetTypeFlow              assetType = "flow"
	assetTypeGroup             assetType = "group"
	assetTypeLabel             assetType = "label"
	assetTypeLocationHierarchy assetType = "location_hierarchy"
	assetTypeResthook          assetType = "resthook"
)

type assetReader func(data json.RawMessage) (interface{}, error)

type assetTypeConfig struct {
	manageAsSet bool
	reader      assetReader
}

var typeConfigs = map[assetType]*assetTypeConfig{}

func registerAssetType(name assetType, manageAsSet bool, reader assetReader) {
	typeConfigs[name] = &assetTypeConfig{manageAsSet: manageAsSet, reader: reader}
}

func init() {
	registerAssetType(assetTypeChannel, true, func(data json.RawMessage) (interface{}, error) { return flows.ReadChannelSet(data) })
	registerAssetType(assetTypeField, true, func(data json.RawMessage) (interface{}, error) { return flows.ReadFieldSet(data) })
	registerAssetType(assetTypeFlow, false, func(data json.RawMessage) (interface{}, error) { return definition.ReadFlow(data) })
	registerAssetType(assetTypeGroup, true, func(data json.RawMessage) (interface{}, error) { return flows.ReadGroupSet(data) })
	registerAssetType(assetTypeLabel, true, func(data json.RawMessage) (interface{}, error) { return flows.ReadLabelSet(data) })
	registerAssetType(assetTypeLocationHierarchy, true, func(data json.RawMessage) (interface{}, error) { return flows.ReadLocationHierarchySet(data) })
	registerAssetType(assetTypeResthook, true, func(data json.RawMessage) (interface{}, error) { return flows.ReadResthookSet(data) })
}
