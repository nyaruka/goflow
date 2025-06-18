package static_test

import (
	"testing"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/static"
	"github.com/stretchr/testify/assert"
)

func TestCampaign(t *testing.T) {
	camp := static.NewCampaign(
		"58e9b092-fe42-4173-876c-ff45a14a24fe",
		"Reminders",
		assets.NewGroupReference("b9c1eaa5-4258-4f14-8aa4-ec0f5550100d", "Registered Mothers"),
	)
	assert.Equal(t, assets.CampaignUUID("58e9b092-fe42-4173-876c-ff45a14a24fe"), camp.UUID())
	assert.Equal(t, "Reminders", camp.Name())
	assert.Equal(t, "Registered Mothers", camp.Group().Name)
}
