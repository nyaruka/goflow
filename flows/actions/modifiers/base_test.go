package modifiers_test

import (
	"encoding/json"
	"fmt"
	"github.com/nyaruka/goflow/assets/static"
	"io/ioutil"
	"testing"
	"time"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/actions/modifiers"
	"github.com/nyaruka/goflow/flows/engine"
	"github.com/nyaruka/goflow/test"
	"github.com/nyaruka/goflow/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestModifierTypes(t *testing.T) {
	assetsJSON, err := ioutil.ReadFile("testdata/_assets.json")
	require.NoError(t, err)

	source, err := static.NewStaticSource(assetsJSON)
	require.NoError(t, err)

	assets, err := engine.NewSessionAssets(source)
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
		eventLog := make([]flows.Event, 0)
		modifier.Apply(utils.NewDefaultEnvironment(), assets, contact, func(e flows.Event) { eventLog = append(eventLog, e) })

		// check contact is in the expected state
		contactJSON, _ := json.Marshal(contact)
		test.AssertEqualJSON(t, tc.ContactAfter, contactJSON, "contact mismatch in %s", testName)

		// check events are what we expected
		actualEventsJSON, _ := json.Marshal(eventLog)
		expectedEventsJSON, _ := json.Marshal(tc.Events)
		test.AssertEqualJSON(t, expectedEventsJSON, actualEventsJSON, "events mismatch in %s", testName)

		// try marshaling the modifier back to JSON
		modifierJSON, err := json.Marshal(modifier)
		test.AssertEqualJSON(t, tc.Modifier, modifierJSON, "marshal mismatch in %s", testName)
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
