package utils

import (
	"bytes"
)

// Typed is an interface of objects that are marshalled as typed envelopes
type Typed interface {
	Type() string
}

type typeOnly struct {
	Type string `json:"type" validate:"required"`
}

// TypedEnvelope represents a json blob with a type property
type TypedEnvelope struct {
	Type string `json:"type" validate:"required"`
	Data []byte `json:"-"`
}

// ReadTypeFromJSON reads a field called `type` from the given JSON
func ReadTypeFromJSON(data []byte) (string, error) {
	t := &typeOnly{}
	if err := UnmarshalAndValidate(data, t); err != nil {
		return "", err
	}
	return t.Type, nil
}

// UnmarshalJSON unmarshals a typed envelope from the given JSON
func (e *TypedEnvelope) UnmarshalJSON(b []byte) error {
	t := &typeOnly{}
	if err := UnmarshalAndValidate(b, t); err != nil {
		return err
	}
	e.Type = t.Type
	e.Data = make([]byte, len(b))
	copy(e.Data, b)
	return nil
}

// MarshalJSON marshals this envelope into JSON
func (e *TypedEnvelope) MarshalJSON() ([]byte, error) {
	// we want the insert the type into our parent data and return that
	t := &typeOnly{Type: e.Type}
	typeBytes, err := JSONMarshal(t)
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

// EnvelopeFromTyped marshals the give object into a typed envelope
func EnvelopeFromTyped(typed Typed) (*TypedEnvelope, error) {
	if typed == nil {
		return nil, nil
	}

	typedData, err := JSONMarshal(typed)
	if err != nil {
		return nil, err
	}

	return &TypedEnvelope{typed.Type(), typedData}, nil
}
