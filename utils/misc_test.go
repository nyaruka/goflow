package utils_test

import (
	"testing"

	"github.com/nyaruka/goflow/utils"

	"github.com/stretchr/testify/assert"
)

func TestIsNil(t *testing.T) {
	assert.True(t, utils.IsNil(nil))
	assert.True(t, utils.IsNil(error(nil)))
	assert.False(t, utils.IsNil(""))
}

func TestMax(t *testing.T) {
	assert.Equal(t, 1, utils.Max(0, 1))
	assert.Equal(t, 1, utils.Max(1, 0))
	assert.Equal(t, 1, utils.Max(1, -1))

	assert.Equal(t, uint16(1), utils.Max(uint16(0), uint16(1)))
}

func TestMin(t *testing.T) {
	assert.Equal(t, 0, utils.Min(0, 1))
	assert.Equal(t, 0, utils.Min(1, 0))
	assert.Equal(t, -1, utils.Min(1, -1))

	assert.Equal(t, uint16(0), utils.Min(uint16(0), uint16(1)))
}

func TestReadTypeFromJSON(t *testing.T) {
	_, err := utils.ReadTypeFromJSON([]byte(`{}`))
	assert.EqualError(t, err, "field 'type' is required")

	_, err = utils.ReadTypeFromJSON([]byte(`{"type": ""}`))
	assert.EqualError(t, err, "field 'type' is required")

	typeName, err := utils.ReadTypeFromJSON([]byte(`{"thing": 2, "type": "foo"}`))
	assert.NoError(t, err)
	assert.Equal(t, "foo", typeName)
}
