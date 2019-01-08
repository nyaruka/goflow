package waits_test

import (
	"encoding/json"
	"testing"

	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/triggers"
	"github.com/nyaruka/goflow/flows/waits"
	"github.com/nyaruka/goflow/flows/waits/hints"
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
                    "wait": {
                        "type": "msg"
                    },
                    "exits": [
                        {
                            "uuid": "598ae7a5-2f81-48f1-afac-595262514aa1",
                            "name": "All Responses"
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
	timeout := 5
	wait = waits.NewMsgWait(&timeout, hints.NewImageHint())
	marshaled, _ = json.Marshal(wait)
	assert.Equal(t, `{"type":"msg","timeout":5,"hint":{"type":"image"}}`, string(marshaled))
}

func TestMsgWaitSkipIfInitial(t *testing.T) {
	env := utils.NewDefaultEnvironment()
	contact := flows.NewEmptyContact("Ben Haggerty", utils.Language("eng"), nil)
	session, flow := initializeSession(t)

	// a manual trigger will wait at the initial wait
	trigger := triggers.NewManualTrigger(env, flow.Reference(), contact, nil, nil, utils.Now())

	sprint, err := session.Start(trigger)
	require.NoError(t, err)

	assert.Equal(t, flows.SessionStatusWaiting, session.Status())
	assert.Equal(t, 1, len(sprint.Events()))
	assert.Equal(t, "msg_wait", sprint.Events()[0].Type())

	session, flow = initializeSession(t)

	// whereas a msg trigger will skip over it
	msg := flows.NewMsgIn(flows.MsgUUID(utils.NewUUID()), urns.NilURN, nil, "Hi there", nil)
	trigger = triggers.NewMsgTrigger(env, flow.Reference(), contact, msg, nil, utils.Now())

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
