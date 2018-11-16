package mobile_test

import (
	"io/ioutil"
	"testing"

	"github.com/nyaruka/goflow/mobile"
	"github.com/nyaruka/goflow/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMobileBindings(t *testing.T) {
	assert.True(t, mobile.IsSpecVersionSupported("12"))
	assert.False(t, mobile.IsSpecVersionSupported("11.6"))

	// error if we try to create assets from invalid JSON
	_, err := mobile.NewAssetsSource("{")
	assert.Error(t, err)

	// can load a standard assets file
	assetsJSON, err := ioutil.ReadFile("../test/testdata/flows/two_questions_offline.json")
	require.NoError(t, err)

	source, err := mobile.NewAssetsSource(string(assetsJSON))
	require.NoError(t, err)

	// and create a new session assets
	sessionAssets, err := mobile.NewSessionAssets(source)
	require.NoError(t, err)

	environment, err := mobile.NewEnvironment("DD-MM-YYYY", "tt:mm", "Africa/Kigali", "eng", []string{"eng", "fra"})
	require.NoError(t, err)

	contact := mobile.NewEmptyContact()

	trigger := mobile.NewManualTrigger(environment, contact, "7c3db26f-e12a-48af-9673-e2feefdf8516", "Two Questions")

	session := mobile.NewSession(sessionAssets, "mobile-test")

	events, err := session.Start(trigger)
	require.NoError(t, err)

	assert.Equal(t, 2, len(events))
	assert.Equal(t, "msg_created", events[0].Type())
	assert.Equal(t, "msg_wait", events[1].Type())

	resume := mobile.NewMsgResume(nil, nil, mobile.NewMsgIn("8e6f0213-a122-4c50-a430-442085754c16", "Hi there", nil))

	events, err = session.Resume(resume)
	require.NoError(t, err)

	assert.Equal(t, 4, len(events))
	assert.Equal(t, "msg_received", events[0].Type())
	assert.Equal(t, "run_result_changed", events[1].Type())
	assert.Equal(t, "msg_created", events[2].Type())
	assert.Equal(t, "msg_wait", events[3].Type())
}

func TestMigrateLegacyFlow(t *testing.T) {
	// error if legacy definition isn't valid
	_, err := mobile.MigrateLegacyFlow(`{"metadata": {}}`)
	assert.EqualError(t, err, `unable to read legacy flow: field 'metadata.uuid' is required`)

	migrated, err := mobile.MigrateLegacyFlow(`{
		"flow_type": "S", 
		"action_sets": [],
		"rule_sets": [],
		"base_language": "eng",
		"metadata": {
			"uuid": "061be894-4507-470c-a20b-34273bf915be",
			"name": "Survey"
		}
	}`)
	assert.NoError(t, err)
	test.AssertEqualJSON(t, []byte(`{
		"uuid": "061be894-4507-470c-a20b-34273bf915be",
		"name": "Survey",
		"spec_version": "12.0",
		"type": "messaging_offline",
		"expire_after_minutes": 0,
		"language": "eng",
		"localization": {},
		"nodes": [],
		"revision": 0
	}`), []byte(migrated), "migrated flow mismatch")
}
