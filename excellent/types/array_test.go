package types_test

import (
	"testing"

	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/utils"

	"github.com/stretchr/testify/assert"
)

func TestXArray(t *testing.T) {
	env := utils.NewEnvironmentBuilder().Build()

	arr1 := types.NewXArray(types.NewXText("abc"), types.NewXNumberFromInt(123), types.XBooleanFalse)
	assert.Equal(t, 3, arr1.Length())
	assert.Equal(t, types.NewXText("abc"), arr1.Get(0))
	assert.Equal(t, types.NewXNumberFromInt(123), arr1.Get(1))

	assert.Equal(t, types.NewXText(`[abc, 123, false]`), arr1.ToXText(env))
	assert.Equal(t, `XArray[XText("abc"), XNumber(123), XBoolean(false)]`, arr1.String())
	assert.Equal(t, "array", arr1.Describe())

	asJSON, _ := types.ToXJSON(arr1)
	assert.Equal(t, types.NewXText(`["abc",123,false]`), asJSON)

	// test equality
	assert.Equal(t, types.NewXArray(types.NewXText("abc"), types.NewXNumberFromInt(123)), types.NewXArray(types.NewXText("abc"), types.NewXNumberFromInt(123)))
	assert.NotEqual(t, types.NewXArray(types.NewXText("abc")), types.NewXArray(types.NewXText("abc"), types.NewXNumberFromInt(123)))
}

func TestXLazyArray(t *testing.T) {
	env := utils.NewEnvironmentBuilder().Build()
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

	assert.Equal(t, 3, arr1.Length())
	assert.Equal(t, types.NewXText("abc"), arr1.Get(0))
	assert.Equal(t, types.NewXNumberFromInt(123), arr1.Get(1))
	assert.Equal(t, types.NewXText(`[abc, 123, false]`), arr1.ToXText(env))
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

	env := utils.NewEnvironmentBuilder().Build()

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
