package utils_test

import (
	"encoding/json"
	"testing"

	"github.com/nyaruka/goflow/utils"

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
	hierarchy, err := utils.ReadLocationHierarchy(json.RawMessage(locationHierarchyJSON))
	assert.NoError(t, err)

	rwanda := hierarchy.Root()
	assert.Equal(t, utils.LocationLevel(0), rwanda.Level())
	assert.Equal(t, "Rwanda", rwanda.Name())
	assert.Equal(t, "Rwanda", rwanda.Path())
	assert.Equal(t, []string{"Ruanda"}, rwanda.Aliases())
	assert.Nil(t, rwanda.Parent())
	assert.Equal(t, 2, len(rwanda.Children()))

	kigali := rwanda.Children()[0]
	assert.Equal(t, utils.LocationLevel(1), kigali.Level())
	assert.Equal(t, "Kigali City", kigali.Name())
	assert.Equal(t, "Rwanda > Kigali City", kigali.Path())
	assert.Equal(t, []string{"Kigali", "Kigari"}, kigali.Aliases())
	assert.Equal(t, rwanda, kigali.Parent())
	assert.Equal(t, 2, len(kigali.Children()))

	gasabo := kigali.Children()[0]
	assert.Equal(t, utils.LocationLevel(2), gasabo.Level())
	assert.Equal(t, "Gasabo", gasabo.Name())
	assert.Equal(t, "Rwanda > Kigali City > Gasabo", gasabo.Path())
	assert.Equal(t, kigali, gasabo.Parent())
	assert.Equal(t, 2, len(gasabo.Children()))

	ndera := gasabo.Children()[1]
	assert.Equal(t, utils.LocationLevel(3), ndera.Level())
	assert.Equal(t, "Ndera", ndera.Name())
	assert.Equal(t, "Rwanda > Kigali City > Gasabo > Ndera", ndera.Path())
	assert.Equal(t, gasabo, ndera.Parent())
	assert.Equal(t, 0, len(ndera.Children()))

	assert.Equal(t, []*utils.Location{rwanda}, hierarchy.FindByName("RWaNdA", utils.LocationLevel(0), nil))
	assert.Equal(t, []*utils.Location{kigali}, hierarchy.FindByName("kigari", utils.LocationLevel(1), nil))
	assert.Equal(t, []*utils.Location{kigali}, hierarchy.FindByName("rwanda > kigali city", utils.LocationLevel(1), nil))
	assert.Equal(t, []*utils.Location{kigali}, hierarchy.FindByName("kigari", utils.LocationLevel(1), rwanda))
	assert.Equal(t, []*utils.Location{gasabo}, hierarchy.FindByName("GASABO", utils.LocationLevel(2), nil))
	assert.Equal(t, []*utils.Location{gasabo}, hierarchy.FindByName("GASABO", utils.LocationLevel(2), kigali))
	assert.Equal(t, []*utils.Location{ndera}, hierarchy.FindByName("RWANDA > kigali city > gasabo > ndera", utils.LocationLevel(3), nil))

	assert.Equal(t, []*utils.Location{}, hierarchy.FindByName("boston", utils.LocationLevel(1), nil))    // no such name
	assert.Equal(t, []*utils.Location{}, hierarchy.FindByName("kigari", utils.LocationLevel(8), nil))    // no such level
	assert.Equal(t, []*utils.Location{}, hierarchy.FindByName("kigari", utils.LocationLevel(2), nil))    // wrong level
	assert.Equal(t, []*utils.Location{}, hierarchy.FindByName("kigari", utils.LocationLevel(2), gasabo)) // wrong parent

	assert.Equal(t, rwanda, hierarchy.FindByPath("RWANDA"))
	assert.Equal(t, kigali, hierarchy.FindByPath("RWANDA > KIGALI CITY"))
	assert.Equal(t, gasabo, hierarchy.FindByPath("rwanda > kigali city > gasabo"))
	assert.Equal(t, ndera, hierarchy.FindByPath("rwanda > kigali city > gasabo > ndera"))
}
