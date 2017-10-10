package triggers

import (
	"encoding/json"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

// TypeRun is a constant for incoming messages
const TypeRun string = "run"

// RunTrigger is used when another session triggered this run using a trigger_flow action.
//
// ```
//   {
//     "type": "run",
//     "flow_uuid": "ea7d8b6b-a4b2-42c1-b9cf-c0370a95a721",
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
type RunTrigger struct {
	baseTrigger
	run flows.RunSummary
}

// Type returns the type of this trigger
func (t *RunTrigger) Type() string { return TypeRun }

func (t *RunTrigger) Run() flows.RunSummary { return t.run }

var _ flows.Trigger = (*RunTrigger)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type runTriggerEnvelope struct {
	baseTriggerEnvelope
	Run json.RawMessage `json:"run"`
}

func ReadRunTrigger(session flows.Session, envelope *utils.TypedEnvelope) (flows.Trigger, error) {
	var err error
	trigger := &RunTrigger{}
	e := runTriggerEnvelope{}
	if err := utils.UnmarshalAndValidate(envelope.Data, &e, "trigger[type=run]"); err != nil {
		return nil, err
	}

	if err := readBaseTrigger(session, &trigger.baseTrigger, &e.baseTriggerEnvelope); err != nil {
		return nil, err
	}

	if trigger.run, err = flows.ReadRunSummary(session, e.Run); err != nil {
		return nil, err
	}

	return trigger, nil
}

func (t *RunTrigger) MarshalJSON() ([]byte, error) {
	var envelope runTriggerEnvelope
	var err error

	envelope.TriggeredOn = t.triggeredOn
	envelope.FlowUUID = t.flow.UUID()

	if envelope.Run, err = json.Marshal(t.run); err != nil {
		return nil, err
	}

	return json.Marshal(envelope)
}
