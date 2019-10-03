package engine

import (
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
)

// WebhookServiceFactory resolves a session to a webhook service
type WebhookServiceFactory func(flows.Session) flows.WebhookService

// NLUServiceFactory resolves a session and classifier to an NLU service
type NLUServiceFactory func(flows.Session, assets.Classifier) flows.NLUService

// AirtimeServiceFactory resolves a session to an airtime service
type AirtimeServiceFactory func(flows.Session) flows.AirtimeService

type services struct {
	webhook WebhookServiceFactory
	nlu     NLUServiceFactory
	airtime AirtimeServiceFactory
}

func newEmptyServices() *services {
	return &services{
		webhook: func(flows.Session) flows.WebhookService { return nil },
		nlu:     func(flows.Session, assets.Classifier) flows.NLUService { return nil },
		airtime: func(flows.Session) flows.AirtimeService { return nil },
	}
}

func (s *services) Webhook(session flows.Session) flows.WebhookService {
	return s.webhook(session)
}

func (s *services) NLU(session flows.Session, classifier assets.Classifier) flows.NLUService {
	return s.nlu(session, classifier)
}

func (s *services) Airtime(session flows.Session) flows.AirtimeService {
	return s.airtime(session)
}
