package engine

import (
	"encoding/json"
	"github.com/nyaruka/goflow/flows"
)

// the configuration options for the flow engine
type config struct {
	disableWebhooks         bool
	maxWebhookResponseBytes int
}

// NewConfig returns a new engine configuration
func NewConfig(disableWebhooks bool, maxWebhookResponseBytes int) flows.EngineConfig {
	return &config{
		disableWebhooks:         disableWebhooks,
		maxWebhookResponseBytes: maxWebhookResponseBytes,
	}
}

// NewDefaultConfig returns the default engine configuration
func NewDefaultConfig() flows.EngineConfig {
	return &config{disableWebhooks: false, maxWebhookResponseBytes: 10000}
}

func (c *config) DisableWebhooks() bool        { return c.disableWebhooks }
func (c *config) MaxWebhookResponseBytes() int { return c.maxWebhookResponseBytes }

type configEnvelope struct {
	DisableWebhooks         *bool `json:"disable_webhooks"`
	MaxWebhookResponseBytes *int  `json:"max_webhook_response_bytes"`
}

func ReadConfig(data json.RawMessage, base flows.EngineConfig) (flows.EngineConfig, error) {
	var envelope configEnvelope
	if err := json.Unmarshal(data, &envelope); err != nil {
		return nil, err
	}

	config := base.(*config)

	if envelope.DisableWebhooks != nil {
		config.disableWebhooks = *envelope.DisableWebhooks
	}
	if envelope.MaxWebhookResponseBytes != nil {
		config.maxWebhookResponseBytes = *envelope.MaxWebhookResponseBytes
	}

	return config, nil
}
