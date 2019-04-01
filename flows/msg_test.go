package flows_test

import (
	"encoding/json"
	"testing"

	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/test"
	"github.com/nyaruka/goflow/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMsgIn(t *testing.T) {
	msg := flows.NewMsgIn(
		flows.MsgUUID("48c32bd4-ed68-4a21-b540-9da96217b022"),
		urns.URN("tel:+1234567890"),
		assets.NewChannelReference(assets.ChannelUUID("61f38f46-a856-4f90-899e-905691784159"), "My Android"),
		"Hi there",
		[]utils.Attachment{
			utils.Attachment("image/jpeg:https://example.com/test.jpg"),
			utils.Attachment("audio/mp3:https://example.com/test.mp3"),
		},
	)
	msg.SetID(123)
	msg.SetExternalID("EX346436734")

	// test marshaling our msg
	marshaled, err := json.Marshal(msg)
	require.NoError(t, err)

	test.AssertEqualJSON(t, []byte(`{
		"uuid":"48c32bd4-ed68-4a21-b540-9da96217b022",
		"id":123,"urn":"tel:+1234567890",
		"channel":{"uuid":"61f38f46-a856-4f90-899e-905691784159",
		"name":"My Android"},
		"text":"Hi there",
		"attachments":["image/jpeg:https://example.com/test.jpg",
		"audio/mp3:https://example.com/test.mp3"],
		"external_id":"EX346436734"
	}`), marshaled, "JSON mismatch")

	// test unmarshaling
	msg = &flows.MsgIn{}
	err = utils.UnmarshalAndValidate(marshaled, msg)
	require.NoError(t, err)
	assert.Equal(t, flows.MsgUUID("48c32bd4-ed68-4a21-b540-9da96217b022"), msg.UUID())
	assert.Equal(t, flows.MsgID(123), msg.ID())
	assert.Equal(t, urns.URN("tel:+1234567890"), msg.URN())
	assert.Equal(t, "Hi there", msg.Text())
	assert.Equal(t, assets.ChannelUUID("61f38f46-a856-4f90-899e-905691784159"), msg.Channel().UUID)
	assert.Equal(t, "My Android", msg.Channel().Name)
	assert.Equal(t, "EX346436734", msg.ExternalID())
}
