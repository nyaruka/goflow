package utils_test

import (
	"regexp"
	"testing"

	"github.com/nyaruka/goflow/utils"

	"github.com/stretchr/testify/assert"
)

var uuid4Regex = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)

func assertIsUUID4(t *testing.T, value utils.UUID) {
	assert.Regexp(t, uuid4Regex, value, "Value %s is not a valid UUID v4", value)
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
