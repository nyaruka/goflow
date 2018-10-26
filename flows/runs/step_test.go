package runs_test

import (
	"testing"
	"time"

	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/definition"
	"github.com/nyaruka/goflow/flows/runs"
	"github.com/nyaruka/goflow/utils"

	"github.com/stretchr/testify/assert"
)

func TestStep(t *testing.T) {
	utils.SetUUIDGenerator(utils.NewSeededUUID4Generator(1234))
	defer utils.SetUUIDGenerator(utils.DefaultUUIDGenerator)

	node := definition.NewNode(flows.NodeUUID("5fb4f555-7662-4c4c-8387-226e359526e4"), nil, nil, nil, nil)

	d := time.Date(2018, 10, 26, 14, 50, 30, 1234567890, time.UTC)
	step := runs.NewStep(node, d)

	assert.Equal(t, flows.StepUUID("c00e5d67-c275-4389-aded-7d8b151cbd5b"), step.UUID())
	assert.Equal(t, flows.NodeUUID("5fb4f555-7662-4c4c-8387-226e359526e4"), step.NodeUUID())
	assert.Equal(t, d, step.ArrivedOn())
	assert.Equal(t, flows.ExitUUID(""), step.ExitUUID())

	// test use in expressions
	env := utils.NewDefaultEnvironment()
	assert.Equal(t, "step", step.Describe())
	assert.Equal(t, types.NewXText("c00e5d67-c275-4389-aded-7d8b151cbd5b"), step.Resolve(env, "UUID"))
	assert.Equal(t, types.NewXDateTime(d), step.Resolve(env, "Arrived_On"))
	assert.Equal(t, types.NewXText("c00e5d67-c275-4389-aded-7d8b151cbd5b"), step.Reduce(env))
	assert.Equal(t, types.NewXText(`{"arrived_on":"2018-10-26T14:50:31.234567Z","exit_uuid":"","node_uuid":"5fb4f555-7662-4c4c-8387-226e359526e4","uuid":"c00e5d67-c275-4389-aded-7d8b151cbd5b"}`), step.ToXJSON(env))
}
