package tests

import (
	"fmt"
	"strconv"

	"github.com/nyaruka/goflow/excellent/types"
)

// XTestResult encapsulates not only if the test was true but what the match was
type XTestResult struct {
	matched bool
	match   interface{}
}

// Matched returns whether the test matched
func (t XTestResult) Matched() bool { return t.matched }

// Match returns the item which was matched
func (t XTestResult) Match() interface{} { return t.match }

// Resolve resolves the given key when this result is referenced in an expression
func (t XTestResult) Resolve(key string) interface{} {
	switch key {
	case "matched":
		return t.matched
	case "match":
		return t.match
	}
	return fmt.Errorf("no such key '%s' on test result", key)
}

// Atomize is called when this object needs to be reduced to a primitive
func (t XTestResult) Atomize() interface{} {
	return strconv.FormatBool(t.matched)
}

// XFalseResult can be used as a singleton for false result values
var XFalseResult = XTestResult{}

var _ types.Atomizable = XTestResult{}
var _ types.Resolvable = XTestResult{}
