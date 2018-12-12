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

var registeredTypes = map[string]readFunc{}

// RegisterType registers a new type of wait
func RegisterType(name string, f readFunc) {
	registeredTypes[name] = f
}

// the base of all wait types
type baseWait struct {
	type_ string

	timeout   *int
	timeoutOn *time.Time
}

func newBaseWait(typeName string, timeout *int) baseWait {
	return baseWait{type_: typeName, timeout: timeout}
}

// Type returns the type of this wait
func (w *baseWait) Type() string { return w.type_ }

// Timeout returns the timeout of this wait in seconds or nil if no timeout is set
func (w *baseWait) Timeout() *int { return w.timeout }

// TimeoutOn returns when this wait times out
func (w *baseWait) TimeoutOn() *time.Time { return w.timeoutOn }

// Begin beings waiting
func (w *baseWait) Begin(run flows.FlowRun) bool {
	if w.timeout != nil {
		timeoutOn := utils.Now().Add(time.Second * time.Duration(*w.timeout))

		w.timeoutOn = &timeoutOn
	}
	return true
}

// End ends this wait or returns an error
func (w *baseWait) End(resume flows.Resume, node flows.Node) error {
	switch resume.Type() {
	case resumes.TypeRunExpiration:
		// expired runs always end a wait
		return nil
	case resumes.TypeWaitTimeout:
		if node.Wait().Timeout() == nil {
			return errors.Errorf("can't end with timeout as node no longer has a wait timeout")
		}
		if w.Timeout() == nil || w.TimeoutOn() == nil {
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
	Type      string     `json:"type" validate:"required"`
	Timeout   *int       `json:"timeout,omitempty"`
	TimeoutOn *time.Time `json:"timeout_on,omitempty"`
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
	w.timeoutOn = e.TimeoutOn
	return nil
}

func (w *baseWait) marshal(e *baseWaitEnvelope) error {
	e.Type = w.type_
	e.Timeout = w.timeout
	e.TimeoutOn = w.timeoutOn
	return nil
}
