package random

import (
	"math/rand"
	"time"

	"github.com/shopspring/decimal"
)

// DefaultGenerator is the default generator for calls to Rand()
var DefaultGenerator = rand.New(rand.NewSource(time.Now().UnixNano()))
var currentGenerator = DefaultGenerator

// NewSeededGenerator creates a new seeded generator
func NewSeededGenerator(seed int64) *rand.Rand {
	return rand.New(rand.NewSource(seed))
}

// SetGenerator sets the rand used by Rand()
func SetGenerator(rnd *rand.Rand) {
	currentGenerator = rnd
}

// Decimal returns a random decimal in the range [0.0, 1.0)
func Decimal() decimal.Decimal {
	return decimal.NewFromFloat(currentGenerator.Float64())
}

// IntN returns a random integer in the range [0, n)
func IntN(n int) int {
	return currentGenerator.Intn(n)
}
