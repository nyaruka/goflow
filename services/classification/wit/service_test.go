package wit_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/nyaruka/gocommon/dates"
	"github.com/nyaruka/gocommon/httpx"
	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/services/classification/wit"
	"github.com/nyaruka/goflow/test"

	"github.com/shopspring/decimal"
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
		"https://api.wit.ai/message?v=20200513&q=book+flight+to+Quito": {
			httpx.NewMockResponse(200, nil, `{
				"text": "I want to book a flight to Quito",
				"intents": [
				  {
					"id": "754569408690533",
					"name": "book_flight",
					"confidence": 0.9024
				  }
				],
				"entities": {
				  "Destination:Location": [
					{
					  "id": "285857329187179",
					  "name": "Destination",
					  "role": "Location",
					  "start": 27,
					  "end": 33,
					  "body": "Quito",
					  "confidence": 0.9648,
					  "entities": [],
					  "value": "Quito",
					  "type": "value"
					}
				  ]
				},
				"traits": {
				  "wit$sentiment": [
					{
					  "id": "5ac2b50a-44e4-466e-9d49-bad6bd40092c",
					  "value": "neutral",
					  "confidence": 0.5816
					}
				  ]
				}
			}`),
		},
	}))

	svc := wit.NewService(
		http.DefaultClient,
		nil,
		test.NewClassifier("Booking", "wit", []string{"book_flight", "book_hotel"}),
		"23532624376",
	)

	httpLogger := &flows.HTTPLogger{}

	classification, err := svc.Classify(session, "book flight to Quito", httpLogger.Log)
	assert.NoError(t, err)
	assert.Equal(t, []flows.ExtractedIntent{
		{Name: "book_flight", Confidence: decimal.RequireFromString(`0.9024`)},
	}, classification.Intents)
	assert.Equal(t, map[string][]flows.ExtractedEntity{
		"Destination": {{Value: "Quito", Confidence: decimal.RequireFromString(`0.9648`)}},
	}, classification.Entities)

	assert.Equal(t, 1, len(httpLogger.Logs))
	assert.Equal(t, "https://api.wit.ai/message?v=20200513&q=book+flight+to+Quito", httpLogger.Logs[0].URL)
	assert.Equal(t, "GET /message?v=20200513&q=book+flight+to+Quito HTTP/1.1\r\nHost: api.wit.ai\r\nUser-Agent: Go-http-client/1.1\r\nAuthorization: Bearer ****************\r\nAccept-Encoding: gzip\r\n\r\n", httpLogger.Logs[0].Request)
}
