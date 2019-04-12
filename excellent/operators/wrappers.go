package operators

import (
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/utils"
)

// UnaryOperator is an operator which takes a single argument
type UnaryOperator func(utils.Environment, types.XValue) types.XValue

// BinaryOperator is an operator which takes two arguments
type BinaryOperator func(utils.Environment, types.XValue, types.XValue) types.XValue

func textualBinary(f func(utils.Environment, types.XText, types.XText) types.XValue) BinaryOperator {
	return func(env utils.Environment, arg1 types.XValue, arg2 types.XValue) types.XValue {
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

func numericalUnary(f func(utils.Environment, types.XNumber) types.XValue) UnaryOperator {
	return func(env utils.Environment, arg types.XValue) types.XValue {
		num, xerr := types.ToXNumber(env, arg)
		if xerr != nil {
			return xerr
		}

		return f(env, num)
	}
}

func numericalBinary(f func(utils.Environment, types.XNumber, types.XNumber) types.XValue) BinaryOperator {
	return func(env utils.Environment, arg1 types.XValue, arg2 types.XValue) types.XValue {
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
