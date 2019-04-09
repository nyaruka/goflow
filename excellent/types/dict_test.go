package types_test

import (
	"testing"

	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/utils"

	"github.com/stretchr/testify/assert"
)

func TestXDict(t *testing.T) {
	env := utils.NewEnvironmentBuilder().Build()

	dict := types.NewXDict(map[string]types.XValue{
		"foo": types.NewXText("abc"),
		"bar": types.NewXNumberFromInt(123),
	})
	assert.Equal(t, 2, dict.Length())
	assert.ElementsMatch(t, []string{"foo", "bar"}, dict.Keys())

	dict.Put("zed", types.XBooleanFalse)

	assert.Equal(t, 3, dict.Length())
	assert.Equal(t, types.XBooleanFalse, dict.Get("zed"))
	assert.Equal(t, types.NewXNumberFromInt(123), dict.Resolve(nil, "bar"))
	assert.True(t, types.IsXError(dict.Resolve(nil, "xxxx")))

	assert.Equal(t, types.NewXText("{bar: 123, foo: abc, zed: false}"), dict.ToXText(env))
	assert.Equal(t, types.NewXText(`{"bar":123,"foo":"abc","zed":false}`), dict.ToXJSON(env))
	assert.Equal(t, "{bar: 123, foo: abc, zed: false}", dict.String())
	assert.Equal(t, "dict", dict.Describe())

	// test equality
	assert.Equal(t, dict, types.NewXDict(map[string]types.XValue{
		"foo": types.NewXText("abc"),
		"bar": types.NewXNumberFromInt(123),
		"zed": types.XBooleanFalse,
	}))
	assert.NotEqual(t, dict, types.NewXDict(map[string]types.XValue{
		"bar": types.NewXNumberFromInt(123),
		"zed": types.XBooleanFalse,
	}))
}

func TestToXDict(t *testing.T) {
	var tests = []struct {
		value    types.XValue
		asDict   *types.XDict
		hasError bool
	}{
		{nil, types.XDictEmpty, false},
		{types.NewXErrorf("Error"), types.XDictEmpty, true},
		{types.NewXNumberFromInt(123), types.XDictEmpty, true},
		{types.NewXText(""), types.XDictEmpty, true},
		{types.NewXDict(map[string]types.XValue{"foo": types.NewXText("bar")}), types.NewXDict(map[string]types.XValue{"foo": types.NewXText("bar")}), false},
	}

	env := utils.NewEnvironmentBuilder().Build()

	for _, test := range tests {
		dict, err := types.ToXDict(env, test.value)

		if test.hasError {
			assert.Error(t, err, "expected error for input %T{%s}", test.value, test.value)
		} else {
			assert.NoError(t, err, "unexpected error for input %T{%s}", test.value, test.value)
			assert.Equal(t, test.asDict, dict, "dict mismatch for input %T{%s}", test.value, test.value)
		}
	}
}
