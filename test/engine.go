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
		WithWebhookService(webhooks.NewService("goflow-testing", 10000)).
		WithAirtimeService(func(flows.Session) flows.AirtimeProvider { return newAirtimeProvider("RWF") }).
		Build()
}

// implementation of AirtimeProvider for testing which uses a fixed currency
type airtimeProvider struct {
	fixedCurrency string
}

func newAirtimeProvider(currency string) *airtimeProvider {
	return &airtimeProvider{fixedCurrency: currency}
}

func (s *airtimeProvider) Transfer(session flows.Session, sender urns.URN, recipient urns.URN, amounts map[string]decimal.Decimal) (*flows.AirtimeTransfer, error) {
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

var _ flows.AirtimeProvider = (*airtimeProvider)(nil)
