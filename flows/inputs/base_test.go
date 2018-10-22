package inputs_test

import (
	"testing"

	"github.com/nyaruka/goflow/flows/inputs"

	"github.com/stretchr/testify/assert"
)

func TestReadInput(t *testing.T) {
	// error if no type field
	_, err := inputs.ReadInput(nil, []byte(`{"foo": "bar"}`))
	assert.EqualError(t, err, "field 'type' is required")

	// error if we don't recognize action type
	_, err = inputs.ReadInput(nil, []byte(`{"type": "do_the_foo", "foo": "bar"}`))
	assert.EqualError(t, err, "unknown type: 'do_the_foo'")
}
