package flow

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/flows/inputs"
	"github.com/nyaruka/goflow/flows/waits"
	"github.com/nyaruka/goflow/utils"
	uuid "github.com/satori/go.uuid"
)

type run struct {
	uuid flows.RunUUID

	contact     flows.Contact
	contactUUID flows.ContactUUID

	flow     flows.Flow
	flowUUID flows.FlowUUID

	channel     flows.Channel
	channelUUID flows.ChannelUUID

	results results
	context flows.Context
	status  flows.RunStatus
	wait    flows.Wait
	webhook utils.RequestResponse
	input   flows.Input
	event   flows.Event

	parent flows.FlowRun
	child  flows.FlowRun

	output flows.RunOutput

	path             []flows.Step
	flowTranslations flows.FlowTranslations
	translations     flows.Translations
	environment      flows.FlowEnvironment
	language         flows.Language

	createdOn  time.Time
	modifiedOn time.Time
	expiresOn  *time.Time
	timesOutOn *time.Time
	exitedOn   *time.Time
}

func (r *run) UUID() flows.RunUUID { return r.uuid }

func (r *run) ContactUUID() flows.ContactUUID { return r.contactUUID }
func (r *run) Contact() flows.Contact         { return r.contact }
func (r *run) SetContact(contact flows.Contact) error {
	if contact.UUID() != r.contactUUID {
		return fmt.Errorf("Cannot change contact on an existing run")
	}
	r.contact = contact
	return nil
}

func (r *run) FlowUUID() flows.FlowUUID { return r.flowUUID }
func (r *run) Flow() flows.Flow         { return r.flow }
func (r *run) SetFlow(flow flows.Flow) error {
	if flow.UUID() != r.flowUUID {
		return fmt.Errorf("Cannot change flow on an existing run")
	}
	r.flow = flow
	return nil
}

func (r *run) ChannelUUID() flows.ChannelUUID { return r.channelUUID }
func (r *run) Channel() flows.Channel         { return r.channel }
func (r *run) SetChannel(channel flows.Channel) {
	r.channelUUID = channel.UUID()
	r.channel = channel
}

func (r *run) Context() flows.Context             { return r.context }
func (r *run) Environment() flows.FlowEnvironment { return r.environment }
func (r *run) Results() flows.Results             { return r.results }

func (r *run) Output() flows.RunOutput { return r.output }
func (r *run) ResetOutput() {
	r.output = newRunOutput()
	r.output.AddRun(r)
}

func (r *run) IsComplete() bool {
	return r.status != flows.RunActive
}
func (r *run) setStatus(status flows.RunStatus) {
	now := time.Now()
	r.status = status
	r.exitedOn = &now
	r.setModifiedOn(now)
}
func (r *run) Exit(status flows.RunStatus) { r.setStatus(status) }
func (r *run) Status() flows.RunStatus     { return r.status }

func (r *run) Parent() flows.FlowRun { return r.parent }
func (r *run) Child() flows.FlowRun  { return r.child }

func (r *run) Wait() flows.Wait        { return r.wait }
func (r *run) SetWait(wait flows.Wait) { r.wait = wait }

func (r *run) Input() flows.Input         { return r.input }
func (r *run) SetInput(input flows.Input) { r.input = input }

func (r *run) SetEvent(event flows.Event) { r.event = event }
func (r *run) Event() flows.Event         { return r.event }

func (r *run) AddEvent(s flows.Step, e flows.Event) {
	now := time.Now()

	e.SetCreatedOn(now)
	e.SetRun(r.UUID())

	fs := s.(*step)
	fs.addEvent(e)

	r.Output().AddEvent(e)
	r.setModifiedOn(now)
}

func (r *run) AddError(step flows.Step, err error) {
	r.AddEvent(step, &events.ErrorEvent{Text: err.Error()})
}

func (r *run) Path() []flows.Step { return r.path }
func (r *run) CreateStep(node flows.Node) flows.Step {
	now := time.Now()
	step := &step{node: node.UUID(), arrivedOn: now}
	r.path = append(r.path, step)
	r.setModifiedOn(now)
	return step
}
func (r *run) ClearPath() {
	r.path = nil
}

func (r *run) Webhook() utils.RequestResponse { return r.webhook }
func (r *run) SetWebhook(rr utils.RequestResponse) {
	r.webhook = rr
	r.setModifiedOn(time.Now())
}

func (r *run) CreatedOn() time.Time        { return r.createdOn }
func (r *run) ModifiedOn() time.Time       { return r.modifiedOn }
func (r *run) setModifiedOn(now time.Time) { r.modifiedOn = now }
func (r *run) ExpiresOn() *time.Time       { return r.expiresOn }
func (r *run) TimesOutOn() *time.Time      { return r.timesOutOn }
func (r *run) ExitedOn() *time.Time        { return r.exitedOn }

func (r *run) updateTranslations() {
	if r.flowTranslations != nil {
		r.translations = r.flowTranslations.GetTranslations(r.language)
	} else {
		r.translations = nil
	}
}
func (r *run) SetFlowTranslations(ft flows.FlowTranslations) {
	r.flowTranslations = ft
	r.updateTranslations()
}
func (r *run) SetLanguage(lang flows.Language) {
	r.language = lang
	r.updateTranslations()
}
func (r *run) GetText(uuid flows.UUID, key string, backdown string) string {
	if r.translations == nil {
		return backdown
	}
	return r.translations.GetText(uuid, key, backdown)
}

// NewRun initializes a new context and flow run for the passed in flow and contact
func newRun(env flows.FlowEnvironment, flow flows.Flow, contact flows.Contact, parent flows.FlowRun) flows.FlowRun {
	now := time.Now()

	r := &run{
		uuid:        flows.RunUUID(uuid.NewV4().String()),
		flowUUID:    flow.UUID(),
		contactUUID: contact.UUID(),
		results:     make(results, 0),
		environment: env,
		contact:     contact,
		flow:        flow,
		status:      flows.RunActive,
		createdOn:   now,
		modifiedOn:  now,
		parent:      parent,
	}

	// create our new context
	r.context = NewContextForContact(contact, r)

	// set our output
	if parent != nil {
		r.output = parent.Output()
	} else {
		r.output = newRunOutput()

	}
	r.output.AddRun(r)

	// set ourselves as the child to our parent
	if parent != nil {
		parent.(*run).child = r
	}

	return r
}

func (r *run) Resolve(key string) interface{} {
	switch key {

	case "channel":
		return r.Channel()

	case "contact":
		return r.Contact()

	case "flow":
		return r.Flow()

	case "status":
		return r.Status()

	case "results":
		return r.Results()

	case "created_on":
		return r.CreatedOn()

	case "exited_on":
		return r.ExitedOn()

	}

	return fmt.Errorf("No field '%s' on run", key)
}

func (r *run) Default() interface{} {
	return r
}

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
		if run.Child() != nil && !run.Child().IsComplete() {
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

// ReadRun decodes a run from the passed in JSON
func ReadRun(data json.RawMessage) (flows.FlowRun, error) {
	run := &run{}
	err := json.Unmarshal(data, run)
	if err == nil {
		// err = run.Validate()
	}
	return run, err
}

// ReadRunOutput decodes a run output from the passed in JSON
func ReadRunOutput(data json.RawMessage) (flows.RunOutput, error) {
	runOutput := &runOutput{}
	err := json.Unmarshal(data, runOutput)
	if err == nil {
		// err = run.Validate()
	}
	return runOutput, err
}

type runEnvelope struct {
	UUID    flows.RunUUID     `json:"uuid"`
	Flow    flows.FlowUUID    `json:"flow"`
	Channel flows.ChannelUUID `json:"channel"`
	Contact flows.ContactUUID `json:"contact"`
	Path    []*step           `json:"path"`

	Status flows.RunStatus `json:"status"`

	Input *utils.TypedEnvelope `json:"input,omitempty"`
	Wait  *utils.TypedEnvelope `json:"wait,omitempty"`
	Event *utils.TypedEnvelope `json:"event,omitempty"`

	Results results               `json:"results"`
	Webhook utils.RequestResponse `json:"webhook,omitempty"`

	CreatedOn  time.Time  `json:"created_on"`
	ModifiedOn time.Time  `json:"modified_on"`
	ExpiresOn  *time.Time `json:"expires_on"`
	TimesOutOn *time.Time `json:"timesout_on"`
	ExitedOn   *time.Time `json:"exited_on"`
}

func (r *run) UnmarshalJSON(data []byte) error {
	var envelope runEnvelope
	var err error

	err = json.Unmarshal(data, &envelope)
	if err != nil {
		return err
	}

	r.uuid = envelope.UUID
	r.contactUUID = envelope.Contact
	r.flowUUID = envelope.Flow
	r.channelUUID = envelope.Channel
	r.status = envelope.Status
	r.createdOn = envelope.CreatedOn
	r.modifiedOn = envelope.ModifiedOn
	r.expiresOn = envelope.ExpiresOn
	r.timesOutOn = envelope.TimesOutOn
	r.exitedOn = envelope.ExitedOn

	if envelope.Input != nil {
		r.input, err = inputs.InputFromEnvelope(envelope.Input)
		if err != nil {
			return err
		}
	}

	if envelope.Wait != nil {
		r.wait, err = waits.WaitFromEnvelope(envelope.Wait)
		if err != nil {
			return err
		}
	}

	if envelope.Event != nil {
		r.event, err = events.EventFromEnvelope(envelope.Event)
		if err != nil {
			return err
		}
	}

	if envelope.Results != nil {
		r.results = envelope.Results
	}

	if envelope.Webhook != nil {
		r.webhook = envelope.Webhook
	}

	// read in our path
	r.path = make([]flows.Step, len(envelope.Path))
	for i, step := range envelope.Path {
		r.path[i] = step
	}

	return err
}

func (r *run) MarshalJSON() ([]byte, error) {
	var re runEnvelope
	var err error

	re.UUID = r.uuid
	re.Flow = r.FlowUUID()
	re.Contact = r.ContactUUID()
	re.Channel = r.ChannelUUID()

	re.Status = r.status
	re.CreatedOn = r.createdOn
	re.ModifiedOn = r.modifiedOn
	re.ExpiresOn = r.expiresOn
	re.TimesOutOn = r.timesOutOn
	re.ExitedOn = r.exitedOn
	re.Results = r.results

	re.Input, err = utils.EnvelopeFromTyped(r.input)
	if err != nil {
		return nil, err
	}

	re.Wait, err = utils.EnvelopeFromTyped(r.wait)
	if err != nil {
		return nil, err
	}

	re.Event, err = utils.EnvelopeFromTyped(r.event)
	if err != nil {
		return nil, err
	}

	re.Path = make([]*step, len(r.path))
	for i, s := range r.path {
		re.Path[i] = s.(*step)
	}

	return json.Marshal(re)
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
