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
func (r *Result) Resolve(key string) types.XValue {
	switch key {
	case "name":
		return types.NewXString(r.Name)
	case "value":
		return types.NewXString(r.Value)
	case "category":
		return types.NewXString(r.Category)
	case "category_localized":
		if r.CategoryLocalized == "" {
			return types.NewXString(r.Category)
		}
		return types.NewXString(r.CategoryLocalized)
	case "created_on":
		return types.NewXTime(r.CreatedOn)
	}

	return types.NewXResolveError(r, key)
}

// Reduce is called when this object needs to be reduced to a primitive
func (r *Result) Reduce() types.XPrimitive {
	return types.NewXString(r.Value)
}

func (r *Result) ToJSON() types.XString { return types.RequireMarshalToXString(r.Value) }

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
func (r Results) Resolve(key string) types.XValue {
	key = utils.Snakify(key)

	result, exists := r[key]
	if !exists {
		return types.NewXResolveError(r, key)
	}
	return result
}

// Reduce is called when this object needs to be reduced to a primitive
func (r Results) Reduce() types.XPrimitive {
	results := make([]string, 0, len(r))
	for _, v := range r {
		results = append(results, fmt.Sprintf("%s: %s", v.Name, v.Value))
	}
	return types.NewXString(strings.Join(results, ", "))
}

func (r Results) ToJSON() types.XString { return types.NewXString("TODO") }

var _ types.XValue = (Results)(nil)
var _ types.XLengthable = (Results)(nil)
var _ types.XResolvable = (Results)(nil)
