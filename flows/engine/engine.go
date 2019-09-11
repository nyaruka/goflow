package engine

import (
	"encoding/json"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils/uuids"
)

// an instance of the engine
type engine struct {
	maxStepsPerSprint int
	services          *services
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

func (e *engine) MaxStepsPerSprint() int   { return e.maxStepsPerSprint }
func (e *engine) Services() flows.Services { return e.services }

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
			maxStepsPerSprint: 100,
			services:          newEmptyServices(),
		},
	}
}

// WithMaxStepsPerSprint sets the maximum number of steps allowed in a single sprint
func (b *Builder) WithMaxStepsPerSprint(max int) *Builder {
	b.eng.maxStepsPerSprint = max
	return b
}

// WithWebhookService sets the webhook service
func (b *Builder) WithWebhookService(svc WebhookService) *Builder {
	b.eng.services.webhook = svc
	return b
}

// WithAirtimeService sets the airtime transfer service
func (b *Builder) WithAirtimeService(svc AirtimeService) *Builder {
	b.eng.services.airtime = svc
	return b
}

// Build returns the final engine
func (b *Builder) Build() flows.Engine { return b.eng }
