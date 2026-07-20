package budget_test

import (
	"context"
	"testing"

	"github.com/nyaruka/goflow/excellent/budget"

	"github.com/stretchr/testify/assert"
)

func TestBudget(t *testing.T) {
	b := budget.New(10)

	assert.True(t, b.Charge(4))  // remaining 6
	assert.True(t, b.Charge(6))  // remaining 0
	assert.False(t, b.Charge(1)) // over

	// once exhausted, stays exhausted
	assert.False(t, b.Charge(0))
}

func TestContext(t *testing.T) {
	// no budget on a bare context
	assert.Nil(t, budget.From(context.Background()))

	// round-trips through the context
	b := budget.New(10)
	ctx := budget.With(context.Background(), b)
	assert.Same(t, b, budget.From(ctx))
}
