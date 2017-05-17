package flows

import (
	"encoding/json"
)

type Label struct {
	UUID LabelUUID
	Name string
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//-------n-----------------------------------------------------------------------------------

type labelEnvelope struct {
	UUID LabelUUID `json:"uuid"`
	Name string    `json:"name"`
}

func (l *Label) UnmarshalJSON(data []byte) error {
	var le labelEnvelope
	var err error

	err = json.Unmarshal(data, &le)
	l.UUID = le.UUID
	l.Name = le.Name

	return err
}

func (l *Label) MarshalJSON() ([]byte, error) {
	var le labelEnvelope

	le.Name = l.Name
	le.UUID = l.UUID

	return json.Marshal(le)
}
