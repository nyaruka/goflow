package httpx_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/nyaruka/goflow/test"
	"github.com/nyaruka/goflow/utils/dates"
	"github.com/nyaruka/goflow/utils/httpx"
	"github.com/stretchr/testify/assert"
)

func TestDoTrace(t *testing.T) {
	defer dates.SetNowSource(dates.DefaultNowSource)

	dates.SetNowSource(dates.NewSequentialNowSource(time.Date(2019, 10, 7, 15, 21, 30, 123456789, time.UTC)))

	server := test.NewTestHTTPServer(52025)

	trace, err := httpx.DoTrace(http.DefaultClient, "GET", server.URL+"?cmd=success", nil, nil, nil)
	assert.NoError(t, err)
	assert.Equal(t, "GET /?cmd=success HTTP/1.1\r\nHost: 127.0.0.1:52025\r\nUser-Agent: Go-http-client/1.1\r\nAccept-Encoding: gzip\r\n\r\n", string(trace.RequestTrace))
	assert.Equal(t, `{ "ok": "true" }`, string(trace.ResponseBody))
	assert.Equal(t, "HTTP/1.1 200 OK\r\nContent-Length: 16\r\nContent-Type: text/plain; charset=utf-8\r\nDate: Wed, 11 Apr 2018 18:24:30 GMT\r\n\r\n{ \"ok\": \"true\" }", string(trace.ResponseTrace))
	assert.Equal(t, time.Date(2019, 10, 7, 15, 21, 30, 123456789, time.UTC), trace.StartTime)
	assert.Equal(t, time.Date(2019, 10, 7, 15, 21, 31, 123456789, time.UTC), trace.EndTime)
}

func TestDoWithRetries(t *testing.T) {
	defer httpx.SetRequestor(httpx.DefaultRequestor)

	mocks := httpx.NewMockRequestor(map[string][]httpx.MockResponse{
		"http://temba.io/1/": []httpx.MockResponse{
			httpx.NewMockResponse(502, "a", nil),
		},
		"http://temba.io/2/": []httpx.MockResponse{
			httpx.NewMockResponse(503, "a", nil),
			httpx.NewMockResponse(504, "b", nil),
			httpx.NewMockResponse(505, "c", nil),
		},
		"http://temba.io/3/": []httpx.MockResponse{
			httpx.NewMockResponse(200, "a", nil),
		},
		"http://temba.io/4/": []httpx.MockResponse{
			httpx.NewMockResponse(502, "a", nil),
		},
		"http://temba.io/5/": []httpx.MockResponse{
			httpx.NewMockResponse(502, "a", nil),
			httpx.NewMockResponse(200, "b", nil),
		},
		"http://temba.io/6/": []httpx.MockResponse{
			httpx.NewMockResponse(429, "a", map[string]string{"Retry-After": "1"}),
			httpx.NewMockResponse(201, "b", nil),
		},
		"http://temba.io/7/": []httpx.MockResponse{
			httpx.NewMockResponse(429, "a", map[string]string{"Retry-After": "100"}),
		},
	})
	httpx.SetRequestor(mocks)

	// no retry config
	trace, err := httpx.DoTrace(http.DefaultClient, "GET", "http://temba.io/1/", nil, nil, nil)
	assert.NoError(t, err)
	assert.Equal(t, 502, trace.Response.StatusCode)

	// a retry config which can make 3 attempts
	retries := &httpx.RetryConfig{
		Delays:      []time.Duration{1 * time.Millisecond, 2 * time.Millisecond},
		ShouldRetry: httpx.DefaultShouldRetry,
	}

	// retrying thats ends with failure
	trace, err = httpx.DoTrace(http.DefaultClient, "GET", "http://temba.io/2/", nil, nil, retries)
	assert.NoError(t, err)
	assert.Equal(t, 505, trace.Response.StatusCode)

	// retrying not needed
	trace, err = httpx.DoTrace(http.DefaultClient, "GET", "http://temba.io/3/", nil, nil, retries)
	assert.NoError(t, err)
	assert.Equal(t, 200, trace.Response.StatusCode)

	// retrying not used for POSTs
	trace, err = httpx.DoTrace(http.DefaultClient, "POST", "http://temba.io/4/", nil, nil, retries)
	assert.NoError(t, err)
	assert.Equal(t, 502, trace.Response.StatusCode)

	// unless idempotency declared via request header
	trace, err = httpx.DoTrace(http.DefaultClient, "POST", "http://temba.io/5/", nil, map[string]string{"Idempotency-Key": "123"}, retries)
	assert.NoError(t, err)
	assert.Equal(t, 200, trace.Response.StatusCode)

	// a retry config which can make 2 attempts (need a longer delay so that the Retry-After header value can be used)
	retries = &httpx.RetryConfig{
		Delays:      []time.Duration{1 * time.Second},
		ShouldRetry: httpx.DefaultShouldRetry,
	}

	// retrying due to Retry-After header
	trace, err = httpx.DoTrace(http.DefaultClient, "POST", "http://temba.io/6/", nil, nil, retries)
	assert.NoError(t, err)
	assert.Equal(t, 201, trace.Response.StatusCode)

	// ignoring Retry-After header when it's too long
	trace, err = httpx.DoTrace(http.DefaultClient, "GET", "http://temba.io/7/", nil, nil, retries)
	assert.NoError(t, err)
	assert.Equal(t, 429, trace.Response.StatusCode)

	assert.False(t, mocks.HasUnused())
}

func TestParseRetryAfter(t *testing.T) {
	defer dates.SetNowSource(dates.DefaultNowSource)

	dates.SetNowSource(dates.NewFixedNowSource(time.Date(2020, 1, 7, 15, 10, 30, 500000000, time.UTC)))

	assert.Equal(t, 0*time.Second, httpx.ParseRetryAfter("x"))
	assert.Equal(t, 0*time.Second, httpx.ParseRetryAfter("0"))
	assert.Equal(t, 10*time.Second, httpx.ParseRetryAfter("10"))
	assert.Equal(t, 10*time.Second, httpx.ParseRetryAfter("10"))
	assert.Equal(t, 4500*time.Millisecond, httpx.ParseRetryAfter("Wed, 07 Jan 2020 15:10:35 GMT")) // 4.5 seconds in future
	assert.Equal(t, 0*time.Second, httpx.ParseRetryAfter("Wed, 07 Jan 2020 15:10:25 GMT"))         // 5.5 seconds in the past
}
