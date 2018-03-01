package triggers

import (
	"encoding/json"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

// TypeFlowAction is a constant for sessions triggered by flow actions in other sessions
const TypeFlowAction string = "flow_action"

// FlowActionTrigger is used when another session triggered this run using a trigger_flow action.
//
// ```
//   {
//     "type": "flow_action",
//     "flow": {"uuid": "ea7d8b6b-a4b2-42c1-b9cf-c0370a95a721", "name": "Registration"},
//     "triggered_on": "2000-01-01T00:00:00.000000000-00:00",
//     "run": {
//       "uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d",
//       "flow_uuid": "93c554a1-b90d-4892-b029-a2a87dec9b87",
//       "contact": {
//         "uuid": "c59b0033-e748-4240-9d4c-e85eb6800151",
//         "name": "Bob",
//         "fields": {"state": {"value": "Azuay", "created_on": "2000-01-01T00:00:00.000000000-00:00"}}
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
// ```
type FlowActionTrigger struct {
	baseTrigger
	run flows.RunSummary
}

// Type returns the type of this trigger
func (t *FlowActionTrigger) Type() string { return TypeFlowAction }

func (t *FlowActionTrigger) Run() flows.RunSummary { return t.run }

var _ flows.Trigger = (*FlowActionTrigger)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type flowActionTriggerEnvelope struct {
	baseTriggerEnvelope
	Run json.RawMessage `json:"run"`
}

func ReadFlowActionTrigger(session flows.Session, envelope *utils.TypedEnvelope) (flows.Trigger, error) {
	var err error
	trigger := &FlowActionTrigger{}
	e := flowActionTriggerEnvelope{}
	if err := utils.UnmarshalAndValidate(envelope.Data, &e, "trigger[type=flow_action]"); err != nil {
		return nil, err
	}

	if err := unmarshalBaseTrigger(session, &trigger.baseTrigger, &e.baseTriggerEnvelope); err != nil {
		return nil, err
	}

	if trigger.run, err = flows.ReadRunSummary(session, e.Run); err != nil {
		return nil, err
	}

	return trigger, nil
}

func (t *FlowActionTrigger) MarshalJSON() ([]byte, error) {
	var envelope flowActionTriggerEnvelope
	var err error

	if err := marshalBaseTrigger(&t.baseTrigger, &envelope.baseTriggerEnvelope); err != nil {
		return nil, err
	}

	if envelope.Run, err = json.Marshal(t.run); err != nil {
		return nil, err
	}

	return json.Marshal(envelope)
}
