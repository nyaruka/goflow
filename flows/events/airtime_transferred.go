package events

import (
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/flows"

	"github.com/shopspring/decimal"
)

func init() {
	registerType(TypeAirtimeTransferred, func() flows.Event { return &AirtimeTransferredEvent{} })
}

// TypeAirtimeTransferred is the type of our airtime transferred event
const TypeAirtimeTransferred string = "airtime_transferred"

// AirtimeTransferredEvent events are created when airtime has been transferred to the contact.
//
//   {
//     "type": "airtime_transferred",
//     "created_on": "2006-01-02T15:04:05Z",
//     "sender": "tel:4748",
//     "recipient": "tel:+1242563637",
//     "currency": "RWF",
//     "amount": 100
//   }
//
// @event airtime_transferred
type AirtimeTransferredEvent struct {
	baseEvent

	Sender    urns.URN        `json:"sender"`
	Recipient urns.URN        `json:"recipient"`
	Currency  string          `json:"currency"`
	Amount    decimal.Decimal `json:"amount"`
}

// NewAirtimeTransferred creates a new airtime transferred event
func NewAirtimeTransferred(t *flows.AirtimeTransfer) *AirtimeTransferredEvent {
	return &AirtimeTransferredEvent{
		baseEvent: newBaseEvent(TypeAirtimeTransferred),
		Sender:    t.Sender,
		Recipient: t.Recipient,
		Currency:  t.Currency,
		Amount:    t.Amount,
	}
}
