package luis_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/nyaruka/gocommon/dates"
	"github.com/nyaruka/gocommon/httpx"
	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/services/classification/luis"
	"github.com/nyaruka/goflow/test"

	"github.com/stretchr/testify/assert"
)

func TestService(t *testing.T) {
	session, _ := test.NewSessionBuilder().MustBuild()

	defer uuids.SetGenerator(uuids.DefaultGenerator)
	defer dates.SetNowSource(dates.DefaultNowSource)
	defer httpx.SetRequestor(httpx.DefaultRequestor)

	uuids.SetGenerator(uuids.NewSeededGenerator(12345))
	dates.SetNowSource(dates.NewSequentialNowSource(time.Date(2019, 10, 7, 15, 21, 30, 123456789, time.UTC)))
	httpx.SetRequestor(httpx.NewMockRequestor(map[string][]httpx.MockResponse{
		"https://luismm2.cognitiveservices.azure.com/luis/prediction/v3.0/apps/f96abf2f-3b53-4766-8ea6-09a655222a02/slots/production/predict?subscription-key=3246231&verbose=true&show-all-intents=true&log=true&query=book+flight+to+Quito": {
			httpx.NewMockResponse(200, nil, `{
				"query": "book a flight to Quito",
				"prediction": {
					"topIntent": "Book Flight",
					"intents": {
						"Book Flight": {
							"score": 0.9106805
						},
						"None": {
							"score": 0.08910245
						},
						"Book Hotel": {
							"score": 0.07790734
						}
					},
					"entities": {
						"City":[
							"Quito"
						],
						"$instance": {
							"City": [
								{
									"type": "location",
									"text": "Quito",
									"startIndex": 12,
									"length": 5,
									"score": 0.9984726,
									"modelTypeId": 1,
									"modelType": "Entity Extractor",
									"recognitionSources": ["model"]
								}
							]
						}
					},
					"sentiment":{
						"label": "positive",
						"score": 0.48264188
					}
				}
			}`),
		},
	}))

	svc := luis.NewService(
		http.DefaultClient,
		nil,
		nil,
		test.NewClassifier("Booking", "luis", []string{"book_flight", "book_hotel"}),
		"https://luismm2.cognitiveservices.azure.com/",
		"f96abf2f-3b53-4766-8ea6-09a655222a02",
		"3246231",
		"production",
	)

	httpLogger := &flows.HTTPLogger{}

	classification, err := svc.Classify(session, "book flight to Quito", httpLogger.Log)
	assert.NoError(t, err)
	assert.Equal(t, []flows.ExtractedIntent{
		{Name: "Book Flight", Confidence: dec(`0.9106805`)},
		{Name: "None", Confidence: dec(`0.08910245`)},
		{Name: "Book Hotel", Confidence: dec(`0.07790734`)},
	}, classification.Intents)
	assert.Equal(t, map[string][]flows.ExtractedEntity{
		"City": {
			flows.ExtractedEntity{Value: "Quito", Confidence: dec(`0.9984726`)},
		},
		"sentiment": {
			flows.ExtractedEntity{Value: "positive", Confidence: dec(`0.48264188`)},
		},
	}, classification.Entities)

	assert.Equal(t, 1, len(httpLogger.Logs))
	assert.Equal(t, "https://luismm2.cognitiveservices.azure.com/luis/prediction/v3.0/apps/f96abf2f-3b53-4766-8ea6-09a655222a02/slots/production/predict?subscription-key=****************&verbose=true&show-all-intents=true&log=true&query=book+flight+to+Quito", httpLogger.Logs[0].URL)
	assert.Equal(t, "GET /luis/prediction/v3.0/apps/f96abf2f-3b53-4766-8ea6-09a655222a02/slots/production/predict?subscription-key=****************&verbose=true&show-all-intents=true&log=true&query=book+flight+to+Quito HTTP/1.1\r\nHost: luismm2.cognitiveservices.azure.com\r\nUser-Agent: Go-http-client/1.1\r\nAccept-Encoding: gzip\r\n\r\n", httpLogger.Logs[0].Request)
}
