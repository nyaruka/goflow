package types_test

import (
	"testing"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"

	"github.com/stretchr/testify/assert"
)

func TestXBoolean(t *testing.T) {
	env := envs.NewBuilder().Build()

	// test equality
	assert.True(t, types.XBooleanFalse.Equals(types.XBooleanFalse))
	assert.True(t, types.XBooleanTrue.Equals(types.XBooleanTrue))
	assert.False(t, types.XBooleanTrue.Equals(types.XBooleanFalse))

	// test comparison
	assert.Equal(t, 0, types.XBooleanFalse.Compare(types.XBooleanFalse))
	assert.Equal(t, 1, types.XBooleanTrue.Compare(types.XBooleanFalse))
	assert.Equal(t, -1, types.XBooleanFalse.Compare(types.XBooleanTrue))

	// test to text
	assert.Equal(t, "false", types.XBooleanFalse.Render())
	assert.Equal(t, "true", types.XBooleanTrue.Render())
	assert.Equal(t, "false", types.XBooleanFalse.Format(env))
	assert.Equal(t, "true", types.XBooleanTrue.Format(env))

	// test truthniess
	assert.False(t, types.XBooleanFalse.Truthy())
	assert.True(t, types.XBooleanTrue.Truthy())

	// test to JSON
	asJSON, _ := types.ToXJSON(types.XBooleanFalse)
	assert.Equal(t, types.NewXText("false"), asJSON)
	asJSON, _ = types.ToXJSON(types.XBooleanTrue)
	assert.Equal(t, types.NewXText("true"), asJSON)

	// test stringify
	assert.Equal(t, "XBoolean(false)", types.XBooleanFalse.String())
	assert.Equal(t, "XBoolean(true)", types.XBooleanTrue.String())

	assert.Equal(t, "true", types.XBooleanTrue.Describe())
	assert.Equal(t, "false", types.XBooleanFalse.Describe())

	// unmarshal
	var val types.XBoolean
	err := jsonx.Unmarshal([]byte(`true`), &val)
	assert.NoError(t, err)
	assert.Equal(t, types.XBooleanTrue, val)
}
