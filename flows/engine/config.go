package engine

import (
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
