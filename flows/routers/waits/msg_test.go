package waits_test

import (
	"testing"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/resumes"
	"github.com/nyaruka/goflow/flows/routers/waits"
	"github.com/nyaruka/goflow/flows/routers/waits/hints"
	"github.com/nyaruka/goflow/test"

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
	session, _, err := test.CreateTestVoiceSession("")
	require.NoError(t, err)
	run := session.Runs()[0]

	// no timeout or media
	wait := waits.NewMsgWait(nil, nil)
	marshaled := jsonx.MustMarshal(wait)
	assert.Equal(t, `{"type":"msg"}`, string(marshaled))

	// try to end with timeout resume type
	assert.False(t, wait.Accepts(resumes.NewWaitTimeout(nil, nil)))

	// timeout and image hint
	wait = waits.NewMsgWait(
		waits.NewTimeout(5, flows.CategoryUUID("63fca57d-5ef6-4afd-9bcd-7bdcf653cea8")),
		hints.NewImageHint(),
	)

	// test marsalling definition wait
	marshaled, err = jsonx.Marshal(wait)
	require.NoError(t, err)
	assert.Equal(t, `{"type":"msg","timeout":{"seconds":5,"category_uuid":"63fca57d-5ef6-4afd-9bcd-7bdcf653cea8"},"hint":{"type":"image"}}`, string(marshaled))

	// try activating the wait
	log := test.NewEventLog()
	begun := wait.Begin(run, log.Log)

	assert.True(t, begun)
	assert.Equal(t, 1, len(log.Events))
	assert.Equal(t, "msg_wait", log.Events[0].Type())

	// try to end with incorrect resume type
	assert.False(t, wait.Accepts(resumes.NewDial(nil, nil, flows.NewDial(flows.DialStatusBusy, 0))))

	// can end with timeout resume type
	assert.True(t, wait.Accepts(resumes.NewWaitTimeout(nil, nil)))
}

func TestMsgWaitSkipIfInitial(t *testing.T) {
	// a manual trigger will wait at the initial wait
	session, sprint := test.NewSessionBuilder().WithAssets([]byte(initialWaitJSON)).
		WithFlow("615b8a0f-588c-4d20-a05f-363b0b4ce6f4").
		MustBuild()

	assert.Equal(t, flows.SessionStatusWaiting, session.Status())
	assert.Equal(t, 1, len(sprint.Events()))
	assert.Equal(t, "msg_wait", sprint.Events()[0].Type())

	// whereas a msg trigger will skip over it
	session, sprint = test.NewSessionBuilder().WithAssets([]byte(initialWaitJSON)).
		WithFlow("615b8a0f-588c-4d20-a05f-363b0b4ce6f4").
		WithTriggerMsg("Hi there").
		MustBuild()

	assert.Equal(t, flows.SessionStatusCompleted, session.Status())
	assert.Equal(t, 1, len(sprint.Events()))
	assert.Equal(t, "msg_received", sprint.Events()[0].Type())
}
