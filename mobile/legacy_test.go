package mobile_test

import (
	"testing"

	"github.com/nyaruka/goflow/mobile"
	"github.com/nyaruka/goflow/test"

	"github.com/stretchr/testify/assert"
)

func TestIsLegacyDefinition(t *testing.T) {
	// try reading empty JSON
	assert.False(t, mobile.IsLegacyDefinition(`{}`))

	// try with new flow
	assert.False(t, mobile.IsLegacyDefinition(`{
		"uuid": "76f0a02f-3b75-4b86-9064-e9195e1b3a02",
		"name": "Simple",
		"spec_version": "13.0",
		"language": "eng",
		"type": "messaging",
		"nodes": []
	}`))

	// try with legacy flow
	assert.True(t, mobile.IsLegacyDefinition(`{
		"metadata": {
			"uuid": "76f0a02f-3b75-4b86-9064-e9195e1b3a02",
			"name": "Simple",
			"revision": 1
		},
		"base_language": "eng",
		"flow_type": "F",
		"version": 11,
		"action_sets": [],
		"rule_sets": []
	}`))
}

func TestMigrateLegacyDefinition(t *testing.T) {
	migrated, err := mobile.MigrateLegacyDefinition(`{
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
		"spec_version": "13.0.0",
		"type": "messaging_offline",
		"expire_after_minutes": 0,
		"language": "eng",
		"localization": {},
		"nodes": [],
		"revision": 0,
		"_ui": {
			"nodes": {},
			"stickies": {}
		}
	}`), []byte(migrated), "migrated flow mismatch")
}
