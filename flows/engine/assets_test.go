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

	label := sa.Labels().Get(assets.LabelUUID("18644b27-fb7f-40e1-b8f4-4ea8999129ef"))
	assert.Equal(t, assets.LabelUUID("18644b27-fb7f-40e1-b8f4-4ea8999129ef"), label.UUID())
	assert.Equal(t, "Spam", label.Name())

	assert.Nil(t, sa.Labels().Get(assets.LabelUUID("xyz")))

	group := sa.Groups().Get(assets.GroupUUID("2aad21f6-30b7-42c5-bd7f-1b720c154817"))
	assert.Equal(t, assets.GroupUUID("2aad21f6-30b7-42c5-bd7f-1b720c154817"), group.UUID())
	assert.Equal(t, "Survey Audience", group.Name())

	assert.Nil(t, sa.Groups().Get(assets.GroupUUID("xyz")))

	resthook := sa.Resthooks().FindBySlug("new-registration")
	assert.Equal(t, "new-registration", resthook.Slug())
	assert.Equal(t, []string{"http://temba.io/"}, resthook.Subscribers())

	assert.Nil(t, sa.Resthooks().FindBySlug("xyz"))
}

func TestSessionAssetsWithSourceErrors(t *testing.T) {
	env := envs.NewBuilder().Build()

	source := &testSource{}

	sa, err := engine.NewSessionAssets(env, source, nil)
	require.NoError(t, err)

	source.currentErrType = "flow"
	_, err = sa.Flows().Get(assets.FlowUUID("ddba5842-252f-4a20-b901-08696fc773e2"))
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

func (s *testSource) Flow(assets.FlowUUID) (assets.Flow, error) {
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

func (s *testSource) Users() ([]assets.User, error) {
	return nil, s.err("users")
}
