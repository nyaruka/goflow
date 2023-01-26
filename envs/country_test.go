package envs_test

import (
	"testing"

	"github.com/nyaruka/goflow/envs"

	"github.com/stretchr/testify/assert"
)

func TestDeriveCountryFromTel(t *testing.T) {
	assert.Equal(t, envs.Country("RW"), envs.DeriveCountryFromTel("+250788383383"))
	assert.Equal(t, envs.Country("EC"), envs.DeriveCountryFromTel("+593979000000"))
	assert.Equal(t, envs.NilCountry, envs.DeriveCountryFromTel("1234"))

	v, err := envs.Country("RW").Value()
	assert.NoError(t, err)
	assert.Equal(t, "RW", v)

	v, err = envs.NilCountry.Value()
	assert.NoError(t, err)
	assert.Nil(t, v)

	var c envs.Country
	assert.NoError(t, c.Scan("RW"))
	assert.Equal(t, envs.Country("RW"), c)

	assert.NoError(t, c.Scan(nil))
	assert.Equal(t, envs.NilCountry, c)
}
