package engine

import (
	"github.com/nyaruka/goflow/flows"
)

// WebhookServiceFactory resolves a session to a webhook service
type WebhookServiceFactory func(flows.Session) flows.WebhookService

// AirtimeServiceFactory resolves a session to an airtime service
type AirtimeServiceFactory func(flows.Session) flows.AirtimeService

type services struct {
	webhook WebhookServiceFactory
	airtime AirtimeServiceFactory
}

func newEmptyServices() *services {
	return &services{
		webhook: func(flows.Session) flows.WebhookService { return nil },
		airtime: func(flows.Session) flows.AirtimeService { return nil },
	}
}

func (s *services) Webhook(session flows.Session) flows.WebhookService {
	return s.webhook(session)
}

func (s *services) Airtime(session flows.Session) flows.AirtimeService {
	return s.airtime(session)
}
