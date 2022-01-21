package engine_test

import (
	"fmt"
	"testing"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/assets/static"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows/engine"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var assetsJSON = `{
	"flows": [
		{
            "uuid": "76f0a02f-3b75-4b86-9064-e9195e1b3a02",
            "name": "Empty",
            "spec_version": "13.0",
            "language": "eng",
            "type": "messaging",
            "nodes": []
        }
	],
	"fields": [
        {"uuid": "d66a7823-eada-40e5-9a3a-57239d4690bf", "key": "gender", "name": "Gender", "type": "text"},
        {"uuid": "f1b5aea6-6586-41c7-9020-1a6326cc6565", "key": "age", "name": "Age", "type": "number"}
    ],
	"groups": [
		{
			"uuid": "2aad21f6-30b7-42c5-bd7f-1b720c154817",
			"name": "Survey Audience"
		}
	],
	"labels": [
		{
			"uuid": "18644b27-fb7f-40e1-b8f4-4ea8999129ef",
			"name": "Spam"
		}
	],
	"resthooks": [
		{
			"slug": "new-registration",
			"subscribers": [
				"http://temba.io/"
			]
		}
	]
}`

func TestSessionAssets(t *testing.T) {
	env := envs.NewBuilder().Build()

	source, err := static.NewSource([]byte(assetsJSON))
	require.NoError(t, err)

	sa, err := engine.NewSessionAssets(env, source, nil)
	require.NoError(t, err)

	assert.Equal(t, source, sa.Source())

	label := sa.Labels().Get("18644b27-fb7f-40e1-b8f4-4ea8999129ef")
	assert.Equal(t, assets.LabelUUID("18644b27-fb7f-40e1-b8f4-4ea8999129ef"), label.UUID())
	assert.Equal(t, "Spam", label.Name())

	assert.Nil(t, sa.Labels().Get("xyz"))

	group := sa.Groups().Get("2aad21f6-30b7-42c5-bd7f-1b720c154817")
	assert.Equal(t, assets.GroupUUID("2aad21f6-30b7-42c5-bd7f-1b720c154817"), group.UUID())
	assert.Equal(t, "Survey Audience", group.Name())

	assert.Nil(t, sa.Groups().Get("xyz"))

	resthook := sa.Resthooks().FindBySlug("new-registration")
	assert.Equal(t, "new-registration", resthook.Slug())
	assert.Equal(t, []string{"http://temba.io/"}, resthook.Subscribers())

	assert.Nil(t, sa.Resthooks().FindBySlug("xyz"))

	// sessions assets are used as a contactql resolver for query parsing
	age := sa.ResolveField("age")
	assert.Equal(t, assets.FieldUUID(`f1b5aea6-6586-41c7-9020-1a6326cc6565`), age.UUID())
	assert.Equal(t, "age", age.Key())
	assert.Equal(t, "Age", age.Name())
	assert.Equal(t, assets.FieldTypeNumber, age.Type())

	audience := sa.ResolveGroup("survey audience")
	assert.Equal(t, assets.GroupUUID(`2aad21f6-30b7-42c5-bd7f-1b720c154817`), audience.UUID())
	assert.Equal(t, "Survey Audience", audience.Name())

	emptyFlow := sa.ResolveFlow("empty")
	assert.Equal(t, assets.FlowUUID(`76f0a02f-3b75-4b86-9064-e9195e1b3a02`), emptyFlow.UUID())
	assert.Equal(t, "Empty", emptyFlow.Name())

	assert.Nil(t, sa.ResolveField("xxx"))
	assert.Nil(t, sa.ResolveGroup("xxx"))
	assert.Nil(t, sa.ResolveFlow("xxx"))
}

func TestSessionAssetsWithSourceErrors(t *testing.T) {
	env := envs.NewBuilder().Build()

	source := &testSource{}

	sa, err := engine.NewSessionAssets(env, source, nil)
	require.NoError(t, err)

	source.currentErrType = "flow"
	_, err = sa.Flows().Get(assets.FlowUUID("ddba5842-252f-4a20-b901-08696fc773e2"))
	assert.EqualError(t, err, "unable to load flow assets")

	source.currentErrType = "flow"
	_, err = sa.Flows().FindByName("Catch All")
	assert.EqualError(t, err, "unable to load flow assets")

	for _, errType := range []string{"channels", "classifiers", "fields", "globals", "groups", "labels", "locations", "resthooks", "templates", "users"} {
		source.currentErrType = errType
		_, err = engine.NewSessionAssets(env, source, nil)
		assert.EqualError(t, err, fmt.Sprintf("unable to load %s assets", errType), "error mismatch for type %s", errType)
	}
}

// a source for testing which will return an err when requested an asset of currentErrType
type testSource struct {
	currentErrType string
}

func (s *testSource) err(t string) error {
	if t == s.currentErrType {
		return errors.Errorf("unable to load %s assets", t)
	}
	return nil
}

func (s *testSource) Channels() ([]assets.Channel, error) {
	return nil, s.err("channels")
}

func (s *testSource) Classifiers() ([]assets.Classifier, error) {
	return nil, s.err("classifiers")
}

func (s *testSource) Fields() ([]assets.Field, error) {
	return nil, s.err("fields")
}

func (s *testSource) FlowByUUID(assets.FlowUUID) (assets.Flow, error) {
	return nil, s.err("flow")
}

func (s *testSource) FlowByName(name string) (assets.Flow, error) {
	return nil, s.err("flow")
}

func (s *testSource) Globals() ([]assets.Global, error) {
	return nil, s.err("globals")
}

func (s *testSource) Groups() ([]assets.Group, error) {
	return nil, s.err("groups")
}

func (s *testSource) Labels() ([]assets.Label, error) {
	return nil, s.err("labels")
}

func (s *testSource) Locations() ([]assets.LocationHierarchy, error) {
	return nil, s.err("locations")
}

func (s *testSource) Resthooks() ([]assets.Resthook, error) {
	return nil, s.err("resthooks")
}

func (s *testSource) Templates() ([]assets.Template, error) {
	return nil, s.err("templates")
}

func (s *testSource) Ticketers() ([]assets.Ticketer, error) {
	return nil, s.err("ticketers")
}

func (s *testSource) Topics() ([]assets.Topic, error) {
	return nil, s.err("topics")
}

func (s *testSource) Users() ([]assets.User, error) {
	return nil, s.err("users")
}
