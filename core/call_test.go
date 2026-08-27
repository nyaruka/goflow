package core_test

import (
	"testing"

	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/static"
	"github.com/nyaruka/goflow/core"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows/engine"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCall(t *testing.T) {
	env := envs.NewBuilder().Build()

	source, err := static.NewSource([]byte(`{
		"channels": [
			{
				"uuid": "a78930fe-6a40-4aa8-99c3-e61b02f45ca1",
				"name": "Twilio Channel",
				"address": "+17036975133",
				"schemes": [
					"tel"
				],
				"roles": [
					"send",
					"receive",
					"call",
					"answer"
				]
			}
		]
	}`))
	require.NoError(t, err)

	sa, err := engine.NewSessionAssets(env, source, nil)
	require.NoError(t, err)

	twilio := sa.Channels().Get("a78930fe-6a40-4aa8-99c3-e61b02f45ca1")

	call := core.NewCall(
		"01978a2f-ad9a-7f2e-ad44-6e7547078cec",
		twilio,
		urns.URN("tel:+1234567890"),
	)

	// test marshaling our call
	ce := &core.CallEnvelope{
		UUID:    "01978a2f-ad9a-7f2e-ad44-6e7547078cec",
		Channel: assets.NewChannelReference("a78930fe-6a40-4aa8-99c3-e61b02f45ca1", "Twilio Channel"),
		URN:     urns.URN("tel:+1234567890"),
	}
	assert.Equal(t, ce, call.Marshal())

	// test unmarshaling
	call = ce.Unmarshal(sa.Channels(), assets.PanicOnMissing)
	assert.Equal(t, assets.ChannelUUID("a78930fe-6a40-4aa8-99c3-e61b02f45ca1"), call.Channel().UUID())
	assert.Equal(t, urns.URN("tel:+1234567890"), call.URN())
}
