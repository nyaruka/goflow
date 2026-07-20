package operators

import (
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
)

// Unary is an operator which takes a single operand
type Unary struct {
	symbol string
	fn     func(envs.Environment, types.XValue) types.XValue
}

// Symbol returns the symbol used for this operator in expressions
func (o *Unary) Symbol() string { return o.symbol }

// Evaluate applies this operator to the given operand
func (o *Unary) Evaluate(env envs.Environment, arg types.XValue) types.XValue {
	return o.fn(env, arg)
}

// Binary is an operator which takes two operands
type Binary struct {
	symbol string
	fn     func(envs.Environment, types.XValue, types.XValue) types.XValue
}

// Symbol returns the symbol used for this operator in expressions
func (o *Binary) Symbol() string { return o.symbol }

// Evaluate applies this operator to the given operands
func (o *Binary) Evaluate(env envs.Environment, arg1 types.XValue, arg2 types.XValue) types.XValue {
	return o.fn(env, arg1, arg2)
}

func textualBinary(symbol string, f func(envs.Environment, *types.XText, *types.XText) types.XValue) *Binary {
	return &Binary{symbol, func(env envs.Environment, arg1 types.XValue, arg2 types.XValue) types.XValue {
		text1, xerr := types.ToXText(env, arg1)
		if xerr != nil {
			return xerr
		}
		text2, xerr := types.ToXText(env, arg2)
		if xerr != nil {
			return xerr
		}

		return f(env, text1, text2)
	}}
}

func numericalUnary(symbol string, f func(envs.Environment, *types.XNumber) types.XValue) *Unary {
	return &Unary{symbol, func(env envs.Environment, arg types.XValue) types.XValue {
		num, xerr := types.ToXNumber(env, arg)
		if xerr != nil {
			return xerr
		}

		return f(env, num)
	}}
}

func numericalBinary(symbol string, f func(envs.Environment, *types.XNumber, *types.XNumber) types.XValue) *Binary {
	return &Binary{symbol, func(env envs.Environment, arg1 types.XValue, arg2 types.XValue) types.XValue {
		num1, xerr := types.ToXNumber(env, arg1)
		if xerr != nil {
			return xerr
		}
		num2, xerr := types.ToXNumber(env, arg2)
		if xerr != nil {
			return xerr
		}

		return f(env, num1, num2)
	}}
}
