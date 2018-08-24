package engine

import (
	"encoding/json"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/definition"
)

const (
	assetTypeChannel           assets.AssetType = "channel"
	assetTypeField             assets.AssetType = "field"
	assetTypeFlow              assets.AssetType = "flow"
	assetTypeGroup             assets.AssetType = "group"
	assetTypeLabel             assets.AssetType = "label"
	assetTypeLocationHierarchy assets.AssetType = "location_hierarchy"
	assetTypeResthook          assets.AssetType = "resthook"
)

func init() {
	assets.RegisterType(assetTypeChannel, true, func(data json.RawMessage) (interface{}, error) { return flows.ReadChannelSet(data) })
	assets.RegisterType(assetTypeField, true, func(data json.RawMessage) (interface{}, error) { return flows.ReadFieldSet(data) })
	assets.RegisterType(assetTypeFlow, false, func(data json.RawMessage) (interface{}, error) { return definition.ReadFlow(data) })
	assets.RegisterType(assetTypeGroup, true, func(data json.RawMessage) (interface{}, error) { return flows.ReadGroupSet(data) })
	assets.RegisterType(assetTypeLabel, true, func(data json.RawMessage) (interface{}, error) { return flows.ReadLabelSet(data) })
	assets.RegisterType(assetTypeLocationHierarchy, true, func(data json.RawMessage) (interface{}, error) { return flows.ReadLocationHierarchySet(data) })
	assets.RegisterType(assetTypeResthook, true, func(data json.RawMessage) (interface{}, error) { return flows.ReadResthookSet(data) })
}
