package engine

import (
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

// the configuration options for the flow engine
type config struct {
	httpClient              *utils.HTTPClient
	disableWebhooks         bool
	maxWebhookResponseBytes int
}

func (c *config) HTTPClient() *utils.HTTPClient { return c.httpClient }
func (c *config) DisableWebhooks() bool         { return c.disableWebhooks }
func (c *config) MaxWebhookResponseBytes() int  { return c.maxWebhookResponseBytes }

//------------------------------------------------------------------------------------------
// Builder
//------------------------------------------------------------------------------------------

// ConfigBuilder is a builder for engine configs
type ConfigBuilder struct {
	config *config
}

// NewConfigBuilder creates a new environment builder
func NewConfigBuilder() *ConfigBuilder {
	return &ConfigBuilder{
		config: &config{
			httpClient:              utils.NewHTTPClient("goflow"),
			disableWebhooks:         false,
			maxWebhookResponseBytes: 10000,
		},
	}
}

// WithDefaultUserAgent sets the default user-agent string used for webhook calls
func (b *ConfigBuilder) WithDefaultUserAgent(userAgent string) *ConfigBuilder {
	b.config.httpClient = utils.NewHTTPClient(userAgent)
	return b
}

// WithDisableWebhooks sets whether webhooks are enabled
func (b *ConfigBuilder) WithDisableWebhooks(disable bool) *ConfigBuilder {
	b.config.disableWebhooks = disable
	return b
}

// WithMaxWebhookResponseBytes sets the maximum webhook request bytes
func (b *ConfigBuilder) WithMaxWebhookResponseBytes(max int) *ConfigBuilder {
	b.config.maxWebhookResponseBytes = max
	return b
}

// Build returns the final config
func (b *ConfigBuilder) Build() flows.EngineConfig { return b.config }
