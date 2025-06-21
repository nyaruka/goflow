package waits_test

import (
	"testing"
	"time"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
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

	// time limits will default if not provided
	wait, err := waits.ReadWait([]byte(`{"type": "dial", "phone": "@(\"+\" & \"593979123456\")"}`))
	assert.NoError(t, err)
	assert.Equal(t, waits.TypeDial, wait.Type())
	assert.Equal(t, 60*time.Second, wait.(*waits.DialWait).DialLimit())
	assert.Equal(t, 2*time.Hour, wait.(*waits.DialWait).CallLimit())

	// or can be provided explicitly
	wait, err = waits.ReadWait([]byte(`{"type": "dial", "phone": "@(\"+\" & \"593979123456\")", "dial_limit_seconds": 10, "call_limit_seconds": 120}`))
	assert.NoError(t, err)
	assert.Equal(t, waits.TypeDial, wait.Type())
	assert.Equal(t, 10*time.Second, wait.(*waits.DialWait).DialLimit())
	assert.Equal(t, 120*time.Second, wait.(*waits.DialWait).CallLimit())

	// test marsalling definition wait
	marshaled, err := jsonx.Marshal(wait)
	assert.NoError(t, err)
	assert.Equal(t, `{"type":"dial","phone":"@(\"+\" & \"593979123456\")","dial_limit_seconds":10,"call_limit_seconds":120}`, string(marshaled))

	// try activating the wait
	log := test.NewEventLog()
	begun := wait.Begin(run, log.Log)

	assert.True(t, begun)
	assert.Equal(t, 1, len(log.Events))
	assert.Equal(t, "dial_wait", log.Events[0].Type())

	// try to end with incorrect resume type
	assert.False(t, wait.Accepts(resumes.NewWaitTimeout()))

	// try to end with dial resume type
	assert.True(t, wait.Accepts(resumes.NewDial(flows.NewDial(flows.DialStatusAnswered, 5))))

	// try when wait has expression error but still generates valid tel URN
	wait, err = waits.ReadWait([]byte(`{"type": "dial", "phone": "+593979123456@(1 / 0)", "dial_limit_seconds": 10, "call_limit_seconds": 120}`))
	assert.NoError(t, err)

	log = test.NewEventLog()
	begun = wait.Begin(run, log.Log)

	assert.True(t, begun)
	assert.Equal(t, 2, len(log.Events))
	assert.Equal(t, "error", log.Events[0].Type())
	assert.Equal(t, "dial_wait", log.Events[1].Type())
	assert.Equal(t, urns.URN("tel:+593979123456"), log.Events[1].(*events.DialWaitEvent).URN)
	assert.Equal(t, 10, log.Events[1].(*events.DialWaitEvent).DialLimitSeconds)
	assert.Equal(t, 120, log.Events[1].(*events.DialWaitEvent).CallLimitSeconds)

	// try when wait doesn't generate a valid tel URN
	wait, err = waits.ReadWait([]byte(`{"type": "dial", "phone": "@(\"\")"}`))
	assert.NoError(t, err)

	log = test.NewEventLog()
	begun = wait.Begin(run, log.Log)

	assert.False(t, begun)
	assert.Equal(t, 1, len(log.Events))
	assert.Equal(t, "error", log.Events[0].Type())
}
