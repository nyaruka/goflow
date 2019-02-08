package engine

import (
	"github.com/nyaruka/goflow/flows"
)

// the configuration options for the flow engine
type config struct {
	maxFieldLength          int
	maxResultLength         int
	disableWebhooks         bool
	maxWebhookResponseBytes int
}

// NewConfig returns a new engine configuration
func NewConfig(maxFieldLength int, maxResultLength int, disableWebhooks bool, maxWebhookResponseBytes int) flows.EngineConfig {
	return &config{
		maxFieldLength:          maxFieldLength,
		maxResultLength:         maxResultLength,
		disableWebhooks:         disableWebhooks,
		maxWebhookResponseBytes: maxWebhookResponseBytes,
	}
}

// NewDefaultConfig returns the default engine configuration
func NewDefaultConfig() flows.EngineConfig {
	return &config{
		maxFieldLength:          640,
		maxResultLength:         640,
		disableWebhooks:         false,
		maxWebhookResponseBytes: 10000,
	}
}

func (c *config) MaxFieldLength() int          { return c.maxFieldLength }
func (c *config) MaxResultLength() int         { return c.maxResultLength }
func (c *config) DisableWebhooks() bool        { return c.disableWebhooks }
func (c *config) MaxWebhookResponseBytes() int { return c.maxWebhookResponseBytes }
