package waits_test

import (
	"testing"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/flows/routers/waits"

	"github.com/stretchr/testify/assert"
)

func TestReadWait(t *testing.T) {
	// error if no type field
	_, err := waits.ReadWait([]byte(`{"foo": "bar"}`))
	assert.EqualError(t, err, "field 'type' is required")

	// error if we don't recognize the type
	_, err = waits.ReadWait([]byte(`{"type": "do_the_foo", "foo": "bar"}`))
	assert.EqualError(t, err, "unknown type: 'do_the_foo'")

	// read msg wait without hint
	wait, err := waits.ReadWait([]byte(`{"type": "msg"}`))
	assert.NoError(t, err)
	assert.Equal(t, waits.TypeMsg, wait.Type())

	// marshal back to JSON
	data, err := jsonx.Marshal(wait)
	assert.NoError(t, err)
	assert.Equal(t, `{"type":"msg"}`, string(data))

	// read msg wait with hint
	wait, err = waits.ReadWait([]byte(`{"type": "msg", "hint": {"type": "image"}}`))
	assert.NoError(t, err)
	assert.Equal(t, "msg", wait.Type())
	assert.Equal(t, "image", wait.(*waits.MsgWait).Hint().Type())

	// marshal back to JSON
	data, err = jsonx.Marshal(wait)
	assert.NoError(t, err)
	assert.Equal(t, `{"type":"msg","hint":{"type":"image"}}`, string(data))
}
