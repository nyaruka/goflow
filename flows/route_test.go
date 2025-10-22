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

func TestRoute(t *testing.T) {
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

	sa, err := engine.NewSessionAssets(env, source, nil)
	require.NoError(t, err)

	channel := sa.Channels().Get("57f1078f-88aa-46f4-a59a-948a5739c03d")

	// check that unmarshaling a route properly extracts its channel affinity
	envelope := &flows.RouteEnvelope{URN: "tel:+250781234567", Channel: assets.NewChannelReference("57f1078f-88aa-46f4-a59a-948a5739c03d", "Android")}

	route1 := envelope.Unmarshal(sa, assets.PanicOnMissing)
	assert.NoError(t, err)
	assert.Equal(t, urns.URN("tel:+250781234567"), route1.URN())
	assert.Equal(t, channel, route1.Channel())

	// check equality
	route2 := flows.NewRoute("tel:+250781234567", channel)
	route3 := flows.NewRoute("tel:+250781234567", nil)
	assert.True(t, route1.Equal(route2))
	assert.False(t, route1.Equal(route3))

	// becomes just a URN string in expressions
	assert.Equal(t, types.NewXText("tel:+250781234567"), route1.ToXValue(env))

	// check when URNs have to be redacted
	env = envs.NewBuilder().WithRedactionPolicy(envs.RedactionPolicyURNs).Build()
	assert.Equal(t, types.NewXText("tel:********"), route1.ToXValue(env))

	// we can clear the channel affinity
	route1.SetChannel(nil)
	assert.Equal(t, urns.URN("tel:+250781234567"), route1.URN())
	assert.Nil(t, route1.Channel())

	// and change it
	route1.SetChannel(channel)
	assert.Equal(t, urns.URN("tel:+250781234567"), route1.URN())
	assert.Equal(t, channel, route1.Channel())
}

func TestRouteList(t *testing.T) {
	r1 := flows.NewRoute("tel:+250781234567", nil)
	r2 := flows.NewRoute("twitter:134252511151#billy_bob", nil)
	r3 := flows.NewRoute("tel:+250781111222", nil)
	routes := flows.RouteList{r1, r2, r3}

	env := envs.NewBuilder().Build()

	// check equality
	assert.True(t, routes.Equal(flows.RouteList{r1, r2, r3}))
	assert.False(t, routes.Equal(flows.RouteList{r3, r2, r1}))
	assert.False(t, routes.Equal(flows.RouteList{r1, r2}))

	// check use in expressions
	test.AssertXEqual(t, types.NewXArray(
		types.NewXText("tel:+250781234567"),
		types.NewXText("twitter:134252511151#billy_bob"),
		types.NewXText("tel:+250781111222"),
	), routes.ToXValue(env))
}
