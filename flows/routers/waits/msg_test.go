package waits_test

import (
	"encoding/json"
	"testing"

	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/routers/waits"
	"github.com/nyaruka/goflow/flows/routers/waits/hints"
	"github.com/nyaruka/goflow/flows/triggers"
	"github.com/nyaruka/goflow/test"
	"github.com/nyaruka/goflow/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var initialWaitJSON = `{
	"flows": [
		{
            "uuid": "615b8a0f-588c-4d20-a05f-363b0b4ce6f4",
			"name": "Initial Wait",
			"spec_version": "12.0",
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
	marshaled, _ := json.Marshal(wait)
	assert.Equal(t, `{"type":"msg"}`, string(marshaled))

	// timeout and image hint
	wait = waits.NewMsgWait(
		waits.NewTimeout(5, flows.CategoryUUID("63fca57d-5ef6-4afd-9bcd-7bdcf653cea8")),
		hints.NewImageHint(),
	)
	marshaled, err := json.Marshal(wait)
	require.NoError(t, err)
	assert.Equal(t, `{"type":"msg","timeout":{"seconds":5,"category_uuid":"63fca57d-5ef6-4afd-9bcd-7bdcf653cea8"},"hint":{"type":"image"}}`, string(marshaled))
}

func TestMsgWaitSkipIfInitial(t *testing.T) {
	env := utils.NewEnvironmentBuilder().Build()
	session, flow := initializeSession(t)
	contact := flows.NewEmptyContact(session.Assets(), "Ben Haggerty", utils.Language("eng"), nil)

	// a manual trigger will wait at the initial wait
	trigger := triggers.NewManualTrigger(env, flow.Reference(), contact, nil)

	sprint, err := session.Start(trigger)
	require.NoError(t, err)

	assert.Equal(t, flows.SessionStatusWaiting, session.Status())
	assert.Equal(t, 1, len(sprint.Events()))
	assert.Equal(t, "msg_wait", sprint.Events()[0].Type())

	session, flow = initializeSession(t)

	// whereas a msg trigger will skip over it
	msg := flows.NewMsgIn(flows.MsgUUID(utils.NewUUID()), urns.NilURN, nil, "Hi there", nil)
	trigger = triggers.NewMsgTrigger(env, flow.Reference(), contact, msg, nil)

	sprint, err = session.Start(trigger)
	require.NoError(t, err)

	assert.Equal(t, flows.SessionStatusCompleted, session.Status())
	assert.Equal(t, 1, len(sprint.Events()))
	assert.Equal(t, "msg_received", sprint.Events()[0].Type())
}

func initializeSession(t *testing.T) (flows.Session, flows.Flow) {
	session, err := test.CreateSession([]byte(initialWaitJSON), "")
	require.NoError(t, err)

	flow, err := session.Assets().Flows().Get("615b8a0f-588c-4d20-a05f-363b0b4ce6f4")
	require.NoError(t, err)

	return session, flow
}
