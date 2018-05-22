package flows

import (
	"time"

	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/utils"
)

// Result describes a value captured during a run's execution. It might have been implicitly created by a router, or explicitly
// created by a [set_run_result](#action:set_run_result) action.It renders as its value in a template, and has the following
// properties which can be accessed:
//
//  * `value` the value of the result
//  * `category` the category of the result
//  * `category_localized` the localized category of the result
//  * `input` the input associated with the result
//  * `node_uuid` the UUID of the node where the result was created
//  * `created_on` the time when the result was created
//
// Examples:
//
//   @run.results.favorite_color -> red
//   @run.results.favorite_color.value -> red
//   @run.results.favorite_color.category -> Red
//
// @context result
type Result struct {
	Name              string    `json:"name"`
	Value             string    `json:"value"`
	Category          string    `json:"category,omitempty"`
	CategoryLocalized string    `json:"category_localized,omitempty"`
	NodeUUID          NodeUUID  `json:"node_uuid"`
	Input             *string   `json:"input,omitempty"`
	CreatedOn         time.Time `json:"created_on"`
}

// Resolve resolves the passed in key to a value. Result values have a name, value, category, node and created_on
func (r *Result) Resolve(env utils.Environment, key string) types.XValue {
	switch key {
	case "name":
		return types.NewXText(r.Name)
	case "value":
		return types.NewXText(r.Value)
	case "category":
		return types.NewXText(r.Category)
	case "category_localized":
		if r.CategoryLocalized == "" {
			return types.NewXText(r.Category)
		}
		return types.NewXText(r.CategoryLocalized)
	case "input":
		if r.Input != nil {
			return types.NewXText(*r.Input)
		}
		return nil
	case "node_uuid":
		return types.NewXText(string(r.NodeUUID))
	case "created_on":
		return types.NewXDateTime(r.CreatedOn)
	}

	return types.NewXResolveError(r, key)
}

// Describe returns a representation of this type for error messages
func (r *Result) Describe() string { return "run result" }

// Reduce is called when this object needs to be reduced to a primitive
func (r *Result) Reduce(env utils.Environment) types.XPrimitive {
	return types.NewXText(r.Value)
}

// ToXJSON is called when this type is passed to @(json(...))
func (r *Result) ToXJSON(env utils.Environment) types.XText {
	return types.ResolveKeys(env, r, "name", "value", "category", "category_localized", "input", "node_uuid", "created_on").ToXJSON(env)
}

var _ types.XValue = (*Result)(nil)
var _ types.XResolvable = (*Result)(nil)

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
func (r Results) Save(name string, value string, category string, categoryLocalized string, nodeUUID NodeUUID, input *string, createdOn time.Time) {
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

func (r Results) Get(key string) *Result {
	return r[key]
}

// Length is called to get the length of this object
func (r Results) Length() int {
	return len(r)
}

// Resolve resolves the passed in key, which is snakified before lookup
func (r Results) Resolve(env utils.Environment, key string) types.XValue {
	key = utils.Snakify(key)

	result, exists := r[key]
	if !exists {
		return types.NewXErrorf("no such run result '%s'", key)
	}
	return result
}

// Describe returns a representation of this type for error messages
func (r Results) Describe() string { return "run results" }

// Reduce is called when this object needs to be reduced to a primitive
func (r Results) Reduce(env utils.Environment) types.XPrimitive {
	results := types.NewEmptyXMap()
	for k, v := range r {
		results.Put(k, v)
	}
	return results
}

// ToXJSON is called when this type is passed to @(json(...))
func (r Results) ToXJSON(env utils.Environment) types.XText {
	return r.Reduce(env).ToXJSON(env)
}

var _ types.XValue = (Results)(nil)
var _ types.XLengthable = (Results)(nil)
var _ types.XResolvable = (Results)(nil)
