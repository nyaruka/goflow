package mailgun_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/static/types"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/services/ticket/mailgun"
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
		"https://api.mailgun.net/v3/mr.nyaruka.com/messages": {
			httpx.MockConnectionError,
			httpx.NewMockResponse(200, nil, `{
				"id": "<20200426161758.1.590432020254B2BF@mr.nyaruka.com>",
				"message": "Queued. Thank you."
			}`),
		},
	}))

	ticketer := flows.NewTicketer(types.NewTicketer(assets.TicketerUUID(uuids.New()), "Support", "mailgun"))

	svc := mailgun.NewService(
		http.DefaultClient,
		nil,
		ticketer,
		"mr.nyaruka.com",
		"123456789",
		"bob@acme.com",
	)

	httpLogger := &flows.HTTPLogger{}

	_, err = svc.Open(session, "Need help", "Where are my cookies?", httpLogger.Log)
	assert.EqualError(t, err, "error calling mailgun API: unable to connect to server")

	httpLogger = &flows.HTTPLogger{}

	ticket, err := svc.Open(session, "Need help", "Where are my cookies?", httpLogger.Log)
	assert.NoError(t, err)
	assert.Equal(t, &flows.Ticket{
		UUID:       flows.TicketUUID("9688d21d-95aa-4bed-afc7-f31b35731a3d"),
		Ticketer:   ticketer.Reference(),
		Subject:    "Need help",
		Body:       "Where are my cookies?",
		ExternalID: "thread+9688d21d-95aa-4bed-afc7-f31b35731a3d@mr.nyaruka.com",
	}, ticket)

	assert.Equal(t, 1, len(httpLogger.Logs))
	assert.Equal(t, "https://api.mailgun.net/v3/mr.nyaruka.com/messages", httpLogger.Logs[0].URL)
}
