package engine

import (
	"net/http"

	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/flows"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type WebhookService func(flows.Session) flows.WebhookProvider
type AirtimeService func(flows.Session) flows.AirtimeProvider

type services struct {
	webhook WebhookService
	airtime AirtimeService
}

func newEmptyServices() *services {
	return &services{
		webhook: func(flows.Session) flows.WebhookProvider { return &nilWebhookProvider{} },
		airtime: func(flows.Session) flows.AirtimeProvider { return &nilAirtimeProvider{} },
	}
}

func (s *services) Webhook(session flows.Session) flows.WebhookProvider {
	return s.webhook(session)
}

func (s *services) Airtime(session flows.Session) flows.AirtimeProvider {
	return s.airtime(session)
}

type nilWebhookProvider struct{}

// Call in this case is a failure
func (s *nilWebhookProvider) Call(request *http.Request, resthook string) (*flows.WebhookCall, error) {
	return nil, errors.New("no webhook service available")
}

type nilAirtimeProvider struct{}

// Transfer in this case is a failure
func (s *nilAirtimeProvider) Transfer(session flows.Session, from urns.URN, to urns.URN, amounts map[string]decimal.Decimal) (*flows.AirtimeTransfer, error) {
	return nil, errors.New("no airtime service available")
}
