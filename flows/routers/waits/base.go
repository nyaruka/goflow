package waits

import (
	"encoding/json"
	"time"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/resumes"
	"github.com/nyaruka/goflow/utils"

	"github.com/pkg/errors"
)

type readFunc func(data json.RawMessage) (flows.Wait, error)
type readActivatedFunc func(data json.RawMessage) (flows.ActivatedWait, error)

var registeredTypes = map[string]readFunc{}
var registeredActivatedTypes = map[string]readActivatedFunc{}

// RegisterType registers a new type of wait
func RegisterType(name string, f1 readFunc, f2 readActivatedFunc) {
	registeredTypes[name] = f1
	registeredActivatedTypes[name] = f2
}

// the base of all wait types
type baseWait struct {
	type_ string

	timeout *int
}

func newBaseWait(typeName string, timeout *int) baseWait {
	return baseWait{type_: typeName, timeout: timeout}
}

// Type returns the type of this wait
func (w *baseWait) Type() string { return w.type_ }

// Timeout returns the timeout of this wait in seconds or nil if no timeout is set
func (w *baseWait) Timeout() *int { return w.timeout }

type baseActivatedWait struct {
	type_     string
	timeoutOn *time.Time
}

func (w *baseActivatedWait) Type() string { return w.type_ }

func (w *baseActivatedWait) TimeoutOn() *time.Time { return w.timeoutOn }

// End ends this wait or returns an error
func (w *baseActivatedWait) End(resume flows.Resume, node flows.Node) error {
	switch resume.Type() {
	case resumes.TypeRunExpiration:
		// expired runs always end a wait
		return nil
	case resumes.TypeWaitTimeout:
		if node.Router() == nil || node.Router().Wait() == nil || node.Router().Wait().Timeout() == nil {
			return errors.Errorf("can't end with timeout as node no longer has a wait timeout")
		}
		if w.TimeoutOn() == nil {
			return errors.Errorf("can't end with timeout as session wait has no timeout")
		}
		if utils.Now().Before(*w.TimeoutOn()) {
			return errors.Errorf("can't end with timeout before wait has timed out")
		}
	}
	return nil
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type baseWaitEnvelope struct {
	Type    string `json:"type" validate:"required"`
	Timeout *int   `json:"timeout,omitempty"`
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
	Type      string     `json:"type" validate:"required"`
	TimeoutOn *time.Time `json:"timeout_on,omitempty"`
}

func (w *baseActivatedWait) unmarshal(e *baseActivatedWaitEnvelope) error {
	w.type_ = e.Type
	w.timeoutOn = e.TimeoutOn
	return nil
}

func (w *baseActivatedWait) marshal(e *baseActivatedWaitEnvelope) error {
	e.Type = w.type_
	e.TimeoutOn = w.timeoutOn
	return nil
}
