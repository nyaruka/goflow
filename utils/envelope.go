package utils

import (
	"bytes"
	"encoding/json"
)

type Typed interface {
	Type() string
}

// TypedEnvelope represents a json blob with a type property
type TypedEnvelope struct {
	Type string `json:"type"`
	Data []byte `json:"-"`
}

func (e *TypedEnvelope) UnmarshalJSON(b []byte) (err error) {
	typeE := &struct {
		Type string `json:"type"`
	}{}
	err = json.Unmarshal(b, &typeE)
	if err != nil {
		return err
	}
	e.Type = typeE.Type
	e.Data = make([]byte, len(b))
	copy(e.Data, b)

	return err
}

func (e *TypedEnvelope) MarshalJSON() ([]byte, error) {
	// we want the insert the type into our parent data and return that
	typeE := &struct {
		Type string `json:"type"`
	}{Type: e.Type}
	typeBytes, err := json.Marshal(&typeE)
	if err != nil {
		return nil, err
	}

	// empty case {}
	if len(e.Data) == 2 {
		return typeBytes, nil
	}

	data := bytes.NewBuffer(typeBytes)

	// remove ending }
	data.Truncate(data.Len() - 1)

	// add ,
	data.WriteByte(',')

	// copy in our data, skipping over the leading {
	data.Write(e.Data[1:])

	return data.Bytes(), nil
}

func EnvelopeFromTyped(typed Typed) (*TypedEnvelope, error) {
	if typed == nil {
		return nil, nil
	}

	typedData, err := json.Marshal(typed)
	if err != nil {
		return nil, err
	}

	envelope := TypedEnvelope{typed.Type(), typedData}
	return &envelope, nil
}
