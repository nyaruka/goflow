package flows_test

import (
	"testing"
	"time"

	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/test"

	"github.com/stretchr/testify/assert"
)

func TestResults(t *testing.T) {
	env := envs.NewBuilder().Build()

	result1 := flows.NewResult("Beer", "skol!", "Skol", "", flows.NodeUUID("26493ebb-a254-4461-a28d-c7761784e276"), "", nil, time.Date(2019, 4, 5, 14, 16, 30, 123456, time.UTC))
	result2 := flows.NewResult("Empty", "", "", "", flows.NodeUUID("26493ebb-a254-4461-a28d-c7761784e276"), "", nil, time.Date(2019, 4, 5, 14, 16, 30, 123456, time.UTC))

	results := flows.NewResults()
	results.Save(result1)
	results.Save(result2)

	assert.Equal(t, result1, results.Get("beer"))
	assert.Equal(t, result2, results.Get("empty"))
	assert.Nil(t, results.Get("xxx"))

	resultsAsContext := flows.Context(env, results)

	test.AssertXEqual(t, types.NewXObject(map[string]types.XValue{
		"__default__": types.NewXText("Beer: skol!\nEmpty: "),
		"beer": types.NewXObject(map[string]types.XValue{
			"__default__":          types.NewXText("skol!"),
			"category":             types.NewXText("Skol"),
			"categories":           types.NewXArray(types.NewXText("Skol")),
			"category_localized":   types.NewXText("Skol"),
			"categories_localized": types.NewXArray(types.NewXText("Skol")),
			"created_on":           types.NewXDateTime(time.Date(2019, 4, 5, 14, 16, 30, 123456, time.UTC)),
			"extra":                nil,
			"input":                types.XTextEmpty,
			"name":                 types.NewXText("Beer"),
			"node_uuid":            types.NewXText("26493ebb-a254-4461-a28d-c7761784e276"),
			"value":                types.NewXText("skol!"),
			"values":               types.NewXArray(types.NewXText("skol!")),
		}),
		"empty": types.NewXObject(map[string]types.XValue{
			"__default__":          types.NewXText(""),
			"category":             types.NewXText(""),
			"categories":           types.NewXArray(types.NewXText("")),
			"category_localized":   types.NewXText(""),
			"categories_localized": types.NewXArray(types.NewXText("")),
			"created_on":           types.NewXDateTime(time.Date(2019, 4, 5, 14, 16, 30, 123456, time.UTC)),
			"extra":                nil,
			"input":                types.XTextEmpty,
			"name":                 types.NewXText("Empty"),
			"node_uuid":            types.NewXText("26493ebb-a254-4461-a28d-c7761784e276"),
			"value":                types.NewXText(""),
			"values":               types.NewXArray(types.NewXText("")),
		}),
	}), resultsAsContext)
}
