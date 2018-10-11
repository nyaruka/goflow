package rest

import (
	"encoding/json"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/static/types"
	"github.com/nyaruka/goflow/utils"
)

// reads channels from the given JSON
func readChannels(data json.RawMessage) (interface{}, error) {
	var items []*types.Channel
	if err := utils.UnmarshalAndValidate(data, &items); err != nil {
		return nil, err
	}

	asAssets := make([]assets.Channel, len(items))
	for i := range items {
		asAssets[i] = items[i]
	}

	return asAssets, nil
}

// reads fields from the given JSON
func readFields(data json.RawMessage) (interface{}, error) {
	var items []*types.Field
	if err := utils.UnmarshalAndValidate(data, &items); err != nil {
		return nil, err
	}

	asAssets := make([]assets.Field, len(items))
	for i := range items {
		asAssets[i] = items[i]
	}

	return asAssets, nil
}

// reads a flow from the given JSON
func readFlow(data json.RawMessage) (interface{}, error) {
	f := &types.Flow{Definition_: data}
	if err := utils.UnmarshalAndValidate(data, f); err != nil {
		return nil, err
	}
	return f, nil
}

// reads groups from the given JSON
func readGroups(data json.RawMessage) (interface{}, error) {
	var items []*types.Group
	if err := utils.UnmarshalAndValidate(data, &items); err != nil {
		return nil, err
	}

	asAssets := make([]assets.Group, len(items))
	for i := range items {
		asAssets[i] = items[i]
	}

	return asAssets, nil
}

// reads labels from the given JSON
func readLabels(data json.RawMessage) (interface{}, error) {
	var items []*types.Label
	if err := utils.UnmarshalAndValidate(data, &items); err != nil {
		return nil, err
	}

	asAssets := make([]assets.Label, len(items))
	for i := range items {
		asAssets[i] = items[i]
	}

	return asAssets, nil
}

// reads location hierarchies from the given JSON
func readLocationHierarchies(data json.RawMessage) (interface{}, error) {
	var items []*utils.LocationHierarchy
	if err := utils.UnmarshalAndValidate(data, &items); err != nil {
		return nil, err
	}

	asAssets := make([]assets.LocationHierarchy, len(items))
	for i := range items {
		asAssets[i] = items[i]
	}

	return asAssets, nil
}

// reads a resthook set from the given JSON
func readResthooks(data json.RawMessage) (interface{}, error) {
	var items []*types.Resthook
	if err := utils.UnmarshalAndValidate(data, &items); err != nil {
		return nil, err
	}

	asAssets := make([]assets.Resthook, len(items))
	for i := range items {
		asAssets[i] = items[i]
	}

	return asAssets, nil
}
