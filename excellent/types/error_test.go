package types_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/stretchr/testify/assert"
)

func TestXError(t *testing.T) {
	env := envs.NewBuilder().Build()

	err1 := types.NewXError(fmt.Errorf("I failed"))
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

func TestErrTooComplex(t *testing.T) {
	// the sentinel matches itself
	assert.True(t, errors.Is(types.ErrTooComplex, types.ErrTooComplex))

	// and still matches when re-wrapped in a new XError (identity would miss this)
	rewrapped := types.NewXError(types.ErrTooComplex.Native())
	assert.NotSame(t, types.ErrTooComplex, rewrapped)
	assert.True(t, errors.Is(rewrapped, types.ErrTooComplex))

	// an unrelated error doesn't match
	assert.False(t, errors.Is(types.NewXErrorf("boom"), types.ErrTooComplex))
}
