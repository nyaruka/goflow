package resumes_test

import (
	"testing"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/static"
	"github.com/nyaruka/goflow/flows/engine"
	"github.com/nyaruka/goflow/flows/resumes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadResume(t *testing.T) {
	missingAssets := make([]assets.Reference, 0)
	missing := func(a assets.Reference) { missingAssets = append(missingAssets, a) }

	sessionAssets, err := engine.NewSessionAssets(static.NewEmptySource())
	require.NoError(t, err)

	// error if no type field
	_, err = resumes.ReadResume(sessionAssets, []byte(`{"foo": "bar"}`), missing)
	assert.EqualError(t, err, "field 'type' is required")

	// error if we don't recognize action type
	_, err = resumes.ReadResume(sessionAssets, []byte(`{"type": "do_the_foo", "foo": "bar"}`), missing)
	assert.EqualError(t, err, "unknown type: 'do_the_foo'")
}
