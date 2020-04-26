package mailgun_test

import (
	"net/http"
	"testing"

	"github.com/nyaruka/goflow/services/ticket/mailgun"
	"github.com/nyaruka/goflow/utils/httpx"
	"github.com/nyaruka/goflow/utils/uuids"

	"github.com/stretchr/testify/assert"
)

func TestSendMessage(t *testing.T) {
	defer httpx.SetRequestor(httpx.DefaultRequestor)
	defer uuids.SetGenerator(uuids.DefaultGenerator)

	httpx.SetRequestor(httpx.NewMockRequestor(map[string][]httpx.MockResponse{
		"https://api.mailgun.net/v3/mr.nyaruka.com/messages": {
			httpx.NewMockResponse(200, nil, `{"id":"<20200426161758.1.590432020254B2BF@mr.nyaruka.com>","message":"Queued. Thank you."}`),
		},
	}))
	uuids.SetGenerator(uuids.NewSeededGenerator(12345))

	client := mailgun.NewClient(http.DefaultClient, nil, "mr.nyaruka.com", "123456789")

	id, trace, err := client.SendMessage("Bob <thread+12446@mr.nyaruka.com>", "support@acme.com", "Need help", "Where are my cookies?")
	assert.NoError(t, err)
	assert.Equal(t, "<20200426161758.1.590432020254B2BF@mr.nyaruka.com>", id)
	assert.Equal(t, "POST /v3/mr.nyaruka.com/messages HTTP/1.1\r\nHost: api.mailgun.net\r\nUser-Agent: Go-http-client/1.1\r\nContent-Length: 478\r\nAuthorization: Basic YXBpOjEyMzQ1Njc4OQ==\r\nContent-Type: multipart/form-data; boundary=1ae96956-4b34-433e-8d1a-f05fe6923d6d\r\nAccept-Encoding: gzip\r\n\r\n--1ae96956-4b34-433e-8d1a-f05fe6923d6d\r\nContent-Disposition: form-data; name=\"from\"\r\n\r\nBob <thread+12446@mr.nyaruka.com>\r\n--1ae96956-4b34-433e-8d1a-f05fe6923d6d\r\nContent-Disposition: form-data; name=\"to\"\r\n\r\nsupport@acme.com\r\n--1ae96956-4b34-433e-8d1a-f05fe6923d6d\r\nContent-Disposition: form-data; name=\"subject\"\r\n\r\nNeed help\r\n--1ae96956-4b34-433e-8d1a-f05fe6923d6d\r\nContent-Disposition: form-data; name=\"text\"\r\n\r\nWhere are my cookies?\r\n--1ae96956-4b34-433e-8d1a-f05fe6923d6d--\r\n", string(trace.RequestTrace))
	assert.Equal(t, "HTTP/1.0 200 OK\r\nContent-Length: 90\r\n\r\n", string(trace.ResponseTrace))
	assert.Equal(t, `{"id":"<20200426161758.1.590432020254B2BF@mr.nyaruka.com>","message":"Queued. Thank you."}`, string(trace.ResponseBody))
}
