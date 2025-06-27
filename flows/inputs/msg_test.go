package inputs_test

import (
	"testing"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/flows/inputs"
	"github.com/nyaruka/goflow/test"
	"github.com/nyaruka/goflow/utils"

	"github.com/stretchr/testify/assert"
)

func TestMsgInput(t *testing.T) {
	test.MockUniverse()

	_, session, _ := test.NewSessionBuilder().MustBuild()
	env := session.Environment()

	channel := session.Assets().Channels().Get("57f1078f-88aa-46f4-a59a-948a5739c03d")

	msgEvt := events.NewMsgReceived(flows.NewMsgIn(
		flows.MsgUUID("f51d7220-10b3-4faa-a91c-1ae70beaae3e"),
		urns.URN("tel:+1234567890"),
		assets.NewChannelReference("57f1078f-88aa-46f4-a59a-948a5739c03d", "Nexmo"),
		"Hi there!",
		[]utils.Attachment{
			"image/jpg:http://example.com/test.jpg",
			"video/mp4:http://example.com/test.mp4",
		},
		"ext12345",
	))

	input := inputs.NewMsg(session.Assets(), msgEvt)
	assert.Equal(t, "msg", input.Type())
	assert.Equal(t, string(msgEvt.Msg.UUID()), string(input.UUID()))
	assert.Equal(t, channel, input.Channel())
	assert.Equal(t, msgEvt.CreatedOn(), input.CreatedOn())

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
	assert.Equal(t, `{"type":"msg","uuid":"f51d7220-10b3-4faa-a91c-1ae70beaae3e","channel":{"uuid":"57f1078f-88aa-46f4-a59a-948a5739c03d","name":"My Android Phone"},"created_on":"2025-05-04T12:31:10.123456789Z","urn":"tel:+1234567890","text":"Hi there!","attachments":["image/jpg:http://example.com/test.jpg","video/mp4:http://example.com/test.mp4"],"external_id":"ext12345"}`, string(marshaled))
}
