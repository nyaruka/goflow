package tools_test

import (
	"testing"

	"github.com/nyaruka/goflow/excellent/tools"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/test"
	"github.com/stretchr/testify/assert"
)

func TestContextWalk(t *testing.T) {
	context := types.NewXObject(map[string]types.XValue{
		"__default__": types.NewXText("Bob"),
		"foo": types.NewXArray(
			types.NewXObject(map[string]types.XValue{
				"__default__": types.NewXText("Bob"),
				"bar":         types.NewXNumberFromInt(123),
			}),
			types.NewXNumberFromInt(256),
		),
		"bar": types.NewXNumberFromInt(256),
		"zed": types.NewXObject(map[string]types.XValue{
			"bar": types.NewXNumberFromInt(345),
		}),
		"nil": (*types.XObject)(nil), // non-nil interface to a nil struct
	})

	// test finding all non-nil values
	count := 0
	tools.ContextWalk(context, func(v types.XValue) {
		count++
	})
	assert.Equal(t, 8, count)

	// test finding just objects
	actual := make([]types.XValue, 0)
	tools.ContextWalkObjects(context, func(o *types.XObject) {
		actual = append(actual, o)
	})

	expected := types.NewXArray(
		context,
		types.NewXObject(map[string]types.XValue{
			"__default__": types.NewXText("Bob"),
			"bar":         types.NewXNumberFromInt(123),
		}),
		types.NewXObject(map[string]types.XValue{
			"bar": types.NewXNumberFromInt(345),
		}),
	)

	test.AssertXEqual(t, expected, types.NewXArray(actual...), "found objects mismatch")
}
