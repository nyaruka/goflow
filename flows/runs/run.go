package runs

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

type flowRun struct {
	uuid flows.RunUUID

	contact     *flows.Contact
	contactUUID flows.ContactUUID

	flow     flows.Flow
	flowUUID flows.FlowUUID

	channel     *flows.Channel
	channelUUID flows.ChannelUUID

	extra   json.RawMessage
	results *flows.Results
	context flows.Context
	status  flows.RunStatus
	wait    flows.Wait
	webhook utils.RequestResponse
	input   flows.Input
	event   flows.Event

	parent flows.FlowRunReference
	child  flows.FlowRunReference

	session flows.Session

	path             []flows.Step
	flowTranslations flows.FlowTranslations
	translations     flows.Translations
	environment      flows.FlowEnvironment
	language         utils.Language

	createdOn  time.Time
	modifiedOn time.Time
	expiresOn  *time.Time
	timesOutOn *time.Time
	exitedOn   *time.Time
}

func (r *flowRun) UUID() flows.RunUUID            { return r.uuid }
func (r *flowRun) ContactUUID() flows.ContactUUID { return r.contactUUID }
func (r *flowRun) Contact() *flows.Contact        { return r.contact }

// Hydrate prepares a deserialized run for executions
func (r *flowRun) Hydrate(env flows.FlowEnvironment) error {
	// start with a fresh output if we don't have one
	if r.session == nil {
		r.ResetSession()
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
		r.parent = newReferenceFromRun(parent.(*flowRun))
	}

	if r.child != nil {
		child, err := env.GetRun(r.child.UUID())
		if err != nil {
			return err
		}
		r.child = newReferenceFromRun(child.(*flowRun))
	}

	return nil
}

func (r *flowRun) FlowUUID() flows.FlowUUID       { return r.flowUUID }
func (r *flowRun) Flow() flows.Flow               { return r.flow }
func (r *flowRun) ChannelUUID() flows.ChannelUUID { return r.channelUUID }
func (r *flowRun) Channel() *flows.Channel        { return r.channel }
func (r *flowRun) SetChannel(channel *flows.Channel) {
	r.channelUUID = channel.UUID()
	r.channel = channel
}

func (r *flowRun) Context() flows.Context             { return r.context }
func (r *flowRun) Environment() flows.FlowEnvironment { return r.environment }
func (r *flowRun) Results() *flows.Results            { return r.results }

func (r *flowRun) Session() flows.Session { return r.session }
func (r *flowRun) SetSession(session flows.Session) {
	r.session = session
	r.session.AddRun(r)
}
func (r *flowRun) ResetSession() {
	r.SetSession(newSession())
}

func (r *flowRun) IsComplete() bool {
	return r.status != flows.StatusActive
}
func (r *flowRun) setStatus(status flows.RunStatus) {
	r.status = status
	r.setModifiedOn(time.Now())
}
func (r *flowRun) Exit(status flows.RunStatus) {
	r.setStatus(status)
	r.exitedOn = &r.modifiedOn
}
func (r *flowRun) Status() flows.RunStatus { return r.status }

func (r *flowRun) Parent() flows.FlowRunReference { return r.parent }
func (r *flowRun) Child() flows.FlowRunReference  { return r.child }

func (r *flowRun) Wait() flows.Wait        { return r.wait }
func (r *flowRun) SetWait(wait flows.Wait) { r.wait = wait }

func (r *flowRun) Input() flows.Input         { return r.input }
func (r *flowRun) SetInput(input flows.Input) { r.input = input }

func (r *flowRun) SetEvent(event flows.Event) { r.event = event }
func (r *flowRun) Event() flows.Event         { return r.event }

func (r *flowRun) AddEvent(s flows.Step, e flows.Event) {
	now := time.Now().In(time.UTC)

	e.SetCreatedOn(now)
	e.SetStep(s.UUID())

	fs := s.(*step)
	fs.addEvent(e)

	r.Session().AddEvent(e)
	r.setModifiedOn(now)
}

func (r *flowRun) AddError(step flows.Step, err error) {
	r.AddEvent(step, &events.ErrorEvent{Text: err.Error()})
}

func (r *flowRun) Path() []flows.Step { return r.path }
func (r *flowRun) CreateStep(node flows.Node) flows.Step {
	now := time.Now().In(time.UTC)
	step := &step{stepUUID: flows.StepUUID(uuid.NewV4().String()), nodeUUID: node.UUID(), arrivedOn: now}
	r.path = append(r.path, step)
	r.setModifiedOn(now)
	return step
}
func (r *flowRun) ClearPath() {
	r.path = nil
}

func (r *flowRun) Webhook() utils.RequestResponse { return r.webhook }
func (r *flowRun) SetWebhook(rr utils.RequestResponse) {
	r.webhook = rr
	r.setModifiedOn(time.Now().In(time.UTC))
}

func (r *flowRun) Extra() utils.JSONFragment {
	return utils.NewJSONFragment([]byte(r.extra))
}
func (r *flowRun) SetExtra(extra json.RawMessage) { r.extra = extra }

func (r *flowRun) CreatedOn() time.Time        { return r.createdOn }
func (r *flowRun) ModifiedOn() time.Time       { return r.modifiedOn }
func (r *flowRun) setModifiedOn(now time.Time) { r.modifiedOn = now }
func (r *flowRun) ExpiresOn() *time.Time       { return r.expiresOn }
func (r *flowRun) TimesOutOn() *time.Time      { return r.timesOutOn }
func (r *flowRun) ExitedOn() *time.Time        { return r.exitedOn }

func (r *flowRun) updateTranslations() {
	if r.flowTranslations != nil {
		r.translations = r.flowTranslations.GetTranslations(r.language)
	} else {
		r.translations = nil
	}
}
func (r *flowRun) SetFlowTranslations(ft flows.FlowTranslations) {
	r.flowTranslations = ft
	r.updateTranslations()
}
func (r *flowRun) SetLanguage(lang utils.Language) {
	r.language = lang
	r.updateTranslations()
}
func (r *flowRun) GetText(uuid flows.UUID, key string, backdown string) string {
	if r.translations == nil {
		return backdown
	}
	return r.translations.GetText(uuid, key, backdown)
}

// NewRun initializes a new context and flow run for the passed in flow and contact
func NewRun(env flows.FlowEnvironment, flow flows.Flow, contact *flows.Contact, parent flows.FlowRun) flows.FlowRun {
	now := time.Now()

	r := &flowRun{
		uuid:        flows.RunUUID(uuid.NewV4().String()),
		flowUUID:    flow.UUID(),
		contactUUID: contact.UUID(),
		results:     flows.NewResults(),
		environment: env,
		contact:     contact,
		flow:        flow,
		status:      flows.StatusActive,
		createdOn:   now,
		modifiedOn:  now,
	}

	// create our new context
	r.context = NewContextForContact(contact, r)

	// set our session
	if parent != nil {
		parentRun := parent.(*flowRun)
		r.parent = newReferenceFromRun(parentRun)
		parentRun.child = newReferenceFromRun(r)

		r.session = parent.Session()
	} else {
		r.session = newSession()
	}
	r.session.AddRun(r)

	return r
}

func (r *flowRun) Resolve(key string) interface{} {
	switch key {

	case "channel":
		return r.Channel()

	case "contact":
		return r.Contact()

	case "extra":
		return r.Extra()

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

func (r *flowRun) Default() interface{} {
	return r
}

// String returns the default string value for this run, which is just our status
func (r *flowRun) String() string {
	return string(r.status)
}

// runReference provides a standalone and serializable version of a run reference. When a run is written
// this is what is written and what will be read (and needed) when resuming that run
type runReference struct {
	uuid flows.RunUUID
	run  *flowRun
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

func (r *runReference) String() string {
	return string(r.Status())
}

func (r *runReference) UUID() flows.RunUUID            { return r.uuid }
func (r *runReference) FlowUUID() flows.FlowUUID       { return r.run.flowUUID }
func (r *runReference) ContactUUID() flows.ContactUUID { return r.run.contactUUID }
func (r *runReference) ChannelUUID() flows.ChannelUUID { return r.run.channelUUID }

func (r *runReference) Results() *flows.Results { return r.run.results }
func (r *runReference) Status() flows.RunStatus { return r.run.status }

func (r *runReference) CreatedOn() time.Time   { return r.run.createdOn }
func (r *runReference) ModifiedOn() time.Time  { return r.run.modifiedOn }
func (r *runReference) ExitedOn() *time.Time   { return r.run.exitedOn }
func (r *runReference) ExpiresOn() *time.Time  { return r.run.expiresOn }
func (r *runReference) TimesOutOn() *time.Time { return r.run.timesOutOn }

func newReferenceFromRun(r *flowRun) *runReference {
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
	run := &flowRun{}
	err := json.Unmarshal(data, run)
	if err == nil {
		// err = run.Validate()
	}
	return run, err
}

type runEnvelope struct {
	UUID    flows.RunUUID     `json:"uuid"`
	Flow    flows.FlowUUID    `json:"flow_uuid"`
	Channel flows.ChannelUUID `json:"channel_uuid"`
	Contact flows.ContactUUID `json:"contact_uuid"`
	Path    []*step           `json:"path"`

	Status flows.RunStatus `json:"status"`

	Input *utils.TypedEnvelope `json:"input,omitempty"`
	Wait  *utils.TypedEnvelope `json:"wait,omitempty"`
	Event *utils.TypedEnvelope `json:"event,omitempty"`

	Parent flows.RunUUID `json:"parent_uuid,omitempty"`
	Child  flows.RunUUID `json:"child_uuid,omitempty"`

	Results *flows.Results        `json:"results,omitempty"`
	Webhook utils.RequestResponse `json:"webhook,omitempty"`
	Extra   json.RawMessage       `json:"extra,omitempty"`

	CreatedOn  time.Time  `json:"created_on"`
	ModifiedOn time.Time  `json:"modified_on"`
	ExpiresOn  *time.Time `json:"expires_on"`
	TimesOutOn *time.Time `json:"timesout_on"`
	ExitedOn   *time.Time `json:"exited_on"`
}

func (r *flowRun) UnmarshalJSON(data []byte) error {
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
	r.extra = envelope.Extra

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
	} else {
		r.results = flows.NewResults()
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

func (r *flowRun) MarshalJSON() ([]byte, error) {
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
	re.Extra = r.extra

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
