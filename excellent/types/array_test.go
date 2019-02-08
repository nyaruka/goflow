package types_test

import (
	"testing"

	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/utils"

	"github.com/stretchr/testify/assert"
)

func TestXArray(t *testing.T) {
	env := utils.NewEnvironmentBuilder().Environment()

	arr1 := types.NewXArray(types.NewXText("abc"), types.NewXNumberFromInt(123))
	assert.Equal(t, 2, arr1.Length())

	arr1.Append(types.XBooleanFalse)
	assert.Equal(t, 3, arr1.Length())
	assert.Equal(t, types.NewXNumberFromInt(123), arr1.Index(1))

	assert.Equal(t, types.NewXText(`["abc",123,false]`), arr1.ToXJSON(env))
	assert.Equal(t, types.NewXText(`["abc",123,false]`), arr1.ToXText(env))
	assert.Equal(t, `["abc",123,false]`, arr1.String())
	assert.Equal(t, arr1, arr1.Reduce(utils.NewEnvironmentBuilder().Environment()))
	assert.Equal(t, "array", arr1.Describe())

	// test equality
	assert.Equal(t, types.NewXArray(types.NewXText("abc"), types.NewXNumberFromInt(123)), types.NewXArray(types.NewXText("abc"), types.NewXNumberFromInt(123)))
	assert.NotEqual(t, types.NewXArray(types.NewXText("abc")), types.NewXArray(types.NewXText("abc"), types.NewXNumberFromInt(123)))
}
