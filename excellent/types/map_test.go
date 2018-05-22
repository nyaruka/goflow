package types_test

import (
	"testing"

	"github.com/nyaruka/goflow/excellent/types"

	"github.com/stretchr/testify/assert"
)

func TestXMap(t *testing.T) {
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

	assert.Equal(t, `{"bar":"123","foo":"abc","zed":"false"}`, map1.String())

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
