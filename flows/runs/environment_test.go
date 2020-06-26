package runs_test

import (
	"testing"
	"time"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/static"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/engine"
	"github.com/nyaruka/goflow/flows/triggers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var assetsJSON = `{
    "flows": [
        {
            "uuid": "76f0a02f-3b75-4b86-9064-e9195e1b3a02",
            "name": "Test",
            "spec_version": "13.1.0",
            "language": "eng",
            "type": "messaging",
            "nodes": []
        }
	],
	"channels": [
    	{
			"uuid": "57f1078f-88aa-46f4-a59a-948a5739c03d",
			"name": "Android Channel",
			"address": "+17036975131",
			"schemes": ["tel"],
			"roles": ["send", "receive"],
			"country": "US"
    	}
	  ],
	  "locations": [
        {
            "name": "Rwanda",
            "aliases": ["Ruanda"],		
            "children": [
                {
                    "name": "Kigali City",
                    "aliases": ["Kigali", "Kigari"],
                    "children": [
                        {
                            "name": "Gasabo",
                            "children": [
                                {
                                    "name": "Gisozi"
                                },
                                {
                                    "name": "Ndera"
                                }
                            ]
                        },
                        {
                            "name": "Nyarugenge",
                            "children": []
                        }
                    ]
                }
            ]
        }
    ]
}`

const contactJSON = `{
	"uuid": "ba96bf7f-bc2a-4873-a7c7-254d1927c4e3",
	"id": 1234567,
	"name": "Ben Haggerty",
	"created_on": "2018-01-01T12:00:00.000000000-00:00",
	"fields": {},
	"timezone": "America/Guayaquil",
	"urns": [
		"tel:+12065551212"
	]
}`

func TestRunEnvironment(t *testing.T) {
	tzRW, _ := time.LoadLocation("Africa/Kigali")
	tzEC, _ := time.LoadLocation("America/Guayaquil")

	env := envs.NewBuilder().
		WithDefaultCountry("RW").
		WithTimezone(tzRW).
		Build()
	source, err := static.NewSource([]byte(assetsJSON))
	require.NoError(t, err)

	sa, err := engine.NewSessionAssets(env, source, nil)
	require.NoError(t, err)

	contact, err := flows.ReadContact(sa, []byte(contactJSON), assets.IgnoreMissing)
	require.NoError(t, err)

	trigger := triggers.NewManual(env, assets.NewFlowReference("76f0a02f-3b75-4b86-9064-e9195e1b3a02", "Test"), contact, false, nil)
	eng := engine.NewBuilder().Build()

	session, _, err := eng.NewSession(sa, trigger)
	require.NoError(t, err)

	// environment on the session has the values we started with
	assert.Equal(t, envs.Country("RW"), session.Environment().DefaultCountry())
	assert.Equal(t, tzRW, session.Environment().Timezone())

	// environment on the run has values from the contact
	runEnv := session.Runs()[0].Environment()
	assert.Equal(t, envs.Country("US"), runEnv.DefaultCountry())
	assert.Equal(t, tzEC, runEnv.Timezone())
	assert.NotNil(t, runEnv.LocationResolver())
}
