package engine

import (
	"github.com/nyaruka/goflow/flows"
)

// the configuration options for the flow engine
type config struct {
	maxWebhookResponseBytes int
}

// NewConfig returns a new engine configuration
func NewConfig(maxWebhookResponseBytes int) flows.EngineConfig {
	return &config{maxWebhookResponseBytes: maxWebhookResponseBytes}
}

// NewDefaultConfig returns the default engine configuration
func NewDefaultConfig() flows.EngineConfig {
	return &config{maxWebhookResponseBytes: 10000}
}

func (c *config) MaxWebhookResponseBytes() int { return c.maxWebhookResponseBytes }
