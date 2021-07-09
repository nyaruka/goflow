package flows

import (
	"net/http"
	"time"

	"github.com/nyaruka/gocommon/httpx"
	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/utils"

	"github.com/shopspring/decimal"
)

// Services groups together interfaces for several services whose implementation is provided outside of the flow engine.
type Services interface {
	Email(Session) (EmailService, error)
	Webhook(Session) (WebhookService, error)
	Classification(Session, *Classifier) (ClassificationService, error)
	Ticket(Session, *Ticketer) (TicketService, error)
	Airtime(Session) (AirtimeService, error)
}

// EmailService provides email functionality to the engine
type EmailService interface {
	Send(session Session, addresses []string, subject, body string) error
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
	*httpx.Trace
	ResponseJSON []byte
}

// WebhookService provides webhook functionality to the engine
type WebhookService interface {
	Call(session Session, request *http.Request) (*WebhookCall, error)
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
	Classify(session Session, input string, logHTTP HTTPLogCallback) (*Classification, error)
}

// TicketService provides ticketing functionality to the engine
type TicketService interface {
	// Open tries to open a new ticket
	Open(session Session, subject, body string, logHTTP HTTPLogCallback) (*Ticket, error)
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
	UUID          uuids.UUID
	Sender        urns.URN
	Recipient     urns.URN
	Currency      string
	DesiredAmount decimal.Decimal
	ActualAmount  decimal.Decimal
}

// AirtimeService provides airtime functionality to the engine
type AirtimeService interface {
	// Transfer transfers airtime to the given URN
	Transfer(session Session, sender urns.URN, recipient urns.URN, amounts map[string]decimal.Decimal, logHTTP HTTPLogCallback) (*AirtimeTransfer, error)
}

// HTTPLog describes an HTTP request/response
type HTTPLog struct {
	URL       string     `json:"url" validate:"required"`
	Status    CallStatus `json:"status" validate:"required"`
	Request   string     `json:"request" validate:"required"`
	Response  string     `json:"response,omitempty"`
	CreatedOn time.Time  `json:"created_on" validate:"required"`
	ElapsedMS int        `json:"elapsed_ms"`
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

// HTTPStatusResolver is a function that determines the status of an HTTP log from the response
type HTTPStatusResolver func(t *httpx.Trace) CallStatus

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
func NewHTTPLog(trace *httpx.Trace, statusFn HTTPStatusResolver, redact utils.Redactor) *HTTPLog {
	return newHTTPLogWithStatus(trace, statusFn(trace), redact)
}

// creates a new HTTP log from a trace with an explicit status
func newHTTPLogWithStatus(trace *httpx.Trace, status CallStatus, redact utils.Redactor) *HTTPLog {
	url := trace.Request.URL.String()
	request := string(trace.RequestTrace)
	response := string(utils.ReplaceEscapedNulls(trace.SanitizedResponse("..."), []byte(`�`)))

	if redact != nil {
		url = redact(url)
		request = redact(request)
		response = redact(response)
	}

	return &HTTPLog{
		URL:       url,
		Status:    status,
		Request:   request,
		Response:  response,
		CreatedOn: trace.StartTime,
		ElapsedMS: int((trace.EndTime.Sub(trace.StartTime)) / time.Millisecond),
	}
}
