package random_test

import (
	"testing"

	"github.com/nyaruka/goflow/utils/random"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestRand(t *testing.T) {
	defer random.SetGenerator(random.DefaultGenerator)
	random.SetGenerator(random.NewSeededGenerator(1234))

	assert.Equal(t, 2, random.IntN(10))
	assert.Equal(t, 5, random.IntN(10))
	assert.Equal(t, 9, random.IntN(10))
	assert.Equal(t, decimal.RequireFromString("0.89891152303272914281251360080204904079437255859375"), random.Decimal())
	assert.Equal(t, decimal.RequireFromString("0.6087185537746531149849715802702121436595916748046875"), random.Decimal())
	assert.Equal(t, decimal.RequireFromString("0.302355432890411612856240708424593321979045867919921875"), random.Decimal())
}
