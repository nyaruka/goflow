package core_test

import (
	"testing"
	"time"

	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/core"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/test"
	"github.com/stretchr/testify/assert"
)

func TestChannel(t *testing.T) {
	env := envs.NewBuilder().Build()

	uuids.SetGenerator(uuids.NewSeededGenerator(1234, time.Now))
	defer uuids.SetGenerator(uuids.DefaultGenerator)

	rolesDefault := []assets.ChannelRole{assets.ChannelRoleSend, assets.ChannelRoleReceive}
	ch := test.NewChannel("Android", "+250961111111", []string{"tel"}, rolesDefault, nil)

	assert.Equal(t, assets.ChannelUUID("15a2ee5e-5e45-4711-8e0f-6b2abe4360d8"), ch.UUID())
	assert.Equal(t, "Android", ch.Name())
	assert.Equal(t, []string{"tel"}, ch.Schemes())
	assert.Equal(t, "+250961111111", ch.Address())
	assert.Equal(t, "+250961111111 (Android)", ch.String())

	test.AssertXEqual(t, types.NewXObject(map[string]types.XValue{
		"__default__": types.NewXText("Android"),
		"uuid":        types.NewXText(string(ch.UUID())),
		"name":        types.NewXText("Android"),
		"address":     types.NewXText("+250961111111"),
	}), core.Context(env, ch))

	assert.Equal(t, assets.NewChannelReference(ch.UUID(), "Android"), ch.Reference())
	assert.True(t, ch.HasRole(assets.ChannelRoleSend))
	assert.False(t, ch.HasRole(assets.ChannelRoleCall))
	assert.False(t, ch.HasFeature(assets.ChannelFeatureOptIns))

	// nil object returns nil reference
	assert.Nil(t, (*core.Channel)(nil).Reference())
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

	all := core.NewChannelAssets([]assets.Channel{claro.Asset(), mtn.Asset(), tigo.Asset(), twitter.Asset()})
	rwOnly := core.NewChannelAssets([]assets.Channel{mtn.Asset(), tigo.Asset()})
	twOnly := core.NewChannelAssets([]assets.Channel{twilio.Asset()})
	receiverSet := core.NewChannelAssets([]assets.Channel{receiver.Asset()})

	// nil if no channel
	emptySet := core.NewChannelAssets(nil)
	assert.Nil(t, emptySet.GetForURN(core.NewURN("tel", "+12345678999", "", nil), assets.ChannelRoleSend))

	// nil if not channel with correct role
	assert.Nil(t, receiverSet.GetForURN(core.NewURN("tel", "+12345678999", "", receiver), assets.ChannelRoleSend))

	// can still match URN has a preferred channel with matching role
	assert.Equal(t, receiver, receiverSet.GetForURN(core.NewURN("tel", "+12345678999", "", receiver), assets.ChannelRoleReceive))

	// nil if no channel with correct scheme
	assert.Nil(t, all.GetForURN(core.NewURN("mailto", "rowan@foo.bar", "", nil), assets.ChannelRoleSend))

	// nil if URN has preferred channel but that channel doesn't support the URN's scheme (e.g. bsuid URN with whatsapp channel affinity)
	whatsapp := test.NewChannel("WhatsApp", "+250788000000", []string{"whatsapp"}, rolesDefault, nil)
	waOnly := core.NewChannelAssets([]assets.Channel{whatsapp.Asset()})
	assert.Nil(t, waOnly.GetForURN(core.NewURN("bsuid", "abc123", "", whatsapp), assets.ChannelRoleSend))

	// if URN has a preferred channel with the required role and matching scheme, that is used
	assert.Equal(t, tigo, all.GetForURN(core.NewURN("tel", "+250962222222", "", tigo), assets.ChannelRoleSend))

	// if there's only one channel for that scheme, it's used
	assert.Equal(t, twitter, all.GetForURN(core.NewURN("twitter", "nyaruka2", "", nil), assets.ChannelRoleSend))

	// if there's only one channel for that country, it's used
	assert.Equal(t, claro, all.GetForURN(core.NewURN("tel", "+593971234567", "", nil), assets.ChannelRoleSend))

	// return nil for international send if the channels don't allow it
	assert.Nil(t, rwOnly.GetForURN(core.NewURN("tel", "+593971234567", "", nil), assets.ChannelRoleSend))

	// but use them if they do
	assert.Equal(t, claro, all.GetForURN(core.NewURN("tel", "+57971234567", "", nil), assets.ChannelRoleSend))

	// or if they're implicitly international by having no country
	assert.Equal(t, twilio, twOnly.GetForURN(core.NewURN("tel", "+57971234567", "", nil), assets.ChannelRoleSend))

	// if there's multiple channels, one with longest number overlap wins
	assert.Equal(t, mtn, all.GetForURN(core.NewURN("tel", "+250781234567", "", nil), assets.ChannelRoleSend))
	assert.Equal(t, tigo, all.GetForURN(core.NewURN("tel", "+250721234567", "", nil), assets.ChannelRoleSend))

	// if there's no overlap, then last/newest channel wins
	assert.Equal(t, tigo, all.GetForURN(core.NewURN("tel", "+250962222222", "", nil), assets.ChannelRoleSend))

	// matching prefixes can be explicitly set too
	short1 := test.NewTelChannel("Shortcode 1", "1234", rolesSend, nil, "RW", []string{"25078", "25077"}, false)
	short2 := test.NewTelChannel("Shortcode 2", "1235", rolesSend, nil, "RW", []string{"25072"}, false)
	all = core.NewChannelAssets([]assets.Channel{short1.Asset(), short2.Asset()})

	assert.Equal(t, short1, all.GetForURN(core.NewURN("tel", "+250781234567", "", nil), assets.ChannelRoleSend))
	assert.Equal(t, short1, all.GetForURN(core.NewURN("tel", "+250771234567", "", nil), assets.ChannelRoleSend))
	assert.Equal(t, short2, all.GetForURN(core.NewURN("tel", "+250721234567", "", nil), assets.ChannelRoleSend))
}
