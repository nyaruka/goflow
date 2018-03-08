package flows_test

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

	android := flows.NewChannel(flows.ChannelUUID(utils.NewUUID()), "Android", "+250961111111", []string{"tel"}, rolesDefault, nil)
	twitter := flows.NewChannel(flows.ChannelUUID(utils.NewUUID()), "Twitter", "nyaruka", []string{"twitter", "twitterid"}, rolesDefault, nil)

	// no channel
	assert.Nil(t, flows.GetChannelForURN([]flows.Channel{}, &flows.ContactURN{URN: urns.URN("tel:+12345678999")}))

	// no channel with correct scheme
	assert.Nil(t, flows.GetChannelForURN([]flows.Channel{}, &flows.ContactURN{URN: urns.URN("twitter:rowan")}))

	assert.Equal(t, flows.GetChannelForURN([]flows.Channel{twitter, android}, &flows.ContactURN{URN: urns.URN("tel:+250962222222")}), android)

	// add bulk sender channel
	nexmo := flows.NewChannel(flows.ChannelUUID(utils.NewUUID()), "Nexmo", "+250961111111", []string{"tel"}, rolesSend, android.Reference())

	assert.Equal(t, flows.GetChannelForURN([]flows.Channel{twitter, android, nexmo}, &flows.ContactURN{URN: urns.URN("tel:+250962222222")}), nexmo)
}
