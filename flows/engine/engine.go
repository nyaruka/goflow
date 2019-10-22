package engine

import (
	"encoding/json"
	"net/http"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils/uuids"
)

// an instance of the engine
type engine struct {
	httpClient        *http.Client
	services          *services
	maxStepsPerSprint int
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

func (e *engine) HTTPClient() *http.Client { return e.httpClient }
func (e *engine) Services() flows.Services { return e.services }
func (e *engine) MaxStepsPerSprint() int   { return e.maxStepsPerSprint }

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
			httpClient:        http.DefaultClient,
			services:          newEmptyServices(),
			maxStepsPerSprint: 100,
		},
	}
}

// WithHTTPClient sets the HTTP client
func (b *Builder) WithHTTPClient(client *http.Client) *Builder {
	b.eng.httpClient = client
	return b
}

// WithWebhookServiceFactory sets the webhook service factory
func (b *Builder) WithWebhookServiceFactory(f WebhookServiceFactory) *Builder {
	b.eng.services.webhook = f
	return b
}

// WithClassificationServiceFactory sets the NLU service factory
func (b *Builder) WithClassificationServiceFactory(f ClassificationServiceFactory) *Builder {
	b.eng.services.classification = f
	return b
}

// WithAirtimeServiceFactory sets the airtime service factory
func (b *Builder) WithAirtimeServiceFactory(f AirtimeServiceFactory) *Builder {
	b.eng.services.airtime = f
	return b
}

// WithMaxStepsPerSprint sets the maximum number of steps allowed in a single sprint
func (b *Builder) WithMaxStepsPerSprint(max int) *Builder {
	b.eng.maxStepsPerSprint = max
	return b
}

// Build returns the final engine
func (b *Builder) Build() flows.Engine { return b.eng }
