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

// a run specific environment which allows values to be overridden by the contact
type runEnvironment struct {
	flows.SessionEnvironment
	run *flowRun

	cachedLanguages utils.LanguageList
}

// creates a run environment based on the given run
func newRunEnvironment(base flows.SessionEnvironment, run *flowRun) *runEnvironment {
	env := &runEnvironment{base, run, nil}
	env.refreshLanguagesCache()
	return env
}

func (e *runEnvironment) Timezone() *time.Location {
	contact := e.run.contact

	// if run has a contact with a timezone, that overrides the enviroment's timezone
	if contact != nil && contact.Timezone() != nil {
		return contact.Timezone()
	}
	return e.SessionEnvironment.Timezone()
}

func (e *runEnvironment) Languages() utils.LanguageList {
	// if contact language has changed, rebuild our cached language list
	if e.run.Contact() != nil && e.cachedLanguages[0] != e.run.Contact().Language() {
		e.refreshLanguagesCache()
	}

	return e.cachedLanguages
}

func (e *runEnvironment) refreshLanguagesCache() {
	contact := e.run.contact
	var languages utils.LanguageList

	// if contact has a language, it takes priority
	if contact != nil && contact.Language() != utils.NilLanguage {
		languages = append(languages, contact.Language())
	}

	// next we include any environment languages
	languages = append(languages, e.SessionEnvironment.Languages()...)

	// finally we include the flow native language
	languages = append(languages, e.run.flow.Language())

	e.cachedLanguages = languages.RemoveDuplicates()
}

type flowRun struct {
	uuid flows.RunUUID

	flow    flows.Flow
	contact *flows.Contact

	extra   json.RawMessage
	results *flows.Results
	context flows.Context
	status  flows.RunStatus
	wait    flows.Wait
	webhook *utils.RequestResponse
	input   flows.Input

	parent flows.FlowRunReference
	child  flows.FlowRunReference

	session flows.Session

	path        []flows.Step
	environment flows.SessionEnvironment

	createdOn  time.Time
	modifiedOn time.Time
	expiresOn  *time.Time
	timesOutOn *time.Time
	exitedOn   *time.Time
}

func (r *flowRun) UUID() flows.RunUUID     { return r.uuid }
func (r *flowRun) Flow() flows.Flow        { return r.flow }
func (r *flowRun) Contact() *flows.Contact { return r.contact }

func (r *flowRun) Context() flows.Context                { return r.context }
func (r *flowRun) Environment() flows.SessionEnvironment { return r.environment }
func (r *flowRun) Results() *flows.Results               { return r.results }
func (r *flowRun) Session() flows.Session                { return r.session }

func (r *flowRun) IsComplete() bool {
	return r.status != flows.StatusActive
}
func (r *flowRun) setStatus(status flows.RunStatus) {
	r.status = status
	r.setModifiedOn(time.Now().UTC())
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

func (r *flowRun) ApplyEvent(s flows.Step, a flows.Action, e flows.Event) {
	e.Apply(r)

	fs := s.(*step)
	fs.addEvent(e)

	if !e.FromCaller() {
		r.Session().LogEvent(s, a, e)
		r.setModifiedOn(time.Now().UTC())
	}
}

func (r *flowRun) AddError(step flows.Step, err error) {
	r.ApplyEvent(step, nil, &events.ErrorEvent{Text: err.Error()})
}

func (r *flowRun) Path() []flows.Step { return r.path }
func (r *flowRun) CreateStep(node flows.Node) flows.Step {
	now := time.Now().UTC()
	step := &step{stepUUID: flows.StepUUID(uuid.NewV4().String()), nodeUUID: node.UUID(), arrivedOn: now}
	r.path = append(r.path, step)
	r.setModifiedOn(now)
	return step
}
func (r *flowRun) ClearPath() {
	r.path = nil
}

func (r *flowRun) Webhook() *utils.RequestResponse { return r.webhook }
func (r *flowRun) SetWebhook(rr *utils.RequestResponse) {
	r.webhook = rr
	r.setModifiedOn(time.Now().UTC())
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

func (r *flowRun) GetText(uuid flows.UUID, key string, native string) string {
	textArray := r.GetTextArray(uuid, key, []string{native})
	return textArray[0]
}

func (r *flowRun) GetTextArray(uuid flows.UUID, key string, native []string) []string {
	for _, lang := range r.environment.Languages() {
		if lang == r.Flow().Language() {
			return native
		}

		translations := r.Flow().Translations().GetLanguageTranslations(lang)
		if translations != nil {
			textArray := translations.GetTextArray(uuid, key)
			if textArray != nil && len(textArray) == len(native) {
				return textArray
			}
		}
	}
	return native
}

// NewRun initializes a new context and flow run for the passed in flow and contact
func NewRun(session flows.Session, flow flows.Flow, contact *flows.Contact, parent flows.FlowRun) flows.FlowRun {
	now := time.Now().UTC()

	r := &flowRun{
		uuid:       flows.RunUUID(uuid.NewV4().String()),
		session:    session,
		flow:       flow,
		contact:    contact,
		results:    flows.NewResults(),
		status:     flows.StatusActive,
		createdOn:  now,
		modifiedOn: now,
	}

	r.environment = newRunEnvironment(session.Environment(), r)

	// build our context
	r.context = NewContextForContact(contact, r)

	if parent != nil {
		parentRun := parent.(*flowRun)
		r.parent = newReferenceFromRun(parentRun)
		parentRun.child = newReferenceFromRun(r)
	}

	return r
}

func (r *flowRun) Resolve(key string) interface{} {
	switch key {

	case "contact":
		return r.Contact()

	case "extra":
		return r.Extra()

	case "flow":
		return r.Flow()

	case "input":
		return r.Input()

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

var _ utils.VariableResolver = (*flowRun)(nil)

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

	case "contact":
		return r.Contact()

	case "flow":
		return r.Flow()

	case "input":
		return r.Input()

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

var _ utils.VariableResolver = (*runReference)(nil)

func (r *runReference) String() string {
	return string(r.Status())
}

func (r *runReference) UUID() flows.RunUUID     { return r.uuid }
func (r *runReference) Flow() flows.Flow        { return r.run.flow }
func (r *runReference) Contact() *flows.Contact { return r.run.contact }
func (r *runReference) Input() flows.Input      { return r.run.input }

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

type runEnvelope struct {
	UUID        flows.RunUUID     `json:"uuid"`
	FlowUUID    flows.FlowUUID    `json:"flow_uuid"`
	ContactUUID flows.ContactUUID `json:"contact_uuid"`
	Path        []*step           `json:"path"`

	Status flows.RunStatus `json:"status"`

	Input *utils.TypedEnvelope `json:"input,omitempty"`
	Wait  *utils.TypedEnvelope `json:"wait,omitempty"`

	Parent flows.RunUUID `json:"parent_uuid,omitempty"`
	Child  flows.RunUUID `json:"child_uuid,omitempty"`

	Results *flows.Results         `json:"results,omitempty"`
	Webhook *utils.RequestResponse `json:"webhook,omitempty"`
	Extra   json.RawMessage        `json:"extra,omitempty"`

	CreatedOn  time.Time  `json:"created_on"`
	ModifiedOn time.Time  `json:"modified_on"`
	ExpiresOn  *time.Time `json:"expires_on"`
	TimesOutOn *time.Time `json:"timesout_on"`
	ExitedOn   *time.Time `json:"exited_on"`
}

// ReadRun decodes a run from the passed in JSON
func ReadRun(session flows.Session, data json.RawMessage) (flows.FlowRun, error) {
	r := &flowRun{}
	var envelope runEnvelope
	var err error

	err = json.Unmarshal(data, &envelope)
	if err != nil {
		return nil, err
	}

	r.uuid = envelope.UUID
	r.status = envelope.Status
	r.createdOn = envelope.CreatedOn
	r.modifiedOn = envelope.ModifiedOn
	r.expiresOn = envelope.ExpiresOn
	r.timesOutOn = envelope.TimesOutOn
	r.exitedOn = envelope.ExitedOn
	r.extra = envelope.Extra

	r.flow, err = session.Environment().GetFlow(envelope.FlowUUID)
	if err != nil {
		return nil, err
	}
	r.contact, err = session.Environment().GetContact(envelope.ContactUUID)
	if err != nil {
		return nil, err
	}

	if envelope.Parent != "" {
		r.parent = &runReference{uuid: envelope.Parent}
	}
	if envelope.Child != "" {
		r.child = &runReference{uuid: envelope.Child}
	}

	if envelope.Input != nil {
		r.input, err = inputs.ReadInput(session.Environment(), envelope.Input)
		if err != nil {
			return nil, err
		}
	}

	if envelope.Wait != nil {
		r.wait, err = waits.WaitFromEnvelope(envelope.Wait)
		if err != nil {
			return nil, err
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

	// add ourselves to the environment and save it off
	session.Environment().AddRun(r)
	r.environment = newRunEnvironment(session.Environment(), r)

	// build our context
	r.context = NewContextForContact(r.contact, r)

	return r, nil
}

// ResolveReferences resolves parent/child run references for unmarshaled runs
func (r *flowRun) ResolveReferences(env flows.SessionEnvironment) error {
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

func (r *flowRun) MarshalJSON() ([]byte, error) {
	var re runEnvelope
	var err error

	re.UUID = r.uuid
	re.FlowUUID = r.flow.UUID()
	re.ContactUUID = r.contact.UUID()

	re.Status = r.status
	re.CreatedOn = r.createdOn
	re.ModifiedOn = r.modifiedOn
	re.ExpiresOn = r.expiresOn
	re.TimesOutOn = r.timesOutOn
	re.ExitedOn = r.exitedOn
	re.Results = r.results
	re.Extra = r.extra
	re.Webhook = r.webhook

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

	re.Path = make([]*step, len(r.path))
	for i, s := range r.path {
		re.Path[i] = s.(*step)
	}

	return json.Marshal(re)
}
