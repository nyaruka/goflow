// Package budget provides a cost allowance for a single expression evaluation. It lives below both the
// value types and the evaluator so that either layer can charge against a budget carried on the context,
// without either having to depend on the other.
package budget

import "context"

// budgetKey is the context key under which a Budget is carried.
type budgetKey struct{}

// Budget tracks how much more work a single evaluation is allowed to do. Callers charge the cost of each
// value they produce; how cost is derived from a value (e.g. text by length, everything else by 1) is the
// caller's concern, so that this package needn't know about the value types.
type Budget struct {
	remaining int
}

// New creates a Budget with the given total cost allowance.
func New(total int) *Budget {
	return &Budget{remaining: total}
}

// Charge deducts the given cost, returning false if the budget is now exhausted.
func (b *Budget) Charge(cost int) bool {
	b.remaining -= cost
	return b.remaining >= 0
}

// With returns a copy of ctx carrying the given Budget.
func With(ctx context.Context, b *Budget) context.Context {
	return context.WithValue(ctx, budgetKey{}, b)
}

// From returns the Budget carried by ctx, or nil if there is none. A nil Budget means no accounting, so
// evaluation outside the evaluator (e.g. direct unit tests of builtins) runs unmetered.
func From(ctx context.Context) *Budget {
	b, _ := ctx.Value(budgetKey{}).(*Budget)
	return b
}
