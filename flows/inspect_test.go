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
	assert.Equal(t, &flows.Dependencies{}, flows.NewDependencies([]assets.Reference{}))

	deps := flows.NewDependencies([]assets.Reference{
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
	})

	depsJSON, _ := json.Marshal(deps)

	test.AssertEqualJSON(t, []byte(`{
		"channels": [
			{"uuid": "8286545d-d1a1-4eff-a3ad-a11ddf4bb20a", "name": "Android"}
		],
		"classifiers": [
			{"uuid": "2138cddc-118a-49ae-b290-98e03ad0573b", "name": "Booking"}
		],
		"contacts": [
			{"uuid": "0b099519-0889-4c74-b744-9122272f346a", "name": "Bob"}
		],
		"fields": [
			{"key": "gender", "name": "Gender"}
		],
		"flows": [
			{"uuid": "4f932672-7995-47f0-96e6-faf5abd2d81d", "name": "Registration"}
		],
		"globals": [
			{"key": "org_name", "name": "Org Name"}
		],
		"groups": [
			{"uuid": "46057a92-6580-4e93-af36-2bb9c9d61e51", "name": "Testers"},
			{"uuid": "377c3101-a7fc-47b1-9136-980348e362c0", "name": "Customers"}
		],
		"labels": [
			{"uuid": "31c06b7c-010d-4f91-9590-d3fbdc2fb7ac", "name": "Spam"}
		],
		"templates": [
			{"uuid": "ff958d30-f50e-48ab-a524-37ed1e9620d9", "name": "Welcome"}
		]
	}`), depsJSON, "deps JSON mismatch")

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

	missing := make([]assets.Reference, 0)
	deps.Check(sa, func(ref assets.Reference, err error) {
		missing = append(missing, ref)
	})

	// check the contact reference is not included, and the group which does exist in the assets
	assert.Equal(t, []assets.Reference{
		assets.NewChannelReference("8286545d-d1a1-4eff-a3ad-a11ddf4bb20a", "Android"),
		assets.NewClassifierReference("2138cddc-118a-49ae-b290-98e03ad0573b", "Booking"),
		assets.NewFieldReference("gender", "Gender"),
		assets.NewFlowReference("4f932672-7995-47f0-96e6-faf5abd2d81d", "Registration"),
		assets.NewGlobalReference("org_name", "Org Name"),
		assets.NewGroupReference("46057a92-6580-4e93-af36-2bb9c9d61e51", "Testers"),
		assets.NewLabelReference("31c06b7c-010d-4f91-9590-d3fbdc2fb7ac", "Spam"),
		assets.NewTemplateReference("ff958d30-f50e-48ab-a524-37ed1e9620d9", "Welcome"),
	}, missing)

	// panic if we get a dependency type we don't recognize
	assert.Panics(t, func() {
		flows.NewDependencies([]assets.Reference{&unknownAssetType{}})
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

func TestFlowInfo(t *testing.T) {
	info := &flows.FlowInfo{
		Dependencies: flows.NewDependencies([]assets.Reference{
			assets.NewGroupReference("46057a92-6580-4e93-af36-2bb9c9d61e51", "Testers"),
			assets.NewLabelReference("31c06b7c-010d-4f91-9590-d3fbdc2fb7ac", "Spam"),
		}),
		Results: []*flows.ResultInfo{
			{Key: "response_1", Name: "Response 1", Categories: []string{"Red", "Green", "Blue"}, NodeUUIDs: []flows.NodeUUID{"edcbe7a9-3b1b-4f49-891e-9519f0309e8b"}},
			{Key: "favorite_beer", Name: "Favorite Beer", Categories: []string{}, NodeUUIDs: []flows.NodeUUID{"0a6f263b-6258-4007-954a-23c20bcd333e"}},
		},
		WaitingExits: []flows.ExitUUID{
			"9d098aea-ccc4-4723-8222-9971b64223e4",
			"8c50f16e-35d0-4e08-a725-33ca1c03ef62",
		},
		ParentRefs: []string{"state", "response_2"},
	}

	// test marshaling
	marshaled, err := json.Marshal(info)
	require.NoError(t, err)

	test.AssertEqualJSON(t, []byte(`{
		"dependencies": {
			"groups": [
				{
					"name": "Testers",
					"uuid": "46057a92-6580-4e93-af36-2bb9c9d61e51"
				}
			],
			"labels": [
				{
					"name": "Spam",
					"uuid": "31c06b7c-010d-4f91-9590-d3fbdc2fb7ac"
				}
			]
		},
		"results": [
			{
				"key": "response_1",
				"name": "Response 1",
				"categories": [
					"Red",
					"Green",
					"Blue"
				],
				"node_uuids": [
					"edcbe7a9-3b1b-4f49-891e-9519f0309e8b"
				]
			},
			{
				"key": "favorite_beer",
				"name": "Favorite Beer",
				"categories": [],
				"node_uuids": [
					"0a6f263b-6258-4007-954a-23c20bcd333e"
				]
			}
		],
		"waiting_exits": [
			"9d098aea-ccc4-4723-8222-9971b64223e4",
			"8c50f16e-35d0-4e08-a725-33ca1c03ef62"
		],
		"parent_refs": [
			"state",
			"response_2"
		]
	}`), marshaled, "marshal mismatch")
}

func TestInspectedReferences(t *testing.T) {
	r := &flows.InspectedChannelReference{
		ChannelReference:   assets.NewChannelReference("9d098aea-ccc4-4723-8222-9971b64223e4", "Android"),
		InspectedReference: flows.InspectedReference{Missing: true},
	}

	j, _ := json.Marshal(r)
	assert.Equal(t, `{"uuid":"9d098aea-ccc4-4723-8222-9971b64223e4","name":"Android","missing":true}`, string(j))
}
