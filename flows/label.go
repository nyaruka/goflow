package flows

import (
	"encoding/json"
)

type Label struct {
	uuid LabelUUID
	name string
}

// UUID returns the UUID of this label
func (l *Label) UUID() LabelUUID { return l.uuid }

// Name returns the name of this label
func (l *Label) Name() string { return l.name }

// NewLabel creates a new label given the passed in uuid and name
func NewLabel(uuid LabelUUID, name string) *Label {
	return &Label{uuid, name}
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type labelEnvelope struct {
	UUID LabelUUID `json:"uuid"`
	Name string    `json:"name"`
}

func (l *Label) UnmarshalJSON(data []byte) error {
	var le labelEnvelope
	var err error

	err = json.Unmarshal(data, &le)
	l.uuid = le.UUID
	l.name = le.Name

	return err
}

func (l *Label) MarshalJSON() ([]byte, error) {
	var le labelEnvelope

	le.Name = l.name
	le.UUID = l.uuid

	return json.Marshal(le)
}
