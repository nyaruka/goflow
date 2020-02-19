package jsonx

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
)

// Marshal marshals the given object to JSON
func Marshal(v interface{}) ([]byte, error) {
	return marshal(v, "")
}

// MarshalPretty marshals the given object to pretty JSON
func MarshalPretty(v interface{}) ([]byte, error) {
	return marshal(v, "    ")
}

// MarshalMerged marshals the properties of two objects as one object
func MarshalMerged(v1 interface{}, v2 interface{}) ([]byte, error) {
	b1, err := marshal(v1, "")
	if err != nil {
		return nil, err
	}
	b2, err := marshal(v2, "")
	if err != nil {
		return nil, err
	}
	b := append(b1[0:len(b1)-1], byte(','))
	b = append(b, b2[1:]...)
	return b, nil
}

func marshal(v interface{}, indent string) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false) // see https://github.com/golang/go/issues/8592
	encoder.SetIndent("", indent)

	err := encoder.Encode(v)
	if err != nil {
		return nil, err
	}

	// don't include the final \n that .Encode() adds
	data := buffer.Bytes()
	return data[0 : len(data)-1], nil
}

// Unmarshal is just a shortcut for json.Unmarshal so all calls can be made via the jsonx package
func Unmarshal(data json.RawMessage, v interface{}) error {
	return json.Unmarshal(data, v)
}

// UnmarshalArray unmarshals an array of objects from the given JSON
func UnmarshalArray(data json.RawMessage) ([]json.RawMessage, error) {
	var items []json.RawMessage
	err := Unmarshal(data, &items)
	return items, err
}

// UnmarshalWithLimit unmarsmals a struct with a limit on how many bytes can be read from the given reader
func UnmarshalWithLimit(reader io.ReadCloser, s interface{}, limit int64) error {
	body, err := ioutil.ReadAll(io.LimitReader(reader, limit))
	if err != nil {
		return err
	}
	if err := reader.Close(); err != nil {
		return err
	}
	return Unmarshal(body, &s)
}

// DecodeGeneric decodes the given JSON as a generic map or slice
func DecodeGeneric(data []byte) (interface{}, error) {
	var asGeneric interface{}
	decoder := json.NewDecoder(bytes.NewBuffer(data))
	decoder.UseNumber()
	return asGeneric, decoder.Decode(&asGeneric)
}
