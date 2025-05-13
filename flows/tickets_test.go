package flows_test

import (
	"testing"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/static"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/engine"
	"github.com/nyaruka/goflow/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTickets(t *testing.T) {
	test.MockUniverse()

	env := envs.NewBuilder().Build()

	source, err := static.NewSource([]byte(`{
		"topics": [
			{
				"uuid": "472a7a73-96cb-4736-b567-056d987cc5b4",
				"name": "Weather"
			},
			{
				"uuid": "daa356b6-32af-44f0-9d35-6126d55ec3e9",
				"name": "Computers"
			}
		],
		"users": [
			{
				"uuid": "0c78ef47-7d56-44d8-8f57-96e0f30e8f44",
				"name": "Bob", 
				"email": "bob@nyaruka.com"
			}
		]
	}`))
	require.NoError(t, err)

	sa, err := engine.NewSessionAssets(env, source, nil)
	require.NoError(t, err)

	weather := sa.Topics().Get("472a7a73-96cb-4736-b567-056d987cc5b4")
	assert.Equal(t, assets.TopicUUID("472a7a73-96cb-4736-b567-056d987cc5b4"), weather.UUID())
	assert.Equal(t, "Weather", weather.Name())
	assert.Equal(t, assets.NewTopicReference("472a7a73-96cb-4736-b567-056d987cc5b4", "Weather"), weather.Reference())

	assert.Equal(t, weather, sa.Topics().FindByName("Weather"))
	assert.Equal(t, weather, sa.Topics().FindByName("WEATHER"))
	assert.Nil(t, sa.Topics().FindByName("Not Real"))

	bob := sa.Users().Get("bob@nyaruka.com")

	// nil object returns nil reference
	assert.Nil(t, (*flows.Topic)(nil).Reference())

	missingRefs := make([]assets.Reference, 0)
	missing := func(ref assets.Reference, err error) {
		missingRefs = append(missingRefs, ref)
	}

	_, err = flows.ReadTicket(sa, []byte(`{}`), missing)
	assert.EqualError(t, err, "field 'uuid' is required")

	ticket1, err := flows.ReadTicket(sa, []byte(`{
		"uuid": "0196a645-3f8d-7452-8d1a-f05fe6923d6d", 
		"topic": {"uuid": "fd3ffcf3-c609-423e-b40f-f7f291a91cc6", "name": "Deleted"},
		"assignee": {"email": "dave@nyaruka.com", "name": "Dave"}
	}`), missing)
	require.NoError(t, err)

	assert.Equal(t, flows.TicketUUID("0196a645-3f8d-7452-8d1a-f05fe6923d6d"), ticket1.UUID())
	assert.Nil(t, ticket1.Topic())
	assert.Nil(t, ticket1.Assignee())

	// check that missing topic and assignee are logged as a missing dependencies
	assert.Equal(t, 2, len(missingRefs))
	assert.Equal(t, "fd3ffcf3-c609-423e-b40f-f7f291a91cc6", missingRefs[0].Identity())
	assert.Equal(t, "dave@nyaruka.com", missingRefs[1].Identity())

	missingRefs = make([]assets.Reference, 0)

	ticket2, err := flows.ReadTicket(sa, []byte(`{
		"uuid": "5a4af021-d2c2-47fc-9abc-abbb8635d8c0", 
		"topic": {"uuid": "472a7a73-96cb-4736-b567-056d987cc5b4", "name": "Weather"},
		"assignee": {"email": "bob@nyaruka.com", "name": "Bob"}
	}`), missing)
	require.NoError(t, err)

	assert.Equal(t, 0, len(missingRefs))
	assert.Equal(t, "Bob", ticket2.Assignee().Name())

	ticket3 := flows.OpenTicket(weather, bob)

	assert.Equal(t, flows.TicketUUID("01969b47-0583-76f8-ae7f-f8b243c49ff5"), ticket3.UUID())
	assert.Equal(t, weather, ticket3.Topic())
	assert.Equal(t, "Bob", ticket2.Assignee().Name())
}
