package flows

import (
	"net/http"
	"time"

	"github.com/nyaruka/gocommon/urns"

	"github.com/shopspring/decimal"
)

// Services groups together interfaces for several services whose implementation is provided outside of the flow engine.
type Services interface {
	Webhook(Session) (WebhookService, error)
	Classification(Session, *Classifier) (ClassificationService, error)
	Airtime(Session) (AirtimeService, error)
}

// CallStatus represents the status of a call to an external service
type CallStatus string

const (
	// CallStatusSuccess represents that the webhook was successful
	CallStatusSuccess CallStatus = "success"

	// CallStatusConnectionError represents that the webhook had a connection error
	CallStatusConnectionError CallStatus = "connection_error"

	// CallStatusResponseError represents that the webhook response had a non 2xx status code
	CallStatusResponseError CallStatus = "response_error"

	// CallStatusSubscriberGone represents a special state of resthook responses which indicate the caller must remove that subscriber
	CallStatusSubscriberGone CallStatus = "subscriber_gone"
)

// WebhookCall is the result of a webhook call
type WebhookCall struct {
	URL         string
	Method      string
	StatusCode  int
	Status      CallStatus
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

// ExtractedIntent models an intent match
type ExtractedIntent struct {
	Name       string          `json:"name"`
	Confidence decimal.Decimal `json:"confidence"`
}

// ExtractedEntity models an entity match
type ExtractedEntity struct {
	Value      string          `json:"value"`
	Confidence decimal.Decimal `json:"confidence"`
}

// Classification is the result of an NLU classification
type Classification struct {
	Intents  []ExtractedIntent            `json:"intents,omitempty"`
	Entities map[string][]ExtractedEntity `json:"entities,omitempty"`
}

// ClassificationService provides NLU functionality to the engine
type ClassificationService interface {
	Classify(session Session, input string, logEvent EventCallback) (*Classification, error)
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
	Sender    urns.URN
	Recipient urns.URN
	Currency  string
	Amount    decimal.Decimal
}

// AirtimeService provides airtime functionality to the engine
type AirtimeService interface {
	// Transfer transfers airtime to the given URN
	Transfer(session Session, sender urns.URN, recipient urns.URN, amounts map[string]decimal.Decimal, logEvent EventCallback) (*AirtimeTransfer, error)
}
