package tests

import (
	"strings"

	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/utils"
)

// XTestResult encapsulates not only if the test was true but what the match was
type XTestResult struct {
	match types.XValue
	extra types.XDict
}

// NewTrueResult creates a new matched result
func NewTrueResult(match types.XValue) XTestResult {
	return XTestResult{match, nil}
}

// NewTrueResultWithExtra creates a new matched result with extra info about the match
func NewTrueResultWithExtra(match types.XValue, extra types.XDict) XTestResult {
	return XTestResult{match, extra}
}

// Match returns the item which was matched
func (t XTestResult) Match() types.XValue { return t.match }

// Extra returns the extra data about the match
func (t XTestResult) Extra() types.XDict { return t.extra }

// Resolve resolves the given key when this result is referenced in an expression
func (t XTestResult) Resolve(env utils.Environment, key string) types.XValue {
	switch strings.ToLower(key) {
	case "match":
		return t.match
	}
	return types.NewXResolveError(t, key)
}

// Describe returns a representation of this type for error messages
func (t XTestResult) Describe() string { return "test result" }

// Reduce is called when this object needs to be reduced to a primitive
func (t XTestResult) Reduce(env utils.Environment) types.XPrimitive {
	return types.XBooleanTrue
}

// ToXJSON is called when this type is passed to @(json(...))
func (t XTestResult) ToXJSON(env utils.Environment) types.XText {
	return types.ResolveKeys(env, t, "match").ToXJSON(env)
}

var _ types.XValue = XTestResult{}
var _ types.XResolvable = XTestResult{}
