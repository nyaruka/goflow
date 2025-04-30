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
//	{
//	  "uuid": "019688A6-41d2-7366-958a-630e35c62431",
//	  "type": "airtime_transferred",
//	  "created_on": "2006-01-02T15:04:05Z",
//	  "transfer_uuid": "552cd7ee-ccba-404d-9692-c1fe3b8d57c5",
//	  "external_id": "12345678",
//	  "sender": "tel:4748",
//	  "recipient": "tel:+1242563637",
//	  "currency": "RWF",
//	  "amount": 100,
//	  "http_logs": [
//	    {
//	      "url": "https://dvs-api.dtone.com/v1/sync/transactions",
//	      "status": "success",
//	      "request": "POST /v1/sync/transactions HTTP/1.1\r\n\r\n{}",
//	      "response": "HTTP/1.1 200 OK\r\n\r\n{}",
//	      "created_on": "2006-01-02T15:04:05Z",
//	      "elapsed_ms": 123
//	    }
//	  ]
//	}
//
// @event airtime_transferred
type AirtimeTransferredEvent struct {
	BaseEvent

	TransferUUID flows.AirtimeTransferUUID `json:"transfer_uuid"`
	ExternalID   string                    `json:"external_id"`
	Sender       urns.URN                  `json:"sender"`
	Recipient    urns.URN                  `json:"recipient"`
	Currency     string                    `json:"currency"`
	Amount       decimal.Decimal           `json:"amount"`
	HTTPLogs     []*flows.HTTPLog          `json:"http_logs"`
}

// NewAirtimeTransferred creates a new airtime transferred event
func NewAirtimeTransferred(t *flows.AirtimeTransfer, httpLogs []*flows.HTTPLog) *AirtimeTransferredEvent {
	return &AirtimeTransferredEvent{
		BaseEvent:    NewBaseEvent(TypeAirtimeTransferred),
		TransferUUID: t.UUID,
		ExternalID:   t.ExternalID,
		Sender:       t.Sender,
		Recipient:    t.Recipient,
		Currency:     t.Currency,
		Amount:       t.Amount,
		HTTPLogs:     httpLogs,
	}
}
