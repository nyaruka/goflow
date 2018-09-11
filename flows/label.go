package flows

import (
	"encoding/json"
	"fmt"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/utils"
)

// Label represents a message label
type Label struct {
	uuid assets.LabelUUID
	name string
}

// NewLabel creates a new label given the passed in uuid and name
func NewLabel(uuid assets.LabelUUID, name string) *Label {
	return &Label{uuid, name}
}

// UUID returns the UUID of this label
func (l *Label) UUID() assets.LabelUUID { return l.uuid }

// Name returns the name of this label
func (l *Label) Name() string { return l.name }

// Reference returns a reference to this label
func (l *Label) Reference() *LabelReference { return NewLabelReference(l.uuid, l.name) }

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type labelEnvelope struct {
	UUID assets.LabelUUID `json:"uuid" validate:"required,uuid4"`
	Name string           `json:"name"`
}

// ReadLabel reads a label asset from the given JSON
func ReadLabel(data json.RawMessage) (assets.Label, error) {
	var le labelEnvelope
	if err := utils.UnmarshalAndValidate(data, &le); err != nil {
		return nil, fmt.Errorf("unable to read label: %s", err)
	}

	return NewLabel(le.UUID, le.Name), nil
}

// ReadLabels reads an array of labels from the given JSON
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
