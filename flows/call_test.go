package flows_test

import (
	"testing"

	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/static"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/engine"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCall(t *testing.T) {
	env := envs.NewBuilder().Build()

	source, err := static.NewSource([]byte(`{
		"channels": [
			{
				"uuid": "3a05eaf5-cb1b-4246-bef1-f277419c83a7",
				"name": "Nexmo",
				"address": "+16055742523",
				"schemes": [
					"tel"
				],
				"roles": [
					"call",
					"answer"
				]
			}
		]
	}`))
	require.NoError(t, err)

	sa, err := engine.NewSessionAssets(env, source, nil)
	require.NoError(t, err)

	vonage := sa.Channels().Get("3a05eaf5-cb1b-4246-bef1-f277419c83a7")

	call := flows.NewCall(
		"01978a2f-ad9a-7f2e-ad44-6e7547078cec",
		vonage,
		urns.URN("tel:+1234567890"),
	)

	// test marshaling our call
	ce := &flows.CallEnvelope{
		UUID:    "01978a2f-ad9a-7f2e-ad44-6e7547078cec",
		Channel: assets.NewChannelReference("3a05eaf5-cb1b-4246-bef1-f277419c83a7", "Nexmo"),
		URN:     urns.URN("tel:+1234567890"),
	}
	assert.Equal(t, ce, call.Marshal())

	// test unmarshaling
	call = ce.Unmarshal(sa, assets.PanicOnMissing)
	assert.Equal(t, assets.ChannelUUID("3a05eaf5-cb1b-4246-bef1-f277419c83a7"), call.Channel().UUID())
	assert.Equal(t, urns.URN("tel:+1234567890"), call.URN())
}
