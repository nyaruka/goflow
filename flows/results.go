package flows

import (
	"fmt"
	"strings"
	"time"

	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/utils"
)

// Result represents a result value in our flow run. Results have a name for which they are the result for,
// the value itself of the result, optional category and the date and node the result was collected on
type Result struct {
	Name              string    `json:"name"`
	Value             string    `json:"value"`
	Category          string    `json:"category,omitempty"`
	CategoryLocalized string    `json:"category_localized,omitempty"`
	NodeUUID          NodeUUID  `json:"node_uuid"`
	Input             string    `json:"input"`
	CreatedOn         time.Time `json:"created_on"`
}

// Resolve resolves the passed in key to a value. Result values have a name, value, category, node and created_on
func (r *Result) Resolve(key string) interface{} {
	switch key {
	case "name":
		return r.Name
	case "value":
		return r.Value
	case "category":
		return r.Category
	case "category_localized":
		if r.CategoryLocalized == "" {
			return r.Category
		}
		return r.CategoryLocalized
	case "created_on":
		return r.CreatedOn
	}

	return fmt.Errorf("no field '%s' on result", key)
}

// Atomize is called when this object needs to be reduced to a primitive
func (r *Result) Atomize() interface{} {
	return r.Value
}

var _ types.Atomizable = (*Result)(nil)
var _ types.Resolvable = (*Result)(nil)

// Results is our wrapper around a map of snakified result names to result objects
type Results map[string]*Result

// NewResults creates a new empty set of results
func NewResults() Results {
	return make(Results, 0)
}

// Clone returns a clone of this results set
func (r Results) Clone() Results {
	clone := make(Results, len(r))
	for k, v := range r {
		clone[k] = v
	}
	return clone
}

// Save saves a new result in our map. The key is saved in a snakified format
func (r Results) Save(name string, value string, category string, categoryLocalized string, nodeUUID NodeUUID, input string, createdOn time.Time) {
	r[utils.Snakify(name)] = &Result{
		Name:              name,
		Value:             value,
		Category:          category,
		CategoryLocalized: categoryLocalized,
		NodeUUID:          nodeUUID,
		Input:             input,
		CreatedOn:         createdOn,
	}
}

// Length is called to get the length of this object
func (r Results) Length() int {
	return len(r)
}

// Resolve resolves the passed in key, which is snakified before lookup
func (r Results) Resolve(key string) interface{} {
	key = utils.Snakify(key)

	result, exists := r[key]
	if !exists {
		return fmt.Errorf("no such run result '%s'", key)
	}
	return result
}

// Atomize is called when this object needs to be reduced to a primitive
func (r Results) Atomize() interface{} {
	results := make([]string, 0, len(r))
	for _, v := range r {
		results = append(results, fmt.Sprintf("%s: %s", v.Name, v.Value))
	}
	return strings.Join(results, ", ")
}

var _ types.Atomizable = (Results)(nil)
var _ types.Lengthable = (Results)(nil)
var _ types.Resolvable = (Results)(nil)
