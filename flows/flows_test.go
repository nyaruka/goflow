package flows_test

import (
	"testing"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/actions"
	"github.com/nyaruka/goflow/flows/routers/waits"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFlowTypeAllows(t *testing.T) {
	webhookAction, err := actions.ReadAction([]byte(`{
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
