package events

import (
	"github.com/nyaruka/goflow/flows"

	"github.com/shopspring/decimal"
)

func init() {
	RegisterType(TypeAirtimeTransferred, func() flows.Event { return &AirtimeTransferredEvent{} })
}

// TypeAirtimeTransferred is the type of our airtime transferred event
const TypeAirtimeTransferred string = "airtime_transferred"

// AirtimeTransferredEvent events are created when airtime has been transferred to the contact.
//
//   {
//     "type": "airtime_transferred",
//     "created_on": "2006-01-02T15:04:05Z",
//     "currency": "RWF",
//     "amount": 100,
//     "status": "success"
//   }
//
// @event airtime_transferred
type AirtimeTransferredEvent struct {
	BaseEvent

	Currency string          `json:"currency"`
	Amount   decimal.Decimal `json:"amount"`
	Status   string          `json:"status"`
}

// NewAirtimeTransferredEvent creates a new airtime transferred event
func NewAirtimeTransferredEvent(t *flows.AirtimeTransfer) *AirtimeTransferredEvent {
	return &AirtimeTransferredEvent{
		BaseEvent: NewBaseEvent(TypeAirtimeTransferred),
		Currency:  t.Currency,
		Amount:    t.ActualAmount,
		Status:    string(t.Status),
	}
}
