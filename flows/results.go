package flows

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/utils"
)

// Result describes a value captured during a run's execution. It might have been implicitly created by a router, or explicitly
// created by a [set_run_result](#action:set_run_result) action.
type Result struct {
	Name              string          `json:"name" validate:"required"`
	Value             string          `json:"value"`
	Category          string          `json:"category,omitempty"`
	CategoryLocalized string          `json:"category_localized,omitempty"`
	NodeUUID          NodeUUID        `json:"node_uuid"`
	Input             string          `json:"input,omitempty"`
	Extra             json.RawMessage `json:"extra,omitempty"`
	CreatedOn         time.Time       `json:"created_on" validate:"required"`
}

// NewResult creates a new result
func NewResult(name string, value string, category string, categoryLocalized string, nodeUUID NodeUUID, input string, extra json.RawMessage, createdOn time.Time) *Result {
	return &Result{
		Name:              name,
		Value:             value,
		Category:          category,
		CategoryLocalized: categoryLocalized,
		NodeUUID:          nodeUUID,
		Input:             input,
		Extra:             extra,
		CreatedOn:         createdOn,
	}
}

// Context returns the properties available in expressions
//
//   __default__:text -> the value
//   name:text -> the name of the result
//   value:text -> the value of the result
//   category:text -> the category of the result
//   category_localized:text -> the localized category of the result
//   input:text -> the input of the result
//   extra:any -> the extra data of the result such as a webhook response
//   node_uuid:text -> the UUID of the node in the flow that generated the result
//   created_on:datetime -> the creation date of the result
//
// @context result
func (r *Result) Context(env envs.Environment) map[string]types.XValue {
	categoryLocalized := r.CategoryLocalized
	if categoryLocalized == "" {
		categoryLocalized = r.Category
	}

	return map[string]types.XValue{
		"__default__":          types.NewXText(r.Value),
		"name":                 types.NewXText(r.Name),
		"value":                types.NewXText(r.Value),
		"values":               types.NewXArray(types.NewXText(r.Value)),
		"category":             types.NewXText(r.Category),
		"categories":           types.NewXArray(types.NewXText(r.Category)),
		"category_localized":   types.NewXText(categoryLocalized),
		"categories_localized": types.NewXArray(types.NewXText(categoryLocalized)),
		"input":                types.NewXText(r.Input),
		"extra":                types.JSONToXValue(r.Extra),
		"node_uuid":            types.NewXText(string(r.NodeUUID)),
		"created_on":           types.NewXDateTime(r.CreatedOn),
	}
}

// Results is our wrapper around a map of snakified result names to result objects
type Results map[string]*Result

// NewResults creates a new empty set of results
func NewResults() Results {
	return make(Results)
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
func (r Results) Save(result *Result) {
	r[utils.Snakify(result.Name)] = result
}

// Get returns the result with the given key
func (r Results) Get(key string) *Result {
	return r[key]
}

// Context returns the properties available in expressions
func (r Results) Context(env envs.Environment) map[string]types.XValue {
	entries := make(map[string]types.XValue, len(r)+1)
	entries["__default__"] = types.NewXText(r.format())

	for k, v := range r {
		entries[k] = Context(env, v)
	}
	return entries
}

func (r Results) format() string {
	lines := make([]string, 0, len(r))
	for _, v := range r {
		lines = append(lines, fmt.Sprintf("%s: %s", v.Name, v.Value))
	}

	sort.Strings(lines)
	return strings.Join(lines, "\n")
}
