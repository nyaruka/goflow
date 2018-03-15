package runs

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/flows/inputs"
	"github.com/nyaruka/goflow/flows/triggers"
	"github.com/nyaruka/goflow/utils"

	log "github.com/sirupsen/logrus"
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

func (e *runEnvironment) Locations() (*utils.LocationHierarchy, error) {
	sessionAssets := e.run.Session().Assets()
	if sessionAssets.HasLocations() {
		return sessionAssets.GetLocationHierarchy()
	}

	return nil, nil
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
	parent  flows.FlowRun

	results flows.Results
	path    []flows.Step
	status  flows.RunStatus

	createdOn time.Time
	expiresOn *time.Time
	exitedOn  *time.Time
}

// NewRun initializes a new context and flow run for the passed in flow and contact
func NewRun(session flows.Session, flow flows.Flow, contact *flows.Contact, parent flows.FlowRun) flows.FlowRun {
	r := &flowRun{
		uuid:      flows.RunUUID(utils.NewUUID()),
		session:   session,
		flow:      flow,
		contact:   contact,
		results:   flows.NewResults(),
		status:    flows.RunStatusActive,
		createdOn: time.Now().UTC(),
	}

	r.environment = newRunEnvironment(session.Environment(), r)
	r.context = newRunContext(r)
	r.parent = parent

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
func (r *flowRun) Results() flows.Results          { return r.results }

func (r *flowRun) Exit(status flows.RunStatus) {
	r.SetStatus(status)
	now := time.Now().UTC()
	r.exitedOn = &now
}
func (r *flowRun) Status() flows.RunStatus { return r.status }
func (r *flowRun) SetStatus(status flows.RunStatus) {
	r.status = status
}

// SessionParent returns the parent of the run within the same session if one exists
func (r *flowRun) SessionParent() flows.FlowRun { return r.parent }

// Parent returns either the same session parent or if this session was triggered from a trigger_flow action
// in another session, that run
func (r *flowRun) Parent() flows.RunSummary {
	if r.parent == nil && r.session.Trigger() != nil && r.session.Trigger().Type() == triggers.TypeFlowAction {
		runTrigger := r.session.Trigger().(*triggers.FlowActionTrigger)

		return runTrigger.Run()
	}
	return r.SessionParent()
}

func (r *flowRun) Ancestors() []flows.FlowRun {
	ancestors := make([]flows.FlowRun, 0)
	if r.parent != nil {
		run := r.parent.(*flowRun)
		ancestors = append(ancestors, run)

		for {
			if run.parent != nil {
				run = run.parent.(*flowRun)
				ancestors = append(ancestors, run)
			} else {
				break
			}
		}
	}

	return ancestors
}

func (r *flowRun) Input() flows.Input         { return r.input }
func (r *flowRun) SetInput(input flows.Input) { r.input = input }

func (r *flowRun) ApplyEvent(s flows.Step, action flows.Action, event flows.Event) error {
	if err := event.Apply(r); err != nil {
		return fmt.Errorf("unable to apply event[type=%s]: %s", event.Type(), err)
	}

	if s != nil {
		fs := s.(*step)
		fs.addEvent(event)
	}

	if !event.FromCaller() {
		r.Session().LogEvent(s, action, event)
	}

	if log.GetLevel() >= log.DebugLevel {
		var origin string
		if event.FromCaller() {
			origin = "caller"
		} else {
			origin = "engine"
		}
		eventEnvelope, _ := utils.EnvelopeFromTyped(event)
		eventJSON, _ := json.Marshal(eventEnvelope)
		log.WithField("event_type", event.Type()).WithField("payload", string(eventJSON)).WithField("run", r.UUID()).Debugf("%s event applied", origin)
	}

	return nil
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
	step := &step{stepUUID: flows.StepUUID(utils.NewUUID()), nodeUUID: node.UUID(), arrivedOn: now}
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

	if r.SessionParent() != nil {
		r.SessionParent().ResetExpiration(r.expiresOn)
	}
}

func (r *flowRun) ExitedOn() *time.Time { return r.exitedOn }

func (r *flowRun) GetText(uuid utils.UUID, key string, native string) string {
	textArray := r.GetTextArray(uuid, key, []string{native})
	return textArray[0]
}

func (r *flowRun) GetTextArray(uuid utils.UUID, key string, native []string) []string {
	return r.GetTranslatedTextArray(uuid, key, native, r.environment.Languages())
}

func (r *flowRun) GetTranslatedTextArray(uuid utils.UUID, key string, native []string, languages utils.LanguageList) []string {
	for _, lang := range languages {
		if lang == r.Flow().Language() {
			return native
		}

		translations := r.Flow().Localization().GetTranslations(lang)
		if translations != nil {
			textArray := translations.GetTextArray(uuid, key)
			if textArray == nil {
				return native
			}

			merged := make([]string, len(native))
			for s := range native {
				if textArray[s] != "" {
					merged[s] = textArray[s]
				} else {
					merged[s] = native[s]
				}
			}
			return merged
		}
	}
	return native
}

// Resolve resolves the given key when this run is referenced in an expression
func (r *flowRun) Resolve(key string) interface{} {
	switch key {
	case "uuid":
		return r.UUID()
	case "contact":
		return r.Contact()
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

	return fmt.Errorf("no field '%s' on run", key)
}

// Default returns the value of this run when it is the result of an expression
func (r *flowRun) Default() interface{} {
	return r
}

// String returns the default string value for this run, which is just our UUID
func (r *flowRun) String() string {
	return string(r.uuid)
}

func (r *flowRun) Snapshot() flows.RunSummary {
	return flows.NewRunSummaryFromRun(r)
}

var _ utils.VariableResolver = (*flowRun)(nil)
var _ flows.FlowRun = (*flowRun)(nil)
var _ flows.RunSummary = (*flowRun)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type runEnvelope struct {
	UUID     flows.RunUUID  `json:"uuid" validate:"required,uuid4"`
	FlowUUID flows.FlowUUID `json:"flow_uuid" validate:"required,uuid4"`
	Path     []*step        `json:"path" validate:"dive"`

	Status     flows.RunStatus `json:"status"`
	ParentUUID flows.RunUUID   `json:"parent_uuid,omitempty" validate:"omitempty,uuid4"`

	Results flows.Results          `json:"results,omitempty" validate:"omitempty,dive"`
	Input   *utils.TypedEnvelope   `json:"input,omitempty" validate:"omitempty,dive"`
	Webhook *utils.RequestResponse `json:"webhook,omitempty" validate:"omitempty,dive"`

	CreatedOn time.Time  `json:"created_on"`
	ExpiresOn *time.Time `json:"expires_on"`
	ExitedOn  *time.Time `json:"exited_on"`
}

// ReadRun decodes a run from the passed in JSON. Parent run UUID is returned separately as the
// run in question might be loaded yet from the session.
func ReadRun(session flows.Session, data json.RawMessage) (flows.FlowRun, error) {
	r := &flowRun{}
	var envelope runEnvelope
	var err error

	if err = utils.UnmarshalAndValidate(data, &envelope, "run"); err != nil {
		return nil, err
	}

	r.session = session
	r.contact = session.Contact()
	r.uuid = envelope.UUID
	r.status = envelope.Status
	r.webhook = envelope.Webhook
	r.createdOn = envelope.CreatedOn
	r.expiresOn = envelope.ExpiresOn
	r.exitedOn = envelope.ExitedOn

	// lookup flow
	if r.flow, err = session.Assets().GetFlow(envelope.FlowUUID); err != nil {
		return nil, err
	}

	// lookup parent run
	if envelope.ParentUUID != "" {
		if r.parent, err = session.GetRun(envelope.ParentUUID); err != nil {
			return nil, err
		}
	}

	if envelope.Input != nil {
		if r.input, err = inputs.ReadInput(session, envelope.Input); err != nil {
			return nil, err
		}
	}

	if envelope.Results != nil {
		r.results = envelope.Results
	} else {
		r.results = flows.NewResults()
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

// MarshalJSON marshals this flow run into JSON
func (r *flowRun) MarshalJSON() ([]byte, error) {
	var re runEnvelope
	var err error

	re.UUID = r.uuid
	re.FlowUUID = r.flow.UUID()
	re.Status = r.status
	re.CreatedOn = r.createdOn
	re.ExpiresOn = r.expiresOn
	re.ExitedOn = r.exitedOn
	re.Results = r.results
	re.Webhook = r.webhook

	if r.parent != nil {
		re.ParentUUID = r.parent.UUID()
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
