package flows

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type Field struct {
	field     FieldUUID
	name      string
	value     string
	createdOn time.Time
}

type Fields map[string]*Field

func newFields() Fields {
	return make(Fields)
}

func (f Fields) Resolve(key string) interface{} {
	field, ok := f[key]
	if !ok {
		return nil
	}

	return field
}

func (f Fields) Default() interface{} {
	return f
}

func (f Fields) Save(uuid FieldUUID, name string, value string, createdOn time.Time) {
	field := Field{uuid, name, value, createdOn}
	f[strings.ToLower(name)] = &field
}

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

func (f *Field) Default() interface{} {
	return f.value
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type fieldEnvelope struct {
	Field     FieldUUID `json:"field"`
	Name      string    `json:"name"`
	Value     string    `json:"value"`
	CreatedOn time.Time `json:"created_on"`
}

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

func (f *Field) MarshalJSON() ([]byte, error) {
	var fe fieldEnvelope

	fe.Field = f.field
	fe.Name = f.name
	fe.Value = f.value
	fe.CreatedOn = f.createdOn

	return json.Marshal(fe)
}
