package dtone

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/nyaruka/gocommon/httpx"
	"github.com/nyaruka/gocommon/stringsx"
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/flows"

	"github.com/shopspring/decimal"
)

type service struct {
	client   *Client
	redactor stringsx.Redactor
}

// NewService creates a new DTOne airtime service
func NewService(httpClient *http.Client, httpRetries *httpx.RetryConfig, key, secret string) flows.AirtimeService {
	return &service{
		client:   NewClient(httpClient, httpRetries, key, secret),
		redactor: stringsx.NewRedactor(flows.RedactionMask, secret),
	}
}

func (s *service) Transfer(ctx context.Context, sender urns.URN, recipient urns.URN, amounts map[string]decimal.Decimal, logHTTP flows.HTTPLogCallback) (*flows.AirtimeTransfer, error) {
	transfer := &flows.AirtimeTransfer{
		UUID:      flows.AirtimeTransferUUID(uuids.NewV4()),
		Sender:    sender,
		Recipient: recipient,
		Currency:  "",
		Amount:    decimal.Zero,
	}
	recipientPhoneNumber := recipient.Path()
	if !strings.HasPrefix(recipientPhoneNumber, "+") {
		recipientPhoneNumber = "+" + recipientPhoneNumber
	}

	operators, trace, err := s.client.LookupMobileNumber(ctx, recipientPhoneNumber)
	if trace != nil {
		logHTTP(flows.NewHTTPLog(trace, flows.HTTPStatusFromCode, s.redactor))
	}
	if err != nil {
		return transfer, fmt.Errorf("number lookup failed: %w", err)
	}

	// look for an exact match
	var operator *Operator
	for _, op := range operators {
		if op.Identified {
			operator = op
			break
		}
	}
	if operator == nil {
		return transfer, fmt.Errorf("unable to find operator for number %s", recipientPhoneNumber)
	}

	// fetch available products for this operator
	products, trace, err := s.client.Products(ctx, "FIXED_VALUE_RECHARGE", operator.ID)
	if trace != nil {
		logHTTP(flows.NewHTTPLog(trace, flows.HTTPStatusFromCode, s.redactor))
	}
	if err != nil {
		return transfer, fmt.Errorf("product fetch failed: %w", err)
	}

	// find a matching product in any currency we have a desired amount for
	var product *Product
	for currency, desiredAmount := range amounts {
		for _, p := range products {
			if p.Destination.Unit == currency {
				if p.Destination.Amount.Equal(desiredAmount) {
					product = p
					break
				}
			}
		}
	}
	if product == nil {
		return transfer, fmt.Errorf("unable to find a suitable product for operator '%s'", operator.Name)
	}

	transfer.Currency = product.Destination.Unit
	transfer.Amount = product.Destination.Amount

	// request asynchronous confirmed transaction for this product
	tx, trace, err := s.client.TransactionAsync(ctx, string(transfer.UUID), product.ID, recipientPhoneNumber)
	if trace != nil {
		logHTTP(flows.NewHTTPLog(trace, flows.HTTPStatusFromCode, s.redactor))
	}
	if err != nil {
		return transfer, fmt.Errorf("transaction creation failed: %w", err)
	}

	if tx.Status.Class.ID != StatusCIDConfirmed && tx.Status.Class.ID != StatusCIDSubmitted && tx.Status.Class.ID != StatusCIDCompleted {
		return transfer, fmt.Errorf("transaction to send product %d on operator %d ended with status %s", product.ID, operator.ID, tx.Status.Message)
	}

	transfer.ExternalID = fmt.Sprintf("%d", tx.ID)

	return transfer, nil
}
