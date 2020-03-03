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
}
