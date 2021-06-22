package types_test

import (
	"testing"

	"github.com/nyaruka/goflow/assets/static/types"

	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {
	ticketer := types.NewUser("bob@nyaruka.com", "Bob")
	assert.Equal(t, "bob@nyaruka.com", ticketer.Email())
	assert.Equal(t, "Bob", ticketer.Name())
}
