package triggers_test

import (
	"testing"

	"github.com/nyaruka/goflow/flows/triggers"

	"github.com/stretchr/testify/assert"
)

func TestReadTrigger(t *testing.T) {
	// error if no type field
	_, err := triggers.ReadTrigger(nil, []byte(`{"foo": "bar"}`))
	assert.EqualError(t, err, "field 'type' is required")

	// error if we don't recognize action type
	_, err = triggers.ReadTrigger(nil, []byte(`{"type": "do_the_foo", "foo": "bar"}`))
	assert.EqualError(t, err, "unknown type: 'do_the_foo'")
}
