package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/nyaruka/goflow/test"

	"github.com/stretchr/testify/suite"
)

type ServerTestSuite struct {
	suite.Suite
	flowServer *FlowServer
	httpServer *httptest.Server
}

func (ts *ServerTestSuite) SetupSuite() {
	ts.flowServer = NewFlowServer(NewDefaultConfig())
	ts.flowServer.Start()

	// wait for server to come up
	time.Sleep(100 * time.Millisecond)

	ts.httpServer = test.NewTestHTTPServer(49993)
}

func (ts *ServerTestSuite) TearDownSuite() {
	ts.flowServer.Stop()

	ts.httpServer.Close()
}

func (ts *ServerTestSuite) testHTTPRequest(method string, url string, data string) (int, []byte) {
	var reqBody io.Reader
	if data != "" {
		reqBody = strings.NewReader(data)
	}

	req, err := http.NewRequest(method, url, reqBody)
	resp, err := http.DefaultClient.Do(req)
	ts.Require().NoError(err)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	ts.Require().NoError(err)
	return resp.StatusCode, body
}

func (ts *ServerTestSuite) assertErrorResponse(body []byte, expectedErrors []string) {
	errResp := &errorResponse{}
	err := json.Unmarshal(body, &errResp)
	ts.Require().NoError(err)
	ts.Equal(expectedErrors, errResp.Text)
}

func (ts *ServerTestSuite) TestMigrate() {
	// try to GET the migrate endpoint
	status, body := ts.testHTTPRequest("GET", "http://localhost:8800/flow/migrate", "")
	ts.Equal(405, status)
	ts.assertErrorResponse(body, []string{"method not allowed"})

	// try POSTing nothing to the migrate endpoint
	status, body = ts.testHTTPRequest("POST", "http://localhost:8800/flow/migrate", "")
	ts.Equal(400, status)
	ts.assertErrorResponse(body, []string{"unexpected end of JSON input"})

	// try POSTing empty JSON to the start endpoint
	status, body = ts.testHTTPRequest("POST", "http://localhost:8800/flow/migrate", "{}")
	ts.Equal(400, status)
	ts.assertErrorResponse(body, []string{"missing flow element"})
}

func TestServerSuite(t *testing.T) {
	suite.Run(t, new(ServerTestSuite))
}
