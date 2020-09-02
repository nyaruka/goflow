package hints_test

import (
	"testing"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/flows/routers/waits/hints"

	"github.com/stretchr/testify/assert"
)

func TestReadHint(t *testing.T) {
	// error if no type field
	_, err := hints.ReadHint([]byte(`{"foo": "bar"}`))
	assert.EqualError(t, err, "field 'type' is required")

	// error if we don't recognize the type
	_, err = hints.ReadHint([]byte(`{"type": "do_the_foo", "foo": "bar"}`))
	assert.EqualError(t, err, "unknown type: 'do_the_foo'")

	// read image hint
	hint, err := hints.ReadHint([]byte(`{"type": "image"}`))
	assert.NoError(t, err)
	assert.Equal(t, "image", hint.Type())

	// marshal back to JSON
	data, err := jsonx.Marshal(hint)
	assert.NoError(t, err)
	assert.Equal(t, `{"type":"image"}`, string(data))

	// read video hint
	hint, err = hints.ReadHint([]byte(`{"type": "video"}`))
	assert.NoError(t, err)
	assert.Equal(t, "video", hint.Type())

	// marshal back to JSON
	data, err = jsonx.Marshal(hint)
	assert.NoError(t, err)
	assert.Equal(t, `{"type":"video"}`, string(data))

	// read audio hint
	hint, err = hints.ReadHint([]byte(`{"type": "audio"}`))
	assert.NoError(t, err)
	assert.Equal(t, "audio", hint.Type())

	// marshal back to JSON
	data, err = jsonx.Marshal(hint)
	assert.NoError(t, err)
	assert.Equal(t, `{"type":"audio"}`, string(data))

	// read location hint
	hint, err = hints.ReadHint([]byte(`{"type": "location"}`))
	assert.NoError(t, err)
	assert.Equal(t, "location", hint.Type())

	// marshal back to JSON
	data, err = jsonx.Marshal(hint)
	assert.NoError(t, err)
	assert.Equal(t, `{"type":"location"}`, string(data))

	// read digits hint
	hint, err = hints.ReadHint([]byte(`{"type": "digits", "count": 1}`))
	assert.NoError(t, err)
	assert.Equal(t, "digits", hint.Type())
	assert.Equal(t, 1, *hint.(*hints.DigitsHint).Count)

	// marshal back to JSON
	data, err = jsonx.Marshal(hint)
	assert.NoError(t, err)
	assert.Equal(t, `{"type":"digits","count":1}`, string(data))
}
