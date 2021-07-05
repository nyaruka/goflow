package flows_test

import (
	"testing"

	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/test"

	"github.com/stretchr/testify/assert"
)

func TestChannel(t *testing.T) {
	env := envs.NewBuilder().Build()

	uuids.SetGenerator(uuids.NewSeededGenerator(1234))
	defer uuids.SetGenerator(uuids.DefaultGenerator)

	rolesDefault := []assets.ChannelRole{assets.ChannelRoleSend, assets.ChannelRoleReceive}
	ch := test.NewChannel("Android", "+250961111111", []string{"tel"}, rolesDefault, nil)

	assert.Equal(t, assets.ChannelUUID("c00e5d67-c275-4389-aded-7d8b151cbd5b"), ch.UUID())
	assert.Equal(t, "Android", ch.Name())
	assert.Equal(t, []string{"tel"}, ch.Schemes())
	assert.Equal(t, "+250961111111", ch.Address())
	assert.Equal(t, "+250961111111 (Android)", ch.String())

	test.AssertXEqual(t, types.NewXObject(map[string]types.XValue{
		"__default__": types.NewXText("Android"),
		"uuid":        types.NewXText(string(ch.UUID())),
		"name":        types.NewXText("Android"),
		"address":     types.NewXText("+250961111111"),
	}), flows.Context(env, ch))

	assert.Equal(t, assets.NewChannelReference(ch.UUID(), "Android"), ch.Reference())
	assert.True(t, ch.HasRole(assets.ChannelRoleSend))
	assert.False(t, ch.HasRole(assets.ChannelRoleCall))

	// nil object returns nil reference
	assert.Nil(t, (*flows.Channel)(nil).Reference())
}

func TestChannelSetGetForURN(t *testing.T) {
	rolesSend := []assets.ChannelRole{assets.ChannelRoleSend}
	rolesReceive := []assets.ChannelRole{assets.ChannelRoleReceive}
	rolesDefault := []assets.ChannelRole{assets.ChannelRoleSend, assets.ChannelRoleReceive}

	claro := test.NewTelChannel("Claro", "+593971111111", rolesDefault, nil, "EC", nil, true)
	mtn := test.NewTelChannel("MTN", "+250782222222", rolesDefault, nil, "RW", nil, false)
	tigo := test.NewTelChannel("Tigo", "+250723333333", rolesDefault, nil, "RW", nil, false)
	twilio := test.NewTelChannel("Twilio", "+17036975131", rolesDefault, nil, "", nil, false)
	twitter := test.NewChannel("Twitter", "nyaruka", []string{"twitter", "twitterid"}, rolesDefault, nil)
	receiver := test.NewTelChannel("Receiver", "+250724444444", rolesReceive, nil, "RW", nil, false)

	all := flows.NewChannelAssets([]assets.Channel{claro.Asset(), mtn.Asset(), tigo.Asset(), twitter.Asset()})
	rwOnly := flows.NewChannelAssets([]assets.Channel{mtn.Asset(), tigo.Asset()})
	twOnly := flows.NewChannelAssets([]assets.Channel{twilio.Asset()})
	receiverSet := flows.NewChannelAssets([]assets.Channel{receiver.Asset()})

	// nil if no channel
	emptySet := flows.NewChannelAssets(nil)
	assert.Nil(t, emptySet.GetForURN(flows.NewContactURN(urns.URN("tel:+12345678999"), nil), assets.ChannelRoleSend))

	// nil if not channel with correct role
	assert.Nil(t, receiverSet.GetForURN(flows.NewContactURN(urns.URN("tel:+12345678999"), receiver), assets.ChannelRoleSend))

	// can still match URN has a preferred channel with matching role
	assert.Equal(t, receiver, receiverSet.GetForURN(flows.NewContactURN(urns.URN("tel:+12345678999"), receiver), assets.ChannelRoleReceive))

	// nil if no channel with correct scheme
	assert.Nil(t, all.GetForURN(flows.NewContactURN(urns.URN("mailto:rowan@foo.bar"), nil), assets.ChannelRoleSend))

	// if URN has a preferred channel, that is always used
	assert.Equal(t, tigo, all.GetForURN(flows.NewContactURN(urns.URN("tel:+250962222222"), tigo), assets.ChannelRoleSend))

	// if there's only one channel for that scheme, it's used
	assert.Equal(t, twitter, all.GetForURN(flows.NewContactURN(urns.URN("twitter:nyaruka2"), nil), assets.ChannelRoleSend))

	// if there's only one channel for that country, it's used
	assert.Equal(t, claro, all.GetForURN(flows.NewContactURN(urns.URN("tel:+593971234567"), nil), assets.ChannelRoleSend))

	// return nil for international send if the channels don't allow it
	assert.Nil(t, rwOnly.GetForURN(flows.NewContactURN(urns.URN("tel:+593971234567"), nil), assets.ChannelRoleSend))

	// but use them if they do
	assert.Equal(t, claro, all.GetForURN(flows.NewContactURN(urns.URN("tel:+57971234567"), nil), assets.ChannelRoleSend))

	// or if they're implicitly international by having no country
	assert.Equal(t, twilio, twOnly.GetForURN(flows.NewContactURN(urns.URN("tel:+57971234567"), nil), assets.ChannelRoleSend))

	// if there's multiple channels, one with longest number overlap wins
	assert.Equal(t, mtn, all.GetForURN(flows.NewContactURN(urns.URN("tel:+250781234567"), nil), assets.ChannelRoleSend))
	assert.Equal(t, tigo, all.GetForURN(flows.NewContactURN(urns.URN("tel:+250721234567"), nil), assets.ChannelRoleSend))

	// if there's no overlap, then last/newest channel wins
	assert.Equal(t, tigo, all.GetForURN(flows.NewContactURN(urns.URN("tel:+250962222222"), nil), assets.ChannelRoleSend))

	// channels can be delegates for other channels
	android := test.NewChannel("Android", "+250723333333", []string{"tel"}, rolesDefault, nil)
	bulk := test.NewChannel("Bulk Sender", "1234", []string{"tel"}, rolesSend, android.Reference())
	all = flows.NewChannelAssets([]assets.Channel{android.Asset(), bulk.Asset()})

	// delegate will always be used if it has the requested role
	assert.Equal(t, android, all.GetForURN(flows.NewContactURN(urns.URN("tel:+250721234567"), nil), assets.ChannelRoleReceive))
	assert.Equal(t, bulk, all.GetForURN(flows.NewContactURN(urns.URN("tel:+250721234567"), nil), assets.ChannelRoleSend))

	// matching prefixes can be explicitly set too
	short1 := test.NewTelChannel("Shortcode 1", "1234", rolesSend, nil, "RW", []string{"25078", "25077"}, false)
	short2 := test.NewTelChannel("Shortcode 2", "1235", rolesSend, nil, "RW", []string{"25072"}, false)
	all = flows.NewChannelAssets([]assets.Channel{short1.Asset(), short2.Asset()})

	assert.Equal(t, short1, all.GetForURN(flows.NewContactURN(urns.URN("tel:+250781234567"), nil), assets.ChannelRoleSend))
	assert.Equal(t, short1, all.GetForURN(flows.NewContactURN(urns.URN("tel:+250771234567"), nil), assets.ChannelRoleSend))
	assert.Equal(t, short2, all.GetForURN(flows.NewContactURN(urns.URN("tel:+250721234567"), nil), assets.ChannelRoleSend))
}
