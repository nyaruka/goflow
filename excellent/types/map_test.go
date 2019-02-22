package types_test

import (
	"testing"

	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/utils"

	"github.com/stretchr/testify/assert"
)

func TestXMap(t *testing.T) {
	env := utils.NewEnvironmentBuilder().Build()

	map1 := types.NewXMap(map[string]types.XValue{
		"foo": types.NewXText("abc"),
		"bar": types.NewXNumberFromInt(123),
	})
	assert.Equal(t, 2, map1.Length())
	assert.ElementsMatch(t, []string{"foo", "bar"}, map1.Keys())

	map1.Put("zed", types.XBooleanFalse)
	assert.Equal(t, 3, map1.Length())
	assert.Equal(t, types.NewXNumberFromInt(123), map1.Resolve(nil, "bar"))
	assert.True(t, types.IsXError(map1.Resolve(nil, "xxxx")))

	assert.Equal(t, types.NewXText("bar: 123\nfoo: abc\nzed: false"), map1.ToXText(env))
	assert.Equal(t, types.NewXText(`{"bar":123,"foo":"abc","zed":false}`), map1.ToXJSON(env))
	assert.Equal(t, "bar: 123\nfoo: abc\nzed: false", map1.String())
	assert.Equal(t, map1, map1.Reduce(utils.NewEnvironmentBuilder().Build()))
	assert.Equal(t, "map", map1.Describe())

	// test equality
	assert.Equal(t, map1, types.NewXMap(map[string]types.XValue{
		"foo": types.NewXText("abc"),
		"bar": types.NewXNumberFromInt(123),
		"zed": types.XBooleanFalse,
	}))
	assert.NotEqual(t, map1, types.NewXMap(map[string]types.XValue{
		"bar": types.NewXNumberFromInt(123),
		"zed": types.XBooleanFalse,
	}))
}
