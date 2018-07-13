package engine

import (
	"encoding/json"

	"github.com/nyaruka/goflow/flows"

	"github.com/mitchellh/mapstructure"
)

type coreConfig struct {
	DisableWebhooks         bool                 `mapstructure:"disable_webhooks"`
	WebhookMocks            []*flows.WebhookMock `mapstructure:"webhook_mocks"`
	MaxWebhookResponseBytes int                  `mapstructure:"max_webhook_response_bytes"`
}

// the configuration options for the flow engine
type config struct {
	core coreConfig
	raw  map[string]interface{}
}

// NewConfig returns a new engine configuration
func NewConfig(disableWebhooks bool, webhookMocks []*flows.WebhookMock, maxWebhookResponseBytes int) flows.EngineConfig {
	return &config{
		core: coreConfig{
			DisableWebhooks:         disableWebhooks,
			WebhookMocks:            webhookMocks,
			MaxWebhookResponseBytes: maxWebhookResponseBytes,
		},
	}
}

// NewDefaultConfig returns the default engine configuration
func NewDefaultConfig() flows.EngineConfig {
	return NewConfig(false, nil, 10000)
}

func (c *config) DisableWebhooks() bool              { return c.core.DisableWebhooks }
func (c *config) WebhookMocks() []*flows.WebhookMock { return c.core.WebhookMocks }
func (c *config) MaxWebhookResponseBytes() int       { return c.core.MaxWebhookResponseBytes }

func (c *config) ReadInto(s interface{}) error {
	return mapstructure.Decode(c.raw, s)
}

// ReadConfig reads an engine configuration from the given JSON, using the provided config as a base
func ReadConfig(data json.RawMessage, defaults map[string]interface{}) (flows.EngineConfig, error) {
	// unmarshal as map initially
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	// add any missing values from our defaults
	for key, val := range defaults {
		_, defined := raw[key]
		if !defined {
			raw[key] = val
		}
	}

	// now unmarshall that into the core config
	c := &config{core: coreConfig{}, raw: raw}
	if err := mapstructure.Decode(raw, &c.core); err != nil {
		return nil, err
	}

	return c, nil
}
