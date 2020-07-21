package flows_test

import (
	"testing"

	"github.com/nyaruka/gocommon/urns"
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
				"name": "Inception",
				"spec_version": "13.1",
				"language": "eng",
				"type": "messaging",
				"nodes": [
					{
						"uuid": "cc49453a-78ed-48a6-8b94-318b46517071",
						"actions": [
							{
								"uuid": "cdf981ae-a9cf-4c32-98f3-65bac07bf990",
								"type": "start_session",
								"flow": {
									"uuid": "5472a1c3-63e1-484f-8485-cc8ecb16a058", 
									"name": "Inception"
								},
								"contacts": [
									{
										"uuid": "51b41bf2-b2e1-439b-ab9b-dd4c9cef6266", 
										"name": "Bob"
									}
								]
							}
						],
						"exits": [
							{
								"uuid": "717ee506-7b2d-4a18-b142-eafed0c5e9d8"
							}
						]
					}
				]
			}
		]
	}`))
	require.NoError(t, err)

	sa, err := engine.NewSessionAssets(env, source, nil)
	require.NoError(t, err)

	flow := assets.NewFlowReference("5472a1c3-63e1-484f-8485-cc8ecb16a058", "Inception")

	contact := flows.NewEmptyContact(sa, "Bob", envs.Language("eng"), nil)
	contact.AddURN(urns.URN("tel:+12065551212"), nil)

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
