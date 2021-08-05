package luis_test

import (
	"net/http"
	"testing"

	"github.com/nyaruka/gocommon/httpx"
	"github.com/nyaruka/goflow/services/classification/luis"
	"github.com/nyaruka/goflow/test"
	"github.com/shopspring/decimal"

	"github.com/stretchr/testify/assert"
)

var dec = decimal.RequireFromString

func TestPredict(t *testing.T) {
	httpx.SetRequestor(httpx.NewMockRequestor(map[string][]httpx.MockResponse{
		"https://luismm2.cognitiveservices.azure.com/luis/prediction/v3.0/apps/f96abf2f-3b53-4766-8ea6-09a655222a02/slots/production/predict?subscription-key=3246231&verbose=true&show-all-intents=true&log=true&query=book+flight+to+Quito": {
			httpx.NewMockResponse(200, nil, `xx`), // non-JSON response
			httpx.NewMockResponse(200, nil, `{}`), // invalid JSON response
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
	defer httpx.SetRequestor(httpx.DefaultRequestor)

	client := luis.NewClient(
		http.DefaultClient,
		nil,
		nil,
		"https://luismm2.cognitiveservices.azure.com/",
		"f96abf2f-3b53-4766-8ea6-09a655222a02",
		"3246231",
		"production",
	)

	response, trace, err := client.Predict("book flight to Quito")
	assert.EqualError(t, err, `invalid character 'x' looking for beginning of value`)
	test.AssertSnapshot(t, "predict_request", string(trace.RequestTrace))
	assert.Equal(t, "HTTP/1.0 200 OK\r\nContent-Length: 2\r\n\r\n", string(trace.ResponseTrace))
	assert.Equal(t, "xx", string(trace.ResponseBody))
	assert.Nil(t, response)

	response, trace, err = client.Predict("book flight to Quito")
	assert.EqualError(t, err, `field 'prediction' is required`)
	assert.NotNil(t, trace)
	assert.Nil(t, response)

	response, trace, err = client.Predict("book flight to Quito")
	assert.NoError(t, err)
	assert.NotNil(t, trace)
	assert.Equal(t, "book a flight to Quito", response.Query)
	assert.Equal(t, "Book Flight", response.Prediction.TopIntent)

	assert.Equal(t, map[string]luis.Intent{
		"Book Flight": {Score: dec(`0.9106805`)},
		"None":        {Score: dec(`0.08910245`)},
		"Book Hotel":  {Score: dec(`0.07790734`)},
	}, response.Prediction.Intents)

	assert.Equal(t, map[string][]string{"City": {"Quito"}}, response.Prediction.Entities.Values)
	assert.Equal(t, map[string][]luis.Entity{
		"City": {
			{Type: "location", Text: "Quito", Length: 5, Score: dec(`0.9984726`)},
		},
	}, response.Prediction.Entities.Instance)

	assert.Equal(t, &luis.Sentiment{Label: "positive", Score: dec(`0.48264188`)}, response.Prediction.Sentiment)
}
