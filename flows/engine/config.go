package engine

import (
	"encoding/json"

	"github.com/nyaruka/goflow/flows"
)

// the configuration options for the flow engine
type config struct {
	disableWebhooks         bool
	webhookMocks            []*flows.WebhookMock
	maxWebhookResponseBytes int
}

// NewConfig returns a new engine configuration
func NewConfig(disableWebhooks bool, webhookMocks []*flows.WebhookMock, maxWebhookResponseBytes int) flows.EngineConfig {
	return &config{
		disableWebhooks:         disableWebhooks,
		webhookMocks:            webhookMocks,
		maxWebhookResponseBytes: maxWebhookResponseBytes,
	}
}

// NewDefaultConfig returns the default engine configuration
func NewDefaultConfig() flows.EngineConfig {
	return &config{disableWebhooks: false, webhookMocks: nil, maxWebhookResponseBytes: 10000}
}

func (c *config) DisableWebhooks() bool              { return c.disableWebhooks }
func (c *config) WebhookMocks() []*flows.WebhookMock { return c.webhookMocks }
func (c *config) MaxWebhookResponseBytes() int       { return c.maxWebhookResponseBytes }

type configEnvelope struct {
	DisableWebhooks         *bool                `json:"disable_webhooks"`
	WebhookMocks            []*flows.WebhookMock `json:"webhook_mocks"`
	MaxWebhookResponseBytes *int                 `json:"max_webhook_response_bytes"`
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
	if envelope.WebhookMocks != nil {
		config.webhookMocks = envelope.WebhookMocks
	}
	if envelope.MaxWebhookResponseBytes != nil {
		config.maxWebhookResponseBytes = *envelope.MaxWebhookResponseBytes
	}

	return config, nil
}
