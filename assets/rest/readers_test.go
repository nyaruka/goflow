package rest

import (
	"testing"

	"github.com/nyaruka/goflow/assets"

	"github.com/stretchr/testify/assert"
)

func TestReadChannels(t *testing.T) {
	asset, err := readChannels([]byte(`[
		{
            "uuid": "57f1078f-88aa-46f4-a59a-948a5739c03d",
            "name": "My Android Phone",
            "address": "+12345671111",
            "schemes": [
                "tel"
            ],
            "roles": [
                "send",
                "receive"
            ]
        },
        {
            "uuid": "3a05eaf5-cb1b-4246-bef1-f277419c83a7",
            "name": "Nexmo",
            "address": "+12345672222",
            "schemes": [
                "tel"
            ],
            "roles": [
                "send",
                "receive"
            ]
        }
	]`))
	assert.NoError(t, err)

	channels := asset.([]assets.Channel)
	assert.Equal(t, 2, len(channels))
	assert.Equal(t, assets.ChannelUUID("57f1078f-88aa-46f4-a59a-948a5739c03d"), channels[0].UUID())
	assert.Equal(t, "My Android Phone", channels[0].Name())
	assert.Equal(t, "+12345671111", channels[0].Address())
	assert.Equal(t, []string{"tel"}, channels[0].Schemes())
	assert.Equal(t, []assets.ChannelRole{assets.ChannelRoleSend, assets.ChannelRoleReceive}, channels[0].Roles())
}

func TestReadFields(t *testing.T) {
	asset, err := readFields([]byte(`[
		{
            "key": "gender",
            "name": "Gender",
            "value_type": "text"
        },
        {
            "key": "age",
            "name": "Age",
            "value_type": "number"
        }
	]`))
	assert.NoError(t, err)

	fields := asset.([]assets.Field)
	assert.Equal(t, 2, len(fields))
	assert.Equal(t, "gender", fields[0].Key())
	assert.Equal(t, "Gender", fields[0].Name())
	assert.Equal(t, assets.FieldType("text"), fields[0].Type())
}

func TestReadGroups(t *testing.T) {
	asset, err := readGroups([]byte(`[
		{
            "uuid": "1e1ce1e1-9288-4504-869e-022d1003c72a",
            "name": "Customers"
        },
        {
            "uuid": "0ec97956-c451-48a0-a180-1ce766623e31",
            "name": "Males",
            "query": "gender = male"
        }
	]`))
	assert.NoError(t, err)

	groups := asset.([]assets.Group)
	assert.Equal(t, 2, len(groups))
	assert.Equal(t, assets.GroupUUID("1e1ce1e1-9288-4504-869e-022d1003c72a"), groups[0].UUID())
	assert.Equal(t, "Customers", groups[0].Name())
	assert.Equal(t, "", groups[0].Query())
}

func TestReadLabels(t *testing.T) {
	asset, err := readLabels([]byte(`[
		{
            "uuid": "3f65d88a-95dc-4140-9451-943e94e06fea",
            "name": "Spam"
		},
		{
            "uuid": "0ec97956-c451-48a0-a180-1ce766623e31",
            "name": "Important"
        }
	]`))
	assert.NoError(t, err)

	labels := asset.([]assets.Label)
	assert.Equal(t, 2, len(labels))
	assert.Equal(t, assets.LabelUUID("3f65d88a-95dc-4140-9451-943e94e06fea"), labels[0].UUID())
	assert.Equal(t, "Spam", labels[0].Name())
}

func TestReadResthooks(t *testing.T) {
	asset, err := readResthooks([]byte(`[
		{
            "slug": "new-registration",
            "subscribers": [
                "http://localhost/?cmd=success",
                "http://localhost/?cmd=unavailable"
            ]
		},
		{
			"slug": "end-registration",
			"subscribers": []
		}
	]`))
	assert.NoError(t, err)

	resthooks := asset.([]assets.Resthook)
	assert.Equal(t, 2, len(resthooks))
	assert.Equal(t, "new-registration", resthooks[0].Slug())
	assert.Equal(t, []string{"http://localhost/?cmd=success", "http://localhost/?cmd=unavailable"}, resthooks[0].Subscribers())
}
