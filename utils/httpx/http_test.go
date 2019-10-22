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

	trace, err := httpx.DoTrace(http.DefaultClient, "GET", server.URL+"?cmd=success", nil, nil)
	assert.NoError(t, err)
	assert.Equal(t, "GET /?cmd=success HTTP/1.1\r\nHost: 127.0.0.1:52025\r\nUser-Agent: Go-http-client/1.1\r\nAccept-Encoding: gzip\r\n\r\n", string(trace.RequestTrace))
	assert.Equal(t, `{ "ok": "true" }`, string(trace.ResponseBody))
	assert.Equal(t, "HTTP/1.1 200 OK\r\nContent-Length: 16\r\nContent-Type: text/plain; charset=utf-8\r\nDate: Wed, 11 Apr 2018 18:24:30 GMT\r\n\r\n{ \"ok\": \"true\" }", string(trace.ResponseTrace))
	assert.Equal(t, time.Date(2019, 10, 7, 15, 21, 30, 123456789, time.UTC), trace.StartTime)
	assert.Equal(t, time.Date(2019, 10, 7, 15, 21, 31, 123456789, time.UTC), trace.EndTime)
}
