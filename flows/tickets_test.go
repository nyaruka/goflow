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
	computers := sa.Topics().Get("daa356b6-32af-44f0-9d35-6126d55ec3e9")
	assert.Equal(t, assets.TopicUUID("472a7a73-96cb-4736-b567-056d987cc5b4"), weather.UUID())
	assert.Equal(t, "Weather", weather.Name())
	assert.Equal(t, assets.NewTopicReference("472a7a73-96cb-4736-b567-056d987cc5b4", "Weather"), weather.Reference())

	assert.Equal(t, weather, sa.Topics().FindByName("Weather"))
	assert.Equal(t, weather, sa.Topics().FindByName("WEATHER"))
	assert.Nil(t, sa.Topics().FindByName("Not Real"))

	bob := sa.Users().Get("0c78ef47-7d56-44d8-8f57-96e0f30e8f44")

	// nil object returns nil reference
	assert.Nil(t, (*flows.Topic)(nil).Reference())

	missingRefs := make([]assets.Reference, 0)
	missing := func(ref assets.Reference, err error) {
		missingRefs = append(missingRefs, ref)
	}

	ticket1 := (&flows.TicketEnvelope{
		UUID:     flows.TicketUUID("0196a645-3f8d-7452-8d1a-f05fe6923d6d"),
		Topic:    assets.NewTopicReference("fd3ffcf3-c609-423e-b40f-f7f291a91cc6", "Missing Topic"),
		Assignee: assets.NewUserReference("b8cfc330-4634-45d1-90bc-7b4658221834", "Dave"),
	}).Unmarshal(sa, missing)

	assert.Equal(t, flows.TicketUUID("0196a645-3f8d-7452-8d1a-f05fe6923d6d"), ticket1.UUID())
	assert.Nil(t, ticket1.Topic())
	assert.Nil(t, ticket1.Assignee())

	// check that missing topic and assignee are logged as a missing dependencies
	assert.Equal(t, 2, len(missingRefs))
	assert.Equal(t, "fd3ffcf3-c609-423e-b40f-f7f291a91cc6", missingRefs[0].Identity())
	assert.Equal(t, "b8cfc330-4634-45d1-90bc-7b4658221834", missingRefs[1].Identity())

	missingRefs = make([]assets.Reference, 0)

	ticket2 := (&flows.TicketEnvelope{
		UUID:     flows.TicketUUID("5a4af021-d2c2-47fc-9abc-abbb8635d8c0"),
		Topic:    weather.Reference(),
		Assignee: bob.Reference(),
	}).Unmarshal(sa, missing)

	assert.Equal(t, 0, len(missingRefs))
	assert.Equal(t, "Bob", ticket2.Assignee().Name())

	ticket3 := flows.OpenTicket(weather, bob)

	assert.Equal(t, flows.TicketUUID("01969b47-0583-76f8-ae7f-f8b243c49ff5"), ticket3.UUID())
	assert.Equal(t, weather, ticket3.Topic())
	assert.Equal(t, "Bob", ticket2.Assignee().Name())

	prevLastActivity := ticket3.LastActivityOn()

	ticket3.SetStatus(flows.TicketStatusClosed)
	ticket3.SetTopic(computers)
	ticket3.SetAssignee(nil)
	ticket3.RecordActivity()

	assert.Equal(t, flows.TicketStatusClosed, ticket3.Status())
	assert.Equal(t, computers, ticket3.Topic())
	assert.Nil(t, ticket3.Assignee())
	assert.True(t, ticket3.LastActivityOn().After(prevLastActivity))
}
