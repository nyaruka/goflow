package events_test

import (
	"bytes"
	"encoding/json"
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
	"github.com/nyaruka/goflow/core"
	"github.com/nyaruka/goflow/core/events"
	"github.com/nyaruka/goflow/core/hints"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
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
	ticket := core.NewTicket("7481888c-07dd-47dc-bf22-ef7448696ffe", core.TicketStatusOpen, weather, user)
	gpt4 := session.Assets().LLMs().Get("14115c03-b4c5-49e2-b9ac-390c43e9d7ce")
	call := core.NewCall("0198ce92-ff2f-7b07-b158-b21ab168ebba", facebook, "tel:+12065551212")

	eventTests := []struct {
		event    func() events.Event
		snapshot string
	}{
		{
			func() events.Event {
				return events.NewAirtimeCreated(
					events.NewEventUUID(),
					&core.AirtimeTransfer{
						ExternalID: "98765432",
						Sender:     urns.URN("tel:+593979099111"),
						Recipient:  urns.URN("tel:+593979099222"),
						Currency:   "USD",
						Amount:     decimal.RequireFromString("1.00"),
					},
					[]*core.HTTPLog{
						{
							HTTPLogWithoutTime: &core.HTTPLogWithoutTime{
								LogWithoutTime: &httpx.LogWithoutTime{
									URL:        "https://send.money.com/topup",
									StatusCode: 200,
									Request:    "POST /topup HTTP/1.1\r\nHost: send.money.com\r\nUser-Agent: Go-http-client/1.1\r\nAccept-Encoding: gzip\r\n\r\n",
									Response:   "HTTP/1.0 200 OK\r\nContent-Length: 14\r\n\r\n{\"errors\":[]}",
									ElapsedMS:  12,
								},
								Status: core.CallStatusSuccess,
							},
							CreatedOn: dates.Now(),
						},
					},
				)
			},
			`airtime_created`,
		},
		{
			func() events.Event {
				return events.NewBroadcastCreated(
					core.BroadcastTranslations{
						"eng": {Text: "Hello", Attachments: nil, QuickReplies: nil},
						"spa": {Text: "Hola", Attachments: nil, QuickReplies: nil},
					},
					i18n.Language("eng"),
					[]*assets.GroupReference{
						assets.NewGroupReference(assets.GroupUUID("5f9fd4f7-4b0f-462a-a598-18bfc7810412"), "Supervisors"),
					},
					[]*core.ContactReference{
						core.NewContactReference(core.ContactUUID("b2aaf598-1bb3-4c7d-b6bb-1f8dbe2ac16f"), "Jim"),
					},
					"name = \"Bob\"",
					[]urns.URN{urns.URN("tel:+12345678900")},
					nil,
					nil,
				)
			},
			`broadcast_created`,
		},
		{
			func() events.Event {
				return events.NewCallCreated(call.Marshal())
			},
			`call_created`,
		},
		{
			func() events.Event {
				return events.NewCallMissed(facebook.Reference())
			},
			`call_missed`,
		},
		{
			func() events.Event {
				return events.NewCallReceived(call.Marshal())
			},
			`call_received`,
		},
		{
			func() events.Event {
				return events.NewChatStarted(facebook.Reference(), map[string]string{"referrer_id": "acme"})
			},
			`chat_started`,
		},
		{
			func() events.Event {
				return events.NewContactFieldChanged(
					gender.Reference(),
					core.NewValue(types.NewXText("male"), nil, nil, "", "", ""),
				)
			},
			`contact_field_changed`,
		},
		{
			func() events.Event {
				return events.NewContactFieldChanged(
					gender.Reference(),
					nil, // value being cleared
				)
			},
			`contact_field_changed_clear`,
		},
		{
			func() events.Event {
				return events.NewContactGroupsChanged(
					core.GroupReferences([]*core.Group{session.Assets().Groups().FindByName("Customers")}),
					nil,
				)
			},
			`contact_groups_changed`,
		},
		{
			func() events.Event {
				return events.NewContactLanguageChanged(i18n.Language("fra"))
			},
			`contact_language_changed`,
		},
		{
			func() events.Event {
				return events.NewContactLastSeenChanged(time.Date(2022, 2, 3, 13, 45, 30, 0, time.UTC))
			},
			`contact_last_seen_changed`,
		},
		{
			func() events.Event {
				return events.NewContactNameChanged("Bryan")
			},
			`contact_name_changed`,
		},
		{
			func() events.Event {
				return events.NewContactStatusChanged(core.ContactStatusActive)
			},
			`contact_status_changed_active`,
		},
		{
			func() events.Event {
				return events.NewContactStatusChanged(core.ContactStatusBlocked)
			},
			`contact_status_changed_blocked`,
		},
		{
			func() events.Event {
				return events.NewContactStatusChanged(core.ContactStatusStopped)
			},
			`contact_status_changed_stopped`,
		},
		{
			func() events.Event {
				return events.NewContactTimezoneChanged(tz)
			},
			`contact_timezone_changed`,
		},
		{
			func() events.Event {
				return events.NewContactURNsChanged([]urns.URN{
					urns.URN("tel:+12345678900"),
					urns.URN("twitterid:8764843252522#bob"),
				})
			},
			`contact_urns_changed`,
		},
		{
			func() events.Event {
				return events.NewDialEnded(core.NewDial(core.DialStatusBusy, 0))
			},
			`dial_ended`,
		},
		{
			func() events.Event {
				return events.NewDialWait(urns.URN("tel:+1234567890"), 20, 120, time.Date(2022, 2, 3, 13, 45, 30, 0, time.UTC))
			},
			`dial_wait`,
		},
		{
			func() events.Event {
				return events.NewEmailSent([]string{"bob@nyaruka.com", "jim@nyaruka.com"}, "Update", "Flows are great!")
			},
			`email_sent`,
		},
		{
			func() events.Event {
				return events.NewError("I'm an error", "test_error")
			},
			`error`,
		},
		{
			func() events.Event {
				return events.NewDependencyError(assets.NewFieldReference("age", "Age"))
			},
			`error_dependency`,
		},
		{
			func() events.Event {
				return events.NewFailure("503 is an failure")
			},
			`failure`,
		},
		{
			func() events.Event {
				return events.NewIVRCreated(
					core.NewIVRMsgOut(
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
			func() events.Event {
				return events.NewLLMCalled(
					gpt4.Reference(),
					"Categorize the following text as Positive or Negative",
					"Please stop messaging me",
					&core.LLMResponse{Output: "Positive", TokensInput: 234, TokensOutput: 333},
					123*time.Millisecond,
				)
			},
			`llm_called`,
		},
		{
			func() events.Event {
				return events.NewMsgReceived(
					core.NewMsgIn(
						urns.URN("tel:+12065551212"),
						assets.NewChannelReference(assets.ChannelUUID("57f1078f-88aa-46f4-a59a-948a5739c03d"), "My Android Phone"),
						"hi there",
						nil,
						"",
					),
					"",
				)
			},
			`msg_received`,
		},
		{
			func() events.Event {
				return events.NewMsgReceived(
					core.NewMsgIn(
						urns.URN("tel:+12065551212"),
						assets.NewChannelReference(assets.ChannelUUID("57f1078f-88aa-46f4-a59a-948a5739c03d"), "My Android Phone"),
						"hi there",
						[]utils.Attachment{"image/jpeg:https://s3.amazon.com/mybucket/attachment.jpg"},
						"ext-id-123",
					),
					"7481888c-07dd-47dc-bf22-ef7448696ffe",
				)
			},
			`msg_received_rich`,
		},
		{
			func() events.Event {
				return events.NewMsgCreated(
					core.NewMsgOut(
						urns.URN("tel:+12345678900"),
						assets.NewChannelReference(assets.ChannelUUID("57f1078f-88aa-46f4-a59a-948a5739c03d"), "My Android Phone"),
						&core.MsgContent{Text: "Hi there"},
						nil,
						i18n.NilLocale,
						"",
					),
					"",
					"",
				)
			},
			`msg_created`,
		},
		{
			func() events.Event {
				return events.NewMsgCreated(
					core.NewMsgOut(
						urns.URN("tel:+12345678900"),
						assets.NewChannelReference(assets.ChannelUUID("57f1078f-88aa-46f4-a59a-948a5739c03d"), "My Android Phone"),
						&core.MsgContent{
							Text:         "Hi there",
							Attachments:  []utils.Attachment{"image/jpeg:http://s3.amazon.com/bucket/test.jpg"},
							QuickReplies: []core.QuickReply{{Text: "yes"}, {Text: "no"}},
						},
						nil,
						"eng-US",
						core.UnsendableReasonContactBlocked,
					),
					"",
					"01990b6d-de7e-7d28-8e40-806ac2c2f3f2",
				)
			},
			`msg_created_rich`,
		},
		{
			func() events.Event {
				return events.NewMsgDeleted("01990b6d-de7e-7d28-8e40-806ac2c2f3f2", false)
			},
			`msg_deleted`,
		},
		{
			func() events.Event {
				return events.NewMsgDeleted("01990b6d-de7e-7d28-8e40-806ac2c2f3f2", true)
			},
			`msg_deleted_by_contact`,
		},
		{
			func() events.Event {
				return events.NewMsgStatusChanged("01990b6d-de7e-7d28-8e40-806ac2c2f3f2", "sent", "")
			},
			`msg_status_changed_sent`,
		},
		{
			func() events.Event {
				return events.NewMsgStatusChanged("01990b6d-de7e-7d28-8e40-806ac2c2f3f2", "failed", "error_limit")
			},
			`msg_status_changed_failed`,
		},
		{
			func() events.Event {
				return events.NewMsgWait(&timeout, time.Date(2022, 2, 3, 13, 45, 30, 0, time.UTC), hints.NewImage())
			},
			`msg_wait`,
		},
		{
			func() events.Event {
				return events.NewOptInRequested(jotd.Reference(), facebook.Reference(), urns.URN("facebook:1234567890"))
			},
			`optin_requested`,
		},
		{
			func() events.Event {
				return events.NewOptInStarted(jotd.Reference(), facebook.Reference())
			},
			`optin_started`,
		},
		{
			func() events.Event {
				return events.NewOptInStopped(jotd.Reference(), facebook.Reference())
			},
			`optin_stopped`,
		},
		{
			func() events.Event {
				return events.NewRunResultChanged(core.NewResult("Age", "44", "", "", "78c4513d-61a1-428b-80d7-3bffd39b74f2", "", nil, time.Date(2025, 9, 1, 13, 45, 30, 0, time.UTC)), nil)
			},
			`run_result_changed`,
		},
		{
			func() events.Event {
				return events.NewRunResultChanged(
					core.NewResult("Age", "44", "", "", "78c4513d-61a1-428b-80d7-3bffd39b74f2", "", nil, time.Date(2025, 9, 1, 13, 45, 30, 0, time.UTC)),
					core.NewResult("Age", "43", "", "", "78c4513d-61a1-428b-80d7-3bffd39b74f2", "", nil, time.Date(2024, 9, 1, 13, 45, 30, 0, time.UTC)),
				)
			},
			`run_result_changed_with_previous`,
		},
		{
			func() events.Event {
				run := session.Runs()[0]
				return events.NewRunStarted(run.FlowReference(), run.UUID(), run.Parent().UUID(), true)
			},
			`run_started`,
		},
		{
			func() events.Event {
				return events.NewRunEnded(
					"01990b6d-de7e-7d28-8e40-806ac2c2f3f2",
					assets.NewFlowReference(assets.FlowUUID("e4d441f0-24e3-4627-85fb-1e99e733baf0"), "Collect Age"),
					core.RunStatusCompleted,
				)
			},
			`run_ended`,
		},
		{
			func() events.Event {
				return events.NewSessionTriggered(
					assets.NewFlowReference(assets.FlowUUID("e4d441f0-24e3-4627-85fb-1e99e733baf0"), "Collect Age"),
					[]*assets.GroupReference{
						assets.NewGroupReference(assets.GroupUUID("5f9fd4f7-4b0f-462a-a598-18bfc7810412"), "Supervisors"),
					},
					[]*core.ContactReference{
						core.NewContactReference(core.ContactUUID("b2aaf598-1bb3-4c7d-b6bb-1f8dbe2ac16f"), "Jim"),
					},
					"age > 20",
					events.Exclusions{InAFlow: true},
					false,
					[]urns.URN{urns.URN("tel:+12345678900")},
					json.RawMessage(`{"uuid": "779eaf3f-1c59-4374-a7cb-0eae9c5e8800"}`),
					&core.SessionHistory{ParentUUID: "418a704c-f33e-4924-a00e-1763d1498a13", Ancestors: 2, AncestorsSinceInput: 0},
				)
			},
			`session_triggered`,
		},
		{
			func() events.Event {
				return events.NewTicketAssigneeChanged("019905d4-5f7b-71b8-bcb8-6a68de2d91d2", user.Reference(), nil)
			},
			`ticket_assignee_changed`,
		},
		{
			func() events.Event {
				return events.NewTicketAssigneeChanged("019905d4-5f7b-71b8-bcb8-6a68de2d91d2", nil, user.Reference())
			},
			`ticket_assignee_changed_nobody`,
		},
		{
			func() events.Event {
				return events.NewTicketClosed(ticket.UUID())
			},
			`ticket_closed`,
		},
		{
			func() events.Event {
				return events.NewTicketNoteAdded("019905d4-5f7b-71b8-bcb8-6a68de2d91d2", "This looks important!")
			},
			`ticket_note_added`,
		},
		{
			func() events.Event {
				return events.NewTicketOpened(ticket.Marshal())
			},
			`ticket_opened`,
		},
		{
			func() events.Event {
				return events.NewTicketReopened("019905d4-5f7b-71b8-bcb8-6a68de2d91d2")
			},
			`ticket_reopened`,
		},
		{
			func() events.Event {
				return events.NewTicketTopicChanged("019905d4-5f7b-71b8-bcb8-6a68de2d91d2", weather.Reference())
			},
			`ticket_topic_changed`,
		},
		{
			func() events.Event {
				return events.NewTypingStarted(events.DirectionIncoming, facebook.Reference(), urns.URN("facebook:1234567890"), "EX12345")
			},
			`typing_started_incoming`,
		},
		{
			func() events.Event {
				return events.NewTypingStarted(events.DirectionOutgoing, nil, urns.NilURN, "")
			},
			`typing_started_outgoing`,
		},
		{
			func() events.Event {
				return events.NewTypingStopped(events.DirectionIncoming, facebook.Reference(), urns.URN("facebook:1234567890"), "EX12345")
			},
			`typing_stopped_incoming`,
		},
		{
			func() events.Event {
				return events.NewTypingStopped(events.DirectionOutgoing, nil, urns.NilURN, "")
			},
			`typing_stopped_outgoing`,
		},
		{
			func() events.Event {
				return events.NewWaitTimedOut()
			},
			`wait_timed_out`,
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

	// check setting of user reference field and non marshaling of steps
	var evt events.Event = events.NewTicketClosed(ticket.UUID())
	evt.SetUser(user.Reference(), "ui")
	evt.SetStep(&events.Step{Flow: assets.NewFlowReference("50c3706e-fedb-42c0-8eab-dda3335714b7", "Registration"), Node: "72ecb927-db78-4acf-b947-db2f29bf6662"})

	eventJSON, err := jsonx.MarshalPretty(evt)
	assert.NoError(t, err)

	test.AssertSnapshot(t, "ticket_closed_with_user_ref", string(eventJSON))
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
	client, _ := test.MockedHTTP(map[string][]*httpx.MockResponse{
		"http://temba.io/": {
			httpx.NewMockResponse(200, nil, bytes.Repeat([]byte("Y"), 20000)),
		},
	})

	request, _ := http.NewRequest("GET", "http://temba.io/", strings.NewReader(strings.Repeat("X", 20000)))

	svc := webhooks.NewService(client, nil, 1024*1024)
	call, err := svc.Call(request)
	require.NoError(t, err)

	assert.Equal(t, 42, len(call.ResponseTrace))
	assert.Equal(t, 20000, len(call.ResponseBody))

	event := events.NewWebhookCalled(call, core.CallStatusSuccess, "")

	assert.Equal(t, "http://temba.io/", event.URL)
	assert.Equal(t, 10000, len(event.Request))
	assert.Equal(t, "XXXXXXX...", event.Request[9990:])
	assert.Equal(t, 10000, len(event.Response))
	assert.Equal(t, "YYYYYYY...", event.Response[9990:])
}

func TestWebhookCalledEventValid(t *testing.T) {
	client, _ := test.MockedHTTP(map[string][]*httpx.MockResponse{
		"http://temba.io/": {
			httpx.NewMockResponse(200, map[string]string{"Header": "hello"}, []byte(`{"foo": "bar"}`)),
		},
	})

	request, _ := http.NewRequest("GET", "http://temba.io/", nil)

	svc := webhooks.NewService(client, nil, 1024*1024)
	call, err := svc.Call(request)
	require.NoError(t, err)

	event := events.NewWebhookCalled(call, core.CallStatusSuccess, "")

	assert.Equal(t, "http://temba.io/", event.URL)
	assert.Equal(t, "HTTP/1.0 200 OK\r\nContent-Length: 14\r\nHeader: hello\r\n\r\n{\"foo\": \"bar\"}", event.Response)
	assert.True(t, utf8.ValidString(event.Response))
}

func TestWebhookCalledEventNullChar(t *testing.T) {
	client, _ := test.MockedHTTP(map[string][]*httpx.MockResponse{
		"http://temba.io/": {
			httpx.NewMockResponse(200, nil, []byte("abc \x00 \\u0000 \\\u0000 \\\\u0000")),
		},
	})

	request, _ := http.NewRequest("GET", "http://temba.io/", nil)

	svc := webhooks.NewService(client, nil, 1024*1024)
	call, err := svc.Call(request)
	require.NoError(t, err)

	event := events.NewWebhookCalled(call, core.CallStatusSuccess, "")

	// actual null will have been stripped, escaped null will remain
	assert.Equal(t, "http://temba.io/", event.URL)
	assert.Equal(t, "HTTP/1.0 200 OK\r\nContent-Length: 23\r\n\r\nabc � � \\� \\\\u0000", event.Response)
	assert.True(t, utf8.ValidString(event.Response))
}

func TestWebhookCalledEventBadUTF8(t *testing.T) {
	client, _ := test.MockedHTTP(map[string][]*httpx.MockResponse{
		"http://temba.io/": {
			httpx.NewMockResponse(200, map[string]string{"Bad-Header": "\xa0\xa1"}, []byte("{\"foo\": \"\xa0\xa1\"}")),
		},
	})

	request, _ := http.NewRequest("GET", "http://temba.io/", nil)

	svc := webhooks.NewService(client, nil, 1024*1024)
	call, err := svc.Call(request)
	require.NoError(t, err)

	event := events.NewWebhookCalled(call, core.CallStatusSuccess, "")

	assert.Equal(t, "http://temba.io/", event.URL)
	assert.Equal(t, "HTTP/1.0 200 OK\r\nContent-Length: 13\r\nBad-Header: �\r\n\r\n...", event.Response)
	assert.True(t, utf8.ValidString(event.Response))
}
