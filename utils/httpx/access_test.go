package httpx_test

import (
	"net/http"
	"testing"

	"github.com/nyaruka/goflow/utils/httpx"
	"github.com/stretchr/testify/assert"
)

func TestDisallowedHosts(t *testing.T) {
	defer httpx.SetRequestor(httpx.DefaultRequestor)

	disallowedHosts := []string{"localhost", "127.0.0.1", "::1"}

	httpx.SetRequestor(httpx.NewMockRequestor(map[string][]httpx.MockResponse{
		"https://temba.io": []httpx.MockResponse{
			httpx.NewMockResponse(200, nil, ``, 1),
		},
	}))

	call := func(url string) (*httpx.Trace, error) {
		request, _ := http.NewRequest("GET", url, nil)
		return httpx.DoTrace(http.DefaultClient, request, nil, httpx.NewAccessConfig(disallowedHosts), -1)
	}

	_, err := call("https://temba.io")
	assert.NoError(t, err)

	_, err = call("https://localhost/path")
	assert.EqualError(t, err, "request to localhost denied")

	_, err = call("https://LOCALHOST:80")
	assert.EqualError(t, err, "request to LOCALHOST denied")

	_, err = call("https://127.0.0.1")
	assert.EqualError(t, err, "request to 127.0.0.1 denied")

	_, err = call("https://127.0.00.1")
	assert.EqualError(t, err, "request to 127.0.00.1 denied")

	_, err = call("https://[::1]:80")
	assert.EqualError(t, err, "request to ::1 denied")

	_, err = call("https://[0:0:0:0:0:0:0:1]:80")
	assert.EqualError(t, err, "request to 0:0:0:0:0:0:0:1 denied")

	_, err = call("https://[0000:0000:0000:0000:0000:0000:0000:0001]:80")
	assert.EqualError(t, err, "request to 0000:0000:0000:0000:0000:0000:0000:0001 denied")
}
