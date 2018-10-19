package routers_test

import (
	"testing"

	"github.com/nyaruka/goflow/flows/routers"

	"github.com/stretchr/testify/assert"
)

func TestReadRouter(t *testing.T) {
	// error if no type field
	_, err := routers.ReadRouter([]byte(`{"foo": "bar"}`))
	assert.EqualError(t, err, "field 'type' is required")

	// error if we don't recognize action type
	_, err = routers.ReadRouter([]byte(`{"type": "do_the_foo", "foo": "bar"}`))
	assert.EqualError(t, err, "unknown type: 'do_the_foo'")
}
