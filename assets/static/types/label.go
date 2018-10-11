package types

import (
	"encoding/json"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/utils"
)

// Label is a JSON serializable implementation of a label asset
type Label struct {
	UUID_ assets.LabelUUID `json:"uuid" validate:"required,uuid4"`
	Name_ string           `json:"name"`
}

// NewLabel creates a new label from the passed in UUID and name
func NewLabel(uuid assets.LabelUUID, name string) assets.Label {
	return &Label{UUID_: uuid, Name_: name}
}

// UUID returns the UUID of the label
func (l *Label) UUID() assets.LabelUUID { return l.UUID_ }

// Name returns the name of the label
func (l *Label) Name() string { return l.Name_ }

// ReadLabels reads labels from the given JSON
func ReadLabels(data json.RawMessage) ([]assets.Label, error) {
	var items []*Label
	if err := utils.UnmarshalAndValidate(data, &items); err != nil {
		return nil, err
	}

	asAssets := make([]assets.Label, len(items))
	for i := range items {
		asAssets[i] = items[i]
	}

	return asAssets, nil
}
