package engine

import (
	"encoding/json"

	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
)

// an instance of the engine
type engine struct {
	services             *services
	maxStepsPerSprint    int
	maxResumesPerSession int
	maxTemplateChars     int
}

// NewSession creates a new session
func (e *engine) NewSession(sa flows.SessionAssets, trigger flows.Trigger) (flows.Session, flows.Sprint, error) {
	s := &session{
		uuid:       flows.SessionUUID(uuids.New()),
		engine:     e,
		assets:     sa,
		trigger:    trigger,
		status:     flows.SessionStatusActive,
		batchStart: trigger.Batch(),
		runsByUUID: make(map[flows.RunUUID]flows.Run),
	}

	sprint, err := s.start(trigger)

	return s, sprint, err
}

// ReadSession reads an existing session
func (e *engine) ReadSession(sa flows.SessionAssets, data json.RawMessage, missing assets.MissingCallback) (flows.Session, error) {
	return readSession(e, sa, data, missing)
}

func (e *engine) Services() flows.Services  { return e.services }
func (e *engine) MaxStepsPerSprint() int    { return e.maxStepsPerSprint }
func (e *engine) MaxResumesPerSession() int { return e.maxResumesPerSession }
func (e *engine) MaxTemplateChars() int     { return e.maxTemplateChars }

var _ flows.Engine = (*engine)(nil)

//------------------------------------------------------------------------------------------
// Builder
//------------------------------------------------------------------------------------------

// Builder is a builder for engine configs
type Builder struct {
	eng *engine
}

// NewBuilder creates a new engine builder
func NewBuilder() *Builder {
	return &Builder{
		eng: &engine{
			services:             newEmptyServices(),
			maxStepsPerSprint:    100,
			maxResumesPerSession: 500,
			maxTemplateChars:     10000,
		},
	}
}

// WithEmailServiceFactory sets the email service factory
func (b *Builder) WithEmailServiceFactory(f EmailServiceFactory) *Builder {
	b.eng.services.email = f
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

// WithTicketServiceFactory sets the ticket service factory
func (b *Builder) WithTicketServiceFactory(f TicketServiceFactory) *Builder {
	b.eng.services.ticket = f
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

// WithMaxResumesPerSession sets the maximum number of resumes allowed in a single session
func (b *Builder) WithMaxResumesPerSession(max int) *Builder {
	b.eng.maxResumesPerSession = max
	return b
}

// WithMaxTemplateChars sets the maximum number of characters allowed from an evaluated template
func (b *Builder) WithMaxTemplateChars(max int) *Builder {
	b.eng.maxTemplateChars = max
	return b
}

// Build returns the final engine
func (b *Builder) Build() flows.Engine { return b.eng }
