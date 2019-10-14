package bothub_test

import (
	"net/http"
	"testing"

	"github.com/nyaruka/goflow/services/classification/bothub"
	"github.com/nyaruka/goflow/utils/httpx"
	"github.com/shopspring/decimal"

	"github.com/stretchr/testify/assert"
)

func TestPredict(t *testing.T) {
	defer httpx.SetRequestor(httpx.DefaultRequestor)

	httpx.SetRequestor(httpx.NewMockRequestor(map[string][]*httpx.MockResponse{
		"https://nlp.bothub.it/parse": []*httpx.MockResponse{
			httpx.NewMockResponse(200, `xx`), // non-JSON response
			httpx.NewMockResponse(200, `{}`), // invalid JSON response
			httpx.NewMockResponse(200, `{
				"intent": {
					"name": "greet",
					"confidence": 0.8341536248216568
				},
				"intent_ranking": [
					{
						"name": "greet",
						"confidence": 0.8341536248216568
					},
					{
						"name": "bye",
						"confidence": 0.16584637517834322
					}
				],
				"labels_list": [
					"hi"
				],
				"entities_list": [
					"hello"
				],
				"entities": {
					"hi": {
						"hello": [
							{
								"value": "hello",
								"entity": "hello",
								"confidence": 0.7979280788804916
							}
						]
					}
				},
				"text": "book a flight to Quito",
				"update_id": 4786,
				"language": "pt_br"
			}`),
		},
	}))

	client := bothub.NewClient(http.DefaultClient, "123e4567-e89b-12d3-a456-426655440000")

	response, trace, err := client.Parse("Hello")
	assert.EqualError(t, err, `invalid character 'x' looking for beginning of value`)
	assert.Equal(t, "POST /parse HTTP/1.1\r\nHost: nlp.bothub.it\r\nUser-Agent: Go-http-client/1.1\r\nContent-Length: 10\r\nAuthorization: Bearer 123e4567-e89b-12d3-a456-426655440000\r\nContent-Type: application/x-www-form-urlencoded\r\nAccept-Encoding: gzip\r\n\r\ntext=Hello", string(trace.RequestTrace))
	assert.Equal(t, "HTTP/1.0 200 OK\r\nContent-Length: 2\r\n\r\nxx", string(trace.ResponseTrace))
	assert.Nil(t, response)

	response, trace, err = client.Parse("Hello")
	assert.EqualError(t, err, `field 'intent_ranking' is required`)
	assert.NotNil(t, trace)
	assert.Nil(t, response)

	response, trace, err = client.Parse("book a flight to Quito")
	assert.NoError(t, err)
	assert.NotNil(t, trace)
	assert.Equal(t, bothub.IntentMatch{"greet", decimal.RequireFromString(`0.8341536248216568`)}, response.Intent)
	assert.Equal(t, []bothub.IntentMatch{
		bothub.IntentMatch{"greet", decimal.RequireFromString(`0.8341536248216568`)},
		bothub.IntentMatch{"bye", decimal.RequireFromString(`0.16584637517834322`)},
	}, response.IntentRanking)
	assert.Equal(t, []string{"hi"}, response.LabelsList)
	assert.Equal(t, []string{"hello"}, response.EntitiesList)
	assert.Equal(t, map[string]map[string][]bothub.EntityMatch{
		"hi": map[string][]bothub.EntityMatch{
			"hello": []bothub.EntityMatch{
				bothub.EntityMatch{Value: "hello", Entity: "hello", Confidence: decimal.RequireFromString(`0.7979280788804916`)},
			},
		},
	}, response.Entities)
	assert.Equal(t, "book a flight to Quito", response.Text)
	assert.Equal(t, 4786, response.UpdateID)
	assert.Equal(t, "pt_br", response.Language)
}
