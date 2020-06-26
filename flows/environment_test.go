package flows_test

import (
	"testing"

	"github.com/nyaruka/goflow/assets/static"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/engine"

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

func TestEnvironment(t *testing.T) {
	env := envs.NewBuilder().WithDefaultCountry("RW").Build()
	source, err := static.NewSource([]byte(assetsJSON))
	require.NoError(t, err)

	sa, err := engine.NewSessionAssets(env, source, nil)
	require.NoError(t, err)

	fenv := flows.NewEnvironment(env, sa.Locations())
	assert.Equal(t, envs.Country("RW"), fenv.DefaultCountry())
	require.NotNil(t, fenv.LocationResolver())

	kigali := fenv.LocationResolver().LookupLocation("Rwanda > Kigali City")
	assert.Equal(t, "Kigali City", kigali.Name())

	matches := fenv.LocationResolver().FindLocationsFuzzy("gisozi town", flows.LocationLevelWard, nil)
	assert.Equal(t, 1, len(matches))
	assert.Equal(t, "Gisozi", matches[0].Name())
}
