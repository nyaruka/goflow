package httpx_test

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/nyaruka/goflow/utils/httpx"

	"github.com/stretchr/testify/assert"
)

func TestMockRequestor(t *testing.T) {
	defer httpx.SetRequestor(httpx.DefaultRequestor)

	httpx.SetRequestor(httpx.NewMockRequestor(map[string][]*http.Response{
		"http://google.com": []*http.Response{
			httpx.NewMockResponse(200, "this is google"),
			httpx.NewMockResponse(201, "this is google again"),
		},
		"http://yahoo.com": []*http.Response{
			httpx.NewMockResponse(202, "this is yahoo"),
		},
	}))

	req1, _ := http.NewRequest("GET", "http://google.com", nil)
	response1, err := httpx.Do(http.DefaultClient, req1)
	assert.NoError(t, err)
	assert.Equal(t, 200, response1.StatusCode)

	body, err := ioutil.ReadAll(response1.Body)
	assert.NoError(t, err)
	assert.Equal(t, "this is google", string(body))

	// request another mocked URL
	req2, _ := http.NewRequest("GET", "http://yahoo.com", nil)
	response2, err := httpx.Do(http.DefaultClient, req2)
	assert.NoError(t, err)
	assert.Equal(t, 202, response2.StatusCode)

	// request second mock for first URL
	req3, _ := http.NewRequest("GET", "http://google.com", nil)
	response3, err := httpx.Do(http.DefaultClient, req3)
	assert.NoError(t, err)
	assert.Equal(t, 201, response3.StatusCode)

	// error if we've run out of mocks for a URL
	req4, _ := http.NewRequest("GET", "http://google.com", nil)
	response4, err := httpx.Do(http.DefaultClient, req4)
	assert.EqualError(t, err, "missing mock for URL http://google.com")
	assert.Nil(t, response4)
}
