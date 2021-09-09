package types_test

import (
	"testing"

	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/stretchr/testify/assert"
)

func TestXArray(t *testing.T) {
	env := envs.NewBuilder().Build()

	arr1 := types.NewXArray(types.NewXText("abc"), types.NewXNumberFromInt(123), types.XBooleanFalse)
	assert.Equal(t, 3, arr1.Count())
	assert.Equal(t, types.NewXText("abc"), arr1.Get(0))
	assert.Equal(t, types.NewXNumberFromInt(123), arr1.Get(1))

	assert.Equal(t, `[abc, 123, false]`, arr1.Render())
	assert.Equal(t, `abc, 123, false`, arr1.Format(env))
	assert.Equal(t, `XArray[XText("abc"), XNumber(123), XBoolean(false)]`, arr1.String())
	assert.Equal(t, "array", arr1.Describe())

	asJSON, _ := types.ToXJSON(arr1)
	assert.Equal(t, types.NewXText(`["abc",123,false]`), asJSON)

	// test equality
	assert.True(t,
		types.NewXArray(types.NewXText("abc"), types.NewXNumberFromInt(123)).Equals(
			types.NewXArray(types.NewXText("abc"), types.NewXNumberFromInt(123)),
		),
	)
	assert.False(t,
		types.NewXArray(types.NewXText("abc")).Equals(
			types.NewXArray(types.NewXText("abc"), types.NewXNumberFromInt(123)),
		),
	)

	arr2 := types.NewXArray(
		types.NewXObject(map[string]types.XValue{
			"foo": types.NewXNumberFromInt(123),
			"bar": types.XBooleanFalse,
		}),
		types.NewXNumberFromInt(123),
	)

	assert.Equal(t, "- bar: false\n  foo: 123\n- 123", arr2.Format(env))
}

func TestXLazyArray(t *testing.T) {
	env := envs.NewBuilder().Build()
	initialized := false

	arr1 := types.NewXLazyArray(func() []types.XValue {
		initialized = true

		return []types.XValue{
			types.NewXText("abc"),
			types.NewXNumberFromInt(123),
			types.XBooleanFalse,
		}
	})

	assert.False(t, initialized)

	assert.Equal(t, 3, arr1.Count())
	assert.Equal(t, types.NewXText("abc"), arr1.Get(0))
	assert.Equal(t, types.NewXNumberFromInt(123), arr1.Get(1))
	assert.Equal(t, `[abc, 123, false]`, arr1.Render())
	assert.Equal(t, `abc, 123, false`, arr1.Format(env))
	assert.Equal(t, `XArray[XText("abc"), XNumber(123), XBoolean(false)]`, arr1.String())
	assert.Equal(t, "array", arr1.Describe())

	assert.True(t, initialized)

	asJSON, _ := types.ToXJSON(arr1)
	assert.Equal(t, types.NewXText(`["abc",123,false]`), asJSON)
}

func TestToXArray(t *testing.T) {
	var tests = []struct {
		value    types.XValue
		asArray  *types.XArray
		hasError bool
	}{
		{nil, types.XArrayEmpty, false},
		{types.NewXErrorf("Error"), types.XArrayEmpty, true},
		{types.NewXNumberFromInt(123), types.XArrayEmpty, true},
		{types.NewXText(""), types.XArrayEmpty, true},
		{types.NewXArray(types.NewXText("foo"), types.NewXText("bar")), types.NewXArray(types.NewXText("foo"), types.NewXText("bar")), false},
	}

	env := envs.NewBuilder().Build()

	for _, test := range tests {
		array, err := types.ToXArray(env, test.value)

		if test.hasError {
			assert.Error(t, err, "expected error for input %T{%s}", test.value, test.value)
		} else {
			assert.NoError(t, err, "unexpected error for input %T{%s}", test.value, test.value)
			assert.Equal(t, test.asArray, array, "array mismatch for input %T{%s}", test.value, test.value)
		}
	}
}
