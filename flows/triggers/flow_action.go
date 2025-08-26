package triggers

import (
	"encoding/json"
	"fmt"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

func init() {
	registerType(TypeFlowAction, readFlowAction)
}

// TypeFlowAction is a constant for sessions triggered by flow actions in other sessions
const TypeFlowAction string = "flow_action"

// FlowAction is used when another session triggered this run using a trigger_flow action.
//
//	{
//	  "type": "flow_action",
//	  "flow": {"uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d", "name": "Collect Age"},
//	  "history": {
//	    "parent_uuid": "a5b25fb0-75fd-4898-a34f-5ff14fc19078",
//	    "ancestors": 3,
//	    "ancestors_since_input": 1
//	  },
//	  "triggered_on": "2000-01-01T00:00:00.000000000-00:00",
//	  "run_summary": {
//	    "uuid": "b7cf0d83-f1c9-411c-96fd-c511a4cfa86d",
//	    "flow": {"uuid": "50c3706e-fedb-42c0-8eab-dda3335714b7", "name": "Registration"},
//	    "contact": {
//	      "uuid": "c59b0033-e748-4240-9d4c-e85eb6800151",
//	      "name": "Bob",
//	      "status": "active",
//	      "fields": {"gender": {"text": "Male"}},
//	      "created_on": "2018-01-01T12:00:00.000000000-00:00"
//	    },
//	    "status": "active",
//	    "results": {
//	      "age": {
//	        "result_name": "Age",
//	        "value": "33",
//	        "node": "cd2be8c4-59bc-453c-8777-dec9a80043b8",
//	        "created_on": "2018-01-01T12:00:00.000000000-00:00"
//	      }
//	    }
//	  }
//	}
//
// @trigger flow_action
type FlowAction struct {
	baseTrigger

	runSummary []byte
}

// RunSummary returns the summary of the run that triggered this session
func (t *FlowAction) RunSummary() []byte { return t.runSummary }

var _ flows.TriggerWithRun = (*FlowAction)(nil)

//------------------------------------------------------------------------------------------
// Builder
//------------------------------------------------------------------------------------------

// FlowActionBuilder is a builder for flow action type triggers
type FlowActionBuilder struct {
	t *FlowAction
}

func (b *Builder) FlowAction(history *flows.SessionHistory, runSummary []byte) *FlowActionBuilder {
	if !json.Valid(runSummary) {
		panic(fmt.Sprintf("invalid run summary JSON: %s", string(runSummary)))
	}

	return &FlowActionBuilder{
		t: &FlowAction{
			baseTrigger: newBaseTrigger(TypeFlowAction, nil, b.flow, false, history),
			runSummary:  runSummary,
		},
	}
}

// AsBatch sets batch mode on for the trigger
func (b *FlowActionBuilder) AsBatch() *FlowActionBuilder {
	b.t.batch = true
	return b
}

// Build builds the trigger
func (b *FlowActionBuilder) Build() *FlowAction {
	return b.t
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type flowActionEnvelope struct {
	baseEnvelope

	RunSummary json.RawMessage `json:"run_summary" validate:"required"`
}

func readFlowAction(sa flows.SessionAssets, data []byte, missing assets.MissingCallback) (flows.Trigger, error) {
	e := &flowActionEnvelope{}
	if err := utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, err
	}

	t := &FlowAction{
		runSummary: e.RunSummary,
	}

	if err := t.unmarshal(sa, &e.baseEnvelope, missing); err != nil {
		return nil, err
	}

	return t, nil
}

// MarshalJSON marshals this trigger into JSON
func (t *FlowAction) MarshalJSON() ([]byte, error) {
	e := &flowActionEnvelope{
		RunSummary: t.runSummary,
	}

	if err := t.marshal(&e.baseEnvelope); err != nil {
		return nil, err
	}

	return jsonx.Marshal(e)
}
