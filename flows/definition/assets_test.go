package definition_test

import (
	"testing"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/static"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows/engine"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var assetsJSON = `{
	"flows": [
		{
            "uuid": "76f0a02f-3b75-4b86-9064-e9195e1b3a02",
            "name": "Zig",
            "spec_version": "13.0",
            "language": "eng",
            "type": "messaging",
            "nodes": []
        },
		{
            "uuid": "5a0b6495-9f34-4d9f-876a-1cfc7f732307",
            "name": "Zag",
            "spec_version": "13.0",
            "language": "eng",
            "type": "messaging",
            "nodes": []
        }
	]
}`

func TestFlowAssets(t *testing.T) {
	env := envs.NewBuilder().Build()

	source, err := static.NewSource([]byte(assetsJSON))
	require.NoError(t, err)

	sa, err := engine.NewSessionAssets(env, source, nil)
	require.NoError(t, err)

	assert.Equal(t, source, sa.Source())

	// fetch flow by valid UUID
	flow, err := sa.Flows().Get("76f0a02f-3b75-4b86-9064-e9195e1b3a02")
	assert.NoError(t, err)
	assert.Equal(t, assets.FlowUUID("76f0a02f-3b75-4b86-9064-e9195e1b3a02"), flow.UUID())
	assert.Equal(t, "Zig", flow.Name())

	// fetching again with same UUID gives same object
	flow1, err := sa.Flows().Get("76f0a02f-3b75-4b86-9064-e9195e1b3a02")
	assert.NoError(t, err)
	assert.Same(t, flow1, flow)

	// and by invalid UUID
	flow, err = sa.Flows().Get("xyz")
	assert.EqualError(t, err, "no such flow with UUID 'xyz'")
	assert.Nil(t, flow)

	// and by valid name
	flow, err = sa.Flows().FindByName("zag")
	assert.NoError(t, err)
	assert.Equal(t, assets.FlowUUID("5a0b6495-9f34-4d9f-876a-1cfc7f732307"), flow.UUID())
	assert.Equal(t, "Zag", flow.Name())

	// fetching again with same name gives same object
	flow2, err := sa.Flows().FindByName("zag")
	assert.NoError(t, err)
	assert.Same(t, flow2, flow)

	flow, err = sa.Flows().FindByName("zog")
	assert.EqualError(t, err, "no such flow with name 'zog'")
	assert.Nil(t, flow)
}
