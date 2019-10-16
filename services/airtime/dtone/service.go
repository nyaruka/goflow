package dtone

import (
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils/httpx"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type service struct {
	login    string
	apiToken string
	currency string
}

// NewService creates a new DTOne airtime service
func NewService(login, apiToken, currency string) flows.AirtimeService {
	return &service{
		login:    login,
		apiToken: apiToken,
		currency: currency,
	}
}

func (s *service) Transfer(session flows.Session, sender urns.URN, recipient urns.URN, amounts map[string]decimal.Decimal) (*flows.AirtimeTransfer, []*httpx.Trace, error) {
	transfer := &flows.AirtimeTransfer{
		Sender:        sender,
		Recipient:     recipient,
		DesiredAmount: decimal.Zero,
		ActualAmount:  decimal.Zero,
	}

	traces := make([]*httpx.Trace, 0, 1)
	client := NewClient(session.Engine().HTTPClient(), s.login, s.apiToken)

	info, trace, err := client.MSISDNInfo(recipient.Path(), s.currency, "1")
	if trace != nil {
		traces = append(traces, trace)
	}
	if err != nil {
		return transfer, traces, err
	}

	transfer.Currency = info.DestinationCurrency

	// look up the amount to send in this currency
	amount, hasAmount := amounts[transfer.Currency]
	if !hasAmount {
		return transfer, traces, errors.Errorf("no amount configured for transfers in %s", transfer.Currency)
	}

	transfer.DesiredAmount = amount

	// find the product closest to our desired amount
	var useProduct string
	useAmount := decimal.Zero
	for p, product := range info.ProductList {
		price := info.LocalInfoValueList[p]
		if price.GreaterThan(useAmount) && price.LessThanOrEqual(amount) {
			useProduct = product
			useAmount = price
		}
	}

	if useAmount == decimal.Zero {
		return transfer, traces, errors.Errorf("amount requested is smaller than the mimimum topup of %s %s", info.LocalInfoValueList[0].String(), info.DestinationCurrency)
	}

	reservedID, trace, err := client.ReserveID()
	if trace != nil {
		traces = append(traces, trace)
	}
	if err != nil {
		return transfer, traces, err
	}

	topup, trace, err := client.Topup(reservedID, sender.Path(), recipient.Path(), useProduct, "")
	if trace != nil {
		traces = append(traces, trace)
	}
	if err != nil {
		return transfer, traces, err
	}

	transfer.ActualAmount = topup.ActualProductSent

	return transfer, traces, nil
}
