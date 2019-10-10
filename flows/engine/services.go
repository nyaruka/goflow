package engine

import (
	"github.com/nyaruka/goflow/flows"
	"github.com/pkg/errors"
)

// WebhookServiceFactory resolves a session to a webhook service
type WebhookServiceFactory func(flows.Session) (flows.WebhookService, error)

// ClassificationServiceFactory resolves a session and classifier to an NLU service
type ClassificationServiceFactory func(flows.Session, *flows.Classifier) (flows.ClassificationService, error)

// AirtimeServiceFactory resolves a session to an airtime service
type AirtimeServiceFactory func(flows.Session) (flows.AirtimeService, error)

type services struct {
	webhook        WebhookServiceFactory
	classification ClassificationServiceFactory
	airtime        AirtimeServiceFactory
}

func newEmptyServices() *services {
	return &services{
		webhook: func(flows.Session) (flows.WebhookService, error) {
			return nil, errors.New("no webhook service factory configured")
		},
		classification: func(flows.Session, *flows.Classifier) (flows.ClassificationService, error) {
			return nil, errors.New("no classification service factory configured")
		},
		airtime: func(flows.Session) (flows.AirtimeService, error) {
			return nil, errors.New("no airtime service factory configured")
		},
	}
}

func (s *services) Webhook(session flows.Session) (flows.WebhookService, error) {
	return s.webhook(session)
}

func (s *services) NLU(session flows.Session, classifier *flows.Classifier) (flows.ClassificationService, error) {
	return s.classification(session, classifier)
}

func (s *services) Airtime(session flows.Session) (flows.AirtimeService, error) {
	return s.airtime(session)
}
