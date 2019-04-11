package flows_test

import (
	"testing"
	"time"

	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/test"
	"github.com/nyaruka/goflow/utils"

	"github.com/stretchr/testify/assert"
)

func TestResults(t *testing.T) {
	env := utils.NewEnvironmentBuilder().Build()

	result := flows.NewResult("Beer", "skol!", "Skol", "", flows.NodeUUID("26493ebb-a254-4461-a28d-c7761784e276"), "", nil, time.Date(2019, 4, 5, 14, 16, 30, 123456, time.UTC))
	results := flows.NewResults()
	results.Save(result)

	assert.Equal(t, result, results.Get("beer"))
	assert.Nil(t, results.Get("xxx"))

	test.AssertXEqual(t, types.NewXDict(map[string]types.XValue{
		"beer": types.NewXDict(map[string]types.XValue{
			"category":           types.NewXText("Skol"),
			"category_localized": types.NewXText("Skol"),
			"created_on":         types.NewXDateTime(time.Date(2019, 4, 5, 14, 16, 30, 123456, time.UTC)),
			"input":              types.XTextEmpty,
			"name":               types.NewXText("Beer"),
			"node_uuid":          types.NewXText("26493ebb-a254-4461-a28d-c7761784e276"),
			"value":              types.NewXText("skol!"),
		}),
	}), flows.ContextFunc(env, results.SimpleContext))

	test.AssertXEqual(t, types.NewXDict(map[string]types.XValue{
		"beer": types.NewXDict(map[string]types.XValue{
			"categories":           types.NewXArray(types.NewXText("Skol")),
			"categories_localized": types.NewXArray(types.NewXText("Skol")),
			"created_on":           types.NewXDateTime(time.Date(2019, 4, 5, 14, 16, 30, 123456, time.UTC)),
			"extra":                nil,
			"input":                types.XTextEmpty,
			"name":                 types.NewXText("Beer"),
			"node_uuid":            types.NewXText("26493ebb-a254-4461-a28d-c7761784e276"),
			"values":               types.NewXArray(types.NewXText("skol!")),
		}),
	}), flows.Context(env, results))
}
