package test

import (
	"fmt"
	"testing"

	"github.com/nyaruka/goflow/excellent/types"

	"github.com/stretchr/testify/assert"
)

// AssertEqual is equivalent to assert.Equal for two XValue instances
func AssertEqual(t *testing.T, expected types.XValue, actual types.XValue, msgAndArgs ...interface{}) bool {
	if !types.Equals(expected, actual) {
		return assert.Fail(t, fmt.Sprintf("Not equal: \n"+
			"expected: %s\n"+
			"actual  : %s", expected, actual), msgAndArgs...)
	}
	return true
}
