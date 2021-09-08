package static_test

import (
	"testing"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/static"

	"github.com/stretchr/testify/assert"
)

func TestGroup(t *testing.T) {
	group := static.NewGroup(assets.GroupUUID("f5263dca-469b-47c2-be4f-845d3a14eedf"), "Spammers", "spam = yes")
	assert.Equal(t, assets.GroupUUID("f5263dca-469b-47c2-be4f-845d3a14eedf"), group.UUID())
	assert.Equal(t, "Spammers", group.Name())
	assert.Equal(t, "spam = yes", group.Query())
}
