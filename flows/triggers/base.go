package triggers

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"

	"github.com/pkg/errors"
)

type readFunc func(sessionAssets flows.SessionAssets, data json.RawMessage) (flows.Trigger, error)

var registeredTypes = map[string]readFunc{}

// RegisterType registers a new type of trigger
func RegisterType(name string, f readFunc) {
	registeredTypes[name] = f
}

type baseTrigger struct {
	type_       string
	environment utils.Environment
	flow        *assets.FlowReference
	contact     *flows.Contact
	params      types.XValue
	triggeredOn time.Time
}

func newBaseTrigger(typeName string, env utils.Environment, flow *assets.FlowReference, contact *flows.Contact, params types.XValue, triggeredOn time.Time) baseTrigger {
	return baseTrigger{type_: typeName, environment: env, flow: flow, contact: contact, params: params, triggeredOn: triggeredOn}
}

// Type returns the type of this trigger
func (t *baseTrigger) Type() string { return t.type_ }

func (t *baseTrigger) Environment() utils.Environment { return t.environment }
func (t *baseTrigger) Flow() *assets.FlowReference    { return t.flow }
func (t *baseTrigger) Contact() *flows.Contact        { return t.contact }
func (t *baseTrigger) Params() types.XValue           { return t.params }
func (t *baseTrigger) TriggeredOn() time.Time         { return t.triggeredOn }

// Initialize initializes the session
func (t *baseTrigger) Initialize(session flows.Session) error {
	// try to load the flow
	flow, err := session.Assets().Flows().Get(t.Flow().UUID)
	if err != nil {
		return errors.Wrapf(err, "unable to load flow[uuid=%s]", t.Flow().UUID)
	}

	// check flow is valid and has everything it needs to run
	if err := flow.Validate(session.Assets(), flows.NewValidationContext()); err != nil {
		return errors.Wrapf(err, "validation failed for flow[uuid=%s]", flow.UUID())
	}

	session.PushFlow(flow, nil, false)

	if t.environment != nil {
		session.SetEnvironment(t.environment)
	}
	if t.contact != nil {
		session.SetContact(t.contact.Clone())

		EnsureDynamicGroups(session)
	}
	return nil
}

// InitializeRun performs additional initialization when we create our first run
func (t *baseTrigger) InitializeRun(run flows.FlowRun, logEvent func(flows.Event)) error {
	return nil
}

// Resolve resolves the given key when this trigger is referenced in an expression
func (t *baseTrigger) Resolve(env utils.Environment, key string) types.XValue {
	switch strings.ToLower(key) {
	case "type":
		return types.NewXText(t.type_)
	case "params":
		return t.params
	}

	return types.NewXResolveError(t, key)
}

// ToXJSON is called when this type is passed to @(json(...))
func (t *baseTrigger) ToXJSON(env utils.Environment) types.XText {
	return types.ResolveKeys(env, t, "type", "params").ToXJSON(env)
}

// Describe returns a representation of this type for error messages
func (t *baseTrigger) Describe() string { return "trigger" }

// Reduce is called when this object needs to be reduced to a primitive
func (t *baseTrigger) Reduce(env utils.Environment) types.XPrimitive {
	return types.NewXText(string(t.flow.UUID))
}

// EnsureDynamicGroups ensures that our session contact is in the correct dynamic groups as
// as far as the engine is concerned
func EnsureDynamicGroups(session flows.Session) {
	allGroups := session.Assets().Groups()
	added, removed, errors := session.Contact().ReevaluateDynamicGroups(session.Environment(), allGroups)

	// add error event for each group we couldn't re-evaluate
	for _, err := range errors {
		session.LogEvent(events.NewErrorEvent(err))
	}

	// add groups changed event for the groups we were added/removed to/from
	if len(added) > 0 || len(removed) > 0 {
		session.LogEvent(events.NewContactGroupsChangedEvent(added, removed))
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
	Params      json.RawMessage       `json:"params,omitempty"`
	TriggeredOn time.Time             `json:"triggered_on" validate:"required"`
}

// ReadTrigger reads a trigger from the given JSON
func ReadTrigger(sessionAssets flows.SessionAssets, data json.RawMessage) (flows.Trigger, error) {
	typeName, err := utils.ReadTypeFromJSON(data)
	if err != nil {
		return nil, err
	}

	f := registeredTypes[typeName]
	if f == nil {
		return nil, errors.Errorf("unknown type: '%s'", typeName)
	}
	return f(sessionAssets, data)
}

func (t *baseTrigger) unmarshal(sessionAssets flows.SessionAssets, e *baseTriggerEnvelope) error {
	var err error

	t.type_ = e.Type
	t.flow = e.Flow
	t.triggeredOn = e.TriggeredOn

	if e.Environment != nil {
		if t.environment, err = utils.ReadEnvironment(e.Environment); err != nil {
			return errors.Wrap(err, "unable to read environment")
		}
	}
	if e.Contact != nil {
		if t.contact, err = flows.ReadContact(sessionAssets, e.Contact, true); err != nil {
			return errors.Wrap(err, "unable to read contact")
		}
	}
	if e.Params != nil {
		t.params = types.JSONToXValue(e.Params)
	}

	return nil
}

func (t *baseTrigger) marshal(e *baseTriggerEnvelope) error {
	var err error
	e.Type = t.type_
	e.Flow = t.flow
	e.TriggeredOn = t.triggeredOn

	if t.environment != nil {
		e.Environment, err = json.Marshal(t.environment)
		if err != nil {
			return err
		}
	}
	if t.contact != nil {
		e.Contact, err = json.Marshal(t.contact)
		if err != nil {
			return err
		}
	}
	if t.params != nil {
		e.Params, err = json.Marshal(t.params)
		if err != nil {
			return err
		}
	}
	return nil
}
