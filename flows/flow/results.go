package flow

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/nyaruka/goflow/flows"
)

type result struct {
	node      flows.NodeUUID
	name      string
	value     string
	category  string
	createdOn time.Time
}

type results map[string]*result

func newResults() results {
	return make(results)
}

func (r results) Resolve(key string) interface{} {
	result, ok := r[key]
	if !ok {
		return nil
	}

	return result
}

func (r results) Default() interface{} {
	return r
}

func (r results) Save(node flows.NodeUUID, name string, value string, category string, createdOn time.Time) {
	result := result{node, name, value, category, createdOn}
	r[strings.ToLower(name)] = &result
}

func (r *result) Resolve(key string) interface{} {
	switch key {

	case "name":
		return r.name

	case "value":
		return r.value

	case "category":
		return r.category

	case "node":
		return r.node
	}

	return fmt.Errorf("No field '%s' on result", key)
}

func (r *result) Default() interface{} {
	return r.value
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type resultEnvelope struct {
	Node      flows.NodeUUID `json:"node"`
	Name      string         `json:"name"`
	Value     string         `json:"value"`
	Category  string         `json:"category,omitempty"`
	CreatedOn time.Time      `json:"created_on"`
}

func (r *result) UnmarshalJSON(data []byte) error {
	var re resultEnvelope
	var err error

	err = json.Unmarshal(data, &re)
	r.node = re.Node
	r.name = re.Name
	r.value = re.Value
	r.category = re.Category
	r.createdOn = re.CreatedOn

	return err
}

func (r *result) MarshalJSON() ([]byte, error) {
	var re resultEnvelope

	re.Node = r.node
	re.Name = r.name
	re.Value = r.value
	re.Category = r.category
	re.CreatedOn = r.createdOn

	return json.Marshal(re)
}
