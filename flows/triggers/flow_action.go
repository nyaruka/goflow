package triggers

import (
	"encoding/json"
	"fmt"

	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/runs"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	RegisterType(TypeFlowAction, ReadFlowActionTrigger)
}

// TypeFlowAction is a constant for sessions triggered by flow actions in other sessions
const TypeFlowAction string = "flow_action"

// FlowActionTrigger is used when another session triggered this run using a trigger_flow action.
//
//   {
//     "type": "flow_action",
//     "flow": {"uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d", "name": "Collect Age"},
//     "triggered_on": "2000-01-01T00:00:00.000000000-00:00",
//     "run": {
//       "uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d",
//       "flow": {"uuid": "50c3706e-fedb-42c0-8eab-dda3335714b7", "name": "Registration"},
//       "contact": {
//         "uuid": "c59b0033-e748-4240-9d4c-e85eb6800151",
//         "name": "Bob",
//         "fields": {"gender": {"text": "Male"}}
//       },
//       "status": "active",
//       "results": {
//         "age": {
//           "result_name": "Age",
//           "value": "33",
//           "node": "cd2be8c4-59bc-453c-8777-dec9a80043b8",
//           "created_on": "2000-01-01T00:00:00.000000000-00:00"
//         }
//       }
//     }
//   }
//
// @trigger flow_action
type FlowActionTrigger struct {
	baseTrigger
	run flows.RunSummary
}

// Type returns the type of this trigger
func (t *FlowActionTrigger) Type() string { return TypeFlowAction }

// Run returns the summary of the run that triggered this session
func (t *FlowActionTrigger) Run() flows.RunSummary { return t.run }

// Resolve resolves the given key when this trigger is referenced in an expression
func (t *FlowActionTrigger) Resolve(env utils.Environment, key string) types.XValue {
	switch key {
	case "type":
		return types.NewXText(TypeFlowAction)
	}

	return t.baseTrigger.Resolve(env, key)
}

// ToXJSON is called when this type is passed to @(json(...))
func (t *FlowActionTrigger) ToXJSON(env utils.Environment) types.XText {
	return types.ResolveKeys(env, t, "type", "params").ToXJSON(env)
}

var _ flows.Trigger = (*FlowActionTrigger)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type flowActionTriggerEnvelope struct {
	baseTriggerEnvelope
	Run json.RawMessage `json:"run"`
}

// ReadFlowActionTrigger reads a flow action trigger
func ReadFlowActionTrigger(session flows.Session, data json.RawMessage) (flows.Trigger, error) {
	e := &flowActionTriggerEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	var err error
	t := &FlowActionTrigger{}
	if t.run, err = runs.ReadRunSummary(session, e.Run); err != nil {
		return nil, fmt.Errorf("unable to read run summary: %s", err)
	}

	if err := t.unmarshal(session, &e.baseTriggerEnvelope); err != nil {
		return nil, err
	}

	return t, nil
}

// MarshalJSON marshals this trigger into JSON
func (t *FlowActionTrigger) MarshalJSON() ([]byte, error) {
	e := &flowActionTriggerEnvelope{}

	var err error
	if e.Run, err = json.Marshal(t.run); err != nil {
		return nil, err
	}

	if err := t.marshal(&e.baseTriggerEnvelope); err != nil {
		return nil, err
	}

	return json.Marshal(e)
}
