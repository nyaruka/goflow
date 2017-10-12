package flows

import (
	"encoding/json"

	"github.com/nyaruka/goflow/utils"
)

type runSummary struct {
	uuid    RunUUID
	flow    Flow
	contact *Contact
	status  RunStatus
	results *Results
}

func (r *runSummary) UUID() RunUUID     { return r.uuid }
func (r *runSummary) Flow() Flow        { return r.flow }
func (r *runSummary) Contact() *Contact { return r.contact }
func (r *runSummary) Status() RunStatus { return r.status }
func (r *runSummary) Results() *Results { return r.results }

func NewRunSummaryFromRun(run FlowRun) RunSummary {
	return &runSummary{
		uuid:    run.UUID(),
		flow:    run.Flow(),
		contact: run.Contact().Clone(),
		status:  run.Status(),
		results: run.Results().Clone(),
	}
}

var _ RunSummary = (*runSummary)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type runSummaryEnvelope struct {
	UUID     RunUUID         `json:"uuid" validate:"uuid4"`
	FlowUUID FlowUUID        `json:"flow_uuid" validate:"uuid4"`
	Contact  json.RawMessage `json:"contact" validate:"required"`
	Status   RunStatus       `json:"status" validate:"required"`
	Results  *Results        `json:"results"`
}

func ReadRunSummary(session Session, data json.RawMessage) (RunSummary, error) {
	var err error
	e := runSummaryEnvelope{}
	if err = utils.UnmarshalAndValidate(data, &e, "runsummary"); err != nil {
		return nil, err
	}

	run := &runSummary{
		uuid:    e.UUID,
		status:  e.Status,
		results: e.Results,
	}

	// lookup the flow
	if run.flow, err = session.Assets().GetFlow(e.FlowUUID); err != nil {
		return nil, err
	}

	// read the contact
	if run.contact, err = ReadContact(session, e.Contact); err != nil {
		return nil, err
	}

	return run, nil
}

func (r *runSummary) MarshalJSON() ([]byte, error) {
	envelope := runSummaryEnvelope{}
	var err error

	envelope.UUID = r.uuid
	envelope.FlowUUID = r.flow.UUID()
	envelope.Status = r.status
	envelope.Results = r.results

	if envelope.Contact, err = r.contact.MarshalJSON(); err != nil {
		return nil, err
	}

	return json.Marshal(envelope)
}
