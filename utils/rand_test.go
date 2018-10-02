package utils_test

import (
	"testing"

	"github.com/nyaruka/goflow/utils"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestRand(t *testing.T) {
	defer utils.SetRand(utils.DefaultRand)
	utils.SetRand(utils.NewSeededRand(1234))

	assert.Equal(t, 2, utils.RandIntN(10))
	assert.Equal(t, 5, utils.RandIntN(10))
	assert.Equal(t, 9, utils.RandIntN(10))
	assert.Equal(t, decimal.RequireFromString("0.89891152303272914281251360080204904079437255859375"), utils.RandDecimal())
	assert.Equal(t, decimal.RequireFromString("0.6087185537746531149849715802702121436595916748046875"), utils.RandDecimal())
	assert.Equal(t, decimal.RequireFromString("0.302355432890411612856240708424593321979045867919921875"), utils.RandDecimal())
}
