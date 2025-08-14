package resumes

import (
	"fmt"
	"time"

	"github.com/nyaruka/gocommon/dates"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

// ReadFunc is a function that can read a resume from JSON
type ReadFunc func(flows.SessionAssets, []byte, assets.MissingCallback) (flows.Resume, error)

var registeredTypes = map[string]ReadFunc{}

// registers a new type of resume
func registerType(name string, f ReadFunc) {
	registeredTypes[name] = f
}

// RegisteredTypes gets the registered types of resumes
func RegisteredTypes() map[string]ReadFunc {
	return registeredTypes
}

// base of all resume types
type baseResume struct {
	type_     string
	resumedOn time.Time
}

// creates a new base resume
func newBaseResume(typeName string) baseResume {
	return baseResume{type_: typeName, resumedOn: dates.Now()}
}

// Type returns the type of this resume
func (r *baseResume) Type() string         { return r.type_ }
func (r *baseResume) Event() flows.Event   { return nil }
func (r *baseResume) ResumedOn() time.Time { return r.resumedOn }

// Apply applies our state changes and saves any events to the run
func (r *baseResume) Apply(run flows.Run, logEvent flows.EventCallback) {
	if run.Status() == flows.RunStatusWaiting {
		run.SetStatus(flows.RunStatusActive)
	}
}

func (r *baseResume) Input(flows.SessionAssets) flows.Input { return nil }

//------------------------------------------------------------------------------------------
// Expressions context
//------------------------------------------------------------------------------------------

// Context is the schema of trigger objects in the context, across all types
type Context struct {
	type_ string
	dial  types.XValue
}

func (c *Context) asMap() map[string]types.XValue {
	return map[string]types.XValue{
		"type": types.NewXText(c.type_),
		"dial": c.dial,
	}
}

func (r *baseResume) context() *Context {
	return &Context{type_: r.type_}
}

// Context returns the properties available in expressions
//
//	type:text -> the type of resume that resumed this session
//
// @context resume
func (r *baseResume) Context(env envs.Environment) map[string]types.XValue {
	return r.context().asMap()
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type baseEnvelope struct {
	Type      string    `json:"type" validate:"required"`
	ResumedOn time.Time `json:"resumed_on" validate:"required"`
}

// Read reads a resume from the given JSON
func Read(sa flows.SessionAssets, data []byte, missing assets.MissingCallback) (flows.Resume, error) {
	typeName, err := utils.ReadTypeFromJSON(data)
	if err != nil {
		return nil, err
	}

	f := registeredTypes[typeName]
	if f == nil {
		return nil, fmt.Errorf("unknown type: '%s'", typeName)
	}
	return f(sa, data, missing)
}

func (r *baseResume) unmarshal(sa flows.SessionAssets, e *baseEnvelope, missing assets.MissingCallback) error {
	r.type_ = e.Type
	r.resumedOn = e.ResumedOn
	return nil
}

func (r *baseResume) marshal(e *baseEnvelope) error {
	e.Type = r.type_
	e.ResumedOn = r.resumedOn
	return nil
}
