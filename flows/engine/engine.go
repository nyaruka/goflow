package engine

import (
	"encoding/json"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
	"github.com/nyaruka/goflow/utils/uuids"
)

// an instance of the engine
type engine struct {
	httpClient              *utils.HTTPClient
	disableWebhooks         bool
	maxWebhookResponseBytes int
	maxStepsPerSprint       int
}

// NewSession creates a new session
func (e *engine) NewSession(sa flows.SessionAssets, trigger flows.Trigger) (flows.Session, flows.Sprint, error) {
	s := &session{
		uuid:       flows.SessionUUID(uuids.New()),
		engine:     e,
		assets:     sa,
		trigger:    trigger,
		status:     flows.SessionStatusActive,
		runsByUUID: make(map[flows.RunUUID]flows.FlowRun),
	}

	sprint, err := s.start(trigger)

	return s, sprint, err
}

// ReadSession reads an existing session
func (e *engine) ReadSession(sa flows.SessionAssets, data json.RawMessage, missing assets.MissingCallback) (flows.Session, error) {
	return readSession(e, sa, data, missing)
}

func (e *engine) HTTPClient() *utils.HTTPClient { return e.httpClient }
func (e *engine) DisableWebhooks() bool         { return e.disableWebhooks }
func (e *engine) MaxWebhookResponseBytes() int  { return e.maxWebhookResponseBytes }
func (e *engine) MaxStepsPerSprint() int        { return e.maxStepsPerSprint }

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
			maxStepsPerSprint:       100,
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

// WithMaxStepsPerSprint sets the maximum number of steps allowed in a single sprint
func (b *Builder) WithMaxStepsPerSprint(max int) *Builder {
	b.eng.maxStepsPerSprint = max
	return b
}

// Build returns the final engine
func (b *Builder) Build() flows.Engine { return b.eng }
