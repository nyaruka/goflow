package zendesk_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/nyaruka/goflow/services/ticket/zendesk"
	"github.com/nyaruka/goflow/utils/httpx"
	"github.com/nyaruka/goflow/utils/uuids"

	"github.com/stretchr/testify/assert"
)

func TestCreateTicket(t *testing.T) {
	defer httpx.SetRequestor(httpx.DefaultRequestor)
	defer uuids.SetGenerator(uuids.DefaultGenerator)

	httpx.SetRequestor(httpx.NewMockRequestor(map[string][]httpx.MockResponse{
		"https://nyaruka.zendesk.com/api/v2/tickets.json": {
			httpx.MockConnectionError,
			httpx.NewMockResponse(400, nil, `{"description": "Something went wrong", "error": "Unknown"}`), // non-200 response
			httpx.NewMockResponse(200, nil, `xx`), // non-JSON response
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
	uuids.SetGenerator(uuids.NewSeededGenerator(12345))

	client := zendesk.NewClient(http.DefaultClient, nil, "nyaruka", "zen@nyaruka.com", "123456789")

	_, _, err := client.CreateTicket("Need help", "Where are my cookies?")
	assert.EqualError(t, err, "unable to connect to server")

	_, _, err = client.CreateTicket("Need help", "Where are my cookies?")
	assert.EqualError(t, err, "Something went wrong")

	_, _, err = client.CreateTicket("Need help", "Where are my cookies?")
	assert.EqualError(t, err, "invalid character 'x' looking for beginning of value")

	zenTicket, trace, err := client.CreateTicket("Need help", "Where are my cookies?")
	assert.NoError(t, err)
	assert.Equal(t, 12345, zenTicket.ID)
	assert.Equal(t, time.Date(2009, 7, 20, 22, 55, 29, 0, time.UTC), zenTicket.CreatedAt)
	assert.Equal(t, "POST /api/v2/tickets.json HTTP/1.1\r\nHost: nyaruka.zendesk.com\r\nUser-Agent: Go-http-client/1.1\r\nContent-Length: 94\r\nAuthorization: Basic emVuQG55YXJ1a2EuY29tOjEyMzQ1Njc4OQ==\r\nAccept-Encoding: gzip\r\n\r\n{\"ticket\":{\"subject\":\"Need help\",\"comment\":{\"body\":\"Where are my cookies?\"},\"external_id\":\"\"}}", string(trace.RequestTrace))
	assert.Equal(t, "HTTP/1.0 201 Created\r\nContent-Length: 245\r\n\r\n", string(trace.ResponseTrace))
}
