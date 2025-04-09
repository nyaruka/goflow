package flows

import (
	"context"
	"net/http"
	"time"

	"github.com/nyaruka/gocommon/httpx"
	"github.com/nyaruka/gocommon/stringsx"
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/envs"

	"github.com/shopspring/decimal"
)

// Services groups together interfaces for several services whose implementation is provided outside of the flow engine.
type Services interface {
	Email(SessionAssets) (EmailService, error)
	Webhook(SessionAssets) (WebhookService, error)
	Classification(*Classifier) (ClassificationService, error)
	LLM(*LLM) (LLMService, error)
	Airtime(SessionAssets) (AirtimeService, error)
}

// EmailService provides email functionality to the engine
type EmailService interface {
	Send(addresses []string, subject, body string) error
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

// WebhookService provides webhook functionality to the engine
type WebhookService interface {
	Call(request *http.Request) (*httpx.Trace, error)
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
	Classify(ctx context.Context, env envs.Environment, input string, logHTTP HTTPLogCallback) (*Classification, error)
}

type LLMResponse struct {
	Output     string
	TokensUsed int64
}

// LLMService provides LLM functionality to the engine
type LLMService interface {
	Response(ctx context.Context, instructions, input string, maxTokens int) (*LLMResponse, error)
}

// AirtimeTransferUUID is the UUID of a airtime transfer
type AirtimeTransferUUID uuids.UUID

// AirtimeTransferStatus is a status of a airtime transfer
type AirtimeTransferStatus string

// possible values for airtime transfer statuses
const (
	AirtimeTransferStatusSuccess AirtimeTransferStatus = "success"
	AirtimeTransferStatusFailed  AirtimeTransferStatus = "failed"
)

// AirtimeTransfer is the result of an attempted airtime transfer
type AirtimeTransfer struct {
	UUID       AirtimeTransferUUID
	ExternalID string
	Sender     urns.URN
	Recipient  urns.URN
	Currency   string
	Amount     decimal.Decimal
}

// AirtimeService provides airtime functionality to the engine
type AirtimeService interface {
	// Transfer transfers airtime to the given URN
	Transfer(ctx context.Context, sender urns.URN, recipient urns.URN, amounts map[string]decimal.Decimal, logHTTP HTTPLogCallback) (*AirtimeTransfer, error)
}

// HTTPLogWithoutTime is an HTTP log no time and status added - used for webhook events which already encode the time
type HTTPLogWithoutTime struct {
	*httpx.LogWithoutTime

	Status CallStatus `json:"status" validate:"required"`
}

// trim request and response traces to 10K chars to avoid bloating serialized sessions
const trimTracesTo = 10000
const trimURLsTo = 2048

// NewHTTPLogWithoutTime creates a new HTTP log from a trace
func NewHTTPLogWithoutTime(trace *httpx.Trace, status CallStatus, redact stringsx.Redactor) *HTTPLogWithoutTime {
	return &HTTPLogWithoutTime{
		LogWithoutTime: httpx.NewLogWithoutTime(trace, trimURLsTo, trimTracesTo, redact),
		Status:         status,
	}
}

// HTTPLog describes an HTTP request/response
type HTTPLog struct {
	*HTTPLogWithoutTime
	CreatedOn time.Time `json:"created_on" validate:"required"`
}

// HTTPLogCallback is a function that handles an HTTP log
type HTTPLogCallback func(*HTTPLog)

// HTTPLogger logs HTTP logs
type HTTPLogger struct {
	Logs []*HTTPLog
}

// Log logs the given HTTP log
func (l *HTTPLogger) Log(h *HTTPLog) {
	l.Logs = append(l.Logs, h)
}

// HTTPLogStatusResolver is a function that determines the status of an HTTP log from the response
type HTTPLogStatusResolver func(t *httpx.Trace) CallStatus

// HTTPStatusFromCode uses the status code to determine status of an HTTP log
func HTTPStatusFromCode(t *httpx.Trace) CallStatus {
	if t.Response == nil {
		return CallStatusConnectionError
	} else if t.Response.StatusCode >= 400 {
		return CallStatusResponseError
	}
	return CallStatusSuccess
}

// RedactionMask is the redaction mask for HTTP service logs
const RedactionMask = "****************"

// NewHTTPLog creates a new HTTP log from a trace
func NewHTTPLog(trace *httpx.Trace, statusFn HTTPLogStatusResolver, redact stringsx.Redactor) *HTTPLog {
	return &HTTPLog{
		HTTPLogWithoutTime: NewHTTPLogWithoutTime(trace, statusFn(trace), redact),
		CreatedOn:          trace.StartTime,
	}
}
