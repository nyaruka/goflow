package events_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/flows/routers/waits/hints"
	"github.com/nyaruka/goflow/test"
	"github.com/nyaruka/goflow/utils/dates"
	"github.com/nyaruka/goflow/utils/httpx"
	"github.com/nyaruka/goflow/utils/uuids"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEventMarshaling(t *testing.T) {
	defer dates.SetNowSource(dates.DefaultNowSource)
	defer uuids.SetGenerator(uuids.DefaultGenerator)

	dates.SetNowSource(dates.NewFixedNowSource(time.Date(2018, 10, 18, 14, 20, 30, 123456, time.UTC)))
	uuids.SetGenerator(uuids.NewSeededGenerator(12345))

	session, _, err := test.CreateTestSession("", envs.RedactionPolicyNone)
	require.NoError(t, err)

	tz, _ := time.LoadLocation("Africa/Kigali")
	timeout := 500
	gender := session.Assets().Fields().Get("gender")

	eventTests := []struct {
		event     flows.Event
		marshaled string
	}{
		{
			events.NewBroadcastCreated(
				map[envs.Language]*events.BroadcastTranslation{
					"eng": {Text: "Hello", Attachments: nil, QuickReplies: nil},
					"spa": {Text: "Hola", Attachments: nil, QuickReplies: nil},
				},
				envs.Language("eng"),
				[]*assets.GroupReference{
					assets.NewGroupReference(assets.GroupUUID("5f9fd4f7-4b0f-462a-a598-18bfc7810412"), "Supervisors"),
				},
				[]*flows.ContactReference{
					flows.NewContactReference(flows.ContactUUID("b2aaf598-1bb3-4c7d-b6bb-1f8dbe2ac16f"), "Jim"),
				},
				[]urns.URN{urns.URN("tel:+12345678900")},
			),
			`{
				"base_language": "eng",
				"contacts": [
					{
						"name": "Jim",
						"uuid": "b2aaf598-1bb3-4c7d-b6bb-1f8dbe2ac16f"
					}
				],
				"created_on": "2018-10-18T14:20:30.000123456Z",
				"groups": [
					{
						"name": "Supervisors",
						"uuid": "5f9fd4f7-4b0f-462a-a598-18bfc7810412"
					}
				],
				"translations": {
					"eng": {
						"text": "Hello"
					},
					"spa": {
						"text": "Hola"
					}
				},
				"type": "broadcast_created",
				"urns": [
					"tel:+12345678900"
				]
			}`,
		},
		{
			events.NewClassifierCalled(
				assets.NewClassifierReference(assets.ClassifierUUID("4b937f49-7fb7-43a5-8e57-14e2f028a471"), "Booking"),
				[]*httpx.Trace{httpx.NewMockTrace("GET", "https://api.wit.ai/message?v=20170307&q=hello", 200, `{"intents":[]}`)},
			),
			`{
				"classifier": {
					"uuid": "4b937f49-7fb7-43a5-8e57-14e2f028a471",
					"name": "Booking"
				},
				"created_on": "2018-10-18T14:20:30.000123456Z",
				"http_logs": [
					{
						"created_on": "2018-10-18T14:20:30.000123456Z",
						"elapsed_ms": 0,
						"request": "GET /message?v=20170307&q=hello HTTP/1.1\r\nHost: api.wit.ai\r\nUser-Agent: Go-http-client/1.1\r\nAccept-Encoding: gzip\r\n\r\n",
						"response": "HTTP/1.0 200 OK\r\nContent-Length: 14\r\n\r\n{\"intents\":[]}",
						"status": "success",
						"url": "https://api.wit.ai/message?v=20170307&q=hello"
					}
				],
				"type": "classifier_called"
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
			events.NewContactLanguageChanged(envs.Language("fra")),
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
					"name": "Ryan Lewis",
					"timezone": "America/Guayaquil",
					"urns": [
						"tel:+12065551212?channel=57f1078f-88aa-46f4-a59a-948a5739c03d",
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
			events.NewEnvironmentRefreshed(session.Environment()),
			`{
				"created_on": "2018-10-18T14:20:30.000123456Z",
				"environment": {
					"allowed_languages": [
						"eng",
						"spa"
					],
					"date_format": "DD-MM-YYYY",
					"default_language": "eng",
					"max_value_length": 640,
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
			events.NewIVRCreated(
				flows.NewMsgOut(
					urns.URN("tel:+12345678900"),
					assets.NewChannelReference(assets.ChannelUUID("57f1078f-88aa-46f4-a59a-948a5739c03d"), "My Android Phone"),
					"Hi there",
					nil,
					nil,
					nil)),
			`{
				"created_on": "2018-10-18T14:20:30.000123456Z",
				"msg": {
					"channel": {
						"name": "My Android Phone",
						"uuid": "57f1078f-88aa-46f4-a59a-948a5739c03d"
					},
					"text": "Hi there",
					"urn": "tel:+12345678900",
					"uuid": "20cc4181-48cf-4344-9751-99419796decd"
				},
				"type": "ivr_created"
			}`,
		},
		{
			events.NewMsgWait(&timeout, hints.NewImageHint()),
			`{
				"created_on": "2018-10-18T14:20:30.000123456Z",
				"hint": {"type": "image"},
				"timeout_seconds": 500,
				"type": "msg_wait"
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
				false,
				[]urns.URN{urns.URN("tel:+12345678900")},
				json.RawMessage(`{"uuid": "779eaf3f-1c59-4374-a7cb-0eae9c5e8800"}`),
			),
			`{
				"contacts": [
					{
						"name": "Jim",
						"uuid": "b2aaf598-1bb3-4c7d-b6bb-1f8dbe2ac16f"
					}
				],
				"contact_query": "age > 20",
				"created_on": "2018-10-18T14:20:30.000123456Z",
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
				"run_summary": {
					"uuid": "779eaf3f-1c59-4374-a7cb-0eae9c5e8800"
				},
				"type": "session_triggered",
				"urns": [
					"tel:+12345678900"
				]
			}`,
		},
	}

	for _, tc := range eventTests {
		eventJSON, err := json.Marshal(tc.event)
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
}
