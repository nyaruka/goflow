package migrations_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sort"
	"testing"

	"github.com/nyaruka/goflow/flows/definition"
	"github.com/nyaruka/goflow/flows/definition/migrations"
	"github.com/nyaruka/goflow/test"
	"github.com/nyaruka/goflow/utils/uuids"

	"github.com/Masterminds/semver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMigrateToVersion(t *testing.T) {
	defer uuids.SetGenerator(uuids.DefaultGenerator)

	// get all versions in order
	versions := make([]*semver.Version, 0, len(migrations.Registered()))
	for v := range migrations.Registered() {
		versions = append(versions, v)
	}
	sort.SliceStable(versions, func(i, j int) bool { return versions[i].LessThan(versions[j]) })

	for _, version := range versions {
		testsJSON, err := ioutil.ReadFile(fmt.Sprintf("testdata/%s.json", version.String()))
		require.NoError(t, err)

		tests := []struct {
			Description string          `json:"description"`
			Original    json.RawMessage `json:"original"`
			Migrated    json.RawMessage `json:"migrated"`
		}{}

		err = json.Unmarshal(testsJSON, &tests)
		require.NoError(t, err, "unable to read tests for version %s", version)

		for _, tc := range tests {
			testName := fmt.Sprintf("version %s with '%s'", version, tc.Description)

			uuids.SetGenerator(uuids.NewSeededGenerator(123456))

			actual, err := migrations.MigrateToVersion(tc.Original, version, nil)
			assert.NoError(t, err, "unexpected error in %s", testName)

			test.AssertEqualJSON(t, tc.Migrated, actual, "migration mismatch in %s", testName)

			// check final flow is valid
			_, err = definition.ReadFlow(actual, nil)
			assert.NoError(t, err, "migrated flow validation error in %s", testName)
		}
	}
}

func TestMigrateToLatest(t *testing.T) {
	defer uuids.SetGenerator(uuids.DefaultGenerator)

	migrated, err := migrations.MigrateToLatest([]byte(`[]`), nil)
	assert.EqualError(t, err, "unable to read flow header: json: cannot unmarshal array into Go value of type migrations.Header13")
	assert.Nil(t, migrated)

	migrated, err = migrations.MigrateToLatest([]byte(`{}`), nil)
	assert.EqualError(t, err, "unable to read flow header: field 'uuid' is required, field 'spec_version' is required")

	migrated, err = migrations.MigrateToLatest([]byte(`{
		"uuid": "76f0a02f-3b75-4b86-9064-e9195e1b3a02",
		"spec_version": "13.0"
	}`), nil)
	require.NoError(t, err)
	test.AssertEqualJSON(t, []byte(`{
		"uuid": "76f0a02f-3b75-4b86-9064-e9195e1b3a02",
		"spec_version": "13.1.0"
	}`), migrated, "flow migration mismatch")

	// migrate valid definition
	migrated, err = migrations.MigrateToLatest([]byte(`{
		"uuid": "76f0a02f-3b75-4b86-9064-e9195e1b3a02",
		"name": "Empty Flow",
		"spec_version": "13.0",
		"language": "eng",
		"type": "messaging",
		"nodes": []
	}`), nil)

	require.NoError(t, err)
	test.AssertEqualJSON(t, []byte(`{
		"uuid": "76f0a02f-3b75-4b86-9064-e9195e1b3a02",
		"name": "Empty Flow",
		"spec_version": "13.1.0",
		"language": "eng",
		"type": "messaging",
		"nodes": []
	}`), migrated, "flow migration mismatch")

	// migrate legacy definition
	migrated, err = migrations.MigrateToLatest([]byte(`{
		"base_language": "eng",
		"entry": "10e483a8-5ffb-4c4f-917b-d43ce86c1d65", 
		"flow_type": "M",
		"action_sets": [],
		"metadata": {
			"uuid": "50c3706e-fedb-42c0-8eab-dda3335714b7",
			"name": "Test Flow"
		}
	}`), &migrations.Config{})

	require.NoError(t, err)
	test.AssertEqualJSON(t, []byte(`{
		"uuid": "50c3706e-fedb-42c0-8eab-dda3335714b7",
		"name": "Test Flow",
		"spec_version": "13.1.0",
		"language": "eng",
		"type": "messaging",
		"expire_after_minutes": 0,
		"revision": 0,
		"localization": {},
		"nodes": [],
		"_ui": {
			"nodes": {},
        	"stickies": {}
		}
	}`), migrated, "flow migration mismatch")

	// try to migrate legacy definition without migration config
	migrated, err = migrations.MigrateToLatest([]byte(`{
		"base_language": "eng",
		"entry": "10e483a8-5ffb-4c4f-917b-d43ce86c1d65", 
		"flow_type": "M",
		"action_sets": [],
		"metadata": {
			"uuid": "50c3706e-fedb-42c0-8eab-dda3335714b7",
			"name": "Test Flow"
		}
	}`), nil)
	assert.EqualError(t, err, "unable to migrate what appears to be a legacy definition without a migration config")
}
