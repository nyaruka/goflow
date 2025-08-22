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
	test.MockUniverse()

	session, _, err := test.CreateTestSession("", envs.RedactionPolicyNone)
	require.NoError(t, err)

	tz, _ := time.LoadLocation("Africa/Kigali")
	timeout := 500
	gender := session.Assets().Fields().Get("gender")
	jotd := session.Assets().OptIns().Get("248be71d-78e9-4d71-a6c4-9981d369e5cb")
	weather := session.Assets().Topics().Get("472a7a73-96cb-4736-b567-056d987cc5b4")
	user := session.Assets().Users().Get("0c78ef47-7d56-44d8-8f57-96e0f30e8f44")
	facebook := session.Assets().Channels().Get("4bb288a0-7fca-4da1-abe8-59a593aff648")
	ticket := flows.NewTicket("7481888c-07dd-47dc-bf22-ef7448696ffe", weather, user)
	gpt4 := session.Assets().LLMs().Get("14115c03-b4c5-49e2-b9ac-390c43e9d7ce")
	call := flows.NewCall("0198ce92-ff2f-7b07-b158-b21ab168ebba", facebook, "tel:+12065551212")

	eventTests := []struct {
		event    func() flows.Event
		snapshot string
	}{
		{
			func() flows.Event {
				return events.NewAirtimeTransferred(
					&flows.AirtimeTransfer{
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
				)
			},
			`airtime_transferred`,
		},
		{
			func() flows.Event {
				return events.NewBroadcastCreated(
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
				)
			},
			`broadcast_created`,
		},
		{
			func() flows.Event {
				return events.NewCallCreated(call)
			},
			`call_created`,
		},
		{
			func() flows.Event {
				return events.NewCallReceived(call)
			},
			`call_received`,
		},
		{
			func() flows.Event {
				return events.NewClassifierCalled(
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
				)
			},
			`service_called`,
		},
		{
			func() flows.Event {
				return events.NewContactFieldChanged(
					gender,
					flows.NewValue(types.NewXText("male"), nil, nil, "", "", ""),
				)
			},
			`contact_field_changed`,
		},
		{
			func() flows.Event {
				return events.NewContactFieldChanged(
					gender,
					nil, // value being cleared
				)
			},
			`contact_field_changed_clear`,
		},
		{
			func() flows.Event {
				return events.NewContactGroupsChanged(
					[]*flows.Group{session.Assets().Groups().FindByName("Customers")},
					nil,
				)
			},
			`contact_groups_changed`,
		},
		{
			func() flows.Event {
				return events.NewContactStatusChanged(flows.ContactStatusActive)
			},
			`contact_status_changed_active`,
		},
		{
			func() flows.Event {
				return events.NewContactStatusChanged(flows.ContactStatusBlocked)
			},
			`contact_status_changed_blocked`,
		},
		{
			func() flows.Event {
				return events.NewContactStatusChanged(flows.ContactStatusStopped)
			},
			`contact_status_changed_stopped`,
		},
		{
			func() flows.Event {
				return events.NewContactLanguageChanged(i18n.Language("fra"))
			},
			`contact_language_changed`,
		},
		{
			func() flows.Event {
				return events.NewContactNameChanged("Bryan")
			},
			`contact_name_changed`,
		},
		{
			func() flows.Event {
				return events.NewContactTimezoneChanged(tz)
			},
			`contact_timezone_changed`,
		},
		{
			func() flows.Event {
				return events.NewContactURNsChanged([]urns.URN{
					urns.URN("tel:+12345678900"),
					urns.URN("twitterid:8764843252522#bob"),
				})
			},
			`contact_urns_changed`,
		},
		{
			func() flows.Event {
				return events.NewEmailSent([]string{"bob@nyaruka.com", "jim@nyaruka.com"}, "Update", "Flows are great!")
			},
			`email_sent`,
		},
		{
			func() flows.Event {
				return events.NewError("I'm an error")
			},
			`error`,
		},
		{
			func() flows.Event {
				return events.NewDependencyError(assets.NewFieldReference("age", "Age"))
			},
			`error_dependency`,
		},
		{
			func() flows.Event {
				return events.NewFailure(errors.New("503 is an failure"))
			},
			`failure`,
		},
		{
			func() flows.Event {
				return events.NewIVRCreated(
					flows.NewIVRMsgOut(
						urns.URN("tel:+12345678900"),
						assets.NewChannelReference(assets.ChannelUUID("57f1078f-88aa-46f4-a59a-948a5739c03d"), "My Android Phone"),
						"Hi there",
						"http://example.com/hi.mp3",
						"eng",
					),
				)
			},
			`ivr_created`,
		},
		{
			func() flows.Event {
				return events.NewLLMCalled(
					gpt4,
					"Categorize the following text as Positive or Negative",
					"Please stop messaging me",
					&flows.LLMResponse{Output: "Positive", TokensUsed: 567},
					123*time.Millisecond,
				)
			},
			`llm_called`,
		},
		{
			func() flows.Event {
				return events.NewMsgCreated(
					flows.NewMsgOut(
						urns.URN("tel:+12345678900"),
						assets.NewChannelReference(assets.ChannelUUID("57f1078f-88aa-46f4-a59a-948a5739c03d"), "My Android Phone"),
						&flows.MsgContent{Text: "Hi there"},
						nil,
						i18n.NilLocale,
						flows.NilUnsendableReason,
					),
				)
			},
			`msg_created`,
		},
		{
			func() flows.Event {
				return events.NewMsgCreated(
					flows.NewMsgOut(
						urns.URN("tel:+12345678900"),
						assets.NewChannelReference(assets.ChannelUUID("57f1078f-88aa-46f4-a59a-948a5739c03d"), "My Android Phone"),
						&flows.MsgContent{
							Text:         "Hi there",
							Attachments:  []utils.Attachment{"image/jpeg:http://s3.amazon.com/bucket/test.jpg"},
							QuickReplies: []flows.QuickReply{{Text: "yes"}, {Text: "no"}},
						},
						nil,
						"eng-US",
						flows.UnsendableReasonContactStatus,
					),
				)
			},
			`msg_created_rich`,
		},
		{
			func() flows.Event {
				return events.NewMsgWait(&timeout, time.Date(2022, 2, 3, 13, 45, 30, 0, time.UTC), hints.NewImage())
			},
			`msg_wait`,
		},
		{
			func() flows.Event {
				return events.NewWaitTimedOut()
			},
			`wait_timed_out`,
		},
		{
			func() flows.Event {
				return events.NewDialEnded(flows.NewDial(flows.DialStatusBusy, 0))
			},
			`dial_ended`,
		},
		{
			func() flows.Event {
				return events.NewDialWait(urns.URN("tel:+1234567890"), 20, 120, time.Date(2022, 2, 3, 13, 45, 30, 0, time.UTC))
			},
			`dial_wait`,
		},
		{
			func() flows.Event {
				return events.NewOptInRequested(jotd, facebook.Reference(), urns.URN("facebook:1234567890"))
			},
			`optin_requested`,
		},
		{
			func() flows.Event {
				return events.NewOptInStarted(jotd, facebook.Reference())
			},
			`optin_started`,
		},
		{
			func() flows.Event {
				return events.NewOptInStopped(jotd, facebook.Reference())
			},
			`optin_stopped`,
		},
		{
			func() flows.Event {
				return events.NewRunStarted(session.Runs()[0], true)
			},
			`run_started`,
		},
		{
			func() flows.Event {
				return events.NewSessionTriggered(
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
				)
			},
			`session_triggered`,
		},
		{
			func() flows.Event {
				return events.NewTicketClosed(ticket)
			},
			`ticket_closed`,
		},
		{
			func() flows.Event {
				return events.NewTicketOpened(ticket, "this is weird")
			},
			`ticket_opened`,
		},
	}

	for _, tc := range eventTests {
		test.MockUniverse()

		eventJSON, err := jsonx.MarshalPretty(tc.event())
		assert.NoError(t, err)

		test.AssertSnapshot(t, tc.snapshot, string(eventJSON))

		// try to read event back
		_, err = events.Read(eventJSON)
		assert.NoError(t, err)
	}
}

func TestReadEvent(t *testing.T) {
	// error if no type field
	_, err := events.Read([]byte(`{"foo": "bar"}`))
	assert.EqualError(t, err, "field 'type' is required")

	// error if we don't recognize action type
	_, err = events.Read([]byte(`{"type": "do_the_foo", "foo": "bar"}`))
	assert.EqualError(t, err, "unknown type: 'do_the_foo'")

	// valid existing type
	event, err := events.Read([]byte(`{"uuid": "0197b335-6ded-79a4-95a6-3af85b57f108", "type": "contact_name_changed", "created_on": "2006-01-02T15:04:05Z", "name": "Bob Smith"}`))
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
}
