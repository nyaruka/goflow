package flows_test

import (
	"encoding/json"
	"testing"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/static"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/definition"
	"github.com/nyaruka/goflow/flows/engine"
	"github.com/nyaruka/goflow/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type unknownAssetType struct{}

func (t *unknownAssetType) String() string   { return "x" }
func (t *unknownAssetType) Type() string     { return "unknown" }
func (t *unknownAssetType) Identity() string { return "unknown[]" }
func (t *unknownAssetType) Variable() bool   { return false }

func TestDependencies(t *testing.T) {
	refs := []assets.Reference{
		assets.NewChannelReference("8286545d-d1a1-4eff-a3ad-a11ddf4bb20a", "Android"),
		assets.NewClassifierReference("2138cddc-118a-49ae-b290-98e03ad0573b", "Booking"),
		flows.NewContactReference("0b099519-0889-4c74-b744-9122272f346a", "Bob"),
		assets.NewFieldReference("gender", "Gender"),
		assets.NewFlowReference("4f932672-7995-47f0-96e6-faf5abd2d81d", "Registration"),
		assets.NewGlobalReference("org_name", "Org Name"),
		assets.NewGroupReference("46057a92-6580-4e93-af36-2bb9c9d61e51", "Testers"),
		assets.NewGroupReference("377c3101-a7fc-47b1-9136-980348e362c0", "Customers"),
		assets.NewLabelReference("31c06b7c-010d-4f91-9590-d3fbdc2fb7ac", "Spam"),
		assets.NewTemplateReference("ff958d30-f50e-48ab-a524-37ed1e9620d9", "Welcome"),
	}

	// can inspect without assets
	deps := flows.InspectReferences(refs, nil)
	depsJSON, _ := json.Marshal(deps)
	test.AssertEqualJSON(t, []byte(`[
		{
			"name": "Android",
			"type": "channel",
			"uuid": "8286545d-d1a1-4eff-a3ad-a11ddf4bb20a"
		},
		{
			"name": "Booking",
			"type": "classifier",
			"uuid": "2138cddc-118a-49ae-b290-98e03ad0573b"
		},
		{
			"name": "Bob",
			"type": "contact",
			"uuid": "0b099519-0889-4c74-b744-9122272f346a"
		},
		{
			"key": "gender",
			"name": "Gender",
			"type": "field"
		},
		{
			"name": "Registration",
			"type": "flow",
			"uuid": "4f932672-7995-47f0-96e6-faf5abd2d81d"
		},
		{
			"key": "org_name",
			"name": "Org Name",
			"type": "global"
		},
		{
			"name": "Testers",
			"type": "group",
			"uuid": "46057a92-6580-4e93-af36-2bb9c9d61e51"
		},
		{
			"name": "Customers",
			"type": "group",
			"uuid": "377c3101-a7fc-47b1-9136-980348e362c0"
		},
		{
			"name": "Spam",
			"type": "label",
			"uuid": "31c06b7c-010d-4f91-9590-d3fbdc2fb7ac"
		},
		{
			"name": "Welcome",
			"type": "template",
			"uuid": "ff958d30-f50e-48ab-a524-37ed1e9620d9"
		}
	]`), depsJSON, "deps JSON mismatch")

	// if our assets only includes a single group, the other assets should be reported as missing
	source, err := static.NewSource([]byte(`{
		"groups": [
			{
				"uuid": "377c3101-a7fc-47b1-9136-980348e362c0",
				"name": "Customers"
			}
		]
	}`))
	require.NoError(t, err)

	sa, err := engine.NewSessionAssets(source, nil)
	require.NoError(t, err)

	deps = flows.InspectReferences(refs, sa)
	depsJSON, _ = json.Marshal(deps)
	test.AssertEqualJSON(t, []byte(`[
		{
			"missing": true,
			"name": "Android",
			"type": "channel",
			"uuid": "8286545d-d1a1-4eff-a3ad-a11ddf4bb20a"
		},
		{
			"missing": true,
			"name": "Booking",
			"type": "classifier",
			"uuid": "2138cddc-118a-49ae-b290-98e03ad0573b"
		},
		{
			"name": "Bob",
			"type": "contact",
			"uuid": "0b099519-0889-4c74-b744-9122272f346a"
		},
		{
			"key": "gender",
			"missing": true,
			"name": "Gender",
			"type": "field"
		},
		{
			"missing": true,
			"name": "Registration",
			"type": "flow",
			"uuid": "4f932672-7995-47f0-96e6-faf5abd2d81d"
		},
		{
			"key": "org_name",
			"missing": true,
			"name": "Org Name",
			"type": "global"
		},
		{
			"missing": true,
			"name": "Testers",
			"type": "group",
			"uuid": "46057a92-6580-4e93-af36-2bb9c9d61e51"
		},
		{
			"name": "Customers",
			"type": "group",
			"uuid": "377c3101-a7fc-47b1-9136-980348e362c0"
		},
		{
			"missing": true,
			"name": "Spam",
			"type": "label",
			"uuid": "31c06b7c-010d-4f91-9590-d3fbdc2fb7ac"
		},
		{
			"missing": true,
			"name": "Welcome",
			"type": "template",
			"uuid": "ff958d30-f50e-48ab-a524-37ed1e9620d9"
		}
	]`), depsJSON, "deps JSON mismatch")

	// panic if we get a dependency type we don't recognize
	assert.Panics(t, func() {
		flows.InspectReferences([]assets.Reference{&unknownAssetType{}}, nil)
	})
}

func TestResultInfos(t *testing.T) {
	assert.Equal(t, []*flows.ResultInfo{}, flows.MergeResultInfos(nil))

	node1 := definition.NewNode(
		flows.NodeUUID("1fb823c3-599a-41e9-b59b-658266af3466"),
		nil,
		nil,
		[]flows.Exit{definition.NewExit(flows.ExitUUID("3c158842-24f3-4a40-bea4-7522952c0131"), "")},
	)
	node2 := definition.NewNode(
		flows.NodeUUID("0ba673a3-63b3-46f9-9246-9c727cf2917f"),
		nil,
		nil,
		[]flows.Exit{definition.NewExit(flows.ExitUUID("434ac29c-abe6-4bd7-b29b-740d517b6bb5"), "")},
	)

	infos := []*flows.ResultInfo{
		flows.NewResultInfo("Response 1", []string{"Red", "Green"}, node1),
		flows.NewResultInfo("Response-1", nil, node1),
		flows.NewResultInfo("Response-1", []string{"Green", "Blue"}, node2),
		flows.NewResultInfo("Favorite Beer", []string{}, node2),
	}

	assert.Equal(t, []*flows.ResultInfo{
		{
			Key:        "response_1",
			Name:       "Response 1",
			Categories: []string{"Red", "Green", "Blue"},
			NodeUUIDs: []flows.NodeUUID{
				flows.NodeUUID("1fb823c3-599a-41e9-b59b-658266af3466"),
				flows.NodeUUID("0ba673a3-63b3-46f9-9246-9c727cf2917f"),
			},
		},
		{
			Key:        "favorite_beer",
			Name:       "Favorite Beer",
			Categories: []string{},
			NodeUUIDs: []flows.NodeUUID{
				flows.NodeUUID("0ba673a3-63b3-46f9-9246-9c727cf2917f"),
			},
		},
	}, flows.MergeResultInfos(infos))

	assert.Equal(t, `key=response_1|name=Response 1|categories=Red,Green`, flows.NewResultInfo("Response 1", []string{"Red", "Green"}, node1).String())
}
