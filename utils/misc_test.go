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

func TestSortedKeys(t *testing.T) {
	assert.Equal(t, []string{}, utils.SortedKeys(map[string]bool{}))
	assert.Equal(t, []string{"a", "x", "y"}, utils.SortedKeys(map[string]bool{"x": true, "y": true, "a": true}))
	assert.Equal(t, []int{3, 5, 6}, utils.SortedKeys(map[int]bool{6: true, 3: true, 5: true}))
}

func TestSet(t *testing.T) {
	assert.Equal(t, map[string]bool{}, utils.Set[string](nil))
	assert.Equal(t, map[string]bool{}, utils.Set([]string{}))
	assert.Equal(t, map[string]bool{"x": true, "y": true, "a": true}, utils.Set([]string{"a", "x", "y"}))
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
