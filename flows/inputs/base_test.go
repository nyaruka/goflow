package inputs_test

import (
	"testing"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/static"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows/engine"
	"github.com/nyaruka/goflow/flows/inputs"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadInput(t *testing.T) {
	env := envs.NewBuilder().Build()

	missingAssets := make([]assets.Reference, 0)
	missing := func(a assets.Reference, err error) { missingAssets = append(missingAssets, a) }

	sessionAssets, err := engine.NewSessionAssets(env, static.NewEmptySource(), nil)
	require.NoError(t, err)

	// error if no type field
	_, err = inputs.ReadInput(sessionAssets, []byte(`{"foo": "bar"}`), missing)
	assert.EqualError(t, err, "field 'type' is required")

	// error if we don't recognize action type
	_, err = inputs.ReadInput(sessionAssets, []byte(`{"type": "do_the_foo", "foo": "bar"}`), missing)
	assert.EqualError(t, err, "unknown type: 'do_the_foo'")

	// channel is optional
	_, err = inputs.ReadInput(sessionAssets, []byte(`{"type": "msg", "text": "Hello", "created_on": "2019-01-30T11:49:30Z"}`), missing)
	assert.NoError(t, err)

	// record of missing asset if channel doesn't exist
	_, err = inputs.ReadInput(sessionAssets, []byte(`{
		"type": "msg", 
		"text": "Hello", 
		"created_on": "2019-01-30T11:49:30Z",
		"channel": {
			"uuid": "2e32a8ef-8f2c-4913-a398-362b5aff9826", 
			"name": "Foo"
		}
	}`), missing)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(missingAssets))
	assert.Equal(t, assets.NewChannelReference(assets.ChannelUUID("2e32a8ef-8f2c-4913-a398-362b5aff9826"), "Foo"), missingAssets[0])
}
