package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func testHTTPRequest(t *testing.T, method string, url string, data string) (int, []byte) {
	var reqBody io.Reader
	if data != "" {
		reqBody = strings.NewReader(data)
	}

	req, err := http.NewRequest(method, url, reqBody)
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	return resp.StatusCode, body
}

func assertErrorResponse(t *testing.T, body []byte, expectedErrors []string) {
	errResp := &errorResponse{}
	err := json.Unmarshal(body, &errResp)
	assert.NoError(t, err)
	assert.Equal(t, expectedErrors, errResp.Text)
}

func assertExpressionResponse(t *testing.T, body []byte, expectedResult string, expectedErrors []string) {
	expResp := &expressionResponse{}
	err := json.Unmarshal(body, &expResp)
	assert.NoError(t, err)
	assert.Equal(t, expectedResult, expResp.Result)
	assert.Equal(t, expectedErrors, expResp.Errors)
}

func TestFlowServer(t *testing.T) {
	flowServer := NewFlowServer(NewTestConfig(), logrus.New())
	flowServer.Start()
	defer flowServer.Stop()

	// wait for server to come up
	time.Sleep(100 * time.Millisecond)

	// hit our home page
	status, body := testHTTPRequest(t, "GET", "http://localhost:8080/", "")
	assert.Equal(t, 200, status)
	assert.Contains(t, string(body), "Echo Flow")

	// hit our version endpoint
	status, body = testHTTPRequest(t, "GET", "http://localhost:8080/version", "")
	assert.Equal(t, 200, status)
	assert.Contains(t, string(body), "version")

	// try the expression endpoint
	status, body = testHTTPRequest(t, "POST", "http://localhost:8080/expression", `{"expression": "@(1 + 2)", "context": {}}`)
	assert.Equal(t, 200, status)
	assertExpressionResponse(t, body, "3", []string{})

	// try the expression endpoint with a bad expression
	status, body = testHTTPRequest(t, "POST", "http://localhost:8080/expression", `{"expression": "@(foo + 2)", "context": {}}`)
	assert.Equal(t, 200, status)
	assertExpressionResponse(t, body, "", []string{"Key path not found"})

	// try to GET the start endpoint
	status, body = testHTTPRequest(t, "GET", "http://localhost:8080/flow/start", "")
	assert.Equal(t, 405, status)
	assertErrorResponse(t, body, []string{"method not allowed"})

	// try POSTing nothing to the start endpoint
	status, body = testHTTPRequest(t, "POST", "http://localhost:8080/flow/start", "")
	assert.Equal(t, 400, status)
	assertErrorResponse(t, body, []string{"unexpected end of JSON input"})

	// try POSTing empty JSON to the start endpoint
	status, body = testHTTPRequest(t, "POST", "http://localhost:8080/flow/start", "{}")
	assert.Equal(t, 400, status)
	assertErrorResponse(t, body, []string{"field 'asset_urls' is required", "field 'trigger' is required"})

	// try POSTing an incomplete trigger to the start endpoint
	status, body = testHTTPRequest(t, "POST", "http://localhost:8080/flow/start", `{"asset_urls": {}, "trigger": {"type": "manual"}}`)
	assert.Equal(t, 400, status)
	assertErrorResponse(t, body, []string{"field 'flow' on 'trigger[type=manual]' is required", "field 'triggered_on' on 'trigger[type=manual]' is required"})

	// try POSTing a trigger with an invalid flow asset
	status, body = testHTTPRequest(t, "POST", "http://localhost:8080/flow/start", `{
		"assets": [
			{ 
				"type": "flow",
				"url": "http://testserver/assets/flow/76f0a02f-3b75-4b86-9064-e9195e1b3a02",
				"content": {
					"uuid": "76f0a02f-3b75-4b86-9064-e9195e1b3a02",
					"name": "Test Flow",
					"language": "eng",
					"nodes": [
						{
							"uuid": "a58be63b-907d-4a1a-856b-0bb5579d7507",
							"exits": [
								{
									"uuid": "37d8813f-1402-4ad2-9cc2-e9054a96525b",
									"label": "Default",
									"destination_node_uuid": "714f1409-486e-4e8e-bb08-23e2943ef9f6"
								}
							]
						}
					]
				}
			}
		],
		"asset_urls": {
			"flow": "http://testserver/assets/flow"
		},
		"trigger": {
			"type": "manual",
			"flow": {"uuid": "76f0a02f-3b75-4b86-9064-e9195e1b3a02", "name": "Test Flow"},
			"triggered_on": "2000-01-01T00:00:00.000000000-00:00"
		}
	}`)
	assert.Equal(t, 400, status)
	assertErrorResponse(t, body, []string{"unable to read asset[url=http://testserver/assets/flow/76f0a02f-3b75-4b86-9064-e9195e1b3a02]: destination 714f1409-486e-4e8e-bb08-23e2943ef9f6 of exit[uuid=37d8813f-1402-4ad2-9cc2-e9054a96525b] isn't a known node"})
}
