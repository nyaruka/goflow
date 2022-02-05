package flows_test

import (
	"testing"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/static"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/engine"
	"github.com/nyaruka/goflow/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGroupList(t *testing.T) {
	env := envs.NewBuilder().Build()

	source, err := static.NewSource([]byte(`{
		"fields": [
			{
				"key": "gender",
				"name": "Gender",
				"type": "text"
			}
		],
		"groups": [
			{
				"uuid": "e25852ea-b014-4ac1-9982-d6dcb0c2a1d5",
				"name": "Customers"
			},
			{
				"uuid": "990e1392-1f49-40c5-9662-f39609324bf9",
				"name": "Testers"
			},
			{
				"uuid": "f4f4b78e-f072-42e2-987d-f5c13da3166d",
				"name": "Males",
				"query": "gender = \"M\""
			},
			{
				"uuid": "f4f4b78e-f072-42e2-987d-f5c13da3166d",
				"name": "Broken",
				"query": "xyz = \"X\""
			}
		]
	}`))
	require.NoError(t, err)

	sa, err := engine.NewSessionAssets(env, source, nil)
	require.NoError(t, err)

	// check we ignored broken group
	assert.Equal(t, 3, len(sa.Groups().All()))

	testers := sa.Groups().Get("990e1392-1f49-40c5-9662-f39609324bf9")
	males := sa.Groups().Get("f4f4b78e-f072-42e2-987d-f5c13da3166d")

	missingRefs := make([]assets.Reference, 0)
	missing := func(ref assets.Reference, err error) {
		missingRefs = append(missingRefs, ref)
	}

	// create empty
	groups := flows.NewGroupList(sa, nil, missing)

	assert.Equal(t, 0, groups.Count())
	assert.Equal(t, 0, len(missingRefs))

	// create with some references
	groups = flows.NewGroupList(sa, []*assets.GroupReference{
		assets.NewGroupReference("990e1392-1f49-40c5-9662-f39609324bf9", "Testers"),
		assets.NewGroupReference("f4f4b78e-f072-42e2-987d-f5c13da3166d", "Males"),
		assets.NewGroupReference("7cb12d0e-e163-492c-95b1-28549cd04fe4", "I don't exist"),
	}, missing)

	assert.Equal(t, 2, groups.Count())
	assert.Equal(t, 1, len(missingRefs))
	assert.Equal(t, assets.NewGroupReference("7cb12d0e-e163-492c-95b1-28549cd04fe4", "I don't exist"), missingRefs[0])

	assert.Equal(t, males, groups.FindByUUID("f4f4b78e-f072-42e2-987d-f5c13da3166d"))
	assert.Nil(t, groups.FindByUUID("7cb12d0e-e163-492c-95b1-28549cd04fe4"))

	// check use in expressions
	test.AssertXEqual(t, types.NewXArray(testers.ToXValue(env), males.ToXValue(env)), groups.ToXValue(env))

	groups.Clear()
	assert.Equal(t, 0, groups.Count())
}
