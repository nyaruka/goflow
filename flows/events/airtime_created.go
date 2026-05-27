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

// AirtimeCreated events are created when a transfer_airtime action successfully initiates an airtime transfer.
// The final outcome of the transfer is determined out of band by the host. `external_id` carries the
// provider's transaction id when it's known at initiation time, otherwise it's empty.
//
//	{
//	  "uuid": "0197b335-6ded-79a4-95a6-3af85b57f108",
//	  "type": "airtime_created",
//	  "created_on": "2006-01-02T15:04:05Z",
//	  "external_id": "2237512891",
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

// NewAirtimeCreated creates a new airtime created event with the given pre-allocated UUID
func NewAirtimeCreated(uuid flows.EventUUID, t *flows.AirtimeTransfer, httpLogs []*flows.HTTPLog) *AirtimeCreated {
	sender := t.Sender
	if sender != "" {
		sender = sender.Identity()
	}

	return &AirtimeCreated{
		BaseEvent:  NewBaseEventWithUUID(uuid, TypeAirtimeCreated),
		ExternalID: t.ExternalID,
		Sender:     sender,
		Recipient:  t.Recipient.Identity(),
		Currency:   t.Currency,
		Amount:     t.Amount,
		HTTPLogs:   httpLogs,
	}
}
