package operators

import (
	"strings"

	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
)

// Concatenate joins two text values together.
//
//	@("hello" & " " & "bar") -> hello bar
//	@("hello" & null) -> hello
//
// @operator concatenate "&"
var Concatenate = textualBinary("&", func(env envs.Environment, text1 *types.XText, text2 *types.XText) types.XValue {
	var buffer strings.Builder
	buffer.WriteString(text1.Native())
	buffer.WriteString(text2.Native())
	return types.NewXText(buffer.String())
})

// Equal returns true if two values are textually equal.
//
//	@("hello" = "hello") -> true
//	@("hello" = "bar") -> false
//	@(1 = "1") -> true
//
// @operator equal "="
var Equal = textualBinary("=", func(env envs.Environment, text1 *types.XText, text2 *types.XText) types.XValue {
	return types.NewXBoolean(text1.Equals(text2))
})

// NotEqual returns true if two values are textually not equal.
//
//	@("hello" != "hello") -> false
//	@("hello" != "bar") -> true
//	@(1 != 2) -> true
//
// @operator notequal "!="
var NotEqual = textualBinary("!=", func(env envs.Environment, text1 *types.XText, text2 *types.XText) types.XValue {
	return types.NewXBoolean(!text1.Equals(text2))
})

// Negate negates a number
//
//	@(-fields.age) -> -23
//
// @operator negate "- (unary)"
var Negate = numericalUnary("-", func(env envs.Environment, num *types.XNumber) types.XValue {
	return num.Neg()
})

// Add adds two numbers.
//
//	@(2 + 3) -> 5
//	@(fields.age + 10) -> 33
//
// @operator add "+"
var Add = numericalBinary("+", func(env envs.Environment, num1 *types.XNumber, num2 *types.XNumber) types.XValue {
	sum, err := num1.Add(num2)
	if err != nil {
		return types.NewXError(err)
	}
	return sum
})

// Subtract subtracts two numbers.
//
//	@(3 - 2) -> 1
//	@(2 - 3) -> -1
//
// @operator subtract "- (binary)"
var Subtract = numericalBinary("-", func(env envs.Environment, num1 *types.XNumber, num2 *types.XNumber) types.XValue {
	diff, err := num1.Sub(num2)
	if err != nil {
		return types.NewXError(err)
	}
	return diff
})

// Multiply multiplies two numbers.
//
//	@(3 * 2) -> 6
//	@(fields.age * 3) -> 69
//
// @operator multiply "*"
var Multiply = numericalBinary("*", func(env envs.Environment, num1 *types.XNumber, num2 *types.XNumber) types.XValue {
	product, err := num1.Mul(num2)
	if err != nil {
		return types.NewXError(err)
	}
	return product
})

// Divide divides a number by another.
//
//	@(4 / 2) -> 2
//	@(3 / 2) -> 1.5
//	@(46 / fields.age) -> 2
//	@(3 / 0) -> ERROR
//
// @operator divide "/"
var Divide = numericalBinary("/", func(env envs.Environment, num1 *types.XNumber, num2 *types.XNumber) types.XValue {
	quotient, err := num1.Div(num2)
	if err != nil {
		return types.NewXError(err)
	}
	return quotient
})

// Exponent raises a number to the power of a another number.
//
//	@(2 ^ 8) -> 256
//	@(2 ^ 400) -> ERROR
//
// @operator exponent "^"
var Exponent = numericalBinary("^", func(env envs.Environment, num1 *types.XNumber, num2 *types.XNumber) types.XValue {
	result, err := num1.Pow(num2)
	if err != nil {
		return types.NewXError(err)
	}
	return result
})

// LessThan returns true if the first number is less than the second.
//
//	@(2 < 3) -> true
//	@(3 < 3) -> false
//	@(4 < 3) -> false
//
// @operator lessthan "<"
var LessThan = numericalBinary("<", func(env envs.Environment, num1 *types.XNumber, num2 *types.XNumber) types.XValue {
	return types.NewXBoolean(num1.Compare(num2) < 0)
})

// LessThanOrEqual returns true if the first number is less than or equal to the second.
//
//	@(2 <= 3) -> true
//	@(3 <= 3) -> true
//	@(4 <= 3) -> false
//
// @operator lessthanorequal "<="
var LessThanOrEqual = numericalBinary("<=", func(env envs.Environment, num1 *types.XNumber, num2 *types.XNumber) types.XValue {
	return types.NewXBoolean(num1.Compare(num2) <= 0)
})

// GreaterThan returns true if the first number is greater than the second.
//
//	@(2 > 3) -> false
//	@(3 > 3) -> false
//	@(4 > 3) -> true
//
// @operator greaterthan ">"
var GreaterThan = numericalBinary(">", func(env envs.Environment, num1 *types.XNumber, num2 *types.XNumber) types.XValue {
	return types.NewXBoolean(num1.Compare(num2) > 0)
})

// GreaterThanOrEqual returns true if the first number is greater than or equal to the second.
//
//	@(2 >= 3) -> false
//	@(3 >= 3) -> true
//	@(4 >= 3) -> true
//
// @operator greaterthanorequal ">="
var GreaterThanOrEqual = numericalBinary(">=", func(env envs.Environment, num1 *types.XNumber, num2 *types.XNumber) types.XValue {
	return types.NewXBoolean(num1.Compare(num2) >= 0)
})
