package luis_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/services/classification/luis"
	"github.com/nyaruka/goflow/test"
	"github.com/nyaruka/goflow/utils/dates"
	"github.com/nyaruka/goflow/utils/httpx"
	"github.com/nyaruka/goflow/utils/uuids"

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
		"https://westus.api.cognitive.microsoft.com/luis/v2.0/apps/f96abf2f-3b53-4766-8ea6-09a655222a02?verbose=true&subscription-key=3246231&q=book+flight+to+Quito": {
			httpx.NewMockResponse(200, nil, `{
				"query": "book a flight to Quito",
				"topScoringIntent": {
				  "intent": "Book Flight",
				  "score": 0.9106805
				},
				"intents": [
				  {
					"intent": "Book Flight",
					"score": 0.9106805
				  },
				  {
					"intent": "None",
					"score": 0.08910245
				  },
				  {
					"intent": "Book Hotel",
					"score": 0.07790734
				  }
				],
				"entities": [
				  {
					"entity": "quito",
					"type": "City",
					"startIndex": 17,
					"endIndex": 21,
					"score": 0.9644149
				  }
				],
				"sentimentAnalysis": {
				  "label": "positive",
				  "score": 0.731448531
				}
			}`),
		},
	}))

	svc := luis.NewService(
		http.DefaultClient,
		nil,
		nil,
		test.NewClassifier("Booking", "luis", []string{"book_flight", "book_hotel"}),
		"https://westus.api.cognitive.microsoft.com/luis/v2.0",
		"f96abf2f-3b53-4766-8ea6-09a655222a02",
		"3246231",
	)

	httpLogger := &flows.HTTPLogger{}

	classification, err := svc.Classify(session, "book flight to Quito", httpLogger.Log)
	assert.NoError(t, err)
	assert.Equal(t, []flows.ExtractedIntent{
		{Name: "Book Flight", Confidence: decimal.RequireFromString(`0.9106805`)},
		{Name: "None", Confidence: decimal.RequireFromString(`0.08910245`)},
		{Name: "Book Hotel", Confidence: decimal.RequireFromString(`0.07790734`)},
	}, classification.Intents)
	assert.Equal(t, map[string][]flows.ExtractedEntity{
		"City": {
			flows.ExtractedEntity{Value: "quito", Confidence: decimal.RequireFromString(`0.9644149`)},
		},
		"sentiment": {
			flows.ExtractedEntity{Value: "positive", Confidence: decimal.RequireFromString(`0.731448531`)},
		},
	}, classification.Entities)

	assert.Equal(t, 1, len(httpLogger.Logs))
	assert.Equal(t, "https://westus.api.cognitive.microsoft.com/luis/v2.0/apps/f96abf2f-3b53-4766-8ea6-09a655222a02?verbose=true&subscription-key=****************&q=book+flight+to+Quito", httpLogger.Logs[0].URL)
	assert.Equal(t, "GET /luis/v2.0/apps/f96abf2f-3b53-4766-8ea6-09a655222a02?verbose=true&subscription-key=****************&q=book+flight+to+Quito HTTP/1.1\r\nHost: westus.api.cognitive.microsoft.com\r\nUser-Agent: Go-http-client/1.1\r\nAccept-Encoding: gzip\r\n\r\n", httpLogger.Logs[0].Request)
}
