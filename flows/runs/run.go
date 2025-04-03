package runs

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/nyaruka/gocommon/dates"
	"github.com/nyaruka/gocommon/i18n"
	"github.com/nyaruka/gocommon/jsonx"
	"github.com/nyaruka/gocommon/stringsx"
	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/flows"
	"github.com/nyaruka/goflow/flows/events"
	"github.com/nyaruka/goflow/utils"
)

type run struct {
	uuid    flows.RunUUID
	session flows.Session

	flow    flows.Flow
	flowRef *assets.FlowReference

	parent  flows.Run
	locals  *flows.Locals
	results flows.Results
	path    Path
	events  []flows.Event
	status  flows.RunStatus

	createdOn  time.Time
	modifiedOn time.Time
	exitedOn   *time.Time

	webhook     *flows.WebhookCall
	legacyExtra *legacyExtra
}

// NewRun initializes a new context and flow run for the passed in flow and contact
func NewRun(session flows.Session, flow flows.Flow, parent flows.Run) flows.Run {
	now := dates.Now()
	r := &run{
		uuid:       flows.RunUUID(uuids.NewV4()),
		session:    session,
		flow:       flow,
		flowRef:    flow.Reference(true),
		parent:     parent,
		locals:     flows.NewLocals(),
		results:    flows.NewResults(),
		status:     flows.RunStatusActive,
		events:     make([]flows.Event, 0),
		createdOn:  now,
		modifiedOn: now,
	}

	r.webhook = nil
	r.legacyExtra = newLegacyExtra(r)

	return r
}

func (r *run) UUID() flows.RunUUID    { return r.uuid }
func (r *run) Session() flows.Session { return r.session }

func (r *run) Flow() flows.Flow                     { return r.flow }
func (r *run) FlowReference() *assets.FlowReference { return r.flowRef }
func (r *run) Contact() *flows.Contact              { return r.session.Contact() }
func (r *run) Events() []flows.Event                { return r.events }

func (r *run) Locals() *flows.Locals  { return r.locals }
func (r *run) Results() flows.Results { return r.results }
func (r *run) SetResult(result *flows.Result) (*flows.Result, bool) {
	// truncate value if necessary
	result.Value = stringsx.Truncate(result.Value, r.session.Engine().Options().MaxResultChars)

	r.modifiedOn = dates.Now()
	r.legacyExtra.addResult(result)

	return r.results.Save(result)
}

func (r *run) Exit(status flows.RunStatus) {
	now := dates.Now()

	r.status = status
	r.exitedOn = &now
	r.modifiedOn = now
}
func (r *run) Status() flows.RunStatus { return r.status }
func (r *run) SetStatus(status flows.RunStatus) {
	r.status = status
	r.modifiedOn = dates.Now()
}

func (r *run) Webhook() *flows.WebhookCall { return r.webhook }
func (r *run) SetWebhook(call *flows.WebhookCall) {
	r.webhook = call
}

// ParentInSession returns the parent of the run within the same session if one exists
func (r *run) ParentInSession() flows.Run { return r.parent }

// Parent returns either the same session parent or if this session was triggered from a trigger_flow action
// in another session, that run
func (r *run) Parent() flows.RunSummary {
	if r.parent == nil {
		return r.session.ParentRun()
	}
	return r.ParentInSession()
}

func (r *run) Ancestors() []flows.Run {
	ancestors := make([]flows.Run, 0)
	if r.parent != nil {
		pr := r.parent.(*run)
		ancestors = append(ancestors, pr)

		for {
			if pr.parent != nil {
				pr = pr.parent.(*run)
				ancestors = append(ancestors, pr)
			} else {
				break
			}
		}
	}

	return ancestors
}

func (r *run) LogEvent(s flows.Step, event flows.Event) {
	if s != nil {
		event.SetStepUUID(s.UUID())
	}

	r.events = append(r.events, event)
	r.modifiedOn = dates.Now()
}

// find the first event matching the given step UUID and type
func (r *run) findEvent(stepUUID flows.StepUUID, eType string) flows.Event {
	for _, e := range r.events {
		if (stepUUID == "" || e.StepUUID() == stepUUID) && e.Type() == eType {
			return e
		}
	}
	return nil
}

func (r *run) ReceivedInput() bool {
	return r.findEvent("", events.TypeMsgReceived) != nil
}

func (r *run) Path() []flows.Step { return r.path }
func (r *run) CreateStep(node flows.Node) flows.Step {
	now := dates.Now()
	step := NewStep(node, now)
	r.path = append(r.path, step)
	r.modifiedOn = now
	return step
}

func (r *run) PathLocation() (flows.Step, flows.Node, error) {
	if r.Path() == nil {
		return nil, nil, fmt.Errorf("run has no location as path is empty")
	}

	step := r.Path()[len(r.Path())-1]

	// check that we still have a node for this step
	var node flows.Node
	if r.Flow() != nil {
		node = r.Flow().GetNode(step.NodeUUID())
	}
	if node == nil {
		return nil, nil, fmt.Errorf("run is located at a flow node that no longer exists")
	}

	return step, node, nil
}

func (r *run) CreatedOn() time.Time  { return r.createdOn }
func (r *run) ModifiedOn() time.Time { return r.modifiedOn }
func (r *run) ExitedOn() *time.Time  { return r.exitedOn }

// RootContext returns the root context for expression evaluation
//
//	contact:contact -> the contact
//	fields:fields -> the custom field values of the contact
//	urns:urns -> the URN values of the contact
//	locals:locals -> the current run local variables
//	results:results -> the current run results
//	input:input -> the current input from the contact
//	run:run -> the current run
//	child:related_run -> the last child run
//	parent:related_run -> the parent of the run
//	ticket:ticket -> the open ticket for the contact
//	webhook:webhook -> the last webhook call (reset after a wait)
//	node:node -> the current node
//	globals:globals -> the global values
//	trigger:trigger -> the trigger that started this session
//	resume:resume -> the current resume that continued this session
//
// @context root
func (r *run) RootContext(env envs.Environment) map[string]types.XValue {
	var urns, fields, ticket, node types.XValue
	if r.Contact() != nil {
		urns = flows.ContextFunc(env, r.Contact().URNs().MapContext)
		fields = flows.Context(env, r.Contact().Fields())

		if r.Contact().Ticket() != nil {
			ticket = flows.Context(env, r.Contact().Ticket())
		}
	}

	var child, parent *relatedRunContext
	if r.Session().GetCurrentChild(r) != nil {
		child = newRelatedRunContext(r.Session().GetCurrentChild(r))
	}
	if r.Parent() != nil {
		parent = newRelatedRunContext(r.Parent())
	}

	_, n, _ := r.PathLocation()
	if n != nil {
		node = flows.ContextFunc(env, r.nodeContext)
	}

	return map[string]types.XValue{
		// the available runs
		"run":    flows.Context(env, r),
		"child":  flows.Context(env, child),
		"parent": flows.Context(env, parent),

		// shortcuts to things on the current run or contact
		"contact": flows.Context(env, r.Contact()),
		"locals":  flows.Context(env, r.Locals()),
		"results": flows.Context(env, r.Results()),
		"urns":    urns,
		"fields":  fields,
		"ticket":  ticket,

		// other
		"trigger":      flows.Context(env, r.Session().Trigger()),
		"resume":       flows.Context(env, r.Session().CurrentResume()),
		"input":        flows.Context(env, r.Session().Input()),
		"globals":      flows.Context(env, r.Session().Assets().Globals()),
		"webhook":      flows.Context(env, r.webhook),
		"node":         node,
		"legacy_extra": r.legacyExtra.ToXValue(env),
	}
}

// Context returns the properties available in expressions
//
//	__default__:text -> the contact name and flow UUID
//	uuid:text -> the UUID of the run
//	contact:contact -> the contact of the run
//	flow:flow -> the flow of the run
//	status:text -> the status of the run
//	locals:locals -> the local variables of the run
//	results:results -> the results saved by the run
//	created_on:datetime -> the creation date of the run
//	exited_on:datetime -> the exit date of the run
//
// @context run
func (r *run) Context(env envs.Environment) map[string]types.XValue {
	var exitedOn types.XValue
	if r.exitedOn != nil {
		exitedOn = types.NewXDateTime(*r.exitedOn)
	}

	return map[string]types.XValue{
		"__default__": types.NewXText(FormatRunSummary(env, r)),
		"uuid":        types.NewXText(string(r.UUID())),
		"contact":     flows.Context(env, r.Contact()),
		"flow":        flows.Context(env, r.Flow()),
		"status":      types.NewXText(string(r.Status())),
		"locals":      flows.Context(env, r.Locals()),
		"results":     flows.Context(env, r.Results()),
		"path":        r.path.ToXValue(env),
		"created_on":  types.NewXDateTime(r.CreatedOn()),
		"exited_on":   exitedOn,
	}
}

// returns the context representation of the current node
//
//	uuid:text -> the UUID of the node
//	categories:[]text -> the category names of the node
//	visit_count:number -> the count of visits to the node in this run
//
// @context node
func (r *run) nodeContext(env envs.Environment) map[string]types.XValue {
	_, node, _ := r.PathLocation()
	visitCount := 0
	for _, s := range r.path {
		if s.NodeUUID() == node.UUID() {
			visitCount++
		}
	}

	var categories []types.XValue
	if node.Router() != nil {
		categories = make([]types.XValue, len(node.Router().Categories()))
		for i, c := range node.Router().Categories() {
			categories[i] = types.NewXText(c.Name())
		}
	}

	return map[string]types.XValue{
		"uuid":        types.NewXText(string(node.UUID())),
		"categories":  types.NewXArray(categories...),
		"visit_count": types.NewXNumberFromInt(visitCount),
	}
}

// EvaluateTemplate evaluates the given template in the context of this run
func (r *run) EvaluateTemplateValue(template string, log flows.EventCallback) (types.XValue, bool) {
	ctx := types.NewXObject(r.RootContext(r.session.MergedEnvironment()))

	value, warnings, err := r.session.Engine().Evaluator().TemplateValue(r.session.MergedEnvironment(), ctx, template)
	if err != nil {
		log(events.NewError(err.Error()))
	}
	for _, w := range warnings {
		log(events.NewWarning(w))
	}
	return value, err == nil
}

// EvaluateTemplateText evaluates the given template as text in the context of this run
func (r *run) EvaluateTemplateText(template string, escaping excellent.Escaping, truncate bool, log flows.EventCallback) (string, bool) {
	ctx := types.NewXObject(r.RootContext(r.session.MergedEnvironment()))

	value, warnings, err := r.session.Engine().Evaluator().Template(r.session.MergedEnvironment(), ctx, template, escaping)
	if err != nil {
		log(events.NewError(err.Error()))
	}
	for _, w := range warnings {
		log(events.NewWarning(w))
	}
	if truncate {
		value = stringsx.TruncateEllipsis(value, r.Session().Engine().Options().MaxTemplateChars)
	}
	return value, err == nil
}

// EvaluateTemplate is a convenience function for evaluating as text with truncating but no escaping
func (r *run) EvaluateTemplate(template string, log flows.EventCallback) (string, bool) {
	return r.EvaluateTemplateText(template, nil, true, log)
}

// get the ordered list of languages to be used for localization in this run
func (r *run) getLanguages() []i18n.Language {
	languages := make([]i18n.Language, 0, 3)

	// if contact has an allowed language, it takes priority
	contactLanguage := r.session.MergedEnvironment().DefaultLanguage()
	if contactLanguage != i18n.NilLanguage {
		languages = append(languages, contactLanguage)
	}

	// next we include the default language if it's different to the contact language
	defaultLanguage := r.session.Environment().DefaultLanguage()
	if defaultLanguage != i18n.NilLanguage && defaultLanguage != contactLanguage {
		languages = append(languages, defaultLanguage)
	}

	// finally we include the flow native language if it isn't an allowed language - because it's the only
	// one guaranteed to have translations
	return append(languages, r.flow.Language())
}

// GetText is a convenience version of GetTextArray for a single text values
func (r *run) GetText(uuid uuids.UUID, key string, native string) (string, i18n.Language) {
	textArray, lang := r.getText(uuid, key, []string{native}, nil)
	return textArray[0], lang
}

// GetTextArray returns the localized value for the given flow definition value
func (r *run) GetTextArray(uuid uuids.UUID, key string, native []string, languages []i18n.Language) ([]string, i18n.Language) {
	return r.getText(uuid, key, native, languages)
}

func (r *run) getText(uuid uuids.UUID, key string, native []string, languages []i18n.Language) ([]string, i18n.Language) {
	nativeLang := r.Flow().Language()

	// if a preferred language list wasn't provided, default to the run preferred languages
	if languages == nil {
		languages = r.getLanguages()
	}

	for _, lang := range languages {
		if lang == r.Flow().Language() {
			return native, nativeLang
		}

		translated := r.Flow().Localization().GetItemTranslation(lang, uuid, key)
		if len(translated) == 0 {
			continue
		}

		return translated, lang
	}

	return native, nativeLang
}

func (r *run) Snapshot() flows.RunSummary {
	return newRunSummaryFromRun(r)
}

var _ flows.RunSummary = (*run)(nil)

//------------------------------------------------------------------------------------------
// JSON Encoding / Decoding
//------------------------------------------------------------------------------------------

type runEnvelope struct {
	UUID       flows.RunUUID         `json:"uuid" validate:"required,uuid4"`
	Flow       *assets.FlowReference `json:"flow" validate:"required"`
	Path       []*step               `json:"path" validate:"dive"`
	Events     []json.RawMessage     `json:"events,omitempty"`
	Locals     *flows.Locals         `json:"locals,omitzero"`
	Results    flows.Results         `json:"results,omitempty" validate:"omitempty,dive"`
	Status     flows.RunStatus       `json:"status" validate:"required"`
	ParentUUID flows.RunUUID         `json:"parent_uuid,omitempty" validate:"omitempty,uuid4"`

	CreatedOn  time.Time  `json:"created_on" validate:"required"`
	ModifiedOn time.Time  `json:"modified_on" validate:"required"`
	ExitedOn   *time.Time `json:"exited_on"`
}

// ReadRun decodes a run from the passed in JSON. Parent run UUID is returned separately as the
// run in question might be loaded yet from the session.
func ReadRun(session flows.Session, data json.RawMessage, missing assets.MissingCallback) (flows.Run, error) {
	e := &runEnvelope{}
	var err error

	if err = utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, fmt.Errorf("unable to read run: %w", err)
	}

	r := &run{
		session:    session,
		uuid:       e.UUID,
		flowRef:    e.Flow,
		status:     e.Status,
		createdOn:  e.CreatedOn,
		modifiedOn: e.ModifiedOn,
		exitedOn:   e.ExitedOn,
	}

	// lookup actual flow
	if r.flow, err = session.Assets().Flows().Get(e.Flow.UUID); err != nil {
		missing(e.Flow, err)
	}

	// lookup parent run
	if e.ParentUUID != "" {
		if r.parent, err = session.GetRun(e.ParentUUID); err != nil {
			return nil, err
		}
	}

	if e.Locals != nil {
		r.locals = e.Locals
	} else {
		r.locals = flows.NewLocals()
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
			return nil, fmt.Errorf("unable to read event: %w", err)
		}
	}

	// create context
	r.webhook = lastWebhookSavedAsExtra(r)
	r.legacyExtra = newLegacyExtra(r)

	return r, nil
}

// MarshalJSON marshals this flow run into JSON
func (r *run) MarshalJSON() ([]byte, error) {
	var err error

	e := &runEnvelope{
		UUID:       r.uuid,
		Flow:       r.flowRef,
		Locals:     r.locals,
		Results:    r.results,
		Status:     r.status,
		CreatedOn:  r.createdOn,
		ModifiedOn: r.modifiedOn,
		ExitedOn:   r.exitedOn,
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
		if e.Events[i], err = jsonx.Marshal(r.events[i]); err != nil {
			return nil, fmt.Errorf("unable to marshal event[type=%s]: %w", r.events[i].Type(), err)
		}
	}

	return jsonx.Marshal(e)
}
