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

func TestPythonRand(t *testing.T) {

	// test values generated using CPython 3.7.0
	rngTests := []struct {
		seed     int64
		sequence []float64
	}{
		{0, []float64{0.8444218515250481, 0.7579544029403025, 0.420571580830845, 0.25891675029296335, 0.5112747213686085}},
		{1234, []float64{0.9664535356921388, 0.4407325991753527, 0.007491470058587191, 0.9109759624491242, 0.939268997363764}},
		{-1234, []float64{0.9664535356921388, 0.4407325991753527, 0.007491470058587191, 0.9109759624491242, 0.939268997363764}},
		{12345678, []float64{0.7202671550185803, 0.6330310001166692, 0.22877255649315598, 0.25254569034434393, 0.6060686820396118}},
		{23523623435353553, []float64{0.7667782846525043, 0.2620460079415977, 0.9385746408320916, 0.10965138305592881, 0.9750957142925043}},
	}

	for _, tc := range rngTests {
		r := &utils.PythonRand{}
		r.Seed(tc.seed)

		for _, v := range tc.sequence {
			assert.Equal(t, v, r.Random())
		}
	}
}
