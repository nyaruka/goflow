package mobile_test

import (
	"io/ioutil"
	"testing"

	"github.com/nyaruka/goflow/mobile"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMobileBindings(t *testing.T) {
	assert.True(t, mobile.IsSpecVersionSupported("12"))
	assert.False(t, mobile.IsSpecVersionSupported("11.6"))

	// error if we try to create assets from invalid JSON
	_, err := mobile.NewAssetsSource("{")
	assert.Error(t, err)

	// can load a standard assets file
	assetsJSON, err := ioutil.ReadFile("../test/testdata/flows/two_questions_offline.json")
	require.NoError(t, err)

	source, err := mobile.NewAssetsSource(string(assetsJSON))
	require.NoError(t, err)

	// and create a new session assets
	sessionAssets, err := mobile.NewSessionAssets(source)
	require.NoError(t, err)

	langs := mobile.NewStringSlice(2)
	langs.Add("eng")
	langs.Add("fra")
	environment, err := mobile.NewEnvironment("DD-MM-YYYY", "tt:mm", "Africa/Kigali", "eng", langs, "none")
	require.NoError(t, err)

	contact := mobile.NewEmptyContact()

	trigger := mobile.NewManualTrigger(environment, contact, mobile.NewFlowReference("7c3db26f-e12a-48af-9673-e2feefdf8516", "Two Questions"))

	session := mobile.NewSession(sessionAssets, "mobile-test")

	events, err := session.Start(trigger)
	require.NoError(t, err)

	assert.Equal(t, "waiting", session.GetStatus())
	assert.Equal(t, 2, events.Length())
	assert.Equal(t, "msg_created", events.Get(0).GetType())
	assert.Equal(t, "msg_wait", events.Get(1).GetType())

	wait := session.GetWait()
	assert.Equal(t, "msg", wait.GetType())
	assert.Equal(t, "", wait.GetMediaHint())

	resume := mobile.NewMsgResume(nil, nil, mobile.NewMsgIn("8e6f0213-a122-4c50-a430-442085754c16", "Hi there", nil))

	events, err = session.Resume(resume)
	require.NoError(t, err)

	assert.Equal(t, 4, events.Length())
	assert.Equal(t, "msg_received", events.Get(0).GetType())
	assert.Equal(t, `{"type":"msg_received","created_`, events.Get(0).GetPayload()[:32])
	assert.Equal(t, "run_result_changed", events.Get(1).GetType())
	assert.Equal(t, "msg_created", events.Get(2).GetType())
	assert.Equal(t, "msg_wait", events.Get(3).GetType())

	// convert session to JSON
	marshaled, err := session.ToJSON()
	require.NoError(t, err)

	assert.Equal(t, `{"environment":{"date_format":"DD-MM-YYYY","time_f`, marshaled[:50])

	// and try to read it back
	session2, err := mobile.ReadSession(sessionAssets, "mobile-test", marshaled)
	require.NoError(t, err)

	assert.Equal(t, "waiting", session2.GetStatus())
}
