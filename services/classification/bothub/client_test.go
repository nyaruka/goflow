package bothub_test

import (
	"net/http"
	"testing"

	"github.com/nyaruka/goflow/services/classification/bothub"
	"github.com/nyaruka/goflow/test"
	"github.com/nyaruka/goflow/utils/httpx"
	"github.com/shopspring/decimal"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	defer httpx.SetRequestor(httpx.DefaultRequestor)

	httpx.SetRequestor(httpx.NewMockRequestor(map[string][]httpx.MockResponse{
		"https://nlp.bothub.it/parse": {
			httpx.NewMockResponse(200, nil, `xx`), // non-JSON response
			httpx.NewMockResponse(200, nil, `{}`), // invalid JSON response
			httpx.NewMockResponse(200, nil, `{
				"intent": {
					"name": "book_flight",
					"confidence": 0.8341536248216568
				},
				"intent_ranking": [
					{
						"name": "book_flight",
						"confidence": 0.8341536248216568
					},
					{
						"name": "book_hotel",
						"confidence": 0.16584637517834322
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
							"confidence": 0.7979280788804916
						}
					]
				},
				"text": "book a flight to Quito",
				"update_id": 4786,
				"language": "pt_br"
			}`),
		},
	}))

	client := bothub.NewClient(http.DefaultClient, nil, "123e4567-e89b-12d3-a456-426655440000")

	response, trace, err := client.Parse("Hello")
	assert.EqualError(t, err, `invalid character 'x' looking for beginning of value`)
	test.AssertSnapshot(t, "parse_request", string(trace.RequestTrace))
	assert.Equal(t, "HTTP/1.0 200 OK\r\nContent-Length: 2\r\n\r\n", string(trace.ResponseTrace))
	assert.Equal(t, "xx", string(trace.ResponseBody))
	assert.Nil(t, response)

	response, trace, err = client.Parse("Hello")
	assert.EqualError(t, err, `field 'intent_ranking' is required`)
	assert.NotNil(t, trace)
	assert.Nil(t, response)

	response, trace, err = client.Parse("book a flight to Quito")
	assert.NoError(t, err)
	assert.NotNil(t, trace)
	assert.Equal(t, bothub.IntentMatch{"book_flight", decimal.RequireFromString(`0.8341536248216568`)}, response.Intent)
	assert.Equal(t, []bothub.IntentMatch{
		{"book_flight", decimal.RequireFromString(`0.8341536248216568`)},
		{"book_hotel", decimal.RequireFromString(`0.16584637517834322`)},
	}, response.IntentRanking)
	assert.Equal(t, []string{"destination"}, response.LabelsList)
	assert.Equal(t, []string{"quito"}, response.EntitiesList)
	assert.Equal(t, map[string][]bothub.EntityMatch{
		"destination": {
			bothub.EntityMatch{Value: "quito", Entity: "quito", Confidence: decimal.RequireFromString(`0.7979280788804916`)},
		},
	}, response.Entities)
	assert.Equal(t, "book a flight to Quito", response.Text)
	assert.Equal(t, 4786, response.UpdateID)
	assert.Equal(t, "pt_br", response.Language)
}
