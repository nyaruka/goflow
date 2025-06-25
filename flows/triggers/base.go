package triggers

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/nyaruka/gocommon/dates"
	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

// ReadFunc is a function that can read a trigger from JSON
type ReadFunc func(flows.SessionAssets, []byte, assets.MissingCallback) (flows.Trigger, error)

var registeredTypes = map[string]ReadFunc{}

// registers a new type of trigger
func registerType(name string, f ReadFunc) {
	registeredTypes[name] = f
}

// RegisteredTypes gets the registered types of trigger
func RegisteredTypes() map[string]ReadFunc {
	return registeredTypes
}

// base of all trigger types
type baseTrigger struct {
	type_       string
	flow        *assets.FlowReference
	batch       bool
	params      *types.XObject
	history     *flows.SessionHistory
	triggeredOn time.Time
}

// create a new base trigger
func newBaseTrigger(typeName string, flow *assets.FlowReference, batch bool, history *flows.SessionHistory) baseTrigger {
	return baseTrigger{
		type_:       typeName,
		flow:        flow,
		batch:       batch,
		history:     history,
		triggeredOn: dates.Now(),
	}
}

func (t *baseTrigger) Type() string                   { return t.type_ }
func (t *baseTrigger) Event() flows.Event             { return nil }
func (t *baseTrigger) Flow() *assets.FlowReference    { return t.flow }
func (t *baseTrigger) Batch() bool                    { return t.batch }
func (t *baseTrigger) Params() *types.XObject         { return t.params }
func (t *baseTrigger) History() *flows.SessionHistory { return t.history }
func (t *baseTrigger) TriggeredOn() time.Time         { return t.triggeredOn }

// Initialize initializes the session
func (t *baseTrigger) Initialize(session flows.Session) error {
	return nil
}

//------------------------------------------------------------------------------------------
// Expressions context
//------------------------------------------------------------------------------------------

// Context is the schema of trigger objects in the context, across all types
type Context struct {
	type_    string
	params   *types.XObject
	keyword  string
	user     types.XValue
	origin   string
	campaign types.XValue
	ticket   types.XValue
	optIn    types.XValue
}

func (c *Context) asMap() map[string]types.XValue {
	return map[string]types.XValue{
		"type":     types.NewXText(c.type_),
		"params":   c.params,
		"keyword":  types.NewXText(c.keyword),
		"user":     c.user,
		"origin":   types.NewXText(c.origin),
		"campaign": c.campaign,
		"ticket":   c.ticket,
		"optin":    c.optIn,
	}
}

func (t *baseTrigger) context() *Context {
	params := t.params
	if params == nil {
		params = types.XObjectEmpty
	}

	return &Context{type_: t.type_, params: params}
}

// Context returns the properties available in expressions
//
//	type:text -> the type of trigger that started this session
//	params:any -> the parameters passed to the trigger
//	keyword:text -> the keyword match if this is a keyword trigger
//	user:user -> the user who started this session if this is a manual trigger
//	origin:text -> the origin of this session if this is a manual trigger
//	ticket:ticket -> the ticket if this is a ticket trigger
//
// @context trigger
func (t *baseTrigger) Context(env envs.Environment) map[string]types.XValue {
	return t.context().asMap()
}

//------------------------------------------------------------------------------------------
// Builder
//------------------------------------------------------------------------------------------

// Builder is a builder for triggers
type Builder struct {
	flow *assets.FlowReference
}

// NewBuilder creates a new trigger builder
func NewBuilder(flow *assets.FlowReference) *Builder {
	return &Builder{
		flow: flow,
	}
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type baseEnvelope struct {
	Type        string                `json:"type"               validate:"required"`
	Flow        *assets.FlowReference `json:"flow"               validate:"required"`
	Batch       bool                  `json:"batch,omitempty"`
	Params      json.RawMessage       `json:"params,omitempty"`
	History     *flows.SessionHistory `json:"history,omitempty"`
	TriggeredOn time.Time             `json:"triggered_on"       validate:"required"`
}

// ReadTrigger reads a trigger from the given JSON
func ReadTrigger(sa flows.SessionAssets, data []byte, missing assets.MissingCallback) (flows.Trigger, error) {
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

func (t *baseTrigger) unmarshal(sa flows.SessionAssets, e *baseEnvelope, missing assets.MissingCallback) error {
	var err error

	t.type_ = e.Type
	t.flow = e.Flow
	t.batch = e.Batch
	t.history = e.History
	t.triggeredOn = e.TriggeredOn

	if e.Params != nil {
		if t.params, err = types.ReadXObject(e.Params); err != nil {
			return fmt.Errorf("unable to read params: %w", err)
		}
	}

	return nil
}

func (t *baseTrigger) marshal(e *baseEnvelope) error {
	var err error
	e.Type = t.type_
	e.Flow = t.flow
	e.Batch = t.batch
	e.History = t.history
	e.TriggeredOn = t.triggeredOn

	if t.params != nil {
		e.Params, err = jsonx.Marshal(t.params)
		if err != nil {
			return err
		}
	}
	return nil
}
