package engine

import (
	"github.com/nyaruka/goflow/flows"
)

// WebhookService resolves a session to a webhook provider
type WebhookService func(flows.Session) flows.WebhookProvider

// AirtimeService resolves a session to an airtime provider
type AirtimeService func(flows.Session) flows.AirtimeProvider

type services struct {
	webhook WebhookService
	airtime AirtimeService
}

func newEmptyServices() *services {
	return &services{
		webhook: func(flows.Session) flows.WebhookProvider { return nil },
		airtime: func(flows.Session) flows.AirtimeProvider { return nil },
	}
}

func (s *services) Webhook(session flows.Session) flows.WebhookProvider {
	return s.webhook(session)
}

func (s *services) Airtime(session flows.Session) flows.AirtimeProvider {
	return s.airtime(session)
}
