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
	run runInfo
}

// Type returns the type of this trigger
func (t *RunTrigger) Type() string { return TypeRun }

func (t *RunTrigger) Run() flows.FlowRunInfo { return &t.run }

type runInfo struct {
	uuid    flows.RunUUID
	flow    flows.Flow
	contact *flows.Contact
	results *flows.Results
}

func (r *runInfo) UUID() flows.RunUUID     { return r.uuid }
func (r *runInfo) Flow() flows.Flow        { return r.flow }
func (r *runInfo) Contact() *flows.Contact { return r.contact }
func (r *runInfo) Status() flows.RunStatus { return flows.RunStatusActive }
func (r *runInfo) Results() *flows.Results { return r.results }

var _ flows.Trigger = (*RunTrigger)(nil)
var _ flows.FlowRunInfo = (*runInfo)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type runInfoEnvelope struct {
	UUID     flows.RunUUID   `json:"uuid" validate:"uuid4"`
	FlowUUID flows.FlowUUID  `json:"flow_uuid" validate:"uuid4"`
	Contact  json.RawMessage `json:"contact" validate:"required"`
	Results  *flows.Results  `json:"results"`
}

type runTriggerEnvelope struct {
	baseTriggerEnvelope
	Run runInfoEnvelope `json:"run"`
}

func ReadRunTrigger(session flows.Session, envelope *utils.TypedEnvelope) (flows.Trigger, error) {
	var err error
	trigger := RunTrigger{}
	e := runTriggerEnvelope{}
	if err := utils.UnmarshalAndValidate(envelope.Data, &e, "trigger[type=run]"); err != nil {
		return nil, err
	}

	if err := readBaseTrigger(session, &trigger.baseTrigger, &e.baseTriggerEnvelope); err != nil {
		return nil, err
	}

	trigger.run = runInfo{uuid: e.Run.UUID, results: e.Run.Results}

	// lookup the run flow
	if e.Run.FlowUUID != "" {
		if trigger.run.flow, err = session.Assets().GetFlow(e.Run.FlowUUID); err != nil {
			return nil, err
		}
	}
	// read the run contact
	if trigger.run.contact, err = flows.ReadContact(session, e.Run.Contact); err != nil {
		return nil, err
	}

	return &trigger, nil
}

func (t *RunTrigger) MarshalJSON() ([]byte, error) {
	var envelope runTriggerEnvelope
	var err error

	envelope.TriggeredOn = t.triggeredOn
	envelope.FlowUUID = t.flow.UUID()
	envelope.Run.UUID = t.run.UUID()
	envelope.Run.FlowUUID = t.run.flow.UUID()
	envelope.Run.Results = t.run.results

	if envelope.Run.Contact, err = t.run.contact.MarshalJSON(); err != nil {
		return nil, err
	}

	return json.Marshal(envelope)
}
