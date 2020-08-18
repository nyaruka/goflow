package flows_test

import (
	"testing"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/test"
	"github.com/nyaruka/goflow/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConnection(t *testing.T) {
	connection := flows.NewConnection(
		assets.NewChannelReference(assets.ChannelUUID("61f38f46-a856-4f90-899e-905691784159"), "My Android"),
		urns.URN("tel:+1234567890"),
	)

	// test marshaling our connection
	marshaled, err := jsonx.Marshal(connection)
	require.NoError(t, err)

	test.AssertEqualJSON(t, []byte(`{
		"channel":{"uuid":"61f38f46-a856-4f90-899e-905691784159","name":"My Android"},
		"urn":"tel:+1234567890"
	}`), marshaled, "JSON mismatch")

	// test unmarshaling
	connection = &flows.Connection{}
	err = utils.UnmarshalAndValidate(marshaled, connection)
	require.NoError(t, err)
	assert.Equal(t, assets.ChannelUUID("61f38f46-a856-4f90-899e-905691784159"), connection.Channel().UUID)
	assert.Equal(t, urns.URN("tel:+1234567890"), connection.URN())
}
