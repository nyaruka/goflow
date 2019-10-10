package luis_test

import (
	"encoding/json"
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
	uuids.SetGenerator(uuids.NewSeededGenerator(12345))
	dates.SetNowSource(dates.NewSequentialNowSource(time.Date(2019, 10, 7, 15, 21, 30, 123456789, time.UTC)))
	httpx.SetRequestor(httpx.NewMockRequestor(map[string][]*httpx.MockResponse{
		"https://westus.api.cognitive.microsoft.com/luis/v2.0/apps/f96abf2f-3b53-4766-8ea6-09a655222a02?verbose=true&subscription-key=3246231&q=book+flight+to+Quito": []*httpx.MockResponse{
			httpx.NewMockResponse(200, `{
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

	defer uuids.SetGenerator(uuids.DefaultGenerator)
	defer dates.SetNowSource(dates.DefaultNowSource)
	defer httpx.SetRequestor(httpx.DefaultRequestor)

	session, _, err := test.CreateTestSession("", nil, envs.RedactionPolicyNone)
	require.NoError(t, err)

	svc := luis.NewService(
		test.NewClassifier("Booking", "luis", []string{"book_flight", "book_hotel"}),
		"https://westus.api.cognitive.microsoft.com/luis/v2.0",
		"f96abf2f-3b53-4766-8ea6-09a655222a02",
		"3246231",
	)

	eventLog := test.NewEventLog()

	classification, err := svc.Classify(session, "book flight to Quito", eventLog.Log)
	assert.NoError(t, err)
	assert.Equal(t, []flows.ExtractedIntent{
		flows.ExtractedIntent{Name: "Book Flight", Confidence: decimal.RequireFromString(`0.9106805`)},
		flows.ExtractedIntent{Name: "None", Confidence: decimal.RequireFromString(`0.08910245`)},
		flows.ExtractedIntent{Name: "Book Hotel", Confidence: decimal.RequireFromString(`0.07790734`)},
	}, classification.Intents)
	assert.Equal(t, map[string][]flows.ExtractedEntity{
		"City": []flows.ExtractedEntity{
			flows.ExtractedEntity{Value: "quito", Confidence: decimal.RequireFromString(`0.9644149`)},
		},
		"sentiment": []flows.ExtractedEntity{
			flows.ExtractedEntity{Value: "positive", Confidence: decimal.RequireFromString(`0.731448531`)},
		},
	}, classification.Entities)

	eventsJSON, _ := json.Marshal(eventLog.Events)
	test.AssertEqualJSON(t, []byte(`[
		{
			"classifier": {
				"name": "Booking",
				"uuid": "20cc4181-48cf-4344-9751-99419796decd"
			},
			"created_on": "2019-10-07T15:22:29.123456789Z",
			"elapsed_ms": 1000,
			"request": "GET /luis/v2.0/apps/f96abf2f-3b53-4766-8ea6-09a655222a02?verbose=true&subscription-key=3246231&q=book+flight+to+Quito HTTP/1.1\r\nHost: westus.api.cognitive.microsoft.com\r\nUser-Agent: Go-http-client/1.1\r\nAccept-Encoding: gzip\r\n\r\n",
			"response": "HTTP/1.0 200 OK\r\nContent-Length: 605\r\n\r\n{\n\t\t\t\t\"query\": \"book a flight to Quito\",\n\t\t\t\t\"topScoringIntent\": {\n\t\t\t\t  \"intent\": \"Book Flight\",\n\t\t\t\t  \"score\": 0.9106805\n\t\t\t\t},\n\t\t\t\t\"intents\": [\n\t\t\t\t  {\n\t\t\t\t\t\"intent\": \"Book Flight\",\n\t\t\t\t\t\"score\": 0.9106805\n\t\t\t\t  },\n\t\t\t\t  {\n\t\t\t\t\t\"intent\": \"None\",\n\t\t\t\t\t\"score\": 0.08910245\n\t\t\t\t  },\n\t\t\t\t  {\n\t\t\t\t\t\"intent\": \"Book Hotel\",\n\t\t\t\t\t\"score\": 0.07790734\n\t\t\t\t  }\n\t\t\t\t],\n\t\t\t\t\"entities\": [\n\t\t\t\t  {\n\t\t\t\t\t\"entity\": \"quito\",\n\t\t\t\t\t\"type\": \"City\",\n\t\t\t\t\t\"startIndex\": 17,\n\t\t\t\t\t\"endIndex\": 21,\n\t\t\t\t\t\"score\": 0.9644149\n\t\t\t\t  }\n\t\t\t\t],\n\t\t\t\t\"sentimentAnalysis\": {\n\t\t\t\t  \"label\": \"positive\",\n\t\t\t\t  \"score\": 0.731448531\n\t\t\t\t}\n\t\t\t}",
			"status": "success",
			"type": "classifier_called",
			"url": "https://westus.api.cognitive.microsoft.com/luis/v2.0/apps/f96abf2f-3b53-4766-8ea6-09a655222a02?verbose=true&subscription-key=3246231&q=book+flight+to+Quito"
		}
	]`), eventsJSON, "events JSON mismatch")
}
