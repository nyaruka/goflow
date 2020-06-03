package triggers

import (
	"encoding/json"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
	"github.com/nyaruka/goflow/utils/jsonx"

	"github.com/pkg/errors"
)

func init() {
	registerType(TypeFlowAction, readFlowActionTrigger)
}

// TypeFlowAction is a constant for sessions triggered by flow actions in other sessions
const TypeFlowAction string = "flow_action"

// FlowActionTrigger is used when another session triggered this run using a trigger_flow action.
//
//   {
//     "type": "flow_action",
//     "flow": {"uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d", "name": "Collect Age"},
//     "triggered_on": "2000-01-01T00:00:00.000000000-00:00",
//     "run_summary": {
//       "uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d",
//       "flow": {"uuid": "50c3706e-fedb-42c0-8eab-dda3335714b7", "name": "Registration"},
//       "contact": {
//         "uuid": "c59b0033-e748-4240-9d4c-e85eb6800151",
//         "name": "Bob",
//         "fields": {"gender": {"text": "Male"}},
//         "created_on": "2018-01-01T12:00:00.000000000-00:00"
//       },
//       "status": "active",
//       "results": {
//         "age": {
//           "result_name": "Age",
//           "value": "33",
//           "node": "cd2be8c4-59bc-453c-8777-dec9a80043b8",
//           "created_on": "2018-01-01T12:00:00.000000000-00:00"
//         }
//       }
//     }
//   }
//
// @trigger flow_action
type FlowActionTrigger struct {
	baseTrigger

	runSummary json.RawMessage
}

// NewFlowAction creates a new flow action trigger with the passed in values
func NewFlowAction(env envs.Environment, flow *assets.FlowReference, contact *flows.Contact, runSummary json.RawMessage, batch bool) (*FlowActionTrigger, error) {
	return newFlowAction(env, flow, contact, nil, runSummary, batch)
}

// NewFlowActionVoice creates a new flow action trigger with the passed in values
func NewFlowActionVoice(env envs.Environment, flow *assets.FlowReference, contact *flows.Contact, connection *flows.Connection, runSummary json.RawMessage, batch bool) (*FlowActionTrigger, error) {
	return newFlowAction(env, flow, contact, connection, runSummary, batch)
}

func newFlowAction(env envs.Environment, flow *assets.FlowReference, contact *flows.Contact, connection *flows.Connection, runSummary json.RawMessage, batch bool) (*FlowActionTrigger, error) {
	if !json.Valid(runSummary) {
		return nil, errors.Errorf("invalid run summary JSON: %s", string(runSummary))
	}

	return &FlowActionTrigger{
		baseTrigger: newBaseTrigger(TypeFlowAction, env, flow, contact, connection, batch, nil),
		runSummary:  runSummary,
	}, nil
}

// RunSummary returns the summary of the run that triggered this session
func (t *FlowActionTrigger) RunSummary() json.RawMessage { return t.runSummary }

var _ flows.TriggerWithRun = (*FlowActionTrigger)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type flowActionTriggerEnvelope struct {
	baseTriggerEnvelope
	RunSummary json.RawMessage `json:"run_summary" validate:"required"`
}

func readFlowActionTrigger(sessionAssets flows.SessionAssets, data json.RawMessage, missing assets.MissingCallback) (flows.Trigger, error) {
	e := &flowActionTriggerEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	t := &FlowActionTrigger{
		runSummary: e.RunSummary,
	}

	if err := t.unmarshal(sessionAssets, &e.baseTriggerEnvelope, missing); err != nil {
		return nil, err
	}

	return t, nil
}

// MarshalJSON marshals this trigger into JSON
func (t *FlowActionTrigger) MarshalJSON() ([]byte, error) {
	e := &flowActionTriggerEnvelope{
		RunSummary: t.runSummary,
	}

	if err := t.marshal(&e.baseTriggerEnvelope); err != nil {
		return nil, err
	}

	return jsonx.Marshal(e)
}
