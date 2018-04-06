package functions

import (
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/utils"
	"strings"
)

// OneStringFunction creates an XFunction from a single string function
func OneStringFunction(name string, f func(utils.Environment, types.XString) types.XValue) XFunction {
	return func(env utils.Environment, args ...types.XValue) types.XValue {
		if len(args) != 1 {
			return types.NewXErrorf("%s takes a single argument, got %d", strings.ToUpper(name), len(args))
		}

		// if argument is an error, return immediately
		if types.IsError(args[0]) {
			return args[0]
		}

		return f(env, types.ToXString(args[0]))
	}
}

// TwoStringFunction creates an XFunction from a function that takes two strings
func TwoStringFunction(name string, f func(utils.Environment, types.XString, types.XString) types.XValue) XFunction {
	return func(env utils.Environment, args ...types.XValue) types.XValue {
		if len(args) != 2 {
			return types.NewXErrorf("%s takes two arguments, got %d", strings.ToUpper(name), len(args))
		}

		// if either argument is an error, return immediately
		if types.IsError(args[0]) {
			return args[0]
		}
		if types.IsError(args[1]) {
			return args[1]
		}

		return f(env, types.ToXString(args[0]), types.ToXString(args[1]))
	}
}

// StringAndIntegerFunction creates an XFunction from a function that takes a string and an integer
func StringAndIntegerFunction(name string, f func(utils.Environment, types.XString, int) types.XValue) XFunction {
	return func(env utils.Environment, args ...types.XValue) types.XValue {
		if len(args) != 2 {
			return types.NewXErrorf("%s takes two arguments, got %d", strings.ToUpper(name), len(args))
		}

		// if either argument is an error, return immediately
		if types.IsError(args[0]) {
			return args[0]
		}
		if types.IsError(args[1]) {
			return args[1]
		}

		num, err := types.ToInteger(args[1])
		if err != nil {
			return types.NewXError(err)
		}

		return f(env, types.ToXString(args[0]), num)
	}
}

// OneNumberFunction creates an XFunction from a single number function
func OneNumberFunction(name string, f func(utils.Environment, types.XNumber) types.XValue) XFunction {
	return func(env utils.Environment, args ...types.XValue) types.XValue {
		if len(args) != 1 {
			return types.NewXErrorf("%s takes a single argument, got %d", strings.ToUpper(name), len(args))
		}

		// if argument is an error, return immediately
		if types.IsError(args[0]) {
			return args[0]
		}

		num, err := types.ToXNumber(args[0])
		if err != nil {
			return types.NewXError(err)
		}

		return f(env, num)
	}
}

// TwoNumberFunction creates an XFunction from a function that takes two numbers
func TwoNumberFunction(name string, f func(utils.Environment, types.XNumber, types.XNumber) types.XValue) XFunction {
	return func(env utils.Environment, args ...types.XValue) types.XValue {
		if len(args) != 2 {
			return types.NewXErrorf("%s takes two arguments, got %d", strings.ToUpper(name), len(args))
		}

		// if either argument is an error, return immediately
		if types.IsError(args[0]) {
			return args[0]
		}
		if types.IsError(args[1]) {
			return args[1]
		}

		num1, err := types.ToXNumber(args[0])
		if err != nil {
			return types.NewXError(err)
		}
		num2, err := types.ToXNumber(args[1])
		if err != nil {
			return types.NewXError(err)
		}

		return f(env, num1, num2)
	}
}

// OneDateFunction creates an XFunction from a single number function
func OneDateFunction(name string, f func(utils.Environment, types.XTime) types.XValue) XFunction {
	return func(env utils.Environment, args ...types.XValue) types.XValue {
		if len(args) != 1 {
			return types.NewXErrorf("%s takes a single argument, got %d", strings.ToUpper(name), len(args))
		}

		// if argument is an error, return immediately
		if types.IsError(args[0]) {
			return args[0]
		}

		date, err := types.ToXTime(env, args[0])
		if err != nil {
			return types.NewXError(err)
		}

		return f(env, date)
	}
}
