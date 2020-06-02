package httpx_test

import (
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/nyaruka/goflow/utils/httpx"
	"github.com/stretchr/testify/assert"
)

func TestAccessConfig(t *testing.T) {
	defer httpx.SetRequestor(httpx.DefaultRequestor)

	access := httpx.NewAccessConfig(30*time.Second, []net.IP{net.ParseIP("127.0.0.1"), net.ParseIP("::1")})

	httpx.SetRequestor(httpx.NewMockRequestor(map[string][]httpx.MockResponse{
		"https://nyaruka.com": {
			httpx.NewMockResponse(200, nil, ``),
		},
	}))

	tests := []struct {
		url string
		err string
	}{
		{"https://nyaruka.com", ""},

		{"https://localhost/path", "request to localhost denied"},
		{"https://LOCALHOST:80", "request to LOCALHOST denied"},
		{"http://foo.localtest.me", "request to foo.localtest.me denied"},
		{"https://127.0.0.1", "request to 127.0.0.1 denied"},
		{"https://127.0.00.1", "request to 127.0.00.1 denied"},
		{"https://[::1]:80", "request to ::1 denied"},
		{"https://[0:0:0:0:0:0:0:1]:80", "request to 0:0:0:0:0:0:0:1 denied"},
		{"https://[0000:0000:0000:0000:0000:0000:0000:0001]:80", "request to 0000:0000:0000:0000:0000:0000:0000:0001 denied"},
	}
	for _, tc := range tests {
		request, _ := http.NewRequest("GET", tc.url, nil)
		_, err := httpx.DoTrace(http.DefaultClient, request, nil, access, -1)

		if tc.err != "" {
			assert.EqualError(t, err, tc.err, "error message mismatch for url %s", tc.url)
		} else {
			assert.NoError(t, err)
		}
	}
}
