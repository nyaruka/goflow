package engine

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

const (
	persistWebhookBytesLimit = 10_000 // max bytes of webhook response we will persist in a run
)

type run struct {
	uuid    flows.RunUUID
	session *session

	flow    flows.Flow
	flowRef *assets.FlowReference

	parent   *run
	locals   *flows.Locals
	results  flows.Results
	path     Path
	hadInput bool
	status   flows.RunStatus
	webhook  *flows.WebhookCall

	createdOn  time.Time
	modifiedOn time.Time
	exitedOn   *time.Time

	legacyExtra     *legacyExtra
	legacyWaitCount int
}

func newRun(session *session, flow flows.Flow, parent *run) *run {
	now := dates.Now()
	r := &run{
		uuid:       flows.NewRunUUID(),
		session:    session,
		flow:       flow,
		flowRef:    flow.Reference(true),
		parent:     parent,
		locals:     flows.NewLocals(),
		results:    flows.NewResults(),
		status:     flows.RunStatusActive,
		createdOn:  now,
		modifiedOn: now,
	}

	r.legacyExtra = newLegacyExtra(r)

	return r
}

func (r *run) UUID() flows.RunUUID    { return r.uuid }
func (r *run) Session() flows.Session { return r.session }

func (r *run) Flow() flows.Flow                     { return r.flow }
func (r *run) FlowReference() *assets.FlowReference { return r.flowRef }
func (r *run) Contact() *flows.Contact              { return r.session.Contact() }
func (r *run) HadInput() bool                       { return r.hadInput }

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
func (r *run) setStatus(status flows.RunStatus) {
	r.status = status
	r.modifiedOn = dates.Now()
}

func (r *run) Webhook() *flows.WebhookCall { return r.webhook }
func (r *run) SetWebhook(call *flows.WebhookCall) {
	r.webhook = call
}

// Parent returns either the same session parent or if this session was triggered from a trigger_flow action
// in another session, that run
func (r *run) Parent() flows.RunSummary {
	if r.parent == nil {
		return r.session.ParentRun()
	}
	return r.parent
}

func (r *run) Ancestors() []flows.Run {
	ancestors := make([]flows.Run, 0)
	if r.parent != nil {
		pr := r.parent
		ancestors = append(ancestors, pr)

		for {
			if pr.parent != nil {
				pr = pr.parent
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

	if event.Type() == events.TypeMsgReceived {
		r.hadInput = true
	}

	r.modifiedOn = dates.Now()
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
	if r.session.findCurrentChild(r) != nil {
		child = newRelatedRunContext(r.session.findCurrentChild(r))
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
	UUID       flows.RunUUID         `json:"uuid"                  validate:"required,uuid"`
	Flow       *assets.FlowReference `json:"flow"                  validate:"required"`
	Path       []*step               `json:"path"                  validate:"dive"`
	Locals     *flows.Locals         `json:"locals,omitzero"`
	Results    flows.Results         `json:"results,omitempty"     validate:"omitempty,dive"`
	Status     flows.RunStatus       `json:"status"                validate:"required"`
	HadInput   bool                  `json:"had_input,omitzero"`
	ParentUUID flows.RunUUID         `json:"parent_uuid,omitempty" validate:"omitempty,uuid"`
	Webhook    *flows.WebhookCall    `json:"webhook,omitempty"`

	CreatedOn  time.Time  `json:"created_on"  validate:"required"`
	ModifiedOn time.Time  `json:"modified_on" validate:"required"`
	ExitedOn   *time.Time `json:"exited_on"`

	// older runs will have events which we can use to infer newer fields
	LegacyEvents []json.RawMessage `json:"events,omitempty"`
}

// decodes a run from the passed in JSON. Parent run UUID is returned separately as the
// run in question might be loaded yet from the session.
func readRun(s *session, data []byte, missing assets.MissingCallback) (*run, error) {
	e := &runEnvelope{}
	var err error

	if err = utils.UnmarshalAndValidate(data, e); err != nil {
		return nil, fmt.Errorf("unable to read run: %w", err)
	}

	r := &run{
		session:    s,
		uuid:       e.UUID,
		flowRef:    e.Flow,
		status:     e.Status,
		hadInput:   e.HadInput,
		webhook:    e.Webhook,
		createdOn:  e.CreatedOn,
		modifiedOn: e.ModifiedOn,
		exitedOn:   e.ExitedOn,
	}

	// lookup actual flow
	if r.flow, err = s.Assets().Flows().Get(e.Flow.UUID); err != nil {
		missing(e.Flow, err)
	}

	// lookup parent run
	if e.ParentUUID != "" {
		if r.parent, err = s.getRun(e.ParentUUID); err != nil {
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

	// older runs will have events which we can use to infer newer fields
	resultEventsWithExtra := make(map[flows.StepUUID]*events.RunResultChanged, 5)
	var lastWebhookEvent *events.WebhookCalled
	for i := range e.LegacyEvents {
		e, err := events.Read(e.LegacyEvents[i])
		if err != nil {
			return nil, fmt.Errorf("unable to read event %d: %w", i, err)
		}

		switch typed := e.(type) {
		case *events.MsgReceived:
			r.hadInput = true
		case *events.RunResultChanged:
			if typed.Extra != nil {
				resultEventsWithExtra[typed.StepUUID()] = typed
			}
		case *events.WebhookCalled:
			lastWebhookEvent = typed
		case *events.MsgWait, *events.DialWait:
			r.legacyWaitCount++
		}
	}

	// if we have a webhook event, look for a result event with extra data at the same step
	if lastWebhookEvent != nil {
		if resultEvent := resultEventsWithExtra[lastWebhookEvent.StepUUID()]; resultEvent != nil {
			r.webhook = &flows.WebhookCall{
				ResponseStatus: lastWebhookEvent.StatusCode,
				ResponseJSON:   resultEvent.Extra,
			}
		}
	}

	r.legacyExtra = newLegacyExtra(r)

	return r, nil
}

// MarshalJSON marshals this flow run into JSON
func (r *run) MarshalJSON() ([]byte, error) {
	e := &runEnvelope{
		UUID:       r.uuid,
		Flow:       r.flowRef,
		Locals:     r.locals,
		Results:    r.results,
		Status:     r.status,
		HadInput:   r.hadInput,
		CreatedOn:  r.createdOn,
		ModifiedOn: r.modifiedOn,
		ExitedOn:   r.exitedOn,
	}

	if r.webhook != nil && len(r.webhook.ResponseJSON) <= persistWebhookBytesLimit {
		e.Webhook = r.webhook
	}

	if r.parent != nil {
		e.ParentUUID = r.parent.UUID()
	}

	e.Path = make([]*step, len(r.path))
	for i, s := range r.path {
		e.Path[i] = s.(*step)
	}

	return jsonx.Marshal(e)
}
