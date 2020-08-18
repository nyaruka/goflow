package inputs_test

import (
	"testing"
	"time"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/inputs"
	"github.com/nyaruka/goflow/test"
	"github.com/nyaruka/goflow/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMsgInput(t *testing.T) {
	session, _, err := test.CreateTestSession("", envs.RedactionPolicyNone)
	require.NoError(t, err)

	env := session.Environment()

	channel := session.Assets().Channels().Get("57f1078f-88aa-46f4-a59a-948a5739c03d")

	msg := flows.NewMsgIn(
		flows.MsgUUID("f51d7220-10b3-4faa-a91c-1ae70beaae3e"),
		urns.URN("tel:+1234567890"),
		assets.NewChannelReference("57f1078f-88aa-46f4-a59a-948a5739c03d", "Nexmo"),
		"Hi there!",
		[]utils.Attachment{
			"image/jpg:http://example.com/test.jpg",
			"video/mp4:http://example.com/test.mp4",
		},
	)
	msg.SetExternalID("ext12345")

	input := inputs.NewMsg(session.Assets(), msg, time.Date(2018, 10, 22, 16, 12, 30, 123456, time.UTC))
	assert.Equal(t, "msg", input.Type())
	assert.Equal(t, flows.InputUUID("f51d7220-10b3-4faa-a91c-1ae70beaae3e"), input.UUID())
	assert.Equal(t, channel, input.Channel())
	assert.Equal(t, time.Date(2018, 10, 22, 16, 12, 30, 123456, time.UTC), input.CreatedOn())

	// check use in expressions
	test.AssertXEqual(t, types.NewXObject(map[string]types.XValue{
		"__default__": types.NewXText("Hi there!\nhttp://example.com/test.jpg\nhttp://example.com/test.mp4"),
		"type":        types.NewXText("msg"),
		"uuid":        types.NewXText("f51d7220-10b3-4faa-a91c-1ae70beaae3e"),
		"channel":     flows.Context(env, channel),
		"created_on":  types.NewXDateTime(input.CreatedOn()),
		"urn":         types.NewXText("tel:+1234567890"),
		"text":        types.NewXText("Hi there!"),
		"attachments": types.NewXArray(types.NewXText("image/jpg:http://example.com/test.jpg"), types.NewXText("video/mp4:http://example.com/test.mp4")),
		"external_id": types.NewXText("ext12345"),
	}), flows.Context(env, input))

	// check marshaling to JSON
	marshaled, err := jsonx.Marshal(input)
	assert.NoError(t, err)
	assert.Equal(t, `{"type":"msg","uuid":"f51d7220-10b3-4faa-a91c-1ae70beaae3e","channel":{"uuid":"57f1078f-88aa-46f4-a59a-948a5739c03d","name":"My Android Phone"},"created_on":"2018-10-22T16:12:30.000123456Z","urn":"tel:+1234567890","text":"Hi there!","attachments":["image/jpg:http://example.com/test.jpg","video/mp4:http://example.com/test.mp4"],"external_id":"ext12345"}`, string(marshaled))
}
