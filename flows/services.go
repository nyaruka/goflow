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
	ResponseJSON    []byte
	ResponseCleaned bool // whether response had to be cleaned to make it valid JSON
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
	Open(session Session, topic *Topic, body string, assignee *User, logHTTP HTTPLogCallback) (*Ticket, error)
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

// HTTPTrace describes an HTTP request/response
type HTTPTrace struct {
	URL        string     `json:"url" validate:"required"`
	StatusCode int        `json:"status_code,omitempty"`
	Status     CallStatus `json:"status" validate:"required"`
	Request    string     `json:"request" validate:"required"`
	Response   string     `json:"response,omitempty"`
	ElapsedMS  int        `json:"elapsed_ms"`
	Retries    int        `json:"retries"`
}

// trim request and response traces to 10K chars to avoid bloating serialized sessions
const trimTracesTo = 10000
const trimURLsTo = 2048

// NewHTTPTrace creates a new HTTP log from a trace
func NewHTTPTrace(trace *httpx.Trace, status CallStatus) *HTTPTrace {
	return newHTTPTraceWithStatus(trace, status, nil)
}

// HTTPLog describes an HTTP request/response
type HTTPLog struct {
	*HTTPTrace
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
func NewHTTPLog(trace *httpx.Trace, statusFn HTTPLogStatusResolver, redact utils.Redactor) *HTTPLog {
	return &HTTPLog{
		newHTTPTraceWithStatus(trace, statusFn(trace), redact),
		trace.StartTime,
	}
}

// creates a new HTTPTrace from a trace with an explicit status
func newHTTPTraceWithStatus(trace *httpx.Trace, status CallStatus, redact utils.Redactor) *HTTPTrace {
	url := trace.Request.URL.String()
	request := string(trace.RequestTrace)
	response := string(utils.ReplaceEscapedNulls(trace.SanitizedResponse("..."), []byte(`�`)))

	statusCode := 0
	if trace.Response != nil {
		statusCode = trace.Response.StatusCode
	}

	if redact != nil {
		url = redact(url)
		request = redact(request)
		response = redact(response)
	}

	return &HTTPTrace{
		URL:        utils.TruncateEllipsis(url, trimURLsTo),
		StatusCode: statusCode,
		Status:     status,
		Request:    utils.TruncateEllipsis(request, trimTracesTo),
		Response:   utils.TruncateEllipsis(response, trimTracesTo),
		ElapsedMS:  int((trace.EndTime.Sub(trace.StartTime)) / time.Millisecond),
		Retries:    trace.Retries,
	}
}
