package static_test

import (
	"testing"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/static"

	"github.com/stretchr/testify/assert"
)

func TestTicketer(t *testing.T) {
	ticketer := static.NewTicketer(
		assets.TicketerUUID("37657cf7-5eab-4286-9cb0-bbf270587bad"),
		"Support Tickets",
		"mailgun",
	)
	assert.Equal(t, assets.TicketerUUID("37657cf7-5eab-4286-9cb0-bbf270587bad"), ticketer.UUID())
	assert.Equal(t, "Support Tickets", ticketer.Name())
	assert.Equal(t, "mailgun", ticketer.Type())
}
