package tests

import (
	"github.com/nyaruka/goflow/excellent/types"
)

// XTestResult encapsulates not only if the test was true but what the match was
type XTestResult struct {
	matched bool
	match   types.XValue
}

// Matched returns whether the test matched
func (t XTestResult) Matched() bool { return t.matched }

// Match returns the item which was matched
func (t XTestResult) Match() types.XValue { return t.match }

// Resolve resolves the given key when this result is referenced in an expression
func (t XTestResult) Resolve(key string) types.XValue {
	switch key {
	case "matched":
		return types.NewXBool(t.matched)
	case "match":
		return t.match
	}
	return types.NewXResolveError(t, key)
}

// Reduce is called when this object needs to be reduced to a primitive
func (t XTestResult) Reduce() types.XPrimitive {
	return types.NewXBool(t.matched)
}

// ToXJSON is called when this type is passed to @(json(...))
func (t XTestResult) ToXJSON() types.XText {
	return types.ResolveKeys(t, "matched", "match").ToXJSON()
}

// XFalseResult can be used as a singleton for false result values
var XFalseResult = XTestResult{}

var _ types.XValue = XTestResult{}
var _ types.XResolvable = XTestResult{}
