package runs

import (
	"encoding/json"
	"time"

	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/excellent"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"

	"github.com/pkg/errors"
)

type flowRun struct {
	uuid        flows.RunUUID
	session     flows.Session
	environment flows.RunEnvironment

	flow    flows.Flow
	parent  flows.FlowRun
	results flows.Results
	path    Path
	events  []flows.Event
	status  flows.RunStatus

	createdOn  time.Time
	modifiedOn time.Time
	expiresOn  *time.Time
	exitedOn   *time.Time

	legacyExtra *legacyExtra
}

// NewRun initializes a new context and flow run for the passed in flow and contact
func NewRun(session flows.Session, flow flows.Flow, parent flows.FlowRun) flows.FlowRun {
	now := utils.Now()
	r := &flowRun{
		uuid:       flows.RunUUID(utils.NewUUID()),
		session:    session,
		flow:       flow,
		parent:     parent,
		results:    flows.NewResults(),
		status:     flows.RunStatusActive,
		events:     make([]flows.Event, 0),
		createdOn:  now,
		modifiedOn: now,
	}

	r.environment = newRunEnvironment(session.Environment(), r)
	r.ResetExpiration(nil)

	r.legacyExtra = newLegacyExtra(r)

	return r
}

func (r *flowRun) UUID() flows.RunUUID               { return r.uuid }
func (r *flowRun) Session() flows.Session            { return r.session }
func (r *flowRun) Environment() flows.RunEnvironment { return r.environment }

func (r *flowRun) Flow() flows.Flow        { return r.flow }
func (r *flowRun) Contact() *flows.Contact { return r.session.Contact() }
func (r *flowRun) Events() []flows.Event   { return r.events }

func (r *flowRun) Results() flows.Results { return r.results }
func (r *flowRun) SaveResult(result *flows.Result) {
	// truncate value if necessary
	if len(result.Value) > r.Environment().MaxValueLength() {
		result.Value = result.Value[0:r.Environment().MaxValueLength()]
	}

	r.results.Save(result)
	r.modifiedOn = utils.Now()

	r.legacyExtra.addResult(result)
}

func (r *flowRun) Exit(status flows.RunStatus) {
	now := utils.Now()

	r.status = status
	r.exitedOn = &now
	r.modifiedOn = now
}
func (r *flowRun) Status() flows.RunStatus { return r.status }
func (r *flowRun) SetStatus(status flows.RunStatus) {
	r.status = status
	r.modifiedOn = utils.Now()
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

func (r *flowRun) LogEvent(s flows.Step, event flows.Event) {
	if s != nil {
		event.SetStepUUID(s.UUID())
	}

	r.events = append(r.events, event)
	r.modifiedOn = utils.Now()
}

func (r *flowRun) LogError(step flows.Step, err error) {
	r.LogEvent(step, events.NewErrorEvent(err))
}

func (r *flowRun) Path() []flows.Step { return r.path }
func (r *flowRun) CreateStep(node flows.Node) flows.Step {
	now := utils.Now()
	step := NewStep(node, now)
	r.path = append(r.path, step)
	r.modifiedOn = now
	return step
}

func (r *flowRun) PathLocation() (flows.Step, flows.Node, error) {
	if r.Path() == nil {
		return nil, nil, errors.Errorf("run has no location as path is empty")
	}

	step := r.Path()[len(r.Path())-1]

	// check that we still have a node for this step
	node := r.Flow().GetNode(step.NodeUUID())
	if node == nil {
		return nil, nil, errors.Errorf("run is located at a flow node that no longer exists")
	}

	return step, node, nil
}

func (r *flowRun) CreatedOn() time.Time  { return r.createdOn }
func (r *flowRun) ModifiedOn() time.Time { return r.modifiedOn }
func (r *flowRun) ExpiresOn() *time.Time { return r.expiresOn }
func (r *flowRun) ResetExpiration(from *time.Time) {
	if r.Flow().ExpireAfterMinutes() >= 0 {
		if from == nil {
			now := utils.Now()
			from = &now
		}

		expiresAfterMinutes := time.Duration(r.Flow().ExpireAfterMinutes())
		expiresOn := from.Add(expiresAfterMinutes * time.Minute)

		r.expiresOn = &expiresOn
		r.modifiedOn = utils.Now()
	}

	if r.ParentInSession() != nil {
		r.ParentInSession().ResetExpiration(r.expiresOn)
	}
}

func (r *flowRun) ExitedOn() *time.Time { return r.exitedOn }

// Context returns the overall context for expression evaluation
func (r *flowRun) RootContext(env utils.Environment) map[string]types.XValue {
	var urns, fields types.XValue
	if r.Contact() != nil {
		urns = flows.ContextFunc(env, r.Contact().URNs().MapContext)
		fields = flows.Context(env, r.Contact().Fields())
	}

	var child = newRelatedRunContext(r.Session().GetCurrentChild(r))
	var parent = newRelatedRunContext(r.Parent())

	return map[string]types.XValue{
		// the available runs
		"run":    flows.Context(env, r),
		"child":  flows.Context(env, child),
		"parent": flows.Context(env, parent),

		// shortcuts to things on the current run
		"contact": flows.Context(env, r.Contact()),
		"results": flows.ContextFunc(env, r.Results().SimpleContext),
		"urns":    urns,
		"fields":  fields,
		"webhook": r.lastWebhookResponse(),

		// other
		"trigger":      flows.Context(env, r.Session().Trigger()),
		"input":        flows.Context(env, r.Session().Input()),
		"legacy_extra": r.legacyExtra.ToXValue(env),
	}
}

func (r *flowRun) lastWebhookResponse() types.XValue {
	for e := len(r.events) - 1; e >= 0; e-- {
		switch typed := r.events[e].(type) {
		case *events.WebhookCalledEvent:
			return types.JSONToXValue(flows.ExtractResponseBody(typed.Response))
		default:
			continue
		}
	}
	return nil
}

// Context returns the properties available in expressions
func (r *flowRun) Context(env utils.Environment) map[string]types.XValue {
	var exitedOn types.XValue
	if r.exitedOn != nil {
		exitedOn = types.NewXDateTime(*r.exitedOn)
	}

	return map[string]types.XValue{
		"__default__": types.NewXText(formatRunSummary(env, r)),
		"uuid":        types.NewXText(string(r.UUID())),
		"contact":     flows.Context(env, r.Contact()),
		"flow":        flows.Context(env, r.Flow()),
		"status":      types.NewXText(string(r.Status())),
		"results":     flows.Context(env, r.Results()),
		"path":        r.path.ToXValue(env),
		"created_on":  types.NewXDateTime(r.CreatedOn()),
		"exited_on":   exitedOn,
	}
}

// EvaluateTemplate evaluates the given template in the context of this run
func (r *flowRun) EvaluateTemplateValue(template string) (types.XValue, error) {
	context := types.NewXObject(r.RootContext(r.Environment()))

	return excellent.EvaluateTemplateValue(r.Environment(), context, template)
}

// EvaluateTemplateAsString evaluates the given template as a string in the context of this run
func (r *flowRun) EvaluateTemplate(template string) (string, error) {
	context := types.NewXObject(r.RootContext(r.Environment()))

	return excellent.EvaluateTemplate(r.Environment(), context, template)
}

// get the ordered list of languages to be used for localization in this run
func (r *flowRun) getLanguages() []utils.Language {
	// TODO cache this this?

	contact := r.Contact()
	languages := make([]utils.Language, 0, 3)

	// if contact has a allowed language, it takes priority
	if contact != nil && contact.Language() != utils.NilLanguage {
		for _, l := range r.Environment().AllowedLanguages() {
			if l == contact.Language() {
				languages = append(languages, contact.Language())
				break
			}
		}
	}

	// next we include the default language if it's different to the contact language
	defaultLanguage := r.Environment().DefaultLanguage()
	if defaultLanguage != utils.NilLanguage && defaultLanguage != contact.Language() {
		languages = append(languages, defaultLanguage)
	}

	// finally we include the flow native language if it isn't an allowed language - because it's the only
	// one guaranteed to have translations
	return append(languages, r.flow.Language())
}

func (r *flowRun) GetText(uuid utils.UUID, key string, native string) string {
	textArray := r.GetTextArray(uuid, key, []string{native})
	return textArray[0]
}

func (r *flowRun) GetTextArray(uuid utils.UUID, key string, native []string) []string {
	return r.GetTranslatedTextArray(uuid, key, native, r.getLanguages())
}

func (r *flowRun) GetTranslatedTextArray(uuid utils.UUID, key string, native []string, languages []utils.Language) []string {
	if languages == nil {
		languages = r.getLanguages()
	}

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

func (r *flowRun) Snapshot() flows.RunSummary {
	return newRunSummaryFromRun(r)
}

var _ flows.RunSummary = (*flowRun)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type runEnvelope struct {
	UUID       flows.RunUUID         `json:"uuid" validate:"required,uuid4"`
	Flow       *assets.FlowReference `json:"flow" validate:"required,dive"`
	Path       []*step               `json:"path" validate:"dive"`
	Events     []json.RawMessage     `json:"events,omitempty"`
	Results    flows.Results         `json:"results,omitempty" validate:"omitempty,dive"`
	Status     flows.RunStatus       `json:"status" validate:"required"`
	ParentUUID flows.RunUUID         `json:"parent_uuid,omitempty" validate:"omitempty,uuid4"`

	CreatedOn  time.Time  `json:"created_on" validate:"required"`
	ModifiedOn time.Time  `json:"modified_on" validate:"required"`
	ExpiresOn  *time.Time `json:"expires_on"`
	ExitedOn   *time.Time `json:"exited_on"`
}

// ReadRun decodes a run from the passed in JSON. Parent run UUID is returned separately as the
// run in question might be loaded yet from the session.
func ReadRun(session flows.Session, data json.RawMessage) (flows.FlowRun, error) {
	e := &runEnvelope{}
	var err error

	if err = utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, errors.Wrap(err, "unable to read run")
	}

	r := &flowRun{
		session:    session,
		uuid:       e.UUID,
		status:     e.Status,
		createdOn:  e.CreatedOn,
		modifiedOn: e.ModifiedOn,
		expiresOn:  e.ExpiresOn,
		exitedOn:   e.ExitedOn,
	}

	// lookup flow
	if r.flow, err = session.Assets().Flows().Get(e.Flow.UUID); err != nil {
		return nil, errors.Wrapf(err, "unable to load %s", e.Flow)
	}

	// lookup parent run
	if e.ParentUUID != "" {
		if r.parent, err = session.GetRun(e.ParentUUID); err != nil {
			return nil, err
		}
	}

	if e.Results != nil {
		r.results = e.Results
	} else {
		r.results = flows.NewResults()
	}

	// read in our path
	r.path = make([]flows.Step, len(e.Path))
	for i, step := range e.Path {
		r.path[i] = step
	}

	// read in our events
	r.events = make([]flows.Event, len(e.Events))
	for i := range r.events {
		if r.events[i], err = events.ReadEvent(e.Events[i]); err != nil {
			return nil, errors.Wrap(err, "unable to read event")
		}
	}

	// create a run specific environment and context
	r.environment = newRunEnvironment(session.Environment(), r)
	r.legacyExtra = newLegacyExtra(r)

	return r, nil
}

// MarshalJSON marshals this flow run into JSON
func (r *flowRun) MarshalJSON() ([]byte, error) {
	var err error

	e := &runEnvelope{
		UUID:       r.uuid,
		Flow:       r.flow.Reference(),
		Status:     r.status,
		CreatedOn:  r.createdOn,
		ModifiedOn: r.modifiedOn,
		ExpiresOn:  r.expiresOn,
		ExitedOn:   r.exitedOn,
		Results:    r.results,
	}

	if r.parent != nil {
		e.ParentUUID = r.parent.UUID()
	}

	e.Path = make([]*step, len(r.path))
	for i, s := range r.path {
		e.Path[i] = s.(*step)
	}

	e.Events = make([]json.RawMessage, len(r.events))
	for i := range r.events {
		if e.Events[i], err = json.Marshal(r.events[i]); err != nil {
			return nil, errors.Wrapf(err, "unable to marshal event[type=%s]", r.events[i].Type())
		}
	}

	return json.Marshal(e)
}
