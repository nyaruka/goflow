package types

import "context"

// budgetKey is the context key under which an evaluation Budget is carried.
type budgetKey struct{}

// Budget tracks how much more work a single evaluation is allowed to do. Cost accrues as values are
// produced: text costs its length in bytes, everything else costs 1. Charging text by length bounds
// memory (every manufactured byte is charged once, where it's created), and the floor of 1 bounds the
// number of operations so that iteration over tiny values is bounded too.
type Budget struct {
	remaining int
}

// NewBudget creates a Budget with the given total cost allowance.
func NewBudget(total int) *Budget {
	return &Budget{remaining: total}
}

// Charge deducts the cost of the given value, returning an error if the budget is exhausted.
func (b *Budget) Charge(v XValue) *XError {
	cost := 1
	if t, ok := v.(*XText); ok {
		cost = len(t.Native())
	}

	b.remaining -= cost
	if b.remaining < 0 {
		return NewXErrorf("expression is too complex to evaluate")
	}
	return nil
}

// WithBudget returns a copy of ctx carrying the given Budget.
func WithBudget(ctx context.Context, b *Budget) context.Context {
	return context.WithValue(ctx, budgetKey{}, b)
}

// BudgetFrom returns the Budget carried by ctx, or nil if there is none. A nil Budget charges nothing,
// so evaluation outside the evaluator (e.g. direct unit tests of builtins) runs unmetered.
func BudgetFrom(ctx context.Context) *Budget {
	b, _ := ctx.Value(budgetKey{}).(*Budget)
	return b
}
