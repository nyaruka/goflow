package flows_test

import (
	"testing"

	"github.com/nyaruka/goflow/flows"

	"github.com/stretchr/testify/assert"
)

func TestContactQueryEscaping(t *testing.T) {
	assert.Equal(t, `""`, flows.ContactQueryEscaping(``))
	assert.Equal(t, `"bobby tables"`, flows.ContactQueryEscaping(`bobby tables`))
	assert.Equal(t, `"\"\" OR (id = 1)"`, flows.ContactQueryEscaping(`"" OR (id = 1)`))
	assert.Equal(t, `"\\\"foo"`, flows.ContactQueryEscaping(`\"foo`))
}
