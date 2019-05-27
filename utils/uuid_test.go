package utils_test

import (
	"testing"

	"github.com/nyaruka/goflow/utils"

	"github.com/stretchr/testify/assert"
)

func TestIsUUIDv4(t *testing.T) {
	assert.False(t, utils.IsUUIDv4(""))
	assert.True(t, utils.IsUUIDv4("182faeb1-eb29-41e5-b288-c1af671ee671"))
	assert.False(t, utils.IsUUIDv4("182faeb1-eb29-41e5-b288-c1af671ee67x"))
	assert.False(t, utils.IsUUIDv4("182faeb1-eb29-41e5-b288-c1af671ee67"))
	assert.False(t, utils.IsUUIDv4("182faeb1-eb29-41e5-b288-c1af671ee6712"))
}

func TestNewUUID(t *testing.T) {
	uuid1 := utils.NewUUID()
	uuid2 := utils.NewUUID()

	assert.True(t, utils.IsUUIDv4(string(uuid1)))
	assert.True(t, utils.IsUUIDv4(string(uuid2)))
	assert.NotEqual(t, uuid1, uuid2)
}
