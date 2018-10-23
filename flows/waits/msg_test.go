package waits_test

import (
	"encoding/json"
	"testing"

	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/triggers"
	"github.com/nyaruka/goflow/flows/waits"
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
	wait := waits.NewMsgWait(nil, "")
	marshaled, _ := json.Marshal(wait)
	assert.Equal(t, `{"type":"msg"}`, string(marshaled))

	// timeout and image media hint
	timeout := 5
	wait = waits.NewMsgWait(&timeout, waits.MediaTypeImage)
	marshaled, _ = json.Marshal(wait)
	assert.Equal(t, `{"type":"msg","timeout":5,"media_hint":"image"}`, string(marshaled))
}

func TestMsgWaitSkipIfInitial(t *testing.T) {
	env := utils.NewDefaultEnvironment()
	contact := flows.NewEmptyContact("Ben Haggerty", utils.Language("eng"), nil)
	session, flow := initializeSession(t)

	// a manual trigger will wait at the initial wait
	trigger := triggers.NewManualTrigger(env, contact, flow.Reference(), nil, utils.Now())

	newEvents, err := session.Start(trigger)
	require.NoError(t, err)

	assert.Equal(t, flows.SessionStatusWaiting, session.Status())
	assert.Equal(t, 1, len(newEvents))
	assert.Equal(t, "msg_wait", newEvents[0].Type())

	session, flow = initializeSession(t)

	// whereas a msg trigger will skip over it
	msg := flows.NewMsgIn(flows.MsgUUID(utils.NewUUID()), flows.NilMsgID, urns.NilURN, nil, "Hi there", nil)
	trigger = triggers.NewMsgTrigger(env, contact, flow.Reference(), msg, nil, utils.Now())

	newEvents, err = session.Start(trigger)
	require.NoError(t, err)

	assert.Equal(t, flows.SessionStatusCompleted, session.Status())
	assert.Equal(t, 1, len(newEvents))
	assert.Equal(t, "msg_received", newEvents[0].Type())
}

func initializeSession(t *testing.T) (flows.Session, flows.Flow) {
	session, err := test.CreateSession([]byte(initialWaitJSON), "")
	require.NoError(t, err)

	flow, err := session.Assets().Flows().Get("615b8a0f-588c-4d20-a05f-363b0b4ce6f4")
	require.NoError(t, err)

	return session, flow
}
