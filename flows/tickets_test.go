package flows_test

import (
	"encoding/json"
	"testing"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/static"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/engine"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTickets(t *testing.T) {
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

	missingRefs := make([]assets.Reference, 0)
	missing := func(ref assets.Reference, err error) {
		missingRefs = append(missingRefs, ref)
	}

	ticketsJSON := `[
		{
			"uuid": "349c851f-3f8e-4353-8bf2-8e90b6d73530", 
			"ticketer": {"uuid": "0a0b5ce4-35c9-47b7-b124-40258f0a5b53", "name": "Deleted"},
			"subject": "Very Old ticket",
			"body": "Ticketer gone!"
		},
		{
			"uuid": "5a4af021-d2c2-47fc-9abc-abbb8635d8c0", 
			"ticketer": {"uuid": "d605bb96-258d-4097-ad0a-080937db2212", "name": "Support Tickets"},
			"subject": "Old ticket",
			"body": "Where are my shoes?"
		}
	]`
	var ticketRefs []*flows.TicketReference
	err = json.Unmarshal([]byte(ticketsJSON), &ticketRefs)
	require.NoError(t, err)

	tickets := flows.NewTicketList(sa, ticketRefs, missing)
	assert.Equal(t, 1, tickets.Count())
	assert.Equal(t, "Old ticket", tickets.All()[0].Subject)

	// check that ticket with missing ticketer is logged as a missing dependency
	assert.Equal(t, 1, len(missingRefs))
	assert.Equal(t, "0a0b5ce4-35c9-47b7-b124-40258f0a5b53", missingRefs[0].Identity())

	ticket := flows.NewTicket(mailgun, "New ticket", "Where are my pants?", "893732")
	tickets.Add(ticket)

	assert.Equal(t, 2, tickets.Count())
}
