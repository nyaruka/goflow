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

func TestMaxInt(t *testing.T) {
	assert.Equal(t, 1, utils.MaxInt(0, 1))
	assert.Equal(t, 1, utils.MaxInt(1, 0))
	assert.Equal(t, 1, utils.MaxInt(1, -1))
}

func TestMinInt(t *testing.T) {
	assert.Equal(t, 0, utils.MinInt(0, 1))
	assert.Equal(t, 0, utils.MinInt(1, 0))
	assert.Equal(t, -1, utils.MinInt(1, -1))
}

func TestFindPhoneNumber(t *testing.T) {
	assert.Equal(t, "", utils.FindPhoneNumber("", ""))
	assert.Equal(t, "", utils.FindPhoneNumber("", "RW"))

	assert.Equal(t, "+250788383383", utils.FindPhoneNumber("+250788383383", ""))

	assert.Equal(t, "+250788383383", utils.FindPhoneNumber("Hi my phone is +250788383383", "RW"))
	assert.Equal(t, "+250788383383", utils.FindPhoneNumber("Hi my phone is +250788383383", ""))
	assert.Equal(t, "+250788383383", utils.FindPhoneNumber("Hi my phone is 0788383383", "RW"))
	assert.Equal(t, "", utils.FindPhoneNumber("Hi my phone is 0788383383", ""))

	assert.Equal(t, "+12024561111", utils.FindPhoneNumber("Hi my phone is +12024561111", "US"))
	assert.Equal(t, "+12024561111", utils.FindPhoneNumber("Hi my phone is +12024561111", ""))
	assert.Equal(t, "+12024561111", utils.FindPhoneNumber("Hi my phone is (202) 456-1111", "US"))
	assert.Equal(t, "", utils.FindPhoneNumber("Hi my phone is (202) 456-1111", ""))
}

func TestDeriveCountryFromTel(t *testing.T) {
	assert.Equal(t, "RW", utils.DeriveCountryFromTel("+250788383383"))
	assert.Equal(t, "EC", utils.DeriveCountryFromTel("+593979000000"))
	assert.Equal(t, "", utils.DeriveCountryFromTel("1234"))
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
