package types_test

import (
	"testing"

	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/test"
	"github.com/stretchr/testify/assert"
)

func TestXObject(t *testing.T) {
	env := envs.NewBuilder().Build()

	dep1 := types.NewXText("old")
	dep1.SetDeprecated("don't use this")

	object := types.NewXObject(map[string]types.XValue{
		"foo": types.NewXText("abc"),
		"bar": types.NewXNumberFromInt(123),
		"zed": types.XBooleanFalse,
		"dep": dep1,
		"xxx": nil,
	})
	assert.Equal(t, 5, object.Count())
	assert.ElementsMatch(t, []string{"foo", "bar", "zed", "dep", "xxx"}, object.Properties())

	val, exists := object.Get("foo")
	assert.True(t, exists)
	assert.Equal(t, types.NewXText("abc"), val)

	val, exists = object.Get("doh")
	assert.False(t, exists)
	assert.Nil(t, val)

	assert.Equal(t, `{bar: 123, dep: old, foo: abc, xxx: , zed: false}`, object.Render())
	assert.Equal(t, "bar: 123\ndep: old\nfoo: abc\nxxx: \nzed: false", object.Format(env))
	assert.Equal(t, `XObject{bar: XNumber(123), dep: XText("old"), foo: XText("abc"), xxx: nil, zed: XBoolean(false)}`, object.String())
	assert.Equal(t, "object", object.Describe())

	// test marshaling to JSON
	asJSON, _ := types.ToXJSON(object)
	assert.Equal(t, types.NewXText(`{"bar":123,"dep":"old","foo":"abc","xxx":null,"zed":false}`), asJSON)

	// if there is no explicit default, it's never included, and we can exclude the deprecated property
	object.SetMarshalOptions(true, false)
	asJSON, _ = types.ToXJSON(object)
	assert.Equal(t, types.NewXText(`{"bar":123,"foo":"abc","xxx":null,"zed":false}`), asJSON)

	// test equality
	test.AssertXEqual(t, object, types.NewXObject(map[string]types.XValue{
		"foo": types.NewXText("abc"),
		"bar": types.NewXNumberFromInt(123),
		"zed": types.XBooleanFalse,
		"dep": types.NewXText("old"),
		"xxx": nil,
	}))
}

func TestReadXObject(t *testing.T) {
	_, err := types.ReadXObject(nil)
	assert.EqualError(t, err, "JSON doesn't contain an object")
	_, err = types.ReadXObject([]byte(`null`))
	assert.EqualError(t, err, "JSON doesn't contain an object")
	_, err = types.ReadXObject([]byte(`[]`))
	assert.EqualError(t, err, "JSON doesn't contain an object")
	_, err = types.ReadXObject([]byte(`{`))
	assert.EqualError(t, err, "invalid JSON")

	obj, err := types.ReadXObject([]byte(`{}`))
	assert.NoError(t, err)
	test.AssertXEqual(t, obj, types.NewXObject(map[string]types.XValue{}))

	obj, err = types.ReadXObject([]byte(`{"foo": "abc", "bar": 123}`))
	assert.NoError(t, err)
	test.AssertXEqual(t, obj, types.NewXObject(map[string]types.XValue{
		"foo": types.NewXText("abc"),
		"bar": types.NewXNumberFromInt(123),
	}))
}

func TestXObjectWithDefault(t *testing.T) {
	env := envs.NewBuilder().Build()

	object := types.NewXObject(map[string]types.XValue{
		"__default__": types.NewXText("abc-123"),
		"foo":         types.NewXText("abc"),
		"bar":         types.NewXNumberFromInt(123),
		"zed":         types.XBooleanFalse,
	})
	assert.Equal(t, 3, object.Count())
	assert.ElementsMatch(t, []string{"foo", "bar", "zed"}, object.Properties())

	val := object.Default()
	assert.Equal(t, types.NewXText("abc-123"), val)

	// can't access default like regular property
	val, exists := object.Get("__default__")
	assert.False(t, exists)
	assert.Nil(t, val)

	assert.Equal(t, `abc-123`, object.Render()) // because of default
	assert.Equal(t, "abc-123", object.Format(env))
	assert.Equal(t, `XObject{__default__: XText("abc-123"), bar: XNumber(123), foo: XText("abc"), zed: XBoolean(false)}`, object.String())
	assert.Equal(t, "object", object.Describe())

	asJSON, _ := types.ToXJSON(object)
	assert.Equal(t, types.NewXText(`{"bar":123,"foo":"abc","zed":false}`), asJSON)

	object.SetMarshalOptions(true, false)
	asJSON, _ = types.ToXJSON(object)
	assert.Equal(t, types.NewXText(`{"__default__":"abc-123","bar":123,"foo":"abc","zed":false}`), asJSON)

	// test equality
	test.AssertXEqual(t, object, types.NewXObject(map[string]types.XValue{
		"__default__": types.NewXText("abc-123"),
		"foo":         types.NewXText("abc"),
		"bar":         types.NewXNumberFromInt(123),
		"zed":         types.XBooleanFalse,
	}))
}

func TestXLazyObject(t *testing.T) {
	env := envs.NewBuilder().Build()
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
	assert.ElementsMatch(t, []string{"foo", "bar", "zed"}, object.Properties())
	assert.Equal(t, `{bar: 123, foo: abc, zed: false}`, object.Render())
	assert.Equal(t, "bar: 123\nfoo: abc\nzed: false", object.Format(env))
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

	env := envs.NewBuilder().Build()

	for _, tc := range tests {
		object, err := types.ToXObject(env, tc.value)

		if tc.hasError {
			assert.Error(t, err.Native(), "expected error for input %s", tc.value)
		} else {
			assert.NoError(t, err.Native(), "unexpected error for input %s", tc.value)
			test.AssertXEqual(t, tc.asObject, object, "object mismatch for input %s", tc.value)
		}
	}
}
