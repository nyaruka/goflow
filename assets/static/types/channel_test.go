package types_test

import (
	"testing"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/static/types"
	"github.com/nyaruka/goflow/utils"

	"github.com/stretchr/testify/assert"
)

func TestChannel(t *testing.T) {
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
	assert.Nil(t, utils.Validate(channel))
}
