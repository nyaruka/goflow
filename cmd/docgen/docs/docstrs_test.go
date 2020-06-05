package docs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveTypeNamePrefix(t *testing.T) {
	assert.Equal(t, "Is a contact\nHere's an example...", removeTypeNamePrefix("Contact is a contact\nHere's an example...", "Contact"))
	assert.Equal(t, "Non-standard comment...", removeTypeNamePrefix("Non-standard comment...", "Contact"))
}
