package engine

import (
	"github.com/nyaruka/goflow/flows"
)

// the configuration options for the flow engine
type config struct {
	maxWebhookResponseBytes int
}

// NewDefaultConfig returns the default engine configuration
func NewDefaultConfig() flows.EngineConfig {
	return &config{maxWebhookResponseBytes: 10000}
}

func (c *config) MaxWebhookResponseBytes() int { return c.maxWebhookResponseBytes }
