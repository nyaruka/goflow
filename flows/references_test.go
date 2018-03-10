package flows

import (
	"testing"

	"github.com/nyaruka/goflow/utils"
	"github.com/stretchr/testify/assert"
)

func TestReferenceValidation(t *testing.T) {
	// channel references must always be concrete
	assert.NoError(t, utils.Validate(&ChannelReference{UUID: "61602f3e-f603-4c70-8a8f-c477505bf4bf", Name: "Nexmo"}))
	assert.EqualError(t, utils.Validate(&ChannelReference{Name: "Nexmo"}), "field 'uuid' is required")

	// contact references must always be concrete
	assert.NoError(t, utils.Validate(&ContactReference{UUID: "61602f3e-f603-4c70-8a8f-c477505bf4bf", Name: "Bob"}))
	assert.EqualError(t, utils.Validate(&ContactReference{Name: "Bob"}), "field 'uuid' is required")

	// group references can be concrete or a name match template
	assert.NoError(t, utils.Validate(&GroupReference{UUID: "61602f3e-f603-4c70-8a8f-c477505bf4bf", Name: "Testers"}))
	assert.NoError(t, utils.Validate(&GroupReference{NameMatch: "@contact.fields.district"}))

	// but they can't be neither or both of those things
	assert.EqualError(t,
		utils.Validate(&GroupReference{}),
		"field 'uuid' is mutually exclusive with 'name_match', field 'name_match' is mutually exclusive with 'uuid'",
	)
	assert.EqualError(t,
		utils.Validate(&GroupReference{UUID: "61602f3e-f603-4c70-8a8f-c477505bf4bf", Name: "Bob", NameMatch: "@contact.fields.district"}),
		"field 'uuid' is mutually exclusive with 'name_match', field 'name_match' is mutually exclusive with 'uuid'",
	)

	// label references can be concrete or a name match template
	assert.NoError(t, utils.Validate(&LabelReference{UUID: "61602f3e-f603-4c70-8a8f-c477505bf4bf", Name: "Spam"}))
	assert.NoError(t, utils.Validate(&LabelReference{NameMatch: "@contact.fields.district"}))

	// but they can't be neither or both of those things
	assert.EqualError(t,
		utils.Validate(&LabelReference{}),
		"field 'uuid' is mutually exclusive with 'name_match', field 'name_match' is mutually exclusive with 'uuid'",
	)
	assert.EqualError(t,
		utils.Validate(&LabelReference{UUID: "61602f3e-f603-4c70-8a8f-c477505bf4bf", Name: "Spam", NameMatch: "@contact.fields.district"}),
		"field 'uuid' is mutually exclusive with 'name_match', field 'name_match' is mutually exclusive with 'uuid'",
	)
}
