package flows_test

import (
	"testing"

	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
	"github.com/stretchr/testify/assert"
)

func TestChannel(t *testing.T) {
	env := utils.NewDefaultEnvironment()

	rolesDefault := []flows.ChannelRole{flows.ChannelRoleSend, flows.ChannelRoleReceive}
	uuid := flows.ChannelUUID("821fe776-b97d-4046-b1dc-a7d9d3b3b9c7")

	ch := flows.NewChannel(uuid, "Android", "+250961111111", []string{"tel"}, rolesDefault, nil)

	assert.Equal(t, uuid, ch.UUID())
	assert.Equal(t, "Android", ch.Name())
	assert.Equal(t, []string{"tel"}, ch.Schemes())
	assert.Equal(t, "+250961111111", ch.Address())
	assert.Equal(t, "channel", ch.Describe())

	assert.Equal(t, types.NewXText(string(uuid)), ch.Resolve(env, "uuid"))
	assert.Equal(t, types.NewXText("Android"), ch.Resolve(env, "name"))
	assert.Equal(t, types.NewXText("+250961111111"), ch.Resolve(env, "address"))
	assert.Equal(t, types.NewXResolveError(ch, "xxx"), ch.Resolve(env, "xxx"))
	assert.Equal(t, types.NewXText("Android"), ch.Reduce(env))
	assert.Equal(t, types.NewXText(`{"address":"+250961111111","name":"Android","uuid":"821fe776-b97d-4046-b1dc-a7d9d3b3b9c7"}`), ch.ToXJSON(env))

	assert.Equal(t, flows.NewChannelReference(uuid, "Android"), ch.Reference())
	assert.True(t, ch.HasRole(flows.ChannelRoleSend))
	assert.False(t, ch.HasRole(flows.ChannelRoleCall))
}

func TestChannelSetGetForURN(t *testing.T) {
	rolesSend := []flows.ChannelRole{flows.ChannelRoleSend}
	rolesDefault := []flows.ChannelRole{flows.ChannelRoleSend, flows.ChannelRoleReceive}

	claro := flows.NewChannel(flows.ChannelUUID(utils.NewUUID()), "Claro", "+593971111111", []string{"tel"}, rolesDefault, nil)
	mtn := flows.NewChannel(flows.ChannelUUID(utils.NewUUID()), "MTN", "+250782222222", []string{"tel"}, rolesDefault, nil)
	tigo := flows.NewChannel(flows.ChannelUUID(utils.NewUUID()), "Tigo", "+250723333333", []string{"tel"}, rolesDefault, nil)
	twitter := flows.NewChannel(flows.ChannelUUID(utils.NewUUID()), "Twitter", "nyaruka", []string{"twitter", "twitterid"}, rolesDefault, nil)

	claro.SetTelMatching("EC", nil)
	mtn.SetTelMatching("RW", nil)
	tigo.SetTelMatching("RW", nil)

	all := flows.NewChannelSet([]flows.Channel{claro, mtn, tigo, twitter})

	// nil if no channel
	emptySet := flows.NewChannelSet([]flows.Channel{})
	assert.Nil(t, emptySet.GetForURN(flows.NewContactURN(urns.URN("tel:+12345678999"), nil), flows.ChannelRoleSend))

	// nil if no channel with correct scheme
	assert.Nil(t, all.GetForURN(flows.NewContactURN(urns.URN("mailto:rowan@foo.bar"), nil), flows.ChannelRoleSend))

	// if URN has a preferred channel, that is always used
	assert.Equal(t, tigo, all.GetForURN(flows.NewContactURN(urns.URN("tel:+250962222222"), tigo), flows.ChannelRoleSend))

	// if there's only one channel for that scheme, it's used
	assert.Equal(t, twitter, all.GetForURN(flows.NewContactURN(urns.URN("twitter:nyaruka2"), nil), flows.ChannelRoleSend))

	// if there's only one channel for that country, it's used
	assert.Equal(t, claro, all.GetForURN(flows.NewContactURN(urns.URN("tel:+593971234567"), nil), flows.ChannelRoleSend))

	// if there's multiple channels, one with longest number overlap wins
	assert.Equal(t, mtn, all.GetForURN(flows.NewContactURN(urns.URN("tel:+250781234567"), nil), flows.ChannelRoleSend))
	assert.Equal(t, tigo, all.GetForURN(flows.NewContactURN(urns.URN("tel:+250721234567"), nil), flows.ChannelRoleSend))

	// if there's no overlap, then last/newest channel wins
	assert.Equal(t, tigo, all.GetForURN(flows.NewContactURN(urns.URN("tel:+250962222222"), nil), flows.ChannelRoleSend))

	// channels can be delegates for other channels
	android := flows.NewChannel(flows.ChannelUUID(utils.NewUUID()), "Android", "+250723333333", []string{"tel"}, rolesDefault, nil)
	bulk := flows.NewChannel(flows.ChannelUUID(utils.NewUUID()), "Bulk Sender", "1234", []string{"tel"}, rolesSend, android.Reference())
	all = flows.NewChannelSet([]flows.Channel{android, bulk})

	// delegate will always be used if it has the requested role
	assert.Equal(t, android, all.GetForURN(flows.NewContactURN(urns.URN("tel:+250721234567"), nil), flows.ChannelRoleReceive))
	assert.Equal(t, bulk, all.GetForURN(flows.NewContactURN(urns.URN("tel:+250721234567"), nil), flows.ChannelRoleSend))

	// matching prefixes can be explicitly set too
	short1 := flows.NewChannel(flows.ChannelUUID(utils.NewUUID()), "Shortcode 1", "1234", []string{"tel"}, rolesSend, nil)
	short1.SetTelMatching("RW", []string{"25078", "25077"})
	short2 := flows.NewChannel(flows.ChannelUUID(utils.NewUUID()), "Shortcode 2", "1235", []string{"tel"}, rolesSend, nil)
	short2.SetTelMatching("RW", []string{"25072"})
	all = flows.NewChannelSet([]flows.Channel{short1, short2})

	assert.Equal(t, short1, all.GetForURN(flows.NewContactURN(urns.URN("tel:+250781234567"), nil), flows.ChannelRoleSend))
	assert.Equal(t, short1, all.GetForURN(flows.NewContactURN(urns.URN("tel:+250771234567"), nil), flows.ChannelRoleSend))
	assert.Equal(t, short2, all.GetForURN(flows.NewContactURN(urns.URN("tel:+250721234567"), nil), flows.ChannelRoleSend))
}

func TestChannelUnmarsal(t *testing.T) {
	// check that UUIDs aren't required to be valid UUID4s
	channel, err := flows.ReadChannel([]byte(`{"uuid": "ffffffff-9b24-92e1-ffff-ffffb207cdb4", "name": "Old Channel", "schemes": ["tel"], "roles": ["send"]}`))
	assert.NoError(t, err)
	assert.Equal(t, flows.ChannelUUID("ffffffff-9b24-92e1-ffff-ffffb207cdb4"), channel.UUID())
	assert.Equal(t, "Old Channel", channel.Name())
	assert.Equal(t, []flows.ChannelRole{flows.ChannelRoleSend}, channel.Roles())
}
