package flows_test

import (
	"testing"

	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/static"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/engine"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTickets(t *testing.T) {
	defer uuids.SetGenerator(uuids.DefaultGenerator)
	uuids.SetGenerator(uuids.NewSeededGenerator(12345))

	env := envs.NewBuilder().Build()

	source, err := static.NewSource([]byte(`{
		"ticketers": [
			{
				"uuid": "d605bb96-258d-4097-ad0a-080937db2212",
				"name": "Support Tickets",
				"type": "internal"
			},
			{
				"uuid": "5885ed52-8d3e-4fd3-be49-57eebe5d4d59",
				"name": "Email Tickets",
				"type": "mailgun"
			}
		],
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
				"email": "bob@nyaruka.com",
				"name": "Bob"
			}
		]
	}`))
	require.NoError(t, err)

	sa, err := engine.NewSessionAssets(env, source, nil)
	require.NoError(t, err)

	mailgun := sa.Ticketers().Get("5885ed52-8d3e-4fd3-be49-57eebe5d4d59")
	assert.Equal(t, assets.TicketerUUID("5885ed52-8d3e-4fd3-be49-57eebe5d4d59"), mailgun.UUID())
	assert.Equal(t, "Email Tickets", mailgun.Name())
	assert.Equal(t, "mailgun", mailgun.Type())
	assert.Equal(t, assets.NewTicketerReference("5885ed52-8d3e-4fd3-be49-57eebe5d4d59", "Email Tickets"), mailgun.Reference())

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
		"uuid": "349c851f-3f8e-4353-8bf2-8e90b6d73530", 
		"ticketer": {"uuid": "0a0b5ce4-35c9-47b7-b124-40258f0a5b53", "name": "Deleted"},
		"topic": {"uuid": "fd3ffcf3-c609-423e-b40f-f7f291a91cc6", "name": "Deleted"},
		"subject": "Very Old Ticket",
		"body": "Ticketer, topic and assignee gone!",
		"external_id": "7654",
		"assignee": {"email": "dave@nyaruka.com", "name": "Dave"}
	}`), missing)
	require.NoError(t, err)

	assert.Equal(t, flows.TicketUUID("349c851f-3f8e-4353-8bf2-8e90b6d73530"), ticket1.UUID())
	assert.Nil(t, ticket1.Ticketer())
	assert.Nil(t, ticket1.Topic())
	assert.Equal(t, "Ticketer, topic and assignee gone!", ticket1.Body())
	assert.Equal(t, "7654", ticket1.ExternalID())
	assert.Nil(t, ticket1.Assignee())

	// check that missing ticketer, topic and assignee are logged as a missing dependencies
	assert.Equal(t, 3, len(missingRefs))
	assert.Equal(t, "0a0b5ce4-35c9-47b7-b124-40258f0a5b53", missingRefs[0].Identity())
	assert.Equal(t, "fd3ffcf3-c609-423e-b40f-f7f291a91cc6", missingRefs[1].Identity())
	assert.Equal(t, "dave@nyaruka.com", missingRefs[2].Identity())

	missingRefs = make([]assets.Reference, 0)

	ticket2, err := flows.ReadTicket(sa, []byte(`{
		"uuid": "5a4af021-d2c2-47fc-9abc-abbb8635d8c0", 
		"ticketer": {"uuid": "d605bb96-258d-4097-ad0a-080937db2212", "name": "Support Tickets"},
		"topic": {"uuid": "472a7a73-96cb-4736-b567-056d987cc5b4", "name": "Weather"},
		"subject": "Old Ticket",
		"body": "Where are my shoes?",
		"assignee": {"email": "bob@nyaruka.com", "name": "Bob"}
	}`), missing)
	require.NoError(t, err)

	assert.Equal(t, 0, len(missingRefs))
	assert.Equal(t, "Support Tickets", ticket2.Ticketer().Name())
	assert.Equal(t, "Bob", ticket2.Assignee().Name())

	tickets := flows.NewTicketList([]*flows.Ticket{ticket1, ticket2})
	assert.Equal(t, 2, tickets.Count())
	assert.Equal(t, flows.TicketUUID("349c851f-3f8e-4353-8bf2-8e90b6d73530"), tickets.All()[0].UUID())
	assert.Equal(t, flows.TicketUUID("5a4af021-d2c2-47fc-9abc-abbb8635d8c0"), tickets.All()[1].UUID())

	ticket3 := flows.OpenTicket(mailgun, weather, "Where are my pants?", bob)
	ticket3.SetExternalID("24567")

	assert.Equal(t, flows.TicketUUID("1ae96956-4b34-433e-8d1a-f05fe6923d6d"), ticket3.UUID())
	assert.Equal(t, mailgun, ticket3.Ticketer())
	assert.Equal(t, weather, ticket3.Topic())
	assert.Equal(t, "Where are my pants?", ticket3.Body())
	assert.Equal(t, "24567", ticket3.ExternalID())
	assert.Equal(t, "Bob", ticket2.Assignee().Name())

	tickets.Add(ticket3)
	assert.Equal(t, 3, tickets.Count())
}
