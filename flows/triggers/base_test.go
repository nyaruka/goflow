package triggers_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/triggers"
	"github.com/nyaruka/goflow/test"
	"github.com/nyaruka/goflow/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTriggerMarshaling(t *testing.T) {
	utils.SetTimeSource(utils.NewFixedTimeSource(time.Date(2018, 10, 18, 14, 20, 30, 123456, time.UTC)))
	defer utils.SetTimeSource(utils.DefaultTimeSource)

	utils.SetUUIDGenerator(utils.NewSeededUUID4Generator(1234))
	defer utils.SetUUIDGenerator(utils.DefaultUUIDGenerator)

	contact := flows.NewEmptyContact("Bob", utils.Language("eng"), nil)
	flow := assets.NewFlowReference(assets.FlowUUID("7c37d7e5-6468-4b31-8109-ced2ef8b5ddc"), "Registration")

	triggerTests := []struct {
		trigger   flows.Trigger
		marshaled string
	}{
		{
			triggers.NewFlowActionTrigger(
				utils.NewDefaultEnvironment(),
				flow,
				contact,
				json.RawMessage(`{"uuid": "084e4bed-667c-425e-82f7-bdb625e6ec9e"}`),
				time.Date(2018, 10, 20, 9, 49, 30, 1234567890, time.UTC),
			),
			`{
				"contact": {
					"created_on": "2018-10-18T14:20:30.000123456Z",
					"id": 0,
					"language": "eng",
					"name": "Bob",
					"urns": [],
					"uuid": "c00e5d67-c275-4389-aded-7d8b151cbd5b"
				},
				"environment": {
					"date_format": "YYYY-MM-DD",
					"redaction_policy": "none",
					"time_format": "tt:mm",
					"timezone": "UTC"
				},
				"flow": {
					"name": "Registration",
					"uuid": "7c37d7e5-6468-4b31-8109-ced2ef8b5ddc"
				},
				"run_summary": {
					"uuid": "084e4bed-667c-425e-82f7-bdb625e6ec9e"
				},
				"triggered_on": "2018-10-20T09:49:31.23456789Z",
				"type": "campaign"
			}`,
		},
	}

	for _, tc := range triggerTests {
		eventJSON, err := json.Marshal(tc.trigger)
		require.NoError(t, err)

		test.AssertEqualJSON(t, []byte(tc.marshaled), eventJSON, "trigger JSON mismatch")

		// also try to unmarshal and validate the JSON
		err = utils.UnmarshalAndValidate(eventJSON, tc.trigger)
		require.NoError(t, err)
	}
}

func TestReadTrigger(t *testing.T) {
	// error if no type field
	_, err := triggers.ReadTrigger(nil, []byte(`{"foo": "bar"}`))
	assert.EqualError(t, err, "field 'type' is required")

	// error if we don't recognize action type
	_, err = triggers.ReadTrigger(nil, []byte(`{"type": "do_the_foo", "foo": "bar"}`))
	assert.EqualError(t, err, "unknown type: 'do_the_foo'")
}
