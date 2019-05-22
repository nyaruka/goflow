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

func assertIsUUID4(t *testing.T, value utils.UUID) {
	assert.Regexp(t, utils.UUID4OnlyRegex, value, "value %s is not a valid UUID v4", value)
}

func TestUUIDGenerators(t *testing.T) {
	uuid1 := utils.NewUUID()
	uuid2 := utils.NewUUID()

	assertIsUUID4(t, uuid1)
	assertIsUUID4(t, uuid2)
	assert.NotEqual(t, uuid1, uuid2)

	utils.SetUUIDGenerator(utils.NewSeededUUID4Generator(123456))
	defer utils.SetUUIDGenerator(utils.DefaultUUIDGenerator)

	uuid3 := utils.NewUUID()
	uuid4 := utils.NewUUID()

	assertIsUUID4(t, uuid3)
	assertIsUUID4(t, uuid4)
	assert.NotEqual(t, uuid3, uuid4)

	utils.SetUUIDGenerator(utils.NewSeededUUID4Generator(123456))

	uuid5 := utils.NewUUID()
	uuid6 := utils.NewUUID()

	assert.Equal(t, uuid3, uuid5)
	assert.Equal(t, uuid4, uuid6)
}
