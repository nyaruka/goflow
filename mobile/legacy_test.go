package mobile_test

import (
	"testing"

	"github.com/nyaruka/goflow/mobile"
	"github.com/nyaruka/goflow/test"

	"github.com/stretchr/testify/assert"
)

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
