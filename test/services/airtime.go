package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/nyaruka/gocommon/httpx"
	"github.com/nyaruka/gocommon/random"
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/flows"
	"github.com/shopspring/decimal"
)

// Airtime is an implementation of an airtime service for testing which uses a fixed currency
type Airtime struct {
	fixedCurrency string
}

func NewAirtime(currency string) *Airtime {
	return &Airtime{fixedCurrency: currency}
}

func (s *Airtime) Transfer(ctx context.Context, sender urns.URN, recipient urns.URN, amounts map[string]decimal.Decimal, logHTTP flows.HTTPLogCallback) (*flows.AirtimeTransfer, error) {
	transfer := &flows.AirtimeTransfer{
		UUID:      flows.AirtimeTransferUUID(uuids.NewV4()),
		Sender:    sender,
		Recipient: recipient,
		Currency:  "",
		Amount:    decimal.Zero,
	}

	if strings.Contains(string(recipient), "666") {
		return transfer, fmt.Errorf("invalid recipient number")
	}

	logHTTP(&flows.HTTPLog{
		HTTPLogWithoutTime: &flows.HTTPLogWithoutTime{
			LogWithoutTime: &httpx.LogWithoutTime{
				URL:        "http://send.airtime.com",
				StatusCode: 200,
				Request:    "GET / HTTP/1.1\r\nHost: send.airtime.com\r\nUser-Agent: Go-http-client/1.1\r\nAccept-Encoding: gzip\r\n\r\n",
				Response:   "HTTP/1.0 200 OK\r\nContent-Length: 15\r\n\r\n{\"status\":\"ok\"}",
				ElapsedMS:  0,
				Retries:    0,
			},
			Status: flows.CallStatusSuccess,
		},
		CreatedOn: time.Date(2019, 10, 16, 13, 59, 30, 123456789, time.UTC),
	})

	amount, hasAmount := amounts[s.fixedCurrency]
	if !hasAmount {
		return nil, fmt.Errorf("no amount configured for transfers in %s", s.fixedCurrency)
	}

	transfer.ExternalID = random.String(10, []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"))
	transfer.Currency = s.fixedCurrency
	transfer.Amount = amount

	return transfer, nil
}

var _ flows.AirtimeService = (*Airtime)(nil)
