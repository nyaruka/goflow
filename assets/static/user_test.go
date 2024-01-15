package static_test

import (
	"testing"

	"github.com/nyaruka/goflow/assets/static"

	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {
	user := static.NewUser("bob@nyaruka.com", "Bob")
	assert.Equal(t, "bob@nyaruka.com", user.Email())
	assert.Equal(t, "Bob", user.Name())
}
