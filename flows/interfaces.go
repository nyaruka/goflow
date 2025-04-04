package flows

import (
	"context"
	"encoding/json"
	"text/template"
	"time"

	"github.com/nyaruka/gocommon/i18n"
	"github.com/nyaruka/gocommon/uuids"
	"github.com/nyaruka/goflow/assets"
	"github.com/nyaruka/goflow/contactql"
	"github.com/nyaruka/goflow/envs"
	"github.com/nyaruka/goflow/excellent"
	"github.com/nyaruka/goflow/excellent/types"
	"github.com/nyaruka/goflow/utils"
)

// NodeUUID is a UUID of a flow node
type NodeUUID uuids.UUID

// CategoryUUID is the UUID of a node category
type CategoryUUID uuids.UUID

// ExitUUID is the UUID of a node exit
type ExitUUID uuids.UUID

// ActionUUID is the UUID of an action
type ActionUUID uuids.UUID

// ContactID is the ID of a contact
type ContactID int64

// ContactUUID is the UUID of a contact
type ContactUUID uuids.UUID

// RunUUID is the UUID of a flow run
type RunUUID uuids.UUID

// StepUUID is the UUID of a run step
type StepUUID uuids.UUID

// InputUUID is the UUID of an input
type InputUUID uuids.UUID

// SessionUUID is the UUID of a session
type SessionUUID uuids.UUID

// SprintUUID is the UUID of a sprint
type SprintUUID uuids.UUID

// MsgID is the ID of a message
type MsgID int64

// NilMsgID is our constant for nil message ids
const NilMsgID = MsgID(0)

// MsgUUID is the UUID of a message
type MsgUUID uuids.UUID

// SessionStatus represents the current status of the engine session
type SessionStatus string

const (
	// SessionStatusActive represents a session that is still active
	SessionStatusActive SessionStatus = "active"

	// SessionStatusCompleted represents a session that has run to completion
	SessionStatusCompleted SessionStatus = "completed"

	// SessionStatusWaiting represents a session which is waiting for something from the caller
	SessionStatusWaiting SessionStatus = "waiting"

	// SessionStatusFailed represents a session that encountered an unrecoverable error
	SessionStatusFailed SessionStatus = "failed"
)

// RunStatus represents the current status of the flow run
type RunStatus string

const (
	// RunStatusActive represents a run that is still active
	RunStatusActive RunStatus = "active"

	// RunStatusCompleted represents a run that has run to completion
	RunStatusCompleted RunStatus = "completed"

	// RunStatusWaiting represents a run which is waiting for something from the caller
	RunStatusWaiting RunStatus = "waiting"

	// RunStatusFailed represents a run that encountered an unrecoverable error
	RunStatusFailed RunStatus = "failed"

	// RunStatusExpired represents a run that expired due to inactivity
	RunStatusExpired RunStatus = "expired"
)

// FlowAssets provides access to flow assets
type FlowAssets interface {
	Get(assets.FlowUUID) (Flow, error)
	FindByName(string) (Flow, error)
}

// SessionAssets is the assets available to a session
type SessionAssets interface {
	contactql.Resolver

	Source() assets.Source

	Channels() *ChannelAssets
	Classifiers() *ClassifierAssets
	Fields() *FieldAssets
	Flows() FlowAssets
	Globals() *GlobalAssets
	Groups() *GroupAssets
	Labels() *LabelAssets
	LLMs() *LLMAssets
	Locations() *LocationAssets
	OptIns() *OptInAssets
	Resthooks() *ResthookAssets
	Templates() *TemplateAssets
	Topics() *TopicAssets
	Users() *UserAssets
}

// Localizable is anything in the flow definition which can be localized and therefore needs a UUID
type Localizable interface {
	LocalizationUUID() uuids.UUID
}

type TemplateEnumerator interface {
	EnumerateTemplates(Localization, func(i18n.Language, string))
}

// Flow describes the ordered logic of actions and routers
type Flow interface {
	Contextable

	// spec properties
	UUID() assets.FlowUUID
	Name() string
	Revision() int
	Language() i18n.Language
	Type() FlowType
	ExpireAfter() time.Duration
	Localization() Localization
	UI() json.RawMessage
	Nodes() []Node
	GetNode(uuid NodeUUID) Node

	Asset() assets.Flow
	Reference(bool) *assets.FlowReference

	Inspect(sa SessionAssets) *Inspection
	ExtractTemplates() []string
	ExtractLocalizables() []string
	ChangeLanguage(i18n.Language) (Flow, error)
}

// Node is a single node in a flow
type Node interface {
	UUID() NodeUUID
	Actions() []Action
	Router() Router
	Exits() []Exit

	Validate(Flow, map[uuids.UUID]bool) error

	EnumerateTemplates(Localization, func(Action, Router, i18n.Language, string))
	EnumerateDependencies(Localization, func(Action, Router, i18n.Language, assets.Reference))
	EnumerateResults(func(Action, Router, *ResultInfo))
	EnumerateLocalizables(func(uuids.UUID, string, []string, func([]string)))
}

// Action is an action within a flow node
type Action interface {
	utils.Typed
	Localizable
	FlowTypeRestricted

	UUID() ActionUUID
	Execute(context.Context, Run, Step, ModifierCallback, EventCallback) error
	Validate() error
}

// Category is how routers map results to exits
type Category interface {
	Localizable

	UUID() CategoryUUID
	Name() string
	ExitUUID() ExitUUID
}

// Router is a router on a note which can pick an exit
type Router interface {
	utils.Typed

	Wait() Wait
	Categories() []Category
	ResultName() string

	Validate(Flow, []Exit) error
	AllowTimeout() bool
	Route(Run, Step, EventCallback) (ExitUUID, string, error)
	RouteTimeout(Run, Step, EventCallback) (ExitUUID, error)

	EnumerateTemplates(Localization, func(i18n.Language, string))
	EnumerateDependencies(Localization, func(i18n.Language, assets.Reference))
	EnumerateResults(func(*ResultInfo))
	EnumerateLocalizables(func(uuids.UUID, string, []string, func([]string)))
}

// Exit is a route out of a node and optionally to another node
type Exit interface {
	UUID() ExitUUID
	DestinationUUID() NodeUUID
}

// Timeout is a way to skip a wait after X amount of time
type Timeout interface {
	Seconds() int
	CategoryUUID() CategoryUUID
}

// Wait tells the engine that the session requires input from the user
type Wait interface {
	utils.Typed
	FlowTypeRestricted

	Timeout() Timeout

	Begin(Run, EventCallback) bool
	Accepts(Resume) bool
}

// Hint tells the caller what type of input the flow is expecting
type Hint interface {
	utils.Typed
}

// Localization provide a way to get the translations for a specific language
type Localization interface {
	GetItemTranslation(i18n.Language, uuids.UUID, string) []string
	SetItemTranslation(i18n.Language, uuids.UUID, string, []string)
	Languages() []i18n.Language
}

// Trigger represents something which can initiate a session with the flow engine
type Trigger interface {
	utils.Typed
	Contextable

	Initialize(Session, EventCallback) error
	InitializeRun(Run, EventCallback) error

	Environment() envs.Environment
	Flow() *assets.FlowReference
	Contact() *Contact
	Call() *Call
	Batch() bool
	Params() *types.XObject
	History() *SessionHistory
	TriggeredOn() time.Time
}

// TriggerWithRun is special case of trigger that provides a parent run to the session
type TriggerWithRun interface {
	Trigger

	RunSummary() json.RawMessage
}

// Resume represents something which can resume a session with the flow engine
type Resume interface {
	utils.Typed
	Contextable

	Apply(Run, EventCallback)

	Environment() envs.Environment
	Contact() *Contact
	ResumedOn() time.Time
}

// Modifier is something which can modify a contact
type Modifier interface {
	utils.Typed

	Apply(Engine, envs.Environment, SessionAssets, *Contact, EventCallback) bool
}

// ModifierCallback is a callback invoked when a modifier has been generated
type ModifierCallback func(Modifier)

// Event describes a state change
type Event interface {
	utils.Typed

	CreatedOn() time.Time
	StepUUID() StepUUID
	SetStepUUID(StepUUID)
}

// EventCallback is a callback invoked when an event has been generated
type EventCallback func(Event)

// Input describes input from the contact and currently we only support one type of input: `msg`
type Input interface {
	utils.Typed
	Contextable

	UUID() InputUUID
	CreatedOn() time.Time
	Channel() *Channel
}

// Step is a single step in the path thru a flow
type Step interface {
	Contextable

	UUID() StepUUID
	NodeUUID() NodeUUID
	ExitUUID() ExitUUID
	ArrivedOn() time.Time

	Leave(ExitUUID)
}

type EngineOptions struct {
	MaxStepsPerSprint    int
	MaxResumesPerSession int
	MaxTemplateChars     int
	MaxFieldChars        int
	MaxResultChars       int
	LLMPrompts           map[string]*template.Template
}

// Engine provides callers with session starting and resuming
type Engine interface {
	NewSession(context.Context, SessionAssets, Trigger) (Session, Sprint, error)
	ReadSession(SessionAssets, json.RawMessage, assets.MissingCallback) (Session, error)

	Evaluator() *excellent.Evaluator
	Services() Services
	Options() *EngineOptions
}

// Segment is a movement on the flow graph from an exit to another node
type Segment interface {
	Flow() Flow
	Node() Node
	Exit() Exit
	Operand() string
	Destination() Node
	Time() time.Time
}

// Sprint is an interaction with the engine - i.e. a start or resume of a session
type Sprint interface {
	UUID() SprintUUID
	Modifiers() []Modifier
	Events() []Event
	Segments() []Segment
}

// Session represents the session of a flow run which may contain many runs
type Session interface {
	Assets() SessionAssets

	UUID() SessionUUID
	Type() FlowType
	SetType(FlowType)

	Environment() envs.Environment
	SetEnvironment(envs.Environment)
	MergedEnvironment() envs.Environment

	Contact() *Contact
	SetContact(*Contact)

	Input() Input
	SetInput(Input)

	Status() SessionStatus
	Trigger() Trigger
	CurrentResume() Resume
	BatchStart() bool
	PushFlow(Flow, Run, bool)

	Resume(context.Context, Resume) (Sprint, error)
	Runs() []Run
	GetRun(RunUUID) (Run, error)
	FindStep(uuid StepUUID) (Run, Step)
	GetCurrentChild(Run) Run
	ParentRun() RunSummary
	CurrentContext() *types.XObject
	History() *SessionHistory

	Engine() Engine
}

// RunSummary represents the minimum information available about all runs (current or related) and is the
// representation of runs made accessible to router tests.
type RunSummary interface {
	UUID() RunUUID
	Contact() *Contact
	Flow() Flow
	Status() RunStatus
	Results() Results
}

// Run is a single contact's journey through a flow. It records the path they have taken,
// and the results that have been collected.
type Run interface {
	Contextable
	RunSummary
	FlowReference() *assets.FlowReference

	Session() Session
	Locals() *Locals
	SetResult(*Result) (*Result, bool)
	SetStatus(RunStatus)
	Webhook() *WebhookCall
	SetWebhook(*WebhookCall)

	CreateStep(Node) Step
	Path() []Step
	PathLocation() (Step, Node, error)

	LogEvent(Step, Event)
	Events() []Event
	ReceivedInput() bool

	EvaluateTemplateValue(string, EventCallback) (types.XValue, bool)
	EvaluateTemplateText(string, excellent.Escaping, bool, EventCallback) (string, bool)
	EvaluateTemplate(string, EventCallback) (string, bool)
	RootContext(envs.Environment) map[string]types.XValue

	GetText(uuids.UUID, string, string) (string, i18n.Language)
	GetTextArray(uuids.UUID, string, []string, []i18n.Language) ([]string, i18n.Language)

	Snapshot() RunSummary
	Parent() RunSummary
	ParentInSession() Run
	Ancestors() []Run

	CreatedOn() time.Time
	ModifiedOn() time.Time
	ExitedOn() *time.Time
	Exit(RunStatus)
}

// LegacyExtraContributor is something which contributes results for constructing @legacy_extra
type LegacyExtraContributor interface {
	LegacyExtra() Results
}

type Dependency interface {
	Reference() assets.Reference
	Type() string
	Missing() bool
}

// Issue is a problem found during flow inspection
type Issue interface {
	utils.Typed

	NodeUUID() NodeUUID
	ActionUUID() ActionUUID
	Language() i18n.Language
	Description() string
}
