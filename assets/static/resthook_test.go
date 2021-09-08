package static_test

import (
	"testing"

	"github.com/nyaruka/goflow/assets/static"
	"github.com/stretchr/testify/assert"
)

func TestResthook(t *testing.T) {
	hook := static.NewResthook("new-contact", []string{"http://example.com"})
	assert.Equal(t, "new-contact", hook.Slug())
	assert.Equal(t, []string{"http://example.com"}, hook.Subscribers())
}
