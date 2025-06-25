package cases

import (
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
)

// HasCategory tests whether the category of a result on of the passed in `categories`
func HasCategory(env envs.Environment, resultObj *types.XObject, categories ...*types.XText) types.XValue {
	result, err := resultFromXObject(resultObj)
	if err != nil {
		return types.NewXErrorf("first argument must be a result")
	}

	category := types.NewXText(result.Category)

	for _, textCategory := range categories {
		if category.Equals(textCategory) {
			return NewTrueResult(category)
		}
	}

	return FalseResult
}

// HasIntent tests whether any intent in a classification result has `name` and minimum `confidence`
func HasIntent(env envs.Environment, result *types.XObject, name *types.XText, confidence *types.XNumber) types.XValue {
	return hasIntent(result, name, confidence, false)
}

// HasTopIntent tests whether the top intent in a classification result has `name` and minimum `confidence`
func HasTopIntent(env envs.Environment, result *types.XObject, name *types.XText, confidence *types.XNumber) types.XValue {
	return hasIntent(result, name, confidence, true)
}
