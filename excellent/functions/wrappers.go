package functions

import (
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/utils"
	"strings"
)

// ArgCountCheck wraps an XFunction and checks the number of args
func ArgCountCheck(name string, count int, f XFunction) XFunction {
	return func(env utils.Environment, args ...types.XValue) types.XValue {
		if len(args) != count {
			return types.NewXErrorf("%s takes %d argument(s), got %d", strings.ToUpper(name), count, len(args))
		}
		return f(env, args...)
	}
}

// OneStringFunction creates an XFunction from a single string function
func OneStringFunction(name string, f func(utils.Environment, types.XString) types.XValue) XFunction {
	return ArgCountCheck(name, 1, func(env utils.Environment, args ...types.XValue) types.XValue {
		str, xerr := types.ToXString(args[0])
		if xerr != nil {
			return xerr
		}
		return f(env, str)
	})
}

// TwoStringFunction creates an XFunction from a function that takes two strings
func TwoStringFunction(name string, f func(utils.Environment, types.XString, types.XString) types.XValue) XFunction {
	return ArgCountCheck(name, 2, func(env utils.Environment, args ...types.XValue) types.XValue {
		str1, xerr := types.ToXString(args[0])
		if xerr != nil {
			return xerr
		}
		str2, xerr := types.ToXString(args[1])
		if xerr != nil {
			return xerr
		}
		return f(env, str1, str2)
	})
}

// StringAndNumberFunction creates an XFunction from a function that takes a string and a number
func StringAndNumberFunction(name string, f func(utils.Environment, types.XString, types.XNumber) types.XValue) XFunction {
	return ArgCountCheck(name, 2, func(env utils.Environment, args ...types.XValue) types.XValue {
		str, xerr := types.ToXString(args[0])
		if xerr != nil {
			return xerr
		}
		num, xerr := types.ToXNumber(args[1])
		if xerr != nil {
			return xerr
		}

		return f(env, str, num)
	})
}

// StringAndIntegerFunction creates an XFunction from a function that takes a string and an integer
func StringAndIntegerFunction(name string, f func(utils.Environment, types.XString, int) types.XValue) XFunction {
	return ArgCountCheck(name, 2, func(env utils.Environment, args ...types.XValue) types.XValue {
		str, xerr := types.ToXString(args[0])
		if xerr != nil {
			return xerr
		}
		num, xerr := types.ToInteger(args[1])
		if xerr != nil {
			return xerr
		}

		return f(env, str, num)
	})
}

// StringAndDateFunction creates an XFunction from a function that takes a string and a date
func StringAndDateFunction(name string, f func(utils.Environment, types.XString, types.XTime) types.XValue) XFunction {
	return ArgCountCheck(name, 2, func(env utils.Environment, args ...types.XValue) types.XValue {
		str, xerr := types.ToXString(args[0])
		if xerr != nil {
			return xerr
		}
		date, xerr := types.ToXTime(env, args[1])
		if xerr != nil {
			return xerr
		}

		return f(env, str, date)
	})
}

// OneNumberFunction creates an XFunction from a single number function
func OneNumberFunction(name string, f func(utils.Environment, types.XNumber) types.XValue) XFunction {
	return ArgCountCheck(name, 1, func(env utils.Environment, args ...types.XValue) types.XValue {
		num, xerr := types.ToXNumber(args[0])
		if xerr != nil {
			return xerr
		}

		return f(env, num)
	})
}

// TwoNumberFunction creates an XFunction from a function that takes two numbers
func TwoNumberFunction(name string, f func(utils.Environment, types.XNumber, types.XNumber) types.XValue) XFunction {
	return ArgCountCheck(name, 2, func(env utils.Environment, args ...types.XValue) types.XValue {
		num1, xerr := types.ToXNumber(args[0])
		if xerr != nil {
			return xerr
		}
		num2, xerr := types.ToXNumber(args[1])
		if xerr != nil {
			return xerr
		}

		return f(env, num1, num2)
	})
}

// OneDateFunction creates an XFunction from a single number function
func OneDateFunction(name string, f func(utils.Environment, types.XTime) types.XValue) XFunction {
	return ArgCountCheck(name, 1, func(env utils.Environment, args ...types.XValue) types.XValue {
		date, xerr := types.ToXTime(env, args[0])
		if xerr != nil {
			return xerr
		}

		return f(env, date)
	})
}
