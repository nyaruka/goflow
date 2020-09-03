package runs_test

import (
	"testing"
	"time"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/definition"
	"github.com/nyaruka/goflow/flows/runs"
	"github.com/nyaruka/goflow/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStep(t *testing.T) {
	uuids.SetGenerator(uuids.NewSeededGenerator(1234))
	defer uuids.SetGenerator(uuids.DefaultGenerator)

	node := definition.NewNode(flows.NodeUUID("5fb4f555-7662-4c4c-8387-226e359526e4"), nil, nil, nil)

	d := time.Date(2018, 10, 26, 14, 50, 30, 1234567890, time.UTC)
	step := runs.NewStep(node, d)

	assert.Equal(t, flows.StepUUID("c00e5d67-c275-4389-aded-7d8b151cbd5b"), step.UUID())
	assert.Equal(t, flows.NodeUUID("5fb4f555-7662-4c4c-8387-226e359526e4"), step.NodeUUID())
	assert.Equal(t, d, step.ArrivedOn())
	assert.Equal(t, flows.ExitUUID(""), step.ExitUUID())

	// test use in expressions
	env := envs.NewBuilder().Build()
	test.AssertXEqual(t, types.NewXObject(map[string]types.XValue{
		"arrived_on": types.NewXDateTime(d),
		"exit_uuid":  types.XTextEmpty,
		"node_uuid":  types.NewXText("5fb4f555-7662-4c4c-8387-226e359526e4"),
		"uuid":       types.NewXText("c00e5d67-c275-4389-aded-7d8b151cbd5b"),
	}), flows.Context(env, step))

	// test marshaling
	marshaled, err := jsonx.Marshal(step)
	require.NoError(t, err)
	test.AssertEqualJSON(t, []byte(`{"arrived_on":"2018-10-26T14:50:31.23456789Z","node_uuid":"5fb4f555-7662-4c4c-8387-226e359526e4","uuid":"c00e5d67-c275-4389-aded-7d8b151cbd5b"}`), marshaled, "JSON mismatch")
}
