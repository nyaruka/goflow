package dtone

import (
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/flows"

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

func (s *service) Transfer(session flows.Session, sender urns.URN, recipient urns.URN, amounts map[string]decimal.Decimal, logEvent flows.EventCallback) (*flows.AirtimeTransfer, error) {
	client := NewClient(session.Engine().HTTPClient(), s.login, s.apiToken)

	info, _, err := client.MSISDNInfo(recipient.Path(), s.currency, "1")
	if err != nil {
		return nil, err
	}

	// look up the amount to send in this currency
	amount, hasAmount := amounts[info.DestinationCurrency]
	if !hasAmount {
		return nil, errors.Errorf("no amount configured for transfers in %s", info.DestinationCurrency)
	}

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
		return nil, errors.Errorf("amount requested is smaller than the mimimum topup of %s %s", info.LocalInfoValueList[0].String(), info.DestinationCurrency)
	}

	reservedID, err := client.ReserveID()
	if err != nil {
		return nil, err
	}

	topup, _, err := client.Topup(reservedID, sender.Path(), recipient.Path(), useProduct, "")
	if err != nil {
		return nil, err
	}

	return &flows.AirtimeTransfer{
		Sender:    sender,
		Recipient: recipient,
		Currency:  info.DestinationCurrency,
		Amount:    topup.ActualProductSent,
	}, nil
}
