package waits

import (
	"encoding/json"
	"time"

	"github.com/nyaruka/gocommon/dates"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"

	"github.com/pkg/errors"
)

type readFunc func(data json.RawMessage) (flows.Wait, error)

var registeredTypes = map[string]readFunc{}

// RegisterType registers a new type of wait
func registerType(name string, f readFunc) {
	registeredTypes[name] = f
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

func (w *baseWait) expiresOn(run flows.Run) *time.Time {
	expiresAfterMins := run.Flow().ExpireAfterMinutes()
	if expiresAfterMins > 0 {
		dt := dates.Now().Add(time.Duration(expiresAfterMins * int(time.Minute)))
		return &dt
	}
	return nil
}

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
