package types_test

import (
	"testing"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestXError(t *testing.T) {
	env := envs.NewBuilder().Build()

	err1 := types.NewXError(errors.Errorf("I failed"))
	assert.Equal(t, "error", err1.Describe())
	assert.Equal(t, "I failed", err1.Render())
	assert.Equal(t, "", err1.Format(env))
	assert.False(t, err1.Truthy())
	assert.Equal(t, `XError("I failed")`, err1.String())
	assert.Equal(t, "I failed", err1.Error())

	asJSON, _ := types.ToXJSON(err1)
	assert.Equal(t, types.NewXText(""), asJSON)

	marshaled, err := jsonx.Marshal(err1)
	assert.NoError(t, err)
	assert.Equal(t, `null`, string(marshaled))
}
