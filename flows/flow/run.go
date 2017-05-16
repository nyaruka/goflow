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

	parent flows.FlowRunReference
	child  flows.FlowRunReference

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

func (r *run) UUID() flows.RunUUID            { return r.uuid }
func (r *run) ContactUUID() flows.ContactUUID { return r.contactUUID }
func (r *run) Contact() flows.Contact         { return r.contact }

// Hydrate prepares a deserialized run for executions
func (r *run) Hydrate(env flows.FlowEnvironment) error {
	// start with a fresh output if we don't have one
	if r.output == nil {
		r.ResetOutput()
	}

	// save off our environment
	r.environment = env

	// set our flow
	runFlow, err := env.GetFlow(r.FlowUUID())
	if err != nil {
		return err
	}
	r.flow = runFlow

	// make sure we have a contact
	runContact, err := env.GetContact(r.ContactUUID())
	if err != nil {
		return err
	}
	r.contact = runContact

	// build our context
	r.context = NewContextForContact(runContact, r)

	// populate our run references
	if r.parent != nil {
		parent, err := env.GetRun(r.parent.UUID())
		if err != nil {
			return err
		}
		r.parent = newReferenceFromRun(parent.(*run))
	}

	if r.child != nil {
		child, err := env.GetRun(r.child.UUID())
		if err != nil {
			return err
		}
		r.child = newReferenceFromRun(child.(*run))
	}

	return nil
}

func (r *run) FlowUUID() flows.FlowUUID       { return r.flowUUID }
func (r *run) Flow() flows.Flow               { return r.flow }
func (r *run) ChannelUUID() flows.ChannelUUID { return r.channelUUID }
func (r *run) Channel() flows.Channel         { return r.channel }
func (r *run) SetChannel(channel flows.Channel) {
	r.channelUUID = channel.UUID()
	r.channel = channel
}

func (r *run) Context() flows.Context             { return r.context }
func (r *run) Environment() flows.FlowEnvironment { return r.environment }
func (r *run) Results() flows.Results             { return r.results }

func (r *run) Output() flows.RunOutput          { return r.output }
func (r *run) SetOutput(output flows.RunOutput) { r.output = output }
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

func (r *run) Parent() flows.FlowRunReference { return r.parent }
func (r *run) Child() flows.FlowRunReference  { return r.child }

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
	}

	// create our new context
	r.context = NewContextForContact(contact, r)

	// set our output
	if parent != nil {
		parentRun := parent.(*run)
		r.parent = newReferenceFromRun(parentRun)
		parentRun.child = newReferenceFromRun(r)

		r.output = parent.Output()
	} else {
		r.output = newRunOutput()
	}
	r.output.AddRun(r)

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

// runReference provides a standalone and serializable version of a run reference. When a run is written
// this is what is written and what will be read (and needed) when resuming that run
type runReference struct {
	uuid flows.RunUUID
	run  *run
}

// Resolve provides a more limited set of results for parent and child references
func (r *runReference) Resolve(key string) interface{} {
	switch key {

	case "channel_uuid":
		return r.ChannelUUID()

	case "contact_uuid":
		return r.ContactUUID()

	case "flow_uuid":
		return r.FlowUUID()

	case "status":
		return r.Status()

	case "results":
		return r.Results()

	case "created_on":
		return r.CreatedOn()

	case "exited_on":
		return r.ExitedOn()

	}

	return fmt.Errorf("No field '%s' on run reference", key)
}

func (r *runReference) Default() interface{} {
	return r
}

func (r *runReference) UUID() flows.RunUUID            { return r.uuid }
func (r *runReference) FlowUUID() flows.FlowUUID       { return r.run.flowUUID }
func (r *runReference) ContactUUID() flows.ContactUUID { return r.run.contactUUID }
func (r *runReference) ChannelUUID() flows.ChannelUUID { return r.run.channelUUID }

func (r *runReference) Results() flows.Results  { return r.run.results }
func (r *runReference) Status() flows.RunStatus { return r.run.status }

func (r *runReference) CreatedOn() time.Time   { return r.run.createdOn }
func (r *runReference) ModifiedOn() time.Time  { return r.run.modifiedOn }
func (r *runReference) ExitedOn() *time.Time   { return r.run.exitedOn }
func (r *runReference) ExpiresOn() *time.Time  { return r.run.expiresOn }
func (r *runReference) TimesOutOn() *time.Time { return r.run.timesOutOn }

func newReferenceFromRun(r *run) *runReference {
	return &runReference{
		uuid: r.UUID(),
		run:  r,
	}
}

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

	Parent flows.RunUUID `json:"parent,omitempty"`
	Child  flows.RunUUID `json:"child,omitempty"`

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

	if envelope.Parent != "" {
		r.parent = &runReference{uuid: envelope.Parent}
	}
	if envelope.Child != "" {
		r.child = &runReference{uuid: envelope.Child}
	}

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

	if r.parent != nil {
		re.Parent = r.parent.UUID()
	}
	if r.child != nil {
		re.Child = r.child.UUID()
	}

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
