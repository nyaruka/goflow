package static_test

import (
	"testing"

	"github.com/nyaruka/goflow/assets/static"

	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {
	ticketer := static.NewUser("bob@nyaruka.com", "Bob")
	assert.Equal(t, "bob@nyaruka.com", ticketer.Email())
	assert.Equal(t, "Bob", ticketer.Name())
}
