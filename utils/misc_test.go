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

func TestMinInt(t *testing.T) {
	assert.Equal(t, 0, utils.MinInt(0, 1))
	assert.Equal(t, 0, utils.MinInt(1, 0))
	assert.Equal(t, -1, utils.MinInt(1, -1))
}
