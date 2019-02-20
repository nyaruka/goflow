package mobile_test

import (
	"io/ioutil"
	"testing"

	"github.com/nyaruka/goflow/mobile"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMobileBindings(t *testing.T) {
	assert.True(t, mobile.IsSpecVersionSupported("11.6"))
	assert.True(t, mobile.IsSpecVersionSupported("12"))
	assert.True(t, mobile.IsSpecVersionSupported("12.5"))
	assert.False(t, mobile.IsSpecVersionSupported("13.3"))

	// error if we try to create assets from invalid JSON
	_, err := mobile.NewAssetsSource("{")
	assert.Error(t, err)

	// can load a standard assets file
	assetsJSON, err := ioutil.ReadFile("../test/testdata/flows/two_questions_offline.json")
	require.NoError(t, err)

	source, err := mobile.NewAssetsSource(string(assetsJSON))
	require.NoError(t, err)

	// and create a new session assets
	sa, err := mobile.NewSessionAssets(source)
	require.NoError(t, err)

	langs := mobile.NewStringSlice(2)
	langs.Add("eng")
	langs.Add("fra")
	environment, err := mobile.NewEnvironment("DD-MM-YYYY", "tt:mm", "Africa/Kigali", "eng", langs, "RW", "none")
	require.NoError(t, err)

	contact := mobile.NewEmptyContact(sa)

	trigger := mobile.NewManualTrigger(environment, contact, mobile.NewFlowReference("7c3db26f-e12a-48af-9673-e2feefdf8516", "Two Questions"))

	eng := mobile.NewEngine("mobile-test")
	session := eng.NewSession(sa)
	assert.Equal(t, sa, session.Assets())

	sprint, err := session.Start(trigger)
	require.NoError(t, err)

	assert.Equal(t, "waiting", session.Status())

	events := sprint.Events()
	assert.Equal(t, 2, events.Length())
	assert.Equal(t, "msg_created", events.Get(0).Type())
	assert.Equal(t, "msg_wait", events.Get(1).Type())

	modifiers := sprint.Modifiers()
	assert.Equal(t, 0, modifiers.Length())

	wait := session.GetWait()
	assert.Equal(t, "msg", wait.Type())
	assert.Nil(t, wait.Hint())

	attachments := mobile.NewStringSlice(1)
	attachments.Add("content://io.rapidpro.surveyor/files/selfie.jpg")
	msg := mobile.NewMsgIn("8e6f0213-a122-4c50-a430-442085754c16", "Hi there", attachments)

	assert.Equal(t, "Hi there", msg.Text())
	assert.Equal(t, 1, msg.Attachments().Length())

	resume := mobile.NewMsgResume(nil, nil, msg)

	sprint, err = session.Resume(resume)
	require.NoError(t, err)

	events = sprint.Events()
	assert.Equal(t, 4, events.Length())
	assert.Equal(t, "msg_received", events.Get(0).Type())
	assert.Equal(t, `{"type":"msg_received","created_`, events.Get(0).Payload()[:32])
	assert.Equal(t, "run_result_changed", events.Get(1).Type())
	assert.Equal(t, "msg_created", events.Get(2).Type())
	assert.Equal(t, "msg_wait", events.Get(3).Type())

	// convert session to JSON
	marshaled, err := session.ToJSON()
	require.NoError(t, err)

	assert.Equal(t, `{"type":"messaging_offline","environment":{"date_f`, marshaled[:50])

	// and try to read it back
	session2, err := eng.ReadSession(sa, marshaled)
	require.NoError(t, err)

	assert.Equal(t, "waiting", session2.Status())
}
