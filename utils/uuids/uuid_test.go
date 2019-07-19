package uuids_test

import (
	"testing"

	"github.com/nyaruka/goflow/utils/uuids"

	"github.com/stretchr/testify/assert"
)

func TestIsV4(t *testing.T) {
	assert.False(t, uuids.IsV4(""))
	assert.True(t, uuids.IsV4("182faeb1-eb29-41e5-b288-c1af671ee671"))
	assert.False(t, uuids.IsV4("182faeb1-eb29-41e5-b288-c1af671ee67x"))
	assert.False(t, uuids.IsV4("182faeb1-eb29-41e5-b288-c1af671ee67"))
	assert.False(t, uuids.IsV4("182faeb1-eb29-41e5-b288-c1af671ee6712"))
}

func TestNew(t *testing.T) {
	uuid1 := uuids.New()
	uuid2 := uuids.New()

	assert.True(t, uuids.IsV4(string(uuid1)))
	assert.True(t, uuids.IsV4(string(uuid2)))
	assert.NotEqual(t, uuid1, uuid2)
}

func TestSeededGenerator(t *testing.T) {
	defer uuids.SetGenerator(uuids.DefaultGenerator)

	uuids.SetGenerator(uuids.NewSeededGenerator(123456))

	uuid1 := uuids.New()
	uuid2 := uuids.New()
	uuid3 := uuids.New()

	assert.True(t, uuids.IsV4(string(uuid1)))
	assert.True(t, uuids.IsV4(string(uuid2)))
	assert.True(t, uuids.IsV4(string(uuid3)))

	assert.Equal(t, uuids.UUID(`d2f852ec-7b4e-457f-ae7f-f8b243c49ff5`), uuid1)
	assert.Equal(t, uuids.UUID(`692926ea-09d6-4942-bd38-d266ec8d3716`), uuid2)
	assert.Equal(t, uuids.UUID(`8720f157-ca1c-432f-9c0b-2014ddc77094`), uuid3)

	uuids.SetGenerator(uuids.NewSeededGenerator(123456))

	// should get same sequence again for same seed
	assert.Equal(t, uuids.UUID(`d2f852ec-7b4e-457f-ae7f-f8b243c49ff5`), uuids.New())
	assert.Equal(t, uuids.UUID(`692926ea-09d6-4942-bd38-d266ec8d3716`), uuids.New())
	assert.Equal(t, uuids.UUID(`8720f157-ca1c-432f-9c0b-2014ddc77094`), uuids.New())
}
