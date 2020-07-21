package flows_test

import (
	"testing"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/static"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/engine"
	"github.com/nyaruka/goflow/flows/triggers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHistory(t *testing.T) {
	env := envs.NewBuilder().Build()

	source, err := static.NewSource([]byte(`{
		"flows": [
			{
				"uuid": "5472a1c3-63e1-484f-8485-cc8ecb16a058",
				"name": "Empty",
				"spec_version": "13.1",
				"language": "eng",
				"type": "messaging",
				"nodes": []
			}
		]
	}`))
	require.NoError(t, err)

	sa, err := engine.NewSessionAssets(env, source, nil)
	require.NoError(t, err)

	flow := assets.NewFlowReference("5472a1c3-63e1-484f-8485-cc8ecb16a058", "Inception")
	contact := flows.NewEmptyContact(sa, "Bob", envs.Language("eng"), nil)

	eng := engine.NewBuilder().Build()
	session, _, err := eng.NewSession(sa, triggers.NewBuilder(env, flow, contact).Manual().Build())
	require.NoError(t, err)

	assert.Equal(t, flows.SessionUUID(""), session.History().ParentUUID)
	assert.Equal(t, 0, session.History().Ancestors)
	assert.Equal(t, 0, session.History().AncestorsSinceInput)

	child := flows.NewChildHistory(session)

	assert.Equal(t, session.UUID(), child.ParentUUID)
	assert.Equal(t, 1, child.Ancestors)
	assert.Equal(t, 1, child.AncestorsSinceInput)
}
