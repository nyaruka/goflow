package flows

import (
	"net/http"
	"time"

	"github.com/nyaruka/gocommon/urns"

	"github.com/shopspring/decimal"
)

// Services groups together interfaces for several services whose implementation is provided outside of the flow engine.
type Services interface {
	Webhook(Session) WebhookService
	Airtime(Session) AirtimeService
}

// WebhookStatus represents the status of a webhook call
type WebhookStatus string

const (
	// WebhookStatusSuccess represents that the webhook was successful
	WebhookStatusSuccess WebhookStatus = "success"

	// WebhookStatusConnectionError represents that the webhook had a connection error
	WebhookStatusConnectionError WebhookStatus = "connection_error"

	// WebhookStatusResponseError represents that the webhook response had a non 2xx status code
	WebhookStatusResponseError WebhookStatus = "response_error"

	// WebhookStatusSubscriberGone represents a special state of resthook responses which indicate the caller must remove that subscriber
	WebhookStatusSubscriberGone WebhookStatus = "subscriber_gone"
)

// WebhookCall is the result of a webhook call
type WebhookCall struct {
	URL         string
	Method      string
	StatusCode  int
	Status      WebhookStatus
	TimeTaken   time.Duration
	Request     []byte
	Response    []byte
	BodyIgnored bool
	Resthook    string
}

// WebhookService provides webhook functionality to the engine
type WebhookService interface {
	Call(session Session, request *http.Request, resthook string) (*WebhookCall, error)
}

// AirtimeTransferStatus is a status of a airtime transfer
type AirtimeTransferStatus string

// possible values for airtime transfer statuses
const (
	AirtimeTransferStatusSuccess AirtimeTransferStatus = "success"
	AirtimeTransferStatusFailed  AirtimeTransferStatus = "failed"
)

// AirtimeTransfer is the result of an attempted airtime transfer
type AirtimeTransfer struct {
	Sender        urns.URN
	Recipient     urns.URN
	Currency      string
	DesiredAmount decimal.Decimal
	ActualAmount  decimal.Decimal
	Status        AirtimeTransferStatus
}

// AirtimeService provides airtime functionality to the engine
type AirtimeService interface {
	// Transfer transfers airtime to the given URN
	Transfer(session Session, sender urns.URN, recipient urns.URN, amounts map[string]decimal.Decimal) (*AirtimeTransfer, error)
}
