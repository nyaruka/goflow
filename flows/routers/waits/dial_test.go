package waits_test

import (
	"testing"

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

	wait, err := waits.ReadWait([]byte(`{"type": "dial", "phone": "@(\"+\" & \"593979123456\")"}`))
	assert.NoError(t, err)
	assert.Equal(t, waits.TypeDial, wait.Type())

	// test marsalling definition wait
	marshaled, err := jsonx.Marshal(wait)
	assert.NoError(t, err)
	assert.Equal(t, `{"type":"dial","phone":"@(\"+\" & \"593979123456\")"}`, string(marshaled))

	// try activating the wait
	log := test.NewEventLog()
	begun := wait.Begin(run, log.Log)

	assert.True(t, begun)
	assert.Equal(t, 1, len(log.Events))
	assert.Equal(t, "dial_wait", log.Events[0].Type())

	// try to end with incorrect resume type
	assert.False(t, wait.Accepts(resumes.NewWaitTimeout(nil, nil)))

	// try to end with dial resume type
	assert.True(t, wait.Accepts(resumes.NewDial(nil, nil, flows.NewDial(flows.DialStatusAnswered, 5))))

	// try when wait has expression error but still generates valid tel URN
	wait, err = waits.ReadWait([]byte(`{"type": "dial", "phone": "+593979123456@(1 / 0)"}`))
	assert.NoError(t, err)

	log = test.NewEventLog()
	begun = wait.Begin(run, log.Log)

	assert.True(t, begun)
	assert.Equal(t, 2, len(log.Events))
	assert.Equal(t, "error", log.Events[0].Type())
	assert.Equal(t, "dial_wait", log.Events[1].Type())
	assert.Equal(t, urns.URN("tel:+593979123456"), log.Events[1].(*events.DialWaitEvent).URN)

	// try when wait doesn't generate a valid tel URN
	wait, err = waits.ReadWait([]byte(`{"type": "dial", "phone": "@(\"\")"}`))
	assert.NoError(t, err)

	log = test.NewEventLog()
	begun = wait.Begin(run, log.Log)

	assert.False(t, begun)
	assert.Equal(t, 1, len(log.Events))
	assert.Equal(t, "error", log.Events[0].Type())
}
