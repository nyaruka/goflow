package flows

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/nyaruka/goflow/utils"
)

// NewResults returns a new empty Results object
func NewResults() *Results {
	return &Results{make(map[string]*Result)}
}

// Results is our wrapper around a map of snakified result names to result objects
type Results struct {
	results map[string]*Result
}

// Save saves a new result in our map. The key is saved in a snakified format
func (r *Results) Save(node NodeUUID, name string, value string, category string, createdOn time.Time) {
	result := Result{node, name, value, category, createdOn}
	r.results[utils.Snakify(name)] = &result
}

// Resolve resolves the passed in key, which is snakified before lookup
func (r *Results) Resolve(key string) interface{} {
	key = utils.Snakify(key)
	result, ok := r.results[key]
	if !ok {
		return nil
	}

	return result
}

// Default returns the default value for our Results, which is the entire map
func (r *Results) Default() interface{} {
	return r
}

// String returns the string representation of our Results, which is a key/value pairing of our fields
func (r *Results) String() string {
	results := make([]string, 0, len(r.results))
	for _, v := range r.results {
		results = append(results, fmt.Sprintf("%s: %s", v.name, v.value))
	}
	return strings.Join(results, ", ")
}

// Result represents a result value in our flow run. Results have a name for which they are the result for,
// the value itself of the result, optional category and the date and node the result was collected on
type Result struct {
	node      NodeUUID
	name      string
	value     string
	category  string
	createdOn time.Time
}

// Resolve resolves the passed in key to a value. Result values have a name, value, category, node and created_on
func (r *Result) Resolve(key string) interface{} {
	switch key {
	case "category":
		return r.category

	case "created_on":
		return r.createdOn

	case "node_uuid":
		return r.node

	case "result_name":
		return r.name

	case "value":
		return r.value
	}

	return fmt.Errorf("No field '%s' on result", key)
}

// Default returns the default value for a result, which is our value
func (r *Result) Default() interface{} {
	return r.value
}

// String returns the string representation of a result, which is our value
func (r *Result) String() string {
	return r.value
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

// UnmarshalJSON is our custom unmarshalling of a Results object, we build our map only with
// with snakified keys
func (r *Results) UnmarshalJSON(data []byte) error {
	r.results = make(map[string]*Result)
	incoming := make(map[string]*Result)
	err := json.Unmarshal(data, &incoming)
	if err != nil {
		return err
	}

	// populate ourselves with the values, but keyed with snakified values
	for k, v := range incoming {
		snaked := utils.Snakify(v.name)
		if k != snaked {
			return fmt.Errorf("invalid results map, key: '%s' does not match snaked result name: '%s'", k, v.name)
		}

		r.results[k] = v
	}
	return nil
}

// MarshalJSON is our custom marshalling of a Results object, we build a map with
// the full names and then marshal that with snakified keys
func (r *Results) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.results)
}

type resultEnvelope struct {
	Node      NodeUUID  `json:"node_uuid"`
	Name      string    `json:"result_name"`
	Value     string    `json:"value"`
	Category  string    `json:"category,omitempty"`
	CreatedOn time.Time `json:"created_on"`
}

// UnmarshalJSON is our custom unmarshalling of a Result object
func (r *Result) UnmarshalJSON(data []byte) error {
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

// MarshalJSON is our custom marshalling of a Result object
func (r *Result) MarshalJSON() ([]byte, error) {
	var re resultEnvelope

	re.Node = r.node
	re.Name = r.name
	re.Value = r.value
	re.Category = r.category
	re.CreatedOn = r.createdOn

	return json.Marshal(re)
}
