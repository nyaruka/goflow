package events

import (
	"encoding/json"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	utils.RegisterValidatorAlias("result_name", "min=1,max=64",
		func(validator.FieldError) string { return "is not a valid result name" },
	)

	// editor enforces max length of 36 for user defined categories but routers like split by group can set the category to a longer value like a group name
	utils.RegisterValidatorAlias("result_category", "min=1,max=64",
		func(validator.FieldError) string { return "is not a valid result category" },
	)
}

// Result describes a value captured during a run's execution. It might have been implicitly created by a router, or explicitly
// created by a [set_run_result](#action:set_run_result) action.
type Result struct {
	Name              string          `json:"name" validate:"required"` // TODO add result_name validation when we're sure sessions no longer have invalid result names
	Value             string          `json:"value"`
	Category          string          `json:"category,omitempty"`
	CategoryLocalized string          `json:"category_localized,omitempty"`
	NodeUUID          NodeUUID        `json:"node_uuid"`
	Input             string          `json:"input,omitempty"` // should be called operand but too late now
	Extra             json.RawMessage `json:"extra,omitempty"`
	CreatedOn         time.Time       `json:"created_on" validate:"required"`
}

// NewResult creates a new result
func NewResult(name string, value string, category string, categoryLocalized string, nodeUUID NodeUUID, input string, extra []byte, createdOn time.Time) *Result {
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
//	__default__:text -> the value
//	name:text -> the name of the result
//	value:text -> the value of the result
//	category:text -> the category of the result
//	category_localized:text -> the localized category of the result
//	input:text -> the input of the result
//	node_uuid:text -> the UUID of the node in the flow that generated the result
//	created_on:datetime -> the creation date of the result
//
// @context result
func (r *Result) Context(env envs.Environment) map[string]types.XValue {
	categoryLocalized := r.CategoryLocalized
	if categoryLocalized == "" {
		categoryLocalized = r.Category
	}

	return map[string]types.XValue{
		"__default__":        types.NewXText(r.Value),
		"name":               types.NewXText(r.Name),
		"value":              types.NewXText(r.Value),
		"category":           types.NewXText(r.Category),
		"category_localized": types.NewXText(categoryLocalized),
		"input":              types.NewXText(r.Input),
		"extra":              types.JSONToXValue(r.Extra),
		"node_uuid":          types.NewXText(string(r.NodeUUID)),
		"created_on":         types.NewXDateTime(r.CreatedOn),
	}
}
