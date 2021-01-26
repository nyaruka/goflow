package waits

import (
	"encoding/json"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"

	"github.com/pkg/errors"
)

type readFunc func(data json.RawMessage) (flows.Wait, error)
type readActivatedFunc func(data json.RawMessage) (flows.ActivatedWait, error)

var registeredTypes = map[string]readFunc{}
var registeredActivatedTypes = map[string]readActivatedFunc{}

// RegisterType registers a new type of wait
func registerType(name string, f1 readFunc, f2 readActivatedFunc) {
	registeredTypes[name] = f1
	registeredActivatedTypes[name] = f2
}

type Timeout struct {
	Seconds_      int                `json:"seconds"       validate:"required"`
	CategoryUUID_ flows.CategoryUUID `json:"category_uuid" validate:"required,uuid4"`
}

func NewTimeout(seconds int, categoryUUID flows.CategoryUUID) *Timeout {
	return &Timeout{Seconds_: seconds, CategoryUUID_: categoryUUID}
}

func (t *Timeout) Seconds() int { return t.Seconds_ }

func (t *Timeout) CategoryUUID() flows.CategoryUUID { return t.CategoryUUID_ }

// the base of all wait types
type baseWait struct {
	type_ string

	timeout *Timeout
}

func newBaseWait(typeName string, timeout *Timeout) baseWait {
	return baseWait{type_: typeName, timeout: timeout}
}

// Type returns the type of this wait
func (w *baseWait) Type() string { return w.type_ }

// Timeout returns the timeout of this wait or nil if no timeout is set
func (w *baseWait) Timeout() flows.Timeout { return w.timeout }

func (w *baseWait) resumeTypeError(r flows.Resume) error {
	return errors.Errorf("can't end a wait of type '%s' with a resume of type '%s'", w.type_, r.Type())
}

type baseActivatedWait struct {
	type_          string
	timeoutSeconds *int
}

func (w *baseActivatedWait) Type() string { return w.type_ }

func (w *baseActivatedWait) TimeoutSeconds() *int { return w.timeoutSeconds }

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type baseWaitEnvelope struct {
	Type    string   `json:"type"              validate:"required"`
	Timeout *Timeout `json:"timeout,omitempty" validate:"omitempty,dive"`
}

// ReadWait reads a wait from the given JSON
func ReadWait(data json.RawMessage) (flows.Wait, error) {
	typeName, err := utils.ReadTypeFromJSON(data)
	if err != nil {
		return nil, err
	}

	f := registeredTypes[typeName]
	if f == nil {
		return nil, errors.Errorf("unknown type: '%s'", typeName)
	}
	return f(data)
}

func (w *baseWait) unmarshal(e *baseWaitEnvelope) error {
	w.type_ = e.Type
	w.timeout = e.Timeout
	return nil
}

func (w *baseWait) marshal(e *baseWaitEnvelope) error {
	e.Type = w.type_
	e.Timeout = w.timeout
	return nil
}

// ReadActivatedWait reads an activated wait from the given JSON
func ReadActivatedWait(data json.RawMessage) (flows.ActivatedWait, error) {
	typeName, err := utils.ReadTypeFromJSON(data)
	if err != nil {
		return nil, err
	}

	f := registeredActivatedTypes[typeName]
	if f == nil {
		return nil, errors.Errorf("unknown type: '%s'", typeName)
	}
	return f(data)
}

type baseActivatedWaitEnvelope struct {
	Type           string `json:"type" validate:"required"`
	TimeoutSeconds *int   `json:"timeout_seconds,omitempty"`
}

func (w *baseActivatedWait) unmarshal(e *baseActivatedWaitEnvelope) error {
	w.type_ = e.Type
	w.timeoutSeconds = e.TimeoutSeconds
	return nil
}

func (w *baseActivatedWait) marshal(e *baseActivatedWaitEnvelope) error {
	e.Type = w.type_
	e.TimeoutSeconds = w.timeoutSeconds
	return nil
}
