package engine

import (
	"github.com/nyaruka/goflow/flows"

	"github.com/pkg/errors"
)

// EmailServiceFactory resolves a session to a email service
type EmailServiceFactory func(flows.SessionAssets) (flows.EmailService, error)

// WebhookServiceFactory resolves a session to a webhook service
type WebhookServiceFactory func(flows.SessionAssets) (flows.WebhookService, error)

// ClassificationServiceFactory resolves a session and classifier to an NLU service
type ClassificationServiceFactory func(*flows.Classifier) (flows.ClassificationService, error)

// AirtimeServiceFactory resolves a session to an airtime service
type AirtimeServiceFactory func(flows.SessionAssets) (flows.AirtimeService, error)

type services struct {
	email          EmailServiceFactory
	webhook        WebhookServiceFactory
	classification ClassificationServiceFactory
	airtime        AirtimeServiceFactory
}

func newEmptyServices() *services {
	return &services{
		email: func(flows.SessionAssets) (flows.EmailService, error) {
			return nil, errors.New("no email service factory configured")
		},
		webhook: func(flows.SessionAssets) (flows.WebhookService, error) {
			return nil, errors.New("no webhook service factory configured")
		},
		classification: func(*flows.Classifier) (flows.ClassificationService, error) {
			return nil, errors.New("no classification service factory configured")
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
	return s.webhook(sa)
}

func (s *services) Classification(classifier *flows.Classifier) (flows.ClassificationService, error) {
	return s.classification(classifier)
}

func (s *services) Airtime(sa flows.SessionAssets) (flows.AirtimeService, error) {
	return s.airtime(sa)
}
