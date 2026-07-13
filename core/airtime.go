package core

import (
	"github.com/nyaruka/gocommon/urns"
	"github.com/shopspring/decimal"
)

// AirtimeTransfer is the result of an attempted airtime transfer
type AirtimeTransfer struct {
	ExternalID string // provider transaction id, when known after initiation
	Sender     urns.URN
	Recipient  urns.URN
	Currency   string
	Amount     decimal.Decimal
}
