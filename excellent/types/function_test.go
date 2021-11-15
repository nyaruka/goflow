package types_test

import (
	"testing"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/stretchr/testify/assert"
)

func TestXFunction(t *testing.T) {
	env := envs.NewBuilder().Build()

	func1 := types.NewXFunction("foo", func(env envs.Environment, args ...types.XValue) types.XValue {
		return types.NewXNumberFromInt(len(args))
	})
	func2 := types.NewXFunction("bad", func(env envs.Environment, args ...types.XValue) types.XValue { return types.NewXErrorf("boom") })
	anon1 := types.NewXFunction("", func(env envs.Environment, args ...types.XValue) types.XValue { return types.NewXText("c") })
	anon2 := types.NewXFunction("", func(env envs.Environment, args ...types.XValue) types.XValue { return types.NewXText("d") })

	assert.True(t, func1.Truthy())
	assert.Equal(t, `foo`, func1.Render())
	assert.Equal(t, `foo`, func1.Format(env))
	assert.Equal(t, `XFunction[foo]`, func1.String())
	assert.Equal(t, `foo(...)`, func1.Describe())
	assert.Equal(t, types.NewXNumberFromInt(0), func1.Call(env, nil))
	assert.Equal(t, types.NewXNumberFromInt(2), func1.Call(env, []types.XValue{types.NewXText("a"), types.NewXText("b")}))

	assert.Equal(t, types.NewXErrorf("error calling bad(...): boom"), func2.Call(env, nil))

	assert.True(t, anon1.Truthy())
	assert.Equal(t, `<anon>`, anon1.Render())
	assert.Equal(t, `<anon>`, anon1.Format(env))
	assert.Equal(t, `XFunction[<anon>]`, anon1.String())
	assert.Equal(t, `<anon>(...)`, anon1.Describe())
	assert.Equal(t, types.NewXText("c"), anon1.Call(env, nil))

	marshaled, err := jsonx.Marshal(func1)
	assert.NoError(t, err)
	assert.Equal(t, `null`, string(marshaled))

	assert.True(t, types.Equals(func1, func1))
	assert.False(t, types.Equals(func1, func2))

	// anonymous functions are never equal
	assert.False(t, types.Equals(anon1, anon1))
	assert.False(t, types.Equals(anon1, anon2))
}
