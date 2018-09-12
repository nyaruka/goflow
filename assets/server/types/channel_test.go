package types_test

import (
	"testing"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/server/types"

	"github.com/stretchr/testify/assert"
)

func TestReadChannel(t *testing.T) {
	// check that UUIDs aren't required to be valid UUID4s
	channels, err := types.ReadChannels([]byte(`[{"uuid": "ffffffff-9b24-92e1-ffff-ffffb207cdb4", "name": "Old Channel", "schemes": ["tel"], "roles": ["send"]}]`))
	assert.NoError(t, err)
	assert.Equal(t, 1, len(channels))
	assert.Equal(t, assets.ChannelUUID("ffffffff-9b24-92e1-ffff-ffffb207cdb4"), channels[0].UUID())
	assert.Equal(t, "Old Channel", channels[0].Name())
	assert.Equal(t, []assets.ChannelRole{assets.ChannelRoleSend}, channels[0].Roles())
}
