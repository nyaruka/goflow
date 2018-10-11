package assets_test

import (
	"testing"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/utils"

	"github.com/stretchr/testify/assert"
)

func TestReferenceValidation(t *testing.T) {
	// channel references must always be concrete
	assert.NoError(t, utils.Validate(assets.NewChannelReference("61602f3e-f603-4c70-8a8f-c477505bf4bf", "Nexmo")))
	assert.EqualError(t, utils.Validate(assets.NewChannelReference("", "Nexmo")), "field 'uuid' is required")

	// field references must have a key
	assert.NoError(t, utils.Validate(assets.NewFieldReference("gender", "Gender")))
	assert.EqualError(t, utils.Validate(assets.NewFieldReference("", "Gender")), "field 'key' is required")

	// flow references must always be concrete
	assert.NoError(t, utils.Validate(assets.NewFlowReference("61602f3e-f603-4c70-8a8f-c477505bf4bf", "Registration")))
	assert.EqualError(t, utils.Validate(assets.NewFlowReference("", "Registration")), "field 'uuid' is required")

	// group references can be concrete or a name match template
	assert.NoError(t, utils.Validate(assets.NewGroupReference("61602f3e-f603-4c70-8a8f-c477505bf4bf", "Testers")))
	assert.NoError(t, utils.Validate(assets.NewVariableGroupReference("@contact.fields.district")))

	// but they can't be neither or both of those things
	assert.EqualError(t,
		utils.Validate(&assets.GroupReference{}),
		"field 'uuid' is mutually exclusive with 'name_match', field 'name_match' is mutually exclusive with 'uuid'",
	)
	assert.EqualError(t,
		utils.Validate(&assets.GroupReference{UUID: "61602f3e-f603-4c70-8a8f-c477505bf4bf", Name: "Bob", NameMatch: "@contact.fields.district"}),
		"field 'uuid' is mutually exclusive with 'name_match', field 'name_match' is mutually exclusive with 'uuid'",
	)

	// label references can be concrete or a name match template
	assert.NoError(t, utils.Validate(assets.NewLabelReference("61602f3e-f603-4c70-8a8f-c477505bf4bf", "Spam")))
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
}

func TestChannelReferenceUnmarsal(t *testing.T) {
	// check that UUIDs aren't required to be valid UUID4s
	channel := &assets.ChannelReference{}
	err := utils.UnmarshalAndValidate([]byte(`{"uuid": "ffffffff-9b24-92e1-ffff-ffffb207cdb4", "name": "Old Channel"}`), channel)
	assert.NoError(t, err)
	assert.Equal(t, assets.ChannelUUID("ffffffff-9b24-92e1-ffff-ffffb207cdb4"), channel.UUID)
	assert.Equal(t, "Old Channel", channel.Name)
}
