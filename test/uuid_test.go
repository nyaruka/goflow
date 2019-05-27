package test_test

import (
	"testing"

	"github.com/nyaruka/goflow/test"
	"github.com/nyaruka/goflow/utils"

	"github.com/stretchr/testify/assert"
)

func TestSeededUUIDGenerator(t *testing.T) {
	defer utils.SetUUIDGenerator(utils.DefaultUUIDGenerator)

	utils.SetUUIDGenerator(test.NewSeededUUIDGenerator(123456))

	uuid1 := utils.NewUUID()
	uuid2 := utils.NewUUID()
	uuid3 := utils.NewUUID()

	assert.True(t, utils.IsUUIDv4(string(uuid1)))
	assert.True(t, utils.IsUUIDv4(string(uuid2)))
	assert.True(t, utils.IsUUIDv4(string(uuid3)))

	assert.Equal(t, utils.UUID(`d2f852ec-7b4e-457f-ae7f-f8b243c49ff5`), uuid1)
	assert.Equal(t, utils.UUID(`692926ea-09d6-4942-bd38-d266ec8d3716`), uuid2)
	assert.Equal(t, utils.UUID(`8720f157-ca1c-432f-9c0b-2014ddc77094`), uuid3)

	utils.SetUUIDGenerator(test.NewSeededUUIDGenerator(123456))

	// should get same sequence again for same seed
	assert.Equal(t, utils.UUID(`d2f852ec-7b4e-457f-ae7f-f8b243c49ff5`), utils.NewUUID())
	assert.Equal(t, utils.UUID(`692926ea-09d6-4942-bd38-d266ec8d3716`), utils.NewUUID())
	assert.Equal(t, utils.UUID(`8720f157-ca1c-432f-9c0b-2014ddc77094`), utils.NewUUID())
}
