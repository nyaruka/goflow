package operators

import (
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
)

// UnaryOperator is an operator which takes a single argument
type UnaryOperator func(envs.Environment, types.XValue) types.XValue

// BinaryOperator is an operator which takes two arguments
type BinaryOperator func(envs.Environment, types.XValue, types.XValue) types.XValue

func textualBinary(f func(envs.Environment, types.XText, types.XText) types.XValue) BinaryOperator {
	return func(env envs.Environment, arg1 types.XValue, arg2 types.XValue) types.XValue {
		text1, xerr := types.ToXText(env, arg1)
		if xerr != nil {
			return xerr
		}
		text2, xerr := types.ToXText(env, arg2)
		if xerr != nil {
			return xerr
		}

		return f(env, text1, text2)
	}
}

func numericalUnary(f func(envs.Environment, types.XNumber) types.XValue) UnaryOperator {
	return func(env envs.Environment, arg types.XValue) types.XValue {
		num, xerr := types.ToXNumber(env, arg)
		if xerr != nil {
			return xerr
		}

		return f(env, num)
	}
}

func numericalBinary(f func(envs.Environment, types.XNumber, types.XNumber) types.XValue) BinaryOperator {
	return func(env envs.Environment, arg1 types.XValue, arg2 types.XValue) types.XValue {
		num1, xerr := types.ToXNumber(env, arg1)
		if xerr != nil {
			return xerr
		}
		num2, xerr := types.ToXNumber(env, arg2)
		if xerr != nil {
			return xerr
		}

		return f(env, num1, num2)
	}
}
