package core

import (
	"time"

	"github.com/nyaruka/gocommon/httpx"
	"github.com/nyaruka/gocommon/stringsx"
)

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

// HTTPSizes are the true sizes in bytes of a request and response, recorded because the traces
// themselves are trimmed for logging
type HTTPSizes struct {
	Request  int `json:"request"`
	Response int `json:"response"`
}

// HTTPLogWithoutTime is an HTTP log no time and status added - used for webhook events which already encode the time
type HTTPLogWithoutTime struct {
	*httpx.LogWithoutTime

	Status CallStatus `json:"status" validate:"required"`
	Sizes  HTTPSizes  `json:"sizes"`
}

// trim request and response traces to 10K chars to avoid bloating serialized sessions
const trimTracesTo = 10000
const trimURLsTo = 2048

// NewHTTPLogWithoutTime creates a new HTTP log from a trace
func NewHTTPLogWithoutTime(trace *httpx.Trace, status CallStatus, redact stringsx.Redactor) *HTTPLogWithoutTime {
	return &HTTPLogWithoutTime{
		LogWithoutTime: httpx.NewLogWithoutTime(trace, trimURLsTo, trimTracesTo, redact),
		Status:         status,
		Sizes: HTTPSizes{
			Request:  len(trace.RequestTrace),
			Response: len(trace.ResponseTrace) + len(trace.ResponseBody),
		},
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
