package types_test

import (
	"testing"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/static/types"

	"github.com/stretchr/testify/assert"
)

func TestChannels(t *testing.T) {
	channel := types.NewChannel(
		assets.ChannelUUID("ffffffff-9b24-92e1-ffff-ffffb207cdb4"),
		"Android",
		"+234151",
		[]string{"tel"},
		[]assets.ChannelRole{assets.ChannelRoleSend},
		nil,
	)
	assert.Equal(t, assets.ChannelUUID("ffffffff-9b24-92e1-ffff-ffffb207cdb4"), channel.UUID())
	assert.Equal(t, "Android", channel.Name())
	assert.Equal(t, "+234151", channel.Address())
	assert.Equal(t, []string{"tel"}, channel.Schemes())
	assert.Equal(t, []assets.ChannelRole{assets.ChannelRoleSend}, channel.Roles())
	assert.Nil(t, channel.Parent())

	// check that UUIDs aren't required to be valid UUID4s
	channels, err := types.ReadChannels([]byte(`[{"uuid": "ffffffff-9b24-92e1-ffff-ffffb207cdb4", "name": "Old Channel", "schemes": ["tel"], "roles": ["send"], "country": "RW"}]`))
	assert.NoError(t, err)
	assert.Equal(t, 1, len(channels))
	assert.Equal(t, assets.ChannelUUID("ffffffff-9b24-92e1-ffff-ffffb207cdb4"), channels[0].UUID())
	assert.Equal(t, "Old Channel", channels[0].Name())
	assert.Equal(t, "", channels[0].Address())
	assert.Equal(t, []string{"tel"}, channels[0].Schemes())
	assert.Equal(t, []assets.ChannelRole{assets.ChannelRoleSend}, channels[0].Roles())
	assert.Nil(t, channels[0].Parent())
	assert.Equal(t, "RW", channels[0].Country())
}
