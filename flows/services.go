package flows

import (
	"net/http"
	"time"

	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/assets"

	"github.com/shopspring/decimal"
)

// Services groups together interfaces for several services whose implementation is provided outside of the flow engine.
type Services interface {
	Webhook(Session) WebhookProvider
	NLU(Session, assets.Classifier) NLUProvider
	Airtime(Session) AirtimeProvider
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

// WebhookProvider provides webhook calling functionality to the engine
type WebhookProvider interface {
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

// NLUClassification is the result of an NLU classification
type NLUClassification struct {
	Intents  []ExtractedIntent            `json:"intents,omitempty"`
	Entities map[string][]ExtractedEntity `json:"entities,omitempty"`
}

// NLUProvider provides NLU functionality to the engine
type NLUProvider interface {
	Classify(session Session, input string) (*NLUClassification, error)
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

// AirtimeProvider is the interface for an airtime transfer provider
type AirtimeProvider interface {
	// Transfer transfers airtime to the given URN
	Transfer(session Session, sender urns.URN, recipient urns.URN, amounts map[string]decimal.Decimal) (*AirtimeTransfer, error)
}
