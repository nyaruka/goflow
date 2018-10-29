package events_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/test"
	"github.com/nyaruka/goflow/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEventMarshaling(t *testing.T) {
	utils.SetTimeSource(utils.NewFixedTimeSource(time.Date(2018, 10, 18, 14, 20, 30, 123456, time.UTC)))
	defer utils.SetTimeSource(utils.DefaultTimeSource)

	session, _, err := test.CreateTestSession("", nil)
	require.NoError(t, err)

	tz, _ := time.LoadLocation("Africa/Kigali")

	eventTests := []struct {
		event     flows.Event
		marshaled string
	}{
		{
			events.NewContactFieldChangedEvent(
				assets.NewFieldReference("gender", "Gender"),
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
			events.NewContactFieldChangedEvent(
				assets.NewFieldReference("gender", "Gender"),
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
			events.NewContactGroupsChangedEvent(
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
			events.NewContactLanguageChangedEvent(utils.Language("fra")),
			`{
				"created_on": "2018-10-18T14:20:30.000123456Z",
				"language": "fra",
				"type": "contact_language_changed"
			}`,
		},
		{
			events.NewContactNameChangedEvent("Bryan"),
			`{
				"created_on": "2018-10-18T14:20:30.000123456Z",
				"name": "Bryan",
				"type": "contact_name_changed"
			}`,
		},
		{
			events.NewContactTimezoneChangedEvent(tz),
			`{
				"created_on": "2018-10-18T14:20:30.000123456Z",
				"timezone": "Africa/Kigali",
				"type": "contact_timezone_changed"
			}`,
		},
		{
			events.NewContactURNsChangedEvent([]urns.URN{
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
			events.NewSessionTriggeredEvent(
				assets.NewFlowReference(assets.FlowUUID("e4d441f0-24e3-4627-85fb-1e99e733baf0"), "Collect Age"),
				[]urns.URN{urns.URN("tel:+12345678900")},
				[]*flows.ContactReference{
					flows.NewContactReference(flows.ContactUUID("b2aaf598-1bb3-4c7d-b6bb-1f8dbe2ac16f"), "Jim"),
				},
				[]*assets.GroupReference{
					assets.NewGroupReference(assets.GroupUUID("5f9fd4f7-4b0f-462a-a598-18bfc7810412"), "Supervisors"),
				},
				false,
				json.RawMessage(`{"uuid": "779eaf3f-1c59-4374-a7cb-0eae9c5e8800"}`),
			),
			`{
				"contacts": [
					{
						"name": "Jim",
						"uuid": "b2aaf598-1bb3-4c7d-b6bb-1f8dbe2ac16f"
					}
				],
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
		require.NoError(t, err)

		test.AssertEqualJSON(t, []byte(tc.marshaled), eventJSON, "event JSON mismatch")

		// also try to unmarshal and validate the JSON
		err = utils.UnmarshalAndValidate(eventJSON, tc.event)
		require.NoError(t, err)
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
