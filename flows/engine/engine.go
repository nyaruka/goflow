package engine

import (
	"encoding/json"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

// an instance of the engine
type engine struct {
	httpClient              *utils.HTTPClient
	disableWebhooks         bool
	maxWebhookResponseBytes int
}

// NewSession creates a new session
func (e *engine) NewSession(sa flows.SessionAssets) flows.Session {
	return &session{
		engine:     e,
		env:        utils.NewEnvironmentBuilder().Build(),
		assets:     sa,
		status:     flows.SessionStatusActive,
		runsByUUID: make(map[flows.RunUUID]flows.FlowRun),
		flowStack:  newFlowStack(),
	}
}

// ReadSession reads an existing session
func (e *engine) ReadSession(sa flows.SessionAssets, data json.RawMessage, missing assets.MissingCallback) (flows.Session, error) {
	return readSession(e, sa, data, missing)
}

func (e *engine) HTTPClient() *utils.HTTPClient { return e.httpClient }
func (e *engine) DisableWebhooks() bool         { return e.disableWebhooks }
func (e *engine) MaxWebhookResponseBytes() int  { return e.maxWebhookResponseBytes }

var _ flows.Engine = (*engine)(nil)

//------------------------------------------------------------------------------------------
// Builder
//------------------------------------------------------------------------------------------

// Builder is a builder for engine configs
type Builder struct {
	eng *engine
}

// NewBuilder creates a new environment builder
func NewBuilder() *Builder {
	return &Builder{
		eng: &engine{
			httpClient:              utils.NewHTTPClient("goflow"),
			disableWebhooks:         false,
			maxWebhookResponseBytes: 10000,
		},
	}
}

// WithDefaultUserAgent sets the default user-agent string used for webhook calls
func (b *Builder) WithDefaultUserAgent(userAgent string) *Builder {
	b.eng.httpClient = utils.NewHTTPClient(userAgent)
	return b
}

// WithDisableWebhooks sets whether webhooks are enabled
func (b *Builder) WithDisableWebhooks(disable bool) *Builder {
	b.eng.disableWebhooks = disable
	return b
}

// WithMaxWebhookResponseBytes sets the maximum webhook request bytes
func (b *Builder) WithMaxWebhookResponseBytes(max int) *Builder {
	b.eng.maxWebhookResponseBytes = max
	return b
}

// Build returns the final engine
func (b *Builder) Build() flows.Engine { return b.eng }
