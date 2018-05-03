package flows_test

import (
	"testing"

	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"

	"github.com/stretchr/testify/assert"
)

func TestChannelSetGetForURN(t *testing.T) {
	rolesSend := []flows.ChannelRole{flows.ChannelRoleSend}
	rolesDefault := []flows.ChannelRole{flows.ChannelRoleSend, flows.ChannelRoleReceive}

	android := flows.NewChannel(flows.ChannelUUID(utils.NewUUID()), "Android", "+250961111111", []string{"tel"}, rolesDefault)
	twitter := flows.NewChannel(flows.ChannelUUID(utils.NewUUID()), "Twitter", "nyaruka", []string{"twitter", "twitterid"}, rolesDefault)
	nexmo := flows.NewChannel(flows.ChannelUUID(utils.NewUUID()), "Nexmo", "+250961111111", []string{"tel"}, rolesSend)

	emptySet := flows.NewChannelSet([]flows.Channel{})
	set := flows.NewChannelSet([]flows.Channel{android, twitter, nexmo})

	// no channel
	assert.Nil(t, emptySet.GetForURN(flows.NewContactURN(urns.URN("tel:+12345678999"), nil)))

	// no channel with correct scheme
	assert.Nil(t, set.GetForURN(flows.NewContactURN(urns.URN("mailto:rowan@foo.bar"), nil)))

	// first channel that supports scheme
	assert.Equal(t, set.GetForURN(flows.NewContactURN(urns.URN("tel:+250962222222"), nil)), android)

	// explicit channel with URN
	assert.Equal(t, set.GetForURN(flows.NewContactURN(urns.URN("tel:+250962222222"), nexmo)), nexmo)
}

func TestChannelUnmarsal(t *testing.T) {
	// check that UUIDs aren't required to be valid UUID4s
	channel, err := flows.ReadChannel([]byte(`{"uuid": "ffffffff-9b24-92e1-ffff-ffffb207cdb4", "name": "Old Channel", "schemes": ["tel"], "roles": ["send"]}`))
	assert.NoError(t, err)
	assert.Equal(t, flows.ChannelUUID("ffffffff-9b24-92e1-ffff-ffffb207cdb4"), channel.UUID())
	assert.Equal(t, "Old Channel", channel.Name())
	assert.Equal(t, []flows.ChannelRole{flows.ChannelRoleSend}, channel.Roles())
}
