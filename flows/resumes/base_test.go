package resumes_test

import (
	"testing"

	"github.com/nyaruka/goflow/flows/resumes"

	"github.com/stretchr/testify/assert"
)

func TestReadResume(t *testing.T) {
	// error if no type field
	_, err := resumes.ReadResume(nil, []byte(`{"foo": "bar"}`))
	assert.EqualError(t, err, "field 'type' is required")

	// error if we don't recognize action type
	_, err = resumes.ReadResume(nil, []byte(`{"type": "do_the_foo", "foo": "bar"}`))
	assert.EqualError(t, err, "unknown type: 'do_the_foo'")
}
