package luis_test

import (
	"net/http"
	"testing"

	"github.com/nyaruka/goflow/services/classification/luis"
	"github.com/nyaruka/goflow/test"
	"github.com/nyaruka/goflow/utils/httpx"
	"github.com/shopspring/decimal"

	"github.com/stretchr/testify/assert"
)

func TestPredict(t *testing.T) {
	httpx.SetRequestor(httpx.NewMockRequestor(map[string][]httpx.MockResponse{
		"https://westus.api.cognitive.microsoft.com/luis/v2.0/apps/f96abf2f-3b53-4766-8ea6-09a655222a02?verbose=true&subscription-key=3246231&q=Hello": {
			httpx.NewMockResponse(200, nil, `xx`), // non-JSON response
			httpx.NewMockResponse(200, nil, `{}`), // invalid JSON response
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
	defer httpx.SetRequestor(httpx.DefaultRequestor)

	client := luis.NewClient(
		http.DefaultClient,
		nil,
		nil,
		"https://westus.api.cognitive.microsoft.com/luis/v2.0",
		"f96abf2f-3b53-4766-8ea6-09a655222a02",
		"3246231",
	)

	response, trace, err := client.Predict("Hello")
	assert.EqualError(t, err, `invalid character 'x' looking for beginning of value`)
	test.AssertSnapshot(t, "predict_request", string(trace.RequestTrace))
	assert.Equal(t, "HTTP/1.0 200 OK\r\nContent-Length: 2\r\n\r\n", string(trace.ResponseTrace))
	assert.Equal(t, "xx", string(trace.ResponseBody))
	assert.Nil(t, response)

	response, trace, err = client.Predict("Hello")
	assert.EqualError(t, err, `field 'intents' is required`)
	assert.NotNil(t, trace)
	assert.Nil(t, response)

	response, trace, err = client.Predict("Hello")
	assert.NoError(t, err)
	assert.NotNil(t, trace)
	assert.Equal(t, "book a flight to Quito", response.Query)
	assert.Equal(t, &luis.ExtractedIntent{"Book Flight", decimal.RequireFromString(`0.9106805`)}, response.TopScoringIntent)
	assert.Equal(t, []luis.ExtractedIntent{
		{"Book Flight", decimal.RequireFromString(`0.9106805`)},
		{"None", decimal.RequireFromString(`0.08910245`)},
		{"Book Hotel", decimal.RequireFromString(`0.07790734`)},
	}, response.Intents)
	assert.Equal(t, []luis.ExtractedEntity{
		{Entity: "quito", Type: "City", StartIndex: 17, EndIndex: 21, Score: decimal.RequireFromString(`0.9644149`)},
	}, response.Entities)
	assert.Equal(t, &luis.SentimentAnalysis{"positive", decimal.RequireFromString(`0.731448531`)}, response.SentimentAnalysis)
}
