package flows_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/nyaruka/goflow/flows"

	"github.com/nyaruka/goflow/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type call struct {
	method string
	url    string
	body   string
}

func (c *call) String() string { return c.method + " " + c.url }

type webhook struct {
	request     string
	response    string
	bodyIgnored bool
	json        []byte
}

func TestWebhookParsing(t *testing.T) {
	server := test.NewTestHTTPServer(49994)
	defer server.Close()

	session, _, err := test.CreateTestSession(server.URL, nil)
	require.NoError(t, err)

	testCases := []struct {
		call    call
		webhook webhook
		isError bool
	}{
		{
			// successful GET
			call: call{"GET", "http://127.0.0.1:49994/?cmd=success", ""},
			webhook: webhook{
				request:  "GET /?cmd=success HTTP/1.1\r\nHost: 127.0.0.1:49994\r\nUser-Agent: goflow-testing\r\nAccept-Encoding: gzip\r\n\r\n",
				response: "HTTP/1.1 200 OK\r\nContent-Length: 16\r\nContent-Type: text/plain; charset=utf-8\r\nDate: Wed, 11 Apr 2018 18:24:30 GMT\r\n\r\n{ \"ok\": \"true\" }",
				json:     []byte(`{ "ok": "true" }`),
			},
		}, {
			// successful GET, text/javascrpt
			call: call{"GET", "http://127.0.0.1:49994/?cmd=textjs", ""},
			webhook: webhook{
				request:  "GET /?cmd=textjs HTTP/1.1\r\nHost: 127.0.0.1:49994\r\nUser-Agent: goflow-testing\r\nAccept-Encoding: gzip\r\n\r\n",
				response: "HTTP/1.1 200 OK\r\nContent-Length: 16\r\nContent-Type: text/javascript; charset=iso-8859-1\r\nDate: Wed, 11 Apr 2018 18:24:30 GMT\r\n\r\n{ \"ok\": \"true\" }",
				json:     []byte(`{ "ok": "true" }`),
			},
		}, {
			// successful POST without body
			call: call{"POST", "http://127.0.0.1:49994/?cmd=success", ""},
			webhook: webhook{
				request:  "POST /?cmd=success HTTP/1.1\r\nHost: 127.0.0.1:49994\r\nUser-Agent: goflow-testing\r\nContent-Length: 0\r\nAccept-Encoding: gzip\r\n\r\n",
				response: "HTTP/1.1 200 OK\r\nContent-Length: 16\r\nContent-Type: text/plain; charset=utf-8\r\nDate: Wed, 11 Apr 2018 18:24:30 GMT\r\n\r\n{ \"ok\": \"true\" }",
				json:     []byte(`{ "ok": "true" }`),
			},
		}, {
			// successful POST with body
			call: call{"POST", "http://127.0.0.1:49994/?cmd=success", `{"contact": "Bob"}`},
			webhook: webhook{
				request:  "POST /?cmd=success HTTP/1.1\r\nHost: 127.0.0.1:49994\r\nUser-Agent: goflow-testing\r\nContent-Length: 18\r\nAccept-Encoding: gzip\r\n\r\n{\"contact\": \"Bob\"}",
				response: "HTTP/1.1 200 OK\r\nContent-Length: 16\r\nContent-Type: text/plain; charset=utf-8\r\nDate: Wed, 11 Apr 2018 18:24:30 GMT\r\n\r\n{ \"ok\": \"true\" }",
				json:     []byte(`{ "ok": "true" }`),
			},
		}, {
			// POST returning 503
			call: call{"POST", "http://127.0.0.1:49994/?cmd=unavailable", ""},
			webhook: webhook{
				request:  "POST /?cmd=unavailable HTTP/1.1\r\nHost: 127.0.0.1:49994\r\nUser-Agent: goflow-testing\r\nContent-Length: 0\r\nAccept-Encoding: gzip\r\n\r\n",
				response: "HTTP/1.1 503 Service Unavailable\r\nContent-Length: 37\r\nContent-Type: text/plain; charset=utf-8\r\nDate: Wed, 11 Apr 2018 18:24:30 GMT\r\n\r\n{ \"errors\": [\"service unavailable\"] }",
				json:     []byte(`{ "errors": ["service unavailable"] }`),
			},
		}, {
			// GET returning non-text content type
			call: call{"GET", "http://127.0.0.1:49994/?cmd=binary", ""},
			webhook: webhook{
				request:     "GET /?cmd=binary HTTP/1.1\r\nHost: 127.0.0.1:49994\r\nUser-Agent: goflow-testing\r\nAccept-Encoding: gzip\r\n\r\n",
				response:    "HTTP/1.1 200 OK\r\nContent-Length: 10\r\nContent-Type: application/octet-stream\r\nDate: Wed, 11 Apr 2018 18:24:30 GMT\r\n\r\n",
				bodyIgnored: true,
				json:        nil,
			},
		}, {
			// GET returning binary body larger than allowed (we ignore binary body so no biggie)
			call: call{"GET", "http://127.0.0.1:49994/?cmd=binary&size=11000", ""},
			webhook: webhook{
				request:     "GET /?cmd=binary&size=11000 HTTP/1.1\r\nHost: 127.0.0.1:49994\r\nUser-Agent: goflow-testing\r\nAccept-Encoding: gzip\r\n\r\n",
				response:    "HTTP/1.1 200 OK\r\nContent-Length: 11000\r\nContent-Type: application/octet-stream\r\nDate: Wed, 11 Apr 2018 18:24:30 GMT\r\n\r\n",
				bodyIgnored: true,
				json:        nil,
			},
		}, {
			// GET returning text body larger than allowed
			call:    call{"GET", "http://127.0.0.1:49994/?cmd=binary&size=11000&type=text%2Fplain", ""},
			isError: true,
		}, {
			// GET returning text body but an empty content-type header
			call: call{"GET", "http://127.0.0.1:49994/?cmd=typeless&content=kthxbai", ""},
			webhook: webhook{
				request:  "GET /?cmd=typeless&content=kthxbai HTTP/1.1\r\nHost: 127.0.0.1:49994\r\nUser-Agent: goflow-testing\r\nAccept-Encoding: gzip\r\n\r\n",
				response: "HTTP/1.1 200 OK\r\nContent-Length: 7\r\nContent-Type: \r\nDate: Wed, 11 Apr 2018 18:24:30 GMT\r\n\r\nkthxbai",
				json:     []byte(`"kthxbai"`),
			},
		}, {
			// GET returning JSON body but an empty content-type header
			call: call{"GET", "http://127.0.0.1:49994/?cmd=typeless&content=%7B%22msg%22%3A%20%22I%27m%20JSON%22%7D", ""},
			webhook: webhook{
				request:  "GET /?cmd=typeless&content=%7B%22msg%22%3A%20%22I%27m%20JSON%22%7D HTTP/1.1\r\nHost: 127.0.0.1:49994\r\nUser-Agent: goflow-testing\r\nAccept-Encoding: gzip\r\n\r\n",
				response: "HTTP/1.1 200 OK\r\nContent-Length: 19\r\nContent-Type: \r\nDate: Wed, 11 Apr 2018 18:24:30 GMT\r\n\r\n{\"msg\": \"I'm JSON\"}",
				json:     []byte(`{"msg": "I'm JSON"}`),
			},
		},
	}

	for _, tc := range testCases {
		request, err := http.NewRequest(tc.call.method, tc.call.url, strings.NewReader(tc.call.body))
		require.NoError(t, err)

		webhook, err := flows.MakeWebhookCall(session, request, "")
		if tc.isError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err, "unexpected error fetching %s", tc.call)

			assert.Equal(t, tc.call.url, webhook.URL(), "URL mismatch for call %s", tc.call)
			assert.Equal(t, "", webhook.Resthook(), "resthook mismatch for call %s", tc.call)
			assert.Equal(t, tc.call.method, webhook.Method(), "method mismatch for call %s", tc.call)
			assert.Equal(t, tc.webhook.request, webhook.Request(), "request trace mismatch for call %s", tc.call)
			assert.Equal(t, tc.webhook.response, webhook.Response(), "response mismatch for call %s", tc.call)
			assert.Equal(t, tc.webhook.bodyIgnored, webhook.BodyIgnored(), "body-ignored mismatch for call %s", tc.call)

			test.AssertEqualJSON(t, tc.webhook.json, flows.ExtractResponseBody(webhook.Response()), "body mismatch for call %s", tc.call)
		}
	}
}
