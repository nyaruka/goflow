package flows_test

import (
	"testing"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/static"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/engine"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDependencies(t *testing.T) {
	assert.Equal(t, &flows.Dependencies{}, flows.NewDependencies([]assets.Reference{}))

	deps := flows.NewDependencies([]assets.Reference{
		assets.NewChannelReference("8286545d-d1a1-4eff-a3ad-a11ddf4bb20a", "Android"),
		flows.NewContactReference("0b099519-0889-4c74-b744-9122272f346a", "Bob"),
		assets.NewFieldReference("gender", "Gender"),
		assets.NewFlowReference("4f932672-7995-47f0-96e6-faf5abd2d81d", "Registration"),
		assets.NewGroupReference("46057a92-6580-4e93-af36-2bb9c9d61e51", "Testers"),
		assets.NewGroupReference("377c3101-a7fc-47b1-9136-980348e362c0", "Customers"),
		assets.NewLabelReference("31c06b7c-010d-4f91-9590-d3fbdc2fb7ac", "Spam"),
		assets.NewTemplateReference("ff958d30-f50e-48ab-a524-37ed1e9620d9", "Welcome"),
	})

	assert.Equal(t, &flows.Dependencies{
		Channels: []*assets.ChannelReference{
			assets.NewChannelReference("8286545d-d1a1-4eff-a3ad-a11ddf4bb20a", "Android"),
		},
		Contacts: []*flows.ContactReference{
			flows.NewContactReference("0b099519-0889-4c74-b744-9122272f346a", "Bob"),
		},
		Fields: []*assets.FieldReference{
			assets.NewFieldReference("gender", "Gender"),
		},
		Flows: []*assets.FlowReference{
			assets.NewFlowReference("4f932672-7995-47f0-96e6-faf5abd2d81d", "Registration"),
		},
		Groups: []*assets.GroupReference{
			assets.NewGroupReference("46057a92-6580-4e93-af36-2bb9c9d61e51", "Testers"),
			assets.NewGroupReference("377c3101-a7fc-47b1-9136-980348e362c0", "Customers"),
		},
		Labels: []*assets.LabelReference{
			assets.NewLabelReference("31c06b7c-010d-4f91-9590-d3fbdc2fb7ac", "Spam"),
		},
		Templates: []*assets.TemplateReference{
			assets.NewTemplateReference("ff958d30-f50e-48ab-a524-37ed1e9620d9", "Welcome"),
		},
	}, deps)

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

	sa, err := engine.NewSessionAssets(source)
	require.NoError(t, err)

	missing := make([]assets.Reference, 0)
	deps.Check(sa, func(ref assets.Reference, err error) {
		missing = append(missing, ref)
	})

	// check the contact reference is not included, and the group which does exist in the assets
	assert.Equal(t, []assets.Reference{
		assets.NewChannelReference("8286545d-d1a1-4eff-a3ad-a11ddf4bb20a", "Android"),
		assets.NewFieldReference("gender", "Gender"),
		assets.NewFlowReference("4f932672-7995-47f0-96e6-faf5abd2d81d", "Registration"),
		assets.NewGroupReference("46057a92-6580-4e93-af36-2bb9c9d61e51", "Testers"),
		assets.NewLabelReference("31c06b7c-010d-4f91-9590-d3fbdc2fb7ac", "Spam"),
		assets.NewTemplateReference("ff958d30-f50e-48ab-a524-37ed1e9620d9", "Welcome"),
	}, missing)
}

func TestResultInfos(t *testing.T) {
	assert.Equal(t, []*flows.ResultInfo{}, flows.MergeResultInfos(nil))

	assert.Equal(t, []*flows.ResultInfo{
		{Key: "response_1", Name: "Response 1", Categories: []string{"Red", "Green", "Blue"}},
		{Key: "favorite_beer", Name: "Favorite Beer", Categories: []string{}},
	}, flows.MergeResultInfos([]*flows.ResultInfo{
		flows.NewResultInfo("Response 1", []string{"Red", "Green"}),
		flows.NewResultInfo("Response-1", nil),
		flows.NewResultInfo("Response-1", []string{"Green", "Blue"}),
		flows.NewResultInfo("Favorite Beer", []string{}),
	}))

	assert.Equal(t, `key=response_1|name=Response 1|categories=Red,Green`, flows.NewResultInfo("Response 1", []string{"Red", "Green"}).String())
}
