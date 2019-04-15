package types_test

import (
	"testing"

	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/utils"

	"github.com/stretchr/testify/assert"
)

func TestXFunction(t *testing.T) {
	env := utils.NewEnvironmentBuilder().Build()

	func1 := types.XFunction(func(env utils.Environment, args ...types.XValue) types.XValue { return nil })

	assert.Equal(t, `function`, func1.Render(env))
	assert.Equal(t, `XFunction`, func1.String())
	assert.Equal(t, `function`, func1.Describe())
}
