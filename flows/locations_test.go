package flows_test

import (
	"encoding/json"
	"testing"

	"github.com/nyaruka/goflow/flows"

	"github.com/stretchr/testify/assert"
)

var locationHierarchyJSON = `
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
}`

func TestLocationHierarchy(t *testing.T) {
	hierarchy, err := flows.ReadLocationHierarchy(json.RawMessage(locationHierarchyJSON))
	assert.NoError(t, err)

	rwanda := hierarchy.Root()
	assert.Equal(t, flows.LocationLevel(0), rwanda.Level())
	assert.Equal(t, "Rwanda", rwanda.Name())
	assert.Equal(t, "Rwanda", rwanda.Path())
	assert.Equal(t, []string{"Ruanda"}, rwanda.Aliases())
	assert.Nil(t, rwanda.Parent())
	assert.Equal(t, 2, len(rwanda.Children()))

	kigali := rwanda.Children()[0]
	assert.Equal(t, flows.LocationLevel(1), kigali.Level())
	assert.Equal(t, "Kigali City", kigali.Name())
	assert.Equal(t, "Rwanda > Kigali City", kigali.Path())
	assert.Equal(t, []string{"Kigali", "Kigari"}, kigali.Aliases())
	assert.Equal(t, rwanda, kigali.Parent())
	assert.Equal(t, 2, len(kigali.Children()))

	gasabo := kigali.Children()[0]
	assert.Equal(t, flows.LocationLevel(2), gasabo.Level())
	assert.Equal(t, "Gasabo", gasabo.Name())
	assert.Equal(t, "Rwanda > Kigali City > Gasabo", gasabo.Path())
	assert.Equal(t, kigali, gasabo.Parent())
	assert.Equal(t, 2, len(gasabo.Children()))

	ndera := gasabo.Children()[1]
	assert.Equal(t, flows.LocationLevel(3), ndera.Level())
	assert.Equal(t, "Ndera", ndera.Name())
	assert.Equal(t, "Rwanda > Kigali City > Gasabo > Ndera", ndera.Path())
	assert.Equal(t, gasabo, ndera.Parent())
	assert.Equal(t, 0, len(ndera.Children()))

	assert.Equal(t, []*flows.Location{kigali}, hierarchy.FindByName("kigari", flows.LocationLevel(1), nil))
	assert.Equal(t, []*flows.Location{kigali}, hierarchy.FindByName("kigari", flows.LocationLevel(1), rwanda))
	assert.Equal(t, []*flows.Location{gasabo}, hierarchy.FindByName("GASABO", flows.LocationLevel(2), nil))
	assert.Equal(t, []*flows.Location{gasabo}, hierarchy.FindByName("GASABO", flows.LocationLevel(2), kigali))

	assert.Equal(t, []*flows.Location{}, hierarchy.FindByName("boston", flows.LocationLevel(1), nil))    // no such name
	assert.Equal(t, []*flows.Location{}, hierarchy.FindByName("kigari", flows.LocationLevel(8), nil))    // no such level
	assert.Equal(t, []*flows.Location{}, hierarchy.FindByName("kigari", flows.LocationLevel(2), nil))    // wrong level
	assert.Equal(t, []*flows.Location{}, hierarchy.FindByName("kigari", flows.LocationLevel(2), gasabo)) // wrong parent
}
