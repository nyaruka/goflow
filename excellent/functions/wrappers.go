package functions

import (
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
)

// NumArgsCheck wraps an XFunction and checks the number of args
func NumArgsCheck(num int, f types.XFunction) types.XFunction {
	return MinAndMaxArgsCheck(num, num, f)
}

// MinArgsCheck wraps an XFunction and checks the minimum number of args
func MinArgsCheck(min int, f types.XFunction) types.XFunction {
	return MinAndMaxArgsCheck(min, -1, f)
}

// MinAndMaxArgsCheck wraps an XFunction and checks the number of args
func MinAndMaxArgsCheck(min int, max int, f types.XFunction) types.XFunction {
	return func(env envs.Environment, args ...types.XValue) types.XValue {
		if min == max {
			// function requires a fixed number of arguments
			if len(args) != min {
				return types.NewXErrorf("need %d argument(s), got %d", min, len(args))
			}
		} else if max < 0 {
			// function requires a minimum number of arguments
			if len(args) < min {
				return types.NewXErrorf("need at least %d argument(s), got %d", min, len(args))
			}
		} else {
			// function requires the given range of arguments
			if len(args) < min || len(args) > max {
				return types.NewXErrorf("need %d to %d argument(s), got %d", min, max, len(args))
			}
		}

		return f(env, args...)
	}
}

// NoArgFunction creates an XFunction from a no-arg function
func NoArgFunction(f func(envs.Environment) types.XValue) types.XFunction {
	return NumArgsCheck(0, func(env envs.Environment, args ...types.XValue) types.XValue {
		return f(env)
	})
}

// OneArgFunction creates an XFunction from a single-arg function
func OneArgFunction(f func(envs.Environment, types.XValue) types.XValue) types.XFunction {
	return NumArgsCheck(1, func(env envs.Environment, args ...types.XValue) types.XValue {
		return f(env, args[0])
	})
}

// TwoArgFunction creates an XFunction from a two-arg function
func TwoArgFunction(f func(envs.Environment, types.XValue, types.XValue) types.XValue) types.XFunction {
	return NumArgsCheck(2, func(env envs.Environment, args ...types.XValue) types.XValue {
		return f(env, args[0], args[1])
	})
}

// ThreeArgFunction creates an XFunction from a three-arg function
func ThreeArgFunction(f func(envs.Environment, types.XValue, types.XValue, types.XValue) types.XValue) types.XFunction {
	return NumArgsCheck(3, func(env envs.Environment, args ...types.XValue) types.XValue {
		return f(env, args[0], args[1], args[2])
	})
}

// OneTextFunction creates an XFunction from a function that takes a single text arg
func OneTextFunction(f func(envs.Environment, types.XText) types.XValue) types.XFunction {
	return NumArgsCheck(1, func(env envs.Environment, args ...types.XValue) types.XValue {
		str, xerr := types.ToXText(env, args[0])
		if xerr != nil {
			return xerr
		}
		return f(env, str)
	})
}

// TwoTextFunction creates an XFunction from a function that takes two text args
func TwoTextFunction(f func(envs.Environment, types.XText, types.XText) types.XValue) types.XFunction {
	return NumArgsCheck(2, func(env envs.Environment, args ...types.XValue) types.XValue {
		str1, xerr := types.ToXText(env, args[0])
		if xerr != nil {
			return xerr
		}
		str2, xerr := types.ToXText(env, args[1])
		if xerr != nil {
			return xerr
		}
		return f(env, str1, str2)
	})
}

// TextAndNumberFunction creates an XFunction from a function that takes a text and a number arg
func TextAndNumberFunction(f func(envs.Environment, types.XText, types.XNumber) types.XValue) types.XFunction {
	return NumArgsCheck(2, func(env envs.Environment, args ...types.XValue) types.XValue {
		str, xerr := types.ToXText(env, args[0])
		if xerr != nil {
			return xerr
		}
		num, xerr := types.ToXNumber(env, args[1])
		if xerr != nil {
			return xerr
		}

		return f(env, str, num)
	})
}

// TextAndIntegerFunction creates an XFunction from a function that takes a text and an integer arg
func TextAndIntegerFunction(f func(envs.Environment, types.XText, int) types.XValue) types.XFunction {
	return NumArgsCheck(2, func(env envs.Environment, args ...types.XValue) types.XValue {
		str, xerr := types.ToXText(env, args[0])
		if xerr != nil {
			return xerr
		}
		num, xerr := types.ToInteger(env, args[1])
		if xerr != nil {
			return xerr
		}

		return f(env, str, num)
	})
}

// TextAndOptionalTextFunction creates an XFunction from a function that takes either one or two text args
func TextAndOptionalTextFunction(f func(envs.Environment, types.XText, types.XText) types.XValue, defaultVal types.XText) types.XFunction {
	return MinAndMaxArgsCheck(1, 2, func(env envs.Environment, args ...types.XValue) types.XValue {
		str1, xerr := types.ToXText(env, args[0])
		if xerr != nil {
			return xerr
		}

		str2 := defaultVal
		if len(args) == 2 {
			str2, xerr = types.ToXText(env, args[1])
			if xerr != nil {
				return xerr
			}
		}

		return f(env, str1, str2)
	})
}

// ThreeIntegerFunction creates an XFunction from a function that takes a text and an integer arg
func ThreeIntegerFunction(f func(envs.Environment, int, int, int) types.XValue) types.XFunction {
	return NumArgsCheck(3, func(env envs.Environment, args ...types.XValue) types.XValue {
		num1, xerr := types.ToInteger(env, args[0])
		if xerr != nil {
			return xerr
		}
		num2, xerr := types.ToInteger(env, args[1])
		if xerr != nil {
			return xerr
		}
		num3, xerr := types.ToInteger(env, args[2])
		if xerr != nil {
			return xerr
		}

		return f(env, num1, num2, num3)
	})
}

// TextAndDateFunction creates an XFunction from a function that takes a text and a date arg
func TextAndDateFunction(f func(envs.Environment, types.XText, types.XDateTime) types.XValue) types.XFunction {
	return NumArgsCheck(2, func(env envs.Environment, args ...types.XValue) types.XValue {
		str, xerr := types.ToXText(env, args[0])
		if xerr != nil {
			return xerr
		}
		date, xerr := types.ToXDateTime(env, args[1])
		if xerr != nil {
			return xerr
		}

		return f(env, str, date)
	})
}

// InitialTextFunction creates an XFunction from a function that takes an initial text arg followed by other args
func InitialTextFunction(minOtherArgs int, maxOtherArgs int, f func(envs.Environment, types.XText, ...types.XValue) types.XValue) types.XFunction {
	return MinAndMaxArgsCheck(minOtherArgs+1, maxOtherArgs+1, func(env envs.Environment, args ...types.XValue) types.XValue {
		str, xerr := types.ToXText(env, args[0])
		if xerr != nil {
			return xerr
		}
		return f(env, str, args[1:]...)
	})
}

// OneNumberFunction creates an XFunction from a single number function
func OneNumberFunction(f func(envs.Environment, types.XNumber) types.XValue) types.XFunction {
	return NumArgsCheck(1, func(env envs.Environment, args ...types.XValue) types.XValue {
		num, xerr := types.ToXNumber(env, args[0])
		if xerr != nil {
			return xerr
		}

		return f(env, num)
	})
}

// OneNumberAndOptionalIntegerFunction creates an XFunction from a function that takes a number and an optional integer
func OneNumberAndOptionalIntegerFunction(f func(envs.Environment, types.XNumber, int) types.XValue, defaultVal int) types.XFunction {
	return MinAndMaxArgsCheck(1, 2, func(env envs.Environment, args ...types.XValue) types.XValue {
		num, xerr := types.ToXNumber(env, args[0])
		if xerr != nil {
			return xerr
		}

		intVal := defaultVal
		if len(args) == 2 {
			intVal, xerr = types.ToInteger(env, args[1])
			if xerr != nil {
				return xerr
			}
		}

		return f(env, num, intVal)
	})
}

// TwoNumberFunction creates an XFunction from a function that takes two numbers
func TwoNumberFunction(f func(envs.Environment, types.XNumber, types.XNumber) types.XValue) types.XFunction {
	return NumArgsCheck(2, func(env envs.Environment, args ...types.XValue) types.XValue {
		num1, xerr := types.ToXNumber(env, args[0])
		if xerr != nil {
			return xerr
		}
		num2, xerr := types.ToXNumber(env, args[1])
		if xerr != nil {
			return xerr
		}

		return f(env, num1, num2)
	})
}

// OneDateFunction creates an XFunction from a single date function
func OneDateFunction(f func(envs.Environment, types.XDate) types.XValue) types.XFunction {
	return NumArgsCheck(1, func(env envs.Environment, args ...types.XValue) types.XValue {
		date, xerr := types.ToXDate(env, args[0])
		if xerr != nil {
			return xerr
		}

		return f(env, date)
	})
}

// OneDateTimeFunction creates an XFunction from a single datetime function
func OneDateTimeFunction(f func(envs.Environment, types.XDateTime) types.XValue) types.XFunction {
	return NumArgsCheck(1, func(env envs.Environment, args ...types.XValue) types.XValue {
		date, xerr := types.ToXDateTime(env, args[0])
		if xerr != nil {
			return xerr
		}

		return f(env, date)
	})
}

// ObjectTextAndNumberFunction creates an XFunction from a function that takes an object, text and a number
func ObjectTextAndNumberFunction(f func(envs.Environment, *types.XObject, types.XText, types.XNumber) types.XValue) types.XFunction {
	return NumArgsCheck(3, func(env envs.Environment, args ...types.XValue) types.XValue {
		object, xerr := types.ToXObject(env, args[0])
		if xerr != nil {
			return xerr
		}
		text, xerr := types.ToXText(env, args[1])
		if xerr != nil {
			return xerr
		}
		num, xerr := types.ToXNumber(env, args[2])
		if xerr != nil {
			return xerr
		}

		return f(env, object, text, num)
	})
}

// ObjectAndTextsFunction creates an XFunction from a function that takes an object and any number of text values
func ObjectAndTextsFunction(f func(envs.Environment, *types.XObject, ...types.XText) types.XValue) types.XFunction {
	return MinArgsCheck(2, func(env envs.Environment, args ...types.XValue) types.XValue {
		object, xerr := types.ToXObject(env, args[0])
		if xerr != nil {
			return xerr
		}

		texts := make([]types.XText, len(args)-1)
		for i, arg := range args[1:] {
			text, xerr := types.ToXText(env, arg)
			if xerr != nil {
				return xerr
			}
			texts[i] = text
		}

		return f(env, object, texts...)
	})
}
