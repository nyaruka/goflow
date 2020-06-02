package bothub_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/services/classification/bothub"
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
		"https://nlp.bothub.it/parse": {
			httpx.NewMockResponse(200, nil, `{
				"intent": {
				  "name": "book_flight",
				  "confidence": 0.9224673593230207
				},
				"intent_ranking": [
				  {
					"name": "book_flight",
					"confidence": 0.9224673593230207
				  },
				  {
					"name": "book_hotel",
					"confidence": 0.07753264067697924
				  }
				],
				"labels_list": [
				  "destination"
				],
				"entities_list": [
				  "quito"
				],
				"entities": {
				  "destination": [
					{
					  "value": "quito",
					  "entity": "quito",
					  "confidence": 0.8824543190522534
					}
				  ]
				},
				"text": "book my flight to Quito",
				"update_id": 13158,
				"language": "en"
			  }`),
		},
	}))

	svc := bothub.NewService(
		http.DefaultClient,
		nil,
		test.NewClassifier("Booking", "bothub", []string{"book_flight", "book_hotel"}),
		"f96abf2f-3b53-4766-8ea6-09a655222a02",
	)

	httpLogger := &flows.HTTPLogger{}

	classification, err := svc.Classify(session, "book my flight to Quito", httpLogger.Log)
	assert.NoError(t, err)
	assert.Equal(t, []flows.ExtractedIntent{
		{Name: "book_flight", Confidence: decimal.RequireFromString(`0.9224673593230207`)},
		{Name: "book_hotel", Confidence: decimal.RequireFromString(`0.07753264067697924`)},
	}, classification.Intents)
	assert.Equal(t, map[string][]flows.ExtractedEntity{
		"destination": {
			flows.ExtractedEntity{Value: "quito", Confidence: decimal.RequireFromString(`0.8824543190522534`)},
		},
	}, classification.Entities)

	assert.Equal(t, 1, len(httpLogger.Logs))
	assert.Equal(t, "https://nlp.bothub.it/parse", httpLogger.Logs[0].URL)
	assert.Equal(t, "POST /parse HTTP/1.1\r\nHost: nlp.bothub.it\r\nUser-Agent: Go-http-client/1.1\r\nContent-Length: 28\r\nAuthorization: Bearer ****************\r\nContent-Type: application/x-www-form-urlencoded\r\nAccept-Encoding: gzip\r\n\r\ntext=book+my+flight+to+Quito", httpLogger.Logs[0].Request)
}
