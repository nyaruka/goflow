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

	// nil object returns nil reference
	assert.Nil(t, (*flows.Ticketer)(nil).Reference())

	missingRefs := make([]assets.Reference, 0)
	missing := func(ref assets.Reference, err error) {
		missingRefs = append(missingRefs, ref)
	}

	_, err = flows.ReadTicket(sa, []byte(`{}`), missing)
	assert.EqualError(t, err, "field 'uuid' is required")

	ticket1, err := flows.ReadTicket(sa, []byte(`{
		"uuid": "349c851f-3f8e-4353-8bf2-8e90b6d73530", 
		"ticketer": {"uuid": "0a0b5ce4-35c9-47b7-b124-40258f0a5b53", "name": "Deleted"},
		"subject": "Very Old Ticket",
		"body": "Ticketer gone!",
		"external_id": "7654",
		"assignee": {"email": "dave@nyaruka.com", "name": "Dave"}
	}`), missing)
	require.NoError(t, err)

	assert.Equal(t, flows.TicketUUID("349c851f-3f8e-4353-8bf2-8e90b6d73530"), ticket1.UUID())
	assert.Nil(t, ticket1.Ticketer())
	assert.Equal(t, "Very Old Ticket", ticket1.Subject())
	assert.Equal(t, "Ticketer gone!", ticket1.Body())
	assert.Equal(t, "7654", ticket1.ExternalID())
	assert.Nil(t, ticket1.Assignee())

	// check that missing ticketer and assignee are logged as a missing dependencies
	assert.Equal(t, 2, len(missingRefs))
	assert.Equal(t, "0a0b5ce4-35c9-47b7-b124-40258f0a5b53", missingRefs[0].Identity())
	assert.Equal(t, "dave@nyaruka.com", missingRefs[1].Identity())

	missingRefs = make([]assets.Reference, 0)

	ticket2, err := flows.ReadTicket(sa, []byte(`{
		"uuid": "5a4af021-d2c2-47fc-9abc-abbb8635d8c0", 
		"ticketer": {"uuid": "d605bb96-258d-4097-ad0a-080937db2212", "name": "Support Tickets"},
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
	assert.Equal(t, "Very Old Ticket", tickets.All()[0].Subject())
	assert.Equal(t, "Old Ticket", tickets.All()[1].Subject())

	ticket3 := flows.OpenTicket(mailgun, "New Ticket", "Where are my pants?")
	ticket3.SetExternalID("24567")

	assert.Equal(t, flows.TicketUUID("1ae96956-4b34-433e-8d1a-f05fe6923d6d"), ticket3.UUID())
	assert.Equal(t, mailgun, ticket3.Ticketer())
	assert.Equal(t, "New Ticket", ticket3.Subject())
	assert.Equal(t, "Where are my pants?", ticket3.Body())
	assert.Equal(t, "24567", ticket3.ExternalID())

	tickets.Add(ticket3)
	assert.Equal(t, 3, tickets.Count())
}
