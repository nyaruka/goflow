package runs

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/nyaruka/goflow/excellent"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/flows/inputs"
	"github.com/nyaruka/goflow/utils"

	log "github.com/sirupsen/logrus"
)

type flowRun struct {
	uuid        flows.RunUUID
	session     flows.Session
	environment flows.RunEnvironment

	flow    flows.Flow
	contact *flows.Contact

	context types.XValue
	webhook *flows.WebhookCall
	input   flows.Input
	parent  flows.FlowRun

	results flows.Results
	path    Path
	events  []flows.Event
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

func (r *flowRun) UUID() flows.RunUUID               { return r.uuid }
func (r *flowRun) Session() flows.Session            { return r.session }
func (r *flowRun) Environment() flows.RunEnvironment { return r.environment }

func (r *flowRun) Flow() flows.Flow                  { return r.flow }
func (r *flowRun) Contact() *flows.Contact           { return r.contact }
func (r *flowRun) SetContact(contact *flows.Contact) { r.contact = contact }

func (r *flowRun) Context() types.XValue  { return r.context }
func (r *flowRun) Results() flows.Results { return r.results }
func (r *flowRun) Events() []flows.Event  { return r.events }

func (r *flowRun) Exit(status flows.RunStatus) {
	r.SetStatus(status)
	now := time.Now().UTC()
	r.exitedOn = &now
}
func (r *flowRun) Status() flows.RunStatus { return r.status }
func (r *flowRun) SetStatus(status flows.RunStatus) {
	r.status = status
}

// ParentInSession returns the parent of the run within the same session if one exists
func (r *flowRun) ParentInSession() flows.FlowRun { return r.parent }

// Parent returns either the same session parent or if this session was triggered from a trigger_flow action
// in another session, that run
func (r *flowRun) Parent() flows.RunSummary {
	if r.parent == nil {
		return r.session.ParentRun()
	}
	return r.ParentInSession()
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
		event.SetStepUUID(s.UUID())
		r.events = append(r.events, event)
	}

	if !event.FromCaller() {
		r.Session().LogEvent(event)
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
	r.ApplyEvent(step, action, events.NewErrorEvent(err))
}

func (r *flowRun) AddFatalError(step flows.Step, action flows.Action, err error) {
	r.ApplyEvent(step, action, events.NewFatalErrorEvent(err))
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

func (r *flowRun) Webhook() *flows.WebhookCall      { return r.webhook }
func (r *flowRun) SetWebhook(rr *flows.WebhookCall) { r.webhook = rr }

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

	if r.ParentInSession() != nil {
		r.ParentInSession().ResetExpiration(r.expiresOn)
	}
}

func (r *flowRun) ExitedOn() *time.Time { return r.exitedOn }

// EvaluateTemplate evaluates the given template in the context of this run
func (r *flowRun) EvaluateTemplate(template string) (types.XValue, error) {
	return excellent.EvaluateTemplate(r.Environment(), r.Context(), template, RunContextTopLevels)
}

// EvaluateTemplateAsString evaluates the given template as a string in the context of this run
func (r *flowRun) EvaluateTemplateAsString(template string, urlEncode bool) (string, error) {
	return excellent.EvaluateTemplateAsString(r.Environment(), r.Context(), template, urlEncode, RunContextTopLevels)
}

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
func (r *flowRun) Resolve(env utils.Environment, key string) types.XValue {
	switch key {
	case "uuid":
		return types.NewXText(string(r.UUID()))
	case "contact":
		return r.Contact()
	case "flow":
		return r.Flow()
	case "input":
		return r.Input()
	case "webhook":
		return r.Webhook()
	case "status":
		return types.NewXText(string(r.Status()))
	case "results":
		return r.Results()
	case "path":
		return r.path
	case "created_on":
		return types.NewXDateTime(r.CreatedOn())
	case "exited_on":
		if r.exitedOn != nil {
			return types.NewXDateTime(*r.exitedOn)
		}
		return nil
	}

	return types.NewXResolveError(r, key)
}

// Describe returns a representation of this type for error messages
func (r *flowRun) Describe() string { return "run" }

// Reduce is called when this object needs to be reduced to a primitive
func (r *flowRun) Reduce() types.XPrimitive {
	return types.NewXText(string(r.uuid))
}

func (r *flowRun) ToXJSON(env utils.Environment) types.XText {
	return types.ResolveKeys(env, r, "uuid", "contact", "flow", "input", "webhook", "status", "results", "created_on", "exited_on").ToXJSON(env)
}

func (r *flowRun) Snapshot() flows.RunSummary {
	return newRunSummaryFromRun(r)
}

var _ flows.FlowRun = (*flowRun)(nil)
var _ flows.RunSummary = (*flowRun)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type runEnvelope struct {
	UUID   flows.RunUUID          `json:"uuid" validate:"required,uuid4"`
	Flow   *flows.FlowReference   `json:"flow" validate:"required,dive"`
	Path   []*step                `json:"path" validate:"dive"`
	Events []*utils.TypedEnvelope `json:"events,omitempty"`

	Status     flows.RunStatus `json:"status"`
	ParentUUID flows.RunUUID   `json:"parent_uuid,omitempty" validate:"omitempty,uuid4"`

	Results flows.Results        `json:"results,omitempty" validate:"omitempty,dive"`
	Input   *utils.TypedEnvelope `json:"input,omitempty" validate:"omitempty,dive"`
	Webhook *flows.WebhookCall   `json:"webhook,omitempty" validate:"omitempty,dive"`

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
	if r.flow, err = session.Assets().GetFlow(envelope.Flow.UUID); err != nil {
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
			return nil, fmt.Errorf("unable to read input: %s", err)
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

	// read in our events
	r.events = make([]flows.Event, len(envelope.Events))
	for i := range r.events {
		if r.events[i], err = events.EventFromEnvelope(envelope.Events[i]); err != nil {
			return nil, err
		}
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
	re.Flow = r.flow.Reference()
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

	re.Events, err = events.EventsToEnvelopes(r.events)
	if err != nil {
		return nil, err
	}

	return json.Marshal(re)
}
