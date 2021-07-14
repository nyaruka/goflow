package migrations_test

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"testing"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/definition"
	"github.com/nyaruka/goflow/flows/definition/migrations"
	"github.com/nyaruka/goflow/test"

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
		testsJSON, err := os.ReadFile(fmt.Sprintf("testdata/migrations/%s.json", version.String()))
		require.NoError(t, err)

		tests := []struct {
			Description string          `json:"description"`
			Original    json.RawMessage `json:"original"`
			Migrated    json.RawMessage `json:"migrated"`
		}{}

		err = jsonx.Unmarshal(testsJSON, &tests)
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

	_, err = migrations.MigrateToLatest([]byte(`{}`), nil)
	assert.EqualError(t, err, "unable to read flow header: field 'uuid' is required, field 'spec_version' is required")

	migrated, err = migrations.MigrateToLatest([]byte(`{
		"uuid": "76f0a02f-3b75-4b86-9064-e9195e1b3a02",
		"spec_version": "13.0"
	}`), nil)
	require.NoError(t, err)

	expected := fmt.Sprintf(`{
		"uuid": "76f0a02f-3b75-4b86-9064-e9195e1b3a02",
		"spec_version": "%s"
	}`, definition.CurrentSpecVersion)
	test.AssertEqualJSON(t, []byte(expected), migrated, "flow migration mismatch")

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

	expected = fmt.Sprintf(`{
		"uuid": "76f0a02f-3b75-4b86-9064-e9195e1b3a02",
		"name": "Empty Flow",
		"spec_version": "%s",
		"language": "eng",
		"type": "messaging",
		"nodes": []
	}`, definition.CurrentSpecVersion)
	test.AssertEqualJSON(t, []byte(expected), migrated, "flow migration mismatch")

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

	expected = fmt.Sprintf(`{
		"uuid": "50c3706e-fedb-42c0-8eab-dda3335714b7",
		"name": "Test Flow",
		"spec_version": "%s",
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
	}`, definition.CurrentSpecVersion)
	test.AssertEqualJSON(t, []byte(expected), migrated, "flow migration mismatch")

	// try to migrate legacy definition without migration config
	_, err = migrations.MigrateToLatest([]byte(`{
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

func TestClone(t *testing.T) {
	env := envs.NewBuilder().Build()

	testCases := []struct {
		path string
		uuid string
	}{
		{"testdata/clone_with_ui.json", "ee765ff2-96b0-440a-b108-393f613466bb"},
		{"../../../test/testdata/runner/two_questions.json", "615b8a0f-588c-4d20-a05f-363b0b4ce6f4"},
		{"../../../test/testdata/runner/all_actions.json", "8ca44c09-791d-453a-9799-a70dd3303306"},
		{"../../../test/testdata/runner/router_tests.json", "615b8a0f-588c-4d20-a05f-363b0b4ce6f4"},
	}

	for _, tc := range testCases {
		uuids.SetGenerator(uuids.NewSeededGenerator(12345))
		defer uuids.SetGenerator(uuids.DefaultGenerator)

		flow, err := test.LoadFlowFromAssets(env, tc.path, assets.FlowUUID(tc.uuid))
		require.NoError(t, err)

		depMappings := map[uuids.UUID]uuids.UUID{
			uuids.UUID(tc.uuid):                    "e0af9907-e0d3-4363-99c6-324ece7f628e", // the flow itself
			"2aad21f6-30b7-42c5-bd7f-1b720c154817": "cd8a68c0-6673-4a02-98a0-7fb3ac788860", // group used in has_group test
		}

		flowJSON, err := jsonx.Marshal(flow)
		require.NoError(t, err)

		cloneJSON, err := migrations.Clone(flowJSON, depMappings)
		require.NoError(t, err)

		clone, err := definition.ReadFlow(cloneJSON, nil)
		require.NoError(t, err)

		assert.Equal(t, assets.FlowUUID("e0af9907-e0d3-4363-99c6-324ece7f628e"), clone.UUID())
		assert.Equal(t, flow.Name(), clone.Name())
		assert.Equal(t, flow.Type(), clone.Type())
		assert.Equal(t, flow.Revision(), clone.Revision())
		assert.Equal(t, len(flow.Nodes()), len(clone.Nodes()))

		// extract all UUIDs from originaland cloned definitions
		originalUUIDs := uuids.V4Regex.FindAllString(string(flowJSON), -1)
		cloneUUIDs := uuids.V4Regex.FindAllString(string(cloneJSON), -1)

		assert.Equal(t, len(originalUUIDs), len(cloneUUIDs))
		assert.NotContains(t, cloneUUIDs, []string{"2aad21f6-30b7-42c5-bd7f-1b720c154817"}) // group used in has_group test

		for _, u1 := range originalUUIDs {
			for _, u2 := range cloneUUIDs {
				if u1 == u2 && depMappings[uuids.UUID(u1)] != "" {
					assert.Fail(t, "uuid", "cloned flow contains non-dependency UUID from original flow: %s", u1)
				}
			}
		}

		// if flow has a UI section, check UI node UUIDs correspond to real nodes
		if len(clone.UI()) > 0 {
			clonedUI, err := jsonx.DecodeGeneric(clone.UI())
			require.NoError(t, err)

			nodeMap := clonedUI.(map[string]interface{})["nodes"].(map[string]interface{})

			for nodeUUID := range nodeMap {
				assert.NotNil(t, clone.GetNode(flows.NodeUUID(nodeUUID)), "UI has node with UUID %s that doesn't exist in flow", nodeUUID)
			}
		}
	}
}

func TestCloneOlderVersion(t *testing.T) {
	uuids.SetGenerator(uuids.NewSeededGenerator(12345))
	defer uuids.SetGenerator(uuids.DefaultGenerator)

	cloneJSON, err := migrations.Clone([]byte(`{
		"uuid": "ee765ff2-96b0-440a-b108-393f613466bb",
		"name": "Older Flow",
		"spec_version": "13.0.0",
		"language": "base",
		"revision": 11,
		"expire_after_minutes": 10080,
		"type": "messaging",
		"nodes": []
	}`), nil)
	require.NoError(t, err)

	// cloned flow should have same spec version but different UUID
	test.AssertEqualJSON(t, []byte(`{
		"uuid": "1ae96956-4b34-433e-8d1a-f05fe6923d6d",
		"name": "Older Flow",
		"spec_version": "13.0.0",
		"language": "base",
		"revision": 11,
		"expire_after_minutes": 10080,
		"type": "messaging",
		"nodes": []
	}`), cloneJSON, "cloned flow mismatch")
}
