package functions_test

import (
	"testing"

	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/functions"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/test"
)

func TestWrappers(t *testing.T) {
	env := envs.NewBuilder().Build()
	result := types.NewXText("Success")
	text := types.NewXText("X")
	num := types.RequireXNumberFromString("1")
	obj := types.XObjectEmpty
	xe := types.NewXErrorf

	f := functions.MinArgsCheck(2, func(envs.Environment, ...types.XValue) types.XValue { return result })
	test.AssertXEqual(t, xe("need at least 2 argument(s), got 0"), f(env))
	test.AssertXEqual(t, xe("need at least 2 argument(s), got 1"), f(env, num))
	test.AssertXEqual(t, result, f(env, num, num))
	test.AssertXEqual(t, result, f(env, num, num, num))

	f = functions.NoArgFunction(func(envs.Environment) types.XValue { return result })
	test.AssertXEqual(t, result, f(env))
	test.AssertXEqual(t, xe("need 0 argument(s), got 1"), f(env, num))

	f = functions.OneArgFunction(func(envs.Environment, types.XValue) types.XValue { return result })
	test.AssertXEqual(t, xe("need 1 argument(s), got 0"), f(env))
	test.AssertXEqual(t, result, f(env, types.NewXText("1")))
	test.AssertXEqual(t, xe("need 1 argument(s), got 2"), f(env, num, num))

	f = functions.TwoArgFunction(func(envs.Environment, types.XValue, types.XValue) types.XValue { return result })
	test.AssertXEqual(t, xe("need 2 argument(s), got 1"), f(env, num))
	test.AssertXEqual(t, result, f(env, num, num))
	test.AssertXEqual(t, xe("need 2 argument(s), got 3"), f(env, num, num, num))

	f = functions.ThreeArgFunction(func(envs.Environment, types.XValue, types.XValue, types.XValue) types.XValue { return result })
	test.AssertXEqual(t, xe("need 3 argument(s), got 2"), f(env, num, num))
	test.AssertXEqual(t, result, f(env, num, num, num))
	test.AssertXEqual(t, xe("need 3 argument(s), got 4"), f(env, num, num, num, num))

	f = functions.OneTextFunction(func(envs.Environment, types.XText) types.XValue { return result })
	test.AssertXEqual(t, xe("need 1 argument(s), got 0"), f(env))
	test.AssertXEqual(t, result, f(env, text))
	test.AssertXEqual(t, result, f(env, num))
	test.AssertXEqual(t, xe("error"), f(env, xe("error")))
	test.AssertXEqual(t, xe("need 1 argument(s), got 2"), f(env, text, text))

	f = functions.TwoTextFunction(func(envs.Environment, types.XText, types.XText) types.XValue { return result })
	test.AssertXEqual(t, xe("need 2 argument(s), got 1"), f(env, text))
	test.AssertXEqual(t, result, f(env, text, text))
	test.AssertXEqual(t, result, f(env, num, num))
	test.AssertXEqual(t, xe("error"), f(env, xe("error"), text))
	test.AssertXEqual(t, xe("error"), f(env, text, xe("error")))
	test.AssertXEqual(t, xe("need 2 argument(s), got 3"), f(env, text, text, text))

	f = functions.OneNumberFunction(func(envs.Environment, types.XNumber) types.XValue { return result })
	test.AssertXEqual(t, xe("need 1 argument(s), got 0"), f(env))
	test.AssertXEqual(t, result, f(env, num))
	test.AssertXEqual(t, result, f(env, types.NewXText("1")))
	test.AssertXEqual(t, xe(`unable to convert "X" to a number`), f(env, text))
	test.AssertXEqual(t, xe("need 1 argument(s), got 2"), f(env, num, num))

	f = functions.TwoNumberFunction(func(envs.Environment, types.XNumber, types.XNumber) types.XValue { return result })
	test.AssertXEqual(t, xe("need 2 argument(s), got 1"), f(env, num))
	test.AssertXEqual(t, result, f(env, num, num))
	test.AssertXEqual(t, result, f(env, types.NewXText("1"), types.NewXText("2")))
	test.AssertXEqual(t, xe(`unable to convert "X" to a number`), f(env, types.NewXText("X"), num))
	test.AssertXEqual(t, xe(`unable to convert "X" to a number`), f(env, num, types.NewXText("X")))
	test.AssertXEqual(t, xe("need 2 argument(s), got 3"), f(env, num, num, num))

	f = functions.TextAndNumberFunction(func(envs.Environment, types.XText, types.XNumber) types.XValue { return result })
	test.AssertXEqual(t, xe("need 2 argument(s), got 1"), f(env, text))
	test.AssertXEqual(t, result, f(env, text, num))
	test.AssertXEqual(t, result, f(env, num, num))
	test.AssertXEqual(t, result, f(env, text, types.NewXText("2")))
	test.AssertXEqual(t, xe("error"), f(env, xe("error"), num))
	test.AssertXEqual(t, xe(`unable to convert "X" to a number`), f(env, text, types.NewXText("X")))

	f = functions.ObjectTextAndNumberFunction(func(envs.Environment, *types.XObject, types.XText, types.XNumber) types.XValue { return result })
	test.AssertXEqual(t, xe("need 3 argument(s), got 2"), f(env, obj, text))
	test.AssertXEqual(t, result, f(env, obj, text, num))
	test.AssertXEqual(t, xe("unable to convert 1 to an object"), f(env, num, text, num))
	test.AssertXEqual(t, xe("error"), f(env, obj, xe("error"), num))
	test.AssertXEqual(t, xe(`unable to convert "X" to a number`), f(env, obj, text, text))

	f = functions.ObjectAndTextsFunction(func(envs.Environment, *types.XObject, ...types.XText) types.XValue { return result })
	test.AssertXEqual(t, xe("need at least 2 argument(s), got 1"), f(env, obj))
	test.AssertXEqual(t, result, f(env, obj, text, text, text))
	test.AssertXEqual(t, xe("unable to convert 1 to an object"), f(env, num, text))
	test.AssertXEqual(t, xe("error"), f(env, obj, xe("error")))
}
