package types_test

import (
	"encoding/json"
	"testing"

	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/utils"

	"github.com/stretchr/testify/assert"
)

func TestXFunction(t *testing.T) {
	env := utils.NewEnvironmentBuilder().Build()

	func1 := types.XFunction(func(env utils.Environment, args ...types.XValue) types.XValue { return nil })
	func2 := types.XFunction(func(env utils.Environment, args ...types.XValue) types.XValue { return nil })

	assert.True(t, func1.Truthy())
	assert.Equal(t, `function`, func1.Render())
	assert.Equal(t, `function`, func1.Format(env))
	assert.Equal(t, `XFunction`, func1.String())
	assert.Equal(t, `function`, func1.Describe())

	marshaled, err := json.Marshal(func1)
	assert.NoError(t, err)
	assert.Equal(t, `null`, string(marshaled))

	assert.True(t, types.Equals(func1, func2))
}
