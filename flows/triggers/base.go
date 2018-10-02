package triggers

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

type readFunc func(session flows.Session, data json.RawMessage) (flows.Trigger, error)

var registeredTypes = map[string]readFunc{}

// RegisterType registers a new type of trigger
func RegisterType(name string, f readFunc) {
	registeredTypes[name] = f
}

type baseTrigger struct {
	environment utils.Environment
	flow        *assets.FlowReference
	contact     *flows.Contact
	params      types.XValue
	triggeredOn time.Time
}

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
		return fmt.Errorf("unable to load flow[uuid=%s]: %s", t.Flow().UUID, err)
	}

	// check flow is valid and has everything it needs to run
	if err := flow.Validate(session.Assets()); err != nil {
		return fmt.Errorf("validation failed for flow[uuid=%s]: %s", flow.UUID(), err)
	}

	session.PushFlow(flow, nil)

	if t.environment != nil {
		session.SetEnvironment(t.environment)
	}
	if t.contact != nil {
		session.SetContact(t.contact.Clone())
	}
	return nil
}

// InitializeRun performs additional initialization when we create our first run
func (t *baseTrigger) InitializeRun(run flows.FlowRun) error {
	return nil
}

// Resolve resolves the given key when this trigger is referenced in an expression
func (t *baseTrigger) Resolve(env utils.Environment, key string) types.XValue {
	switch key {
	case "params":
		return t.params
	}

	return types.NewXResolveError(t, key)
}

// Describe returns a representation of this type for error messages
func (t *baseTrigger) Describe() string { return "trigger" }

// Reduce is called when this object needs to be reduced to a primitive
func (t *baseTrigger) Reduce(env utils.Environment) types.XPrimitive {
	return types.NewXText(string(t.flow.UUID))
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type baseTriggerEnvelope struct {
	Environment json.RawMessage       `json:"environment,omitempty"`
	Flow        *assets.FlowReference `json:"flow" validate:"required"`
	Contact     json.RawMessage       `json:"contact,omitempty"`
	Params      json.RawMessage       `json:"params,omitempty"`
	TriggeredOn time.Time             `json:"triggered_on" validate:"required"`
}

// ReadTrigger reads a trigger from the given typed envelope
func ReadTrigger(session flows.Session, envelope *utils.TypedEnvelope) (flows.Trigger, error) {
	f := registeredTypes[envelope.Type]
	if f == nil {
		return nil, fmt.Errorf("unknown type: %s", envelope.Type)
	}
	return f(session, envelope.Data)
}

func unmarshalBaseTrigger(session flows.Session, base *baseTrigger, envelope *baseTriggerEnvelope) error {
	var err error

	base.flow = envelope.Flow
	base.triggeredOn = envelope.TriggeredOn

	if envelope.Environment != nil {
		if base.environment, err = utils.ReadEnvironment(envelope.Environment); err != nil {
			return fmt.Errorf("unable to read environment: %s", err)
		}
	}
	if envelope.Contact != nil {
		if base.contact, err = flows.ReadContact(session.Assets(), envelope.Contact, true); err != nil {
			return fmt.Errorf("unable to read contact: %s", err)
		}
	}
	if envelope.Params != nil {
		base.params = types.JSONToXValue(envelope.Params)
	}

	return nil
}

func marshalBaseTrigger(t *baseTrigger, envelope *baseTriggerEnvelope) error {
	var err error
	envelope.Flow = t.flow
	envelope.TriggeredOn = t.triggeredOn

	if t.environment != nil {
		envelope.Environment, err = json.Marshal(t.environment)
		if err != nil {
			return err
		}
	}
	if t.contact != nil {
		envelope.Contact, err = json.Marshal(t.contact)
		if err != nil {
			return err
		}
	}
	if t.params != nil {
		envelope.Params, err = json.Marshal(t.params)
		if err != nil {
			return err
		}
	}
	return nil
}
