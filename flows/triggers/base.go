package triggers

import (
	"encoding/json"
	"time"

	"github.com/nyaruka/gocommon/dates"
	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"

	"github.com/pkg/errors"
)

// ReadFunc is a function that can read a trigger from JSON
type ReadFunc func(flows.SessionAssets, json.RawMessage, assets.MissingCallback) (flows.Trigger, error)

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
	environment envs.Environment
	flow        *assets.FlowReference
	contact     *flows.Contact
	connection  *flows.Connection
	batch       bool
	params      *types.XObject
	history     *flows.SessionHistory
	triggeredOn time.Time
}

// create a new base trigger
func newBaseTrigger(typeName string, env envs.Environment, flow *assets.FlowReference, contact *flows.Contact, connection *flows.Connection, batch bool, history *flows.SessionHistory) baseTrigger {
	return baseTrigger{
		type_:       typeName,
		environment: env,
		flow:        flow,
		contact:     contact,
		connection:  connection,
		batch:       batch,
		history:     history,
		triggeredOn: dates.Now(),
	}
}

// Type returns the type of this trigger
func (t *baseTrigger) Type() string { return t.type_ }

func (t *baseTrigger) Environment() envs.Environment  { return t.environment }
func (t *baseTrigger) Flow() *assets.FlowReference    { return t.flow }
func (t *baseTrigger) Contact() *flows.Contact        { return t.contact }
func (t *baseTrigger) Connection() *flows.Connection  { return t.connection }
func (t *baseTrigger) Batch() bool                    { return t.batch }
func (t *baseTrigger) Params() *types.XObject         { return t.params }
func (t *baseTrigger) History() *flows.SessionHistory { return t.history }
func (t *baseTrigger) TriggeredOn() time.Time         { return t.triggeredOn }

// Initialize initializes the session
func (t *baseTrigger) Initialize(session flows.Session, logEvent flows.EventCallback) error {
	// try to load the flow
	flow, err := session.Assets().Flows().Get(t.Flow().UUID)
	if err != nil {
		return errors.Wrapf(err, "unable to load %s", t.Flow())
	}

	if flow.Type() == flows.FlowTypeVoice && t.connection == nil {
		return errors.New("unable to trigger voice flow without connection")
	}

	session.SetType(flow.Type())
	session.PushFlow(flow, nil, false)

	if t.environment != nil {
		session.SetEnvironment(t.environment)
	} else {
		session.SetEnvironment(envs.NewBuilder().Build())
	}

	if t.contact != nil {
		session.SetContact(t.contact.Clone())
	}
	return nil
}

// InitializeRun performs additional initialization when we create our first run
func (t *baseTrigger) InitializeRun(run flows.FlowRun, logEvent flows.EventCallback) error {
	return nil
}

//------------------------------------------------------------------------------------------
// Expressions context
//------------------------------------------------------------------------------------------

// Context is the schema of trigger objects in the context, across all types
type Context struct {
	type_   string
	params  *types.XObject
	keyword string
	user    types.XValue
	origin  string
	ticket  types.XValue
}

func (c *Context) asMap() map[string]types.XValue {
	return map[string]types.XValue{
		"type":    types.NewXText(c.type_),
		"params":  c.params,
		"keyword": types.NewXText(c.keyword),
		"user":    c.user,
		"origin":  types.NewXText(c.origin),
		"ticket":  c.ticket,
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
//   type:text -> the type of trigger that started this session
//   params:any -> the parameters passed to the trigger
//   keyword:text -> the keyword match if this is a keyword trigger
//   user:text -> the user who started this session if this is a manual trigger
//   origin:text -> the origin of this session if this is a manual trigger
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
	environment envs.Environment
	flow        *assets.FlowReference
	contact     *flows.Contact
}

// NewBuilder creates a new trigger builder
func NewBuilder(env envs.Environment, flow *assets.FlowReference, contact *flows.Contact) *Builder {
	return &Builder{
		environment: env,
		flow:        flow,
		contact:     contact,
	}
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type baseTriggerEnvelope struct {
	Type        string                `json:"type" validate:"required"`
	Environment json.RawMessage       `json:"environment,omitempty"`
	Flow        *assets.FlowReference `json:"flow" validate:"required"`
	Contact     json.RawMessage       `json:"contact,omitempty"`
	Connection  *flows.Connection     `json:"connection,omitempty"`
	Batch       bool                  `json:"batch,omitempty"`
	Params      json.RawMessage       `json:"params,omitempty"`
	History     *flows.SessionHistory `json:"history,omitempty"`
	TriggeredOn time.Time             `json:"triggered_on" validate:"required"`
}

// ReadTrigger reads a trigger from the given JSON
func ReadTrigger(sessionAssets flows.SessionAssets, data json.RawMessage, missing assets.MissingCallback) (flows.Trigger, error) {
	typeName, err := utils.ReadTypeFromJSON(data)
	if err != nil {
		return nil, err
	}

	f := registeredTypes[typeName]
	if f == nil {
		return nil, errors.Errorf("unknown type: '%s'", typeName)
	}
	return f(sessionAssets, data, missing)
}

func (t *baseTrigger) unmarshal(sessionAssets flows.SessionAssets, e *baseTriggerEnvelope, missing assets.MissingCallback) error {
	var err error

	t.type_ = e.Type
	t.flow = e.Flow
	t.connection = e.Connection
	t.batch = e.Batch
	t.history = e.History
	t.triggeredOn = e.TriggeredOn

	if e.Environment != nil {
		if t.environment, err = envs.ReadEnvironment(e.Environment); err != nil {
			return errors.Wrap(err, "unable to read environment")
		}
	}
	if e.Contact != nil {
		if t.contact, err = flows.ReadContact(sessionAssets, e.Contact, missing); err != nil {
			return errors.Wrap(err, "unable to read contact")
		}
	}
	if e.Params != nil {
		if t.params, err = types.ReadXObject(e.Params); err != nil {
			return errors.Wrap(err, "unable to read params")
		}
	}

	return nil
}

func (t *baseTrigger) marshal(e *baseTriggerEnvelope) error {
	var err error
	e.Type = t.type_
	e.Flow = t.flow
	e.Connection = t.connection
	e.Batch = t.batch
	e.History = t.history
	e.TriggeredOn = t.triggeredOn

	if t.environment != nil {
		e.Environment, err = jsonx.Marshal(t.environment)
		if err != nil {
			return err
		}
	}
	if t.contact != nil {
		e.Contact, err = jsonx.Marshal(t.contact)
		if err != nil {
			return err
		}
	}
	if t.params != nil {
		e.Params, err = jsonx.Marshal(t.params)
		if err != nil {
			return err
		}
	}
	return nil
}
