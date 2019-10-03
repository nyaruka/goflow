package test

import (
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/engine"
	"github.com/nyaruka/goflow/providers/webhooks"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

// NewEngine creates an engine instance for testing
func NewEngine() flows.Engine {
	return engine.NewBuilder().
		WithWebhookServiceFactory(webhooks.NewServiceFactory("goflow-testing", 10000)).
		WithAirtimeServiceFactory(func(flows.Session) flows.AirtimeService { return newAirtimeService("RWF") }).
		Build()
}

// implementation of an airtime service for testing which uses a fixed currency
type airtimeService struct {
	fixedCurrency string
}

func newAirtimeService(currency string) *airtimeService {
	return &airtimeService{fixedCurrency: currency}
}

func (s *airtimeService) Transfer(session flows.Session, sender urns.URN, recipient urns.URN, amounts map[string]decimal.Decimal) (*flows.AirtimeTransfer, error) {
	t := &flows.AirtimeTransfer{
		Sender:    sender,
		Recipient: recipient,
		Currency:  s.fixedCurrency,
		Status:    flows.AirtimeTransferStatusFailed,
	}

	amount, hasAmount := amounts[s.fixedCurrency]
	if !hasAmount {
		return t, errors.Errorf("no amount configured for transfers in %s", s.fixedCurrency)
	}

	t.DesiredAmount = amount
	t.ActualAmount = amount
	t.Status = flows.AirtimeTransferStatusSuccess
	return t, nil
}

var _ flows.AirtimeService = (*airtimeService)(nil)
