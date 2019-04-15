package types_test

import (
	"testing"

	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/utils"

	"github.com/stretchr/testify/assert"
)

func TestXObject(t *testing.T) {
	env := utils.NewEnvironmentBuilder().Build()

	object := types.NewXObject(map[string]types.XValue{
		"foo": types.NewXText("abc"),
		"bar": types.NewXNumberFromInt(123),
		"zed": types.XBooleanFalse,
		"xxx": nil,
	})
	assert.Equal(t, 4, object.Count())
	assert.ElementsMatch(t, []string{"foo", "bar", "zed", "xxx"}, object.Keys())

	val, exists := object.Get("foo")
	assert.True(t, exists)
	assert.Equal(t, types.NewXText("abc"), val)

	val, exists = object.Get("doh")
	assert.False(t, exists)
	assert.Nil(t, val)

	assert.Equal(t, `{bar: 123, foo: abc, xxx: , zed: false}`, object.Render())
	assert.Equal(t, `{bar: 123, foo: abc, xxx: , zed: false}`, object.Format(env))
	assert.Equal(t, `XObject{bar: XNumber(123), foo: XText("abc"), xxx: nil, zed: XBoolean(false)}`, object.String())
	assert.Equal(t, "object", object.Describe())

	asJSON, _ := types.ToXJSON(object)
	assert.Equal(t, types.NewXText(`{"bar":123,"foo":"abc","xxx":null,"zed":false}`), asJSON)

	// test equality
	assert.Equal(t, object, types.NewXObject(map[string]types.XValue{
		"foo": types.NewXText("abc"),
		"bar": types.NewXNumberFromInt(123),
		"zed": types.XBooleanFalse,
		"xxx": nil,
	}))
	assert.NotEqual(t, object, types.NewXObject(map[string]types.XValue{
		"bar": types.NewXNumberFromInt(123),
		"zed": types.XBooleanFalse,
		"xxx": nil,
	}))
}

func TestXLazyObject(t *testing.T) {
	env := utils.NewEnvironmentBuilder().Build()
	initialized := false

	object := types.NewXLazyObject(func() map[string]types.XValue {
		initialized = true

		return map[string]types.XValue{
			"foo": types.NewXText("abc"),
			"bar": types.NewXNumberFromInt(123),
			"zed": types.XBooleanFalse,
		}
	})

	assert.False(t, initialized)

	assert.Equal(t, 3, object.Count())
	assert.ElementsMatch(t, []string{"foo", "bar", "zed"}, object.Keys())
	assert.Equal(t, `{bar: 123, foo: abc, zed: false}`, object.Render())
	assert.Equal(t, `{bar: 123, foo: abc, zed: false}`, object.Format(env))
	assert.Equal(t, "object", object.Describe())

	assert.True(t, initialized)

	asJSON, _ := types.ToXJSON(object)
	assert.Equal(t, types.NewXText(`{"bar":123,"foo":"abc","zed":false}`), asJSON)
}

func TestToXObject(t *testing.T) {
	var tests = []struct {
		value    types.XValue
		asObject *types.XObject
		hasError bool
	}{
		{nil, types.XObjectEmpty, false},
		{types.NewXErrorf("Error"), types.XObjectEmpty, true},
		{types.NewXNumberFromInt(123), types.XObjectEmpty, true},
		{types.NewXText(""), types.XObjectEmpty, true},
		{types.NewXObject(map[string]types.XValue{"foo": types.NewXText("bar")}), types.NewXObject(map[string]types.XValue{"foo": types.NewXText("bar")}), false},
	}

	env := utils.NewEnvironmentBuilder().Build()

	for _, test := range tests {
		object, err := types.ToXObject(env, test.value)

		if test.hasError {
			assert.Error(t, err, "expected error for input %T{%s}", test.value, test.value)
		} else {
			assert.NoError(t, err, "unexpected error for input %T{%s}", test.value, test.value)
			assert.Equal(t, test.asObject, object, "object mismatch for input %T{%s}", test.value, test.value)
		}
	}
}
