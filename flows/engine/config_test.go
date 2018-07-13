package engine_test

import (
	"testing"

	"github.com/nyaruka/goflow/flows/engine"

	"github.com/stretchr/testify/assert"
)

func TestReadConfig(t *testing.T) {
	base := engine.NewDefaultConfig()

	// can read empty in which case nothing in the base is overridden
	config, err := engine.ReadConfig([]byte(`{}`), base)
	assert.NoError(t, err)
	assert.False(t, config.DisableWebhooks())
	assert.Equal(t, 10000, config.MaxWebhookResponseBytes())

	// or we can override values
	config, err = engine.ReadConfig([]byte(`{"disable_webhooks":true,"max_webhook_response_bytes":1234}`), base)
	assert.NoError(t, err)
	assert.True(t, config.DisableWebhooks())
	assert.Equal(t, 1234, config.MaxWebhookResponseBytes())

	// or add extra properties
	config, err = engine.ReadConfig([]byte(`{"disable_webhooks":true,"max_webhook_response_bytes":1234,"foo":"bar"}`), base)
	assert.NoError(t, err)
	assert.True(t, config.DisableWebhooks())
	assert.Equal(t, 1234, config.MaxWebhookResponseBytes())
	assert.Equal(t, "bar", config.Extra("foo"))
	assert.Equal(t, nil, config.Extra("disable_webhooks")) // core properties aren't duplicated here
}
