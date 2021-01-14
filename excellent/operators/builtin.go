package operators

import (
	"math"
	"strings"

	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/shopspring/decimal"
)

// Concatenate joins two text values together.
//
//   @("hello" & " " & "bar") -> hello bar
//   @("hello" & null) -> hello
//
// @operator concatenate "&"
var Concatenate = textualBinary(func(env envs.Environment, text1 types.XText, text2 types.XText) types.XValue {
	var buffer strings.Builder
	buffer.WriteString(text1.Native())
	buffer.WriteString(text2.Native())
	return types.NewXText(buffer.String())
})

// Equal returns true if two values are textually equal.
//
//   @("hello" = "hello") -> true
//   @("hello" = "bar") -> false
//   @(1 = 1) -> true
//
// @operator equal "="
var Equal = textualBinary(func(env envs.Environment, text1 types.XText, text2 types.XText) types.XValue {
	return types.NewXBoolean(text1.Equals(text2))
})

// NotEqual returns true if two values are textually not equal.
//
//   @("hello" != "hello") -> false
//   @("hello" != "bar") -> true
//   @(1 != 2) -> true
//
// @operator notequal "!="
var NotEqual = textualBinary(func(env envs.Environment, text1 types.XText, text2 types.XText) types.XValue {
	return types.NewXBoolean(!text1.Equals(text2))
})

// Negate negates a number
//
//   @(-fields.age) -> -23
//
// @operator negate "- (unary)"
var Negate = numericalUnary(func(env envs.Environment, num types.XNumber) types.XValue {
	return types.NewXNumber(num.Native().Neg())
})

// Add adds two numbers.
//
//   @(2 + 3) -> 5
//   @(fields.age + 10) -> 33
//
// @operator add "+"
var Add = numericalBinary(func(env envs.Environment, num1 types.XNumber, num2 types.XNumber) types.XValue {
	return types.NewXNumber(num1.Native().Add(num2.Native()))
})

// Subtract subtracts two numbers.
//
//   @(3 - 2) -> 1
//   @(2 - 3) -> -1
//
// @operator subtract "- (binary)"
var Subtract = numericalBinary(func(env envs.Environment, num1 types.XNumber, num2 types.XNumber) types.XValue {
	return types.NewXNumber(num1.Native().Sub(num2.Native()))
})

// Multiply multiplies two numbers.
//
//   @(3 * 2) -> 6
//   @(fields.age * 3) -> 69
//
// @operator multiply "*"
var Multiply = numericalBinary(func(env envs.Environment, num1 types.XNumber, num2 types.XNumber) types.XValue {
	return types.NewXNumber(num1.Native().Mul(num2.Native()))
})

// Divide divides a number by another.
//
//   @(4 / 2) -> 2
//   @(3 / 2) -> 1.5
//   @(46 / fields.age) -> 2
//   @(3 / 0) -> ERROR
//
// @operator divide "/"
var Divide = numericalBinary(func(env envs.Environment, num1 types.XNumber, num2 types.XNumber) types.XValue {
	if num2.Equals(types.XNumberZero) {
		return types.NewXErrorf("division by zero")
	}

	return types.NewXNumber(num1.Native().Div(num2.Native()))
})

// Exponent raises a number to the power of a another number.
//
//   @(2 ^ 8) -> 256
//
// @operator exponent "^"
var Exponent = numericalBinary(func(env envs.Environment, num1 types.XNumber, num2 types.XNumber) types.XValue {
	d1 := num1.Native()
	d2 := num2.Native()

	// TODO there is currently a bug in shopspring/decimal which means that only the integer part of the
	// exponent is considered (see https://github.com/nyaruka/goflow/issues/984). If we have a whole number,
	// we can use the library function, otherwise fallback to float64 math.

	if decimal.New(d2.IntPart(), 0).Equals(d2) {
		return types.NewXNumber(d1.Pow(d2))
	}

	f1, _ := d1.Float64()
	f2, _ := d2.Float64()

	return types.NewXNumber(decimal.NewFromFloat(math.Pow(f1, f2)))
})

// LessThan returns true if the first number is less than the second.
//
//   @(2 < 3) -> true
//   @(3 < 3) -> false
//   @(4 < 3) -> false
//
// @operator lessthan "<"
var LessThan = numericalBinary(func(env envs.Environment, num1 types.XNumber, num2 types.XNumber) types.XValue {
	return types.NewXBoolean(num1.Compare(num2) < 0)
})

// LessThanOrEqual returns true if the first number is less than or equal to the second.
//
//   @(2 <= 3) -> true
//   @(3 <= 3) -> true
//   @(4 <= 3) -> false
//
// @operator lessthanorequal "<="
var LessThanOrEqual = numericalBinary(func(env envs.Environment, num1 types.XNumber, num2 types.XNumber) types.XValue {
	return types.NewXBoolean(num1.Compare(num2) <= 0)
})

// GreaterThan returns true if the first number is greater than the second.
//
//   @(2 > 3) -> false
//   @(3 > 3) -> false
//   @(4 > 3) -> true
//
// @operator greaterthan ">"
var GreaterThan = numericalBinary(func(env envs.Environment, num1 types.XNumber, num2 types.XNumber) types.XValue {
	return types.NewXBoolean(num1.Compare(num2) > 0)
})

// GreaterThanOrEqual returns true if the first number is greater than or equal to the second.
//
//   @(2 >= 3) -> false
//   @(3 >= 3) -> true
//   @(4 >= 3) -> true
//
// @operator greaterthanorequal ">="
var GreaterThanOrEqual = numericalBinary(func(env envs.Environment, num1 types.XNumber, num2 types.XNumber) types.XValue {
	return types.NewXBoolean(num1.Compare(num2) >= 0)
})
