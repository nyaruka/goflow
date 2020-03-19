package webhooks_test

import (
	"net"
	"net/http"
	"strings"
	"testing"

	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/services/webhooks"
	"github.com/nyaruka/goflow/test"
	"github.com/nyaruka/goflow/utils/httpx"

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
	json        string
	bodyIgnored bool
}

func TestWebhookParsing(t *testing.T) {
	server := test.NewTestHTTPServer(49994)
	defer server.Close()

	session, _, err := test.CreateTestSession(server.URL, envs.RedactionPolicyNone)
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
				json:     `{ "ok": "true" }`,
			},
		}, {
			// successful GET with valid JSON response body
			call: call{"GET", "http://127.0.0.1:49994/?cmd=textjs", ""},
			webhook: webhook{
				request:  "GET /?cmd=textjs HTTP/1.1\r\nHost: 127.0.0.1:49994\r\nUser-Agent: goflow-testing\r\nAccept-Encoding: gzip\r\n\r\n",
				response: "HTTP/1.1 200 OK\r\nContent-Length: 16\r\nContent-Type: text/javascript; charset=iso-8859-1\r\nDate: Wed, 11 Apr 2018 18:24:30 GMT\r\n\r\n{ \"ok\": \"true\" }",
				json:     `{ "ok": "true" }`,
			},
		}, {
			// successful POST without request body and valid JSON response body
			call: call{"POST", "http://127.0.0.1:49994/?cmd=success", ""},
			webhook: webhook{
				request:  "POST /?cmd=success HTTP/1.1\r\nHost: 127.0.0.1:49994\r\nUser-Agent: goflow-testing\r\nContent-Length: 0\r\nAccept-Encoding: gzip\r\n\r\n",
				response: "HTTP/1.1 200 OK\r\nContent-Length: 16\r\nContent-Type: text/plain; charset=utf-8\r\nDate: Wed, 11 Apr 2018 18:24:30 GMT\r\n\r\n{ \"ok\": \"true\" }",
				json:     `{ "ok": "true" }`,
			},
		}, {
			// successful POST with request body and valid JSON response body
			call: call{"POST", "http://127.0.0.1:49994/?cmd=success", `{"contact": "Bob"}`},
			webhook: webhook{
				request:  "POST /?cmd=success HTTP/1.1\r\nHost: 127.0.0.1:49994\r\nUser-Agent: goflow-testing\r\nContent-Length: 18\r\nAccept-Encoding: gzip\r\n\r\n{\"contact\": \"Bob\"}",
				response: "HTTP/1.1 200 OK\r\nContent-Length: 16\r\nContent-Type: text/plain; charset=utf-8\r\nDate: Wed, 11 Apr 2018 18:24:30 GMT\r\n\r\n{ \"ok\": \"true\" }",
				json:     `{ "ok": "true" }`,
			},
		}, {
			// POST returning 503
			call: call{"POST", "http://127.0.0.1:49994/?cmd=unavailable", ""},
			webhook: webhook{
				request:  "POST /?cmd=unavailable HTTP/1.1\r\nHost: 127.0.0.1:49994\r\nUser-Agent: goflow-testing\r\nContent-Length: 0\r\nAccept-Encoding: gzip\r\n\r\n",
				response: "HTTP/1.1 503 Service Unavailable\r\nContent-Length: 37\r\nContent-Type: text/plain; charset=utf-8\r\nDate: Wed, 11 Apr 2018 18:24:30 GMT\r\n\r\n{ \"errors\": [\"service unavailable\"] }",
				json:     `{ "errors": ["service unavailable"] }`,
			},
		}, {
			// GET returning text body larger than allowed
			call:    call{"GET", "http://127.0.0.1:49994/?cmd=binary&size=11000", ""},
			isError: true,
		}, {
			// GET returning non-JSON body
			call: call{"GET", "http://127.0.0.1:49994/?cmd=typeless&content=kthxbai", ""},
			webhook: webhook{
				request:     "GET /?cmd=typeless&content=kthxbai HTTP/1.1\r\nHost: 127.0.0.1:49994\r\nUser-Agent: goflow-testing\r\nAccept-Encoding: gzip\r\n\r\n",
				response:    "HTTP/1.1 200 OK\r\nContent-Length: 7\r\nContent-Type: \r\nDate: Wed, 11 Apr 2018 18:24:30 GMT\r\n\r\nkthxbai",
				json:        "",
				bodyIgnored: true,
			},
		}, {
			// GET returning JSON body but an empty content-type header
			call: call{"GET", "http://127.0.0.1:49994/?cmd=typeless&content=%7B%22msg%22%3A%20%22I%27m%20JSON%22%7D", ""},
			webhook: webhook{
				request:  "GET /?cmd=typeless&content=%7B%22msg%22%3A%20%22I%27m%20JSON%22%7D HTTP/1.1\r\nHost: 127.0.0.1:49994\r\nUser-Agent: goflow-testing\r\nAccept-Encoding: gzip\r\n\r\n",
				response: "HTTP/1.1 200 OK\r\nContent-Length: 19\r\nContent-Type: \r\nDate: Wed, 11 Apr 2018 18:24:30 GMT\r\n\r\n{\"msg\": \"I'm JSON\"}",
				json:     `{"msg": "I'm JSON"}`,
			},
		}, {
			// connection error
			call: call{"POST", "http://127.0.0.1:55555/", ""},
			webhook: webhook{
				request:  "POST / HTTP/1.1\r\nHost: 127.0.0.1:55555\r\nUser-Agent: goflow-testing\r\nContent-Length: 0\r\nAccept-Encoding: gzip\r\n\r\n",
				response: "",
				json:     "",
			},
		},
	}

	for _, tc := range testCases {
		request, err := http.NewRequest(tc.call.method, tc.call.url, strings.NewReader(tc.call.body))
		require.NoError(t, err)

		svc, _ := session.Engine().Services().Webhook(session)
		c, err := svc.Call(session, request)

		if tc.isError {
			assert.Error(t, err, "expected error for call %s", tc.call)
		} else {
			assert.NoError(t, err, "unexpected error fetching %s", tc.call)

			assert.Equal(t, tc.call.url, c.Request.URL.String(), "URL mismatch for call %s", tc.call)
			assert.Equal(t, tc.call.method, c.Request.Method, "method mismatch for call %s", tc.call)
			assert.Equal(t, tc.webhook.request, string(c.RequestTrace), "request trace mismatch for call %s", tc.call)
			assert.Equal(t, tc.webhook.response, string(c.ResponseTrace), "response mismatch for call %s", tc.call)
			assert.Equal(t, tc.webhook.bodyIgnored, c.BodyIgnored, "body-ignored mismatch for call %s", tc.call)
			assert.Equal(t, tc.webhook.json, string(c.ResponseBody), "body mismatch for call %s", tc.call)
		}
	}
}

func TestRetries(t *testing.T) {
	session, _, err := test.CreateTestSession("", envs.RedactionPolicyNone)
	require.NoError(t, err)

	defer httpx.SetRequestor(httpx.DefaultRequestor)

	mocks := httpx.NewMockRequestor(map[string][]httpx.MockResponse{
		"http://temba.io/": []httpx.MockResponse{
			httpx.NewMockResponse(502, nil, "a", 1),
			httpx.NewMockResponse(200, nil, "b", 1),
		},
	})
	httpx.SetRequestor(mocks)

	request, err := http.NewRequest("GET", "http://temba.io/", strings.NewReader("BODY"))
	require.NoError(t, err)

	svc, _ := session.Engine().Services().Webhook(session)
	c, err := svc.Call(session, request)

	assert.Equal(t, 200, c.Response.StatusCode)
	assert.Equal(t, "GET / HTTP/1.1\r\nHost: temba.io\r\nUser-Agent: goflow-testing\r\nContent-Length: 4\r\nAccept-Encoding: gzip\r\n\r\nBODY", string(c.RequestTrace))
	assert.Equal(t, "HTTP/1.0 200 OK\r\nContent-Length: 1\r\n\r\nb", string(c.ResponseTrace))
}

func TestAccessRestrictions(t *testing.T) {
	retries := httpx.NewFixedRetries(5, 10)
	access := httpx.NewAccessConfig(10, []net.IP{net.IPv4(127, 0, 0, 1)})

	factory := webhooks.NewServiceFactory(http.DefaultClient, retries, access, map[string]string{"User-Agent": "Foo"}, 12345)
	svc, err := factory(nil)
	assert.NoError(t, err)

	request, _ := http.NewRequest("GET", "http://localhost/foo", nil)
	call, err := svc.Call(nil, request)

	// actual error becomes a call with a connection error
	assert.NoError(t, err)

	// should still have a trace.. just no response part
	assert.Equal(t, "GET /foo HTTP/1.1\r\nHost: localhost\r\nUser-Agent: Foo\r\nAccept-Encoding: gzip\r\n\r\n", string(call.RequestTrace))
	assert.Equal(t, "", string(call.ResponseTrace))
}
