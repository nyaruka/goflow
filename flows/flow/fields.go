package flow

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/nyaruka/goflow/flows"
)

type field struct {
	field     flows.FieldUUID
	name      string
	value     string
	createdOn time.Time
}

type fields map[string]*field

func newFields() fields {
	return make(fields)
}

func (f fields) Resolve(key string) interface{} {
	field, ok := f[key]
	if !ok {
		return nil
	}

	return field
}

func (f fields) Default() interface{} {
	return f
}

func (f fields) Save(uuid flows.FieldUUID, name string, value string, createdOn time.Time) {
	field := field{uuid, name, value, createdOn}
	f[strings.ToLower(name)] = &field
}

func (f *field) Resolve(key string) interface{} {
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

func (f *field) Default() interface{} {
	return f.value
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type fieldEnvelope struct {
	Field     flows.FieldUUID `json:"field"`
	Name      string          `json:"name"`
	Value     string          `json:"value"`
	CreatedOn time.Time       `json:"created_on"`
}

func (f *field) UnmarshalJSON(data []byte) error {
	var fe fieldEnvelope
	var err error

	err = json.Unmarshal(data, &fe)
	f.field = fe.Field
	f.name = fe.Name
	f.value = fe.Value
	f.createdOn = fe.CreatedOn

	return err
}

func (f *field) MarshalJSON() ([]byte, error) {
	var fe fieldEnvelope

	fe.Field = f.field
	fe.Name = f.name
	fe.Value = f.value
	fe.CreatedOn = f.createdOn

	return json.Marshal(fe)
}
