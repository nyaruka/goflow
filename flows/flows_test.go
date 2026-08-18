package flows_test

import (
	"testing"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/actions"
	"github.com/nyaruka/goflow/flows/routers/waits"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFlowTypeCanEnter(t *testing.T) {
	tcs := []struct {
		session  flows.FlowType
		entered  flows.FlowType
		expected bool
	}{
		{flows.FlowTypeMessaging, flows.FlowTypeMessaging, true},
		{flows.FlowTypeMessaging, flows.FlowTypeMessagingBackground, true},
		{flows.FlowTypeMessaging, flows.FlowTypeMessagingOffline, false},
		{flows.FlowTypeMessaging, flows.FlowTypeVoice, false},
		{flows.FlowTypeMessagingOffline, flows.FlowTypeMessaging, false},
		{flows.FlowTypeMessagingOffline, flows.FlowTypeMessagingOffline, true},
		{flows.FlowTypeVoice, flows.FlowTypeMessaging, true}, // voice can drop into messaging
		{flows.FlowTypeVoice, flows.FlowTypeMessagingBackground, true},
		{flows.FlowTypeVoice, flows.FlowTypeMessagingOffline, false},
		{flows.FlowTypeVoice, flows.FlowTypeVoice, true},
	}

	for _, tc := range tcs {
		assert.Equal(t, tc.expected, tc.session.CanEnter(tc.entered), "can enter mismatch for %s entering %s", tc.session, tc.entered)
	}
}

func TestFlowTypeAllows(t *testing.T) {
	webhookAction, err := actions.Read([]byte(`{
		"uuid": "8eebd020-1af5-431c-b943-aa670fc74da9",
		"type": "call_webhook",
		"method": "GET",
		"url": "http://localhost:49998/?cmd=success"
	}`))
	require.NoError(t, err)

	assert.True(t, flows.FlowTypeMessaging.Allows(webhookAction))
	assert.False(t, flows.FlowTypeMessagingOffline.Allows(webhookAction))

	msgWait, err := waits.ReadWait([]byte(`{"type": "msg"}`))
	require.NoError(t, err)

	assert.True(t, flows.FlowTypeMessaging.Allows(msgWait))
	assert.False(t, flows.FlowTypeMessagingBackground.Allows(msgWait))
}
