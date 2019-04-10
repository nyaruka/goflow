package types_test

import (
	"testing"

	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/utils"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestXError(t *testing.T) {
	env := utils.NewEnvironmentBuilder().Build()

	err1 := types.NewXError(errors.Errorf("I failed"))
	assert.Equal(t, types.NewXText("I failed"), err1.ToXText(env))
	assert.Equal(t, types.XBooleanFalse, err1.ToXBoolean())
	assert.Equal(t, `XError("I failed")`, err1.String())
	assert.Equal(t, "I failed", err1.Error())

	asJSON, _ := types.ToXJSON(err1)
	assert.Equal(t, types.NewXText(""), asJSON)
}
