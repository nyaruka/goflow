package static_test

import (
	"testing"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/static"
	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {
	user := static.NewUser("0c78ef47-7d56-44d8-8f57-96e0f30e8f44", "Bob", "bob@nyaruka.com")
	assert.Equal(t, assets.UserUUID("0c78ef47-7d56-44d8-8f57-96e0f30e8f44"), user.UUID())
	assert.Equal(t, "Bob", user.Name())
	assert.Equal(t, "bob@nyaruka.com", user.Email())
}
