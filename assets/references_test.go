package assets_test

import (
	"testing"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/utils"
	"github.com/nyaruka/goflow/utils/jsonx"
	"github.com/nyaruka/goflow/utils/uuids"

	"github.com/stretchr/testify/assert"
)

func TestReferences(t *testing.T) {
	channelRef := assets.NewChannelReference("61602f3e-f603-4c70-8a8f-c477505bf4bf", "Nexmo")
	assert.Equal(t, "channel", channelRef.Type())
	assert.Equal(t, "61602f3e-f603-4c70-8a8f-c477505bf4bf", channelRef.Identity())
	assert.Equal(t, uuids.UUID("61602f3e-f603-4c70-8a8f-c477505bf4bf"), channelRef.GenericUUID())
	assert.Equal(t, "channel[uuid=61602f3e-f603-4c70-8a8f-c477505bf4bf,name=Nexmo]", channelRef.String())
	assert.False(t, channelRef.Variable())
	assert.NoError(t, utils.Validate(channelRef))

	// channel references must always be concrete
	assert.EqualError(t, utils.Validate(assets.NewChannelReference("", "Nexmo")), "field 'uuid' is required")

	classifierRef := assets.NewClassifierReference("61602f3e-f603-4c70-8a8f-c477505bf4bf", "Booking")
	assert.Equal(t, "classifier", classifierRef.Type())
	assert.Equal(t, "61602f3e-f603-4c70-8a8f-c477505bf4bf", classifierRef.Identity())
	assert.Equal(t, uuids.UUID("61602f3e-f603-4c70-8a8f-c477505bf4bf"), classifierRef.GenericUUID())
	assert.Equal(t, "classifier[uuid=61602f3e-f603-4c70-8a8f-c477505bf4bf,name=Booking]", classifierRef.String())
	assert.False(t, classifierRef.Variable())
	assert.NoError(t, utils.Validate(classifierRef))

	// classifier references must always be concrete
	assert.EqualError(t, utils.Validate(assets.NewClassifierReference("", "Booking")), "field 'uuid' is required")

	fieldRef := assets.NewFieldReference("gender", "Gender")
	assert.Equal(t, "field", fieldRef.Type())
	assert.Equal(t, "gender", fieldRef.Identity())
	assert.Equal(t, "field[key=gender,name=Gender]", fieldRef.String())
	assert.False(t, fieldRef.Variable())
	assert.NoError(t, utils.Validate(fieldRef))

	// field references must have a key
	assert.EqualError(t, utils.Validate(assets.NewFieldReference("", "Gender")), "field 'key' is required")

	flowRef := assets.NewFlowReference("61602f3e-f603-4c70-8a8f-c477505bf4bf", "Registration")
	assert.Equal(t, "flow", flowRef.Type())
	assert.Equal(t, "61602f3e-f603-4c70-8a8f-c477505bf4bf", flowRef.Identity())
	assert.Equal(t, uuids.UUID("61602f3e-f603-4c70-8a8f-c477505bf4bf"), flowRef.GenericUUID())
	assert.Equal(t, "flow[uuid=61602f3e-f603-4c70-8a8f-c477505bf4bf,name=Registration]", flowRef.String())
	assert.False(t, flowRef.Variable())
	assert.NoError(t, utils.Validate(flowRef))

	// flow references must always be concrete
	assert.EqualError(t, utils.Validate(assets.NewFlowReference("", "Registration")), "field 'uuid' is required")

	globalRef := assets.NewGlobalReference("org_name", "Org Name")
	assert.Equal(t, "global", globalRef.Type())
	assert.Equal(t, "org_name", globalRef.Identity())
	assert.Equal(t, "global[key=org_name,name=Org Name]", globalRef.String())
	assert.False(t, globalRef.Variable())
	assert.NoError(t, utils.Validate(globalRef))

	// global references must have a key
	assert.EqualError(t, utils.Validate(assets.NewGlobalReference("", "Org Name")), "field 'key' is required")

	groupRef := assets.NewGroupReference("61602f3e-f603-4c70-8a8f-c477505bf4bf", "Testers")
	assert.Equal(t, "group", groupRef.Type())
	assert.Equal(t, "61602f3e-f603-4c70-8a8f-c477505bf4bf", groupRef.Identity())
	assert.Equal(t, uuids.UUID("61602f3e-f603-4c70-8a8f-c477505bf4bf"), groupRef.GenericUUID())
	assert.Equal(t, "group[uuid=61602f3e-f603-4c70-8a8f-c477505bf4bf,name=Testers]", groupRef.String())
	assert.False(t, groupRef.Variable())
	assert.NoError(t, utils.Validate(groupRef))

	// group references can be concrete or a name match template
	assert.NoError(t, utils.Validate(assets.NewVariableGroupReference("@contact.fields.district")))

	// but they can't be neither or both of those things
	assert.EqualError(t,
		utils.Validate(&assets.GroupReference{}),
		"field 'uuid' is mutually exclusive with 'name_match', field 'name_match' is mutually exclusive with 'uuid'",
	)
	assert.EqualError(t,
		utils.Validate(&assets.GroupReference{
			UUID: "61602f3e-f603-4c70-8a8f-c477505bf4bf",
			Name: "Bob", NameMatch: "@contact.fields.district"}),
		"field 'uuid' is mutually exclusive with 'name_match', field 'name_match' is mutually exclusive with 'uuid'",
	)

	labelRef := assets.NewLabelReference("61602f3e-f603-4c70-8a8f-c477505bf4bf", "Spam")
	assert.Equal(t, "label", labelRef.Type())
	assert.Equal(t, "61602f3e-f603-4c70-8a8f-c477505bf4bf", labelRef.Identity())
	assert.Equal(t, uuids.UUID("61602f3e-f603-4c70-8a8f-c477505bf4bf"), labelRef.GenericUUID())
	assert.Equal(t, "label[uuid=61602f3e-f603-4c70-8a8f-c477505bf4bf,name=Spam]", labelRef.String())
	assert.False(t, labelRef.Variable())
	assert.NoError(t, utils.Validate(labelRef))

	// label references can be concrete or a name match template
	assert.NoError(t, utils.Validate(assets.NewVariableLabelReference("@contact.fields.district")))

	// but they can't be neither or both of those things
	assert.EqualError(t,
		utils.Validate(&assets.LabelReference{}),
		"field 'uuid' is mutually exclusive with 'name_match', field 'name_match' is mutually exclusive with 'uuid'",
	)
	assert.EqualError(t,
		utils.Validate(&assets.LabelReference{UUID: "61602f3e-f603-4c70-8a8f-c477505bf4bf", Name: "Spam", NameMatch: "@contact.fields.district"}),
		"field 'uuid' is mutually exclusive with 'name_match', field 'name_match' is mutually exclusive with 'uuid'",
	)

	templateRef := assets.NewTemplateReference("61602f3e-f603-4c70-8a8f-c477505bf4bf", "Affirmation")
	assert.Equal(t, "template", templateRef.Type())
	assert.Equal(t, "61602f3e-f603-4c70-8a8f-c477505bf4bf", templateRef.Identity())
	assert.Equal(t, uuids.UUID("61602f3e-f603-4c70-8a8f-c477505bf4bf"), templateRef.GenericUUID())
	assert.Equal(t, "template[uuid=61602f3e-f603-4c70-8a8f-c477505bf4bf,name=Affirmation]", templateRef.String())
	assert.False(t, templateRef.Variable())
	assert.NoError(t, utils.Validate(templateRef))

	ticketerRef := assets.NewTicketerReference("61602f3e-f603-4c70-8a8f-c477505bf4bf", "Support Tickets")
	assert.Equal(t, "ticketer", ticketerRef.Type())
	assert.Equal(t, "61602f3e-f603-4c70-8a8f-c477505bf4bf", ticketerRef.Identity())
	assert.Equal(t, uuids.UUID("61602f3e-f603-4c70-8a8f-c477505bf4bf"), ticketerRef.GenericUUID())
	assert.Equal(t, "ticketer[uuid=61602f3e-f603-4c70-8a8f-c477505bf4bf,name=Support Tickets]", ticketerRef.String())
	assert.False(t, ticketerRef.Variable())
	assert.NoError(t, utils.Validate(ticketerRef))

	// ticketer references must always be concrete
	assert.EqualError(t, utils.Validate(assets.NewTicketerReference("", "Booking")), "field 'uuid' is required")
}

func TestChannelReferenceUnmarsal(t *testing.T) {
	// check that UUIDs aren't required to be valid UUID4s
	channel := &assets.ChannelReference{}
	err := utils.UnmarshalAndValidate([]byte(`{"uuid": "ffffffff-9b24-92e1-ffff-ffffb207cdb4", "name": "Old Channel"}`), channel)
	assert.NoError(t, err)
	assert.Equal(t, assets.ChannelUUID("ffffffff-9b24-92e1-ffff-ffffb207cdb4"), channel.UUID)
	assert.Equal(t, "Old Channel", channel.Name)
}

func TestTypedReference(t *testing.T) {
	ref := assets.NewGroupReference("61602f3e-f603-4c70-8a8f-c477505bf4bf", "Bobs")
	typed := assets.NewTypedReference(ref)

	refJSON, _ := jsonx.Marshal(ref)
	typedJSON, _ := jsonx.Marshal(typed)

	assert.Equal(t, `{"uuid":"61602f3e-f603-4c70-8a8f-c477505bf4bf","name":"Bobs"}`, string(refJSON))
	assert.Equal(t, `{"uuid":"61602f3e-f603-4c70-8a8f-c477505bf4bf","name":"Bobs","type":"group"}`, string(typedJSON))
}
