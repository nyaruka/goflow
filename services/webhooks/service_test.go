package webhooks_test

import (
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/nyaruka/gocommon/dates"
	"github.com/nyaruka/gocommon/httpx"
	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/services/webhooks"
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
	request  string
	response string
	body     string
}

func TestWebhookService(t *testing.T) {
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
				response: "HTTP/1.1 200 OK\r\nContent-Length: 16\r\nContent-Type: text/plain; charset=utf-8\r\nDate: Wed, 11 Apr 2018 18:24:30 GMT\r\n\r\n",
				body:     `{ "ok": "true" }`,
			},
		}, {
			// successful GET with valid JSON response body
			call: call{"GET", "http://127.0.0.1:49994/?cmd=textjs", ""},
			webhook: webhook{
				request:  "GET /?cmd=textjs HTTP/1.1\r\nHost: 127.0.0.1:49994\r\nUser-Agent: goflow-testing\r\nAccept-Encoding: gzip\r\n\r\n",
				response: "HTTP/1.1 200 OK\r\nContent-Length: 16\r\nContent-Type: text/javascript; charset=iso-8859-1\r\nDate: Wed, 11 Apr 2018 18:24:30 GMT\r\n\r\n",
				body:     `{ "ok": "true" }`,
			},
		}, {
			// successful POST without request body and valid JSON response body
			call: call{"POST", "http://127.0.0.1:49994/?cmd=success", ""},
			webhook: webhook{
				request:  "POST /?cmd=success HTTP/1.1\r\nHost: 127.0.0.1:49994\r\nUser-Agent: goflow-testing\r\nContent-Length: 0\r\nAccept-Encoding: gzip\r\n\r\n",
				response: "HTTP/1.1 200 OK\r\nContent-Length: 16\r\nContent-Type: text/plain; charset=utf-8\r\nDate: Wed, 11 Apr 2018 18:24:30 GMT\r\n\r\n",
				body:     `{ "ok": "true" }`,
			},
		}, {
			// successful POST with request body and valid JSON response body
			call: call{"POST", "http://127.0.0.1:49994/?cmd=success", `{"contact": "Bob"}`},
			webhook: webhook{
				request:  "POST /?cmd=success HTTP/1.1\r\nHost: 127.0.0.1:49994\r\nUser-Agent: goflow-testing\r\nContent-Length: 18\r\nAccept-Encoding: gzip\r\n\r\n{\"contact\": \"Bob\"}",
				response: "HTTP/1.1 200 OK\r\nContent-Length: 16\r\nContent-Type: text/plain; charset=utf-8\r\nDate: Wed, 11 Apr 2018 18:24:30 GMT\r\n\r\n",
				body:     `{ "ok": "true" }`,
			},
		}, {
			// successful GET with JSON response body containing escaped null chars (actual escaped nulls should be replaced with \ufffd)
			call: call{"GET", "http://127.0.0.1:49994/?cmd=badjson", ""},
			webhook: webhook{
				request:  "GET /?cmd=badjson HTTP/1.1\r\nHost: 127.0.0.1:49994\r\nUser-Agent: goflow-testing\r\nAccept-Encoding: gzip\r\n\r\n",
				response: "HTTP/1.1 200 OK\r\nContent-Length: 67\r\nContent-Type: application/json\r\nDate: Wed, 11 Apr 2018 18:24:30 GMT\r\n\r\n",
				body:     "{ \"bad\": \"null=\x00 escaped=\\u0000 double-escaped=\\\\u0000 badseq=\x80\x81\" }",
			},
		}, {
			// successful POST receiving gzipped non-JSON body
			call: call{"POST", "http://127.0.0.1:49994/?cmd=gzipped&content=Hello", ``},
			webhook: webhook{
				request:  "POST /?cmd=gzipped&content=Hello HTTP/1.1\r\nHost: 127.0.0.1:49994\r\nUser-Agent: goflow-testing\r\nContent-Length: 0\r\nAccept-Encoding: gzip\r\n\r\n",
				response: "HTTP/1.1 200 OK\r\nContent-Type: application/x-gzip\r\nDate: Wed, 11 Apr 2018 18:24:30 GMT\r\n\r\n",
				body:     `Hello`,
			},
		}, {
			// successful POST receiving gzipped JSON body
			call: call{"POST", "http://127.0.0.1:49994/?cmd=gzipped&content=%7B%22contact%22%3A%20%22Bob%22%7D", ``},
			webhook: webhook{
				request:  "POST /?cmd=gzipped&content=%7B%22contact%22%3A%20%22Bob%22%7D HTTP/1.1\r\nHost: 127.0.0.1:49994\r\nUser-Agent: goflow-testing\r\nContent-Length: 0\r\nAccept-Encoding: gzip\r\n\r\n",
				response: "HTTP/1.1 200 OK\r\nContent-Type: application/x-gzip\r\nDate: Wed, 11 Apr 2018 18:24:30 GMT\r\n\r\n",
				body:     `{"contact": "Bob"}`,
			},
		}, {
			// POST returning 503
			call: call{"POST", "http://127.0.0.1:49994/?cmd=unavailable", ""},
			webhook: webhook{
				request:  "POST /?cmd=unavailable HTTP/1.1\r\nHost: 127.0.0.1:49994\r\nUser-Agent: goflow-testing\r\nContent-Length: 0\r\nAccept-Encoding: gzip\r\n\r\n",
				response: "HTTP/1.1 503 Service Unavailable\r\nContent-Length: 37\r\nContent-Type: text/plain; charset=utf-8\r\nDate: Wed, 11 Apr 2018 18:24:30 GMT\r\n\r\n",
				body:     `{ "errors": ["service unavailable"] }`,
			},
		}, {
			// GET returning text body larger than allowed
			call:    call{"GET", "http://127.0.0.1:49994/?cmd=binary&size=11000", ""},
			isError: true,
		}, {
			// GET returning non-JSON body
			call: call{"GET", "http://127.0.0.1:49994/?cmd=typeless&content=kthxbai", ""},
			webhook: webhook{
				request:  "GET /?cmd=typeless&content=kthxbai HTTP/1.1\r\nHost: 127.0.0.1:49994\r\nUser-Agent: goflow-testing\r\nAccept-Encoding: gzip\r\n\r\n",
				response: "HTTP/1.1 200 OK\r\nContent-Length: 7\r\nContent-Type: \r\nDate: Wed, 11 Apr 2018 18:24:30 GMT\r\n\r\n",
				body:     "kthxbai",
			},
		}, {
			// GET returning JSON body but an empty content-type header
			call: call{"GET", "http://127.0.0.1:49994/?cmd=typeless&content=%7B%22msg%22%3A%20%22I%27m%20JSON%22%7D", ""},
			webhook: webhook{
				request:  "GET /?cmd=typeless&content=%7B%22msg%22%3A%20%22I%27m%20JSON%22%7D HTTP/1.1\r\nHost: 127.0.0.1:49994\r\nUser-Agent: goflow-testing\r\nAccept-Encoding: gzip\r\n\r\n",
				response: "HTTP/1.1 200 OK\r\nContent-Length: 19\r\nContent-Type: \r\nDate: Wed, 11 Apr 2018 18:24:30 GMT\r\n\r\n",
				body:     `{"msg": "I'm JSON"}`,
			},
		}, {
			// connection error
			call: call{"POST", "http://127.0.0.1:55555/", ""},
			webhook: webhook{
				request:  "POST / HTTP/1.1\r\nHost: 127.0.0.1:55555\r\nUser-Agent: goflow-testing\r\nContent-Length: 0\r\nAccept-Encoding: gzip\r\n\r\n",
				response: "",
				body:     "",
			},
		},
	}

	for _, tc := range testCases {
		request, err := http.NewRequest(tc.call.method, tc.call.url, strings.NewReader(tc.call.body))
		require.NoError(t, err)

		svc, _ := session.Engine().Services().Webhook(session.Assets())
		c, err := svc.Call(request)

		if tc.isError {
			assert.Error(t, err, "expected error for call %s", tc.call)
		} else {
			assert.NoError(t, err, "unexpected error fetching %s", tc.call)

			assert.Equal(t, tc.call.url, c.Request.URL.String(), "URL mismatch for call %s", tc.call)
			assert.Equal(t, tc.call.method, c.Request.Method, "method mismatch for call %s", tc.call)
			assert.Equal(t, tc.webhook.request, string(c.RequestTrace), "request trace mismatch for call %s", tc.call)
			assert.Equal(t, tc.webhook.response, string(c.ResponseTrace), "response mismatch for call %s", tc.call)
			assert.Equal(t, tc.webhook.body, string(c.ResponseBody), "body mismatch for call %s", tc.call)
		}
	}
}

func TestRetries(t *testing.T) {
	_, session, _ := test.NewSessionBuilder().MustBuild()

	defer httpx.SetRequestor(httpx.DefaultRequestor)

	mocks := httpx.NewMockRequestor(map[string][]*httpx.MockResponse{
		"http://temba.io/": {
			httpx.NewMockResponse(502, nil, []byte("a")),
			httpx.NewMockResponse(200, nil, []byte("b")),
		},
	})
	httpx.SetRequestor(mocks)

	request, err := http.NewRequest("GET", "http://temba.io/", strings.NewReader("BODY"))
	require.NoError(t, err)

	svc, _ := session.Engine().Services().Webhook(session.Assets())
	c, err := svc.Call(request)
	require.NoError(t, err)

	assert.Equal(t, 200, c.Response.StatusCode)
	assert.Equal(t, "GET / HTTP/1.1\r\nHost: temba.io\r\nUser-Agent: goflow-testing\r\nContent-Length: 4\r\nAccept-Encoding: gzip\r\n\r\nBODY", string(c.RequestTrace))
	assert.Equal(t, "HTTP/1.0 200 OK\r\nContent-Length: 1\r\n\r\n", string(c.ResponseTrace))
	assert.Equal(t, "b", string(c.ResponseBody))
}

func TestAccessRestrictions(t *testing.T) {
	retries := httpx.NewFixedRetries(5, 10)
	access := httpx.NewAccessConfig(10, []net.IP{net.IPv4(127, 0, 0, 1)}, nil)

	factory := webhooks.NewServiceFactory(http.DefaultClient, retries, access, map[string]string{"User-Agent": "Foo"}, 12345)
	svc, err := factory(nil)
	assert.NoError(t, err)

	request, _ := http.NewRequest("GET", "http://localhost/foo", nil)
	call, err := svc.Call(request)

	// actual error becomes a call with a connection error
	assert.NoError(t, err)

	// should still have a trace.. just no response part
	assert.Equal(t, "GET /foo HTTP/1.1\r\nHost: localhost\r\nUser-Agent: Foo\r\nAccept-Encoding: gzip\r\n\r\n", string(call.RequestTrace))
	assert.Equal(t, "", string(call.ResponseTrace))
}

func TestGzipEncoding(t *testing.T) {
	_, session, _ := test.NewSessionBuilder().MustBuild()

	defer dates.SetNowFunc(time.Now)
	dates.SetNowFunc(dates.NewSequentialNow(time.Date(2019, 10, 7, 15, 21, 30, 123456789, time.UTC), time.Second))

	server := test.NewTestHTTPServer(52025)

	request, err := http.NewRequest("GET", server.URL+"?cmd=gzipped&content=Hello", nil)
	require.NoError(t, err)

	request.Header.Set("Accept-Encoding", "gzip")

	svc, _ := session.Engine().Services().Webhook(session.Assets())
	c, err := svc.Call(request)
	require.NoError(t, err)

	// check that gzip decompression happens transparently
	assert.Equal(t, 200, c.Response.StatusCode)
	assert.Equal(t, "GET /?cmd=gzipped&content=Hello HTTP/1.1\r\nHost: 127.0.0.1:52025\r\nUser-Agent: goflow-testing\r\nAccept-Encoding: gzip\r\n\r\n", string(c.RequestTrace))
	assert.Equal(t, "HTTP/1.1 200 OK\r\nContent-Type: application/x-gzip\r\nDate: Wed, 11 Apr 2018 18:24:30 GMT\r\n\r\n", string(c.ResponseTrace))
	assert.Equal(t, "Hello", string(c.ResponseBody))

	assert.Equal(t, "GET /?cmd=gzipped&content=Hello HTTP/1.1\r\nHost: 127.0.0.1:52025\r\nUser-Agent: goflow-testing\r\nAccept-Encoding: gzip\r\n\r\n", c.SanitizedRequest("..."))
	assert.Equal(t, "HTTP/1.1 200 OK\r\nContent-Type: application/x-gzip\r\nDate: Wed, 11 Apr 2018 18:24:30 GMT\r\n\r\nHello", c.SanitizedResponse("..."))
}

func TestWebhookResponseWithEscapes(t *testing.T) {
	defer httpx.SetRequestor(httpx.DefaultRequestor)

	mocks := httpx.NewMockRequestor(map[string][]*httpx.MockResponse{
		"http://cheapcontactlookups.com": {
			httpx.NewMockResponse(200, nil, []byte(`{"name": "01\\02\\03", "joined": "04\\05\\06"}`)),
		},
	})
	httpx.SetRequestor(mocks)

	_, session, _ := test.NewSessionBuilder().
		WithAssetsPath("testdata/webhook_flow.json").
		WithFlow("bb38eefb-3cd9-4f80-9867-9c84ae276f7a").MustBuild()

	joined := session.Assets().Fields().Get("joined")

	assert.Equal(t, flows.SessionStatusCompleted, session.Status())
	assert.Equal(t, `01\02\03`, session.Contact().Name())
	assert.Equal(t, types.NewXText(`04\05\06`), session.Contact().Fields().Get(joined).Text)

	// check nothing became an escaped NULL
	assert.NotContains(t, string(jsonx.MustMarshal(session)), `\u0000`)
}
