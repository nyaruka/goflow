package inputs_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/inputs"
	"github.com/nyaruka/goflow/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMsgInput(t *testing.T) {
	session, _, err := test.CreateTestSession("", nil)
	require.NoError(t, err)

	env := session.Environment()

	channel, err := session.Assets().Channels().Get("57f1078f-88aa-46f4-a59a-948a5739c03d")
	require.NoError(t, err)

	msg := flows.NewMsgIn(
		flows.MsgUUID("f51d7220-10b3-4faa-a91c-1ae70beaae3e"),
		urns.URN("tel:+1234567890"),
		assets.NewChannelReference("57f1078f-88aa-46f4-a59a-948a5739c03d", "Nexmo"),
		"Hi there!",
		[]flows.Attachment{
			"image/jpg:http://example.com/test.jpg",
			"video/mp4:http://example.com/test.mp4",
		},
	)
	input, err := inputs.NewMsgInput(session.Assets(), msg, time.Date(2018, 10, 22, 16, 12, 30, 123456, time.UTC))
	require.NoError(t, err)

	assert.Equal(t, "msg", input.Type())
	assert.Equal(t, flows.InputUUID("f51d7220-10b3-4faa-a91c-1ae70beaae3e"), input.UUID())
	assert.Equal(t, channel, input.Channel())
	assert.Equal(t, time.Date(2018, 10, 22, 16, 12, 30, 123456, time.UTC), input.CreatedOn())

	// check use in expressions
	assert.Equal(t, "input", input.Describe())
	assert.Equal(t, types.NewXText("Hi there!\nhttp://example.com/test.jpg\nhttp://example.com/test.mp4"), input.Reduce(env))
	assert.Equal(t, types.NewXText("Hi there!"), input.Resolve(env, "text"))
	assert.Equal(t, channel, input.Resolve(env, "channel"))
	assert.Equal(t, flows.AttachmentList{"image/jpg:http://example.com/test.jpg", "video/mp4:http://example.com/test.mp4"}, input.Resolve(env, "attachments"))
	assert.Equal(t, types.NewXText(`{"attachments":[{"content_type":"image/jpg","url":"http://example.com/test.jpg"},{"content_type":"video/mp4","url":"http://example.com/test.mp4"}],"channel":{"address":"+12345671111","name":"My Android Phone","uuid":"57f1078f-88aa-46f4-a59a-948a5739c03d"},"created_on":"2018-10-22T16:12:30.000123Z","text":"Hi there!","type":"msg","urn":{"display":"234567890","path":"+1234567890","scheme":"tel"},"uuid":"f51d7220-10b3-4faa-a91c-1ae70beaae3e"}`), input.ToXJSON(env))

	// check marshaling to JSON
	marshaled, err := json.Marshal(input)
	assert.NoError(t, err)
	assert.Equal(t, `{"type":"msg","uuid":"f51d7220-10b3-4faa-a91c-1ae70beaae3e","channel":{"uuid":"57f1078f-88aa-46f4-a59a-948a5739c03d","name":"My Android Phone"},"created_on":"2018-10-22T16:12:30.000123456Z","urn":"tel:+1234567890","text":"Hi there!","attachments":["image/jpg:http://example.com/test.jpg","video/mp4:http://example.com/test.mp4"]}`, string(marshaled))
}
