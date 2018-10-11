package types_test

import (
	"testing"

	"github.com/nyaruka/goflow/assets/static/types"

	"github.com/stretchr/testify/assert"
)

func TestResthook(t *testing.T) {
	hook := types.NewResthook("new-contact", []string{"http://example.com"})
	assert.Equal(t, "new-contact", hook.Slug())
	assert.Equal(t, []string{"http://example.com"}, hook.Subscribers())
}
