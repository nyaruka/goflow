package dtone

import (
	"net/http"

	"github.com/nyaruka/gocommon/httpx"
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"

	"github.com/shopspring/decimal"
)

type service struct {
	client   *Client
	redactor utils.Redactor
}

// NewService creates a new DTOne airtime service
func NewService(httpClient *http.Client, httpRetries *httpx.RetryConfig, key, secret string) flows.AirtimeService {
	return &service{
		client:   NewClient(httpClient, httpRetries, key, secret),
		redactor: utils.NewRedactor(flows.RedactionMask, secret),
	}
}

func (s *service) Transfer(session flows.Session, sender urns.URN, recipient urns.URN, amounts map[string]decimal.Decimal, logHTTP flows.HTTPLogCallback) (*flows.AirtimeTransfer, error) {
	transfer := &flows.AirtimeTransfer{
		Sender:        sender,
		Recipient:     recipient,
		DesiredAmount: decimal.Zero,
		ActualAmount:  decimal.Zero,
	}

	/*info, trace, err := s.client.LookupMobileNumber(recipient.Path())
	if trace != nil {
		logHTTP(flows.NewHTTPLog(trace, httpLogStatus, s.redactor))
	}
	if err != nil {
		return transfer, err
	}

	transfer.Currency = info.DestinationCurrency

	// look up the amount to send in this currency
	amount, hasAmount := amounts[transfer.Currency]
	if !hasAmount {
		return transfer, errors.Errorf("no amount configured for transfers in %s", transfer.Currency)
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
		return transfer, errors.Errorf("amount requested is smaller than the minimum topup of %s %s", info.LocalInfoValueList[0].String(), info.DestinationCurrency)
	}

	reservedID, trace, err := s.client.ReserveID()
	if trace != nil {
		logHTTP(flows.NewHTTPLog(trace, httpLogStatus, s.redactor))
	}
	if err != nil {
		return transfer, err
	}

	topup, trace, err := s.client.Topup(reservedID.ReservedID, sender.Path(), recipient.Path(), useProduct, "")
	if trace != nil {
		logHTTP(flows.NewHTTPLog(trace, httpLogStatus, s.redactor))
	}
	if err != nil {
		return transfer, err
	}

	transfer.ActualAmount = topup.ActualProductSent*/

	return transfer, nil
}
