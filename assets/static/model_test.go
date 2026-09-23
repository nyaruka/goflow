package static_test

import (
	"testing"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/static"

	"github.com/stretchr/testify/assert"
)

func TestModel(t *testing.T) {
	model := static.NewModel(
		assets.ModelUUID("37657cf7-5eab-4286-9cb0-bbf270587bad"),
		"GPT-4",
		"openai",
		[]assets.ModelRole{assets.ModelRoleEditing, assets.ModelRoleEngine},
	)
	assert.Equal(t, assets.ModelUUID("37657cf7-5eab-4286-9cb0-bbf270587bad"), model.UUID())
	assert.Equal(t, "GPT-4", model.Name())
	assert.Equal(t, "openai", model.Type())
	assert.Equal(t, []assets.ModelRole{assets.ModelRoleEditing, assets.ModelRoleEngine}, model.Roles())
}
