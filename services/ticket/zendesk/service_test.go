package zendesk_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/static/types"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/services/ticket/zendesk"
	"github.com/nyaruka/goflow/test"
	"github.com/nyaruka/goflow/utils/dates"
	"github.com/nyaruka/goflow/utils/httpx"
	"github.com/nyaruka/goflow/utils/uuids"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService(t *testing.T) {
	session, _, err := test.CreateTestSession("", envs.RedactionPolicyNone)
	require.NoError(t, err)

	defer uuids.SetGenerator(uuids.DefaultGenerator)
	defer dates.SetNowSource(dates.DefaultNowSource)
	defer httpx.SetRequestor(httpx.DefaultRequestor)

	uuids.SetGenerator(uuids.NewSeededGenerator(12345))
	dates.SetNowSource(dates.NewSequentialNowSource(time.Date(2019, 10, 7, 15, 21, 30, 123456789, time.UTC)))
	httpx.SetRequestor(httpx.NewMockRequestor(map[string][]httpx.MockResponse{
		"https://nyaruka.zendesk.com/api/v2/tickets.json": {
			httpx.MockConnectionError,
			httpx.NewMockResponse(201, nil, `{
				"ticket":{
					"id": 12345,
					"url": "https://nyaruka.zendesk.com/api/v2/tickets/12345.json",
					"external_id": "a78c5d9d-283a-4be9-ad6d-690e4307c961",
					"created_at": "2009-07-20T22:55:29Z",
					"subject": "Need help"
				}
			}`),
		},
	}))

	ticketer := flows.NewTicketer(types.NewTicketer(assets.TicketerUUID(uuids.New()), "Support", "zendesk"))

	svc := zendesk.NewService(
		http.DefaultClient,
		nil,
		ticketer,
		"nyaruka",
		"zen@nyaruka.com",
		"123456789",
	)

	httpLogger := &flows.HTTPLogger{}

	_, err = svc.Open(session, "Need help", "Where are my cookies?", httpLogger.Log)
	assert.EqualError(t, err, "error calling zendesk API: unable to connect to server")

	httpLogger = &flows.HTTPLogger{}

	ticket, err := svc.Open(session, "Need help", "Where are my cookies?", httpLogger.Log)

	assert.NoError(t, err)
	assert.Equal(t, &flows.Ticket{
		UUID:       flows.TicketUUID("59d74b86-3e2f-4a93-aece-b05d2fdcde0c"),
		Ticketer:   ticketer.Reference(),
		Subject:    "Need help",
		Body:       "Where are my cookies?",
		ExternalID: "12345",
	}, ticket)

	assert.Equal(t, 1, len(httpLogger.Logs))
	assert.Equal(t, "https://nyaruka.zendesk.com/api/v2/tickets.json", httpLogger.Logs[0].URL)
}
