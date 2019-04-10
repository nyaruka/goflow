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
		"zed": types.XBooleanFalse,
		"xxx": nil,
	})
	assert.Equal(t, 4, dict.Length())
	assert.ElementsMatch(t, []string{"foo", "bar", "zed", "xxx"}, dict.Keys())

	val, exists := dict.Get("foo")
	assert.True(t, exists)
	assert.Equal(t, types.NewXText("abc"), val)

	val, exists = dict.Get("doh")
	assert.False(t, exists)
	assert.Nil(t, val)

	assert.Equal(t, types.NewXText("{bar: 123, foo: abc, xxx: , zed: false}"), dict.ToXText(env))
	assert.Equal(t, `XDict{bar: XNumber(123), foo: XText("abc"), xxx: nil, zed: XBoolean(false)}`, dict.String())
	assert.Equal(t, "dict", dict.Describe())

	asJSON, _ := types.ToXJSON(dict)
	assert.Equal(t, types.NewXText(`{"bar":123,"foo":"abc","xxx":null,"zed":false}`), asJSON)

	// test equality
	assert.Equal(t, dict, types.NewXDict(map[string]types.XValue{
		"foo": types.NewXText("abc"),
		"bar": types.NewXNumberFromInt(123),
		"zed": types.XBooleanFalse,
		"xxx": nil,
	}))
	assert.NotEqual(t, dict, types.NewXDict(map[string]types.XValue{
		"bar": types.NewXNumberFromInt(123),
		"zed": types.XBooleanFalse,
		"xxx": nil,
	}))
}

func TestXLazyDict(t *testing.T) {
	env := utils.NewEnvironmentBuilder().Build()

	dict := types.NewXLazyDict(func() map[string]types.XValue {
		return map[string]types.XValue{
			"foo": types.NewXText("abc"),
			"bar": types.NewXNumberFromInt(123),
			"zed": types.XBooleanFalse,
		}
	})

	assert.Equal(t, 3, dict.Length())
	assert.ElementsMatch(t, []string{"foo", "bar", "zed"}, dict.Keys())
	assert.Equal(t, types.NewXText("{bar: 123, foo: abc, zed: false}"), dict.ToXText(env))
	assert.Equal(t, "dict", dict.Describe())

	asJSON, _ := types.ToXJSON(dict)
	assert.Equal(t, types.NewXText(`{"bar":123,"foo":"abc","zed":false}`), asJSON)
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
