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

func TestDeriveCountryFromTel(t *testing.T) {
	assert.Equal(t, "RW", utils.DeriveCountryFromTel("+250788383383"))
	assert.Equal(t, "EC", utils.DeriveCountryFromTel("+593979000000"))
	assert.Equal(t, "", utils.DeriveCountryFromTel("1234"))
}

func TestVersionCompare(t *testing.T) {
	_, err := utils.VersionCompare("x", "12.0")
	assert.EqualError(t, err, "Malformed version: x")

	_, err = utils.VersionCompare("12.0", "x")
	assert.EqualError(t, err, "Malformed version: x")

	c, err := utils.VersionCompare("12.0.0", "12.0")
	assert.NoError(t, err)
	assert.Equal(t, c, 0)

	c, err = utils.VersionCompare("13.0", "12.0")
	assert.NoError(t, err)
	assert.Equal(t, c, 1)

	c, err = utils.VersionCompare("12.0", "12.1")
	assert.NoError(t, err)
	assert.Equal(t, c, -1)
}
