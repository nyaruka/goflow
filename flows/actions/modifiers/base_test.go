package modifiers_test

import (
	"encoding/json"
	"testing"

	"github.com/nyaruka/goflow/flows/actions/modifiers"

	"github.com/stretchr/testify/assert"
)

func TestReadModifier(t *testing.T) {
	// error if no type field
	_, err := modifiers.ReadModifier(nil, []byte(`{"foo": "bar"}`))
	assert.EqualError(t, err, "field 'type' is required")

	// error if we don't recognize the type
	_, err = modifiers.ReadModifier(nil, []byte(`{"type": "do_the_foo", "foo": "bar"}`))
	assert.EqualError(t, err, "unknown type: 'do_the_foo'")

	// read name modifier
	mod, err := modifiers.ReadModifier(nil, []byte(`{"type": "name", "name": "Bob"}`))
	assert.NoError(t, err)
	assert.Equal(t, "name", mod.Type())

	// marshal back to JSON
	data, err := json.Marshal(mod)
	assert.NoError(t, err)
	assert.Equal(t, `{"type":"name","name":"Bob"}`, string(data))
}
