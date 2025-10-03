package events

import (
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/flows"
	"github.com/shopspring/decimal"
)

func init() {
	registerType(TypeAirtimeTransferred, func() flows.Event { return &AirtimeTransferred{} })
}

// TypeAirtimeTransferred is the type of our airtime transferred event
const TypeAirtimeTransferred string = "airtime_transferred"

// AirtimeTransferred events are created when and airtime transfer to the contact has been initiated.
//
//	{
//	  "uuid": "0197b335-6ded-79a4-95a6-3af85b57f108",
//	  "type": "airtime_transferred",
//	  "created_on": "2006-01-02T15:04:05Z",
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
type AirtimeTransferred struct {
	BaseEvent

	ExternalID string           `json:"external_id"`
	Sender     urns.URN         `json:"sender"`
	Recipient  urns.URN         `json:"recipient"`
	Currency   string           `json:"currency"`
	Amount     decimal.Decimal  `json:"amount"`
	HTTPLogs   []*flows.HTTPLog `json:"http_logs"`
}

// NewAirtimeTransferred creates a new airtime transferred event
func NewAirtimeTransferred(t *flows.AirtimeTransfer, httpLogs []*flows.HTTPLog) *AirtimeTransferred {
	sender := t.Sender
	if sender != "" {
		sender = sender.Identity()
	}

	return &AirtimeTransferred{
		BaseEvent:  NewBaseEvent(TypeAirtimeTransferred),
		ExternalID: t.ExternalID,
		Sender:     sender,
		Recipient:  t.Recipient.Identity(),
		Currency:   t.Currency,
		Amount:     t.Amount,
		HTTPLogs:   httpLogs,
	}
}
