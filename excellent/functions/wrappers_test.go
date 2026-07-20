package functions_test

import (
	"context"
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

	f := functions.MinArgsCheck(2, func(context.Context, envs.Environment, ...types.XValue) types.XValue { return result })
	test.AssertXEqual(t, xe("need at least 2 argument(s), got 0"), f(context.Background(), env))
	test.AssertXEqual(t, xe("need at least 2 argument(s), got 1"), f(context.Background(), env, num))
	test.AssertXEqual(t, result, f(context.Background(), env, num, num))
	test.AssertXEqual(t, result, f(context.Background(), env, num, num, num))

	f = functions.NoArgFunction(func(envs.Environment) types.XValue { return result })
	test.AssertXEqual(t, result, f(context.Background(), env))
	test.AssertXEqual(t, xe("need 0 argument(s), got 1"), f(context.Background(), env, num))

	f = functions.OneArgFunction(func(envs.Environment, types.XValue) types.XValue { return result })
	test.AssertXEqual(t, xe("need 1 argument(s), got 0"), f(context.Background(), env))
	test.AssertXEqual(t, result, f(context.Background(), env, types.NewXText("1")))
	test.AssertXEqual(t, xe("need 1 argument(s), got 2"), f(context.Background(), env, num, num))

	f = functions.TwoArgFunction(func(envs.Environment, types.XValue, types.XValue) types.XValue { return result })
	test.AssertXEqual(t, xe("need 2 argument(s), got 1"), f(context.Background(), env, num))
	test.AssertXEqual(t, result, f(context.Background(), env, num, num))
	test.AssertXEqual(t, xe("need 2 argument(s), got 3"), f(context.Background(), env, num, num, num))

	f = functions.ThreeArgFunction(func(envs.Environment, types.XValue, types.XValue, types.XValue) types.XValue { return result })
	test.AssertXEqual(t, xe("need 3 argument(s), got 2"), f(context.Background(), env, num, num))
	test.AssertXEqual(t, result, f(context.Background(), env, num, num, num))
	test.AssertXEqual(t, xe("need 3 argument(s), got 4"), f(context.Background(), env, num, num, num, num))

	f = functions.OneTextFunction(func(envs.Environment, *types.XText) types.XValue { return result })
	test.AssertXEqual(t, xe("need 1 argument(s), got 0"), f(context.Background(), env))
	test.AssertXEqual(t, result, f(context.Background(), env, text))
	test.AssertXEqual(t, result, f(context.Background(), env, num))
	test.AssertXEqual(t, xe("error"), f(context.Background(), env, xe("error")))
	test.AssertXEqual(t, xe("need 1 argument(s), got 2"), f(context.Background(), env, text, text))

	f = functions.TwoTextFunction(func(envs.Environment, *types.XText, *types.XText) types.XValue { return result })
	test.AssertXEqual(t, xe("need 2 argument(s), got 1"), f(context.Background(), env, text))
	test.AssertXEqual(t, result, f(context.Background(), env, text, text))
	test.AssertXEqual(t, result, f(context.Background(), env, num, num))
	test.AssertXEqual(t, xe("error"), f(context.Background(), env, xe("error"), text))
	test.AssertXEqual(t, xe("error"), f(context.Background(), env, text, xe("error")))
	test.AssertXEqual(t, xe("need 2 argument(s), got 3"), f(context.Background(), env, text, text, text))

	f = functions.OneNumberFunction(func(envs.Environment, *types.XNumber) types.XValue { return result })
	test.AssertXEqual(t, xe("need 1 argument(s), got 0"), f(context.Background(), env))
	test.AssertXEqual(t, result, f(context.Background(), env, num))
	test.AssertXEqual(t, result, f(context.Background(), env, types.NewXText("1")))
	test.AssertXEqual(t, xe(`unable to convert "X" to a number`), f(context.Background(), env, text))
	test.AssertXEqual(t, xe("need 1 argument(s), got 2"), f(context.Background(), env, num, num))

	f = functions.TwoNumberFunction(func(envs.Environment, *types.XNumber, *types.XNumber) types.XValue { return result })
	test.AssertXEqual(t, xe("need 2 argument(s), got 1"), f(context.Background(), env, num))
	test.AssertXEqual(t, result, f(context.Background(), env, num, num))
	test.AssertXEqual(t, result, f(context.Background(), env, types.NewXText("1"), types.NewXText("2")))
	test.AssertXEqual(t, xe(`unable to convert "X" to a number`), f(context.Background(), env, types.NewXText("X"), num))
	test.AssertXEqual(t, xe(`unable to convert "X" to a number`), f(context.Background(), env, num, types.NewXText("X")))
	test.AssertXEqual(t, xe("need 2 argument(s), got 3"), f(context.Background(), env, num, num, num))

	f = functions.TextAndNumberFunction(func(envs.Environment, *types.XText, *types.XNumber) types.XValue { return result })
	test.AssertXEqual(t, xe("need 2 argument(s), got 1"), f(context.Background(), env, text))
	test.AssertXEqual(t, result, f(context.Background(), env, text, num))
	test.AssertXEqual(t, result, f(context.Background(), env, num, num))
	test.AssertXEqual(t, result, f(context.Background(), env, text, types.NewXText("2")))
	test.AssertXEqual(t, xe("error"), f(context.Background(), env, xe("error"), num))
	test.AssertXEqual(t, xe(`unable to convert "X" to a number`), f(context.Background(), env, text, types.NewXText("X")))

	f = functions.ObjectTextAndNumberFunction(func(envs.Environment, *types.XObject, *types.XText, *types.XNumber) types.XValue { return result })
	test.AssertXEqual(t, xe("need 3 argument(s), got 2"), f(context.Background(), env, obj, text))
	test.AssertXEqual(t, result, f(context.Background(), env, obj, text, num))
	test.AssertXEqual(t, xe("unable to convert 1 to an object"), f(context.Background(), env, num, text, num))
	test.AssertXEqual(t, xe("error"), f(context.Background(), env, obj, xe("error"), num))
	test.AssertXEqual(t, xe(`unable to convert "X" to a number`), f(context.Background(), env, obj, text, text))

	f = functions.ObjectAndTextsFunction(func(envs.Environment, *types.XObject, ...*types.XText) types.XValue { return result })
	test.AssertXEqual(t, xe("need at least 2 argument(s), got 1"), f(context.Background(), env, obj))
	test.AssertXEqual(t, result, f(context.Background(), env, obj, text, text, text))
	test.AssertXEqual(t, xe("unable to convert 1 to an object"), f(context.Background(), env, num, text))
	test.AssertXEqual(t, xe("error"), f(context.Background(), env, obj, xe("error")))

	f = functions.OneObjectFunction(func(envs.Environment, *types.XObject) types.XValue { return result })
	test.AssertXEqual(t, xe("need 1 argument(s), got 0"), f(context.Background(), env))
	test.AssertXEqual(t, result, f(context.Background(), env, obj))
	test.AssertXEqual(t, xe("unable to convert 1 to an object"), f(context.Background(), env, num))
	test.AssertXEqual(t, xe("error"), f(context.Background(), env, xe("error")))
}
