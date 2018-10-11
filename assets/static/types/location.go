package types

import (
	"encoding/json"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/utils"
)

// ReadLocationHierarchies reads location hierarchies from the given JSON
func ReadLocationHierarchies(data json.RawMessage) ([]assets.LocationHierarchy, error) {
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
