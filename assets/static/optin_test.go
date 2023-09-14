package static_test

import (
	"testing"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/static"

	"github.com/stretchr/testify/assert"
)

func TestOptIn(t *testing.T) {
	optin := static.NewOptIn(
		"37657cf7-5eab-4286-9cb0-bbf270587bad",
		"Weather Updates",
		assets.NewChannelReference("f4366920-cb05-47b9-a974-29be2d528984", "Facebook"),
	)
	assert.Equal(t, assets.OptInUUID("37657cf7-5eab-4286-9cb0-bbf270587bad"), optin.UUID())
	assert.Equal(t, "Weather Updates", optin.Name())
	assert.Equal(t, assets.NewChannelReference("f4366920-cb05-47b9-a974-29be2d528984", "Facebook"), optin.Channel())
}
