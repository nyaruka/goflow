package test

import (
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/engine"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

// NewEngine creates an engine instance for testing
func NewEngine() flows.Engine {
	return engine.NewBuilder().
		WithDefaultUserAgent("goflow-testing").
		WithAirtimeService(newAirtimeService("RWF")).
		Build()
}

// implementation of AirtimeService for testing which uses a fixed currency
type airtimeService struct {
	fixedCurrency string
}

func newAirtimeService(currency string) *airtimeService {
	return &airtimeService{fixedCurrency: currency}
}

func (s *airtimeService) Transfer(session flows.Session, from urns.URN, to urns.URN, amounts map[string]decimal.Decimal) (*flows.AirtimeTransfer, error) {
	t := &flows.AirtimeTransfer{
		Currency: s.fixedCurrency,
		Status:   flows.AirtimeTransferStatusFailed,
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
