package engine

import (
	"encoding/json"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

// the configuration options for the flow engine
type config struct {
	disableWebhooks         bool
	webhookMocks            []*flows.WebhookMock
	maxWebhookResponseBytes int
	extra                   map[string]interface{}
}

// NewConfig returns a new engine configuration
func NewConfig(disableWebhooks bool, webhookMocks []*flows.WebhookMock, maxWebhookResponseBytes int) flows.EngineConfig {
	return &config{
		disableWebhooks:         disableWebhooks,
		webhookMocks:            webhookMocks,
		maxWebhookResponseBytes: maxWebhookResponseBytes,
		extra: make(map[string]interface{}),
	}
}

// NewDefaultConfig returns the default engine configuration
func NewDefaultConfig() flows.EngineConfig {
	return &config{disableWebhooks: false, webhookMocks: nil, maxWebhookResponseBytes: 10000}
}

func (c *config) DisableWebhooks() bool              { return c.disableWebhooks }
func (c *config) WebhookMocks() []*flows.WebhookMock { return c.webhookMocks }
func (c *config) MaxWebhookResponseBytes() int       { return c.maxWebhookResponseBytes }
func (c *config) Extra(name string) interface{}      { return c.extra[name] }

type configEnvelope struct {
	DisableWebhooks         *bool                `json:"disable_webhooks"`
	WebhookMocks            []*flows.WebhookMock `json:"webhook_mocks"`
	MaxWebhookResponseBytes *int                 `json:"max_webhook_response_bytes"`
}

// ReadConfig reads an engine configuration from the given JSON, using the provided config as a base
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

	// unmarshal again as a map to get non-core properies we don't know about
	if err := json.Unmarshal(data, &config.extra); err != nil {
		return nil, err
	}

	// remove the core properties from the map so they're not duplicated
	for _, prop := range utils.GetJSONFields(envelope) {
		delete(config.extra, prop)
	}

	return config, nil
}
