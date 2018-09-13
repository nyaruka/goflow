package types

import (
	"encoding/json"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/utils"
)

// ReadLocationHierarchies reads location hierarchies from the given JSON
func ReadLocationHierarchies(data json.RawMessage) ([]assets.LocationHierarchy, error) {
	items, err := utils.UnmarshalArray(data)
	if err != nil {
		return nil, err
	}

	hierarchies := make([]assets.LocationHierarchy, len(items))
	for d := range items {
		if hierarchies[d], err = utils.ReadLocationHierarchy(items[d]); err != nil {
			return nil, err
		}
	}

	return hierarchies, nil
}
