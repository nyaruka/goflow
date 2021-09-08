package static_test

import (
	"testing"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/static"

	"github.com/stretchr/testify/assert"
)

func TestTopic(t *testing.T) {
	topic := static.NewTopic(
		assets.TopicUUID("37657cf7-5eab-4286-9cb0-bbf270587bad"),
		"Weather",
	)
	assert.Equal(t, assets.TopicUUID("37657cf7-5eab-4286-9cb0-bbf270587bad"), topic.UUID())
	assert.Equal(t, "Weather", topic.Name())
}
