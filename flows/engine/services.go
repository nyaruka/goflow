package engine

import (
	"net/http"

	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/flows"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type services struct {
	airtime flows.AirtimeService
	webhook flows.WebhookService
}

func newEmptyServices() *services {
	return &services{
		webhook: &nilWebhookService{},
		airtime: &nilAirtimeService{},
	}
}

func (s *services) Webhook() flows.WebhookService {
	return s.webhook
}

func (s *services) Airtime() flows.AirtimeService {
	return s.airtime
}

type nilWebhookService struct{}

// Transfer in this case is a failure
func (s *nilWebhookService) Call(request *http.Request, resthook string) (*flows.WebhookCall, error) {
	return nil, errors.New("no webhook service available")
}

type nilAirtimeService struct{}

// Transfer in this case is a failure
func (s *nilAirtimeService) Transfer(session flows.Session, from urns.URN, to urns.URN, amounts map[string]decimal.Decimal) (*flows.AirtimeTransfer, error) {
	return nil, errors.New("no airtime service available")
}
