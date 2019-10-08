package engine

import (
	"github.com/nyaruka/goflow/flows"
)

// WebhookServiceFactory resolves a session to a webhook service
type WebhookServiceFactory func(flows.Session) flows.WebhookService

// ClassificationServiceFactory resolves a session and classifier to an NLU service
type ClassificationServiceFactory func(flows.Session, *flows.Classifier) flows.ClassificationService

// AirtimeServiceFactory resolves a session to an airtime service
type AirtimeServiceFactory func(flows.Session) flows.AirtimeService

type services struct {
	webhook WebhookServiceFactory
	nlu     ClassificationServiceFactory
	airtime AirtimeServiceFactory
}

func newEmptyServices() *services {
	return &services{
		webhook: func(flows.Session) flows.WebhookService { return nil },
		nlu:     func(flows.Session, *flows.Classifier) flows.ClassificationService { return nil },
		airtime: func(flows.Session) flows.AirtimeService { return nil },
	}
}

func (s *services) Webhook(session flows.Session) flows.WebhookService {
	return s.webhook(session)
}

func (s *services) NLU(session flows.Session, classifier *flows.Classifier) flows.ClassificationService {
	return s.nlu(session, classifier)
}

func (s *services) Airtime(session flows.Session) flows.AirtimeService {
	return s.airtime(session)
}
