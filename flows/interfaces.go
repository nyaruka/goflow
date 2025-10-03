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

const (
	MaxNodesPerFlow        = 1000 // max number of nodes in a flow
	MaxActionsPerNode      = 100  // max number of actions in a node
	MaxExitsPerNode        = 100  // max number of exits in a node
	MaxCategoriesPerRouter = 100  // max number of categories a router can have
	MaxCasesPerRouter      = 100  // max number of categories a switch router can have
	MaxArgumentsPerCase    = 100  // max number of test arguments a switch router case can have
)

// NodeUUID is a UUID of a flow node
type NodeUUID uuids.UUID

// CategoryUUID is the UUID of a node category
type CategoryUUID uuids.UUID

// ExitUUID is the UUID of a node exit
type ExitUUID uuids.UUID

// ActionUUID is the UUID of an action
type ActionUUID uuids.UUID

// StepUUID is the UUID of a run step
type StepUUID uuids.UUID

// InputUUID is the UUID of an input
type InputUUID uuids.UUID

// FlowAssets provides access to flow assets
type FlowAssets interface {
	Get(assets.FlowUUID) (Flow, error)
	FindByName(string) (Flow, error)
}

// SessionAssets is the assets available to a session
type SessionAssets interface {
	contactql.Resolver

	Source() assets.Source

	Campaigns() *CampaignAssets
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

	Inspect(sa SessionAssets) *Info
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
	Inspect(func(Action, Router, assets.Reference), func(Action, Router, string), func(Action, Router, *ResultInfo))
	EnumerateTemplates(Localization, func(Action, Router, i18n.Language, string))
	EnumerateLocalizables(func(uuids.UUID, string, []string, func([]string)))
}

// Action is an action within a flow node
type Action interface {
	utils.Typed
	Localizable
	FlowTypeRestricted

	UUID() ActionUUID
	Execute(context.Context, Run, Step, EventLogger) error
	Validate() error
	Inspect(func(assets.Reference), func(string), func(*ResultInfo))
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

	AllowTimeout() bool
	Route(Run, Step, EventLogger) (ExitUUID, string, error)
	RouteTimeout(Run, Step, EventLogger) (ExitUUID, error)

	Validate(Flow, []Exit) error
	Inspect(func(*ResultInfo), func(assets.Reference))
	EnumerateTemplates(Localization, func(i18n.Language, string))
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

	Begin(Run, EventLogger) bool
	Accepts(Resume) bool
}

// Hint tells the caller what type of input the flow is expecting
type Hint interface {
	utils.Typed
}

// Localization provide a way to get the translations for a specific language
type Localization interface {
	Validate() error
	GetItemTranslation(i18n.Language, uuids.UUID, string) []string
	SetItemTranslation(i18n.Language, uuids.UUID, string, []string)
	Languages() []i18n.Language
}

// Trigger represents something which can initiate a session with the flow engine
type Trigger interface {
	utils.Typed
	Contextable

	Event() Event
	Flow() *assets.FlowReference
	Batch() bool
	Params() *types.XObject
	History() *SessionHistory
	TriggeredOn() time.Time

	Input(SessionAssets) Input
}

// TriggerWithRun is special case of trigger that provides a parent run to the session
type TriggerWithRun interface {
	Trigger

	RunSummary() []byte
}

// Resume represents something which can resume a session with the flow engine
type Resume interface {
	utils.Typed
	Contextable

	Event() Event
	ResumedOn() time.Time

	Input(SessionAssets) Input
}

// Modifier is something which can modify a contact
type Modifier interface {
	utils.Typed

	Apply(Engine, envs.Environment, SessionAssets, *Contact, EventLogger) (bool, error)
}

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

type CheckSendableCallback func(SessionAssets, *Contact, *MsgContent) (UnsendableReason, error)

type EngineOptions struct {
	MaxStepsPerSprint    int
	MaxSprintsPerSession int
	MaxTemplateChars     int
	MaxFieldChars        int
	MaxResultChars       int
	LLMPrompts           map[string]*template.Template
	CheckSendable        CheckSendableCallback
}

// Engine provides callers with session starting and resuming
type Engine interface {
	NewSession(context.Context, SessionAssets, envs.Environment, *Contact, Trigger, *Call) (Session, Sprint, error)
	ReadSession(SessionAssets, []byte, envs.Environment, *Contact, *Call, assets.MissingCallback) (Session, error)

	Evaluator() *excellent.Evaluator
	Services() Services
	Options() *EngineOptions
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
