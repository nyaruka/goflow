package types_test

import (
	"encoding/json"
	"testing"

	"github.com/nyaruka/goflow/assets/server/types"

	"github.com/stretchr/testify/assert"
)

var locationsJSON = `[
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
								"id": "575743222",
								"name": "Gisozi"
							},
							{
								"id": "457378732",
								"name": "Ndera"
							}
						]
					},
					{
						"name": "Nyarugenge",
						"children": []
					}
				]
			},
			{
				"name": "Eastern Province"
			}
		]
	}
]`

func TestReadLocations(t *testing.T) {
	locations, err := types.ReadLocations(json.RawMessage(locationsJSON))
	assert.NoError(t, err)
	assert.Equal(t, 1, len(locations))

	rwanda := locations[0]
	assert.Equal(t, "Rwanda", rwanda.Name())
	assert.Equal(t, []string{"Ruanda"}, rwanda.Aliases())
	assert.Equal(t, 2, len(rwanda.Children()))

	kigali := rwanda.Children()[0]
	assert.Equal(t, "Kigali City", kigali.Name())
	assert.Equal(t, []string{"Kigali", "Kigari"}, kigali.Aliases())
	assert.Equal(t, 2, len(kigali.Children()))

	gasabo := kigali.Children()[0]
	assert.Equal(t, "Gasabo", gasabo.Name())
	assert.Equal(t, 2, len(gasabo.Children()))

	ndera := gasabo.Children()[1]
	assert.Equal(t, "Ndera", ndera.Name())
	assert.Equal(t, 0, len(ndera.Children()))
}
