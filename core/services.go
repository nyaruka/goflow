package core

import (
	"github.com/nyaruka/gocommon/urns"
	"github.com/shopspring/decimal"
)

// LLMResponse is the response from an LLM service call
type LLMResponse struct {
	Output       string
	TokensInput  int64
	TokensOutput int64
}

// AirtimeTransfer is the result of an attempted airtime transfer
type AirtimeTransfer struct {
	ExternalID string // provider transaction id, when known after initiation
	Sender     urns.URN
	Recipient  urns.URN
	Currency   string
	Amount     decimal.Decimal
}
