package tests_test

import (
	"testing"

	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows/routers/tests"
	"github.com/nyaruka/goflow/utils"

	"github.com/stretchr/testify/assert"
)

func TestTestResult(t *testing.T) {
	env := utils.NewDefaultEnvironment()

	res := tests.NewTrueResult(types.NewXText("abc"))

	assert.True(t, res.Matched())
	assert.Equal(t, types.NewXText("abc"), res.Match())

	assert.Equal(t, "test result", res.Describe())
	assert.Equal(t, types.XBooleanTrue, res.Resolve(env, "matched"))
	assert.Equal(t, types.NewXText("abc"), res.Resolve(env, "match"))
	assert.Equal(t, types.NewXResolveError(res, "xxx"), res.Resolve(env, "xxx"))
	assert.Equal(t, types.XBooleanTrue, res.Reduce(env))
	assert.Equal(t, types.NewXText(`{"match":"abc","matched":true}`), res.ToXJSON(env))
}
