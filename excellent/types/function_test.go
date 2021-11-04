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

	func1 := types.NewXFunction("nill", func(env envs.Environment, args ...types.XValue) types.XValue { return nil })
	func2 := types.NewXFunction("foo", func(env envs.Environment, args ...types.XValue) types.XValue { return nil })

	assert.True(t, func1.Truthy())
	assert.Equal(t, `nill`, func1.Render())
	assert.Equal(t, `nill`, func1.Format(env))
	assert.Equal(t, `XFunction[nill]`, func1.String())
	assert.Equal(t, `nill(...)`, func1.Describe())

	marshaled, err := jsonx.Marshal(func1)
	assert.NoError(t, err)
	assert.Equal(t, `null`, string(marshaled))

	assert.True(t, types.Equals(func1, func1))
	assert.False(t, types.Equals(func1, func2))
}
