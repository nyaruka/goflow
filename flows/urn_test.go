package flows_test

import (
	"testing"

	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/static"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/engine"
	"github.com/nyaruka/goflow/test"
	"github.com/nyaruka/goflow/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestURNValidation(t *testing.T) {
	type testStruct struct {
		ValidURN      string `json:"valid_urn" validate:"urn"`
		InvalidURN    string `json:"invalid_urn" validate:"urn"`
		ValidScheme   string `json:"valid_scheme" validate:"urnscheme"`
		InvalidScheme string `json:"invalid_scheme" validate:"urnscheme"`
	}

	obj := testStruct{
		ValidURN:      "tel:+123456789",
		InvalidURN:    "xyz",
		ValidScheme:   "viber",
		InvalidScheme: "$%@^#^^!!!",
	}
	err := utils.Validate(obj)
	assert.EqualError(t, err, "field 'invalid_urn' is not a valid URN, field 'invalid_scheme' is not a valid URN scheme")
}

func TestURN(t *testing.T) {
	env := envs.NewBuilder().Build()

	source, err := static.NewSource([]byte(`{
        "channels": [
            {
                "uuid": "57f1078f-88aa-46f4-a59a-948a5739c03d",
                "name": "Android Channel",
				"address": "+17036975131",
				"schemes": [
					"tel"
				],
				"roles": [
					"send",
					"receive"
				],
				"country": "US"
            }
        ]
    }`))
	require.NoError(t, err)

	sessionAssets, err := engine.NewSessionAssets(env, source, nil)
	require.NoError(t, err)

	channels := sessionAssets.Channels()
	channel := channels.Get("57f1078f-88aa-46f4-a59a-948a5739c03d")

	// check that parsing a URN properly extracts its channel affinity
	urn, err := flows.ParseURN(channels, "tel:+250781234567?channel=57f1078f-88aa-46f4-a59a-948a5739c03d", assets.PanicOnMissing)
	assert.NoError(t, err)
	assert.Equal(t, "tel", urn.Scheme)
	assert.Equal(t, "+250781234567", urn.Path)
	assert.Equal(t, "", urn.Display)
	assert.Equal(t, channel, urn.Channel)
	assert.Equal(t, urns.URN("tel:+250781234567?channel=57f1078f-88aa-46f4-a59a-948a5739c03d"), urn.Encode())
	assert.Equal(t, urns.URN("tel:+250781234567"), urn.Identity())

	// check equality
	urn2, _ := flows.ParseURN(channels, "tel:+250781234567?channel=57f1078f-88aa-46f4-a59a-948a5739c03d", assets.PanicOnMissing)
	urn3, _ := flows.ParseURN(channels, "tel:+250781234567", assets.PanicOnMissing)
	assert.True(t, urn.Equal(urn2))
	assert.False(t, urn.Equal(urn3))

	// check using URN in expressions
	assert.Equal(t, types.NewXText("tel:+250781234567"), urn.ToXValue(env))

	// check when URNs have to be redacted
	env = envs.NewBuilder().WithRedactionPolicy(envs.RedactionPolicyURNs).Build()
	assert.Equal(t, types.NewXText("tel:********"), urn.ToXValue(env))

	// we can clear the channel affinity
	urn.Channel = nil
	assert.Equal(t, urns.URN("tel:+250781234567"), urn.Encode())
	assert.Nil(t, urn.Channel)

	// and change it
	urn.Channel = channel
	assert.Equal(t, urns.URN("tel:+250781234567?channel=57f1078f-88aa-46f4-a59a-948a5739c03d"), urn.Encode())
	assert.Equal(t, channel, urn.Channel)
}

func TestURNList(t *testing.T) {
	urn1 := flows.NewURN("tel", "+250781234567", "", nil)
	urn2 := flows.NewURN("twitter", "134252511151", "billy_bob", nil)
	urn3 := flows.NewURN("tel", "+250781111222", "", nil)
	urnList := flows.URNList{urn1, urn2, urn3}

	env := envs.NewBuilder().Build()

	// check equality
	assert.True(t, urnList.Equal(flows.URNList{urn1, urn2, urn3}))
	assert.False(t, urnList.Equal(flows.URNList{urn3, urn2, urn1}))
	assert.False(t, urnList.Equal(flows.URNList{urn1, urn2}))

	// check use in expressions
	test.AssertXEqual(t, types.NewXArray(
		types.NewXText("tel:+250781234567"),
		types.NewXText("twitter:134252511151#billy_bob"),
		types.NewXText("tel:+250781111222"),
	), urnList.ToXValue(env))
}
