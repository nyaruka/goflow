package engine

import (
	"github.com/nyaruka/goflow/flows"

	"github.com/pkg/errors"
)

// EmailServiceFactory resolves a session to a email service
type EmailServiceFactory func(flows.Session) (flows.EmailService, error)

// WebhookServiceFactory resolves a session to a webhook service
type WebhookServiceFactory func(flows.Session) (flows.WebhookService, error)

// ClassificationServiceFactory resolves a session and classifier to an NLU service
type ClassificationServiceFactory func(flows.Session, *flows.Classifier) (flows.ClassificationService, error)

// TicketServiceFactory resolves a session to a ticket service
type TicketServiceFactory func(flows.Session, *flows.Ticketer) (flows.TicketService, error)

// AirtimeServiceFactory resolves a session to an airtime service
type AirtimeServiceFactory func(flows.Session) (flows.AirtimeService, error)

type services struct {
	email          EmailServiceFactory
	webhook        WebhookServiceFactory
	classification ClassificationServiceFactory
	ticket         TicketServiceFactory
	airtime        AirtimeServiceFactory
}

func newEmptyServices() *services {
	return &services{
		email: func(flows.Session) (flows.EmailService, error) {
			return nil, errors.New("no email service factory configured")
		},
		webhook: func(flows.Session) (flows.WebhookService, error) {
			return nil, errors.New("no webhook service factory configured")
		},
		classification: func(flows.Session, *flows.Classifier) (flows.ClassificationService, error) {
			return nil, errors.New("no classification service factory configured")
		},
		ticket: func(flows.Session, *flows.Ticketer) (flows.TicketService, error) {
			return nil, errors.New("no ticket service factory configured")
		},
		airtime: func(flows.Session) (flows.AirtimeService, error) {
			return nil, errors.New("no airtime service factory configured")
		},
	}
}

func (s *services) Email(session flows.Session) (flows.EmailService, error) {
	return s.email(session)
}

func (s *services) Webhook(session flows.Session) (flows.WebhookService, error) {
	return s.webhook(session)
}

func (s *services) Classification(session flows.Session, classifier *flows.Classifier) (flows.ClassificationService, error) {
	return s.classification(session, classifier)
}

func (s *services) Ticket(session flows.Session, ticketer *flows.Ticketer) (flows.TicketService, error) {
	return s.ticket(session, ticketer)
}

func (s *services) Airtime(session flows.Session) (flows.AirtimeService, error) {
	return s.airtime(session)
}
