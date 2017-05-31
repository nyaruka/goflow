package flows

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/nyaruka/goflow/utils"
)

// NewFields returns a new empty field
func NewFields() *Fields {
	return &Fields{make(map[string]*Field)}
}

// Fields represents a set of fields keyed by snakified field names
type Fields struct {
	fields map[string]*Field
}

// Save saves a new field to our map
func (f *Fields) Save(uuid FieldUUID, name string, value string, createdOn time.Time) {
	field := Field{uuid, name, value, createdOn}
	f.fields[utils.Snakify(name)] = &field
}

// Resolve resolves the field for the passed in key which will be snakified
func (f *Fields) Resolve(key string) interface{} {
	key = utils.Snakify(key)
	value, ok := f.fields[key]
	if !ok {
		return nil
	}

	return value
}

// Default returns the default value for Fields, which is ourselves
func (f *Fields) Default() interface{} {
	return f
}

// String returns the string representation of these Fields, which is our JSON representation
func (f *Fields) String() string {
	fieldJSON, err := json.Marshal(f)
	if err != nil {
		return err.Error()
	}
	return string(fieldJSON)
}

// Field represents a contact field and value for a contact
type Field struct {
	field     FieldUUID
	name      string
	value     string
	createdOn time.Time
}

// Resolve resolves one of the fields on a Field
func (f *Field) Resolve(key string) interface{} {
	switch key {

	case "field":
		return f.field

	case "name":
		return f.name

	case "value":
		return f.value

	case "created_on":
		return f.createdOn

	}

	return fmt.Errorf("No field '%s' on contact field", key)
}

// Default returns the default value for a Field, which is the field itself
func (f *Field) Default() interface{} {
	return f
}

// String returns the string value for a field, which is our value
func (f *Field) String() string {
	return f.value
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

// UnmarshalJSON is our custom unmarshalling of a Fields object, we build our map with snakified keys
func (f *Fields) UnmarshalJSON(data []byte) error {
	f.fields = make(map[string]*Field)
	incoming := make(map[string]*Field)
	err := json.Unmarshal(data, &incoming)
	if err != nil {
		return nil
	}

	// populate ourselves with the fields, but keyed with snakified values
	for k, v := range incoming {
		// validate name and key are the same
		if k != v.name {
			return fmt.Errorf("invalid fields map, key: '%s' does not match field name: '%s'", k, v.name)
		}

		f.fields[utils.Snakify(v.name)] = v
	}
	return nil
}

// MarshalJSON is our custom marshalling of a Fields object, we build a map with
// the full names and then marshal that with snakified keys
func (f *Fields) MarshalJSON() ([]byte, error) {
	outgoing := make(map[string]*Field)
	for _, v := range f.fields {
		outgoing[v.name] = v
	}
	return json.Marshal(outgoing)
}

type fieldEnvelope struct {
	Field     FieldUUID `json:"field"`
	Name      string    `json:"name"`
	Value     string    `json:"value"`
	CreatedOn time.Time `json:"created_on"`
}

// UnmarshalJSON is our custom unmarshalling for Field
func (f *Field) UnmarshalJSON(data []byte) error {
	var fe fieldEnvelope
	var err error

	err = json.Unmarshal(data, &fe)
	f.field = fe.Field
	f.name = fe.Name
	f.value = fe.Value
	f.createdOn = fe.CreatedOn

	return err
}

// MarshalJSON is our custom unmarshalling for Field
func (f *Field) MarshalJSON() ([]byte, error) {
	var fe fieldEnvelope

	fe.Field = f.field
	fe.Name = f.name
	fe.Value = f.value
	fe.CreatedOn = f.createdOn

	return json.Marshal(fe)
}
