package wit_test

import (
	"net/http"
	"testing"

	"github.com/nyaruka/gocommon/httpx"
	"github.com/nyaruka/goflow/services/classification/wit"
	"github.com/nyaruka/goflow/test"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestMessage(t *testing.T) {
	defer httpx.SetRequestor(httpx.DefaultRequestor)

	httpx.SetRequestor(httpx.NewMockRequestor(map[string][]httpx.MockResponse{
		"https://api.wit.ai/message?v=20200513&q=Hello": {
			httpx.NewMockResponse(200, nil, `xx`), // non-JSON response
			httpx.NewMockResponse(200, nil, `{}`), // invalid JSON response
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

	client := wit.NewClient(http.DefaultClient, nil, "3246231")

	response, trace, err := client.Message("Hello")
	assert.EqualError(t, err, `invalid character 'x' looking for beginning of value`)
	test.AssertSnapshot(t, "message_request", string(trace.RequestTrace))
	assert.Equal(t, "HTTP/1.0 200 OK\r\nContent-Length: 2\r\n\r\n", string(trace.ResponseTrace))
	assert.Equal(t, "xx", string(trace.ResponseBody))
	assert.Nil(t, response)

	response, trace, err = client.Message("Hello")
	assert.EqualError(t, err, `field 'intents' is required`)
	assert.NotNil(t, trace)
	assert.Nil(t, response)

	response, trace, err = client.Message("Hello")
	assert.NoError(t, err)
	assert.NotNil(t, trace)
	assert.Equal(t, "I want to book a flight to Quito", response.Text)
	assert.Equal(t, []wit.IntentMatch{{ID: "754569408690533", Name: "book_flight", Confidence: decimal.RequireFromString(`0.9024`)}}, response.Intents)
	assert.Equal(t, map[string][]wit.EntityMatch{
		"Destination:Location": {
			wit.EntityMatch{
				ID:         "285857329187179",
				Name:       "Destination",
				Role:       "Location",
				Value:      "Quito",
				Confidence: decimal.RequireFromString(`0.9648`),
			},
		},
	}, response.Entities)
}
