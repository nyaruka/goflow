package engine

import (
	"context"
	"text/template"

	"github.com/nyaruka/gocommon/dates"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent"
	"github.com/nyaruka/goflow/flows"
)

// an instance of the engine
type engine struct {
	evaluator *excellent.Evaluator
	services  *services
	options   *flows.EngineOptions
}

// NewSession creates a new session
func (e *engine) NewSession(ctx context.Context, sa flows.SessionAssets, contact *flows.Contact, trigger flows.Trigger, call *flows.Call) (flows.Session, flows.Sprint, error) {
	s := &session{
		uuid:       flows.NewSessionUUID(),
		createdOn:  dates.Now(),
		env:        envs.NewBuilder().Build(),
		engine:     e,
		assets:     sa,
		contact:    contact,
		trigger:    trigger,
		status:     flows.SessionStatusActive,
		batchStart: trigger.Batch(),
		runsByUUID: make(map[flows.RunUUID]flows.Run),
		call:       call,
	}

	sprint, err := s.start(ctx, trigger)

	return s, sprint, err
}

// ReadSession reads an existing session
func (e *engine) ReadSession(sa flows.SessionAssets, data []byte, contact *flows.Contact, call *flows.Call, missing assets.MissingCallback) (flows.Session, error) {
	return readSession(e, sa, data, contact, call, missing)
}

func (e *engine) Evaluator() *excellent.Evaluator { return e.evaluator }
func (e *engine) Services() flows.Services        { return e.services }
func (e *engine) Options() *flows.EngineOptions   { return e.options }

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
			evaluator: excellent.NewEvaluator(),
			services:  newEmptyServices(),
			options: &flows.EngineOptions{
				MaxStepsPerSprint:    100,
				MaxResumesPerSession: 500,
				MaxTemplateChars:     10000,
				MaxFieldChars:        640,
				MaxResultChars:       640,
				LLMPrompts:           make(map[string]*template.Template),
			},
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

// WithLLMServiceFactory sets the LLM service factory
func (b *Builder) WithLLMServiceFactory(f LLMServiceFactory) *Builder {
	b.eng.services.llm = f
	return b
}

// WithAirtimeServiceFactory sets the airtime service factory
func (b *Builder) WithAirtimeServiceFactory(f AirtimeServiceFactory) *Builder {
	b.eng.services.airtime = f
	return b
}

// WithMaxStepsPerSprint sets the maximum number of steps allowed in a single sprint
func (b *Builder) WithMaxStepsPerSprint(max int) *Builder {
	b.eng.options.MaxStepsPerSprint = max
	return b
}

// WithMaxResumesPerSession sets the maximum number of resumes allowed in a single session
func (b *Builder) WithMaxResumesPerSession(max int) *Builder {
	b.eng.options.MaxResumesPerSession = max
	return b
}

// WithMaxTemplateChars sets the maximum number of characters allowed from an evaluated template
func (b *Builder) WithMaxTemplateChars(max int) *Builder {
	b.eng.options.MaxTemplateChars = max
	return b
}

// WithMaxFieldChars sets the maximum number of characters allowed in a contact field value
func (b *Builder) WithMaxFieldChars(max int) *Builder {
	b.eng.options.MaxFieldChars = max
	return b
}

// WithMaxResultChars sets the maximum number of characters allowed in a result value
func (b *Builder) WithMaxResultChars(max int) *Builder {
	b.eng.options.MaxResultChars = max
	return b
}

// WithLLMPrompts sets the LLM prompts to use with LLM services
func (b *Builder) WithLLMPrompts(prompts map[string]*template.Template) *Builder {
	b.eng.options.LLMPrompts = prompts
	return b
}

// Build returns the final engine
func (b *Builder) Build() flows.Engine { return b.eng }
