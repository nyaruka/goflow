package runs

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/flows/inputs"
	"github.com/nyaruka/goflow/utils"
	uuid "github.com/satori/go.uuid"
)

// a run specific environment which allows values to be overridden by the contact
type runEnvironment struct {
	utils.Environment
	run *flowRun

	cachedLanguages utils.LanguageList
}

// creates a run environment based on the given run
func newRunEnvironment(base utils.Environment, run *flowRun) *runEnvironment {
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
	return e.run.Session().Environment().Timezone()
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
	languages = append(languages, e.run.Session().Environment().Languages()...)

	// finally we include the flow native language
	languages = append(languages, e.run.flow.Language())

	e.cachedLanguages = languages.RemoveDuplicates()
}

type flowRun struct {
	uuid        flows.RunUUID
	session     flows.Session
	environment *runEnvironment

	flow    flows.Flow
	contact *flows.Contact
	extra   utils.JSONFragment

	context utils.VariableResolver
	webhook *utils.RequestResponse
	input   flows.Input

	parent flows.FlowRunReference
	child  flows.FlowRunReference

	results *flows.Results
	path    []flows.Step
	status  flows.RunStatus

	createdOn  time.Time
	expiresOn  *time.Time
	timesOutOn *time.Time
	exitedOn   *time.Time
}

// NewRun initializes a new context and flow run for the passed in flow and contact
func NewRun(session flows.Session, flow flows.Flow, contact *flows.Contact, parent flows.FlowRun) flows.FlowRun {
	r := &flowRun{
		uuid:      flows.RunUUID(uuid.NewV4().String()),
		session:   session,
		flow:      flow,
		contact:   contact,
		results:   flows.NewResults(),
		status:    flows.RunStatusActive,
		createdOn: time.Now().UTC(),
	}

	r.environment = newRunEnvironment(session.Environment(), r)
	r.context = newRunContext(r)

	if parent != nil {
		parentRun := parent.(*flowRun)
		r.parent = newReferenceFromRun(parentRun)
		parentRun.child = newReferenceFromRun(r)
	}

	r.ResetExpiration(nil)

	return r
}

func (r *flowRun) UUID() flows.RunUUID            { return r.uuid }
func (r *flowRun) Session() flows.Session         { return r.session }
func (r *flowRun) Environment() utils.Environment { return r.environment }

func (r *flowRun) Flow() flows.Flow                  { return r.flow }
func (r *flowRun) Contact() *flows.Contact           { return r.contact }
func (r *flowRun) SetContact(contact *flows.Contact) { r.contact = contact }

func (r *flowRun) Context() utils.VariableResolver { return r.context }
func (r *flowRun) Results() *flows.Results         { return r.results }

func (r *flowRun) Exit(status flows.RunStatus) {
	r.SetStatus(status)
	now := time.Now().UTC()
	r.exitedOn = &now
}
func (r *flowRun) Status() flows.RunStatus { return r.status }
func (r *flowRun) SetStatus(status flows.RunStatus) {
	r.status = status
}

func (r *flowRun) Parent() flows.FlowRunReference { return r.parent }
func (r *flowRun) Child() flows.FlowRunReference  { return r.child }

func (r *flowRun) Input() flows.Input         { return r.input }
func (r *flowRun) SetInput(input flows.Input) { r.input = input }

func (r *flowRun) ApplyEvent(s flows.Step, action flows.Action, event flows.Event) {
	event.Apply(r)

	fs := s.(*step)
	fs.addEvent(event)

	if !event.FromCaller() {
		r.Session().LogEvent(s, action, event)
	}

	// eventEnvelope, _ := utils.EnvelopeFromTyped(event)
	// eventJSON, _ := json.Marshal(eventEnvelope)
	// fmt.Printf("⚡︎ in run %s: %s\n", r.UUID(), string(eventJSON))
}

func (r *flowRun) AddError(step flows.Step, action flows.Action, err error) {
	r.ApplyEvent(step, action, &events.ErrorEvent{Text: err.Error(), Fatal: false})
}

func (r *flowRun) AddFatalError(step flows.Step, action flows.Action, err error) {
	r.ApplyEvent(step, action, &events.ErrorEvent{Text: err.Error(), Fatal: true})
}

func (r *flowRun) Path() []flows.Step { return r.path }
func (r *flowRun) CreateStep(node flows.Node) flows.Step {
	now := time.Now().UTC()
	step := &step{stepUUID: flows.StepUUID(uuid.NewV4().String()), nodeUUID: node.UUID(), arrivedOn: now}
	r.path = append(r.path, step)
	return step
}

func (r *flowRun) PathLocation() (flows.Step, flows.Node, error) {
	if r.Path() == nil {
		return nil, nil, fmt.Errorf("run has no location as path is empty")
	}

	step := r.Path()[len(r.Path())-1]

	// check that we still have a node for this step
	node := r.Flow().GetNode(step.NodeUUID())
	if node == nil {
		return nil, nil, fmt.Errorf("run is located at a flow node that no longer exists")
	}

	return step, node, nil
}

func (r *flowRun) Webhook() *utils.RequestResponse      { return r.webhook }
func (r *flowRun) SetWebhook(rr *utils.RequestResponse) { r.webhook = rr }

func (r *flowRun) Extra() utils.JSONFragment         { return r.extra }
func (r *flowRun) SetExtra(extra utils.JSONFragment) { r.extra = extra }

func (r *flowRun) CreatedOn() time.Time  { return r.createdOn }
func (r *flowRun) ExpiresOn() *time.Time { return r.expiresOn }
func (r *flowRun) ResetExpiration(from *time.Time) {
	if r.Flow().ExpireAfterMinutes() >= 0 {
		if from == nil {
			now := time.Now().UTC()
			from = &now
		}

		expiresAfterMinutes := time.Duration(r.Flow().ExpireAfterMinutes())
		expiresOn := from.Add(expiresAfterMinutes * time.Minute)

		r.expiresOn = &expiresOn
	}

	if r.Parent() != nil {
		r.Parent().ResetExpiration(r.expiresOn)
	}
}

func (r *flowRun) TimesOutOn() *time.Time { return r.timesOutOn }
func (r *flowRun) ExitedOn() *time.Time   { return r.exitedOn }

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

	case "webhook":
		return r.Webhook()

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

func (r *runReference) CreatedOn() time.Time            { return r.run.createdOn }
func (r *runReference) ExpiresOn() *time.Time           { return r.run.expiresOn }
func (r *runReference) ResetExpiration(from *time.Time) { r.run.ResetExpiration(from) }
func (r *runReference) ExitedOn() *time.Time            { return r.run.exitedOn }
func (r *runReference) TimesOutOn() *time.Time          { return r.run.timesOutOn }

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
	Parent flows.RunUUID   `json:"parent_uuid,omitempty"`
	Child  flows.RunUUID   `json:"child_uuid,omitempty"`

	Results *flows.Results         `json:"results,omitempty"`
	Input   *utils.TypedEnvelope   `json:"input,omitempty"`
	Webhook *utils.RequestResponse `json:"webhook,omitempty"`
	Extra   json.RawMessage        `json:"extra,omitempty"`

	CreatedOn  time.Time  `json:"created_on"`
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

	r.session = session
	r.uuid = envelope.UUID
	r.status = envelope.Status
	r.createdOn = envelope.CreatedOn
	r.expiresOn = envelope.ExpiresOn
	r.timesOutOn = envelope.TimesOutOn
	r.exitedOn = envelope.ExitedOn
	r.extra = utils.JSONFragment(envelope.Extra)

	// TODO runs with different contact to the session?
	r.contact = session.Contact()

	r.flow, err = session.Assets().GetFlow(envelope.FlowUUID)
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
		r.input, err = inputs.ReadInput(session, envelope.Input)
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

	// create a run specific environment and context
	r.environment = newRunEnvironment(session.Environment(), r)
	r.context = newRunContext(r)

	return r, nil
}

// resolves parent/child run references for unmarshaled runs
func ResolveReferences(session flows.Session, runs []flows.FlowRun) error {
	for _, run := range runs {
		r := run.(*flowRun)

		if r.parent != nil {
			parent, err := session.GetRun(r.parent.UUID())
			if err != nil {
				return err
			}
			r.parent = newReferenceFromRun(parent.(*flowRun))
		}

		if r.child != nil {
			child, err := session.GetRun(r.child.UUID())
			if err != nil {
				return err
			}
			r.child = newReferenceFromRun(child.(*flowRun))
		}
	}

	return nil
}

func (r *flowRun) MarshalJSON() ([]byte, error) {
	var re runEnvelope
	var err error

	re.UUID = r.uuid
	re.FlowUUID = r.flow.UUID()
	re.ContactUUID = r.contact.UUID()
	re.Extra, _ = json.Marshal(r.extra)
	re.Status = r.status
	re.CreatedOn = r.createdOn
	re.ExpiresOn = r.expiresOn
	re.TimesOutOn = r.timesOutOn
	re.ExitedOn = r.exitedOn
	re.Results = r.results
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

	re.Path = make([]*step, len(r.path))
	for i, s := range r.path {
		re.Path[i] = s.(*step)
	}

	return json.Marshal(re)
}
