package waits

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/nyaruka/gocommon/dates"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

type readFunc func(data json.RawMessage) (flows.Wait, error)

var registeredTypes = map[string]readFunc{}

// RegisterType registers a new type of wait
func registerType(name string, f readFunc) {
	registeredTypes[name] = f
}

type Timeout struct {
	Seconds_      int                `json:"seconds"       validate:"required"`
	CategoryUUID_ flows.CategoryUUID `json:"category_uuid" validate:"required,uuid"`
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
func (w *baseWait) Timeout() flows.Timeout {
	if w.timeout == nil {
		return nil
	}
	return w.timeout
}

func (w *baseWait) expiresOn(run flows.Run) time.Time {
	return dates.Now().Add(run.Flow().ExpireAfter())
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type baseEnvelope struct {
	Type    string   `json:"type"              validate:"required"`
	Timeout *Timeout `json:"timeout,omitempty" validate:"omitempty"`
}

// ReadWait reads a wait from the given JSON
func ReadWait(data []byte) (flows.Wait, error) {
	typeName, err := utils.ReadTypeFromJSON(data)
	if err != nil {
		return nil, err
	}

	f := registeredTypes[typeName]
	if f == nil {
		return nil, fmt.Errorf("unknown type: '%s'", typeName)
	}
	return f(data)
}

func (w *baseWait) unmarshal(e *baseEnvelope) error {
	w.type_ = e.Type
	w.timeout = e.Timeout
	return nil
}

func (w *baseWait) marshal(e *baseEnvelope) error {
	e.Type = w.type_
	e.Timeout = w.timeout
	return nil
}
