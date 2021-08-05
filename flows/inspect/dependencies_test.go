package inspect_test

import (
	"testing"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/static"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/actions"
	"github.com/nyaruka/goflow/flows/definition"
	"github.com/nyaruka/goflow/flows/engine"
	"github.com/nyaruka/goflow/flows/inspect"
	"github.com/nyaruka/goflow/flows/routers"
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
	env := envs.NewBuilder().Build()

	action1 := actions.NewSendMsg("ed08e6b9-ed22-4294-9871-c7ac7d82cbd5", "Hi there", nil, nil, false)
	node1 := definition.NewNode("91b20e13-d6e2-42a9-b74f-bce85c9da8c8", []flows.Action{action1}, nil, nil)
	router2 := routers.NewRandom(nil, "", nil)
	node2 := definition.NewNode("7c959933-4c30-4277-9810-adc95a459bd0", nil, router2, nil)

	refs := []flows.ExtractedReference{
		flows.NewExtractedReference(node1, action1, nil, envs.NilLanguage, assets.NewChannelReference("8286545d-d1a1-4eff-a3ad-a11ddf4bb20a", "Android")),
		flows.NewExtractedReference(node1, action1, nil, envs.NilLanguage, assets.NewClassifierReference("2138cddc-118a-49ae-b290-98e03ad0573b", "Booking")),
		flows.NewExtractedReference(node1, action1, nil, envs.NilLanguage, flows.NewContactReference("0b099519-0889-4c74-b744-9122272f346a", "Bob")),
		flows.NewExtractedReference(node1, action1, nil, envs.NilLanguage, assets.NewFieldReference("gender", "Gender")),
		flows.NewExtractedReference(node1, action1, nil, envs.NilLanguage, assets.NewFlowReference("4f932672-7995-47f0-96e6-faf5abd2d81d", "Registration")),
		flows.NewExtractedReference(node1, action1, nil, envs.NilLanguage, assets.NewGlobalReference("org_name", "Org Name")),
		flows.NewExtractedReference(node1, action1, nil, envs.NilLanguage, assets.NewGroupReference("46057a92-6580-4e93-af36-2bb9c9d61e51", "Testers")),
		flows.NewExtractedReference(node1, action1, nil, envs.NilLanguage, assets.NewGroupReference("377c3101-a7fc-47b1-9136-980348e362c0", "Customers")),
		flows.NewExtractedReference(node1, action1, nil, envs.NilLanguage, assets.NewLabelReference("31c06b7c-010d-4f91-9590-d3fbdc2fb7ac", "Spam")),
		flows.NewExtractedReference(node1, action1, nil, envs.NilLanguage, assets.NewTemplateReference("ff958d30-f50e-48ab-a524-37ed1e9620d9", "Welcome")),
		flows.NewExtractedReference(node1, action1, nil, envs.NilLanguage, assets.NewTicketerReference("fb9cab80-4450-4a9d-ba9b-cb8df40dd233", "Support")),
		flows.NewExtractedReference(node1, action1, nil, envs.NilLanguage, assets.NewTopicReference("531d3fc7-64f4-4170-927d-b477e8145dd3", "Weather")),
		flows.NewExtractedReference(node1, action1, nil, envs.NilLanguage, assets.NewUserReference("jim@nyaruka.com", "Jim")),
		flows.NewExtractedReference(node2, nil, router2, envs.NilLanguage, assets.NewGlobalReference("org_name", "Org Name")),
	}

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

	sa, err := engine.NewSessionAssets(env, source, nil)
	require.NoError(t, err)

	deps := inspect.NewDependencies(refs, sa)
	depsJSON := jsonx.MustMarshal(deps)
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
		},
		{
			"missing": true,
			"name": "Support",
			"type": "ticketer",
			"uuid": "fb9cab80-4450-4a9d-ba9b-cb8df40dd233"
		},
		{
			"missing": true,
			"name": "Weather",
			"type": "topic",
			"uuid": "531d3fc7-64f4-4170-927d-b477e8145dd3"
		},
		{
			"missing": true,
			"type": "user",
			"email": "jim@nyaruka.com",
			"name": "Jim"
		}
	]`), depsJSON, "deps JSON mismatch")

	// panic if we get a dependency type we don't recognize
	assert.Panics(t, func() {
		inspect.NewDependencies([]flows.ExtractedReference{
			flows.NewExtractedReference(node1, action1, nil, envs.NilLanguage, &unknownAssetType{}),
		}, sa)
	})
}
