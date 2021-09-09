package operators_test

import (
	"fmt"
	"testing"

	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/operators"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/test"
	"github.com/stretchr/testify/assert"
)

var xs = types.NewXText
var xn = types.RequireXNumberFromString
var xi = types.NewXNumberFromInt
var ERROR = types.NewXErrorf("any error")

func TestBinaryOperators(t *testing.T) {
	env := envs.NewBuilder().Build()

	testCases := []struct {
		operator operators.BinaryOperator
		arg1     types.XValue
		arg2     types.XValue
		expected types.XValue
	}{
		{operators.Concatenate, xs("hello"), xs("world"), xs("helloworld")},
		{operators.Concatenate, xs("hello"), nil, xs("hello")},
		{operators.Concatenate, xs("1"), xs("3"), xs("13")},
		{operators.Concatenate, xi(1), xi(3), xs("13")},
		{operators.Concatenate, ERROR, xi(1), ERROR},
		{operators.Concatenate, xi(1), ERROR, ERROR},

		{operators.Equal, xs("hello"), xs("hello"), types.XBooleanTrue},
		{operators.Equal, xs("hello"), xs("world"), types.XBooleanFalse},
		{operators.Equal, xn("1.5"), xn("1.5"), types.XBooleanTrue},
		{operators.Equal, ERROR, xi(1), ERROR},
		{operators.Equal, xi(1), ERROR, ERROR},

		{operators.NotEqual, xs("hello"), xs("hello"), types.XBooleanFalse},
		{operators.NotEqual, xs("hello"), xs("world"), types.XBooleanTrue},
		{operators.NotEqual, xn("1.5"), xn("1.5"), types.XBooleanFalse},
		{operators.NotEqual, ERROR, xi(1), ERROR},
		{operators.NotEqual, xi(1), ERROR, ERROR},

		{operators.Add, xi(1), xi(3), xi(4)},
		{operators.Add, xn("1.5"), xn("2.3"), xn("3.8")},
		{operators.Add, xs("1"), xs("3"), xi(4)},
		{operators.Add, ERROR, xi(1), ERROR},
		{operators.Add, xi(1), ERROR, ERROR},

		{operators.Subtract, xi(1), xi(3), xi(-2)},
		{operators.Subtract, xi(3), xi(1), xi(2)},
		{operators.Subtract, xs("3"), xs("1"), xi(2)},
		{operators.Subtract, ERROR, xi(1), ERROR},
		{operators.Subtract, xi(1), ERROR, ERROR},

		{operators.Multiply, xi(2), xi(3), xi(6)},
		{operators.Multiply, xn("1.5"), xn("2.3"), xn("3.45")},
		{operators.Multiply, xs("2"), xs("3"), xi(6)},
		{operators.Multiply, ERROR, xi(1), ERROR},
		{operators.Multiply, xi(1), ERROR, ERROR},

		{operators.Divide, xi(3), xi(2), xn("1.5")},
		{operators.Divide, xs("3"), xs("2"), xn("1.5")},
		{operators.Divide, xi(3), xi(0), ERROR},
		{operators.Divide, ERROR, xi(1), ERROR},
		{operators.Divide, xi(1), ERROR, ERROR},

		{operators.Exponent, xi(3), xi(2), xi(9)},
		{operators.Exponent, xs("3"), xs("2"), xi(9)},
		{operators.Exponent, xn("2"), xn("32.000"), xn("4294967296")},
		{operators.Exponent, xn("9"), xn("0.5"), xn("3")},
		{operators.Exponent, xn("4"), xn("2.5"), xn("32")},
		{operators.Exponent, ERROR, xi(1), ERROR},
		{operators.Exponent, xi(1), ERROR, ERROR},

		{operators.LessThan, xi(2), xi(3), types.XBooleanTrue},
		{operators.LessThan, xi(3), xi(3), types.XBooleanFalse},
		{operators.LessThan, xi(4), xi(3), types.XBooleanFalse},
		{operators.LessThan, ERROR, xi(1), ERROR},
		{operators.LessThan, xi(1), ERROR, ERROR},

		{operators.LessThanOrEqual, xi(2), xi(3), types.XBooleanTrue},
		{operators.LessThanOrEqual, xi(3), xi(3), types.XBooleanTrue},
		{operators.LessThanOrEqual, xi(4), xi(3), types.XBooleanFalse},
		{operators.LessThanOrEqual, ERROR, xi(1), ERROR},
		{operators.LessThanOrEqual, xi(1), ERROR, ERROR},

		{operators.GreaterThan, xi(2), xi(3), types.XBooleanFalse},
		{operators.GreaterThan, xi(3), xi(3), types.XBooleanFalse},
		{operators.GreaterThan, xi(4), xi(3), types.XBooleanTrue},
		{operators.GreaterThan, ERROR, xi(1), ERROR},
		{operators.GreaterThan, xi(1), ERROR, ERROR},

		{operators.GreaterThanOrEqual, xi(2), xi(3), types.XBooleanFalse},
		{operators.GreaterThanOrEqual, xi(3), xi(3), types.XBooleanTrue},
		{operators.GreaterThanOrEqual, xi(4), xi(3), types.XBooleanTrue},
		{operators.GreaterThanOrEqual, ERROR, xi(1), ERROR},
		{operators.GreaterThanOrEqual, xi(1), ERROR, ERROR},
	}

	for _, tc := range testCases {
		testID := fmt.Sprintf("%v(%s, %s)", tc.operator, tc.arg1, tc.arg2)

		result := tc.operator(env, tc.arg1, tc.arg2)

		// don't check error equality - just check that we got an error if we expected one
		if tc.expected == ERROR {
			assert.True(t, types.IsXError(result), "expecting error, got %T{%s} for ", result, result, testID)
		} else {
			test.AssertXEqual(t, tc.expected, result, "result mismatch for %s", testID)
		}
	}
}

func TestUnaryOperators(t *testing.T) {
	env := envs.NewBuilder().Build()

	testCases := []struct {
		operator operators.UnaryOperator
		arg      types.XValue
		expected types.XValue
	}{
		{operators.Negate, xi(123), xi(-123)},
		{operators.Negate, xs("123"), xi(-123)},
		{operators.Negate, xn("123.45"), xn("-123.45")},
		{operators.Negate, ERROR, ERROR},
	}

	for _, tc := range testCases {
		testID := fmt.Sprintf("%v(%s)", tc.operator, tc.arg)

		result := tc.operator(env, tc.arg)

		// don't check error equality - just check that we got an error if we expected one
		if tc.expected == ERROR {
			assert.True(t, types.IsXError(result), "expecting error, got %T{%s} for ", result, result, testID)
		} else {
			test.AssertXEqual(t, tc.expected, result, "result mismatch for %s", testID)
		}
	}
}
