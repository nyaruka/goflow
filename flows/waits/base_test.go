package waits_test

import (
	"testing"

	"github.com/nyaruka/goflow/flows/waits"

	"github.com/stretchr/testify/assert"
)

func TestReadWait(t *testing.T) {
	// error if no type field
	_, err := waits.ReadWait([]byte(`{"foo": "bar"}`))
	assert.EqualError(t, err, "field 'type' is required")

	// error if we don't recognize action type
	_, err = waits.ReadWait([]byte(`{"type": "do_the_foo", "foo": "bar"}`))
	assert.EqualError(t, err, "unknown type: 'do_the_foo'")
}
