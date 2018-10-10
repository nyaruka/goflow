package types

import (
	"encoding/json"
	"fmt"

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

// ReadLabel reads a label from the given JSON
func ReadLabel(data json.RawMessage) (assets.Label, error) {
	l := &Label{}
	if err := utils.UnmarshalAndValidate(data, l); err != nil {
		return nil, fmt.Errorf("unable to read label: %s", err)
	}
	return l, nil
}

// ReadLabels reads labels from the given JSON
func ReadLabels(data json.RawMessage) ([]assets.Label, error) {
	items, err := utils.UnmarshalArray(data)
	if err != nil {
		return nil, err
	}

	labels := make([]assets.Label, len(items))
	for d := range items {
		if labels[d], err = ReadLabel(items[d]); err != nil {
			return nil, err
		}
	}

	return labels, nil
}
