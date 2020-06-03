package waits_test

import (
	"testing"

	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/routers/waits"
	"github.com/nyaruka/goflow/flows/routers/waits/hints"
	"github.com/nyaruka/goflow/flows/triggers"
	"github.com/nyaruka/goflow/test"
	"github.com/nyaruka/goflow/utils/jsonx"
	"github.com/nyaruka/goflow/utils/uuids"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var initialWaitJSON = `{
	"flows": [
		{
            "uuid": "615b8a0f-588c-4d20-a05f-363b0b4ce6f4",
			"name": "Initial Wait",
			"spec_version": "13.0",
            "language": "eng",
            "type": "messaging",
            "nodes": [
                {
					"uuid": "46d51f50-58de-49da-8d13-dadbf322685d",
					"router": {
						"type": "switch",
						"wait": {
							"type": "msg"
						},
						"categories": [
							{
								"uuid": "c82e161f-fa2d-4e7d-a338-c27f6c349445",
								"name": "All Responses",
								"exit_uuid": "598ae7a5-2f81-48f1-afac-595262514aa1"
							}
						],
						"operand": "@input.text",
						"default_category_uuid": "c82e161f-fa2d-4e7d-a338-c27f6c349445"
					},
                    "exits": [
                        {
                            "uuid": "598ae7a5-2f81-48f1-afac-595262514aa1"
                        }
                    ]
                }
            ]
        }
	]
}`

func TestMsgWait(t *testing.T) {
	// no timeout or media
	wait := waits.NewMsgWait(nil, nil)
	marshaled, _ := jsonx.Marshal(wait)
	assert.Equal(t, `{"type":"msg"}`, string(marshaled))

	// timeout and image hint
	wait = waits.NewMsgWait(
		waits.NewTimeout(5, flows.CategoryUUID("63fca57d-5ef6-4afd-9bcd-7bdcf653cea8")),
		hints.NewImageHint(),
	)
	marshaled, err := jsonx.Marshal(wait)
	require.NoError(t, err)
	assert.Equal(t, `{"type":"msg","timeout":{"seconds":5,"category_uuid":"63fca57d-5ef6-4afd-9bcd-7bdcf653cea8"},"hint":{"type":"image"}}`, string(marshaled))
}

func TestMsgWaitSkipIfInitial(t *testing.T) {
	eng := test.NewEngine()
	env := envs.NewBuilder().Build()
	sa, flow := initializeSessionAssets(t)
	contact := flows.NewEmptyContact(sa, "Ben Haggerty", envs.Language("eng"), nil)

	// a manual trigger will wait at the initial wait
	trigger := triggers.NewManual(env, flow.Reference(), contact, false, nil)

	session, sprint, err := eng.NewSession(sa, trigger)
	require.NoError(t, err)

	assert.Equal(t, flows.SessionStatusWaiting, session.Status())
	assert.Equal(t, 1, len(sprint.Events()))
	assert.Equal(t, "msg_wait", sprint.Events()[0].Type())

	sa, flow = initializeSessionAssets(t)

	// whereas a msg trigger will skip over it
	msg := flows.NewMsgIn(flows.MsgUUID(uuids.New()), urns.NilURN, nil, "Hi there", nil)
	trigger = triggers.NewMsg(env, flow.Reference(), contact, msg, nil)

	session, sprint, err = eng.NewSession(sa, trigger)
	require.NoError(t, err)

	assert.Equal(t, flows.SessionStatusCompleted, session.Status())
	assert.Equal(t, 1, len(sprint.Events()))
	assert.Equal(t, "msg_received", sprint.Events()[0].Type())
}

func initializeSessionAssets(t *testing.T) (flows.SessionAssets, flows.Flow) {
	sa, err := test.CreateSessionAssets([]byte(initialWaitJSON), "")
	require.NoError(t, err)

	flow, err := sa.Flows().Get("615b8a0f-588c-4d20-a05f-363b0b4ce6f4")
	require.NoError(t, err)

	return sa, flow
}
