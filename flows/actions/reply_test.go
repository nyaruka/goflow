package actions

import (
	"testing"

	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"

	"github.com/stretchr/testify/assert"
)

func TestGetChannelForURN(t *testing.T) {
	rolesSend := []flows.ChannelRole{flows.ChannelRoleSend}
	rolesDefault := []flows.ChannelRole{flows.ChannelRoleSend, flows.ChannelRoleReceive}

	android := flows.NewChannel(flows.ChannelUUID(utils.NewUUID()), "Android", "+250961111111", []string{"tel"}, rolesDefault)
	twitter := flows.NewChannel(flows.ChannelUUID(utils.NewUUID()), "Twitter", "nyaruka", []string{"twitter", "twitterid"}, rolesDefault)
	nexmo := flows.NewChannel(flows.ChannelUUID(utils.NewUUID()), "Nexmo", "+250961111111", []string{"tel"}, rolesSend)

	// no channel
	assert.Nil(t, getChannelForURN([]flows.Channel{}, flows.NewContactURN(urns.URN("tel:+12345678999"), nil)))

	// no channel with correct scheme
	assert.Nil(t, getChannelForURN([]flows.Channel{}, flows.NewContactURN(urns.URN("twitter:rowan"), nil)))

	// first channel that supports scheme
	assert.Equal(t, getChannelForURN([]flows.Channel{twitter, android}, flows.NewContactURN(urns.URN("tel:+250962222222"), nil)), android)

	// explicit channel with URN
	assert.Equal(t, getChannelForURN([]flows.Channel{twitter, android, nexmo}, flows.NewContactURN(urns.URN("tel:+250962222222"), nexmo)), nexmo)
}
