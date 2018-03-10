package engine

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/nyaruka/goflow/excellent"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/flows/triggers"
	"github.com/nyaruka/goflow/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testRequest struct {
	Trigger *utils.TypedEnvelope   `json:"trigger"`
	Events  []*utils.TypedEnvelope `json:"events"`
}

func TestEvaluateTemplate(t *testing.T) {
	tests := []struct {
		template string
		expected string
		hasError bool
	}{
		{"@contact.uuid", "ba96bf7f-bc2a-4873-a7c7-254d1927c4e3", false},
		{"@contact.name", "Ben Haggerty", false},
		{"@contact.first_name", "Ben", false},
		{"@contact.language", "eng", false},
		{"@contact.timezone", "America/Guayaquil", false},
		{"@contact.urns", "tel:+12065551212", false}, // TODO should be list
		{"@contact.urns.tel", "tel:+12065551212", false},
		{"@contact.urns.0", "tel:+12065551212", false},
		{"@contact.urns.0.scheme", "tel", false},
		{"@contact.urns.0.path", "+12065551212", false},
		{"@contact.urns.0.display", "", false},
		{"@contact.urns.0.channel", "Nexmo", false},
		{"@contact.urns.0.channel.uuid", "57f1078f-88aa-46f4-a59a-948a5739c03d", false},
		{"@contact.urns.0.channel.name", "Nexmo", false},
		{"@contact.urns.0.channel.address", "+12345671111", false},
		{"@contact.urns.1", "facebook:1122334455667788", false},
		{"@contact.urns.1.channel", "", false},
		{"@(format_urn(contact.urns.0))", "(206) 555-1212", false},
		{"@contact.groups", "Survey Audience", false}, // TODO should be list
		{"@contact.fields.state", "Azuay", false},
		{"@contact.fields.favorite_icecream", "", false},                   // TODO should be error
		{"@(has_error(contact.fields.favorite_icecream))", "false", false}, // TODO should be true

		{"@run.input", "Hi there\nhttp://s3.amazon.com/bucket/test_en.jpg?a=Azuay", false},
		{"@run.input.text", "Hi there", false},
		{"@run.input.attachments", "", false}, // TODO should be list
		{"@run.input.attachments.0", "http://s3.amazon.com/bucket/test_en.jpg?a=Azuay", false},
		{"@run.input.created_on", "2000-01-01T00:00:00.000000Z", false},
		{"@run.input.channel.name", "Nexmo", false},
		{"@run.results", "", false},                                     // TODO should be empty dict?
		{"@run.results.favorite_icecream", "", false},                   // TODO should be error
		{"@(has_error(run.results.favorite_icecream))", "false", false}, // TODO should be true
		{"@run.exited_on", "", false},

		{"@trigger.params", "{\n            \"coupons\": [\n                {\n                    \"code\": \"AAA-BBB-CCC\",\n                    \"expiration\": \"2000-01-01T00:00:00.000000000-00:00\"\n                }\n            ]\n        }", false},
		{"@trigger.params.coupons.0.code", "AAA-BBB-CCC", false},
	}

	assetsJSON, err := ioutil.ReadFile("testdata/assets.json")
	require.NoError(t, err)

	// build our session
	assetCache := NewAssetCache(100, 5, "testing/1.0")
	err = assetCache.Include(assetsJSON)
	require.NoError(t, err)

	session := NewSession(assetCache, NewMockAssetServer())
	require.NoError(t, err)

	// read trigger from file
	requestJSON, err := ioutil.ReadFile("testdata/trigger.json")
	require.NoError(t, err)

	testRequest := testRequest{}
	err = json.Unmarshal(requestJSON, &testRequest)
	require.NoError(t, err)

	trigger, err := triggers.ReadTrigger(session, testRequest.Trigger)
	require.NoError(t, err)

	initialEvents, err := events.ReadEvents(testRequest.Events)
	require.NoError(t, err)

	session.Start(trigger, initialEvents)
	run := session.Runs()[0]

	for _, test := range tests {
		eval, err := excellent.EvaluateTemplateAsString(session.Environment(), run.Context(), test.template, false)
		if err != nil {
			assert.True(t, test.hasError, "Received error evaluating '%s': %s", test.template, err)
		} else {
			assert.False(t, test.hasError, "Did not receive error evaluating '%s'", test.template)
		}

		assert.Equal(t, test.expected, eval, "Actual '%s' does not match expected '%s' evaluating template: '%s'", eval, test.expected, test.template)
	}
}
