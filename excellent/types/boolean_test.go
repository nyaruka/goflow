package types_test

import (
	"encoding/json"
	"testing"

	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/utils"

	"github.com/stretchr/testify/assert"
)

func TestXBoolean(t *testing.T) {
	env := utils.NewEnvironmentBuilder().Build()

	// test equality
	assert.True(t, types.XBooleanFalse.Equals(types.XBooleanFalse))
	assert.True(t, types.XBooleanTrue.Equals(types.XBooleanTrue))
	assert.False(t, types.XBooleanTrue.Equals(types.XBooleanFalse))

	// test comparison
	assert.Equal(t, 0, types.XBooleanFalse.Compare(types.XBooleanFalse))
	assert.Equal(t, 1, types.XBooleanTrue.Compare(types.XBooleanFalse))
	assert.Equal(t, -1, types.XBooleanFalse.Compare(types.XBooleanTrue))

	// test stringify
	assert.Equal(t, "false", types.XBooleanFalse.String())
	assert.Equal(t, "true", types.XBooleanTrue.String())

	assert.Equal(t, "true", types.XBooleanTrue.Describe(env))
	assert.Equal(t, "false", types.XBooleanFalse.Describe(env))

	// unmarshal
	var val types.XBoolean
	err := json.Unmarshal([]byte(`true`), &val)
	assert.NoError(t, err)
	assert.Equal(t, types.XBooleanTrue, val)
}
