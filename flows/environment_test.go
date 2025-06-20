package flows_test

import (
	"context"
	"testing"
	"time"

	"github.com/nyaruka/gocommon/i18n"
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

func TestLocationResolver(t *testing.T) {
	env := envs.NewBuilder().WithDefaultCountry("RW").Build()
	source, err := static.NewSource([]byte(assetsJSON))
	require.NoError(t, err)

	sa, err := engine.NewSessionAssets(env, source, nil)
	require.NoError(t, err)

	contact := flows.NewEmptyContact(sa, "", i18n.NilLanguage, nil)

	trigger := triggers.NewBuilder(env, assets.NewFlowReference("76f0a02f-3b75-4b86-9064-e9195e1b3a02", "Test"), contact).Manual().Build()
	eng := engine.NewBuilder().Build()

	session, _, err := eng.NewSession(context.Background(), sa, trigger, nil)
	require.NoError(t, err)

	resolver := session.Assets().Locations()

	kigali := resolver.LookupLocation("Rwanda > Kigali City")
	assert.Equal(t, "Kigali City", kigali.Name())

	matches := resolver.FindLocationsFuzzy(env, "gisozi town", flows.LocationLevelWard, nil)
	assert.Equal(t, 1, len(matches))
	assert.Equal(t, "Gisozi", matches[0].Name())
}

const contactJSON = `{
	"uuid": "ba96bf7f-bc2a-4873-a7c7-254d1927c4e3",
	"id": 1234567,
	"name": "Ben Haggerty",
	"created_on": "2018-01-01T12:00:00.000000000-00:00",
	"fields": {},
	"language": "fra",
	"timezone": "America/Guayaquil",
	"urns": [
		"tel:+12065551212"
	]
}`

func TestSessionEnvironment(t *testing.T) {
	tzRW, _ := time.LoadLocation("Africa/Kigali")
	tzEC, _ := time.LoadLocation("America/Guayaquil")
	tzUK, _ := time.LoadLocation("Europe/London")

	env := envs.NewBuilder().
		WithAllowedLanguages("eng", "fra", "kin").
		WithDefaultCountry("RW").
		WithTimezone(tzRW).
		Build()
	source, err := static.NewSource([]byte(assetsJSON))
	require.NoError(t, err)

	sa, err := engine.NewSessionAssets(env, source, nil)
	require.NoError(t, err)

	contact, err := flows.ReadContact(sa, []byte(contactJSON), assets.IgnoreMissing)
	require.NoError(t, err)

	trigger := triggers.NewBuilder(env, assets.NewFlowReference("76f0a02f-3b75-4b86-9064-e9195e1b3a02", "Test"), contact).Manual().Build()
	eng := engine.NewBuilder().Build()

	session, _, err := eng.NewSession(context.Background(), sa, trigger, nil)
	require.NoError(t, err)

	// main environment on the session has the values we started with
	serializedEnv := session.Environment()
	assert.Equal(t, i18n.Language("eng"), serializedEnv.DefaultLanguage())
	assert.Equal(t, []i18n.Language{"eng", "fra", "kin"}, serializedEnv.AllowedLanguages())
	assert.Equal(t, i18n.Country("RW"), serializedEnv.DefaultCountry())
	assert.Equal(t, i18n.Locale("eng-RW"), serializedEnv.DefaultLocale())
	assert.Equal(t, tzRW, serializedEnv.Timezone())

	// merged environment on the session has values from the contact
	mergedEnv := session.MergedEnvironment()
	assert.Equal(t, i18n.Language("fra"), mergedEnv.DefaultLanguage())
	assert.Equal(t, []i18n.Language{"eng", "fra", "kin"}, mergedEnv.AllowedLanguages())
	assert.Equal(t, i18n.Country("US"), mergedEnv.DefaultCountry())
	assert.Equal(t, i18n.Locale("fra-US"), mergedEnv.DefaultLocale())
	assert.Equal(t, tzEC, mergedEnv.Timezone())
	assert.NotNil(t, mergedEnv.LocationResolver())

	// can make changes to contact
	session.Contact().SetLanguage(i18n.Language("kin"))
	session.Contact().SetTimezone(tzUK)

	// and environment reflects those changes
	assert.Equal(t, i18n.Language("kin"), mergedEnv.DefaultLanguage())
	assert.Equal(t, tzUK, mergedEnv.Timezone())

	// if contact language is not an allowed language it won't be used
	session.Contact().SetLanguage(i18n.Language("spa"))
	assert.Equal(t, i18n.Language("eng"), mergedEnv.DefaultLanguage())
	assert.Equal(t, i18n.Locale("eng-US"), mergedEnv.DefaultLocale())
}
