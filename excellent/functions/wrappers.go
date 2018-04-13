package functions

import (
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/utils"
)

// ArgCountCheck wraps an XFunction and checks the number of args
func ArgCountCheck(min int, max int, f XFunction) XFunction {
	return func(env utils.Environment, args ...types.XValue) types.XValue {
		if min == max {
			if len(args) != min {
				return types.NewXErrorf("need %d argument(s), got %d", min, len(args))
			}
		} else {
			if len(args) < min || len(args) > max {
				return types.NewXErrorf("need %d to %d argument(s), got %d", min, max, len(args))
			}
		}

		return f(env, args...)
	}
}

// NoArgFunction creates an XFunction from a no-arg function
func NoArgFunction(f func(utils.Environment) types.XValue) XFunction {
	return ArgCountCheck(0, 0, func(env utils.Environment, args ...types.XValue) types.XValue {
		return f(env)
	})
}

// OneArgFunction creates an XFunction from a single-arg function
func OneArgFunction(f func(utils.Environment, types.XValue) types.XValue) XFunction {
	return ArgCountCheck(1, 1, func(env utils.Environment, args ...types.XValue) types.XValue {
		return f(env, args[0])
	})
}

// TwoArgFunction creates an XFunction from a two-arg function
func TwoArgFunction(f func(utils.Environment, types.XValue, types.XValue) types.XValue) XFunction {
	return ArgCountCheck(2, 2, func(env utils.Environment, args ...types.XValue) types.XValue {
		return f(env, args[0], args[1])
	})
}

// ThreeArgFunction creates an XFunction from a three-arg function
func ThreeArgFunction(f func(utils.Environment, types.XValue, types.XValue, types.XValue) types.XValue) XFunction {
	return ArgCountCheck(3, 3, func(env utils.Environment, args ...types.XValue) types.XValue {
		return f(env, args[0], args[1], args[2])
	})
}

// OneStringFunction creates an XFunction from a single string function
func OneStringFunction(f func(utils.Environment, types.XString) types.XValue) XFunction {
	return ArgCountCheck(1, 1, func(env utils.Environment, args ...types.XValue) types.XValue {
		str, xerr := types.ToXString(args[0])
		if xerr != nil {
			return xerr
		}
		return f(env, str)
	})
}

// TwoStringFunction creates an XFunction from a function that takes two strings
func TwoStringFunction(f func(utils.Environment, types.XString, types.XString) types.XValue) XFunction {
	return ArgCountCheck(2, 2, func(env utils.Environment, args ...types.XValue) types.XValue {
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

// ThreeStringFunction creates an XFunction from a function that takes three strings
func ThreeStringFunction(f func(utils.Environment, types.XString, types.XString, types.XString) types.XValue) XFunction {
	return ArgCountCheck(3, 3, func(env utils.Environment, args ...types.XValue) types.XValue {
		str1, xerr := types.ToXString(args[0])
		if xerr != nil {
			return xerr
		}
		str2, xerr := types.ToXString(args[1])
		if xerr != nil {
			return xerr
		}
		str3, xerr := types.ToXString(args[2])
		if xerr != nil {
			return xerr
		}
		return f(env, str1, str2, str3)
	})
}

// StringAndNumberFunction creates an XFunction from a function that takes a string and a number
func StringAndNumberFunction(f func(utils.Environment, types.XString, types.XNumber) types.XValue) XFunction {
	return ArgCountCheck(2, 2, func(env utils.Environment, args ...types.XValue) types.XValue {
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
func StringAndIntegerFunction(f func(utils.Environment, types.XString, int) types.XValue) XFunction {
	return ArgCountCheck(2, 2, func(env utils.Environment, args ...types.XValue) types.XValue {
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
func StringAndDateFunction(f func(utils.Environment, types.XString, types.XDate) types.XValue) XFunction {
	return ArgCountCheck(2, 2, func(env utils.Environment, args ...types.XValue) types.XValue {
		str, xerr := types.ToXString(args[0])
		if xerr != nil {
			return xerr
		}
		date, xerr := types.ToXDate(env, args[1])
		if xerr != nil {
			return xerr
		}

		return f(env, str, date)
	})
}

// OneNumberFunction creates an XFunction from a single number function
func OneNumberFunction(f func(utils.Environment, types.XNumber) types.XValue) XFunction {
	return ArgCountCheck(1, 1, func(env utils.Environment, args ...types.XValue) types.XValue {
		num, xerr := types.ToXNumber(args[0])
		if xerr != nil {
			return xerr
		}

		return f(env, num)
	})
}

// OneNumberAndOptionalIntegerFunction creates an XFunction from a function that takes a number and an optional integer
func OneNumberAndOptionalIntegerFunction(f func(utils.Environment, types.XNumber, int) types.XValue, defaultVal int) XFunction {
	return ArgCountCheck(1, 2, func(env utils.Environment, args ...types.XValue) types.XValue {
		num, xerr := types.ToXNumber(args[0])
		if xerr != nil {
			return xerr
		}

		intVal := defaultVal
		if len(args) == 2 {
			intVal, xerr = types.ToInteger(args[1])
			if xerr != nil {
				return xerr
			}
		}

		return f(env, num, intVal)
	})
}

// TwoNumberFunction creates an XFunction from a function that takes two numbers
func TwoNumberFunction(f func(utils.Environment, types.XNumber, types.XNumber) types.XValue) XFunction {
	return ArgCountCheck(2, 2, func(env utils.Environment, args ...types.XValue) types.XValue {
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
func OneDateFunction(f func(utils.Environment, types.XDate) types.XValue) XFunction {
	return ArgCountCheck(1, 1, func(env utils.Environment, args ...types.XValue) types.XValue {
		date, xerr := types.ToXDate(env, args[0])
		if xerr != nil {
			return xerr
		}

		return f(env, date)
	})
}
