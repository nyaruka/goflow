package utils

import (
	"math/rand"
	"time"

	"github.com/shopspring/decimal"
)

// DefaultRand is the default rand for calls to Rand()
var DefaultRand = rand.New(rand.NewSource(time.Now().UnixNano()))
var currentRand = DefaultRand

// NewSeededRand creates a new seeded rand
func NewSeededRand(seed int64) *rand.Rand {
	return rand.New(rand.NewSource(seed))
}

// RandDecimal returns a random decimal in the range [0.0, 1.0)
func RandDecimal() decimal.Decimal {
	return decimal.NewFromFloat(currentRand.Float64())
}

// RandIntN returns a random integer in the range [0, n)
func RandIntN(n int) int {
	return currentRand.Intn(n)
}

// SetRand sets the rand used by Rand()
func SetRand(rnd *rand.Rand) {
	currentRand = rnd
}
