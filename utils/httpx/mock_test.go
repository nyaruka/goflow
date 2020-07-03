package httpx_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/nyaruka/goflow/utils/httpx"
	"github.com/nyaruka/goflow/utils/jsonx"

	"github.com/stretchr/testify/assert"
)

func TestMockRequestor(t *testing.T) {
	defer httpx.SetRequestor(httpx.DefaultRequestor)

	// can create requestor with constructor
	requestor1 := httpx.NewMockRequestor(map[string][]httpx.MockResponse{
		"http://google.com": {
			httpx.NewMockResponse(200, nil, "this is google"),
			httpx.NewMockResponse(201, nil, "this is google again"),
		},
		"http://yahoo.com": {
			httpx.NewMockResponse(202, nil, "this is yahoo"),
			httpx.MockConnectionError,
		},
	})

	httpx.SetRequestor(requestor1)

	req1, _ := http.NewRequest("GET", "http://google.com", nil)
	response1, err := httpx.Do(http.DefaultClient, req1, nil, nil)
	assert.NoError(t, err)
	assert.Equal(t, 200, response1.StatusCode)

	body, err := ioutil.ReadAll(response1.Body)
	assert.NoError(t, err)
	assert.Equal(t, "this is google", string(body))

	assert.True(t, requestor1.HasUnused())

	// request another mocked URL
	req2, _ := http.NewRequest("GET", "http://yahoo.com", nil)
	response2, err := httpx.Do(http.DefaultClient, req2, nil, nil)
	assert.NoError(t, err)
	assert.Equal(t, 202, response2.StatusCode)

	// request second mock for first URL
	req3, _ := http.NewRequest("GET", "http://google.com", nil)
	response3, err := httpx.Do(http.DefaultClient, req3, nil, nil)
	assert.NoError(t, err)
	assert.Equal(t, 201, response3.StatusCode)

	// request mocked connection error
	req4, _ := http.NewRequest("GET", "http://yahoo.com", nil)
	response4, err := httpx.Do(http.DefaultClient, req4, nil, nil)
	assert.EqualError(t, err, "unable to connect to server")
	assert.Nil(t, response4)

	assert.False(t, requestor1.HasUnused())

	// panic if we've run out of mocks for a URL
	req5, _ := http.NewRequest("GET", "http://google.com", nil)
	assert.Panics(t, func() { httpx.Do(http.DefaultClient, req5, nil, nil) })
}

func TestMockRequestorMarshaling(t *testing.T) {
	// can create requestor with constructor
	requestor1 := httpx.NewMockRequestor(map[string][]httpx.MockResponse{
		"http://google.com": {
			httpx.NewMockResponse(200, nil, "this is google"),
			httpx.NewMockResponse(201, nil, "this is google again"),
		},
		"http://yahoo.com": {
			httpx.NewMockResponse(202, nil, "this is yahoo"),
			httpx.MockConnectionError,
		},
	})

	asJSON := []byte(`{
		"http://google.com": [
			{"status": 200, "body": "this is google"},
			{"status": 201, "body": "this is google again"}
		],
		"http://yahoo.com": [
			{"status": 202, "body": "this is yahoo"},
			{"status": 0, "body": ""}
		]
	}`)

	// test unmarshaling
	requestor2 := &httpx.MockRequestor{}
	err := json.Unmarshal(asJSON, requestor2)
	assert.NoError(t, err)
	assert.Equal(t, requestor1, requestor2)

	// test re-marshaling
	marshaled, err := jsonx.Marshal(requestor2)
	assert.JSONEq(t, string(asJSON), string(marshaled))
}
