package static_test

import (
	"testing"

	"github.com/nyaruka/gocommon/i18n"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/static"
	"github.com/nyaruka/goflow/utils"

	"github.com/stretchr/testify/assert"
)

func TestChannel(t *testing.T) {
	channel := static.NewChannel(
		assets.ChannelUUID("ffffffff-9b24-92e1-ffff-ffffb207cdb4"),
		"Android",
		"+234151",
		[]string{"tel"},
		[]assets.ChannelRole{assets.ChannelRoleSend},
		[]assets.ChannelFeature{assets.ChannelFeatureOptIns},
	)
	assert.Equal(t, assets.ChannelUUID("ffffffff-9b24-92e1-ffff-ffffb207cdb4"), channel.UUID())
	assert.Equal(t, "Android", channel.Name())
	assert.Equal(t, "+234151", channel.Address())
	assert.Equal(t, []string{"tel"}, channel.Schemes())
	assert.Equal(t, []assets.ChannelRole{assets.ChannelRoleSend}, channel.Roles())
	assert.Equal(t, []assets.ChannelFeature{assets.ChannelFeatureOptIns}, channel.Features())
	assert.Equal(t, i18n.NilCountry, channel.Country())
	assert.Nil(t, channel.MatchPrefixes())
	assert.True(t, channel.AllowInternational())

	// check that UUIDs aren't required to be valid UUID4s
	assert.Nil(t, utils.Validate(channel))

	channel = static.NewTelChannel(
		assets.ChannelUUID("ffffffff-9b24-92e1-ffff-ffffb207cdb4"),
		"Android",
		"+234151",
		[]assets.ChannelRole{assets.ChannelRoleSend},
		nil,
		"RW",
		[]string{"+25079"},
		false,
	)

	assert.Equal(t, i18n.Country("RW"), channel.Country())
	assert.Equal(t, []string{"+25079"}, channel.MatchPrefixes())
	assert.False(t, channel.AllowInternational())
}
