package engine

import (
	"errors"

	"github.com/nyaruka/goflow/core"
	"github.com/nyaruka/goflow/flows"
)

// EmailServiceFactory resolves a session to a email service
type EmailServiceFactory func(flows.SessionAssets) (flows.EmailService, error)

// WebhookServiceFactory resolves a session to a webhook service, using the engine's HTTP client and options
type WebhookServiceFactory func(flows.Engine, flows.SessionAssets) (flows.WebhookService, error)

// LLMServiceFactory resolves an LLM asset to to an LLM service
type LLMServiceFactory func(*core.LLM) (flows.LLMService, error)

// AirtimeServiceFactory resolves a session to an airtime service
type AirtimeServiceFactory func(flows.SessionAssets) (flows.AirtimeService, error)

type services struct {
	engine  flows.Engine
	email   EmailServiceFactory
	webhook WebhookServiceFactory
	llm     LLMServiceFactory
	airtime AirtimeServiceFactory
}

func newEmptyServices() *services {
	return &services{
		email: func(flows.SessionAssets) (flows.EmailService, error) {
			return nil, errors.New("no email service factory configured")
		},
		webhook: func(flows.Engine, flows.SessionAssets) (flows.WebhookService, error) {
			return nil, errors.New("no webhook service factory configured")
		},
		llm: func(*core.LLM) (flows.LLMService, error) {
			return nil, errors.New("no LLM service factory configured")
		},
		airtime: func(flows.SessionAssets) (flows.AirtimeService, error) {
			return nil, errors.New("no airtime service factory configured")
		},
	}
}

func (s *services) Email(sa flows.SessionAssets) (flows.EmailService, error) {
	return s.email(sa)
}

func (s *services) Webhook(sa flows.SessionAssets) (flows.WebhookService, error) {
	return s.webhook(s.engine, sa)
}

func (s *services) LLM(llm *core.LLM) (flows.LLMService, error) {
	return s.llm(llm)
}

func (s *services) Airtime(sa flows.SessionAssets) (flows.AirtimeService, error) {
	return s.airtime(sa)
}
