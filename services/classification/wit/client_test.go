package wit_test

import (
	"net/http"
	"testing"

	"github.com/nyaruka/goflow/services/classification/wit"
	"github.com/nyaruka/goflow/utils/httpx"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestMessage(t *testing.T) {
	httpx.SetRequestor(httpx.NewMockRequestor(map[string][]*http.Response{
		"https://api.wit.ai/message?v=20170307&q=Hello": []*http.Response{
			httpx.NewMockResponse(200, `xx`), // non-JSON response
			httpx.NewMockResponse(200, `{}`), // invalid JSON response
			httpx.NewMockResponse(200, `{"_text":"book flight","entities":{"intent":[{"confidence":0.84709152161066,"value":"book_flight"}]},"msg_id":"1M7fAcDWag76OmgDI"}`),
		},
	}))
	defer httpx.SetRequestor(httpx.DefaultRequestor)

	client := wit.NewClient(http.DefaultClient, "3246231")

	response, trace, err := client.Message("Hello")
	assert.EqualError(t, err, `invalid character 'x' looking for beginning of value`)
	assert.Equal(t, "GET /message?v=20170307&q=Hello HTTP/1.1\r\nHost: api.wit.ai\r\nUser-Agent: Go-http-client/1.1\r\nAuthorization: Bearer 3246231\r\nAccept-Encoding: gzip\r\n\r\n", string(trace.RequestTrace))
	assert.Equal(t, "HTTP/1.0 200 OK\r\nContent-Length: 2\r\n\r\nxx", string(trace.ResponseTrace))
	assert.Nil(t, response)

	response, trace, err = client.Message("Hello")
	assert.EqualError(t, err, `field 'entities' is required`)
	assert.NotNil(t, trace)
	assert.Nil(t, response)

	response, trace, err = client.Message("Hello")
	assert.NoError(t, err)
	assert.NotNil(t, trace)
	assert.Equal(t, "1M7fAcDWag76OmgDI", response.MsgID)
	assert.Equal(t, "book flight", response.Text)
	assert.Equal(t, map[string][]wit.EntityCandidate{"intent": []wit.EntityCandidate{
		wit.EntityCandidate{Value: "book_flight", Confidence: decimal.RequireFromString(`0.84709152161066`)},
	}}, response.Entities)
}
