package mobile_test

import (
	"io/ioutil"
	"testing"

	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/flows/definition"
	"github.com/nyaruka/goflow/mobile"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMobileBindings(t *testing.T) {
	defer uuids.SetGenerator(uuids.DefaultGenerator)
	uuids.SetGenerator(uuids.NewSeededGenerator(1234))

	assert.Equal(t, definition.CurrentSpecVersion.String(), mobile.CurrentSpecVersion())

	assert.False(t, mobile.IsVersionSupported("x"))
	assert.True(t, mobile.IsVersionSupported("11.12"))
	assert.True(t, mobile.IsVersionSupported("13"))
	assert.True(t, mobile.IsVersionSupported("13.3"))
	assert.False(t, mobile.IsVersionSupported("14.0"))

	// error if we try to create assets from invalid JSON
	_, err := mobile.NewAssetsSource("{")
	assert.Error(t, err)

	// can load a standard assets file
	assetsJSON, err := ioutil.ReadFile("../test/testdata/runner/two_questions_offline.json")
	require.NoError(t, err)

	source, err := mobile.NewAssetsSource(string(assetsJSON))
	require.NoError(t, err)

	langs := mobile.NewStringSlice(2)
	langs.Add("eng")
	langs.Add("fra")
	environment, err := mobile.NewEnvironment("DD-MM-YYYY", "tt:mm", "Africa/Kigali", langs, "RW", "none")
	require.NoError(t, err)

	// and create a new session assets
	sa, err := mobile.NewSessionAssets(environment, source)
	require.NoError(t, err)

	contact := mobile.NewEmptyContact(sa)

	trigger := mobile.NewManualTrigger(environment, contact, mobile.NewFlowReference("7c3db26f-e12a-48af-9673-e2feefdf8516", "Two Questions"))

	eng := mobile.NewEngine()
	ss, err := eng.NewSession(sa, trigger)
	session := ss.Session()
	sprint := ss.Sprint()
	require.NoError(t, err)
	assert.Equal(t, sa, session.Assets())
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

	assert.Equal(t, `{"uuid":"cdf7ed27-5ad5-4028-b664-880fc7581c77","ty`, marshaled[:50])

	// and try to read it back
	session2, err := eng.ReadSession(sa, marshaled)
	require.NoError(t, err)

	assert.Equal(t, "waiting", session2.Status())
}
