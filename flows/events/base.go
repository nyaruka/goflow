package events

import (
	"encoding/json"
	"time"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
	"github.com/nyaruka/goflow/utils/dates"
	"github.com/nyaruka/goflow/utils/httpx"

	"github.com/pkg/errors"
)

var registeredTypes = map[string](func() flows.Event){}

// registers a new type of event
func registerType(name string, initFunc func() flows.Event) {
	registeredTypes[name] = initFunc
}

// base of all event types
type baseEvent struct {
	Type_      string         `json:"type" validate:"required"`
	CreatedOn_ time.Time      `json:"created_on" validate:"required"`
	StepUUID_  flows.StepUUID `json:"step_uuid,omitempty" validate:"omitempty,uuid4"`
}

// creates a new base event
func newBaseEvent(typeName string) baseEvent {
	return baseEvent{Type_: typeName, CreatedOn_: dates.Now()}
}

// Type returns the type of this event
func (e *baseEvent) Type() string { return e.Type_ }

// CreatedOn returns the created on time of this event
func (e *baseEvent) CreatedOn() time.Time { return e.CreatedOn_ }

// StepUUID returns the UUID of the step in the path where this event occured
func (e *baseEvent) StepUUID() flows.StepUUID { return e.StepUUID_ }

// SetStepUUID sets the UUID of the step in the path where this event occured
func (e *baseEvent) SetStepUUID(stepUUID flows.StepUUID) { e.StepUUID_ = stepUUID }

// utility for events which describe calls to external services
type externalCallEvent struct {
	URL       string           `json:"url" validate:"required"`
	Status    flows.CallStatus `json:"status" validate:"required"`
	Request   string           `json:"request" validate:"required"`
	Response  string           `json:"response,omitempty"`
	ElapsedMS int              `json:"elapsed_ms"`
}

//------------------------------------------------------------------------------------------
// HTTP logging
//------------------------------------------------------------------------------------------

// HTTPLog describes an HTTP request/response
type HTTPLog struct {
	URL       string           `json:"url" validate:"required"`
	Status    flows.CallStatus `json:"status" validate:"required"`
	Request   string           `json:"request" validate:"required"`
	Response  string           `json:"response,omitempty"`
	CreatedOn time.Time        `json:"created_on" validate:"required"`
	ElapsedMS int              `json:"elapsed_ms"`
}

// creates a new HTTP log from a trace
func httpLogFromTrace(trace *httpx.Trace) *HTTPLog {
	status := flows.CallStatusSuccess
	if trace.Response == nil {
		status = flows.CallStatusConnectionError
	} else if trace.Response.StatusCode >= 400 {
		status = flows.CallStatusResponseError
	}

	return &HTTPLog{
		URL:       trace.Request.URL.String(),
		Status:    status,
		Request:   string(trace.RequestTrace),
		Response:  string(trace.ResponseTrace),
		CreatedOn: trace.StartTime,
		ElapsedMS: int((trace.EndTime.Sub(trace.StartTime)) / time.Millisecond),
	}
}

func httpLogsFromTraces(traces []*httpx.Trace) []*HTTPLog {
	logs := make([]*HTTPLog, len(traces))
	for i := range traces {
		logs[i] = httpLogFromTrace(traces[i])
	}
	return logs
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

// ReadEvent reads a single event from the given JSON
func ReadEvent(data json.RawMessage) (flows.Event, error) {
	typeName, err := utils.ReadTypeFromJSON(data)
	if err != nil {
		return nil, err
	}

	f := registeredTypes[typeName]
	if f == nil {
		return nil, errors.Errorf("unknown type: '%s'", typeName)
	}

	event := f()
	return event, utils.UnmarshalAndValidate(data, event)
}
