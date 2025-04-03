package events_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"testing"
	"time"
	"unicode/utf8"

	"github.com/nyaruka/gocommon/dates"
	"github.com/nyaruka/gocommon/httpx"
	"github.com/nyaruka/gocommon/i18n"
	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/flows/routers/waits/hints"
	"github.com/nyaruka/goflow/services/webhooks"
	"github.com/nyaruka/goflow/test"
	"github.com/nyaruka/goflow/utils"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEventMarshaling(t *testing.T) {
	defer dates.SetNowFunc(time.Now)
	defer uuids.SetGenerator(uuids.DefaultGenerator)

	dates.SetNowFunc(dates.NewFixedNow(time.Date(2018, 10, 18, 14, 20, 30, 123456, time.UTC)))
	uuids.SetGenerator(uuids.NewSeededGenerator(12345, time.Now))

	session, _, err := test.CreateTestSession("", envs.RedactionPolicyNone)
	require.NoError(t, err)

	tz, _ := time.LoadLocation("Africa/Kigali")
	timeout := 500
	gender := session.Assets().Fields().Get("gender")
	jotd := session.Assets().OptIns().Get("248be71d-78e9-4d71-a6c4-9981d369e5cb")
	weather := session.Assets().Topics().Get("472a7a73-96cb-4736-b567-056d987cc5b4")
	user := session.Assets().Users().Get("bob@nyaruka.com")
	facebook := session.Assets().Channels().Get("4bb288a0-7fca-4da1-abe8-59a593aff648")
	ticket := flows.NewTicket("7481888c-07dd-47dc-bf22-ef7448696ffe", weather, user)
	gpt4 := session.Assets().LLMs().Get("14115c03-b4c5-49e2-b9ac-390c43e9d7ce")

	eventTests := []struct {
		event     flows.Event
		marshaled string
	}{
		{
			events.NewAirtimeTransferred(
				&flows.AirtimeTransfer{
					UUID:       "4c2d9b7a-e02c-4e6a-ab18-06df4cb5666d",
					ExternalID: "98765432",
					Sender:     urns.URN("tel:+593979099111"),
					Recipient:  urns.URN("tel:+593979099222"),
					Currency:   "USD",
					Amount:     decimal.RequireFromString("1.00"),
				},
				[]*flows.HTTPLog{
					{
						HTTPLogWithoutTime: &flows.HTTPLogWithoutTime{
							LogWithoutTime: &httpx.LogWithoutTime{
								URL:        "https://send.money.com/topup",
								StatusCode: 200,
								Request:    "POST /topup HTTP/1.1\r\nHost: send.money.com\r\nUser-Agent: Go-http-client/1.1\r\nAccept-Encoding: gzip\r\n\r\n",
								Response:   "HTTP/1.0 200 OK\r\nContent-Length: 14\r\n\r\n{\"errors\":[]}",
								ElapsedMS:  12,
							},
							Status: flows.CallStatusSuccess,
						},
						CreatedOn: dates.Now(),
					},
				},
			),
			`{
				"amount": 1,
        	    "created_on": "2018-10-18T14:20:30.000123456Z",
        	    "currency": "USD",
				"external_id": "98765432",
				"http_logs": [
					{
						"url": "https://send.money.com/topup",
						"status_code": 200,
						"status": "success",
						"request": "POST /topup HTTP/1.1\r\nHost: send.money.com\r\nUser-Agent: Go-http-client/1.1\r\nAccept-Encoding: gzip\r\n\r\n",
						"response": "HTTP/1.0 200 OK\r\nContent-Length: 14\r\n\r\n{\"errors\":[]}",
						"elapsed_ms": 12,
						"retries": 0,
						"created_on": "2018-10-18T14:20:30.000123456Z"
					}
				],
				"recipient": "tel:+593979099222",
        	    "sender": "tel:+593979099111",
				"type": "airtime_transferred",
				"transfer_uuid": "4c2d9b7a-e02c-4e6a-ab18-06df4cb5666d"
			}`,
		},
		{
			events.NewBroadcastCreated(
				flows.BroadcastTranslations{
					"eng": {Text: "Hello", Attachments: nil, QuickReplies: nil},
					"spa": {Text: "Hola", Attachments: nil, QuickReplies: nil},
				},
				i18n.Language("eng"),
				[]*assets.GroupReference{
					assets.NewGroupReference(assets.GroupUUID("5f9fd4f7-4b0f-462a-a598-18bfc7810412"), "Supervisors"),
				},
				[]*flows.ContactReference{
					flows.NewContactReference(flows.ContactUUID("b2aaf598-1bb3-4c7d-b6bb-1f8dbe2ac16f"), "Jim"),
				},
				"name = \"Bob\"",
				[]urns.URN{urns.URN("tel:+12345678900")},
			),
			`{
				"type": "broadcast_created",
				"created_on": "2018-10-18T14:20:30.000123456Z",
				"base_language": "eng",
				"translations": {
					"eng": {
						"text": "Hello"
					},
					"spa": {
						"text": "Hola"
					}
				},
				"groups": [
					{
						"name": "Supervisors",
						"uuid": "5f9fd4f7-4b0f-462a-a598-18bfc7810412"
					}
				],
				"contacts": [
					{
						"name": "Jim",
						"uuid": "b2aaf598-1bb3-4c7d-b6bb-1f8dbe2ac16f"
					}
				],
				"contact_query": "name = \"Bob\"",
				"urns": [
					"tel:+12345678900"
				]
			}`,
		},
		{
			events.NewClassifierCalled(
				assets.NewClassifierReference(assets.ClassifierUUID("4b937f49-7fb7-43a5-8e57-14e2f028a471"), "Booking"),
				[]*flows.HTTPLog{
					{
						HTTPLogWithoutTime: &flows.HTTPLogWithoutTime{
							LogWithoutTime: &httpx.LogWithoutTime{
								URL:        "https://api.wit.ai/message?v=20200513&q=hello",
								StatusCode: 200,
								Request:    "GET /message?v=20200513&q=hello HTTP/1.1\r\nHost: api.wit.ai\r\nUser-Agent: Go-http-client/1.1\r\nAccept-Encoding: gzip\r\n\r\n",
								Response:   "HTTP/1.0 200 OK\r\nContent-Length: 14\r\n\r\n{\"intents\":[]}",
								ElapsedMS:  12,
							},
							Status: flows.CallStatusSuccess,
						},
						CreatedOn: dates.Now(),
					},
				},
			),
			`{
				"type": "service_called",
				"created_on": "2018-10-18T14:20:30.000123456Z",
				"service": "classifier",
				"classifier": {
					"uuid": "4b937f49-7fb7-43a5-8e57-14e2f028a471",
					"name": "Booking"
				},
				"http_logs": [
					{
						"url": "https://api.wit.ai/message?v=20200513&q=hello",
						"status_code": 200,
						"status": "success",
						"request": "GET /message?v=20200513&q=hello HTTP/1.1\r\nHost: api.wit.ai\r\nUser-Agent: Go-http-client/1.1\r\nAccept-Encoding: gzip\r\n\r\n",
						"response": "HTTP/1.0 200 OK\r\nContent-Length: 14\r\n\r\n{\"intents\":[]}",
						"elapsed_ms": 12,
						"retries": 0,
						"created_on": "2018-10-18T14:20:30.000123456Z"
					}
				]
			}`,
		},
		{
			events.NewContactFieldChanged(
				gender,
				flows.NewValue(types.NewXText("male"), nil, nil, "", "", ""),
			),
			`{
				"created_on": "2018-10-18T14:20:30.000123456Z",
				"field": {
					"key": "gender",
					"name": "Gender"
				},
				"type": "contact_field_changed",
				"value": {
					"text": "male"
				}
			}`,
		},
		{
			events.NewContactFieldChanged(
				gender,
				nil, // value being cleared
			),
			`{
				"created_on": "2018-10-18T14:20:30.000123456Z",
				"field": {
					"key": "gender",
					"name": "Gender"
				},
				"type": "contact_field_changed",
				"value": null
			}`,
		},
		{
			events.NewContactGroupsChanged(
				[]*flows.Group{session.Assets().Groups().FindByName("Customers")},
				nil,
			),
			`{
				"created_on": "2018-10-18T14:20:30.000123456Z",
				"groups_added": [
					{
						"name": "Customers",
						"uuid": "1e1ce1e1-9288-4504-869e-022d1003c72a"
					}
				],
				"type": "contact_groups_changed"
			}`,
		},
		{
			events.NewContactStatusChanged(flows.ContactStatusActive),
			`{
				"created_on": "2018-10-18T14:20:30.000123456Z",
				"type": "contact_status_changed",
				"status": "active"
			}`,
		},
		{
			events.NewContactStatusChanged(flows.ContactStatusBlocked),
			`{
				"created_on": "2018-10-18T14:20:30.000123456Z",
				"type": "contact_status_changed",
				"status": "blocked"
			}`,
		},
		{
			events.NewContactStatusChanged(flows.ContactStatusStopped),
			`{
				"created_on": "2018-10-18T14:20:30.000123456Z",
				"type": "contact_status_changed",
				"status": "stopped"
			}`,
		},
		{
			events.NewContactLanguageChanged(i18n.Language("fra")),
			`{
				"created_on": "2018-10-18T14:20:30.000123456Z",
				"language": "fra",
				"type": "contact_language_changed"
			}`,
		},
		{
			events.NewContactRefreshed(session.Contact()),
			`{
				"contact": {
					"created_on": "2018-06-20T11:40:30.123456789Z",
					"fields": {
						"activation_token": {
							"text": "AACC55"
						},
						"age": {
							"number": 23,
							"text": "23"
						},
						"gender": {
							"text": "Male"
						},
						"join_date": {
							"datetime": "2017-12-02T00:00:00.000000-02:00",
							"text": "2017-12-02"
						}
					},
					"groups": [
						{
							"name": "Testers",
							"uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d"
						},
						{
							"name": "Males",
							"uuid": "4f1f98fc-27a7-4a69-bbdb-24744ba739a9"
						}
					],
					"id": 1234567,
					"language": "eng",
					"last_seen_on": "2017-12-31T11:35:10.035757258-02:00",
					"name": "Ryan Lewis",
					"status": "active",
					"ticket": {
						"uuid": "78d1fe0d-7e39-461e-81c3-a6a25f15ed69",
						"topic": {
							"uuid": "472a7a73-96cb-4736-b567-056d987cc5b4",
							"name": "Weather"
						},
						"assignee": {
							"email": "bob@nyaruka.com",
							"name": "Bob"
						}
					},
					"timezone": "America/Guayaquil",
					"urns": [
						"tel:+12024561111?channel=57f1078f-88aa-46f4-a59a-948a5739c03d",
						"twitterid:54784326227#nyaruka",
						"mailto:foo@bar.com"
					],
					"uuid": "5d76d86b-3bb9-4d5a-b822-c9d86f5d8e4f"
				},
				"created_on": "2018-10-18T14:20:30.000123456Z",
				"type": "contact_refreshed"
			}`,
		},
		{
			events.NewContactNameChanged("Bryan"),
			`{
				"created_on": "2018-10-18T14:20:30.000123456Z",
				"name": "Bryan",
				"type": "contact_name_changed"
			}`,
		},
		{
			events.NewContactTimezoneChanged(tz),
			`{
				"created_on": "2018-10-18T14:20:30.000123456Z",
				"timezone": "Africa/Kigali",
				"type": "contact_timezone_changed"
			}`,
		},
		{
			events.NewContactURNsChanged([]urns.URN{
				urns.URN("tel:+12345678900"),
				urns.URN("twitterid:8764843252522#bob"),
			}),
			`{
				"created_on": "2018-10-18T14:20:30.000123456Z",
				"type": "contact_urns_changed",
				"urns": [
					"tel:+12345678900",
					"twitterid:8764843252522#bob"
				]
			}`,
		},
		{
			events.NewEmailSent([]string{"bob@nyaruka.com", "jim@nyaruka.com"}, "Update", "Flows are great!"),
			`{
				"created_on": "2018-10-18T14:20:30.000123456Z",
				"type": "email_sent",
				"to": ["bob@nyaruka.com", "jim@nyaruka.com"],
				"subject": "Update",
				"body": "Flows are great!"
			}`,
		},
		{
			events.NewEnvironmentRefreshed(session.Environment()),
			`{
				"created_on": "2018-10-18T14:20:30.000123456Z",
				"environment": {
					"allowed_languages": [
						"eng",
						"spa"
					],
					"date_format": "DD-MM-YYYY",
					"default_country": "US",
					"input_collation": "default",
					"number_format": {
						"decimal_symbol": ".",
						"digit_grouping_symbol": ","
					},
					"redaction_policy": "none",
					"time_format": "tt:mm",
					"timezone": "America/Guayaquil"
				},
				"type": "environment_refreshed"
			}`,
		},
		{
			events.NewError("I'm an error"),
			`{
				"created_on": "2018-10-18T14:20:30.000123456Z",
				"text": "I'm an error",
				"type": "error"
			}`,
		},
		{
			events.NewFailure(errors.New("503 is an failure")),
			`{
				"created_on": "2018-10-18T14:20:30.000123456Z",
				"text": "503 is an failure",
				"type": "failure"
			}`,
		},
		{
			events.NewDependencyError(assets.NewFieldReference("age", "Age")),
			`{
				"created_on": "2018-10-18T14:20:30.000123456Z",
				"text": "missing dependency: field[key=age,name=Age]",
				"type": "error"
			}`,
		},
		{
			events.NewIVRCreated(
				flows.NewIVRMsgOut(
					urns.URN("tel:+12345678900"),
					assets.NewChannelReference(assets.ChannelUUID("57f1078f-88aa-46f4-a59a-948a5739c03d"), "My Android Phone"),
					"Hi there",
					"http://example.com/hi.mp3",
					"eng",
				),
			),
			`{
				"type": "ivr_created",
				"created_on": "2018-10-18T14:20:30.000123456Z",
				"msg": {
					"uuid": "94f0e964-be11-4d7b-866b-323926b4c6a0",
					"urn": "tel:+12345678900",
					"channel": {
						"name": "My Android Phone",
						"uuid": "57f1078f-88aa-46f4-a59a-948a5739c03d"
					},
					"text": "Hi there",
					"attachments": ["audio:http://example.com/hi.mp3"],
					"locale": "eng"
				}
			}`,
		},
		{
			events.NewLLMCalled(
				gpt4,
				"Categorize the following text as Positive or Negative",
				"Please stop messaging me",
				&flows.LLMResponse{
					Output:     "Positive",
					TokensUsed: 567,
				},
				123*time.Millisecond,
			),
			`{
				"type": "llm_called",
				"created_on": "2018-10-18T14:20:30.000123456Z",
				"llm": {
					"uuid": "14115c03-b4c5-49e2-b9ac-390c43e9d7ce", 
					"name": "GPT-4"
				},
				"instructions": "Categorize the following text as Positive or Negative",
				"input": "Please stop messaging me",
				"output": "Positive",
				"tokens_used": 567,
				"elapsed_ms": 123
			}`,
		},
		{
			events.NewMsgCreated(
				flows.NewMsgOut(
					urns.URN("tel:+12345678900"),
					assets.NewChannelReference(assets.ChannelUUID("57f1078f-88aa-46f4-a59a-948a5739c03d"), "My Android Phone"),
					&flows.MsgContent{Text: "Hi there"},
					nil,
					flows.NilMsgTopic,
					i18n.NilLocale,
					flows.NilUnsendableReason,
				),
			),
			`{
				"type": "msg_created",
				"created_on": "2018-10-18T14:20:30.000123456Z",
				"msg": {
					"uuid": "5b835baa-3607-48cb-a489-7cc248dc15c5",
					"urn": "tel:+12345678900",
					"channel": {
						"name": "My Android Phone",
						"uuid": "57f1078f-88aa-46f4-a59a-948a5739c03d"
					},
					"text": "Hi there"
				}
			}`,
		},
		{
			events.NewMsgCreated(
				flows.NewMsgOut(
					urns.URN("tel:+12345678900"),
					assets.NewChannelReference(assets.ChannelUUID("57f1078f-88aa-46f4-a59a-948a5739c03d"), "My Android Phone"),
					&flows.MsgContent{
						Text:         "Hi there",
						Attachments:  []utils.Attachment{"image/jpeg:http://s3.amazon.com/bucket/test.jpg"},
						QuickReplies: []flows.QuickReply{{Text: "yes"}, {Text: "no"}},
					},
					nil,
					flows.MsgTopicAgent,
					"eng-US",
					flows.UnsendableReasonContactStatus,
				),
			),
			`{
				"type": "msg_created",
				"created_on": "2018-10-18T14:20:30.000123456Z",
				"msg": {
					"uuid": "5078c828-5e46-4bac-8c96-e8696b9ca2d2",
					"urn": "tel:+12345678900",
					"channel": {
						"name": "My Android Phone",
						"uuid": "57f1078f-88aa-46f4-a59a-948a5739c03d"
					},
					"text": "Hi there",
					"attachments": ["image/jpeg:http://s3.amazon.com/bucket/test.jpg"],
					"quick_replies": [{"text": "yes"}, {"text": "no"}],
					"topic": "agent",
					"locale": "eng-US",
					"unsendable_reason": "contact_status"
				}
			}`,
		},
		{
			events.NewMsgWait(&timeout, time.Date(2022, 2, 3, 13, 45, 30, 0, time.UTC), hints.NewImageHint()),
			`{
				"type": "msg_wait",
				"created_on": "2018-10-18T14:20:30.000123456Z",
				"timeout_seconds": 500,
				"expires_on": "2022-02-03T13:45:30Z",
				"hint": {"type": "image"}
			}`,
		},
		{
			events.NewWaitTimedOut(),
			`{
				"type": "wait_timed_out",
				"created_on": "2018-10-18T14:20:30.000123456Z"
			}`,
		},
		{
			events.NewDialEnded(flows.NewDial(flows.DialStatusBusy, 0)),
			`{
				"type": "dial_ended",
				"created_on": "2018-10-18T14:20:30.000123456Z",
				"dial": {
					"status": "busy",
					"duration": 0
				}
			}`,
		},
		{
			events.NewDialWait(urns.URN("tel:+1234567890"), 20, 120, time.Date(2022, 2, 3, 13, 45, 30, 0, time.UTC)),
			`{
				"type": "dial_wait",
				"created_on": "2018-10-18T14:20:30.000123456Z",
				"urn": "tel:+1234567890",
				"dial_limit_seconds": 20,
				"call_limit_seconds": 120,
				"expires_on": "2022-02-03T13:45:30Z"
			}`,
		},
		{
			events.NewOptInRequested(jotd, facebook, urns.URN("facebook:1234567890")),
			`{
				"type": "optin_requested",
				"created_on": "2018-10-18T14:20:30.000123456Z",
				"optin": {
					"uuid": "248be71d-78e9-4d71-a6c4-9981d369e5cb",
					"name": "Joke Of The Day"
				},
				"channel": {
					"uuid": "4bb288a0-7fca-4da1-abe8-59a593aff648",
					"name": "Facebook Channel"
				},
				"urn": "facebook:1234567890"
			}`,
		},
		{
			events.NewSessionTriggered(
				assets.NewFlowReference(assets.FlowUUID("e4d441f0-24e3-4627-85fb-1e99e733baf0"), "Collect Age"),
				[]*assets.GroupReference{
					assets.NewGroupReference(assets.GroupUUID("5f9fd4f7-4b0f-462a-a598-18bfc7810412"), "Supervisors"),
				},
				[]*flows.ContactReference{
					flows.NewContactReference(flows.ContactUUID("b2aaf598-1bb3-4c7d-b6bb-1f8dbe2ac16f"), "Jim"),
				},
				"age > 20",
				events.Exclusions{InAFlow: true},
				false,
				[]urns.URN{urns.URN("tel:+12345678900")},
				json.RawMessage(`{"uuid": "779eaf3f-1c59-4374-a7cb-0eae9c5e8800"}`),
				&flows.SessionHistory{ParentUUID: "418a704c-f33e-4924-a00e-1763d1498a13", Ancestors: 2, AncestorsSinceInput: 0},
			),
			`{
				"type": "session_triggered",
				"created_on": "2018-10-18T14:20:30.000123456Z",
				"contacts": [
					{
						"name": "Jim",
						"uuid": "b2aaf598-1bb3-4c7d-b6bb-1f8dbe2ac16f"
					}
				],
				"flow": {
					"name": "Collect Age",
					"uuid": "e4d441f0-24e3-4627-85fb-1e99e733baf0"
				},
				"groups": [
					{
						"name": "Supervisors",
						"uuid": "5f9fd4f7-4b0f-462a-a598-18bfc7810412"
					}
				],
				"urns": [
					"tel:+12345678900"
				],
				"contact_query": "age > 20",
				"exclusions": {"in_a_flow": true},
				"run_summary": {
					"uuid": "779eaf3f-1c59-4374-a7cb-0eae9c5e8800"
				},
				"history": {
					"parent_uuid": "418a704c-f33e-4924-a00e-1763d1498a13",
					"ancestors": 2,
					"ancestors_since_input": 0
				}
			}`,
		},
		{
			events.NewTicketOpened(ticket, "this is weird"),
			`{
				"type": "ticket_opened",
				"created_on": "2018-10-18T14:20:30.000123456Z",
				"ticket": {
					"uuid": "7481888c-07dd-47dc-bf22-ef7448696ffe",
					"topic": {
						"uuid": "472a7a73-96cb-4736-b567-056d987cc5b4",
         				"name": "Weather"
					},
					"assignee": {
						"email": "bob@nyaruka.com",
						"name": "Bob"
					}
				},
				"note": "this is weird"
			}`,
		},
	}

	for _, tc := range eventTests {
		eventJSON, err := jsonx.Marshal(tc.event)
		assert.NoError(t, err)

		test.AssertEqualJSON(t, []byte(tc.marshaled), eventJSON, "event JSON mismatch")

		// try to read event back
		_, err = events.ReadEvent(eventJSON)
		assert.NoError(t, err)
	}
}

func TestReadEvent(t *testing.T) {
	// error if no type field
	_, err := events.ReadEvent([]byte(`{"foo": "bar"}`))
	assert.EqualError(t, err, "field 'type' is required")

	// error if we don't recognize action type
	_, err = events.ReadEvent([]byte(`{"type": "do_the_foo", "foo": "bar"}`))
	assert.EqualError(t, err, "unknown type: 'do_the_foo'")

	// valid existing type
	event, err := events.ReadEvent([]byte(`{"type": "contact_name_changed", "created_on": "2006-01-02T15:04:05Z", "name": "Bob Smith"}`))
	require.NoError(t, err)

	assert.Equal(t, events.TypeContactNameChanged, event.Type())
	eventTime, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
	assert.Equal(t, eventTime, event.CreatedOn())

}

func TestWebhookCalledEventTrimming(t *testing.T) {
	defer httpx.SetRequestor(httpx.DefaultRequestor)

	httpx.SetRequestor(httpx.NewMockRequestor(map[string][]*httpx.MockResponse{
		"http://temba.io/": {
			httpx.NewMockResponse(200, nil, bytes.Repeat([]byte("Y"), 20000)),
		},
	}))

	request, _ := http.NewRequest("GET", "http://temba.io/", strings.NewReader(strings.Repeat("X", 20000)))

	svc := webhooks.NewService(http.DefaultClient, nil, nil, nil, 1024*1024)
	call, err := svc.Call(request)
	require.NoError(t, err)

	assert.Equal(t, 42, len(call.ResponseTrace))
	assert.Equal(t, 20000, len(call.ResponseBody))

	event := events.NewWebhookCalled(call, flows.CallStatusSuccess, "")

	assert.Equal(t, "http://temba.io/", event.URL)
	assert.Equal(t, 10000, len(event.Request))
	assert.Equal(t, "XXXXXXX...", event.Request[9990:])
	assert.Equal(t, 10000, len(event.Response))
	assert.Equal(t, "YYYYYYY...", event.Response[9990:])
}

func TestWebhookCalledEventValid(t *testing.T) {
	defer httpx.SetRequestor(httpx.DefaultRequestor)

	httpx.SetRequestor(httpx.NewMockRequestor(map[string][]*httpx.MockResponse{
		"http://temba.io/": {
			httpx.NewMockResponse(200, map[string]string{"Header": "hello"}, []byte(`{"foo": "bar"}`)),
		},
	}))

	request, _ := http.NewRequest("GET", "http://temba.io/", nil)

	svc := webhooks.NewService(http.DefaultClient, nil, nil, nil, 1024*1024)
	call, err := svc.Call(request)
	require.NoError(t, err)

	event := events.NewWebhookCalled(call, flows.CallStatusSuccess, "")

	assert.Equal(t, "http://temba.io/", event.URL)
	assert.Equal(t, "HTTP/1.0 200 OK\r\nContent-Length: 14\r\nHeader: hello\r\n\r\n{\"foo\": \"bar\"}", event.Response)
	assert.True(t, utf8.ValidString(event.Response))
	assert.Equal(t, events.ExtractionValid, event.Extraction)
}

func TestWebhookCalledEventNullChar(t *testing.T) {
	defer httpx.SetRequestor(httpx.DefaultRequestor)

	httpx.SetRequestor(httpx.NewMockRequestor(map[string][]*httpx.MockResponse{
		"http://temba.io/": {
			httpx.NewMockResponse(200, nil, []byte("abc \x00 \\u0000 \\\u0000 \\\\u0000")),
		},
	}))

	request, _ := http.NewRequest("GET", "http://temba.io/", nil)

	svc := webhooks.NewService(http.DefaultClient, nil, nil, nil, 1024*1024)
	call, err := svc.Call(request)
	require.NoError(t, err)

	event := events.NewWebhookCalled(call, flows.CallStatusSuccess, "")

	// actual null will have been stripped, escaped null will remain
	assert.Equal(t, "http://temba.io/", event.URL)
	assert.Equal(t, "HTTP/1.0 200 OK\r\nContent-Length: 23\r\n\r\nabc � � \\� \\\\u0000", event.Response)
	assert.True(t, utf8.ValidString(event.Response))
	assert.Equal(t, events.ExtractionIgnored, event.Extraction)
}

func TestWebhookCalledEventBadUTF8(t *testing.T) {
	defer httpx.SetRequestor(httpx.DefaultRequestor)

	httpx.SetRequestor(httpx.NewMockRequestor(map[string][]*httpx.MockResponse{
		"http://temba.io/": {
			httpx.NewMockResponse(200, map[string]string{"Bad-Header": "\xa0\xa1"}, []byte("{\"foo\": \"\xa0\xa1\"}")),
		},
	}))

	request, _ := http.NewRequest("GET", "http://temba.io/", nil)

	svc := webhooks.NewService(http.DefaultClient, nil, nil, nil, 1024*1024)
	call, err := svc.Call(request)
	require.NoError(t, err)

	event := events.NewWebhookCalled(call, flows.CallStatusSuccess, "")

	assert.Equal(t, "http://temba.io/", event.URL)
	assert.Equal(t, "HTTP/1.0 200 OK\r\nContent-Length: 13\r\nBad-Header: �\r\n\r\n...", event.Response)
	assert.True(t, utf8.ValidString(event.Response))
	assert.Equal(t, events.ExtractionCleaned, event.Extraction)
}

func TestDeprecatedEvents(t *testing.T) {
	eventJSON := []byte(`{
		"type": "classifier_called",
		"created_on": "2006-01-02T15:04:05Z",
		"classifier": {"uuid": "1c06c884-39dd-4ce4-ad9f-9a01cbe6c000", "name": "Booking"},
		"http_logs": [
			{
				"url": "https://api.wit.ai/message?v=20170307&q=hello",
				"status_code": 200,
				"status": "success",
				"request": "GET /message?v=20170307&q=hello HTTP/1.1",
				"response": "HTTP/1.1 200 OK\r\n\r\n{\"intents\":[]}",
				"elapsed_ms": 123,
				"retries": 0,
				"created_on": "2006-01-02T15:04:05Z"
			}
		]
	}`)

	e, err := events.ReadEvent(eventJSON)
	assert.NoError(t, err)
	assert.Equal(t, events.TypeClassifierCalled, e.Type())

	marshaled, err := jsonx.Marshal(e)
	assert.NoError(t, err)
	test.AssertEqualJSON(t, eventJSON, marshaled)
}
