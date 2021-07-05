package waits_test

import (
	"testing"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/resumes"
	"github.com/nyaruka/goflow/flows/routers/waits"
	"github.com/nyaruka/goflow/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDialWait(t *testing.T) {
	session, _, err := test.CreateTestVoiceSession("")
	require.NoError(t, err)
	run := session.Runs()[0]

	// phone field required
	_, err = waits.ReadWait([]byte(`{"type": "dial"}`))
	assert.EqualError(t, err, "field 'phone' is required")

	wait, err := waits.ReadWait([]byte(`{"type": "dial", "phone": "@(\"+\" & \"593979123456\")"}`))
	assert.NoError(t, err)
	assert.Equal(t, waits.TypeDial, wait.Type())

	// test marsalling definition wait
	marshaled, err := jsonx.Marshal(wait)
	assert.NoError(t, err)
	assert.Equal(t, `{"type":"dial","phone":"@(\"+\" & \"593979123456\")"}`, string(marshaled))

	// try activating the wait
	log := test.NewEventLog()
	activated := wait.Begin(run, log.Log)

	assert.Equal(t, "dial", activated.Type())
	assert.Equal(t, 1, len(log.Events))
	assert.Equal(t, "dial_wait", log.Events[0].Type())

	// test marsalling activated wait
	marshaled, err = jsonx.Marshal(activated)
	assert.NoError(t, err)
	assert.Equal(t, `{"type":"dial","urn":"tel:+593979123456"}`, string(marshaled))

	// try to end with incorrect resume type
	err = wait.End(resumes.NewWaitTimeout(nil, nil))
	assert.EqualError(t, err, "can't end a wait of type 'dial' with a resume of type 'wait_timeout'")

	// try to end with dial resume type
	err = wait.End(resumes.NewDial(nil, nil, flows.NewDial(flows.DialStatusAnswered, 5)))
	assert.NoError(t, err)

	// try when wait has expression error but still generates valid tel URN
	wait, err = waits.ReadWait([]byte(`{"type": "dial", "phone": "+593979123456@(1 / 0)"}`))
	assert.NoError(t, err)

	log = test.NewEventLog()
	activated = wait.Begin(run, log.Log)

	assert.Equal(t, "dial", activated.Type())
	assert.Equal(t, urns.URN("tel:+593979123456"), activated.(*waits.ActivatedDialWait).URN())
	assert.Equal(t, 2, len(log.Events))
	assert.Equal(t, "error", log.Events[0].Type())
	assert.Equal(t, "dial_wait", log.Events[1].Type())

	// try when wait doesn't generate a valid tel URN
	wait, err = waits.ReadWait([]byte(`{"type": "dial", "phone": "@(\"\")"}`))
	assert.NoError(t, err)

	log = test.NewEventLog()
	activated = wait.Begin(run, log.Log)

	assert.Nil(t, activated)
	assert.Equal(t, 1, len(log.Events))
	assert.Equal(t, "error", log.Events[0].Type())
}
