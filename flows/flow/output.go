package flow

import (
	"encoding/json"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

type runOutput struct {
	runs   []flows.FlowRun
	events []flows.Event
}

func newRunOutput() *runOutput {
	output := runOutput{}
	return &output
}

func (o *runOutput) AddRun(run flows.FlowRun) { o.runs = append(o.runs, run) }
func (o *runOutput) Runs() []flows.FlowRun    { return o.runs }

func (o *runOutput) ActiveRun() flows.FlowRun {
	var active flows.FlowRun
	mostRecent := utils.ZeroTime

	for _, run := range o.runs {
		// We are complete, therefore can't be active
		if run.IsComplete() {
			continue
		}

		// We have a child, and it isn't complete, we can't be active
		if run.Child() != nil && run.Child().Status() == flows.RunActive {
			continue
		}

		// this is more recent than our most recent flow
		if run.ModifiedOn().After(mostRecent) {
			active = run
			mostRecent = run.ModifiedOn()
		}
	}
	return active
}

func (o *runOutput) AddEvent(event flows.Event) { o.events = append(o.events, event) }
func (o *runOutput) Events() []flows.Event      { return o.events }

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

// ReadRunOutput decodes a run output from the passed in JSON
func ReadRunOutput(data json.RawMessage) (flows.RunOutput, error) {
	runOutput := &runOutput{}
	err := json.Unmarshal(data, runOutput)
	if err == nil {
		// err = run.Validate()
	}
	return runOutput, err
}

type outputEnvelope struct {
	Runs   []*run                 `json:"runs"`
	Events []*utils.TypedEnvelope `json:"events"`
}

func (o *runOutput) UnmarshalJSON(data []byte) error {
	var oe outputEnvelope
	var err error

	err = json.Unmarshal(data, &oe)
	if err != nil {
		return err
	}

	o.runs = make([]flows.FlowRun, len(oe.Runs))
	for i := range o.runs {
		o.runs[i] = oe.Runs[i]
	}

	o.events = make([]flows.Event, len(oe.Events))
	for i := range o.events {
		o.events[i], err = events.EventFromEnvelope(oe.Events[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func (o *runOutput) MarshalJSON() ([]byte, error) {
	var oe outputEnvelope

	oe.Events = make([]*utils.TypedEnvelope, len(o.events))
	for i, event := range o.events {
		eventData, err := json.Marshal(event)
		if err != nil {
			return nil, err
		}
		oe.Events[i] = &utils.TypedEnvelope{Type: event.Type(), Data: eventData}
	}

	oe.Runs = make([]*run, len(o.runs))
	for i := range o.runs {
		oe.Runs[i] = o.runs[i].(*run)
	}

	return json.Marshal(oe)
}
