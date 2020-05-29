package wit_test

import (
	"net/http"
	"testing"

	"github.com/nyaruka/goflow/services/classification/wit"
	"github.com/nyaruka/goflow/test"
	"github.com/nyaruka/goflow/utils/httpx"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestMessage(t *testing.T) {
	defer httpx.SetRequestor(httpx.DefaultRequestor)

	httpx.SetRequestor(httpx.NewMockRequestor(map[string][]httpx.MockResponse{
		"https://api.wit.ai/message?v=20170307&q=Hello": {
			httpx.NewMockResponse(200, nil, `xx`), // non-JSON response
			httpx.NewMockResponse(200, nil, `{}`), // invalid JSON response
			httpx.NewMockResponse(200, nil, `{"_text":"book flight","entities":{"intent":[{"confidence":0.84709152161066,"value":"book_flight"}]},"msg_id":"1M7fAcDWag76OmgDI"}`),
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
	assert.EqualError(t, err, `field 'entities' is required`)
	assert.NotNil(t, trace)
	assert.Nil(t, response)

	response, trace, err = client.Message("Hello")
	assert.NoError(t, err)
	assert.NotNil(t, trace)
	assert.Equal(t, "1M7fAcDWag76OmgDI", response.MsgID)
	assert.Equal(t, "book flight", response.Text)
	assert.Equal(t, map[string][]wit.EntityCandidate{"intent": {
		wit.EntityCandidate{Value: "book_flight", Confidence: decimal.RequireFromString(`0.84709152161066`)},
	}}, response.Entities)
}
