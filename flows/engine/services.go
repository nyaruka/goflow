package engine

import (
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
)

// WebhookService resolves a session to a webhook provider
type WebhookService func(flows.Session) flows.WebhookProvider

// NLUService resolves a session and classifier to an NLU provider
type NLUService func(flows.Session, assets.Classifier) flows.NLUProvider

// AirtimeService resolves a session to an airtime provider
type AirtimeService func(flows.Session) flows.AirtimeProvider

type services struct {
	webhook WebhookService
	nlu     NLUService
	airtime AirtimeService
}

func newEmptyServices() *services {
	return &services{
		webhook: func(flows.Session) flows.WebhookProvider { return nil },
		nlu:     func(flows.Session, assets.Classifier) flows.NLUProvider { return nil },
		airtime: func(flows.Session) flows.AirtimeProvider { return nil },
	}
}

func (s *services) Webhook(session flows.Session) flows.WebhookProvider {
	return s.webhook(session)
}

func (s *services) NLU(session flows.Session, classifier assets.Classifier) flows.NLUProvider {
	return s.nlu(session, classifier)
}

func (s *services) Airtime(session flows.Session) flows.AirtimeProvider {
	return s.airtime(session)
}
