package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/engine"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/test"
	"github.com/nyaruka/goflow/utils"

	"github.com/stretchr/testify/suite"
)

var testStructurallyInvalidFlowAssets = `[
	{
		"type": "flow",
		"url": "http://testserver/assets/flow/76f0a02f-3b75-4b86-9064-e9195e1b3a02",
		"content": {
			"uuid": "76f0a02f-3b75-4b86-9064-e9195e1b3a02",
			"name": "Test Flow",
			"language": "eng",
			"type": "messaging",
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
	},
	{
		"type": "channel",
		"url": "http://testserver/assets/channel",
		"content": []
	},
	{
		"type": "group",
		"url": "http://testserver/assets/group",
		"content": []
	},
	{
		"type": "label",
		"url": "http://testserver/assets/label",
		"content": []
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
			"type": "messaging",
			"nodes": [
				{
					"uuid": "a58be63b-907d-4a1a-856b-0bb5579d7507",
					"actions": [
						{
							"uuid": "ad154980-7bf7-4ab8-8728-545fd6378912",
							"type": "add_contact_groups",
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
		"type": "channel",
		"url": "http://testserver/assets/channel",
		"content": []
	},
	{
		"type": "group",
		"url": "http://testserver/assets/group",
		"content": [
			{
				"uuid": "2aad21f6-30b7-42c5-bd7f-1b720c154817",
				"name": "Survey Audience"
			}
		]
	},
	{
		"type": "label",
		"url": "http://testserver/assets/label",
		"content": []
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
			"type": "messaging",
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
	},
	{
		"type": "channel",
		"url": "http://testserver/assets/channel",
		"content": []
	},
	{
		"type": "group",
		"url": "http://testserver/assets/group",
		"content": []
	},
	{
		"type": "label",
		"url": "http://testserver/assets/label",
		"content": []
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
			"type": "messaging",
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
	},
	{
		"type": "channel",
		"url": "http://testserver/assets/channel",
		"content": []
	},
	{
		"type": "group",
		"url": "http://testserver/assets/group",
		"content": []
	},
	{
		"type": "label",
		"url": "http://testserver/assets/label",
		"content": []
	}
]`

var testValidFlowWithWebhook = `[
	{
		"type": "flow",
		"url": "http://testserver/assets/flow/76f0a02f-3b75-4b86-9064-e9195e1b3a02",
		"content": {
			"uuid": "76f0a02f-3b75-4b86-9064-e9195e1b3a02",
			"name": "Webhook Flow",
			"language": "eng",
			"type": "messaging",
			"nodes": [
				{
					"uuid": "a58be63b-907d-4a1a-856b-0bb5579d7507",
					"actions": [
						{
                            "uuid": "06153fbd-3e2c-413a-b0df-ed15d631835a",
                            "type": "call_webhook",
                            "method": "GET",
                            "url": "http://localhost:49993/?cmd=success"
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
		"type": "channel",
		"url": "http://testserver/assets/channel",
		"content": []
	},
	{
		"type": "group",
		"url": "http://testserver/assets/group",
		"content": []
	},
	{
		"type": "label",
		"url": "http://testserver/assets/label",
		"content": []
	}
]`

var assetServerConfig = `{
	"type_urls": {
		"channel": "http://testserver/assets/channel/",
		"flow": "http://testserver/assets/flow/",
		"field": "http://testserver/assets/field/",
		"group": "http://testserver/assets/group/",
		"label": "http://testserver/assets/label/"
	}
}`

var startRequestTemplate = `{
	"assets": %s,
	"asset_server": %s,
	"trigger": {
		"type": "manual",
		"flow": {"uuid": "76f0a02f-3b75-4b86-9064-e9195e1b3a02", "name": "Test Flow"},
		"triggered_on": "2000-01-01T00:00:00.000000000-00:00"
	},
	"config": %s
}`

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

	ts.httpServer, _ = test.NewTestHTTPServer(49993)
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

func (ts *ServerTestSuite) assertExpressionResponse(body []byte, expectedResult string, expectedErrors []string) {
	expResp := &expressionResponse{}
	err := json.Unmarshal(body, &expResp)
	ts.Require().NoError(err)
	ts.Equal(expectedResult, expResp.Result)
	ts.Equal(expectedErrors, expResp.Errors)
}

func (ts *ServerTestSuite) parseSessionResponse(body []byte) (flows.Session, []map[string]interface{}) {
	envelope := struct {
		Session json.RawMessage
		Log     []map[string]interface{}
	}{}
	err := json.Unmarshal(body, &envelope)
	ts.Require().NoError(err)

	assets, err := engine.NewSessionAssets(engine.NewMockServerSource(ts.flowServer.assetCache))
	ts.Require().NoError(err)

	session, err := engine.ReadSession(assets, engine.NewDefaultConfig(), test.TestHTTPClient, envelope.Session)
	ts.Require().NoError(err)

	return session, envelope.Log
}

func (ts *ServerTestSuite) buildResumeRequest(assetsJSON string, session flows.Session, events []flows.Event) string {
	sessionJSON, err := utils.JSONMarshal(session)
	ts.Require().NoError(err)

	eventEnvelopes := make([]*utils.TypedEnvelope, len(events))
	for e := range events {
		eventEnvelopes[e], err = utils.EnvelopeFromTyped(events[e])
		ts.Require().NoError(err)
	}

	assetsData := json.RawMessage(assetsJSON)
	assetServer, _ := utils.JSONMarshal(engine.NewMockServerSource(nil))

	request := &resumeRequest{
		sessionRequest: sessionRequest{
			Assets:      &assetsData,
			AssetServer: assetServer,
		},
		Session: sessionJSON,
		Events:  eventEnvelopes,
	}

	requestJSON, err := utils.JSONMarshal(request)
	ts.Require().NoError(err)
	return string(requestJSON)
}

func (ts *ServerTestSuite) TestHomePages() {
	// hit our home page
	status, body := ts.testHTTPRequest("GET", "http://localhost:8800/", "")
	ts.Equal(200, status)
	ts.Contains(string(body), "Start")
	ts.Contains(string(body), "Resume")

	// test an example start request on the home page
	startJSON, err := ioutil.ReadFile("testdata/start.json")
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
	ts.assertExpressionResponse(body, "", []string{"error evaluating @(foo + 2): json object has no property 'foo'"})
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
	status, body = ts.testHTTPRequest("POST", "http://localhost:8800/flow/start", fmt.Sprintf(`{"assets": %s, "asset_server": %s, "trigger": {"type": "manual"}}`, testValidFlowWithNoWaitAssets, assetServerConfig))
	ts.Equal(400, status)
	ts.assertErrorResponse(body, []string{"unable to read trigger[type=manual]: field 'flow' is required, field 'triggered_on' is required"})

	// try POSTing to the start endpoint a structurally invalid flow asset
	requestBody := fmt.Sprintf(startRequestTemplate, testStructurallyInvalidFlowAssets, assetServerConfig, `{}`)
	status, body = ts.testHTTPRequest("POST", "http://localhost:8800/flow/start", requestBody)
	ts.Equal(400, status)
	ts.assertErrorResponse(body, []string{"unable to read asset[url=http://testserver/assets/flow/76f0a02f-3b75-4b86-9064-e9195e1b3a02]: destination 714f1409-486e-4e8e-bb08-23e2943ef9f6 of exit[uuid=37d8813f-1402-4ad2-9cc2-e9054a96525b] isn't a known node"})

	// try POSTing to the start endpoint a flow asset that references a non-existent group asset
	requestBody = fmt.Sprintf(startRequestTemplate, testFlowMissingGroupAssets, assetServerConfig, `{}`)
	status, body = ts.testHTTPRequest("POST", "http://localhost:8800/flow/start", requestBody)
	ts.Equal(400, status)
	ts.assertErrorResponse(body, []string{"validation failed for flow[uuid=76f0a02f-3b75-4b86-9064-e9195e1b3a02]: validation failed for action[uuid=ad154980-7bf7-4ab8-8728-545fd6378912, type=add_contact_groups]: no such group with uuid '77a1bb5c-92f7-42bc-8a54-d21c1536ebc0'"})

	// POST to the start endpoint with a valid flow with no wait (it should complete)
	requestBody = fmt.Sprintf(startRequestTemplate, testValidFlowWithNoWaitAssets, assetServerConfig, `{}`)
	status, body = ts.testHTTPRequest("POST", "http://localhost:8800/flow/start", requestBody)
	ts.Equal(200, status)

	session, _ := ts.parseSessionResponse(body)
	ts.Equal(flows.SessionStatus("completed"), session.Status())

	// try to resume this completed session but with no caller events
	status, body = ts.testHTTPRequest("POST", "http://localhost:8800/flow/resume", ts.buildResumeRequest(`[]`, session, []flows.Event{}))
	ts.Equal(400, status)
	ts.assertErrorResponse(body, []string{"field 'events' must have a minimum of 1 items"})

	// try to resume this completed session
	tgURN, _ := urns.NewTelegramURN(1234567, "bob")
	msg := flows.NewMsgIn(flows.MsgUUID(utils.NewUUID()), flows.NilMsgID, tgURN, nil, "hello", []flows.Attachment{})
	status, body = ts.testHTTPRequest("POST", "http://localhost:8800/flow/resume", ts.buildResumeRequest(`[]`, session, []flows.Event{
		events.NewMsgReceivedEvent(msg),
	}))
	ts.Equal(400, status)
	ts.assertErrorResponse(body, []string{"only waiting sessions can be resumed"})

	// start another session on a flow that will wait for input
	requestBody = fmt.Sprintf(startRequestTemplate, testValidFlowWithWaitAssets, assetServerConfig, `{}`)
	status, body = ts.testHTTPRequest("POST", "http://localhost:8800/flow/start", requestBody)
	ts.Equal(200, status)

	waitingSession, _ := ts.parseSessionResponse(body)
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
	ts.assertErrorResponse(body, []string{"validation failed for flow[uuid=76f0a02f-3b75-4b86-9064-e9195e1b3a02]: validation failed for action[uuid=ad154980-7bf7-4ab8-8728-545fd6378912, type=add_contact_groups]: no such group with uuid '77a1bb5c-92f7-42bc-8a54-d21c1536ebc0'"})

	// check we can resume if we include a fixed version of the flow as an asset
	status, body = ts.testHTTPRequest("POST", "http://localhost:8800/flow/resume", ts.buildResumeRequest(testValidFlowWithWaitAssets, waitingSession, []flows.Event{
		events.NewMsgReceivedEvent(msg),
	}))
	ts.Equal(200, status)

	// check we got back a completed session
	completedSession, _ := ts.parseSessionResponse(body)
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
	erroredSession, _ := ts.parseSessionResponse(body)
	ts.Equal(flows.SessionStatus("errored"), erroredSession.Status())
}

func (ts *ServerTestSuite) TestWebhookMocking() {
	testCases := []struct {
		config         string
		expectedStatus int
		expectedBody   string
	}{
		// default config is make webhook calls
		{`{}`, 200, `{ "ok": "true" }`},

		// explicitly disabled or enabled
		{`{"disable_webhooks": false}`, 200, `{ "ok": "true" }`},
		{`{"disable_webhooks": true}`, 200, "DISABLED"},

		// a matching mock will always be used and matching is case-insensitive
		{`{"webhook_mocks":[{"method":"GET","url":"http://localhost:49993/?cmd=success","status":201,"body":"I'm mocked"}]}`, 201, "I'm mocked"},
		{`{"webhook_mocks":[{"method":"get","url":"http://LOCALHOST:49993/?cmd=success","status":201,"body":"I'm mocked"}]}`, 201, "I'm mocked"},

		// no matching mock means we fall back to whether disable_webhooks is set
		{`{"webhook_mocks":[{"method":"POST","url":"http://xxxxxx/?cmd=success","status":201,"body":"I'm mocked"}]}`, 200, `{ "ok": "true" }`},
		{`{"disable_webhooks": true, "webhook_mocks":[{"method":"POST","url":"http://xxxxxx/?cmd=success","status":201,"body":"I'm mocked"}]}`, 200, "DISABLED"},
	}

	for _, tc := range testCases {
		// POST to the start endpoint with a flow with a webhook call, with webhooks disabled
		requestBody := fmt.Sprintf(startRequestTemplate, testValidFlowWithWebhook, assetServerConfig, tc.config)
		status, body := ts.testHTTPRequest("POST", "http://localhost:8800/flow/start", requestBody)
		ts.Equal(200, status)

		session, _ := ts.parseSessionResponse(body)
		run := session.Runs()[0]

		ts.NotNil(run.Webhook())
		ts.Equal(tc.expectedStatus, run.Webhook().StatusCode())
		ts.Equal(tc.expectedBody, run.Webhook().Body())
	}
}

func TestServerSuite(t *testing.T) {
	suite.Run(t, new(ServerTestSuite))
}
