package modifiers_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/actions/modifiers"
	"github.com/nyaruka/goflow/test"
	"github.com/nyaruka/goflow/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestModifierTypes(t *testing.T) {
	assets, err := test.LoadSessionAssets("testdata/_assets.json")
	require.NoError(t, err)

	for typeName := range modifiers.RegisteredTypes {
		testModifierType(t, assets, typeName)
	}
}

func testModifierType(t *testing.T, assets flows.SessionAssets, typeName string) {
	testFile, err := ioutil.ReadFile(fmt.Sprintf("testdata/%s.json", typeName))
	require.NoError(t, err)

	tests := []struct {
		Description   string            `json:"description"`
		ContactBefore json.RawMessage   `json:"contact_before"`
		Modifier      json.RawMessage   `json:"modifier"`
		ContactAfter  json.RawMessage   `json:"contact_after"`
		Events        []json.RawMessage `json:"events"`
	}{}

	err = json.Unmarshal(testFile, &tests)
	require.NoError(t, err)

	defer utils.SetTimeSource(utils.DefaultTimeSource)
	defer utils.SetUUIDGenerator(utils.DefaultUUIDGenerator)

	for _, tc := range tests {
		utils.SetTimeSource(utils.NewFixedTimeSource(time.Date(2018, 10, 18, 14, 20, 30, 123456, time.UTC)))
		utils.SetUUIDGenerator(utils.NewSeededUUID4Generator(12345))

		testName := fmt.Sprintf("test '%s' for modifier type '%s'", tc.Description, typeName)

		// read the modifier to be tested
		modifier, err := modifiers.ReadModifier(assets, tc.Modifier)
		require.NoError(t, err, "error loading modifier in %s", testName)
		assert.Equal(t, typeName, modifier.Type())

		// read the initial contact state
		contact, err := flows.ReadContact(assets, tc.ContactBefore, true)
		require.NoError(t, err, "error loading contact_before in %s", testName)

		// apply the modifier
		logEvent := make([]flows.Event, 0)
		modifier.Apply(utils.NewDefaultEnvironment(), assets, contact, func(e flows.Event) { logEvent = append(logEvent, e) })

		// check contact is in the expected state
		contactJSON, _ := json.Marshal(contact)
		test.AssertEqualJSON(t, tc.ContactAfter, contactJSON, "contact mismatch in %s", testName)

		// check events are what we expected
		actualEventsJSON, _ := json.Marshal(logEvent)
		expectedEventsJSON, _ := json.Marshal(tc.Events)
		test.AssertEqualJSON(t, expectedEventsJSON, actualEventsJSON, "events mismatch in %s", testName)

		// try marshaling the modifier back to JSON
		modifierJSON, err := json.Marshal(modifier)
		require.NoError(t, err)
		test.AssertEqualJSON(t, tc.Modifier, modifierJSON, "marshal mismatch in %s", testName)
	}
}

func TestConstructors(t *testing.T) {
	assets, err := test.LoadSessionAssets("testdata/_assets.json")
	require.NoError(t, err)

	nexmo, _ := assets.Channels().Get("3a05eaf5-cb1b-4246-bef1-f277419c83a7")
	age, _ := assets.Fields().Get("age")
	ageValue := types.NewXNumberFromInt(37)
	testers, _ := assets.Groups().Get("b7cf0d83-f1c9-411c-96fd-c511a4cfa86d")
	la, _ := time.LoadLocation("America/Los_Angeles")

	tests := []struct {
		modifier flows.Modifier
		json     string
	}{
		{
			modifiers.NewChannelModifier(nexmo),
			`{
				"type": "channel",
				"channel": {
					"uuid": "3a05eaf5-cb1b-4246-bef1-f277419c83a7",
					"name": "Nexmo"
				}
			}`,
		},
		{
			modifiers.NewFieldModifier(age, flows.NewValue(types.NewXText("37 years"), nil, &ageValue, "", "", "")),
			`{
				"type": "field",
				"field": {
					"key": "age",
					"name": "Age"
				},
				"value": {
					"text": "37 years",
					"number": 37
				}
			}`,
		},
		{
			modifiers.NewGroupsModifier([]*flows.Group{testers}, modifiers.GroupsAdd),
			`{
				"type": "groups",
				"groups": [
					{
						"uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d",
						"name": "Testers"
					}
				],
				"modification": "add"
			}`,
		},
		{
			modifiers.NewLanguageModifier(utils.Language("fra")),
			`{
				"type": "language",
				"language": "fra"
			}`,
		},
		{
			modifiers.NewNameModifier("Bob"),
			`{
				"type": "name",
				"name": "Bob"
			}`,
		},
		{
			modifiers.NewTimezoneModifier(la),
			`{
				"type": "timezone",
				"timezone": "America/Los_Angeles"
			}`,
		},
		{
			modifiers.NewURNModifier(urns.URN("tel:+1234567890"), modifiers.URNAppend),
			`{
				"type": "urn",
				"urn": "tel:+1234567890",
				"modification": "append"
			}`,
		},
	}

	for _, tc := range tests {
		modifierJSON, err := json.Marshal(tc.modifier)
		require.NoError(t, err)
		test.AssertEqualJSON(t, []byte(tc.json), modifierJSON, "marshal mismatch for modifier %s", string(modifierJSON))
	}
}

func TestReadModifier(t *testing.T) {
	// error if no type field
	_, err := modifiers.ReadModifier(nil, []byte(`{"foo": "bar"}`))
	assert.EqualError(t, err, "field 'type' is required")

	// error if we don't recognize the type
	_, err = modifiers.ReadModifier(nil, []byte(`{"type": "do_the_foo", "foo": "bar"}`))
	assert.EqualError(t, err, "unknown type: 'do_the_foo'")
}
