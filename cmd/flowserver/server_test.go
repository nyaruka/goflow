package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/flows/inputs"

	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/engine"
	"github.com/nyaruka/goflow/utils"
)

var testStructurallyInvalidFlowAssets = `[
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
]`

var testFlowMissingGroupAssets = `[
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
					"actions": [
						{
							"uuid": "ad154980-7bf7-4ab8-8728-545fd6378912",
							"type": "add_to_group",
							"groups": [
								{
									"uuid": "77a1bb5c-92f7-42bc-8a54-d21c1536ebc0",
									"name": "Testers"
								}
							]
						}
					],
					"exits": [
						{
							"uuid": "37d8813f-1402-4ad2-9cc2-e9054a96525b",
							"label": "Default",
							"destination_node_uuid": null
						}
					]
				}
			]
		}
	},
	{
		"type": "group",
		"url": "http://testserver/assets/group",
		"content": [
			{
				"uuid": "2aad21f6-30b7-42c5-bd7f-1b720c154817",
				"name": "Survey Audience"
			}
		],
		"is_set": true
	}
]`

var testValidFlowWithNoWaitAssets = `[
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
					"actions": [],
					"exits": [
						{
							"uuid": "37d8813f-1402-4ad2-9cc2-e9054a96525b",
							"label": "Default",
							"destination_node_uuid": null
						}
					]
				}
			]
		}
	}
]`

var testValidFlowWithWaitAssets = `[
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
					"wait": {
						"type": "msg",
						"timeout": 600
					},
					"router": {
						"type": "switch",
						"default_exit_uuid": "0680b01f-ba0b-48f4-a688-d2f963130126",
						"result_name": "name",
						"operand": "@run.input",
						"cases": [
							{
								"uuid": "5d6abc80-39e7-4620-9988-a2447bffe526",
								"type": "has_text",
								"exit_uuid": "37d8813f-1402-4ad2-9cc2-e9054a96525b"
							}
						]
					},
					"exits": [
						{
							"uuid": "37d8813f-1402-4ad2-9cc2-e9054a96525b",
							"label": "Not Empty",
							"destination_node_uuid": null
						},
						{
							"uuid": "0680b01f-ba0b-48f4-a688-d2f963130126",
							"label": "Other",
							"destination_node_uuid": null
						}
					]
				}
			]
		}
	}
]`

var startRequestTemplate = `{
	"assets": %s,
	"asset_server": {
		"type_urls": {
			"flow": "http://testserver/assets/flow",
			"group": "http://testserver/assets/group"
		}
	},
	"trigger": {
		"type": "manual",
		"flow": {"uuid": "76f0a02f-3b75-4b86-9064-e9195e1b3a02", "name": "Test Flow"},
		"triggered_on": "2000-01-01T00:00:00.000000000-00:00"
	}
}`

type ServerTestSuite struct {
	suite.Suite
	flowServer  *FlowServer
	assetServer engine.AssetServer
}

func (ts *ServerTestSuite) SetupSuite() {
	ts.assetServer = engine.NewMockAssetServer()

	ts.flowServer = NewFlowServer(NewDefaultConfig(), logrus.New())
	ts.flowServer.Start()

	// wait for server to come up
	time.Sleep(100 * time.Millisecond)
}

func (ts *ServerTestSuite) TearDownSuite() {
	ts.flowServer.Stop()
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

func (ts *ServerTestSuite) assertExpressionResponse(body []byte, expectedResult string, expectedErrors []string) {
	expResp := &expressionResponse{}
	err := json.Unmarshal(body, &expResp)
	ts.Require().NoError(err)
	ts.Equal(expectedResult, expResp.Result)
	ts.Equal(expectedErrors, expResp.Errors)
}

func (ts *ServerTestSuite) parseSessionResponse(assetCache *engine.AssetCache, body []byte) (flows.Session, []map[string]interface{}) {
	envelope := struct {
		Session json.RawMessage
		Log     []map[string]interface{}
	}{}
	err := json.Unmarshal(body, &envelope)
	ts.Require().NoError(err)

	session, err := engine.ReadSession(assetCache, ts.assetServer, envelope.Session)
	ts.Require().NoError(err)

	return session, envelope.Log
}

func (ts *ServerTestSuite) buildResumeRequest(assetsJSON string, session flows.Session, events []flows.Event) string {
	sessionJSON, err := json.Marshal(session)
	ts.Require().NoError(err)

	eventEnvelopes := make([]*utils.TypedEnvelope, len(events))
	for e := range events {
		eventEnvelopes[e], err = utils.EnvelopeFromTyped(events[e])
		ts.Require().NoError(err)
	}

	assets := json.RawMessage(assetsJSON)
	assetServer, _ := json.Marshal(engine.NewMockAssetServer())

	request := &resumeRequest{
		Assets:      &assets,
		AssetServer: assetServer,
		Session:     sessionJSON,
		Events:      eventEnvelopes,
	}

	requestJSON, err := json.Marshal(request)
	ts.Require().NoError(err)
	return string(requestJSON)
}

func (ts *ServerTestSuite) TestHomePages() {
	// hit our home page
	status, body := ts.testHTTPRequest("GET", "http://localhost:8800/", "")
	ts.Equal(200, status)
	ts.Contains(string(body), "Echo Flow")

	// test the example start request on the home page
	startJSON, err := ioutil.ReadFile("static/start.json")
	ts.Require().NoError(err)

	status, body = ts.testHTTPRequest("POST", "http://localhost:8800/flow/start", string(startJSON))
	ts.Equal(200, status)
	ts.Contains(string(body), "You said 'Let's go thrifting!'")

	// hit our version endpoint
	status, body = ts.testHTTPRequest("GET", "http://localhost:8800/version", "")
	ts.Equal(200, status)
	ts.Contains(string(body), "version")
}

func (ts *ServerTestSuite) TestExpression() {
	// try the expression endpoint with a valid expression
	status, body := ts.testHTTPRequest("POST", "http://localhost:8800/expression", `{"expression": "@(1 + 2)", "context": {}}`)
	ts.Equal(200, status)
	ts.assertExpressionResponse(body, "3", []string{})

	// try the expression endpoint with an unparseable expression... which we treat as not an expression
	status, body = ts.testHTTPRequest("POST", "http://localhost:8800/expression", `{"expression": "@(1 + 2", "context": {}}`)
	ts.Equal(200, status)
	ts.assertExpressionResponse(body, "@(1 + 2", []string{})

	// try the expression endpoint with a missing variable
	status, body = ts.testHTTPRequest("POST", "http://localhost:8800/expression", `{"expression": "@(foo + 2)", "context": {}}`)
	ts.Equal(200, status)
	ts.assertExpressionResponse(body, "", []string{"no such variable: foo"})
}

func (ts *ServerTestSuite) TestFlowStartAndResume() {
	// try to GET the start endpoint
	status, body := ts.testHTTPRequest("GET", "http://localhost:8800/flow/start", "")
	ts.Equal(405, status)
	ts.assertErrorResponse(body, []string{"method not allowed"})

	// ry POSTing nothing to the start endpoint
	status, body = ts.testHTTPRequest("POST", "http://localhost:8800/flow/start", "")
	ts.Equal(400, status)
	ts.assertErrorResponse(body, []string{"unexpected end of JSON input"})

	// try POSTing empty JSON to the start endpoint
	status, body = ts.testHTTPRequest("POST", "http://localhost:8800/flow/start", "{}")
	ts.Equal(400, status)
	ts.assertErrorResponse(body, []string{"field 'asset_server' is required", "field 'trigger' is required"})

	// try POSTing an incomplete trigger to the start endpoint
	status, body = ts.testHTTPRequest("POST", "http://localhost:8800/flow/start", `{"asset_server": {}, "trigger": {"type": "manual"}}`)
	ts.Equal(400, status)
	ts.assertErrorResponse(body, []string{"field 'flow' on 'trigger[type=manual]' is required", "field 'triggered_on' on 'trigger[type=manual]' is required"})

	// try POSTing to the start endpoint a structurally invalid flow asset
	requestBody := fmt.Sprintf(startRequestTemplate, testStructurallyInvalidFlowAssets)
	status, body = ts.testHTTPRequest("POST", "http://localhost:8800/flow/start", requestBody)
	ts.Equal(400, status)
	ts.assertErrorResponse(body, []string{"unable to read asset[url=http://testserver/assets/flow/76f0a02f-3b75-4b86-9064-e9195e1b3a02]: destination 714f1409-486e-4e8e-bb08-23e2943ef9f6 of exit[uuid=37d8813f-1402-4ad2-9cc2-e9054a96525b] isn't a known node"})

	// try POSTing to the start endpoint a flow asset that references a non-existent group asset
	requestBody = fmt.Sprintf(startRequestTemplate, testFlowMissingGroupAssets)
	status, body = ts.testHTTPRequest("POST", "http://localhost:8800/flow/start", requestBody)
	ts.Equal(400, status)
	ts.assertErrorResponse(body, []string{"validation failed for flow[uuid=76f0a02f-3b75-4b86-9064-e9195e1b3a02]: validation failed for action[uuid=ad154980-7bf7-4ab8-8728-545fd6378912, type=add_to_group]: no such group with uuid '77a1bb5c-92f7-42bc-8a54-d21c1536ebc0'"})

	// POST to the start endpoint with a valid flow with no wait (it should complete)
	requestBody = fmt.Sprintf(startRequestTemplate, testValidFlowWithNoWaitAssets)
	status, body = ts.testHTTPRequest("POST", "http://localhost:8800/flow/start", requestBody)
	ts.Equal(200, status)

	session, _ := ts.parseSessionResponse(ts.flowServer.assetCache, body)
	ts.Equal(flows.SessionStatus("completed"), session.Status())

	// try to resume this completed session but with no caller events
	status, body = ts.testHTTPRequest("POST", "http://localhost:8800/flow/resume", ts.buildResumeRequest(`[]`, session, []flows.Event{}))
	ts.Equal(400, status)
	ts.assertErrorResponse(body, []string{"field 'events' must have a minimum of 1 items"})

	// try to resume this completed session
	msg := inputs.NewMsgInput(flows.InputUUID(uuid.NewV4().String()), nil, time.Now(), urns.NewTelegramURN(1234567, "bob"), "hello", nil)
	status, body = ts.testHTTPRequest("POST", "http://localhost:8800/flow/resume", ts.buildResumeRequest(`[]`, session, []flows.Event{
		events.NewMsgReceivedEvent(msg),
	}))
	ts.Equal(400, status)
	ts.assertErrorResponse(body, []string{"only waiting sessions can be resumed"})

	// start another session on a flow that will wait for input
	requestBody = fmt.Sprintf(startRequestTemplate, testValidFlowWithWaitAssets)
	status, body = ts.testHTTPRequest("POST", "http://localhost:8800/flow/start", requestBody)
	ts.Equal(200, status)

	waitingSession, _ := ts.parseSessionResponse(ts.flowServer.assetCache, body)
	ts.Equal(flows.SessionStatus("waiting"), waitingSession.Status())

	// try to resume this session with a structurally invalid version of the flow passed as an asset
	status, body = ts.testHTTPRequest("POST", "http://localhost:8800/flow/resume", ts.buildResumeRequest(testStructurallyInvalidFlowAssets, waitingSession, []flows.Event{
		events.NewMsgReceivedEvent(msg),
	}))
	ts.Equal(400, status)
	ts.assertErrorResponse(body, []string{"unable to read asset[url=http://testserver/assets/flow/76f0a02f-3b75-4b86-9064-e9195e1b3a02]: destination 714f1409-486e-4e8e-bb08-23e2943ef9f6 of exit[uuid=37d8813f-1402-4ad2-9cc2-e9054a96525b] isn't a known node"})

	// try to resume this session with a invalid version of the flow which is missing a group
	status, body = ts.testHTTPRequest("POST", "http://localhost:8800/flow/resume", ts.buildResumeRequest(testFlowMissingGroupAssets, waitingSession, []flows.Event{
		events.NewMsgReceivedEvent(msg),
	}))
	ts.Equal(400, status)
	ts.assertErrorResponse(body, []string{"validation failed for flow[uuid=76f0a02f-3b75-4b86-9064-e9195e1b3a02]: validation failed for action[uuid=ad154980-7bf7-4ab8-8728-545fd6378912, type=add_to_group]: no such group with uuid '77a1bb5c-92f7-42bc-8a54-d21c1536ebc0'"})

	// check we can resume if we include a fixed version of the flow as an asset
	status, body = ts.testHTTPRequest("POST", "http://localhost:8800/flow/resume", ts.buildResumeRequest(testValidFlowWithWaitAssets, waitingSession, []flows.Event{
		events.NewMsgReceivedEvent(msg),
	}))
	ts.Equal(200, status)

	// check we got back a completed session
	completedSession, _ := ts.parseSessionResponse(ts.flowServer.assetCache, body)
	ts.Equal(flows.SessionStatus("completed"), completedSession.Status())

	// mess with our waiting session JSON so we appear to be waiting on a node that doesn't exist in the flow
	sessionJSON := ts.buildResumeRequest(`[]`, waitingSession, []flows.Event{
		events.NewMsgReceivedEvent(msg),
	})
	sessionJSON = strings.Replace(sessionJSON, "a58be63b-907d-4a1a-856b-0bb5579d7507", "626daa56-2fac-48eb-825d-af9a7ab23a2c", -1)

	// and try to resume that
	status, body = ts.testHTTPRequest("POST", "http://localhost:8800/flow/resume", sessionJSON)
	ts.Equal(200, status)

	// check we got back an errored session
	erroredSession, _ := ts.parseSessionResponse(ts.flowServer.assetCache, body)
	ts.Equal(flows.SessionStatus("errored"), erroredSession.Status())
}

func TestServerSuite(t *testing.T) {
	suite.Run(t, new(ServerTestSuite))
}
