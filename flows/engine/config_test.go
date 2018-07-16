package engine_test

import (
	"testing"

	"github.com/nyaruka/goflow/flows/engine"

	"github.com/stretchr/testify/assert"
)

func TestReadConfig(t *testing.T) {
	defaults := map[string]interface{}{
		"disable_webhooks":           false,
		"webhook_mocks":              nil,
		"max_webhook_response_bytes": 10000,
	}

	// can read empty in which case nothing in the base is overridden
	config, err := engine.ReadConfig([]byte(`{}`), defaults)
	assert.NoError(t, err)
	assert.False(t, config.DisableWebhooks())
	assert.Equal(t, 10000, config.MaxWebhookResponseBytes())

	// or we can override values
	config, err = engine.ReadConfig([]byte(`{"disable_webhooks":true,"max_webhook_response_bytes":1234}`), defaults)
	assert.NoError(t, err)
	assert.True(t, config.DisableWebhooks())
	assert.Equal(t, 1234, config.MaxWebhookResponseBytes())
}
