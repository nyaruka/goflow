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

// HasIntent is a deprecated NLU classification test that always returns a false result.
func HasIntent(env envs.Environment, result *types.XObject, name *types.XText, confidence *types.XNumber) types.XValue {
	return FalseResult
}

// HasTopIntent is a deprecated NLU classification test that always returns a false result.
func HasTopIntent(env envs.Environment, result *types.XObject, name *types.XText, confidence *types.XNumber) types.XValue {
	return FalseResult
}
