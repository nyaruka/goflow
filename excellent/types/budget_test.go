package types_test

import (
	"context"
	"testing"

	"github.com/nyaruka/goflow/excellent/types"

	"github.com/stretchr/testify/assert"
)

func TestBudget(t *testing.T) {
	b := types.NewBudget(100)

	// text costs its length in bytes
	assert.Nil(t, b.Charge(types.NewXText("hello"))) // 5, remaining 95
	assert.Nil(t, b.Charge(types.NewXText("world"))) // 5, remaining 90

	// non-text costs 1
	assert.Nil(t, b.Charge(types.NewXNumberFromInt(123456789))) // 1, remaining 89
	assert.Nil(t, b.Charge(types.XBooleanTrue))                 // 1, remaining 88
	assert.Nil(t, b.Charge(nil))                                // 1, remaining 87

	// exhausting the budget returns an error
	assert.Nil(t, b.Charge(types.NewXText(str(87)))) // remaining 0
	xerr := b.Charge(types.NewXText("x"))            // over
	assert.NotNil(t, xerr)
	assert.EqualError(t, xerr, "expression is too complex to evaluate")

	// once exhausted, stays exhausted
	assert.NotNil(t, b.Charge(types.XBooleanTrue))
}

func TestBudgetFrom(t *testing.T) {
	// no budget on a bare context
	assert.Nil(t, types.BudgetFrom(context.Background()))

	// round-trips through the context
	b := types.NewBudget(10)
	ctx := types.WithBudget(context.Background(), b)
	assert.Same(t, b, types.BudgetFrom(ctx))
}

func str(n int) string {
	s := make([]byte, n)
	for i := range s {
		s[i] = 'x'
	}
	return string(s)
}
