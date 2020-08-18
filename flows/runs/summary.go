package runs

import (
	"encoding/json"
	"fmt"

	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/utils"
)

// concrete run summary which might be stored on a trigger or event
type runSummary struct {
	uuid    flows.RunUUID
	flow    flows.Flow
	flowRef *assets.FlowReference
	contact *flows.Contact
	status  flows.RunStatus
	results flows.Results
}

// creates a new run summary from the given run
func newRunSummaryFromRun(run flows.FlowRun) flows.RunSummary {
	return &runSummary{
		uuid:    run.UUID(),
		flow:    run.Flow(),
		flowRef: run.Flow().Reference(),
		contact: run.Contact().Clone(),
		status:  run.Status(),
		results: run.Results().Clone(),
	}
}

func (r *runSummary) UUID() flows.RunUUID     { return r.uuid }
func (r *runSummary) Flow() flows.Flow        { return r.flow }
func (r *runSummary) Contact() *flows.Contact { return r.contact }
func (r *runSummary) Status() flows.RunStatus { return r.status }
func (r *runSummary) Results() flows.Results  { return r.results }

var _ flows.RunSummary = (*runSummary)(nil)

// wrapper for a run summary (concrete like runSummary or view of child run via interface)
type relatedRunContext struct {
	run flows.RunSummary
}

func newRelatedRunContext(run flows.RunSummary) *relatedRunContext {
	if utils.IsNil(run) {
		return nil
	}
	return &relatedRunContext{run: run}
}

// Context returns the properties available in expressions for @parent and @child
//
//   __default__:text -> the contact name and flow UUID
//   uuid:text -> the UUID of the run
//   contact:contact -> the contact of the run
//   flow:flow -> the flow of the run
//   fields:fields -> the custom field values of the run's contact
//   urns:urns -> the URN values of the run's contact
//   results:any -> the results saved by the run
//   status:text -> the current status of the run
//
// @context related_run
func (c *relatedRunContext) Context(env envs.Environment) map[string]types.XValue {
	var urns, fields types.XValue
	if c.run.Contact() != nil {
		urns = flows.ContextFunc(env, c.run.Contact().URNs().MapContext)
		fields = flows.Context(env, c.run.Contact().Fields())
	}

	return map[string]types.XValue{
		"__default__": types.NewXText(FormatRunSummary(env, c.run)),
		"uuid":        types.NewXText(string(c.run.UUID())),
		"run":         flows.ContextFunc(env, c.RunContext), // deprecated to be removed in 13.1
		"contact":     flows.Context(env, c.run.Contact()),
		"flow":        flows.Context(env, c.run.Flow()),
		"urns":        urns,
		"fields":      fields,
		"results":     flows.Context(env, c.run.Results()),
		"status":      types.NewXText(string(c.run.Status())),
	}
}

// Context returns the properties available in expressions. Run summaries expose a
// subset of the properties exposed by a real run.
func (c *relatedRunContext) RunContext(env envs.Environment) map[string]types.XValue {
	return map[string]types.XValue{
		"__default__": types.NewXText(FormatRunSummary(env, c.run)),
		"uuid":        types.NewXText(string(c.run.UUID())),
		"contact":     flows.Context(env, c.run.Contact()),
		"flow":        flows.Context(env, c.run.Flow()),
		"status":      types.NewXText(string(c.run.Status())),
		"results":     flows.Context(env, c.run.Results()),
	}
}

// FormatRunSummary formats an instance of the RunSummary interface
func FormatRunSummary(env envs.Environment, run flows.RunSummary) string {
	var flow, contact string

	if run.Flow() != nil {
		flow = run.Flow().Name()
	} else {
		flow = "<missing>"
	}

	if run.Contact() != nil {
		contact = run.Contact().Format(env)
	} else {
		contact = "<nocontact>"
	}

	return fmt.Sprintf("%s@%s", contact, flow)
}

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type runSummaryEnvelope struct {
	UUID    flows.RunUUID         `json:"uuid" validate:"uuid4"`
	Flow    *assets.FlowReference `json:"flow" validate:"required,dive"`
	Contact json.RawMessage       `json:"contact"`
	Status  flows.RunStatus       `json:"status" validate:"required"`
	Results flows.Results         `json:"results"`
}

// ReadRunSummary reads a run summary from the given JSON
func ReadRunSummary(sessionAssets flows.SessionAssets, data json.RawMessage, missing assets.MissingCallback) (flows.RunSummary, error) {
	var err error
	e := runSummaryEnvelope{}
	if err = utils.UnmarshalAndValidate(data, &e); err != nil {
		return nil, err
	}

	run := &runSummary{
		uuid:    e.UUID,
		flowRef: e.Flow,
		status:  e.Status,
		results: e.Results,
	}

	// lookup the actual flow
	if run.flow, err = sessionAssets.Flows().Get(e.Flow.UUID); err != nil {
		missing(e.Flow, err)
	}

	// read the contact
	if e.Contact != nil {
		if run.contact, err = flows.ReadContact(sessionAssets, e.Contact, missing); err != nil {
			return nil, err
		}
	}

	return run, nil
}

// MarshalJSON marshals this run summary into JSON
func (r *runSummary) MarshalJSON() ([]byte, error) {
	envelope := runSummaryEnvelope{}
	var err error

	envelope.UUID = r.uuid
	envelope.Flow = r.flowRef
	envelope.Status = r.status
	envelope.Results = r.results

	if r.contact != nil {
		if envelope.Contact, err = r.contact.MarshalJSON(); err != nil {
			return nil, err
		}
	}

	return jsonx.Marshal(envelope)
}
