package simple_test

import (
	"testing"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/simple"

	"github.com/stretchr/testify/assert"
)

func TestReadChannel(t *testing.T) {
	// check that UUIDs aren't required to be valid UUID4s
	channel, err := simple.ReadChannel([]byte(`{"uuid": "ffffffff-9b24-92e1-ffff-ffffb207cdb4", "name": "Old Channel", "schemes": ["tel"], "roles": ["send"]}`))
	assert.NoError(t, err)
	assert.Equal(t, assets.ChannelUUID("ffffffff-9b24-92e1-ffff-ffffb207cdb4"), channel.UUID())
	assert.Equal(t, "Old Channel", channel.Name())
	assert.Equal(t, []assets.ChannelRole{assets.ChannelRoleSend}, channel.Roles())
}
