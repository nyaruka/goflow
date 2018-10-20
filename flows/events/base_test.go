package events_test

import (
	"testing"

	"github.com/nyaruka/goflow/flows/events"

	"github.com/stretchr/testify/assert"
)

func TestEventMarshaling(t *testing.T) {

}

func TestReadEvent(t *testing.T) {
	// error if no type field
	_, err := events.ReadEvent([]byte(`{"foo": "bar"}`))
	assert.EqualError(t, err, "field 'type' is required")

	// error if we don't recognize action type
	_, err = events.ReadEvent([]byte(`{"type": "do_the_foo", "foo": "bar"}`))
	assert.EqualError(t, err, "unknown type: 'do_the_foo'")
}
