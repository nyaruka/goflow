package engine_test

import (
	"testing"
	"time"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/core"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows/definition"
	"github.com/nyaruka/goflow/flows/engine"
	"github.com/nyaruka/goflow/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStep(t *testing.T) {
	node := definition.NewNode(core.NodeUUID("5fb4f555-7662-4c4c-8387-226e359526e4"), nil, nil, nil)

	d := time.Date(2018, 10, 26, 14, 50, 30, 1234567890, time.UTC)
	step := engine.NewStep(nil, node, d)

	assert.Equal(t, core.NodeUUID("5fb4f555-7662-4c4c-8387-226e359526e4"), step.NodeUUID())
	assert.Equal(t, d, step.ArrivedOn())

	// test use in expressions
	env := envs.NewBuilder().Build()
	test.AssertXEqual(t, types.NewXObject(map[string]types.XValue{
		"arrived_on": types.NewXDateTime(d),
		"node_uuid":  types.NewXText("5fb4f555-7662-4c4c-8387-226e359526e4"),
	}), core.Context(env, step))

	// test marshaling
	marshaled, err := jsonx.Marshal(step)
	require.NoError(t, err)
	test.AssertEqualJSON(t, []byte(`{"arrived_on":"2018-10-26T14:50:31.23456789Z","node_uuid":"5fb4f555-7662-4c4c-8387-226e359526e4"}`), marshaled, "JSON mismatch")

	// test unmarshaling a step from an older session which has legacy uuid and exit_uuid fields
	step = engine.NewStep(nil, node, time.Now())
	err = jsonx.Unmarshal([]byte(`{"uuid":"7a4c2ea9-f47f-4a02-b83b-56b6c1e2bb27","node_uuid":"faf2b807-9871-4785-9b02-d6d2d7a91998","exit_uuid":"9c2c9dc6-2b02-4b83-8bd0-e08a91f89464","arrived_on":"2018-10-26T14:50:31.23456789Z"}`), step)
	require.NoError(t, err)

	assert.Equal(t, core.NodeUUID("faf2b807-9871-4785-9b02-d6d2d7a91998"), step.NodeUUID())
	assert.Equal(t, time.Date(2018, 10, 26, 14, 50, 31, 234567890, time.UTC), step.ArrivedOn())
}
