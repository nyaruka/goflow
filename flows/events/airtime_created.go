package events

import (
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/flows"
	"github.com/shopspring/decimal"
)

func init() {
	registerType(TypeAirtimeCreated, func() flows.Event { return &AirtimeCreated{} })
}

// TypeAirtimeCreated is the type of our airtime created event
const TypeAirtimeCreated string = "airtime_created"

// AirtimeCreated events are created when an airtime transfer to the contact has been initiated. The transfer's
// final outcome is determined out of band by the host, so the event represents the request, not the delivery.
// `external_id` may be empty when the provider's transaction id isn't known at the time of event creation.
//
//	{
//	  "uuid": "0197b335-6ded-79a4-95a6-3af85b57f108",
//	  "type": "airtime_created",
//	  "created_on": "2006-01-02T15:04:05Z",
//	  "external_id": "",
//	  "sender": "tel:4748",
//	  "recipient": "tel:+1242563637",
//	  "currency": "RWF",
//	  "amount": 100,
//	  "http_logs": [
//	    {
//	      "url": "https://dvs-api.dtone.com/v1/lookup/mobile-number",
//	      "status": "success",
//	      "request": "POST /v1/lookup/mobile-number HTTP/1.1\r\n\r\n{}",
//	      "response": "HTTP/1.1 200 OK\r\n\r\n{}",
//	      "created_on": "2006-01-02T15:04:05Z",
//	      "elapsed_ms": 123
//	    }
//	  ]
//	}
//
// @event airtime_created
type AirtimeCreated struct {
	BaseEvent

	ExternalID string           `json:"external_id"`
	Sender     urns.URN         `json:"sender"`
	Recipient  urns.URN         `json:"recipient"`
	Currency   string           `json:"currency"`
	Amount     decimal.Decimal  `json:"amount"`
	HTTPLogs   []*flows.HTTPLog `json:"http_logs"`
}

// NewAirtimeCreated creates a new airtime created event
func NewAirtimeCreated(t *flows.AirtimeTransfer, httpLogs []*flows.HTTPLog) *AirtimeCreated {
	sender := t.Sender
	if sender != "" {
		sender = sender.Identity()
	}

	return &AirtimeCreated{
		BaseEvent:  NewBaseEvent(TypeAirtimeCreated),
		ExternalID: t.ExternalID,
		Sender:     sender,
		Recipient:  t.Recipient.Identity(),
		Currency:   t.Currency,
		Amount:     t.Amount,
		HTTPLogs:   httpLogs,
	}
}
