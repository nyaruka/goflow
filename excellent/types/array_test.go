package types_test

import (
	"testing"

	"github.com/nyaruka/goflow/excellent/types"

	"github.com/stretchr/testify/assert"
)

func TestXArray(t *testing.T) {
	arr1 := types.NewXArray(types.NewXText("abc"), types.NewXNumberFromInt(123))
	assert.Equal(t, 2, arr1.Length())

	arr1.Append(types.XBooleanFalse)
	assert.Equal(t, 3, arr1.Length())
	assert.Equal(t, types.NewXNumberFromInt(123), arr1.Index(1))

	assert.Equal(t, `["abc","123","false"]`, arr1.String())

	// test equality
	assert.Equal(t, types.NewXArray(types.NewXText("abc"), types.NewXNumberFromInt(123)), types.NewXArray(types.NewXText("abc"), types.NewXNumberFromInt(123)))
	assert.NotEqual(t, types.NewXArray(types.NewXText("abc")), types.NewXArray(types.NewXText("abc"), types.NewXNumberFromInt(123)))
}
