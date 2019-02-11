package engine

import (
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

// an instance of the engine
type engine struct {
	httpClient              *utils.HTTPClient
	disableWebhooks         bool
	maxWebhookResponseBytes int
}

func (e *engine) HTTPClient() *utils.HTTPClient { return e.httpClient }
func (e *engine) DisableWebhooks() bool         { return e.disableWebhooks }
func (e *engine) MaxWebhookResponseBytes() int  { return e.maxWebhookResponseBytes }

//------------------------------------------------------------------------------------------
// Builder
//------------------------------------------------------------------------------------------

// EngineBuilder is a builder for engine configs
type EngineBuilder struct {
	eng *engine
}

// NewEngineBuilder creates a new environment builder
func NewEngineBuilder() *EngineBuilder {
	return &EngineBuilder{
		eng: &engine{
			httpClient:              utils.NewHTTPClient("goflow"),
			disableWebhooks:         false,
			maxWebhookResponseBytes: 10000,
		},
	}
}

// WithDefaultUserAgent sets the default user-agent string used for webhook calls
func (b *EngineBuilder) WithDefaultUserAgent(userAgent string) *EngineBuilder {
	b.eng.httpClient = utils.NewHTTPClient(userAgent)
	return b
}

// WithDisableWebhooks sets whether webhooks are enabled
func (b *EngineBuilder) WithDisableWebhooks(disable bool) *EngineBuilder {
	b.eng.disableWebhooks = disable
	return b
}

// WithMaxWebhookResponseBytes sets the maximum webhook request bytes
func (b *EngineBuilder) WithMaxWebhookResponseBytes(max int) *EngineBuilder {
	b.eng.maxWebhookResponseBytes = max
	return b
}

// Build returns the final engine
func (b *EngineBuilder) Build() flows.Engine { return b.eng }
