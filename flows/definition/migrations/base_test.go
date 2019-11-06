package migrations_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sort"
	"testing"

	"github.com/Masterminds/semver"
	"github.com/nyaruka/goflow/flows/definition"
	"github.com/nyaruka/goflow/flows/definition/migrations"
	"github.com/nyaruka/goflow/test"
	"github.com/nyaruka/goflow/utils/uuids"

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

	prevVersion := semver.MustParse(`13.0.0`)

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

			actual, err := migrations.MigrateToVersion(tc.Original, prevVersion, version)
			assert.NoError(t, err, "unexpected error in %s", testName)

			test.AssertEqualJSON(t, tc.Migrated, actual, "migration mismatch in %s", testName)

			// check final flow is valid
			_, err = definition.ReadFlow(actual, nil)
			assert.NoError(t, err, "migrated flow validation error in %s", testName)
		}

		prevVersion = version
	}
}

func TestMigrateToLatest(t *testing.T) {
	defer uuids.SetGenerator(uuids.DefaultGenerator)

	migrated, err := migrations.MigrateToLatest([]byte(`[]`), semver.MustParse(`13.0.0`))
	assert.EqualError(t, err, "can't migrate definition which isn't a flow")
	assert.Nil(t, migrated)

	migrated, err = migrations.MigrateToLatest([]byte(`{}`), semver.MustParse(`13.0.0`))
	require.NoError(t, err)
	test.AssertEqualJSON(t, []byte(`{
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
	}`), semver.MustParse(`13.0.0`))

	require.NoError(t, err)
	test.AssertEqualJSON(t, []byte(`{
		"uuid": "76f0a02f-3b75-4b86-9064-e9195e1b3a02",
		"name": "Empty Flow",
		"spec_version": "13.1.0",
		"language": "eng",
		"type": "messaging",
		"nodes": []
	}`), migrated, "flow migration mismatch")
}

func TestMigrationPrimitives(t *testing.T) {
	f := migrations.Flow{map[string]interface{}{}} // nodes not set
	assert.Equal(t, []migrations.Node{}, f.Nodes())

	f = migrations.Flow{map[string]interface{}{"nodes": nil}} // nodes is nil
	assert.Equal(t, []migrations.Node{}, f.Nodes())

	f = migrations.Flow{map[string]interface{}{"nodes": []interface{}{}}} // nodes is empty
	assert.Equal(t, []migrations.Node{}, f.Nodes())

	f = migrations.Flow{map[string]interface{}{"nodes": []interface{}{
		map[string]interface{}{},
	}}}
	assert.Equal(t, []migrations.Node{migrations.Node{map[string]interface{}{}}}, f.Nodes())

	n := migrations.Node{map[string]interface{}{}} // actions and router are not set
	assert.Equal(t, []migrations.Action{}, n.Actions())
	assert.Nil(t, n.Router())

	n = migrations.Node{map[string]interface{}{"actions": nil, "router": nil}} // actions and router are nil
	assert.Equal(t, []migrations.Action{}, n.Actions())
	assert.Nil(t, n.Router())

	n = migrations.Node{map[string]interface{}{
		"actions": []interface{}{},
		"router":  map[string]interface{}{},
	}}
	assert.Equal(t, []migrations.Action{}, n.Actions())
	assert.Equal(t, &migrations.Router{map[string]interface{}{}}, n.Router())

	a := migrations.Action{map[string]interface{}{}} // type not set
	assert.Equal(t, "", a.Type())

	a = migrations.Action{map[string]interface{}{"type": "foo"}} // type set
	assert.Equal(t, "foo", a.Type())
}
